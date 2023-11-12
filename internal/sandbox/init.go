package sandbox

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unicode"

	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/msqtt/sb-judger/internal/pkg/sleep"
	res "github.com/msqtt/sb-judger/internal/sandbox/resource"
	"google.golang.org/protobuf/proto"
)

const (
	ArgInit = iota
	ArgLaunch
)

const inteErrCode = 500

// Entry function is the intro of sandbox program.
func InitEntry() error {
	input, err := readInputFromPipe()
	if err != nil {
		return err
	}

	langConf, err := json.GetLangConfs("")
	if err != nil {
		return fmt.Errorf("cannot get lang config: %w", err)
	}

	lc := langConf[input.GetLang().String()]
	if lc == nil {
		return errors.New("not supported language")
	}

	resManager, err := res.NewCgroupV2(input.GetHashName())
	if err != nil {
		return fmt.Errorf("cannot new cgroupv2: %w", err)
	}
	defer resManager.Destroy()

	err = resManager.Config(&res.ResourceConfig{
		CpuLimit:    100000,
		MemoryLimit: input.GetMemory(),
		PidsLimit:   lc.Pids,
	})
	if err != nil {
		return fmt.Errorf("cannot config for resourceManager: %w", err)
	}

	collectOuts := make([]*pb_sb.Output, len(input.GetCases()))

	// todo: using goroutine to fork per case sub process.
	for i, cas := range input.GetCases() {
		// add sub cgroup for per case.
		cv, err := resManager.AddSubGroup(strconv.FormatUint(uint64(cas.GetCaseId()), 10))
		if err != nil {
			return fmt.Errorf("cannot add sub group: %w", err)
		}

		// passing hashname and input string
		bytes, code, err := forkParentProcess(cv, input.GetTime(),
			input.GetMntPath(), lc.Out, cas.In)

		// reads resource state after launch process done.
		usage, err := cv.ReadState()
		if err != nil {
			return fmt.Errorf("cannot read states from case %d process: %w", cas.GetCaseId(), err)
		}

		// check status
		status := getStatus(bytes, code, err, []byte(cas.GetOut()), usage,
			input.GetTime(), input.GetMemory())
		if err != nil {
			return fmt.Errorf("Case %d failed to check status: [%w]", cas.GetCaseId(), err)
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

	ret := &pb_sb.CollectOutput{CaseOuts: collectOuts}
	b, err := proto.Marshal(ret)
	if err != nil {
		return fmt.Errorf("cannot marshal collect output: %w", err)
	}

	for n, err := os.Stdout.Write(b); err != nil; n, err = os.Stdout.Write(b) {
		b = b[n:]
	}

	return nil
}

// forkParentProcess forks a launch process then set the namespace for it.
func forkParentProcess(cv *res.CgroupV2, sec uint32, mntPath, binaryName, inputContent string) (
	[]byte, int, error) {
	// passing launch arg to execute launch entry parts.
	cmd, writePipe, err := maskFork([]string{strconv.Itoa(ArgLaunch)})
	if err != nil {
		return nil, inteErrCode, err
	}

	// stage1: write infomation and program input to son process.
	mes := fmt.Sprintf("%s %s", mntPath, binaryName)
	_, err = writePipe.WriteString(mes)
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot write pipe: %w", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot get stdin: %w", err)
	}
	_, err = stdin.Write([]byte(inputContent))
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot write input content to stdin: %w", err)
	}

	// stage2: start two goroutine, one for starting son process then wait to it exit,
	// the other for counting total seconds to kill son process.
	errChan := make(chan error)
	dataChan := make(chan []byte)

	// start process
	err = cmd.Start()
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot start sub process: %w", err)
	}

	// waiting for sub procss start signal then add its pid to cgroups
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGUSR1)
	<-sign
	// start to adding pid to cgroups
	err = cv.Apply(cmd.Process.Pid)
	if err != nil {
		return nil, inteErrCode, fmt.Errorf("cannot apply cgroups for sub process: %w", err)
	}
	// after adding
	syscall.Kill(cmd.Process.Pid, syscall.SIGUSR1)

	go func() {
		var data bytes.Buffer
		cmd.Stdout = &data
		cmd.Stderr = &data
		err := cmd.Wait()
		if err != nil {
			errChan <- err
		}
		dataChan <- data.Bytes()
	}()

	go func() {
		// sleep total seconds
		for i := 0; i < int(sec*10); i++ {
			err = sleep.NanoSleep(100)
			if err != nil {
				errChan <- err
			}
		}
		errChan <- os.ErrDeadlineExceeded
	}()

	var data []byte
	// successfully exited or exceed time to be killed.
	select {
	// case: any errors happen
	case err = <-errChan:
		// kill pgroup
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	// case: process sucessfully exists
	case data = <-dataChan:
	}
	exitCode := cmd.ProcessState.ExitCode()
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
		return pb_sb.Status_UE
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

// Read all input data from number 3 pipe.
func readInputFromPipe() (*pb_sb.Input, error) {
	if bytes, err := ioutil.ReadAll(os.NewFile(uintptr(3), "pipe")); err != nil {
		return nil, fmt.Errorf("parent process cannot read pipe: [%w]", err)
	} else {
		input := pb_sb.Input{}
		return &input, proto.Unmarshal(bytes, &input)
	}
}
