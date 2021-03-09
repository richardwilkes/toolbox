// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/xio"
)

// Copy src to dst. src may be a directory, file, or symlink.
func Copy(src, dst string) error {
	return CopyWithMask(src, dst, 0777) //nolint:gocritic // File modes are octal
}

// CopyWithMask src to dst. src may be a directory, file, or symlink.
func CopyWithMask(src, dst string, mask os.FileMode) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return generalCopy(src, dst, info, mask)
}

func generalCopy(src, dst string, info os.FileInfo, mask os.FileMode) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return linkCopy(src, dst, info)
	}
	if info.IsDir() {
		return dirCopy(src, dst, info, mask)
	}
	return fileCopy(src, dst, info, mask)
}

func fileCopy(src, dst string, info os.FileInfo, mask os.FileMode) (err error) {
	if err = os.MkdirAll(filepath.Dir(dst), 0755&mask); err != nil {
		return err
	}
	var f *os.File
	if f, err = os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644&mask); err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
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

func dirCopy(srcDir, dstDir string, info os.FileInfo, mask os.FileMode) error {
	if err := os.MkdirAll(dstDir, info.Mode()&mask); err != nil {
		return err
	}
	list, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, one := range list {
		name := one.Name()
		if err = generalCopy(filepath.Join(srcDir, name), filepath.Join(dstDir, name), one, mask); err != nil {
			return err
		}
	}
	return nil
}

func linkCopy(src, dst string, info os.FileInfo) error {
	s, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(s, dst)
}
