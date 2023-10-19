package fs

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
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
