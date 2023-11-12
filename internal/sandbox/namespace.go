package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const selfExe = "/proc/self/exe"

// maskFork executes self to clone a process and make a namespace for it.
// returns cmd of clone process,  writePipe file and an error.
func maskFork(args []string) (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create pipe: [%w]", err)
	}
	cmd := exec.Command(selfExe, args...)
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// put sub sub process to a pgroup (if have)
		Setpgid: true,
		// needed namespace
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
		// damn, not enough permittion to mount fs
		// syscall.CLONE_NEWUSER
	}
	return cmd, writePipe, nil
}
