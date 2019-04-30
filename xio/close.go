// Package xio provides i/o utilities.
package xio

import "io"

// CloseIgnoringErrors closes the closer and ignores any error it might
// produce. Should only be used for read-only streams of data where closing
// should never cause an error.
func CloseIgnoringErrors(closer io.Closer) {
	// The extra code here is just to quiet the linter about not checking
	// for an error.
	if err := closer.Close(); err != nil {
		return
	}
}
