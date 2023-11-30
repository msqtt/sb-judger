package sandbox

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	msgq "github.com/msqtt/sb-judger/internal/message"
	"github.com/msqtt/sb-judger/internal/sandbox/fs"
)

// LaunchEntry is the entry of executing compiled program.
func LaunchEntry() (err error) {
	// stage1: read args from pipe
	mntPath, runCmd, msqidStr, err := readArgsFromPipe()
	if err != nil {
		return
	}
  // setup args and lookup path.
  cmdArgs := strings.Split(runCmd, " ")
  runCmd, _ = exec.LookPath(cmdArgs[0])
  if runCmd == "" {
    runCmd = cmdArgs[0]
  }

	// stage2: mount to mask
  err = fs.ChrootMaskPath(mntPath)
	if err != nil {
		return
	}

  // chdir to workspace
  syscall.Chdir("/tmp")

  // set rootless for process
  syscall.Setgid(65534)
  syscall.Setgroups([]int{65534})
  syscall.Setuid(65534)

	// stage3: send a start message to parent process then wait to exec to binary program.
	// telling parent process im ready.
  d, _ := strconv.Atoi(msqidStr)
  msqid := uintptr(d)

  err2 := msgq.SndMsg(msqid, msgq.NewMsg(1, nil))
  if err2 != nil {
    return err2
  }

	// waiting for parent's message(rm queue) then launch code program.
  err = msgq.RcvMsg(msqid, msgq.NewMsg(2, nil))
  if err.Error() != "identifier removed" {
    return err
  }

  if err := syscall.Exec(runCmd, cmdArgs, nil); err != nil {
		return fmt.Errorf("cannot exec command in the container: %w", err)
	}

	return nil
}

// readArgsFromPipe reads args from pipe passing by parent process.
func readArgsFromPipe() (mntPath, runCmd, msqid string, err error) {
	var tmp []byte
  pipe := os.NewFile(uintptr(3), "pipe")
  defer pipe.Close()

	if tmp, err = ioutil.ReadAll(pipe); err != nil {
		err = fmt.Errorf("cannot read pipe: [%w]", err)
		return
	} else {
		s := strings.Split(string(tmp), "#")
		if len(s) < 2 {
			err = errors.New("invalid launch args read from pipe")
			return
		}
		mntPath = s[0]
		runCmd = s[1]
		msqid = s[2]
		return
	}
}
