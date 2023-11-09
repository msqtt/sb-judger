package compile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/msqtt/sb-judger/internal/pkg/json"
)

// CreateCompileCmd saves codei which will be compiled into a temp dir inside named destPath.
// it returns a command type, tempPath, and an error, if any.
func CreateCompileCmd(tempPath, lang, code string, conf json.LanguageConfig) (
	cmd *exec.Cmd, err error) {

	filePath := filepath.Join(tempPath, conf.Src)
	data := []byte(code)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return
	}

	compileFlags := strings.Join(conf.Compile.Flags, " ")
	
	src := filepath.Join(tempPath, conf.Src)
	out := filepath.Join(tempPath, conf.Out)
	if lang == "golang" {
		src, out = out, src
	}
	// magic code makes my days.
	compileFlags = fmt.Sprintf(compileFlags, src, out)
	compileFlags = strings.Split(compileFlags, "%!")[0]
	flags := strings.Split(compileFlags, " ")

	cmd = exec.Command(conf.Cmd, flags...)
	return
}
