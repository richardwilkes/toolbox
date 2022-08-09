// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/xio/fs/safe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbortNonExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	require.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "abort.txt")
	f, err := safe.CreateWithMode(filename, 0o600)
	require.NoError(t, err)
	n, err := f.WriteString("abort")
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	err = f.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filename)
	assert.Error(t, err)
}

func removeAll(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.RemoveAll(path))
}

func TestCommitNonExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	require.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "commit.txt")
	f, err := safe.CreateWithMode(filename, 0o600)
	require.NoError(t, err)
	n, err := f.WriteString("commit")
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	err = f.Commit()
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filename)
	assert.NoError(t, err)
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestAbortExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	require.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	require.NoError(t, os.WriteFile(filename, originalData, 0o600))
	f, err := safe.CreateWithMode(filename, 0o600)
	require.NoError(t, err)
	n, err := f.WriteString("bad")
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	err = f.Close()
	assert.NoError(t, err)
	data, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, originalData, data)
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestCommitExisting(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "safe_test_")
	require.NoError(t, err)
	defer removeAll(t, tmpdir)
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	replacement := []byte("replaced")
	require.NoError(t, os.WriteFile(filename, originalData, 0o600))
	f, err := safe.CreateWithMode(filename, 0o600)
	require.NoError(t, err)
	n, err := f.Write(replacement)
	assert.NoError(t, err)
	assert.Equal(t, len(replacement), n)
	err = f.Commit()
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	data, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, replacement, data)
	err = os.Remove(filename)
	assert.NoError(t, err)
}
