// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package jot

// LoggerWriter provides a bridge between the standard log.Logger and the jot package. You can use it like this:
//
// log.New(&jot.LoggerWriter{}, "", 0)
//
// This will send all output for this logger to the jot.Error() call.
//
// You can also set the Filter function to direct the output to a particular jot logging method:
//
// log.New(&jot.LoggerWriter{Filter: jot.Info}), "", 0)
//
// Deprecated: Use slog instead. August 28, 2023
type LoggerWriter struct {
	Filter func(v ...any)
}

// Write implements the io.Writer interface required by log.Logger.
//
// Deprecated: Use slog instead. August 28, 2023
func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		filter := w.Filter
		if filter == nil {
			filter = Error
		}
		filter(string(p[:len(p)-1]))
		Flush() // To ensure the output is recorded.
	}
	return len(p), nil
}
