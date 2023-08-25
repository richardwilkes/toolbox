// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package errs

import (
	"context"
	"log/slog"
)

// Log an error with a stack trace to the default logger.
func Log(err error) {
	LogTo(slog.Default(), err)
}

// LogTo logs an error with a stack trace to the specified logger.
func LogTo(logger *slog.Logger, err error) {
	e := WrapTyped(err)
	logger.Error(e.Message(), "stack", e.slogStackTrace())
}

// LogToWithLevel logs an error with a stack trace to the specified logger with the specified level.
func LogToWithLevel(logger *slog.Logger, level slog.Level, err error) {
	e := WrapTyped(err)
	logger.Log(context.Background(), level, e.Message(), "stack", e.slogStackTrace())
}
