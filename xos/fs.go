// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

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

// FileIsReadable returns true if the path points to a regular file whose owner-read permission bit is set.
func FileIsReadable(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular() && mode.Perm()&0o400 != 0
	}
	return false
}

// MoveFile moves src to dst, using rename if possible and falling back to copy-then-remove otherwise. It errors if
// src, or an existing dst, is not a regular file. If src and dst are the same file, nothing is done.
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

// CopyWithMask copies src to dst, masking the permission bits of everything it creates with mask. src may be a
// directory, file, or symlink.
func CopyWithMask(src, dst string, mask fs.FileMode) error {
	info, err := os.Lstat(src)
	if err != nil {
		return errs.Wrap(err)
	}
	return generalCopy(src, dst, info.Mode(), mask)
}

func generalCopy(src, dst string, srcMode, mask fs.FileMode) error {
	if srcMode&os.ModeSymlink != 0 {
		return linkCopy(src, dst, mask)
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
		return errs.Wrap(err)
	}
	// Registered first so it runs last, after the defer below has closed f: on failure, remove the destination rather
	// than leave a truncated, empty, or incorrectly-permissioned file behind.
	defer func() {
		if err != nil {
			_ = os.Remove(dst) //nolint:errcheck // best-effort cleanup; the original error is what matters
		}
	}()
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = errs.Wrap(closeErr)
		}
		if err == nil {
			if (srcMode&mask)|0o200 != srcMode&mask {
				err = chmodIfSupported(dst, srcMode&mask)
			}
		}
	}()
	var s *os.File
	if s, err = os.Open(src); err != nil {
		err = errs.Wrap(err)
		return err
	}
	if _, err = io.Copy(f, s); err != nil {
		err = errs.Wrap(err)
	}
	xio.CloseIgnoringErrors(s)
	return err
}

func dirCopy(srcDir, dstDir string, srcMode, mask fs.FileMode) (err error) {
	dstMode := srcMode & mask
	// Force owner rwx while populating the directory so children can be created even when the mask clears the owner's
	// write or execute bits, then restore the intended mode, as fileCopy does with owner write.
	if err = os.MkdirAll(dstDir, dstMode|0o700); err != nil {
		return errs.Wrap(err)
	}
	if dstMode|0o700 != dstMode {
		defer func() {
			if chmodErr := chmodIfSupported(dstDir, dstMode); chmodErr != nil && err == nil {
				err = chmodErr
			}
		}()
	}
	var list []os.DirEntry
	if list, err = os.ReadDir(srcDir); err != nil {
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

func linkCopy(src, dst string, mask fs.FileMode) error {
	s, err := os.Readlink(src)
	if err != nil {
		return errs.Wrap(err)
	}
	if err = os.MkdirAll(filepath.Dir(dst), 0o755&mask); err != nil {
		return errs.Wrap(err)
	}
	// os.Symlink fails with EEXIST if dst exists, so remove it first to match the overwrite behavior of the regular-file
	// path (which opens with O_TRUNC).
	if err = os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return errs.Wrap(err)
	}
	if err = os.Symlink(s, dst); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// chmodIfSupported changes the mode of the named path, treating an "operation not supported" error as success, since
// some filesystems (e.g. certain SMB/CIFS mounts) don't support chmod and that shouldn't abort an otherwise-valid copy.
func chmodIfSupported(name string, mode fs.FileMode) error {
	if err := os.Chmod(name, mode); err != nil && !errors.Is(err, syscall.ENOTSUP) &&
		!errors.Is(err, syscall.EOPNOTSUPP) {
		return errs.Wrap(err)
	}
	return nil
}
