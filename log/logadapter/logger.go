// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package logadapter defines an API to use for logging, which actual logging implementations can implement directly or
// provide an adapter to use.
package logadapter

import "time"

// DebugLogger defines an API to use for logging debugging messages, which actual logging implementations can implement
// directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type DebugLogger interface {
	// Debug logs a debugging message. Arguments are handled in the manner of fmt.Print.
	Debug(v ...any)
	// Debugf logs a debugging message. Arguments are handled in the manner of fmt.Printf.
	Debugf(format string, v ...any)
}

// InfoLogger defines an API to use for logging informational messages, which actual logging implementations can
// implement directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type InfoLogger interface {
	// Info logs an informational message. Arguments are handled in the manner of fmt.Print.
	Info(v ...any)
	// Infof logs an informational message. Arguments are handled in the manner of fmt.Print.
	Infof(format string, v ...any)
}

// WarnLogger defines an API to use for logging warning messages, which actual logging implementations can implement
// directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type WarnLogger interface {
	// Warn logs a warning message. Arguments are handled in the manner of fmt.Print.
	Warn(v ...any)
	// Warnf logs a warning message. Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, v ...any)
}

// ErrorLogger defines an API to use for logging error messages, which actual logging implementations can implement
// directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type ErrorLogger interface {
	// Error logs an error message. Arguments are handled in the manner of fmt.Print.
	Error(v ...any)
	// Errorf logs an error message. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, v ...any)
}

// FatalLogger defines an API to use for logging fatal error messages, which actual logging implementations can
// implement directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type FatalLogger interface {
	// Fatal logs a fatal error message. Arguments other than the status are handled in the manner of fmt.Print.
	Fatal(status int, v ...any)
	// Fatalf logs a fatal error message. Arguments other than the status are handled in the manner of fmt.Printf.
	Fatalf(status int, format string, v ...any)
}

// Timing is used to record the duration between two events. One of End(), EndWithMsg(), or EndWithMsgf() should be
// called when the event has finished.
//
// Deprecated: Use slog instead. August 28, 2023
type Timing interface {
	// End finishes timing an event and logs an informational message.
	End() time.Duration
	// EndWithMsg finishes timing an event and logs an informational message. Arguments are handled in the manner of
	// fmt.Print.
	EndWithMsg(v ...any) time.Duration
	// EndWithMsgf finishes timing an event and logs an informational message. Arguments are handled in the manner of
	// fmt.Printf.
	EndWithMsgf(format string, v ...any) time.Duration
}

// TimingLogger defines an API to use for logging timed data, which actual logging implementations can implement
// directly or provide an adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type TimingLogger interface {
	// Time starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Print.
	Time(v ...any) Timing
	// Timef starts timing an event and logs an informational message. Arguments are handled in the manner of
	// fmt.Printf.
	Timef(format string, v ...any) Timing
}

// Logger defines an API to use for logging, which actual logging implementations can implement directly or provide an
// adapter to use.
//
// Deprecated: Use slog instead. August 28, 2023
type Logger interface {
	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
	FatalLogger
	TimingLogger
}
