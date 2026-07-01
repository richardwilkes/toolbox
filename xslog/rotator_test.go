// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestRotator(t *testing.T) {
	const (
		maxSize    = 100
		maxBackups = 2
	)

	tmpdir := t.TempDir()

	logFiles := append(make([]string, 0, 1+maxBackups), filepath.Join(tmpdir, "test"))
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

func TestRotatorPrunesStaleBackupsOnShrink(t *testing.T) {
	const maxSize = 100
	tmpdir := t.TempDir()
	path := filepath.Join(tmpdir, "test")
	c := check.New(t)

	// Simulate a previous run that retained more backups than the current configuration allows.
	const priorBackups = 5
	for i := 1; i <= priorBackups; i++ {
		c.NoError(os.WriteFile(fmt.Sprintf("%s-%d%s", path, i, xslog.LogFileExt), []byte("old"), 0o644))
	}

	// Run with a smaller MaxBackups and write enough to force at least one rotation.
	const maxBackups = 2
	cfg := xslog.Rotator{
		Path:       path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
	}
	r := cfg.NewWriteCloser()
	c.NotNil(r)
	for i := range maxSize * (2 + maxBackups) {
		_, err := fmt.Fprintln(r, i)
		c.NoError(err)
	}
	c.NoError(r.Close())

	// The retained backups (1..maxBackups) must exist and the stale higher-numbered ones must be gone.
	for i := 1; i <= priorBackups; i++ {
		_, err := os.Stat(fmt.Sprintf("%s-%d%s", path, i, xslog.LogFileExt))
		if i <= maxBackups {
			c.NoError(err)
		} else {
			c.True(os.IsNotExist(err))
		}
	}
}

func TestRotatorOversizedWrite(t *testing.T) {
	const maxSize = 100
	tmpdir := t.TempDir()
	path := filepath.Join(tmpdir, "test")
	cfg := xslog.Rotator{
		Path:       path,
		MaxSize:    maxSize,
		MaxBackups: 2,
	}
	r := cfg.NewWriteCloser()
	c := check.New(t)
	c.NotNil(r)

	// A single write larger than MaxSize must still succeed and not loop forever. Run it under a watchdog so a
	// regression of the infinite-loop bug fails fast instead of hanging the suite.
	big := bytes.Repeat([]byte("a"), maxSize*3)
	var n int
	done := make(chan error, 1)
	go func() {
		var wErr error
		n, wErr = r.Write(big)
		done <- wErr
	}()
	select {
	case wErr := <-done:
		c.NoError(wErr)
		c.Equal(len(big), n)
	case <-time.After(10 * time.Second):
		t.Fatal("Write of a record larger than MaxSize did not return; rotator is looping")
	}
	c.NoError(r.Close())

	// The oversized record should have been written in full to the primary log file.
	fi, err := os.Stat(path + xslog.LogFileExt)
	c.NoError(err)
	c.Equal(int64(len(big)), fi.Size())

	// A subsequent normal write must still rotate that oversized file out of the way rather than appending to it.
	r = cfg.NewWriteCloser()
	c.NotNil(r)
	_, err = fmt.Fprintln(r, "small")
	c.NoError(err)
	c.NoError(r.Close())
	fi, err = os.Stat(path + xslog.LogFileExt)
	c.NoError(err)
	c.True(fi.Size() <= maxSize)
}

func TestRotatorDefaults(t *testing.T) {
	var r xslog.Rotator
	r.Normalize()
	c := check.New(t)
	c.Equal(filepath.Join(xos.AppLogDir(true), xos.AppCmdName+xslog.LogFileExt), r.Path)
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
