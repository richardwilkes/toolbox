// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/errs"
)

const utf8BOM = '\uFEFF'

// NewBOMStripper strips a leading UTF-8 BOM marker from the input. The reader that is returned will be the same as the
// one passed in if it was a *bufio.Reader, otherwise, the original reader will be wrapped with a *bufio.Reader and
// returned.
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
