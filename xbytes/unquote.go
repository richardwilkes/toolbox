// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes

import (
	"unicode/utf8"
)

// Unquote strips up to one set of surrounding single or double quotes from the bytes and returns them. For a more
// capable version that supports different quoting types and unescaping, consider using strconv.Unquote().
func Unquote(text []byte) []byte {
	if len(text) > 1 {
		if ch1, _ := utf8.DecodeRune(text); ch1 == '"' || ch1 == '\'' {
			if ch2, _ := utf8.DecodeLastRune(text); ch1 == ch2 {
				text = text[1 : len(text)-1]
			}
		}
	}
	return text
}
