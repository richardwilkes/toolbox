// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xunicode

import "unicode"

// VowelChecker returns true if the rune is to be considered a vowel.
type VowelChecker func(rune) bool

// IsVowel is a VowelChecker that accepts a, e, i, o, u and their common accented forms, ignoring case.
func IsVowel(ch rune) bool {
	if unicode.IsUpper(ch) {
		ch = unicode.ToLower(ch)
	}
	switch ch {
	case 'a', 'à', 'á', 'â', 'ä', 'æ', 'ã', 'å', 'ā',
		'e', 'è', 'é', 'ê', 'ë', 'ē', 'ė', 'ę',
		'i', 'î', 'ï', 'í', 'ī', 'į', 'ì',
		'o', 'ô', 'ö', 'ò', 'ó', 'œ', 'ø', 'ō', 'õ',
		'u', 'û', 'ü', 'ù', 'ú', 'ū':
		return true
	default:
		return false
	}
}

// IsVowely is like IsVowel, but also accepts y and ÿ.
func IsVowely(ch rune) bool {
	if unicode.IsUpper(ch) {
		ch = unicode.ToLower(ch)
	}
	switch ch {
	case 'y', 'ÿ':
		return true
	default:
		return IsVowel(ch)
	}
}
