// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import (
	"unicode"
)

// Options represents a set of options.
type Options []*Option

// Len implements the sort.Interface interface.
func (op Options) Len() int {
	return len(op)
}

// Less implements the sort.Interface interface.
func (op Options) Less(i, j int) bool {
	in := op[i].single
	jn := op[j].single
	if in == 0 && jn != 0 {
		return false
	}
	if in != 0 && jn == 0 {
		return true
	}
	in = swapCase(in)
	jn = swapCase(jn)
	if in < jn {
		return true
	}
	if in == jn {
		return op[i].name < op[j].name
	}
	return false
}

// Swap implements the sort.Interface interface.
func (op Options) Swap(i, j int) {
	op[i], op[j] = op[j], op[i]
}

func swapCase(ch rune) rune {
	if unicode.IsUpper(ch) {
		return unicode.ToLower(ch)
	}
	if unicode.IsLower(ch) {
		return unicode.ToUpper(ch)
	}
	return ch
}
