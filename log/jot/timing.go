package jot

import (
	"fmt"
	"time"
)

// Timing is used to record the duration between two events. One of End(),
// EndWithMsg(), or EndWithMsgf() should be called when the event has
// finished.
type Timing struct {
	started time.Time
	msg     string
}

// End finishes timing an event and logs an informational message.
func (timing *Timing) End() time.Duration {
	elapsed := time.Since(timing.started)
	Infof("Finished %s | %v elapsed", timing.msg, elapsed)
	return elapsed
}

// EndWithMsg finishes timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Print.
func (timing *Timing) EndWithMsg(v ...interface{}) time.Duration {
	elapsed := time.Since(timing.started)
	Infof("Finished %s | %s | %v elapsed", timing.msg, fmt.Sprint(v...), elapsed)
	return elapsed
}

// EndWithMsgf finishes timing an event and logs an informational message.
// Arguments are handled in the manner of fmt.Printf.
func (timing *Timing) EndWithMsgf(format string, v ...interface{}) time.Duration {
	elapsed := time.Since(timing.started)
	Infof("Finished %s | %s | %v elapsed", timing.msg, fmt.Sprintf(format, v...), elapsed)
	return elapsed
}

// Time starts timing an event and logs an informational message. Arguments
// are handled in the manner of fmt.Print.
func Time(v ...interface{}) *Timing {
	msg := fmt.Sprint(v...)
	Infof("Starting %s", msg)
	return &Timing{
		started: time.Now(),
		msg:     msg,
	}
}

// Timef starts timing an event and logs an informational message. Arguments
// are handled in the manner of fmt.Printf.
func Timef(format string, v ...interface{}) *Timing {
	msg := fmt.Sprintf(format, v...)
	Infof("Starting %s", msg)
	return &Timing{
		started: time.Now(),
		msg:     msg,
	}
}
