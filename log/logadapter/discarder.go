// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package logadapter

import (
	"time"

	"github.com/richardwilkes/toolbox/atexit"
)

// Discarder discards all data given to it.
type Discarder struct{}

// Debug logs a debug message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Debug(_ ...any) {
}

// Debugf logs a debug message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Debugf(_ string, _ ...any) {
}

// Info logs an informational message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Info(_ ...any) {
}

// Infof logs an informational message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Infof(_ string, _ ...any) {
}

// Warn logs a warning message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Warn(_ ...any) {
}

// Warnf logs a warning message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Warnf(_ string, _ ...any) {
}

// Error logs an error message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Error(_ ...any) {
}

// Errorf logs an error message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Errorf(_ string, _ ...any) {
}

// Fatal logs a fatal error message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Fatal(status int, _ ...any) {
	atexit.Exit(status)
}

// Fatalf logs a fatal error message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Fatalf(status int, _ string, _ ...any) {
	atexit.Exit(status)
}

type discarderTiming struct {
	started time.Time
}

func (d *discarderTiming) End() time.Duration {
	return time.Since(d.started)
}

func (d *discarderTiming) EndWithMsg(_ ...any) time.Duration {
	return time.Since(d.started)
}

func (d *discarderTiming) EndWithMsgf(_ string, _ ...any) time.Duration {
	return time.Since(d.started)
}

// Time starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Time(_ ...any) Timing {
	return &discarderTiming{started: time.Now()}
}

// Timef starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Timef(_ string, _ ...any) Timing {
	return &discarderTiming{started: time.Now()}
}
