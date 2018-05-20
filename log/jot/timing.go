package jot

import (
	"fmt"
	"time"

	"github.com/richardwilkes/toolbox/log/logadapter"
)

type timing struct {
	started time.Time
	msg     string
}

func (t *timing) End() time.Duration {
	elapsed := time.Since(t.started)
	Infof("Finished %s | %v elapsed", t.msg, elapsed)
	return elapsed
}

func (t *timing) EndWithMsg(v ...interface{}) time.Duration {
	elapsed := time.Since(t.started)
	Infof("Finished %s | %s | %v elapsed", t.msg, fmt.Sprint(v...), elapsed)
	return elapsed
}

func (t *timing) EndWithMsgf(format string, v ...interface{}) time.Duration {
	elapsed := time.Since(t.started)
	Infof("Finished %s | %s | %v elapsed", t.msg, fmt.Sprintf(format, v...), elapsed)
	return elapsed
}

// Time starts timing an event and logs an informational message. Arguments
// are handled in the manner of fmt.Print.
func Time(v ...interface{}) logadapter.Timing {
	msg := fmt.Sprint(v...)
	Infof("Starting %s", msg)
	return &timing{
		started: time.Now(),
		msg:     msg,
	}
}

// Timef starts timing an event and logs an informational message. Arguments
// are handled in the manner of fmt.Printf.
func Timef(format string, v ...interface{}) logadapter.Timing {
	msg := fmt.Sprintf(format, v...)
	Infof("Starting %s", msg)
	return &timing{
		started: time.Now(),
		msg:     msg,
	}
}
