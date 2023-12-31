package sandbox

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	msgq "github.com/msqtt/sb-judger/internal/message"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	res "github.com/msqtt/sb-judger/internal/sandbox/resource"
	"golang.org/x/sys/unix"
)

const inteErrCode = 500

type processResult struct {
	outPut   []byte
	exitCode int
	err      error
}

// Entry function is the intro of sandbox program.
func InitEntry(lang, id, mntPath string, outLimit, mem, time uint32, cases []*pb_sb.Case) (
	*pb_sb.CollectOutput, error) {
	langConf, err := json.GetLangConfs("")
	if err != nil {
		return nil, fmt.Errorf("cannot get lang config: %w", err)
	}

	lc := langConf[lang]
	if lc == nil {
		return nil, errors.New("not supported language")
	}

	resManager, err := res.NewCgroupV2(id)
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

	outPath := filepath.Join("/tmp", lc.Out)
	// still, lazy codes 🤓
	runCmd := strings.Join(lc.Run, " ")
	runCmd = fmt.Sprintf(runCmd, outPath)
	runCmd = strings.Split(runCmd, "%!")[0]

	collectOuts := make([]*pb_sb.Output, len(cases))
	// todo: using goroutine to fork per case sub process.
	for i, cas := range cases {
		// add sub cgroup for per case.
		cv, err := resManager.AddSubGroup(strconv.FormatUint(uint64(cas.GetCaseId()), 10))
		if err != nil {
			return nil, fmt.Errorf("cannot add sub group: %w", err)
		}

		// passing id and input string
		bytes, code, realTime, err := execLaunchProcess(cv, lc, outLimit, time, mntPath, runCmd, cas.In)
		if code == inteErrCode {
			// internal errror happens just cut off
			return nil, fmt.Errorf("cannot collect sub process: [%w, msg: %s]", err, string(bytes))
		}

		// reads resource state after launch process done.
		usage, err1 := cv.ReadState()
		if err1 != nil {
			return nil, fmt.Errorf("cannot read states from case %d process: %w", cas.GetCaseId(), err)
		}

		// check status
		status := checkState(bytes, code, err, []byte(cas.GetOut()), usage, time, mem)

		// collection result
		collectOuts[i] = &pb_sb.Output{
			CaseId:        cas.CaseId,
			CpuTimeUsage:  usage.CpuTime,
			RealTimeUsage: uint32(realTime),
			MemoryUsage:   usage.MemoryUsage,
			State:         status,
			OutPut:        string(bytes),
		}
	}

	return &pb_sb.CollectOutput{CaseOuts: collectOuts}, nil
}

// execLaunchProcess exec a process then set the namespaces for it.
func execLaunchProcess(cv *res.CgroupV2, lc *json.LanguageConfig, outLim, sec uint32, mntPath,
	runCmd, inputContent string) (
	[]byte, int, int, error) {
	// passing launch arg to execute launch entry parts.
	cmd, writePipe, err := maskFork()
	if err != nil {
		return nil, inteErrCode, 0, err
	}

	// write args, input content to launch process
	// and binding stdout and stderr.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, inteErrCode, 0, fmt.Errorf("cannot get stdin: %w", err)
	}
	_, err = stdin.Write([]byte(inputContent))
	if err != nil {
		return nil, inteErrCode, 0, fmt.Errorf("cannot write input content to stdin: %w", err)
	}
	stdin.Close()

	// binding stdout and stderr
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	retCh := make(chan processResult)

	// start launch process.
	err = cmd.Start()
	if err != nil {
		return nil, inteErrCode, 0, fmt.Errorf("cannot start sub process: %w", err)
	}

	msqid, err := msgq.OpenQueue(cmd.Process.Pid)
	if err != nil {
		log.Println(err)
	}

	mes := fmt.Sprintf("%s#%s#%s", mntPath, runCmd, fmt.Sprint(msqid))
	_, err = writePipe.WriteString(mes)
	if err != nil {
		return nil, inteErrCode, 0, fmt.Errorf("cannot write pipe: %w", err)
	}
	writePipe.Close()

	go func() {
		err := cmd.Wait()
		var outMsg []byte
		// erase outrange messages.
		if buf.Len() > int(outLim) {
			outMsg = make([]byte, outLim)
			buf.Read(outMsg)
			// fix utf-8 inter split
			r := bytes.Runes(outMsg)
			outMsg = append([]byte(string(r)), []byte("\n...\n")...)
		} else {
			outMsg = buf.Bytes()
		}
		retCh <- processResult{
			outPut:   outMsg,
			exitCode: cmd.ProcessState.ExitCode(),
			err:      err,
		}
	}()

	var ret processResult
	var startAt time.Time
	var endAt time.Time
	var realTime int
