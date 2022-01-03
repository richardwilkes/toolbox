// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import "io"

// TeeWriter is a writer that writes to multiple other writers.
type TeeWriter struct {
	Writers []io.Writer
}

// Write to each of the underlying streams.
func (t *TeeWriter) Write(p []byte) (n int, err error) {
	var curErr error
	for _, w := range t.Writers {
		if n, curErr = w.Write(p); curErr != nil {
			err = curErr
		}
	}
	return
}
