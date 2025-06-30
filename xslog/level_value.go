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
	"log/slog"
	"strings"

	"github.com/richardwilkes/toolbox/v2/cmdline"
	"github.com/richardwilkes/toolbox/v2/i18n"
)

var (
	_ slog.Leveler  = LevelValue{}
	_ cmdline.Value = &LevelValue{} // For command line parsing
)

// LevelValue is a cmdline.Value that holds a slog.Level and can be used to set the log level from the command line.
type LevelValue struct {
	level slog.Level
}

// Level implements the slog.Leveler interface.
func (v LevelValue) Level() slog.Level {
	return v.level
}

// Set implements the cmdline.Value interface.
func (v *LevelValue) Set(value string) error {
	return v.level.UnmarshalText([]byte(value))
}

func (v LevelValue) String() string {
	return v.level.String()
}

// AddStdCmdLineOptions adds the standard command-line options for controlling the log level.
func (v *LevelValue) AddStdCmdLineOptions(cl *cmdline.CmdLine) {
	cl.NewOption(v).SetName("log-level").
		SetUsage(i18n.Text("The level of logging to use. Valid values are: ") +
			strings.Join([]string{
				slog.LevelDebug.String(),
				slog.LevelInfo.String(),
				slog.LevelWarn.String(),
				slog.LevelError.String(),
			}, ", "))
}
