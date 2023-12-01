package service

import (
	"context"
	"errors"
	"log"
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
}

// JudgeCode implements pb_jg.CodeServer.
func (js *JudgerServer) JudgeCode(ctx context.Context, req *pb_jg.JudgeCodeRequest) (
	resp *pb_jg.JudgeCodeResponse, err error) {
	panic("unimplemented")
}

// RunCode implements pb_jg.CodeServer.
func (js *JudgerServer) RunCode(ctx context.Context, req *pb_jg.RunCodeRequest) (
	*pb_jg.RunCodeResponse, error) {
	code := req.GetCode()
	lang := req.GetLang()
	time := req.GetTime()
	mem := req.GetMemory()
	inputContent := req.GetInput()
	outMsgLimit := req.GetOutMsgLimit()

	if strings.TrimRightFunc(code, unicode.IsSpace) == "" {
		return nil, status.Error(codes.InvalidArgument, "code cannot be none")
	}

	if time > 2000 {
		return nil, status.Error(codes.InvalidArgument, "time limit should be in [0, 2000]")
	}

	if mem > 256 || mem < 1 {
		return nil, status.Error(codes.InvalidArgument, "memory limit should be in [1, 256]")
	}
 
  if outMsgLimit > 1024 {
		return nil, status.Error(codes.InvalidArgument, "output limit should be in [0, 1024]")
  }


	lc := js.langConfMap[lang.String()]
	if lc == nil {
		log.Println(ErrNotSupportedLang)
		return nil, status.Error(codes.InvalidArgument, "not supported language")
	}

	// make a tempPath for working
	tempPath, err := os.MkdirTemp(js.conf.WorkDir, "sb-judger*")
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to mkdir temp")
	}

	// stage1: compiling code
	cmd, err := compile.CreateCompileCmd(tempPath, lang.String(), code, *lc)
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
		return &pb_jg.RunCodeResponse{OutPut: msg, State: "CE"}, nil
	}
	// chmod 755 for program
	compileOutPath := filepath.Join(tempPath, lc.Out)
	err = syscall.Chmod(compileOutPath, 0755)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to chmod for compile out file")
	}

	// stage2: builds a rootfs for running program.
	overlayfs, err := fs.NewOverlayfs(js.conf.RootFsDir)
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

	// stage3: run process
  id, err := gonanoid.New()
	mntPath := path.Join(tempPath, "mnt")

	// prepare outputContent limit.
	var outContentLimit uint32
	if outMsgLimit <= 0 {
		outContentLimit = uint32(js.conf.OutLenLimit << 10)
	} else {
		outContentLimit = outMsgLimit << 10
	}

	collectOut, err := sandbox.InitEntry(lang.String(), id, mntPath,
		outContentLimit, mem, time,
		[]*pb_sb.Case{
			{
				CaseId: 1,
				In:     inputContent,
				Out:    "",
			},
		},
	)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to run program")
	}
	outs := collectOut.GetCaseOuts()
	if len(outs) <= 0 {
		return nil, status.Error(codes.Internal, "failed to collect output")
	}
	out := outs[0]
  state := out.Status.String()
  if state == "WA" || state == "AC" {
    state = ""
  }
	return &pb_jg.RunCodeResponse{
		OutPut:        out.OutPut,
		CpuTimeUsage:  float32(out.CpuTimeUsage) / 1000,
		RealTimeUsage: float32(out.RealTimeUsage) / 1000,
		MemoryUsage:   float32(out.MemoryUsage) / 1024,
    State: state,
	}, nil
}

func str2pbLang(str string) pb_sb.Language {
	i, ok := pb_sb.Status_value[str]
	if !ok {
		return -1
	}
	return pb_sb.Language(i)
}

var _ pb_jg.CodeServer = (*JudgerServer)(nil)

func NewJudgerServer(conf config.Config, langMap json.LangConfMap) *JudgerServer {
	return &JudgerServer{
		conf:        conf,
		langConfMap: langMap}
}
