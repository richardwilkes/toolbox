// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package jotrotate provides a pre-canned way to add jot logging with file
// rotation, along with command-line options for controlling it.
package jotrotate

import (
	"io"
	"os"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/rotation"
	"github.com/richardwilkes/toolbox/xio"
)

// ParseAndSetup adds command-line options for controlling logging, parses the
// command line, then instantiates a rotator and attaches it to jot. Returns
// the remaining arguments that weren't used for option content.
func ParseAndSetup(cl *cmdline.CmdLine) []string {
	logFile := rotation.DefaultPath()
	var maxSize int64 = rotation.DefaultMaxSize
	maxBackups := rotation.DefaultMaxBackups
	logToConsole := false
	cl.NewStringOption(&logFile).SetSingle('l').SetName("log-file").SetUsage("The file to write logs to")
	cl.NewInt64Option(&maxSize).SetName("log-file-size").SetUsage("The maximum number of bytes to write to a log file before rotating it")
	cl.NewIntOption(&maxBackups).SetName("log-file-backups").SetUsage("The maximum number of old logs files to retain")
	cl.NewBoolOption(&logToConsole).SetSingle('C').SetName("log-to-console").SetUsage("Copy the log output to the console")
	remainingArgs := cl.Parse(os.Args[1:])
	if rotator, err := rotation.New(rotation.Path(logFile), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups)); err == nil {
		if logToConsole {
			jot.SetWriter(&xio.TeeWriter{Writers: []io.Writer{rotator, os.Stdout}})
		} else {
			jot.SetWriter(rotator)
		}
	} else {
		jot.Error(err)
	}
	return remainingArgs
}
