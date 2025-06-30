// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/xio/fs/paths"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestRotator(t *testing.T) {
	const (
		maxSize    = 100
		maxBackups = 2
	)

	tmpdir := t.TempDir()

	logFiles := []string{filepath.Join(tmpdir, "test")}
	for i := range maxBackups {
		logFiles = append(logFiles, filepath.Join(tmpdir, fmt.Sprintf("test-%d", i+1)))
	}

	cfg := xslog.Rotator{
		Path:       logFiles[0],
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
	}
	r := cfg.NewWriteCloser()
	c := check.New(t)
	c.NotNil(r)
	_, err := os.Stat(logFiles[0] + xslog.LogFileExt)
	c.HasError(err)
	c.True(os.IsNotExist(err))
	for i := range maxSize * (2 + maxBackups) {
		_, err = fmt.Fprintln(r, i)
		c.NoError(err)
	}
	_, err = fmt.Fprintln(r, "goodbye")
	c.NoError(err)
	c.NoError(r.Close())
	for _, f := range logFiles {
		fi, fErr := os.Stat(f + xslog.LogFileExt)
		c.NoError(fErr)
		c.True(fi.Size() <= maxSize)
	}

	// Verify that we can start again
	r = cfg.NewWriteCloser()
	c.NotNil(r)
	_, err = fmt.Fprintln(r, "hello")
	c.NoError(err)
	c.NoError(r.Close())
}

func TestRotatorDefaults(t *testing.T) {
	var r xslog.Rotator
	r.Normalize()
	c := check.New(t)
	c.Equal(filepath.Join(paths.AppLogDir(), cmdline.AppCmdName+xslog.LogFileExt), r.Path)
	c.Equal(int64(10*1024*1024), r.MaxSize) //
	c.Equal(1, r.MaxBackups)
	c.Equal(os.FileMode(0o644), r.FileMode)
	c.Equal(os.FileMode(0o755), r.DirMode)
}

func TestRotatorWithNilConfig(t *testing.T) {
	c := check.New(t)
	c.NotNil(((*xslog.Rotator)(nil)).NewWriteCloser())
}

func TestRotatorCmdLineOpts(t *testing.T) {
	cl := cmdline.New(false)
	var r xslog.Rotator
	r.AddStdCmdLineOptions(cl)
	if os.Getenv("ROTATOR_CMDLINE_TEST") == "1" {
		// This is the subprocess
		cl.DisplayUsage()
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestRotatorCmdLineOpts")
	cmd.Env = append(os.Environ(), "ROTATOR_CMDLINE_TEST=1")
	output, err := cmd.CombinedOutput()
	c := check.New(t)
	c.NoError(err)
	c.Contains(string(output), "--log-file <value>")
	c.Contains(string(output), "--log-file-backups <value>")
	c.Contains(string(output), "--log-file-size <value>")
}
