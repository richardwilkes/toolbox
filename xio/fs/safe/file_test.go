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

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xio/fs/safe"
)

func TestAbortNonExisting(t *testing.T) {
	c := check.New(t)
	tmpdir := c.TempDir()
	filename := filepath.Join(tmpdir, "abort.txt")
	f, err := safe.CreateWithMode(filename, 0o600)
	c.NoError(err)
	var n int
	n, err = f.WriteString("abort")
	c.NoError(err)
	c.Equal(5, n)
	c.NoError(f.Close())
	_, err = os.Stat(filename)
	c.HasError(err)
}

func TestCommitNonExisting(t *testing.T) {
	c := check.New(t)
	tmpdir := c.TempDir()
	filename := filepath.Join(tmpdir, "commit.txt")
	f, err := safe.CreateWithMode(filename, 0o600)
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

func TestAbortExisting(t *testing.T) {
	c := check.New(t)
	tmpdir := c.TempDir()
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	c.NoError(os.WriteFile(filename, originalData, 0o600))
	f, err := safe.CreateWithMode(filename, 0o600)
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

func TestCommitExisting(t *testing.T) {
	c := check.New(t)
	tmpdir := c.TempDir()
	filename := filepath.Join(tmpdir, "safe.txt")
	originalData := []byte("safe")
	replacement := []byte("replaced")
	c.NoError(os.WriteFile(filename, originalData, 0o600))
	f, err := safe.CreateWithMode(filename, 0o600)
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
