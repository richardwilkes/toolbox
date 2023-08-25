// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import "io"

// RuneReader implements io.RuneReader
type RuneReader struct {
	Src []rune
	Pos int
}

// ReadRune returns the next rune and its size in bytes.
func (rr *RuneReader) ReadRune() (r rune, size int, err error) {
	if rr.Pos >= len(rr.Src) {
		return -1, 0, io.EOF
	}
	nextRune := rr.Src[rr.Pos]
	rr.Pos++
	return nextRune, 1, nil
}
