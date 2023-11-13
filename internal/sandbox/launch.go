package sandbox

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/msqtt/sb-judger/internal/sandbox/fs"
)

// LaunchEntry is the entry of executing compiled program.
func LaunchEntry() (err error) {
	// stage1: read args from pipe
	mntPath, binaryName, err := readArgsFromPipe()	
	if err != nil {
		return
	}
	// stage2: mount something to mark
	err = setUpMount(mntPath)
	if err != nil {
		return
	}
	// stage3: send a start signal to parent process then wait to exec to binary program.
	binaryPath := filepath.Join("/tmp", binaryName)

	// telling parent process im ready.
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGUSR1)
	syscall.Kill(os.Getppid(), syscall.SIGUSR1)
	// wait to start.
	<-sign
	if err := syscall.Exec(binaryPath, nil, os.Environ()); err != nil {
		return fmt.Errorf("cannot exec command in the container: %w", err)
	}

	return nil
}

// readArgsFromPipe reads args from pipe passing by parent process.
func readArgsFromPipe() (mntPath, binary string, err error) {
	var tmp []byte
	if tmp, err = ioutil.ReadAll(os.NewFile(uintptr(3), "pipe")); err != nil {
		err = fmt.Errorf("parent process cannot read pipe: [%w]", err)
		return
	} else {
		s := strings.Split(string(tmp), " ")
		if len(s) < 2 {
			err = errors.New("invalid launch args read from pipe")
			return
		}
		mntPath = s[0]
		binary = s[1]
		return
	}
}

// setUpMount pivots rootfs and mark /dev path.
func setUpMount(mnt string) (err error) {
	err = fs.ChrootMaskPath(mnt)
	if err != nil {
		return
	}
	// mark /dev
	if err = syscall.Mount("tmpfs", "/dev", "tmpfs",
			syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		err = fmt.Errorf("cannot mount tmpfs: %w", err)
		return
	}
	return
}
