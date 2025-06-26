// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rotation

import (
	"io"
	"log"
	"os"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
)

// PathToLog holds the path to the log file that was configured on the command line when using ParseAndSetup().
var PathToLog string

// ParseAndSetupLogging adds command-line options for controlling logging, parses the command line, then instantiates a
// rotator and attaches it to slog. Returns the remaining arguments that weren't used for option content. If
// consoleOnByDefault is true, then logs will also go to the console by default, but an option to turn them off will be
// added to the command line flags. Conversely, if it is false, an option to turn them on will be added to the command
// line flags.
func ParseAndSetupLogging(cl *cmdline.CmdLine, consoleOnByDefault bool) []string {
	logFile := DefaultPath()
	var maxSize int64 = DefaultMaxSize
	maxBackups := DefaultMaxBackups
	consoleOption := false
	cl.NewGeneralOption(&logFile).SetSingle('l').SetName("log-file").SetUsage(i18n.Text("The file to write logs to"))
	cl.NewGeneralOption(&maxSize).SetName("log-file-size").SetUsage(i18n.Text("The maximum number of bytes to write to a log file before rotating it"))
	cl.NewGeneralOption(&maxBackups).SetName("log-file-backups").SetUsage(i18n.Text("The maximum number of old logs files to retain"))
	opt := cl.NewGeneralOption(&consoleOption)
	if consoleOnByDefault {
		opt.SetSingle('q').SetName("quiet").SetUsage(i18n.Text("Suppress the log output to the console"))
	} else {
		opt.SetSingle('C').SetName("log-to-console").SetUsage(i18n.Text("Copy the log output to the console"))
	}
	remainingArgs := cl.Parse(os.Args[1:])
	if rotator, err := New(Path(logFile), MaxSize(maxSize), MaxBackups(maxBackups)); err == nil {
		if consoleOnByDefault == consoleOption {
			log.SetOutput(rotator)
		} else {
			log.SetOutput(&xio.TeeWriter{Writers: []io.Writer{rotator, os.Stdout}})
		}
		PathToLog = rotator.PathToLog()
	} else {
		errs.Log(err)
	}
	return remainingArgs
}
