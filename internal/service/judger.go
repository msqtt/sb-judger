package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	pb_jg "github.com/msqtt/sb-judger/api/pb/v1/judger"
	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	"github.com/msqtt/sb-judger/internal/compile"
	"github.com/msqtt/sb-judger/internal/pkg/config"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/msqtt/sb-judger/internal/sandbox"
	"github.com/msqtt/sb-judger/internal/sandbox/fs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
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
	resp *pb_jg.RunCodeResponse, err error) {
	// var err error
	inputContent := req.GetInput()
	code := req.GetCode()
	lang := req.GetLang()
	time := req.GetTime()
	mem := req.GetMemory()

	lc := js.langConfMap[lang.String()]
	if lc == nil {
		log.Println(ErrNotSupportedLang)
		err = status.Error(codes.InvalidArgument, "not supported language")
		return
	}

	// make a tempPath for working
	tempPath, err := os.MkdirTemp(js.conf.WorkDir, "sb-judger-*")
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to mkdir temp")
		return
	}
	
	// stage1: compiling code
	cmd, err := compile.CreateCompileCmd(tempPath, lang.String(), code, *lc)
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to create compile cmd")
	}

	log.Println(cmd.String())
	msg, exitCode, err := compile.RunCmdCombinded(cmd)
	// if compile fails, return error message directly.
	if exitCode != 0 {
		errorMsg := strings.ReplaceAll(msg, tempPath, "...")
		log.Println(err)
		err = nil
		resp = &pb_jg.RunCodeResponse{OutPut: errorMsg}
		return
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
		err = status.Error(codes.Internal, "failed to make overlayfs")
		return
	}
	defer overlayfs.Remove()

	err = overlayfs.Move(filepath.Join(tempPath, lc.Out), filepath.Join("/tmp", lc.Out))
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to move binary file")
		return
	}

	// stage3: run process
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to create pipe")
		return
	}
	h := sha1.New()
	h.Write([]byte(code))
	hashName := base64.URLEncoding.EncodeToString(h.Sum(nil))
	inputMsg := &pb_sb.Input{
		HashName: hashName,
		Lang: str2pbLang(lang.String()),
		Time: time,
		Memory: mem,
		MntPath: tempPath,
		Cases: []*pb_sb.Case{
			{
				CaseId: 1,
				In: inputContent,
				Out: "",
			},
		},
	}
	b, err := proto.Marshal(inputMsg)
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to marshal input message")
		return
	}
	_, err = writePipe.Write(b)
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to write to pipe")
		return
	}
	cmd = exec.Command("./sandbox", strconv.Itoa(sandbox.ArgInit))
	cmd.ExtraFiles = []*os.File{readPipe}
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to run cmd")
		return
	}
	collectOut := new(pb_sb.CollectOutput)
	err = proto.Unmarshal(data, collectOut)
	if err != nil {
		log.Println(err)
		err = status.Error(codes.Internal, "failed to unmarshal out")
		return
	}
	out := collectOut.GetCaseOuts()[0]
	resp = &pb_jg.RunCodeResponse{
		OutPut: out.OutPut,
		TimeUsage: out.TimeUsage,
		MemoryUsage: out.MemoryUsage,
	}
	return
}

func str2pbLang(str string) (pb_sb.Language) {
	switch str {
	case "c":
		return pb_sb.Language_c
	case "cpp":
		return pb_sb.Language_cpp
	case "golang":
		return pb_sb.Language_golang
	case "java":
		return pb_sb.Language_java
	case "python":
		return pb_sb.Language_python
	case "rust":
		return pb_sb.Language_rust
	default:
		return -1
	}
}

var _ pb_jg.CodeServer = (*JudgerServer)(nil)

func NewJudgerServer(conf config.Config, langMap json.LangConfMap) *JudgerServer {
	return &JudgerServer{
		conf:        conf,
		langConfMap: langMap}
}
