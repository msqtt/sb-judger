package fs

import (
	_ "io/ioutil"
	_ "os"
	_ "testing"

	_ "github.com/stretchr/testify/require"
)

const testRootfsPath = "../../../rootfs"
const testMntPath = "./test"

// well, it's pretty tough to test this funciton.
// todo
//
// func TestChrootMaskPath(t *testing.T) {
// 	overlay, err := NewOverlayfs(testRootfsPath)
// 	require.NoError(t, err)
// 	require.NotNil(t, overlay)
//
// 	dir := makeTestMnt(t, overlay)
// 	
// 	mnt, err := overlay.MountPoint()
// 	require.NoError(t, err)
//
// 	err = ChrootMaskPath(mnt)
// 	require.NoError(t, err)
//
// 	rootDir, err := ioutil.ReadDir(".")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, rootDir)
// 	require.ElementsMatch(t, rootDir, dir)
//
// 	err = overlay.Remove()
// 	require.NoError(t, err)
// }
