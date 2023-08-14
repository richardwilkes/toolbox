// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xio provides i/o utilities.
package xio

import "io"

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
