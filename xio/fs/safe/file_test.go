// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package safe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

func TestAbortNonExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	check.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "abort.txt")
	var f *safe.File
	f, err = safe.CreateWithMode(filename, 0o600)
	check.NoError(t, err)
	var n int
	n, err = f.WriteString("abort")
	check.NoError(t, err)
	check.Equal(t, 5, n)
	check.NoError(t, f.Close())
	_, err = os.Stat(filename)
	check.Error(t, err)
}

func removeAll(t *testing.T, path string) {
	t.Helper()
	check.NoError(t, os.RemoveAll(path))
}

func TestCommitNonExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	check.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "commit.txt")
	var f *safe.File
	f, err = safe.CreateWithMode(filename, 0o600)
	check.NoError(t, err)
	var n int
	n, err = f.WriteString("commit")
	check.NoError(t, err)
	check.Equal(t, 6, n)
	check.NoError(t, f.Commit())
	check.NoError(t, f.Close())
	_, err = os.Stat(filename)
	check.NoError(t, err)
	check.NoError(t, os.Remove(filename))
}

func TestAbortExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	check.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	check.NoError(t, os.WriteFile(filename, originalData, 0o600))
	var f *safe.File
	f, err = safe.CreateWithMode(filename, 0o600)
	check.NoError(t, err)
	var n int
	n, err = f.WriteString("bad")
	check.NoError(t, err)
	check.Equal(t, 3, n)
	err = f.Close()
	check.NoError(t, err)
	var data []byte
	data, err = os.ReadFile(filename)
	check.NoError(t, err)
	check.Equal(t, originalData, data)
	check.NoError(t, os.Remove(filename))
}

func TestCommitExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	check.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	replacement := []byte("replaced")
	check.NoError(t, os.WriteFile(filename, originalData, 0o600))
	var f *safe.File
	f, err = safe.CreateWithMode(filename, 0o600)
	check.NoError(t, err)
	var n int
	n, err = f.Write(replacement)
	check.NoError(t, err)
	check.Equal(t, len(replacement), n)
	check.NoError(t, f.Commit())
	check.NoError(t, f.Close())
	var data []byte
	data, err = os.ReadFile(filename)
	check.NoError(t, err)
	check.Equal(t, replacement, data)
	check.NoError(t, os.Remove(filename))
}
