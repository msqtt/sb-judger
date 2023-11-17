package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// Linux Overlayfs implements rootfs.
type Overlayfs struct {
	lowerPath string
	upperPath string
	workPath  string
	mntPath   string
	rootPath  string
}

// Delete implements Rootfs.
func (fs *Overlayfs) Delete(path string, rec bool) error {
	path = filepath.Join(fs.upperPath, path)
	if rec {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

// Move implements Rootfs.
func (fs *Overlayfs) Move(srcPath, destPath string) error {
	dest := filepath.Join(fs.upperPath, destPath)
	s := strings.Split(dest, "/")
	dirPath := filepath.Join(s[:len(s)-1]...)
  os.MkdirAll("/"+dirPath, 0755)
	return os.Rename(srcPath, dest)
}

// PutFile implements Rootfs.
func (fs *Overlayfs) PutFile(path, fileName string, file io.Reader) error {
	loc := filepath.Join(fs.upperPath, path)
	_, err := os.Stat(loc)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("cannot identify given path: [%w]", err)
		}
	}
	err = os.MkdirAll(loc, 0755)
	if err != nil {
		return fmt.Errorf("cannot mkdir dir: [%w]", err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read from file: [%w]", err)
	}
	err = os.WriteFile(filepath.Join(loc, fileName), b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file: [%w]", err)
	}
	return nil
}

var _ Rootfs = (*Overlayfs)(nil)

// NewOverlayFS returns an OverlayFS struct with a specific rootfs path.
func NewOverlayfs(rootfsPath string) (*Overlayfs, error) {
	// check rootfs
	fi, err := os.Stat(rootfsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot found rootfs: [%w]", err)
		}
		return nil, fmt.Errorf("cannot identify given rootfs: [%w]", err)
	}
	if !fi.IsDir() {
		return nil, errors.New("given roots is not a dir!")
	}
	return &Overlayfs{
		lowerPath: rootfsPath,
	}, nil
}

// Make Overlay Rootfs in dest path.
// returns an error
func (fs *Overlayfs) Make(dest string) error {
	// makedir
	folders := []string{"upper", "worker", "mnt"}
	for i, folder := range folders {
		folders[i] = filepath.Join(dest, folder)
		if err := syscall.Mkdir(folders[i], 0755); err != nil {
			return fmt.Errorf("cannot mkdir %s folder: [%w]", folders[i], err)
		}
	}
	fs.upperPath = folders[0]
	fs.workPath = folders[1]
	fs.mntPath = folders[2]
	fs.rootPath = dest

	option := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		fs.lowerPath,
		fs.upperPath,
		fs.workPath,
	)
	if err := syscall.Mount("overlay", fs.mntPath,
		"overlay", syscall.MS_NODEV, option); err != nil {
		return fmt.Errorf("cannot mount overlayfs: [%w]", err)
	}
	return nil
}

// MountPoint returns mount path.
func (fs *Overlayfs) MountPoint() (string, error) {
	if fs.mntPath == "" {
		return "", errors.New("this rootfs have not been made before!")
	}
	return fs.mntPath, nil
}

// Remove will umount the rootfs and delete the paths.
func (fs *Overlayfs) Remove() error {
	mnt, err := fs.MountPoint()
	if err != nil {
		return fmt.Errorf("cannot get mount point: [%w]", err)
	}
	err = syscall.Unmount(mnt, syscall.MNT_DETACH)
	if err != nil {
		return fmt.Errorf("cannot unmount rootfs: [%w]", err)
	}
	if err := os.RemoveAll(fs.rootPath); err != nil {
		return fmt.Errorf("cannot remove all folders: [%w]", err)
	}
	return nil
}
