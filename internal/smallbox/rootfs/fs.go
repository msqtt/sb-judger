package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// Rootfs supports some functions to control the virt filesystem.
type Rootfs interface {
	Make(dest string) error
	Remove() error
	MountPoint() (string, error)
	// put file to a named path, if path do not exist, it will be created.
	PutFile(path, fileName string, file io.Reader) error
	// move srcPath to destPath.
	Move(srcPath, destPath string) error
	// delete path, if enable recursion.
	Delete(path string, rec bool) error
}

// ChrootMaskPath will mount rootfs, chroot and mask some paths.
// that will be dangerous if executing on local machine directly.
// ChrootMaskPath should be executed after do a mount namespace for process.
func ChrootMaskPath(mnt string) error {
	// private mount
	if err := syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return fmt.Errorf("cannot set private mount: [%w]", err)
	}
	// pivot root
	oldRootfs := filepath.Join(mnt, ".pivot_root")
	if err := os.Mkdir(oldRootfs, 0755); err != nil {
		return fmt.Errorf("cannot mkdir for oldrootfs: [%w]", err)
	}
	if err := syscall.Mount(mnt, mnt, "bind",
		syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("cannot mount pivotPath itself: [%w]", err)
	}
	if err := syscall.PivotRoot(mnt, oldRootfs); err != nil {
		return fmt.Errorf("cannot pivot to %s from %s: [%w]", mnt, oldRootfs, err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("cannot chdir to root: [%w]", err)
	}
	if err := syscall.Unmount(".pivot_root", syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount old pivotRootPath: [%w]", err)
	}
	if err := os.Remove(".pivot_root"); err != nil {
		return fmt.Errorf("remove old pivotRootPath: [%w]", err)
	}
	// mask proc
	if err := syscall.Mount("proc", "/proc", "proc",
		syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, ""); err != nil {
		return fmt.Errorf("cannot mount proc in the container: [%w]", err)
	}
	return nil
}
