// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package jot

import (
	"io"

	"github.com/richardwilkes/toolbox/log/logadapter"
)

// Logger wraps the various jot function calls into a struct that can be
// passed around, typically for the sake of satisfying one or more logging
// interfaces.
type Logger struct {
}

// SetWriter sets the io.Writer to use when writing log messages. Default is
// os.Stderr.
func (lgr *Logger) SetWriter(w io.Writer) {
	SetWriter(w)
}

// SetMinimumLevel sets the minimum log level that will be output. Default is
// DEBUG.
func (lgr *Logger) SetMinimumLevel(level Level) {
	SetMinimumLevel(level)
}

// Debug logs a debug message. Arguments are handled in the manner of
// fmt.Print.
func (lgr *Logger) Debug(v ...interface{}) {
	Debug(v...)
}

// Debugf logs a debug message. Arguments are handled in the manner of
// fmt.Printf.
func (lgr *Logger) Debugf(format string, v ...interface{}) {
	Debugf(format, v...)
}

// Info logs an informational message. Arguments are handled in the manner of
// fmt.Print.
func (lgr *Logger) Info(v ...interface{}) {
	Info(v...)
}

// Infof logs an informational message. Arguments are handled in the manner of
// fmt.Printf.
func (lgr *Logger) Infof(format string, v ...interface{}) {
	Infof(format, v...)
}

// Warn logs a warning message. Arguments are handled in the manner of
// fmt.Print.
func (lgr *Logger) Warn(v ...interface{}) {
	Warn(v...)
}

// Warnf logs a warning message. Arguments are handled in the manner of
// fmt.Printf.
func (lgr *Logger) Warnf(format string, v ...interface{}) {
	Warnf(format, v...)
}

// Error logs an error message. Arguments are handled in the manner of
// fmt.Print.
func (lgr *Logger) Error(v ...interface{}) {
	Error(v...)
}

// Errorf logs an error message. Arguments are handled in the manner of
// fmt.Printf.
func (lgr *Logger) Errorf(format string, v ...interface{}) {
	Errorf(format, v...)
}

// Fatal logs a fatal error message. Arguments other than the status are
// handled in the manner of fmt.Print.
func (lgr *Logger) Fatal(status int, v ...interface{}) {
	Fatal(status, v...)
}

// Fatalf logs a fatal error message. Arguments other than the status are
// handled in the manner of fmt.Printf.
func (lgr *Logger) Fatalf(status int, format string, v ...interface{}) {
	Fatalf(status, format, v...)
}

// Time starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Print.
func (lgr *Logger) Time(v ...interface{}) logadapter.Timing {
	return Time(v...)
}

// Timef starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Printf.
func (lgr *Logger) Timef(format string, v ...interface{}) logadapter.Timing {
	return Timef(format, v...)
}

// Flush waits for all current log entries to be written before returning.
func (lgr *Logger) Flush() {
	Flush()
}

// Writer logs the data as an error after casting it to a string.
func (lgr *Logger) Write(data []byte) (int, error) {
	Error(string(data))
	return len(data), nil
}
