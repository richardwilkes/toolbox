// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog

import (
	"io"
	"log/slog"
	"os"

	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xio/term"
)

// SetupStd adds command-line options for controlling logging, parses the command line, instantiates a PrettyHandler and
// a log Rotator, then attaches it to slog. Returns the remaining arguments that weren't used for option content. An
// option to also send logs to the console will be added to the command line flags.
//
// Note that this function sends colored output to the log file if os.Stdout supports colors.
func SetupStd(cl *cmdline.CmdLine) (logFilePath string, remainingArgs []string) {
	var logLevel LevelValue
	logLevel.AddStdCmdLineOptions(cl)

	var rotatorCfg Rotator
	rotatorCfg.AddStdCmdLineOptions(cl)

	console := false
	opt := cl.NewGeneralOption(&console)
	opt.SetSingle('C').SetName("log-to-console").SetUsage(i18n.Text("Copy the log output to the console"))

	remainingArgs = cl.Parse(os.Args[1:])

	w := io.Writer(rotatorCfg.NewWriteCloser())
	if console {
		w = io.MultiWriter(w, os.Stdout)
	}
	slog.SetDefault(slog.New(NewPrettyHandler(w, &PrettyOptions{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		},
		ColorSupportOverride: term.DetectKind(os.Stdout),
	})))

	return rotatorCfg.Path, remainingArgs
}

// SetupStdToConsole adds command-line options for controlling logging, parses the command line, instantiates a
// PrettyHandler and a log Rotator, then attaches it to slog. Returns the remaining arguments that weren't used for
// option content. The logs will go to both the console and a file by default, but an option to suppress the console
// output will be added to the command line flags.
//
// Note that this function sends colored output to the log file if os.Stdout supports colors.
func SetupStdToConsole(cl *cmdline.CmdLine) (logFilePath string, remainingArgs []string) {
	var logLevel LevelValue
	logLevel.AddStdCmdLineOptions(cl)

	var rotatorCfg Rotator
	rotatorCfg.AddStdCmdLineOptions(cl)

	var quiet bool
	opt := cl.NewGeneralOption(&quiet)
	opt.SetSingle('q').SetName("quiet").SetUsage(i18n.Text("Suppress the log output to the console"))

	remainingArgs = cl.Parse(os.Args[1:])

	w := io.Writer(rotatorCfg.NewWriteCloser())
	if !quiet {
		w = io.MultiWriter(w, os.Stdout)
	}
	slog.SetDefault(slog.New(NewPrettyHandler(w, &PrettyOptions{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		},
		ColorSupportOverride: term.DetectKind(os.Stdout),
	})))

	return rotatorCfg.Path, remainingArgs
}
