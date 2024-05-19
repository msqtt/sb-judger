package compile

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/msqtt/sb-judger/internal/pkg/json"
	"github.com/stretchr/testify/require"
)

func TestCreateCompileCmd(t *testing.T) {
	const testGoodPath = "./testcode/compile/good"
	const testBadPath = "./testcode/compile/bad"

	lcm, err := json.GetLangConfs("./lang_test.json")
	require.NoError(t, err)
	require.NotEmpty(t, lcm)

	testCases := []struct {
		lang string
		code string
	}{
		{
			lang: "c",
		},
		{
			lang: "cpp",
		},
		{
			lang: "golang",
		},
		{
			lang: "java",
		},
		{
			lang: "python",
		},
		{
			lang: "rust",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		conf := lcm[tc.lang]
		t.Run(tc.lang+"good", func(t *testing.T) {
			t.Parallel()
			runCreateCmd(t, true, testGoodPath, *conf, tc)
		})
		t.Run(tc.lang+"bad", func(t *testing.T) {
			t.Parallel()
			runCreateCmd(t, false, testBadPath, *conf, tc)
		})
	}
}

func runCreateCmd(t *testing.T, isGood bool, path string, conf json.LanguageConfig, tc struct{ lang, code string }) {
	srcPath := filepath.Join(path, conf.Src)
	tempPath, err := os.MkdirTemp(path, "sb-judger*")
	require.NoError(t, err)
	require.NotEmpty(t, tempPath)

	b, err := os.ReadFile(srcPath)
	require.NoError(t, err)
	require.NotEmpty(t, b)
	tc.code = string(b)
	c, err2 := CreateCompileCmd(tempPath, tc.lang, tc.code, conf)
	require.NoError(t, err2)
	require.NotNil(t, c)
	require.NotEmpty(t, tempPath)
	log.Print(tempPath)

	fi, err3 := os.Stat(tempPath)
	require.NoError(t, err3)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir())
	require.NotZero(t, fi.Size())

	duplicatedFileName := filepath.Join(tempPath, conf.Src)
	log.Println(duplicatedFileName)
	b2, err := os.ReadFile(duplicatedFileName)
	require.NoError(t, err)
	require.NotEmpty(t, b)
	require.ElementsMatch(t, b, b2)

	if isGood {
		testGoodCompile(t, c, filepath.Join(tempPath, conf.Out))
	} else {
		testBadCompile(t, c, filepath.Join(tempPath, conf.Out))
	}

	err4 := os.RemoveAll(tempPath)
	require.NoError(t, err4)
}

func testBadCompile(t *testing.T, cmd *exec.Cmd, outPath string) {
	b, err := cmd.CombinedOutput()
	require.Error(t, err)
	require.NotEmpty(t, b)
	log.Println("------------------------")
	log.Println(err)
	log.Println(string(b))
	exitCode := cmd.ProcessState.ExitCode()
	require.NotZero(t, exitCode)
	log.Println(exitCode)
	fi, err2 := os.Stat(outPath)
	require.Error(t, err2)
	require.Nil(t, fi)
}

func testGoodCompile(t *testing.T, cmd *exec.Cmd, outPath string) {
	b, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Empty(t, b)
	exitCode := cmd.ProcessState.ExitCode()
	require.Zero(t, exitCode)
	fi, err2 := os.Stat(outPath)
	require.NoError(t, err2)
	require.NotEmpty(t, fi)
}
