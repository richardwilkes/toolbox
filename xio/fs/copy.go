// Package fs provides filesystem-related utilities.
package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/xio"
)

// Copy src to dst. src may be a directory, file, or symlink.
func Copy(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copy(src, dst, info)
}

func copy(src, dst string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return linkCopy(src, dst, info)
	}
	if info.IsDir() {
		return dirCopy(src, dst, info)
	}
	return fileCopy(src, dst, info)
}

func fileCopy(src, dst string, info os.FileInfo) (err error) {
	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	var f *os.File
	if f, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		if lerr := f.Close(); lerr != nil && err == nil {
			err = lerr
		}
	}()
	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}
	var s *os.File
	if s, err = os.Open(src); err != nil {
		return err
	}
	defer xio.CloseIgnoringErrors(s)
	_, err = io.Copy(f, s)
	return err
}

func dirCopy(srcdir, dstdir string, info os.FileInfo) error {
	if err := os.MkdirAll(dstdir, info.Mode()); err != nil {
		return err
	}
	list, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}
	for _, one := range list {
		name := one.Name()
		if err := copy(filepath.Join(srcdir, name), filepath.Join(dstdir, name), one); err != nil {
			return err
		}
	}
	return nil
}

func linkCopy(src, dst string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dst)
}
