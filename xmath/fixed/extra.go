// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed

import (
	"errors"
	"unicode/utf8"
)

var (
	errDoesNotFitInFloat64 = errors.New("does not fit in float64")
	errDoesNotFitInInt64   = errors.New("does not fit in int64")
)

func unquote(text []byte) string {
	if len(text) > 1 {
		if ch, _ := utf8.DecodeRune(text); ch == '"' {
			if ch, _ = utf8.DecodeLastRune(text); ch == '"' {
				text = text[1 : len(text)-1]
			}
		}
	}
	return string(text)
}
