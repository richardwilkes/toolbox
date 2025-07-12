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

// LogFlagPriority is the priority used when adding the automatic post-parse function for the log flags.
const LogFlagPriority = -5

// Config holds the configuration for logging.
type Config struct {
	RotatorCfg Rotator
	LogLevel   LevelValue
	Console    bool
}

// AddFlags for configuring logging to the command line. Adds a post-parse function that will setup logging. Note that
// this automatic handling only works if xflag.Parse is used and not flag.Parse directly.
func (c *Config) AddFlags() {
	c.RotatorCfg.AddFlags()
	c.LogLevel.AddFlags()
	if c.Console {
		quietUsage := i18n.Text("Suppress the log output to the console")
		quietFunc := func(_ string) error {
			c.Console = false
			return nil
		}
		flag.BoolFunc("quiet", quietUsage, quietFunc)
		flag.BoolFunc("q", quietUsage, quietFunc)
	} else {
		flag.BoolVar(&c.Console, "console", false, i18n.Text("Copy the log output to the console"))
	}
	xflag.AddPostParseFunc(LogFlagPriority, func() {
		w := io.Writer(c.RotatorCfg.NewWriteCloser())
		if c.Console {
			w = io.MultiWriter(w, os.Stdout)
		}
		slog.SetDefault(slog.New(NewPrettyHandler(w, &PrettyOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: true,
				Level:     c.LogLevel,
			},
			ColorSupportOverride: xterm.DetectKind(os.Stdout),
		})))
	})
}
