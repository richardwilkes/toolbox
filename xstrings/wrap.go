// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import (
	"strings"
	"unicode/utf8"
)

// Wrap text to a certain length, giving it an optional prefix on each line. Words will not be broken, even if they
// exceed the maximum column size and instead will extend past the desired length. Column counts are measured in runes,
// not bytes, so multibyte text wraps at its visible width.
func Wrap(prefix, text string, maxColumns int) string {
	var buffer strings.Builder
	prefixLen := utf8.RuneCountInString(prefix)
	for i, line := range strings.Split(text, "\n") {
		if i != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(prefix)
		avail := maxColumns - prefixLen
		for j, token := range strings.Fields(line) {
			tokenLen := utf8.RuneCountInString(token)
			if j != 0 {
				if 1+tokenLen > avail {
					buffer.WriteByte('\n')
					buffer.WriteString(prefix)
					avail = maxColumns - prefixLen
				} else {
					buffer.WriteByte(' ')
					avail--
				}
			}
			buffer.WriteString(token)
			avail -= tokenLen
		}
	}
	return buffer.String()
}