loop:
	// one round for message coming, the other for checking that
	// process successfully exited or exceed time to be killed
	// or bad things happened.
	for i := 0; i < 2; i++ {
		select {
		case ret = <-retCh:
			endAt = time.Now()
			realTime = int(endAt.Sub(startAt).Microseconds())
			// case: internal err happens
			if ret.exitCode == inteErrCode {
				err := cmd.Process.Kill()
				if err != nil {
					ret.err = fmt.Errorf("cannot kill process: %w", ret.err)
				}
				return ret.outPut, ret.exitCode, realTime, ret.err
			}
			// case: program exceeds time limit.
			if ret.err == os.ErrDeadlineExceeded {
				err := cmd.Process.Kill()
				if err != nil {
					ret.err = fmt.Errorf("cannot kill process: %w", ret.err)
				}
				// collect program output
				retWait := <-retCh
				ret.outPut = retWait.outPut
			}
			// case: program exists successfully.
			break loop

		case <-msgq.MsgChan(msqid, msgq.NewMsg(1, nil)):
			// set pid and cgroup namespace for sub process.
			err := setNsFor(cmd.Process.Pid, unix.CLONE_NEWIPC)
			if err != nil {
				err = fmt.Errorf("cannot setns for sub process: %w", err)
				err2 := cmd.Process.Kill()
				if err2 != nil {
					log.Println(err2)
					err = fmt.Errorf("cannot kill process: %w", err)
				}
				return nil, inteErrCode, realTime, err
			}
			// start to adding pid to cgroups
			err = cv.Apply(cmd.Process.Pid)
			if err != nil {
				err = fmt.Errorf("cannot apply cgroups for sub process: %w", err)
				err2 := cmd.Process.Kill()
				if err2 != nil {
					log.Println(err2)
					err = fmt.Errorf("cannot kill process: %w", err)
				}
				return nil, inteErrCode, 0, err
			}
			// set timer to watch
			go func() {
				// sleep total seconds
				startAt = time.Now()
				time.Sleep(time.Duration(sec) * time.Millisecond)
				retCh <- processResult{
					outPut:   []byte("time limit exceeded"),
					exitCode: 1,
					err:      os.ErrDeadlineExceeded,
				}
			}()
			// tell sub process it is ok to run program after adding.
			// err = msgq.SndMsg(msqid, msgq.NewMsg(2, nil))
			err = msgq.DestroyQueue(msqid)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return trimMessage(mntPath, ret.outPut, lc), ret.exitCode, realTime, ret.err
}

// trimMessage delete the message
func trimMessage(mntPath string, out []byte, lc *json.LanguageConfig) []byte {
	path := strings.TrimRight(mntPath, "/mnt")
	out = bytes.ReplaceAll(out, []byte(path), []byte("..."))

	if len(lc.TrimMsg) <= 0 {
		return out
	}
	for _, v := range lc.TrimMsg {
		out = bytes.ReplaceAll(out, []byte(v), []byte(""))
	}
	return out
}

// checkState returns running state depending on output、exitcode and err
func checkState(outCont []byte, code int, outErr error, ans []byte,
	usage *res.RunState, timeLimit, memLimit uint32) pb_sb.State {
	// trim right space (is it needed ?)
	outCont = bytes.TrimRightFunc(outCont, unicode.IsSpace)
	ans = bytes.TrimRightFunc(ans, unicode.IsSpace)

	if usage.CpuTime/1000 > timeLimit {
		return pb_sb.State_TLE
	}
	if usage.MemoryUsage > memLimit<<20 {
		return pb_sb.State_MLE
	}
	if usage.OOMKill > 0 {
		return pb_sb.State_MLE
	}

	if errors.Is(outErr, os.ErrDeadlineExceeded) {
		return pb_sb.State_TLE
	}

	switch code {
	case 0:
		// diff between answer and user printout.
		if bytes.Compare(outCont, ans) == 0 {
			return pb_sb.State_AC
		} else {
			return pb_sb.State_WA
		}
	case 1, 2, 136, 139:
		return pb_sb.State_RE
	default:
		return pb_sb.State_UE
	}
}
