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
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// CaptureCrashOutput arranges for the Go runtime's crash report -- the message and goroutine dump an unrecovered panic
// or a fatal runtime error produces -- to be written to a file beside the log.
//
// That report otherwise only goes to standard error, which is discarded for an application launched from the Finder or
// its equivalent, and the log receives nothing because the process dies before any of it reaches the logging handler.
// The result is a crash report from a user that no one can act on. The runtime writes here in addition to standard
// error rather than instead of it, so running from a terminal is unaffected.
//
// The file is appended to, so a crash is still recoverable if the application is relaunched before anyone collects it,
// and it is only ever written when the process actually dies, so it stays absent or small in normal use. It is rotated
// on the same terms as the log; see RotateCrashOutput for why that happens here rather than as it is written.
func CaptureCrashOutput(cfg Rotator) {
	cfg.Normalize()
	if cfg.Path == "" {
		return
	}
	// Both the file and its backups are named from this base, matching how the log's rotating writer derives its own.
	base := strings.TrimSuffix(cfg.Path, LogFileExt) + "-crash"
	path := base + LogFileExt
	if err := os.MkdirAll(filepath.Dir(path), cfg.DirMode); err != nil {
		errs.Log(errs.NewWithCause("unable to create the directory for the crash output file", err), "path", path)
		return
	}
	RotateCrashOutput(base, cfg)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, cfg.FileMode)
	if err != nil {
		errs.Log(errs.NewWithCause("unable to create the crash output file", err), "path", path)
		return
	}
	// SetCrashOutput duplicates the descriptor, so ours is no longer needed once it returns.
	defer func() {
		if err = f.Close(); err != nil {
			errs.Log(errs.NewWithCause("unable to close the crash output file", err), "path", path)
		}
	}()
	if err = debug.SetCrashOutput(f, debug.CrashOptions{}); err != nil {
		errs.Log(errs.NewWithCause("unable to direct crash output to a file", err), "path", path)
	}
}

// RotateCrashOutput retires the crash output file, and prunes its backups, once it has reached the size the log's
// rotation is configured with, so that repeated crashes cannot grow it without bound.
//
// This runs when the file is opened rather than when it is written, because there is no opportunity to run Go code at
// the point it is written: the runtime writes the crash report directly to the file descriptor while the process is
// dying. Rotating here is enough, since a single run can only ever append one report -- the process does not survive
// it -- so the file can only exceed the limit between runs, which is exactly when this is called.
func RotateCrashOutput(base string, cfg Rotator) {
	fi, err := os.Stat(CrashOutputPathFor(base, 0))
	if err != nil || fi.Size() < cfg.MaxSize {
		return
	}
	// Remove the oldest retained backup along with any higher-numbered ones left behind by a previous run configured to
	// retain more. Backups are always consecutive, so the first missing slot ends the sweep.
	for n := cfg.MaxBackups; ; n++ {
		if err = os.Remove(CrashOutputPathFor(base, n)); err != nil {
			if !os.IsNotExist(err) {
				errs.Log(errs.NewWithCause("unable to prune a crash output backup", err), "path",
					CrashOutputPathFor(base, n))
			}
			break
		}
	}
	for i := cfg.MaxBackups; i > 0; i-- {
		if err = os.Rename(CrashOutputPathFor(base, i-1), CrashOutputPathFor(base, i)); err != nil &&
			!os.IsNotExist(err) {
			errs.Log(errs.NewWithCause("unable to rotate the crash output file", err), "path",
				CrashOutputPathFor(base, i-1))
			return
		}
	}
}

// CrashOutputPathFor returns the crash output file itself for n <= 0, or its n-th backup, using the same naming the
// log's rotating writer uses for its own backups.
func CrashOutputPathFor(base string, n int) string {
	if n <= 0 {
		return base + LogFileExt
	}
	return fmt.Sprintf("%s-%d%s", base, n, LogFileExt)
}
