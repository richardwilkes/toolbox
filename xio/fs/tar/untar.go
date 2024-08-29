// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package tar provides simple tar extraction.
package tar

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
)

// ExtractArchive extracts the contents of a tar archive at 'src' into the 'dst' directory.
func ExtractArchive(src, dst string) error {
	return ExtractArchiveWithMask(src, dst, 0o777)
}

// ExtractArchiveWithMask extracts the contents of a tar archive at 'src' into the 'dst' directory.
func ExtractArchiveWithMask(src, dst string, mask os.FileMode) error {
	f, err := os.Open(src)
	if err != nil {
		return errs.Wrap(err)
	}
	r := tar.NewReader(f)
	defer xio.CloseIgnoringErrors(f)
	return ExtractWithMask(r, dst, mask)
}

// Extract the contents of a tar reader into the 'dst' directory.
func Extract(tr *tar.Reader, dst string) error {
	return ExtractWithMask(tr, dst, 0o777)
}

// ExtractWithMask the contents of a tar reader into the 'dst' directory.
func ExtractWithMask(tr *tar.Reader, dst string, mask os.FileMode) error {
	root, err := filepath.Abs(dst)
	if err != nil {
		return errs.Wrap(err)
	}
	rootWithTrailingSep := fmt.Sprintf("%s%c", root, filepath.Separator)
	for {
		var hdr *tar.Header
		if hdr, err = tr.Next(); errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return errs.Wrap(err)
		}
		path := filepath.Join(root, hdr.Name) //nolint:gosec // We check for path traversal below
		if !strings.HasPrefix(path, rootWithTrailingSep) {
			return errs.Newf("Path outside of root is not permitted: %s", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeReg:
			if err = extractFile(tr, path, hdr.FileInfo().Mode().Perm(), mask); err != nil {
				return err
			}
		case tar.TypeLink:
			if err = os.MkdirAll(filepath.Dir(path), 0o755&mask); err != nil {
				return errs.Wrap(err)
			}
			if err = os.Link(hdr.Linkname, path); err != nil {
				return errs.Wrap(err)
			}
		case tar.TypeSymlink:
			if err = os.MkdirAll(filepath.Dir(path), 0o755&mask); err != nil {
				return errs.Wrap(err)
			}
			if err = os.Symlink(hdr.Linkname, path); err != nil {
				return errs.Wrap(err)
			}
		case tar.TypeDir:
			if err = os.MkdirAll(path, hdr.FileInfo().Mode().Perm()&mask); err != nil {
				return errs.Wrap(err)
			}
		}
	}
}

func extractFile(r io.Reader, dst string, mode, mask os.FileMode) (err error) {
	if err = os.MkdirAll(filepath.Dir(dst), 0o755&mask); err != nil {
		return errs.Wrap(err)
	}
	var file *os.File
	if file, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode&mask); err != nil {
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
	return nil
}
