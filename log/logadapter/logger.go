// Package logadapter defines an API to use for logging, which actual logging
// implementations can implement directly or provide an adapter to use.
package logadapter

import "time"

// DebugLogger defines an API to use for logging debugging messages, which
// actual logging implementations can implement directly or provide an
// adapter to use.
type DebugLogger interface {
	// Debug logs a debugging message. Arguments are handled in the manner of
	// fmt.Print.
	Debug(v ...interface{})
	// Debugf logs a debugging message. Arguments are handled in the manner of
	// fmt.Printf.
	Debugf(format string, v ...interface{})
}

// InfoLogger defines an API to use for logging informational messages, which
// actual logging implementations can implement directly or provide an
// adapter to use.
type InfoLogger interface {
	// Info logs an informational message. Arguments are handled in the manner
	// of fmt.Print.
	Info(v ...interface{})
	// Infof logs an informational message. Arguments are handled in the
	// manner of fmt.Print.
	Infof(format string, v ...interface{})
}

// WarnLogger defines an API to use for logging warning messages, which actual
// logging implementations can implement directly or provide an adapter to
// use.
type WarnLogger interface {
	// Warn logs a warning message. Arguments are handled in the manner of
	// fmt.Print.
	Warn(v ...interface{})
	// Warnf logs a warning message. Arguments are handled in the manner of
	// fmt.Printf.
	Warnf(format string, v ...interface{})
}

// ErrorLogger defines an API to use for logging error messages, which actual
// logging implementations can implement directly or provide an adapter to
// use.
type ErrorLogger interface {
	// Error logs an error message. Arguments are handled in the manner of
	// fmt.Print.
	Error(v ...interface{})
	// Errorf logs an error message. Arguments are handled in the manner of
	// fmt.Printf.
	Errorf(format string, v ...interface{})
}

// FatalLogger defines an API to use for logging fatal error messages, which
// actual logging implementations can implement directly or provide an
// adapter to use.
type FatalLogger interface {
	// Fatal logs a fatal error message. Arguments other than the status are
	// handled in the manner of fmt.Print.
	Fatal(status int, v ...interface{})
	// Fatalf logs a fatal error message. Arguments other than the status are
	// handled in the manner of fmt.Printf.
	Fatalf(status int, format string, v ...interface{})
}

// Timing is used to record the duration between two events. One of End(),
// EndWithMsg(), or EndWithMsgf() should be called when the event has
// finished.
type Timing interface {
	// End finishes timing an event and logs an informational message.
	End() time.Duration
	// EndWithMsg finishes timing an event and logs an informational message.
	// Arguments are handled in the manner of fmt.Print.
	EndWithMsg(v ...interface{}) time.Duration
	// EndWithMsgf finishes timing an event and logs an informational message.
	// Arguments are handled in the manner of fmt.Printf.
	EndWithMsgf(format string, v ...interface{}) time.Duration
}

// TimingLogger defines an API to use for logging timed data, which actual
// logging implementations can implement directly or provide an adapter to
// use.
type TimingLogger interface {
	// Time starts timing an event and logs an informational message.
	// Arguments are handled in the manner of fmt.Print.
	Time(v ...interface{}) Timing
	// Timef starts timing an event and logs an informational message.
	// Arguments are handled in the manner of fmt.Printf.
	Timef(format string, v ...interface{}) Timing
}

// Logger defines an API to use for logging, which actual logging
// implementations can implement directly or provide an adapter to use.
type Logger interface {
	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
	FatalLogger
	TimingLogger
}
