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

// Time starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Print.
func Time(v ...interface{}) logadapter.Timing {
	msg := fmt.Sprint(v...)
	Infof("Starting %s", msg)
	return &timing{
		started: time.Now(),
		msg:     msg,
	}
}

// Timef starts timing an event and logs an informational message. Arguments are handled in the manner of fmt.Printf.
func Timef(format string, v ...interface{}) logadapter.Timing {
	msg := fmt.Sprintf(format, v...)
	Infof("Starting %s", msg)
	return &timing{
		started: time.Now(),
		msg:     msg,
	}
}
