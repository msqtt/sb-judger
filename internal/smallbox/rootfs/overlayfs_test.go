package fs

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOverlayfs(t *testing.T) {
	testCases := []struct {
		name       string
		rootfsPath string
	}{
		{
			name:       "success",
			rootfsPath: testRootfsPath,
		},
		{
			name:       "not_found",
			rootfsPath: "./hpquyquq",
		},
		{
			name:       "not_dir",
			rootfsPath: "./fs.go",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			overlayfs, err := NewOverlayfs(tc.rootfsPath)

			if err != nil {
				switch tc.name {
				case "not_found":
					require.Nil(t, overlayfs)
					require.True(t, strings.Contains(err.Error(), "cannot found rootfs: "))
				case "not_dir":
					require.Nil(t, overlayfs)
					require.True(t, strings.Contains(err.Error(), "given roots is not a dir!"))
				default:
					t.FailNow()
				}
				t.Skip()
			}

			// success case
			require.NotNil(t, overlayfs)

			dir := makeTestMnt(t, overlayfs)
			for _, fi := range dir {
				log.Println(fi.Name())
			}

			err = overlayfs.Remove()
			require.NoError(t, err)
		})
	}
}

func makeTestMnt(t *testing.T, overlayfs *Overlayfs) []fs.FileInfo {
	err := os.Mkdir(testMntPath, 0755)
	require.NoError(t, err)

	err = overlayfs.Make(testMntPath)
	require.NoError(t, err)

	mnt, err := overlayfs.MountPoint()
	require.NoError(t, err)

	dir, err := ioutil.ReadDir(mnt)
	require.NoError(t, err)
	require.NotEmpty(t, dir)

	return dir
}

func TestPutFile(t *testing.T) {
	o, err := NewOverlayfs(testRootfsPath)
	require.NoError(t, err)
	require.NotNil(t, o)
	
	testPath := "./test"
	err = os.Mkdir(testPath, 0755)
	require.NoError(t, err)

	err2 := o.Make(testPath)
	require.NoError(t, err2)

	b := []byte("hello")
	err3 := o.PutFile("/tmp", "1.txt", bytes.NewReader(b))
	require.NoError(t, err3)

	b2, err4 := os.ReadFile(filepath.Join(o.upperPath, "tmp", "1.txt"))
	require.NoError(t, err4)
	require.ElementsMatch(t, b, b2)

	err5 := o.Remove()
	require.NoError(t, err5)
}
