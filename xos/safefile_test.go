// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestCreateSafeFile_AbortNonExisting(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "abort.txt")
	f, err := xos.CreateSafeFile(filename)
	c.NoError(err)
	var n int
	n, err = f.WriteString("abort")
	c.NoError(err)
	c.Equal(5, n)
	c.NoError(f.Close())
	_, err = os.Stat(filename)
	c.HasError(err)
}

func TestCreateSafeFile_CommitNonExisting(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "commit.txt")
	f, err := xos.CreateSafeFile(filename)
	c.NoError(err)
	var n int
	n, err = f.WriteString("commit")
	c.NoError(err)
	c.Equal(6, n)
	c.NoError(f.Commit())
	c.NoError(f.Close())
	_, err = os.Stat(filename)
	c.NoError(err)
	c.NoError(os.Remove(filename))
}

func TestCreateSafeFile_AbortExisting(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "safe.txt")
	originalData := []byte("safe")
	c.NoError(os.WriteFile(filename, originalData, 0o640))
	f, err := xos.CreateSafeFile(filename)
	c.NoError(err)
	var n int
	n, err = f.WriteString("bad")
	c.NoError(err)
	c.Equal(3, n)
	err = f.Close()
	c.NoError(err)
	var data []byte
	data, err = os.ReadFile(filename)
	c.NoError(err)
	c.Equal(originalData, data)
	c.NoError(os.Remove(filename))
}

func TestCreateSafeFile_CommitExisting(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "safe.txt")
	originalData := []byte("safe")
	replacement := []byte("replaced")
	c.NoError(os.WriteFile(filename, originalData, 0o640))
	f, err := xos.CreateSafeFile(filename)
	c.NoError(err)
	var n int
	n, err = f.Write(replacement)
	c.NoError(err)
	c.Equal(len(replacement), n)
	c.NoError(f.Commit())
	c.NoError(f.Close())
	var data []byte
	data, err = os.ReadFile(filename)
	c.NoError(err)
	c.Equal(replacement, data)
	c.NoError(os.Remove(filename))
}

func TestWriteSafeFile_Success(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "writesafe.txt")
	testData := "test data for WriteSafeFile"
	err := xos.WriteSafeFile(filename, func(w io.Writer) error {
		_, writeErr := w.Write([]byte(testData))
		return writeErr
	})
	c.NoError(err)
	data, err := os.ReadFile(filename)
	c.NoError(err)
	c.Equal(testData, string(data))
	c.NoError(os.Remove(filename))
}

func TestWriteSafeFile_WriterError(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "writesafe_error.txt")
	expectedErr := errors.New("writer error")
	err := xos.WriteSafeFile(filename, func(_ io.Writer) error {
		return expectedErr
	})
	c.HasError(err)
	c.Equal(expectedErr, err)
	_, err = os.Stat(filename)
	c.HasError(err)
}

func TestWriteSafeFile_ReplaceExisting(t *testing.T) {
	c := check.New(t)
	filename := filepath.Join(c.TempDir(), "replace.txt")
	originalData := "original content"
	c.NoError(os.WriteFile(filename, []byte(originalData), 0o640))
	newData := "new content"
	err := xos.WriteSafeFile(filename, func(w io.Writer) error {
		_, writeErr := w.Write([]byte(newData))
		return writeErr
	})
	c.NoError(err)
	data, err := os.ReadFile(filename)
	c.NoError(err)
	c.Equal(newData, string(data))
	c.NoError(os.Remove(filename))
}

func TestWriteSafeFile_InvalidFilename(t *testing.T) {
	c := check.New(t)
	err := xos.WriteSafeFile("", func(_ io.Writer) error {
		return nil
	})
	c.HasError(err)
}
