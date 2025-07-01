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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xflag"
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
	c.Equal(filepath.Join(paths.AppLogDir(), xflag.AppCmdName+xslog.LogFileExt), r.Path)
	c.Equal(int64(10*1024*1024), r.MaxSize) //
	c.Equal(1, r.MaxBackups)
	c.Equal(os.FileMode(0o644), r.FileMode)
	c.Equal(os.FileMode(0o755), r.DirMode)
}

func TestRotatorWithNilConfig(t *testing.T) {
	c := check.New(t)
	c.NotNil(((*xslog.Rotator)(nil)).NewWriteCloser())
}

func TestRotatorAddFlags(t *testing.T) {
	var r xslog.Rotator
	r.AddFlags()
	hasFile := false
	hasBackups := false
	hasSize := false
	flag.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "log-file":
			hasFile = true
		case "log-file-backups":
			hasBackups = true
		case "log-file-size":
			hasSize = true
		}
	})
	c := check.New(t)
	c.True(hasFile)
	c.True(hasBackups)
	c.True(hasSize)
}
