// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xio provides i/o utilities.
package xio

import (
	"io"
	"log/slog"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// CloseIgnoringErrors closes the closer and ignores any error it might produce. Should only be used for read-only
// streams of data where closing should never cause an error.
func CloseIgnoringErrors(closer io.Closer) {
	_ = closer.Close() //nolint:errcheck // intentionally ignoring any error
}

// DiscardAndCloseIgnoringErrors reads any content remaining in the body and discards it, then closes the body.
func DiscardAndCloseIgnoringErrors(rc io.ReadCloser) {
	_, _ = io.Copy(io.Discard, rc) //nolint:errcheck // intentionally ignoring any error
	CloseIgnoringErrors(rc)
}

// CloseLoggingErrors closes the closer and logs any errors that occur to the default logger.
func CloseLoggingErrors(closer io.Closer) {
	CloseLoggingErrorsTo(slog.Default(), closer)
}

// CloseLoggingErrorsTo closes the closer and logs any errors that occur to the provided logger.
func CloseLoggingErrorsTo(logger *slog.Logger, closer io.Closer) {
	if err := closer.Close(); err != nil {
		errs.LogTo(logger, err)
	}
}
