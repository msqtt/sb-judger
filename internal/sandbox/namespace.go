package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

const exe = "./sandbox"

var NSTypeItoa = map[int]string{
	unix.CLONE_NEWCGROUP: "cgroup",
	unix.CLONE_NEWIPC:    "ipc",
	unix.CLONE_NEWNS:     "mnt",
	unix.CLONE_NEWNET:    "net",
	unix.CLONE_NEWPID:    "pid",
	unix.CLONE_NEWUSER:   "user",
	unix.CLONE_NEWUTS:    "uts",
	unix.CLONE_NEWTIME:   "time",
}

// maskFork executes self to clone a process and make a namespace for it.
// returns cmd of clone process,  writePipe file and an error.
func maskFork() (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create pipe: [%w]", err)
	}
	cmd := exec.Command(exe)
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// put sub sub process to a pgroup (if have)
		Setpgid: true,
		// needed namespace
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			// need to send message to ppid, so use setns to set instead
			// syscall.CLONE_NEWIPC |
			// damn, can not set this
			// syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}
	return cmd, writePipe, nil
}

func getNsFD(pid int, nsType int) (fd int, err error) {
	nsPath := filepath.Join("/proc", strconv.Itoa(pid), "ns", NSTypeItoa[nsType])
	if _, err = os.Stat(nsPath); os.IsNotExist(err) {
		err = fmt.Errorf("namespace file does not exist: %s", nsPath)
		return
	}
	fd, err = unix.Open(nsPath, unix.O_RDONLY, 0)
	return
}

func setNsFor(pid int, nsTypes ...int) error {
	for _, ns := range nsTypes {
		fd, err := getNsFD(pid, ns)
		if err != nil {
			return fmt.Errorf("cannot get %s fd for %d: %w", NSTypeItoa[ns], pid, err)
		}
		err = unix.Setns(fd, ns)
		if err != nil {
			return fmt.Errorf("cannot set %s namespace for %d: %w", NSTypeItoa[ns], pid, err)
		}
	}
	return nil
}
