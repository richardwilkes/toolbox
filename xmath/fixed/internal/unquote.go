// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package internal

import (
	"unicode/utf8"
)

// Unquote strips up to one set of surrounding double quotes from the bytes and returns them as a string.
func Unquote(text []byte) string {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRune(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRune(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return string(text)
}