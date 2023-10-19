package res

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var hashPath = hex.EncodeToString(sha1.New().Sum([]byte("testhashpath")))

func TestCgroupV2CRD(t *testing.T) {
	cgroupv2, err := NewCgroupV2(hashPath)
	require.NoError(t, err)
	require.NotEmpty(t, cgroupv2)

	testreadPath(t, cgroupv2.path)

	config := &ResourceConfig{
		CpuLimit:    100000,
		MemoryLimit: 100,
		PidsLimit:   1,
	}

	err = cgroupv2.Config(config)
	require.NoError(t, err)

	testCompareFileConent(t, filepath.Join(cgroupv2.path, "pids.max"),
		strconv.Itoa(int(config.PidsLimit)))

	mems := []string{"memory.max", "memory.high", "memory.swap.max", "memory.swap.high"}

	for _, memFile := range mems {
		testCompareFileConent(t, filepath.Join(cgroupv2.path, memFile),
			strconv.Itoa(int(config.MemoryLimit<<20)))
	}

	num := 3
	bk := make([]*CgroupV2, 0)
	for i := 1; i <= num; i++ {
		cv, err := cgroupv2.AddSubGroup(fmt.Sprintf("testsub%d", i))
		require.NoError(t, err)
		require.NotEmpty(t, cv)
		testreadPath(t, cv.path)
		bk = append(bk, cv)
	}

	cvs := cgroupv2.GetSons()
	require.Len(t, cvs, 3)
	require.ElementsMatch(t, bk, cvs)

	err = cgroupv2.Destroy()
	require.NoError(t, err)

	_, err = os.Stat(cgroupv2.path)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestCgroupV2Apply(t *testing.T) {
	cgroupv2, err := NewCgroupV2(hashPath)
	require.NoError(t, err)
	require.NotEmpty(t, cgroupv2)

	config := &ResourceConfig{
		CpuLimit:    100000,
		MemoryLimit: 100,
		PidsLimit:   1,
	}

	err = cgroupv2.Config(config)
	require.NoError(t, err)

	cv, err := cgroupv2.AddSubGroup("1")
	require.NoError(t, err)
	require.NotEmpty(t, cv)

	cmd := exec.Command("sleep", "1s")
	err = cmd.Start()
	require.NoError(t, err)

	cpid := cmd.Process.Pid
	err = cv.Apply(cpid)
	require.NoError(t, err)

	testCompareFileConent(t, filepath.Join(cv.path, "cgroup.procs"), strconv.Itoa(cpid))

	err = cmd.Wait()
	require.NoError(t, err)

	err = cgroupv2.Destroy()
	require.NoError(t, err)
}

func testreadPath(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err)

	dirs, err := os.ReadDir(path)
	require.NoError(t, err)
	require.NotEmpty(t, dirs)
}

func testCompareFileConent(t *testing.T, path string, content string) {
	byte, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.NotEmpty(t, byte)
	require.Equal(t, content, strings.Trim(string(byte), "\n"))
}
