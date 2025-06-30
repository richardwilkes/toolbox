// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xio/fs/paths"
)

// LogFileExt is the extension used for log files.
const LogFileExt = ".log"

var _ io.WriteCloser = &rotatingWriter{}

// Rotator provides configuration for creating a io.WriteCloser via a call to its NewWriteCloser() method that will
// rotate a log file when it exceeds a certain size. You can also use the method AddStdCmdLineOptions() to add
// command-line options to control the log rotation.
type Rotator struct {
	// Path specifies the file to write logs to. Backup log files will be retained in the same directory. Leave empty to
	// use the default log path.
	Path string
	// MaxSize sets the maximum size of the log file before it gets rotated. Defaults to 10 MiB.
	MaxSize int64
	// MaxBackups sets the maximum number of old log files to retain. Defaults to 1.
	MaxBackups int
	// DirMode sets the permission bits to use when creating directories. Defaults to 0o755.
	DirMode os.FileMode
	// FileMode sets the permission bits to use when creating files. Defaults to 0o644.
	FileMode os.FileMode
}

// NewWriteCloser creates a new io.WriteCloser that will write to the log file specified in the configuration. It will
// create the file if it does not exist when needed, and will rotate the log file when it exceeds the maximum size.
func (c *Rotator) NewWriteCloser() io.WriteCloser {
	var r rotatingWriter
	if c != nil {
		c.Normalize()
		r.cfg = *c
	} else {
		r.cfg.Normalize()
	}
	r.cfg.Path = strings.TrimSuffix(r.cfg.Path, LogFileExt)
	return &r
}

// AddStdCmdLineOptions adds the standard command-line options for controlling log rotation.
func (c *Rotator) AddStdCmdLineOptions(cl *cmdline.CmdLine) {
	c.Normalize()
	cl.NewGeneralOption(&c.Path).SetName("log-file").SetUsage(i18n.Text("The file to write logs to"))
	cl.NewGeneralOption(&c.MaxSize).SetName("log-file-size").
		SetUsage(i18n.Text("The maximum number of bytes to write to a log file before rotating it"))
	cl.NewGeneralOption(&c.MaxBackups).SetName("log-file-backups").
		SetUsage(i18n.Text("The maximum number of old logs files to retain"))
}

// Normalize ensures that the configuration is valid. It sets defaults for any fields that are not set. It is not
// necessary to call this, but might be useful if you want to programmatically determine the default values.
func (c *Rotator) Normalize() {
	if c.Path == "" {
		c.Path = filepath.Join(paths.AppLogDir(), cmdline.AppCmdName+LogFileExt)
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 10 * 1024 * 1024
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 1
	}
	if c.DirMode == 0 {
		c.DirMode = 0o755
	}
	if c.FileMode == 0 {
		c.FileMode = 0o644
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
	if r.size+int64(len(b)) >= r.cfg.MaxSize {
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
	if err := os.Remove(r.pathFor(r.cfg.MaxBackups)); err != nil && !os.IsNotExist(err) {
		return errs.Wrap(err)
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
