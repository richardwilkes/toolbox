// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio"
)

// IsDir returns true if the specified path exists and is a directory.
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// FileExists returns true if the path points to a regular file.
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular()
	}
	return false
}

// FileIsReadable returns true if the path points to a regular file that we have permission to read.
func FileIsReadable(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular() && mode.Perm()&0o400 != 0
	}
	return false
}

// MoveFile moves a file in the file system or across volumes, using rename if possible, but falling back to copying the
// file if not. This will error if either src or dst are not regular files.
func MoveFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errs.Wrap(err)
	}
	if !srcInfo.Mode().IsRegular() {
		return errs.Newf("%s is not a regular file", src)
	}
	var dstInfo os.FileInfo
	dstInfo, err = os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
	} else {
		if !dstInfo.Mode().IsRegular() {
			return errs.Newf("%s is not a regular file", dst)
		}
		if os.SameFile(srcInfo, dstInfo) {
			return nil
		}
	}
	if os.Rename(src, dst) == nil {
		return nil
	}
	if err = Copy(src, dst); err != nil {
		return err
	}
	if err = os.Remove(src); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// Copy src to dst. src may be a directory, file, or symlink.
func Copy(src, dst string) error {
	return CopyWithMask(src, dst, 0o777)
}

// CopyWithMask src to dst. src may be a directory, file, or symlink.
func CopyWithMask(src, dst string, mask fs.FileMode) error {
	info, err := os.Lstat(src)
	if err != nil {
		return errs.Wrap(err)
	}
	return generalCopy(src, dst, info.Mode(), mask)
}

func generalCopy(src, dst string, srcMode, mask fs.FileMode) error {
	if srcMode&os.ModeSymlink != 0 {
		return linkCopy(src, dst)
	}
	if srcMode.IsDir() {
		return dirCopy(src, dst, srcMode, mask)
	}
	return fileCopy(src, dst, srcMode, mask)
}

func fileCopy(src, dst string, srcMode, mask fs.FileMode) (err error) {
	if err = os.MkdirAll(filepath.Dir(dst), 0o755&mask); err != nil {
		return errs.Wrap(err)
	}
	var f *os.File
	if f, err = os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, (srcMode&mask)|0o200); err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = errs.Wrap(closeErr)
		}
		if err == nil {
			if (srcMode&mask)|0o200 != srcMode&mask {
				if err = os.Chmod(dst, srcMode&mask); err != nil {
					err = errs.Wrap(err)
				}
			}
		}
	}()
	var s *os.File
	if s, err = os.Open(src); err != nil {
		err = errs.Wrap(err)
		return
	}
	if _, err = io.Copy(f, s); err != nil {
		err = errs.Wrap(err)
	}
	xio.CloseIgnoringErrors(s)
	return
}

func dirCopy(srcDir, dstDir string, srcMode, mask fs.FileMode) error {
	if err := os.MkdirAll(dstDir, srcMode&mask); err != nil {
		return errs.Wrap(err)
	}
	list, err := os.ReadDir(srcDir)
	if err != nil {
		return errs.Wrap(err)
	}
	for _, one := range list {
		name := one.Name()
		var fi os.FileInfo
		if fi, err = one.Info(); err != nil {
			return errs.Wrap(err)
		}
		if err = generalCopy(filepath.Join(srcDir, name), filepath.Join(dstDir, name), fi.Mode(), mask); err != nil {
			return err
		}
	}
	return nil
}

func linkCopy(src, dst string) error {
	s, err := os.Readlink(src)
	if err != nil {
		return errs.Wrap(err)
	}
	if err = os.Symlink(s, dst); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
