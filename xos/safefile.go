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
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// SafeFile provides safe overwriting of files. Instead of truncating and overwriting the destination file, it creates a
// temporary file in the same directory, writes to it, then renames the temporary file to the original name when
// Commit() is called. If Close() is called without calling Commit(), or the Commit() fails, then the original file is
// left untouched.
type SafeFile struct {
	*os.File
	name      string
	committed bool
	closed    bool
}

// WriteSafeFile creates a SafeFile, calls 'writer' to write data into it, then commits it.
func WriteSafeFile(filename string, writer func(io.Writer) error) (err error) {
	var f *SafeFile
	f, err = CreateSafeFile(filename)
	if err != nil {
		return
	}
	w := bufio.NewWriterSize(f, 1<<16)
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if err = writer(w); err != nil {
		return
	}
	if err = w.Flush(); err != nil {
		return
	}
	if err = f.Commit(); err != nil {
		return
	}
	return
}

// CreateSafeFile creates a temporary file in the same directory as filename, which will be renamed to the given
// filename when calling Commit.
func CreateSafeFile(filename string) (*SafeFile, error) {
	filename = filepath.Clean(filename)
	if filename == "" || filename[len(filename)-1] == filepath.Separator {
		return nil, os.ErrInvalid
	}
	f, err := os.CreateTemp(filepath.Dir(filename), "safe")
	if err != nil {
		return nil, err
	}
	return &SafeFile{
		File: f,
		name: filename,
	}, nil
}

// OriginalName returns the original filename passed into CreateSafeFile().
func (f *SafeFile) OriginalName() string {
	return f.name
}

// Commit the data into the original file and remove the temporary file from disk. Close() may still be called, but will
// do nothing.
func (f *SafeFile) Commit() error {
	if f.committed {
		return nil
	}
	if f.closed {
		return os.ErrInvalid
	}
	f.committed = true
	f.closed = true
	var err error
	name := f.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(name) //nolint:errcheck // no need to report this error, too
		}
	}()
	if err = f.File.Close(); err != nil {
		return errs.Wrap(err)
	}
	if err = os.Rename(name, f.name); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// Close the temporary file and remove it, if it hasn't already been committed. If it has been committed, nothing
// happens.
func (f *SafeFile) Close() error {
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
