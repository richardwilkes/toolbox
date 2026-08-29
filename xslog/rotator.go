// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// LogFileExt is the extension used for log files.
const LogFileExt = ".log"

var _ io.WriteCloser = &rotatingWriter{}

// Rotator configures the size-rotated log file io.WriteCloser created by NewWriteCloser. AddFlags adds command-line
// options for these settings.
type Rotator struct {
	// Path is the log file to write to; LogFileExt is appended if not already present. Backups are kept in the same
	// directory. Leave empty to use the default log path.
	Path string
	// MaxSize is the size in bytes at which the log file is rotated. Defaults to 10 MiB.
	MaxSize int64
	// MaxBackups is the maximum number of old log files to retain. Defaults to 1.
	MaxBackups int
	// DirMode is the permission bits used when creating directories. Defaults to 0o755.
	DirMode os.FileMode
	// FileMode is the permission bits used when creating files. Defaults to 0o644.
	FileMode os.FileMode
}

// NewWriteCloser returns an io.WriteCloser that creates the log file on first write and rotates it when a write would
// reach MaxSize. The receiver is normalized first; a nil receiver uses the defaults.
func (r *Rotator) NewWriteCloser() io.WriteCloser {
	var w rotatingWriter
	if r != nil {
		r.Normalize()
		w.cfg = *r
	} else {
		w.cfg.Normalize()
	}
	w.cfg.Path = strings.TrimSuffix(w.cfg.Path, LogFileExt)
	return &w
}

// AddFlags adds command-line flags for controlling log rotation.
func (r *Rotator) AddFlags() {
	r.Normalize()
	flag.StringVar(&r.Path, "log-file", r.Path, i18n.Text("The `file` to write logs to"))
	flag.Int64Var(&r.MaxSize, "log-file-size", r.MaxSize,
		i18n.Text("The maximum number of `bytes` to write to a log file before rotating it"))
	flag.IntVar(&r.MaxBackups, "log-file-backups", r.MaxBackups,
		i18n.Text("The maximum `number` of old logs files to retain"))
}

// Normalize fills in defaults for any unset fields. Calling it is optional, but useful for determining the default
// values.
func (r *Rotator) Normalize() {
	if r.Path == "" {
		r.Path = filepath.Join(xos.AppLogDir(true), xos.AppCmdName+LogFileExt)
	}
	if r.MaxSize <= 0 {
		r.MaxSize = 10 * 1024 * 1024
	}
	if r.MaxBackups <= 0 {
		r.MaxBackups = 1
	}
	if r.DirMode == 0 {
		r.DirMode = 0o755
	}
	if r.FileMode == 0 {
		r.FileMode = 0o644
	}
}

type rotatingWriter struct {
	file *os.File
	cfg  Rotator
	size int64
	lock sync.Mutex
}

// Write implements io.Writer.
func (r *rotatingWriter) Write(b []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
retry:
	if r.file == nil {
		if err := os.MkdirAll(filepath.Dir(r.cfg.Path), r.cfg.DirMode); err != nil {
			return 0, errs.Wrap(err)
		}
		p := r.pathFor(0)
		if fi, err := os.Stat(p); err == nil {
			r.size = fi.Size()
		} else {
			r.size = 0
		}
		file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, r.cfg.FileMode)
		if err != nil {
			return 0, errs.Wrap(err)
		}
		r.file = file
	}
	// Only rotate a file that already has content. A single write >= MaxSize into an empty file must still be written;
	// rotating would reset the size to 0 and loop forever without making progress.
	if r.size > 0 && r.size+int64(len(b)) >= r.cfg.MaxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
		goto retry
	}
	n, err := r.file.Write(b)
	if err != nil {
		err = errs.Wrap(err)
	}
	r.size += int64(n)
	return n, err
}

// Close implements io.Closer.
func (r *rotatingWriter) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.file == nil {
		return nil
	}
	file := r.file
	r.file = nil
	return errs.Wrap(file.Close())
}

func (r *rotatingWriter) rotate() error {
	if r.file != nil {
		err := r.file.Close()
		r.file = nil
		if err != nil {
			return errs.Wrap(err)
		}
	}
	// Remove the oldest backup and any higher-numbered ones left by a previous run with a larger MaxBackups. Backups
	// are consecutive, so the first missing slot ends the sweep.
	for n := r.cfg.MaxBackups; ; n++ {
		if err := os.Remove(r.pathFor(n)); err != nil {
			if os.IsNotExist(err) {
				break
			}
			return errs.Wrap(err)
		}
	}
	for i := r.cfg.MaxBackups; i > 0; i-- {
		if err := os.Rename(r.pathFor(i-1), r.pathFor(i)); err != nil && !os.IsNotExist(err) {
			return errs.Wrap(err)
		}
	}
	r.file = nil
	r.size = 0
	return nil
}

func (r *rotatingWriter) pathFor(n int) string {
	if n <= 0 {
		return r.cfg.Path + LogFileExt
	}
	return fmt.Sprintf("%s-%d%s", r.cfg.Path, n, LogFileExt)
}
