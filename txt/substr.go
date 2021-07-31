// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
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
