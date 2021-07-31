// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rotation

import (
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
)

// Constants for defaults.
const (
	DefaultMaxSize    = 10 * 1024 * 1024
	DefaultMaxBackups = 1
)

// DefaultPath returns the default path that will be used. This will use
// cmdline.AppIdentifier (if set) to better isolate the log location.
func DefaultPath() string {
	return filepath.Join(paths.AppLogDir(), cmdline.AppCmdName+".log")
}

// Path specifies the file to write logs to. Backup log files will be retained
// in the same directory. Defaults to the value of DefaultPath().
func Path(path string) func(*Rotator) error {
	return func(r *Rotator) error {
		if path == "" {
			return errs.New("Must specify a path")
		}
		r.path = path
		return nil
	}
}

// MaxSize sets the maximum size of the log file before it gets rotated.
// Defaults to DefaultMaxSize.
func MaxSize(maxSize int64) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxSize = maxSize
		return nil
	}
}

// MaxBackups sets the maximum number of old log files to retain.  Defaults
// to DefaultMaxBackups.
func MaxBackups(maxBackups int) func(*Rotator) error {
	return func(r *Rotator) error {
		r.maxBackups = maxBackups
		return nil
	}
}

// WithMask sets the mask when creating files, which have the unmasked mode of 0644, and directories, which have the
// unmasked mode of 0755. Defaults to 0777.
func WithMask(mask os.FileMode) func(*Rotator) error {
	return func(r *Rotator) error {
		r.mask = mask
		return nil
	}
}
