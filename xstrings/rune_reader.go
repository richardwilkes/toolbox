// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import "io"

var _ io.RuneReader = &RuneReader{}

// RuneReader implements io.RuneReader over a []rune.
//
// Unlike the io.RuneReader contract, ReadRune reports a size of 1 for every rune rather than its UTF-8 byte length.
// This is deliberate: it makes the regexp *Reader methods (e.g. FindReaderIndex) report positions as rune indices into
// Src, which callers indexing back into the []rune need (see ToCamelCaseWithExceptions). Do not change this without
// updating those callers.
type RuneReader struct {
	Src []rune
	Pos int
}

// ReadRune returns the next rune with a size of 1 (see the RuneReader documentation), or -1, 0, io.EOF at the end.
func (rr *RuneReader) ReadRune() (r rune, size int, err error) {
	if rr.Pos >= len(rr.Src) {
		return -1, 0, io.EOF
	}
	nextRune := rr.Src[rr.Pos]
	rr.Pos++
	return nextRune, 1, nil
}
