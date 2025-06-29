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
	"os"
)

// SetupStd creates a new logger with some standard configuration and sets it as the default.
func SetupStd(forDev bool) {
	var h slog.Handler
	opts := slog.HandlerOptions{AddSource: true}
	if forDev {
		opts.Level = slog.LevelDebug
		h = NewPrettyHandler(os.Stdout, &PrettyOptions{HandlerOptions: opts})
	} else {
		h = slog.NewJSONHandler(os.Stdout, &opts)
	}
	slog.SetDefault(slog.New(h))
}
