// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"log/slog"
	"os"

	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xflag"
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
		// Fan out at the handler level rather than the writer level so each destination auto-detects color support from
		// its own writer. An io.MultiWriter broadcasts identical bytes to every destination, so a single handler over
		// it can only ever be all-colored or all-plain -- it cannot produce colored console output plus a plain file.
		// Each PrettyHandler here auto-detects: DetectKind returns Dumb (no escapes) for the non-terminal rotating file
		// and a color kind for a TTY stdout.
		newOpts := func() *PrettyOptions {
			return &PrettyOptions{
				HandlerOptions: slog.HandlerOptions{
					AddSource: true,
					Level:     c.LogLevel,
				},
			}
		}
		var handler slog.Handler = NewPrettyHandler(c.RotatorCfg.NewWriteCloser(), newOpts())
		if c.Console {
			handler = NewMultiHandler(handler, NewPrettyHandler(os.Stdout, newOpts()))
		}
		slog.SetDefault(slog.New(handler))
	})
}
