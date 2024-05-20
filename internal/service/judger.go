package service

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unicode"

	gonanoid "github.com/matoous/go-nanoid/v2"
	pb_jg "github.com/msqtt/sb-judger/api/pb/v1/judger"
	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	"github.com/msqtt/sb-judger/internal/compile"
	"github.com/msqtt/sb-judger/internal/pkg/config"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/msqtt/sb-judger/internal/sandbox"
	"github.com/msqtt/sb-judger/internal/sandbox/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotSupportedLang = errors.New("not supported language.")

type JudgerServer struct {
	pb_jg.UnimplementedCodeServer
	conf        config.Config
	langConfMap json.LangConfMap
	limit       chan struct{}
}

type judgeCodeResult struct {
	resp *pb_jg.JudgeCodeResponse
	err  error
}

// JudgeCode implements pb_jg.CodeServer.
func (js *JudgerServer) JudgeCode(ctx context.Context, req *pb_jg.JudgeCodeRequest) (
	*pb_jg.JudgeCodeResponse, error) {
	code := req.GetCode()
	lang := req.GetLang()
	timeLimit := req.GetTime()
	memLimit := req.GetMemory()
	outLimit := req.GetOutMsgLimit()
	cases := req.GetCase()

	retCh := make(chan judgeCodeResult)

	// long time actions
	go func() {
		collectOut, err := js.runCode(lang.String(), code,
			memLimit, timeLimit, outLimit, cases)

		if err != nil {
			var errCompile *compile.ErrCompileMsg
			if errors.As(err, &errCompile) {
				retCh <- judgeCodeResult{
					resp: &pb_jg.JudgeCodeResponse{OutPut: err.Error(), State: pb_sb.State_CE},
					err:  nil}
				return
			}
			retCh <- judgeCodeResult{resp: nil, err: err}
			return
		}

		outs := collectOut.GetCaseOuts()
		if len(outs) <= 0 {
			retCh <- judgeCodeResult{
				resp: nil,
				err:  status.Error(codes.Internal, "failed to collect output")}
			return
		}

		log.Println(outs)

		cr := make([]*pb_jg.CodeResult, len(outs))
		var finalState pb_sb.State = -1
		var finalTimeUsage float64
		var finalMemoryUsage float64

		for i, o := range outs {
			if o.State != pb_sb.State_AC {
				finalState = o.State
			}

			cpuTimeUsage := float64(o.CpuTimeUsage)
			realTimeUsage := float64(o.RealTimeUsage)
			memoryTimeUsage := float64(o.MemoryUsage)

			finalTimeUsage = math.Max(finalTimeUsage, realTimeUsage)
			finalMemoryUsage = math.Max(finalMemoryUsage, memoryTimeUsage)

			cr[i] = &pb_jg.CodeResult{
				CaseId:        o.CaseId,
				CpuTimeUsage:  cpuTimeUsage / 1000,
				RealTimeUsage: realTimeUsage / 1000,
				MemoryUsage:   memoryTimeUsage / 1024,
				State:         o.State,
			}
		}
		if finalState == -1 {
			finalState = pb_sb.State_AC
		}

		retCh <- judgeCodeResult{
			resp: &pb_jg.JudgeCodeResponse{
				State:          finalState,
				MaxTimeUsage:   finalTimeUsage / 1000,
				MaxMemoryUsage: finalMemoryUsage / 1024,
				OutPut:         "",
				CodeResults:    cr},
			err: nil,
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("request cancel...")
		return nil, nil
	case ret := <-retCh:
		return ret.resp, ret.err
	}
}

type runCodeResult struct {
	resp *pb_jg.RunCodeResponse
	err  error
}

// RunCode implements pb_jg.CodeServer.
func (js *JudgerServer) RunCode(ctx context.Context, req *pb_jg.RunCodeRequest) (
	*pb_jg.RunCodeResponse, error) {
	code := req.GetCode()
	lang := req.GetLang()
	timeLimit := req.GetTime()
	memLimit := req.GetMemory()
	outLimit := req.GetOutMsgLimit()
	inputContent := req.GetInput()

	retCh := make(chan runCodeResult)

	// long time actions
	go func() {
		collectOut, err := js.runCode(lang.String(), code,
			memLimit, timeLimit, outLimit,
			[]*pb_sb.Case{
				{
					CaseId: 1,
					In:     inputContent,
					Out:    "",
				},
			},
		)
		if err != nil {
			var errCompile *compile.ErrCompileMsg
			if errors.As(err, &errCompile) {
				retCh <- runCodeResult{
					resp: &pb_jg.RunCodeResponse{OutPut: err.Error(), State: pb_sb.State_CE},
					err:  nil,
				}
				return
			}

			retCh <- runCodeResult{
				resp: nil, err: err,
			}
			return
		}

		outs := collectOut.GetCaseOuts()
		if len(outs) <= 0 {
			retCh <- runCodeResult{
				resp: nil, err: status.Error(codes.Internal, "failed to collect output"),
			}
			return
		}
		out := outs[0]
		if out.State == pb_sb.State_WA {
			out.State = pb_sb.State_AC
		}

		retCh <- runCodeResult{
			resp: &pb_jg.RunCodeResponse{
				OutPut:        out.OutPut,
				CpuTimeUsage:  float64(out.CpuTimeUsage) / 1000,
				RealTimeUsage: float64(out.RealTimeUsage) / 1000,
				MemoryUsage:   float64(out.MemoryUsage) / 1024,
				State:         out.State,
			}, err: nil}
	}()

	select {
	case <-ctx.Done():
		log.Println("request cancel...")
		return nil, nil
	case ret := <-retCh:
		return ret.resp, ret.err
	}
}

// runCode makes a sandbox to run program and judges each cases then returns outputs.
// runCode will limit the specific number(js.limit) of concurrency process.
func (js *JudgerServer) runCode(lang, code string,
	memLimit, timeLimit, outLimit uint32, cases []*pb_sb.Case) (*pb_sb.CollectOutput, error) {

	if strings.TrimRightFunc(code, unicode.IsSpace) == "" {
		return nil, status.Error(codes.InvalidArgument, "code cannot be none")
	}

	if timeLimit > 2000 {
		return nil, status.Error(codes.InvalidArgument, "time limit should be in [0, 2000]")
	}

	if memLimit > 256 || memLimit < 1 {
		return nil, status.Error(codes.InvalidArgument, "memory limit should be in [1, 256]")
	}

	if outLimit > 1024 {
		return nil, status.Error(codes.InvalidArgument, "output limit should be in [0, 1024]")
	}

	conf := js.conf
	lc := js.langConfMap[lang]
	if lc == nil {
		log.Println(ErrNotSupportedLang)
		return nil, status.Error(codes.InvalidArgument, "not supported language")
	}

	// make a tempPath for working
	tempPath, err := os.MkdirTemp(conf.WorkDir, "sb-judger*")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to mkdir temp")
	}

	// stage1: compiling code
	cmd, err := compile.CreateCompileCmd(tempPath, lang, code, *lc)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to create compile cmd")
	}

	log.Println(cmd.String())
	msg, exitCode, err := compile.RunCmdCombinded(cmd)
	// if compile fails, return error message directly.
	if exitCode != 0 || err != nil {
		defer os.RemoveAll(tempPath)
		r, _ := regexp.Compile("/.[^\\s]*/")
		msg = r.ReplaceAllString(msg, ".../")
		log.Println(err)
		err = nil
		return nil, &compile.ErrCompileMsg{Msg: msg}
	}
	// chmod 755 for compile result
	compileOutPath := filepath.Join(tempPath, lc.Out)
	err = syscall.Chmod(compileOutPath, 0755)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to chmod for compile out file")
	}

	// stage2: builds a rootfs for running program.
	overlayfs, err := fs.NewOverlayfs(conf.RootFsDir)
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to new overlayfs")
	}
	err = overlayfs.Make(tempPath)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to make overlayfs")
	}
	defer overlayfs.Remove()

	err = overlayfs.Move(compileOutPath, filepath.Join("/tmp", lc.Out))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to move binary file")
	}

	// before runing, check the limit number of concurrencies
	<-js.limit
	// release the number
	defer func() { js.limit <- struct{}{} }()

	// stage3: run process
	id, err := gonanoid.New()
	mntPath := path.Join(tempPath, "mnt")

	// prepare outputContent limit.
	var outContentLimit uint32
	if outLimit <= 0 {
		outContentLimit = uint32(conf.OutLenLimit << 10)
	} else {
		outContentLimit = outLimit << 10
	}

	collectOut, err := sandbox.InitEntry(lang, id, mntPath,
		outContentLimit, memLimit, timeLimit, cases,
	)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to run program")
	}

	return collectOut, nil
}

func str2pbLang(str string) pb_sb.Language {
	i, ok := pb_sb.State_value[str]
	if !ok {
		return -1
	}
	return pb_sb.Language(i)
}

var _ pb_jg.CodeServer = (*JudgerServer)(nil)

func NewJudgerServer(conf config.Config, langMap json.LangConfMap) *JudgerServer {
	ret := &JudgerServer{
		conf:        conf,
		langConfMap: langMap,
		limit:       make(chan struct{}, conf.ConcurrencyLimit),
	}
	// setup limit of concurrencies.
	for i := 0; i < conf.ConcurrencyLimit; i++ {
		ret.limit <- struct{}{}
	}
	return ret
}
