// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

// FirstN returns the first n runes of a string.
func FirstN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[:n])
}

// LastN returns the last n runes of a string.
func LastN(s string, n int) string {
	if n < 1 {
		return ""
	}
	r := []rune(s)
	if n > len(r) {
		return s
	}
	return string(r[len(r)-n:])
}

// Truncate the input string to count runes, trimming from the end if keepFirst is true or the start if not. If trimming
// occurs, a … will be added in place of the trimmed characters.
func Truncate(s string, count int, keepFirst bool) string {
	var result string
	if keepFirst {
		result = FirstN(s, count)
	} else {
		result = LastN(s, count)
	}
	if result != s {
		if keepFirst {
			result += "…"
		} else {
			result = "…" + result
		}
	}
	return result
}
