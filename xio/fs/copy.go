// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package fs provides filesystem-related utilities.
package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/xio"
)

// Copy src to dst. src may be a directory, file, or symlink.
func Copy(src, dst string) error {
	return CopyWithMask(src, dst, 0o777)
}

// CopyWithMask src to dst. src may be a directory, file, or symlink.
func CopyWithMask(src, dst string, mask fs.FileMode) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
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
		return err
	}
	var f *os.File
	if f, err = os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, (srcMode&mask)|0o200); err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	var s *os.File
	if s, err = os.Open(src); err != nil {
		return err
	}
	defer xio.CloseIgnoringErrors(s)
	_, err = io.Copy(f, s)
	return err
}

func dirCopy(srcDir, dstDir string, srcMode, mask fs.FileMode) error {
	if err := os.MkdirAll(dstDir, srcMode&mask); err != nil {
		return err
	}
	list, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, one := range list {
		name := one.Name()
		if err = generalCopy(filepath.Join(srcDir, name), filepath.Join(dstDir, name), one.Type(), mask); err != nil {
			return err
		}
	}
	return nil
}

func linkCopy(src, dst string) error {
	s, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(s, dst)
}
