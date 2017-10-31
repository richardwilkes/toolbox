package logadapter

import "github.com/richardwilkes/gokit/atexit"

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
