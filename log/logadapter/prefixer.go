package logadapter

import "fmt"

// Prefixer adds a prefix to another logger's output.
type Prefixer struct {
	Logger Logger
	Prefix string
}

// Debug logs a debug message. Arguments are handled in the manner of
// fmt.Print.
func (p *Prefixer) Debug(v ...interface{}) {
	p.Logger.Debugf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Debugf logs a debug message. Arguments are handled in the manner of
// fmt.Printf.
func (p *Prefixer) Debugf(format string, v ...interface{}) {
	p.Logger.Debugf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Info logs an informational message. Arguments are handled in the manner of
// fmt.Print.
func (p *Prefixer) Info(v ...interface{}) {
	p.Logger.Infof("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Infof logs an informational message. Arguments are handled in the manner of
// fmt.Printf.
func (p *Prefixer) Infof(format string, v ...interface{}) {
	p.Logger.Infof("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Warn logs a warning message. Arguments are handled in the manner of
// fmt.Print.
func (p *Prefixer) Warn(v ...interface{}) {
	p.Logger.Warnf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Warnf logs a warning message. Arguments are handled in the manner of
// fmt.Printf.
func (p *Prefixer) Warnf(format string, v ...interface{}) {
	p.Logger.Warnf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Error logs an error message. Arguments are handled in the manner of
// fmt.Print.
func (p *Prefixer) Error(v ...interface{}) {
	p.Logger.Errorf("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Errorf logs an error message. Arguments are handled in the manner of
// fmt.Printf.
func (p *Prefixer) Errorf(format string, v ...interface{}) {
	p.Logger.Errorf("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Fatal logs a fatal error message. Arguments are handled in the manner of
// fmt.Print.
func (p *Prefixer) Fatal(status int, v ...interface{}) {
	p.Logger.Fatalf(status, "%s%s", p.Prefix, fmt.Sprint(v...))
}

// Fatalf logs a fatal error message. Arguments are handled in the manner of
// fmt.Printf.
func (p *Prefixer) Fatalf(status int, format string, v ...interface{}) {
	p.Logger.Fatalf(status, "%s%s", p.Prefix, fmt.Sprintf(format, v...))
}

// Time starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Print.
func (p *Prefixer) Time(v ...interface{}) Timing {
	return p.Logger.Timef("%s%s", p.Prefix, fmt.Sprint(v...))
}

// Timef starts timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Printf.
func (p *Prefixer) Timef(format string, v ...interface{}) Timing {
	return p.Logger.Timef("%s%s", p.Prefix, fmt.Sprintf(format, v...))
}
