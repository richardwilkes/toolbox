// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// Bounds applied by DiscardAndCloseIgnoringErrors when draining an unread reader before closing it. They mirror the
// limits net/http uses when it drains a response body on Close (Go 1.27): stop after drainOnCloseByteLimit bytes or
// drainOnCloseTimeLimit of wall-clock time, whichever comes first, so a huge or slow-trickling source can't make
// closing expensive.
const (
	drainOnCloseByteLimit = 256 * 1024
	drainOnCloseTimeLimit = 50 * time.Millisecond
)

// CloseIgnoringErrors closes the closer and ignores any error it might produce. Should only be used for read-only
// streams of data where closing should never cause an error.
func CloseIgnoringErrors(closer io.Closer) {
	_ = closer.Close() //nolint:errcheck // intentionally ignoring any error
}

// DiscardAndCloseIgnoringErrors reads and discards any content remaining in the reader, then closes it. Draining lets a
// keep-alive HTTP/1.1 connection be reused by a later request rather than being discarded. The drain is bounded to
// drainOnCloseByteLimit bytes and drainOnCloseTimeLimit of wall-clock time, whichever is reached first, so an oversized
// or slow-trickling source can't make closing expensive; these are the same bounds net/http applies when it drains a
// response body on Close (Go 1.27). The time bound is enforced by closing the reader out from under an in-flight read,
// so the reader must tolerate a concurrent Close (net/http bodies and os.File do).
func DiscardAndCloseIgnoringErrors(rc io.ReadCloser) {
	// If the drain outlasts the time bound (e.g. a slow-trickling source), close rc to interrupt the in-progress read;
	// the trailing Close below is then a harmless second close.
	timer := time.AfterFunc(drainOnCloseTimeLimit, func() { CloseIgnoringErrors(rc) })
	_, _ = io.CopyN(io.Discard, rc, drainOnCloseByteLimit) //nolint:errcheck // intentionally ignoring any error
	timer.Stop()
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
