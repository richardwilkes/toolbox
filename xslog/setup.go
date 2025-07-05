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
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xflag"
	"github.com/richardwilkes/toolbox/v2/xterm"
)

// SetupStd adds command-line options for controlling logging (including one to also send logs to the console), parses
// the command line, instantiates a PrettyHandler and a log Rotator, then attaches it to slog. Returns the path to the
// log file.
//
// Note that this function sends colored output to the log file if os.Stdout supports colors.
func SetupStd(description, argsUsage string) string {
	var logLevel LevelValue
	logLevel.AddFlags()

	var rotatorCfg Rotator
	rotatorCfg.AddFlags()

	console := flag.Bool("console", false, i18n.Text("Copy the log output to the console"))
	xflag.SetUsage(description, argsUsage)
	flag.Parse()

	w := io.Writer(rotatorCfg.NewWriteCloser())
	if *console {
		w = io.MultiWriter(w, os.Stdout)
	}
	slog.SetDefault(slog.New(NewPrettyHandler(w, &PrettyOptions{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		},
		ColorSupportOverride: xterm.DetectKind(os.Stdout),
	})))

	return rotatorCfg.Path
}

// SetupStdToConsole adds command-line options for controlling logging (including one to suppress the console output),
// parses the command line, instantiates a PrettyHandler and a log Rotator, then attaches it to slog. Returns the path
// to the log file.
//
// Note that this function sends colored output to the log file if os.Stdout supports colors.
func SetupStdToConsole(description, argsUsage string) string {
	var logLevel LevelValue
	logLevel.AddFlags()

	var rotatorCfg Rotator
	rotatorCfg.AddFlags()

	quiet := false
	quietUsage := i18n.Text("Suppress the log output to the console")
	flag.BoolVar(&quiet, "quiet", quiet, quietUsage)
	flag.BoolVar(&quiet, "q", quiet, quietUsage)
	xflag.SetUsage(description, argsUsage)
	flag.Parse()

	w := io.Writer(rotatorCfg.NewWriteCloser())
	if !quiet {
		w = io.MultiWriter(w, os.Stdout)
	}
	slog.SetDefault(slog.New(NewPrettyHandler(w, &PrettyOptions{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		},
		ColorSupportOverride: xterm.DetectKind(os.Stdout),
	})))

	return rotatorCfg.Path
}
