package logadapter

import (
	"time"

	"github.com/richardwilkes/toolbox/atexit"
)

// Discarder discards all data given to it.
type Discarder struct {
}

// Debug logs a debug message. Arguments are handled in the manner of
// fmt.Print.
func (d *Discarder) Debug(v ...interface{}) {
}

// Debugf logs a debug message. Arguments are handled in the manner of
// fmt.Printf.
func (d *Discarder) Debugf(format string, v ...interface{}) {
}

// Info logs an informational message. Arguments are handled in the manner of
// fmt.Print.
func (d *Discarder) Info(v ...interface{}) {
}

// Infof logs an informational message. Arguments are handled in the manner of
// fmt.Printf.
func (d *Discarder) Infof(format string, v ...interface{}) {
}

// Warn logs a warning message. Arguments are handled in the manner of
// fmt.Print.
func (d *Discarder) Warn(v ...interface{}) {
}

// Warnf logs a warning message. Arguments are handled in the manner of
// fmt.Printf.
func (d *Discarder) Warnf(format string, v ...interface{}) {
}

// Error logs an error message. Arguments are handled in the manner of
// fmt.Print.
func (d *Discarder) Error(v ...interface{}) {
}

// Errorf logs an error message. Arguments are handled in the manner of
// fmt.Printf.
func (d *Discarder) Errorf(format string, v ...interface{}) {
}

// Fatal logs a fatal error message. Arguments are handled in the manner of
// fmt.Print.
func (d *Discarder) Fatal(status int, v ...interface{}) {
	atexit.Exit(status)
}

// Fatalf logs a fatal error message. Arguments are handled in the manner of
// fmt.Printf.
func (d *Discarder) Fatalf(status int, format string, v ...interface{}) {
	atexit.Exit(status)
}

type discarderTiming struct {
	started time.Time
}

func (d *discarderTiming) End() time.Duration {
	return time.Since(d.started)
}

func (d *discarderTiming) EndWithMsg(v ...interface{}) time.Duration {
	return time.Since(d.started)
}

func (d *discarderTiming) EndWithMsgf(format string, v ...interface{}) time.Duration {
	return time.Since(d.started)
}

// Time starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Print.
func (d *Discarder) Time(v ...interface{}) Timing {
	return &discarderTiming{started: time.Now()}
}

// Timef starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Printf.
func (d *Discarder) Timef(format string, v ...interface{}) Timing {
	return &discarderTiming{started: time.Now()}
}
