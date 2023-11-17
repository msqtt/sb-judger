package sandbox

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/msqtt/sb-judger/internal/pkg/sleep"
	res "github.com/msqtt/sb-judger/internal/sandbox/resource"
)

const inteErrCode = 500

// Entry function is the intro of sandbox program.
func InitEntry(lang, hashName, mntPath string, mem, time uint32, cases []*pb_sb.Case) (
	*pb_sb.CollectOutput, error) {
	langConf, err := json.GetLangConfs("")
	if err != nil {
		return nil, fmt.Errorf("cannot get lang config: %w", err)
	}

	lc := langConf[lang]
	if lc == nil {
		return nil, errors.New("not supported language")
	}

	resManager, err := res.NewCgroupV2(hashName)
	if err != nil {
		return nil, fmt.Errorf("cannot new cgroupv2: %w", err)
	}
	defer resManager.Destroy()

	err = resManager.Config(&res.ResourceConfig{
		CpuLimit:    100000,
		MemoryLimit: mem,
		PidsLimit:   lc.Pids,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot config for resourceManager: %w", err)
	}

	collectOuts := make([]*pb_sb.Output, len(cases))

	// todo: using goroutine to fork per case sub process.
	for i, cas := range cases {
		// add sub cgroup for per case.
		cv, err := resManager.AddSubGroup(strconv.FormatUint(uint64(cas.GetCaseId()), 10))
		if err != nil {
			return nil, fmt.Errorf("cannot add sub group: %w", err)
		}

    outPath := filepath.Join("/tmp", lc.Out)
    // still, lazy codes 🤓
    runCmd := strings.Join(lc.Run, " ")
    runCmd = fmt.Sprintf(runCmd, outPath)
    runCmd = strings.Split(runCmd, "%!")[0]

		// passing hashname and input string
		bytes, code, err := execLaunchProcess(cv, time,
			mntPath, runCmd, cas.In)

    if code == inteErrCode {
      return nil, fmt.Errorf("cannot collect sub process: [%w, msg: %s]", err, string(bytes))
    }

		// reads resource state after launch process done.
		usage, err1 := cv.ReadState()
		if err1 != nil {
			return nil, fmt.Errorf("cannot read states from case %d process: %w", cas.GetCaseId(),
				err)
		}

		// check status
		status := getStatus(bytes, code, err, []byte(cas.GetOut()), usage,
			time, mem)
		if err != nil {
			return nil, fmt.Errorf("Case %d failed to check status: [%w]", cas.GetCaseId(), err)
		}
		// collection result
		collectOuts[i] = &pb_sb.Output{
			CaseId:      cas.CaseId,
			TimeUsage:   usage.CpuTime,
			MemoryUsage: usage.MemoryUsage,
			Status:      status,
			OutPut:      string(bytes),
		}
	}

	return &pb_sb.CollectOutput{CaseOuts: collectOuts}, nil
}

// execLaunchProcess exec a process then set the namespaces for it.
func execLaunchProcess(cv *res.CgroupV2, sec uint32, mntPath, runCmd, inputContent string) (
	[]byte, int, error) {
	// passing launch arg to execute launch entry parts.
	cmd, writePipe, err := maskFork()
	if err != nil {
		return nil, inteErrCode, err
	}

	// stage1: write args, input content to launch process
  // and binding stdout and stderr.
	mes := fmt.Sprintf("%s#%s", mntPath, runCmd)
	_, err = writePipe.WriteString(mes)
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot write pipe: %w", err)
	}
	writePipe.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot get stdin: %w", err)
	}
	_, err = stdin.Write([]byte(inputContent))
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot write input content to stdin: %w", err)
	}
	stdin.Close()

  // binding stdout and stderr
  var buf bytes.Buffer
  cmd.Stdout = &buf
  cmd.Stderr = &buf

	// stage2: start two goroutine, one for starting son process then wait to it exit,
	// the other for counting total seconds to kill son process.
	errChan := make(chan error)
	exitChan := make(chan int)
	dataChan := make(chan []byte)
	defer func() {
		close(errChan)
		close(exitChan)
		close(dataChan)
	}()

	// start launch process.
	err = cmd.Start()
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot start sub process: %w", err)
	}

	// waiting for launch process's start signal then add its pid to cgroups
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGUSR1)


	go func() {
		err := cmd.Wait()
		if err != nil {
			errChan <- err
			exitChan <- inteErrCode
		}
		dataChan <- buf.Bytes()
	}()

	go func() {
		// sleep total seconds
		for i := 0; i < int(sec*10); i++ {
			err = sleep.NanoSleep(100)
			if err != nil {
				errChan <- err
				exitChan <- inteErrCode
				dataChan <- []byte("cannot usleep")
				return
			}
		}
		errChan <- os.ErrDeadlineExceeded
		dataChan <- []byte("time limit exceeded")
	}()

	var data []byte
	var exitCode int

loop:
	// one round for signal coming, the other for checking that
	// process successfully exited or exceed time to be killed
	// or bad things happened.
	for i := 0; i < 2; i++ {
		select {

		// case: any errors happen
		case err = <-errChan:
			// kill pgroup
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			exitCode = <-exitChan
			data = <-dataChan
			break loop

		// case: process sucessfully exists
		case data = <-dataChan:
			break loop

		case <-sign:
			// start to adding pid to cgroups
			err = cv.Apply(cmd.Process.Pid)
			if err != nil {
				return nil, inteErrCode, fmt.Errorf("cannot apply cgroups for sub process: %w", err)
			}
			// after adding telling sub process is ok to run code program.
			syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)
		}
	}
	if exitCode != inteErrCode {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return data, exitCode, err
}

func getStatus(outCont []byte, code int, outErr error, ans []byte,
	usage *res.RunState, timeLimit, memLimit uint32) pb_sb.Status {
	// trim right space (is it needed ?)
	outCont = bytes.TrimRightFunc(outCont, unicode.IsSpace)
	ans = bytes.TrimRightFunc(ans, unicode.IsSpace)

	if usage.CpuTime > timeLimit {
		return pb_sb.Status_TLE
	}
	if usage.MemoryUsage > memLimit {
		return pb_sb.Status_MLE
	}
	if usage.OOMKill > 0 {
		return pb_sb.Status_MLE
	}

	switch code {
	case -1:
		if errors.Is(outErr, os.ErrDeadlineExceeded) {
			return pb_sb.Status_TLE
		}
		return pb_sb.Status_UE
	case 0:
		// diff between answer and user printout.
		if bytes.Compare(outCont, ans) == 0 {
			return pb_sb.Status_AC
		} else {
			return pb_sb.Status_WA
		}
	case 1, 2, 136, 139:
		return pb_sb.Status_RE
	default:
		return pb_sb.Status_UE
	}
}
