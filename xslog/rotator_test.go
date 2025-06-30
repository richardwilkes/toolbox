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
	check.NotNil(t, r)
	_, err := os.Stat(logFiles[0] + xslog.LogFileExt)
	check.Error(t, err)
	check.True(t, os.IsNotExist(err))
	for i := range maxSize * (2 + maxBackups) {
		_, err = fmt.Fprintln(r, i)
		check.NoError(t, err)
	}
	_, err = fmt.Fprintln(r, "goodbye")
	check.NoError(t, err)
	check.NoError(t, r.Close())
	for _, f := range logFiles {
		fi, fErr := os.Stat(f + xslog.LogFileExt)
		check.NoError(t, fErr)
		check.True(t, fi.Size() <= maxSize)
	}

	// Verify that we can start again
	r = cfg.NewWriteCloser()
	check.NotNil(t, r)
	_, err = fmt.Fprintln(r, "hello")
	check.NoError(t, err)
	check.NoError(t, r.Close())
}

func TestRotatorDefaults(t *testing.T) {
	var r xslog.Rotator
	r.Normalize()
	check.Equal(t, filepath.Join(paths.AppLogDir(), cmdline.AppCmdName+xslog.LogFileExt), r.Path)
	check.Equal(t, int64(10*1024*1024), r.MaxSize) //
	check.Equal(t, 1, r.MaxBackups)
	check.Equal(t, os.FileMode(0o644), r.FileMode)
	check.Equal(t, os.FileMode(0o755), r.DirMode)
}

func TestRotatorWithNilConfig(t *testing.T) {
	check.NotNil(t, ((*xslog.Rotator)(nil)).NewWriteCloser())
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
	check.NoError(t, err)
	check.Contains(t, string(output), "--log-file <value>")
	check.Contains(t, string(output), "--log-file-backups <value>")
	check.Contains(t, string(output), "--log-file-size <value>")
}
