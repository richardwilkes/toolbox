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
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

var crashName = "test-crash"

func TestCrashOutputPathFor(t *testing.T) {
	c := check.New(t)
	prefix := filepath.Join("/logs", crashName)
	c.Equal(prefix+".log", xslog.CrashOutputPathFor(prefix, 0))
	c.Equal(prefix+".log", xslog.CrashOutputPathFor(prefix, -1))
	c.Equal(prefix+"-1.log", xslog.CrashOutputPathFor(prefix, 1))
	c.Equal(prefix+"-3.log", xslog.CrashOutputPathFor(prefix, 3))
}

// writeCrashFile creates one of the crash output files with the requested size.
func writeCrashFile(t *testing.T, base string, n, size int) string {
	t.Helper()
	path := xslog.CrashOutputPathFor(base, n)
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatalf("unable to write %s: %v", path, err)
	}
	return path
}

// crashFileSize returns the size of one of the crash output files, or -1 if it does not exist.
func crashFileSize(t *testing.T, base string, n int) int64 {
	t.Helper()
	fi, err := os.Stat(xslog.CrashOutputPathFor(base, n))
	if err != nil {
		if os.IsNotExist(err) {
			return -1
		}
		t.Fatalf("unable to stat: %v", err)
	}
	return fi.Size()
}

// TestRotateCrashOutput covers the bound on the crash output file. The Go runtime writes the crash report straight to
// the file descriptor as the process dies, so rotation can only happen when the file is opened at startup; without it
// the file would accumulate every crash report the installation ever produced.
func TestRotateCrashOutput(t *testing.T) {
	cfg := xslog.Rotator{MaxSize: 1000, MaxBackups: 2}
	cfg.Normalize()

	t.Run("no file is a no-op", func(t *testing.T) {
		c := check.New(t)
		base := filepath.Join(t.TempDir(), crashName)
		xslog.RotateCrashOutput(base, cfg)
		c.Equal(int64(-1), crashFileSize(t, base, 0))
		c.Equal(int64(-1), crashFileSize(t, base, 1))
	})

	t.Run("under the limit is left alone", func(t *testing.T) {
		c := check.New(t)
		base := filepath.Join(t.TempDir(), crashName)
		writeCrashFile(t, base, 0, 999)
		xslog.RotateCrashOutput(base, cfg)
		c.Equal(int64(999), crashFileSize(t, base, 0), "the file should not have been rotated")
		c.Equal(int64(-1), crashFileSize(t, base, 1))
	})

	t.Run("at the limit rotates", func(t *testing.T) {
		c := check.New(t)
		base := filepath.Join(t.TempDir(), crashName)
		writeCrashFile(t, base, 0, 1000)
		xslog.RotateCrashOutput(base, cfg)
		c.Equal(int64(-1), crashFileSize(t, base, 0), "the crash file should have been moved aside")
		c.Equal(int64(1000), crashFileSize(t, base, 1), "its content should be retained as the first backup")
	})

	t.Run("backups shift and the oldest is dropped", func(t *testing.T) {
		c := check.New(t)
		base := filepath.Join(t.TempDir(), crashName)
		writeCrashFile(t, base, 0, 1000)
		writeCrashFile(t, base, 1, 11)
		writeCrashFile(t, base, 2, 22)
		xslog.RotateCrashOutput(base, cfg)
		c.Equal(int64(-1), crashFileSize(t, base, 0))
		c.Equal(int64(1000), crashFileSize(t, base, 1))
		c.Equal(int64(11), crashFileSize(t, base, 2), "the first backup should have become the second")
		c.Equal(int64(-1), crashFileSize(t, base, 3), "retention must not exceed MaxBackups")
	})

	t.Run("backups beyond the retention limit are pruned", func(t *testing.T) {
		c := check.New(t)
		base := filepath.Join(t.TempDir(), crashName)
		writeCrashFile(t, base, 0, 1000)
		writeCrashFile(t, base, 1, 11)
		writeCrashFile(t, base, 2, 22)
		// Left behind by a run configured to retain more than this one does.
		writeCrashFile(t, base, 3, 33)
		writeCrashFile(t, base, 4, 44)
		xslog.RotateCrashOutput(base, cfg)
		c.Equal(int64(-1), crashFileSize(t, base, 3), "a stale backup should have been pruned")
		c.Equal(int64(-1), crashFileSize(t, base, 4), "the sweep should continue past the first stale backup")
	})
}

// TestCaptureCrashOutputRotatesAndOpens exercises the whole path: the crash file is rotated when it is already at the
// limit, and a fresh one is left in place for the runtime to append to.
func TestCaptureCrashOutputRotatesAndOpens(t *testing.T) {
	c := check.New(t)
	dir := t.TempDir()
	cfg := xslog.Rotator{Path: filepath.Join(dir, "test"+xslog.LogFileExt), MaxSize: 1000, MaxBackups: 1}
	cfg.Normalize()
	base := filepath.Join(dir, crashName)
	writeCrashFile(t, base, 0, 1000)

	xslog.CaptureCrashOutput(cfg)
	// Stop directing crash output at the temporary directory once this test is done with it.
	t.Cleanup(func() { c.NoError(debug.SetCrashOutput(nil, debug.CrashOptions{})) })

	c.Equal(int64(0), crashFileSize(t, base, 0), "a fresh, empty crash file should be ready for the runtime")
	c.Equal(int64(1000), crashFileSize(t, base, 1), "the previous content should have been retained as a backup")
	// The crash file must sit beside the log rather than replacing it.
	c.Equal(dir, filepath.Dir(xslog.CrashOutputPathFor(base, 0)))
}
