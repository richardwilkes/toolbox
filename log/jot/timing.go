package jot

import (
	"fmt"
	"time"
)

// Timing is used to record the duration between two events.
type Timing struct {
	started time.Time
	msg     string
}

// End should be called when the event has finished.
func (timing *Timing) End() time.Duration {
	elapsed := time.Since(timing.started)
	Infof("Finished %s | %v elapsed", timing.msg, elapsed)
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
