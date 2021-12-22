// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package zip provides simple zip extraction.
package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// ExtractArchive extracts the contents of a zip archive at 'src' into the 'dst' directory.
func ExtractArchive(src, dst string) error {
	return ExtractArchiveWithMask(src, dst, 0o777)
}

// ExtractArchiveWithMask extracts the contents of a zip archive at 'src' into the 'dst' directory.
func ExtractArchiveWithMask(src, dst string, mask os.FileMode) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	return ExtractWithMask(&r.Reader, dst, mask)
}

// Extract the contents of a zip reader into the 'dst' directory.
func Extract(zr *zip.Reader, dst string) error {
	return ExtractWithMask(zr, dst, 0o777)
}

// ExtractWithMask the contents of a zip reader into the 'dst' directory.
func ExtractWithMask(zr *zip.Reader, dst string, mask os.FileMode) error {
	root, err := filepath.Abs(dst)
	if err != nil {
		return errs.Wrap(err)
	}
	rootWithTrailingSep := fmt.Sprintf("%s%c", root, filepath.Separator)
	for _, f := range zr.File {
		path := filepath.Join(root, f.Name)
		if !strings.HasPrefix(path, rootWithTrailingSep) {
			return errs.Newf("Path outside of root is not permitted: %s", f.Name)
		}
		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode&os.ModeSymlink != 0:
			if err = extractSymLink(f, path, mask); err != nil {
				return err
			}
		case fi.IsDir():
			if err = os.MkdirAll(path, mode.Perm()&mask); err != nil {
				return errs.Wrap(err)
			}
		default:
			if err = extractFile(f, path, mask); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractSymLink(f *zip.File, dst string, mask os.FileMode) error {
	r, err := f.Open()
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	var buffer []byte
	if buffer, err = ioutil.ReadAll(r); err != nil {
		return errs.Wrap(err)
	}
	if err = os.MkdirAll(filepath.Dir(dst), 0o755&mask); err != nil {
		return errs.Wrap(err)
	}
	if err = os.Symlink(string(buffer), dst); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func extractFile(f *zip.File, dst string, mask os.FileMode) (err error) {
	var r io.ReadCloser
	if r, err = f.Open(); err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	if err = os.MkdirAll(filepath.Dir(dst), 0o755&mask); err != nil {
		return errs.Wrap(err)
	}
	var file *os.File
	if file, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.FileInfo().Mode().Perm()&mask); err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = errs.Wrap(closeErr)
		}
	}()
	if _, err = io.Copy(file, r); err != nil {
		err = errs.Wrap(err)
	}
	return
}
