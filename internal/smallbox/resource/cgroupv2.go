package res

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const cgroupRoot = "/sys/fs/cgroup"

// CgroupV2 m.anage resource by writing to linux cgroup file.
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

// ReadState implements ResourceManager.
func (cg *CgroupV2) ReadState() (*RunState, error) {
	rs := new(RunState)
	byte, err := os.ReadFile(filepath.Join(cg.path, "cpu.stat"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile cpu.stat: %w", err)
	}
	cpuStatMap, err := makeStatMap(byte)
	rs.CpuTime = cpuStatMap["user_usec"] / 1000

	byte, err = os.ReadFile(filepath.Join(cg.path, "memory.peak"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.peak: %w", err)
	}
	integer, err := strconv.Atoi(string(byte))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.peak: %w", err)
	}
	rs.MemoryUsage = uint(integer >> 20)

	byte, err = os.ReadFile(filepath.Join(cg.path, "memory.events"))
	if err != nil {
		return nil, fmt.Errorf("cannot readfile memory.events: %w", err)
	}
	eventStat, err := makeStatMap(byte)
	rs.OOMKill = eventStat["oom_kill"]

	return rs, nil
}

func makeStatMap(byte []byte) (map[string]uint, error) {
	all := string(byte)
	lines := strings.Split(all, "\n")
	m := make(map[string]uint)
	for _, line := range lines {
		item := strings.Split(line, " ")
		integer, err := strconv.Atoi(item[1])
		if err != nil {
			return nil, err
		}
		m[item[0]] = uint(integer)
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
