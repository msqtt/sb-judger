package service

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"testing"

	pb_jg "github.com/msqtt/sb-judger/api/pb/v1/judger"
	pb_sb "github.com/msqtt/sb-judger/api/pb/v1/sandbox"
	"github.com/msqtt/sb-judger/internal/pkg/config"
	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestJudgerCode(t *testing.T) {
	conf, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalln(err)
	}
	lcm, err := json.GetLangConfs("./configs/lang.json")
	if err != nil {
		log.Fatalln(err)
	}
	addr := startTestJudgerServer(t, conf, lcm)
	client := newTestJudgerClient(t, addr)

	testCases := []struct {
		ctx       context.Context
		name      string
		lang      pb_sb.Language
		code      string
		memLimit  uint32
		timeLimit uint32
		cases     []*pb_sb.Case
		state     []pb_sb.State
	}{
		{
			ctx:       context.Background(),
			name:      "python-hello",
			lang:      pb_sb.Language_python,
			code:      readCode(t, "./testcode/service/python/hello.py"),
			memLimit:  64,
			timeLimit: 1000,
			cases: newCasesType().
				buildCases(1, "sandbox", "hello sandbox").
				buildCases(2, " ", "hello").
				buildCases(3, "  ", "hello   ").
				buildCases(4, "cat", "hello dog"),
			state: []pb_sb.State{
				pb_sb.State_WA,
				pb_sb.State_AC,
				pb_sb.State_AC,
				pb_sb.State_AC,
				pb_sb.State_WA,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.code)
			req := &pb_jg.JudgeCodeRequest{
				Lang:        tc.lang,
				Code:        tc.code,
				Time:        tc.timeLimit,
				Memory:      tc.memLimit,
				OutMsgLimit: 0,
				Case:        tc.cases,
			}
			jcr, err2 := client.JudgeCode(tc.ctx, req)
			require.NoError(t, err2)
			require.NotNil(t, jcr)

			require.Equal(t, tc.state[0], jcr.State)
			crs := jcr.GetCodeResults()
			require.NotEmpty(t, crs)
			for i, cr := range crs {
				t.Log(cr)
				require.Equal(t, tc.state[i+1], cr.GetState())
				require.Less(t, cr.MemoryUsage, float64(tc.memLimit<<10))
				require.Less(t, cr.RealTimeUsage, float64(tc.timeLimit))
				require.Less(t, cr.CpuTimeUsage, float64(tc.timeLimit))
			}
		})
	}

}

type casesType []*pb_sb.Case

func newCasesType() casesType {
	return make(casesType, 0)
}
func (cs casesType) buildCases(id uint32, in, out string) casesType {
	return append(cs, &pb_sb.Case{CaseId: id, In: in, Out: out})
}

func readCode(t *testing.T, path string) string {
	f, err := os.Open(path)
	require.NoError(t, err)
	b, err := io.ReadAll(f)
	require.NoError(t, err)
	return string(b)
}

func startTestJudgerServer(t *testing.T, conf config.Config, lcm json.LangConfMap) string {
	js := NewJudgerServer(conf, lcm)
	server := grpc.NewServer()
	pb_jg.RegisterCodeServer(server, js)
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go server.Serve(l)
	return l.Addr().String()
}

func newTestJudgerClient(t *testing.T, addr string) pb_jg.CodeClient {
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb_jg.NewCodeClient(cc)
}
