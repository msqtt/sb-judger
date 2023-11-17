package res

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const cgroupRoot = "/sys/fs/cgroup"

// init function write the cgroupv2 controllers to sub_controller.
// make sure that sub cgroup cant inherit the controllers.
// every error causes panic.
func init() {
	procPath := filepath.Join(cgroupRoot, "cgroup.procs")
	subTreePath := filepath.Join(cgroupRoot, "cgroup.subtree_control")
	ctrlPath := filepath.Join(cgroupRoot, "cgroup.controllers")
	tmpPath := filepath.Join(cgroupRoot, "tmp")
	tmpProcPath := filepath.Join(tmpPath, "cgroup.procs")

	// if subtrees already have controllers, skip.
	b, err := os.ReadFile(subTreePath)
	if err != nil {
		err = fmt.Errorf("cannot read system subtree controllers: %w", err)
		panic(err)
	}
	if bytes.Contains(b, []byte("cpu")) &&
		bytes.Contains(b, []byte("memory")) &&
		bytes.Contains(b, []byte("pids")) {
		return
	}

	// if cgroup have not needed controllers, panic.
	b1, err := os.ReadFile(ctrlPath)
	if err != nil {
		err = fmt.Errorf("cannot read system cgroup controllers: %w", err)
		panic(err)
	}
	if !bytes.Contains(b1, []byte("cpu")) ||
		!bytes.Contains(b1, []byte("memory")) ||
		!bytes.Contains(b1, []byte("pids")) {
		panic(errors.New("system cgroup controllers do not support"))
	}

	// inherit controllers
	b2, err := os.ReadFile(procPath)
	if err != nil {
		err = fmt.Errorf("cannot read system cgroup procs: %w", err)
		panic(err)
	}
	if len(b2) >= 0 {
		// root cgroup is busy
		// stage1: make a temp cgroup
		err = os.Mkdir(tmpPath, 0644)
		if err != nil {
			if !errors.Is(err, fs.ErrExist)  {
				err = fmt.Errorf("cannot make temp cgroup: %w", err)
				panic(err)
			}
		}
		// stage2: move root cgroup's pids to temp one.
    // if you start this program in a container and the first running process is itself
    // that might be no problem, other cases would cause panic deal to invalid arguments.
    // todo 💩
		err = os.WriteFile(tmpProcPath, b2, 0644)
		if err != nil {
			err = fmt.Errorf("cannot mv pid to temp cgroup: %w", err)
			panic(err)
		}
	}

	// root cgroup is not busy now
	// stage3: write sub-tree controllers
	err = addSubCtrls(cgroupRoot)
	if err != nil {
		err = fmt.Errorf("cannot add sub controllors: %w", err)
		panic(err)
	}
}

// CgroupV2 manage resource by writing to linux cgroup file.
// Make sure that the host is using cgroupv2 before using this struct.
type CgroupV2 struct {
	path string
	sons []*CgroupV2
}

var _ ResourceManager = (*CgroupV2)(nil)

// NewCgroupV2 mkdir in cgroup path by given name.
func NewCgroupV2(name string) (*CgroupV2, error) {
	path := filepath.Join(cgroupRoot, name)
	err := os.MkdirAll(path, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot new cgroup: [%w]", err)
	}
	return &CgroupV2{path: path}, nil
}

// Apply cgroup for pid.
func (cg *CgroupV2) Apply(pid int) error {
	filePath := filepath.Join(cg.path, "cgroup.procs")
	err := ioutil.WriteFile(filePath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("cannot apply pid %d: [%w]", pid, err)
	}
	return nil
}

