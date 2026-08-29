// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import (
	"bufio"
	"io"

	"github.com/richardwilkes/toolbox/v2/errs"
)

const utf8BOM = '\uFEFF'

// NewBOMStripper strips a leading UTF-8 BOM from r. The returned reader is r itself if it is already a *bufio.Reader;
// otherwise r is wrapped in a new one. An error (wrapping io.EOF for empty input) is returned if the first rune cannot
// be read.
func NewBOMStripper(r io.Reader) (*bufio.Reader, error) {
	buffer, ok := r.(*bufio.Reader)
	if !ok {
		buffer = bufio.NewReader(r)
	}
	ch, _, err := buffer.ReadRune()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if ch != utf8BOM {
		if err = buffer.UnreadRune(); err != nil {
			return nil, errs.Wrap(err)
		}
	}
	return buffer, nil
}
