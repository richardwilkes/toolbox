// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package safe provides safe, atomic saving of files.
package safe

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/xio"
)

// File provides safe, atomic saving of files. Instead of truncating and overwriting the destination file, it creates a
// temporary file in the same directory, writes to it, and then renames the temporary file to the original name when
// Commit() is called. If Close() is called without calling Commit(), or the Commit() fails, then the original file is
// left untouched.
type File struct {
	*os.File
	originalName string
	committed    bool
	closed       bool
}

// Create creates a temporary file in the same directory as filename, which will be renamed to the given filename when
// calling Commit.
func Create(filename string) (*File, error) {
	return CreateWithMode(filename, 0o644)
}

// CreateWithMode creates a temporary file in the same directory as filename, which will be renamed to the given
// filename when calling Commit.
func CreateWithMode(filename string, mode os.FileMode) (*File, error) {
	filename = filepath.Clean(filename)
	if filename == "" || filename[len(filename)-1] == filepath.Separator {
		return nil, os.ErrInvalid
	}
	f, err := ioutil.TempFile(filepath.Dir(filename), "safe")
	if err != nil {
		return nil, err
	}
	if runtime.GOOS != toolbox.WindowsOS { // Windows doesn't support changing the mode
		if err = f.Chmod(mode); err != nil {
			xio.CloseIgnoringErrors(f)
			_ = os.Remove(f.Name()) //nolint:errcheck // no need to report this error, too
			return nil, err
		}
	}
	return &File{
		File:         f,
		originalName: filename,
	}, nil
}

// OriginalName returns the original filename passed into Create().
func (f *File) OriginalName() string {
	return f.originalName
}

// Commit the data into the original file and remove the temporary file from disk. Close() may still be called, but will
// do nothing.
func (f *File) Commit() error {
	if f.committed {
		return nil
	}
	if f.closed {
		return os.ErrInvalid
	}
	f.committed = true
	f.closed = true
	err := f.Sync()
	if closeErr := f.File.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	name := f.Name()
	if err == nil {
		err = os.Rename(name, f.originalName)
	}
	if err != nil {
		_ = os.Remove(name) //nolint:errcheck // no need to report this error, too
	}
	return err
}

// Close the temporary file and remove it, if it hasn't already been committed. If it has been committed, nothing
// happens.
func (f *File) Close() error {
	if f.committed {
		return nil
	}
	if f.closed {
		return os.ErrInvalid
	}
	f.closed = true
	err := f.File.Close()
	if removeErr := os.Remove(f.Name()); removeErr != nil && err == nil {
		err = removeErr
	}
	return err
}