// Config set the cgroup resource limit.
func (cg *CgroupV2) Config(config *ResourceConfig) error {
	// memory limit
	memorys := []string{
		"memory.max", "memory.high", "memory.swap.high", "memory.swap.max"}
	for _, mem := range memorys {
		filePath := filepath.Join(cg.path, mem)
		memLimit := int(config.MemoryLimit << 20)
		err := ioutil.WriteFile(filePath, []byte(strconv.Itoa(memLimit)), 0644)
		if err != nil {
			return fmt.Errorf("cannot write limit for %s: [%w]", mem, err)
		}
	}
	// cpu limit
	err := ioutil.WriteFile(filepath.Join(cg.path, "cpu.max"),
		[]byte(strconv.Itoa(int(config.CpuLimit))), 0644)
	if err != nil {
		return fmt.Errorf("cannot write limit for cpu.max: [%w]", err)
	}
	// pids limit
	err = ioutil.WriteFile(filepath.Join(cg.path, "pids.max"),
		[]byte(strconv.Itoa(int(config.PidsLimit))), 0644)
	if err != nil {
		return fmt.Errorf("cannot write limit for pids.max: [%w]", err)
	}
	return nil
}

// ReadState returns the usage data product by process running.
func (cg *CgroupV2) ReadState() (*RunState, error) {
	rs := new(RunState)
	byte, err := os.ReadFile(filepath.Join(cg.path, "cpu.stat"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile cpu.stat: %w", err)
	}
	cpuStatMap, err := makeStatMap(byte)
	rs.CpuTime = cpuStatMap["user_usec"]

	byte, err = os.ReadFile(filepath.Join(cg.path, "memory.peak"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.peak: %w", err)
	}
	memUsc, err := strconv.Atoi(strings.TrimSpace(string(byte)))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.peak: %w", err)
	}
	rs.MemoryUsage = uint32(memUsc >> 20)

	byte, err = os.ReadFile(filepath.Join(cg.path, "memory.events"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.events: %w", err)
	}
	eventStat, err := makeStatMap(byte)
	rs.OOMKill = eventStat["oom_kill"]

	return rs, nil
}

func makeStatMap(byte []byte) (map[string]uint32, error) {
	all := string(byte)
	lines := strings.Split(all, "\n")
	m := make(map[string]uint32)
	for _, line := range lines {
		item := strings.Split(line, " ")
		if len(item) < 2 {
			continue
		}
		integer, err := strconv.Atoi(item[1])
		if err != nil {
			return nil, err
		}
		m[item[0]] = uint32(integer)
	}
	return m, nil
}

// AddSubGroup will create a cgroup dir in current cgroup then add this to
func (cg *CgroupV2) AddSubGroup(name string) (*CgroupV2, error) {
	if cg.sons == nil {
		cg.sons = make([]*CgroupV2, 0)
		err := addSubCtrls(cg.path)
		if err != nil {
			return nil, fmt.Errorf("cannot set sub ctrl tree: %w", err)
		}
	}

	path := filepath.Join(cg.path, name)
	err := os.Mkdir(path, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot add subgroup: [%w]", err)
	}
	subCg := &CgroupV2{path: path}
	cg.sons = append(cg.sons, subCg)

	return subCg, nil
}

func addSubCtrls(path string) error {
	return os.WriteFile(filepath.Join(path, "cgroup.subtree_control"), []byte("+cpu +memory +pids"), 0644)
}

// Destroy remove current cgroup path, include its subcgroups.
// Do not use the CgroupV2 struct after executing Destroy or return an error.
func (cg *CgroupV2) Destroy() error {
	var err error
	if cg.sons != nil && len(cg.sons) > 0 {
		for i := range cg.sons {
			err1 := cg.sons[i].Destroy()
			if err1 != nil {
				err = fmt.Errorf("%w\n%s", err1, err)
			}
			cg.sons[i] = nil
		}
		if err != nil {
			return fmt.Errorf("cannot destroy cgroups: [%w]", err)
		}
	}

	err = os.RemoveAll(cg.path)
	if err != nil {
		return fmt.Errorf("cannot destroy %s: [%w]", cg.path, err)
	}
	return nil
}

// GetSons returns all subcgroups of current cgroup.
func (cg *CgroupV2) GetSons() []*CgroupV2 {
	return cg.sons
}
