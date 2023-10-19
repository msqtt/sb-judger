package smallbox

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const selfExe = "/proc/self/exe"

// maskFork executes self to clone a process and make a namespace for it.
// returns cmd of clone process,  writePipe file and an error.
func maskFork() (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create pipe: [%w]", err)
	}
	cmd := exec.Command(selfExe, strconv.Itoa(ArgLaunch))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
		// damn, not enough permittion to mount fs
		// syscall.CLONE_NEWUSER
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe, nil
}
