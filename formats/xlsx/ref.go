// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xlsx

import (
	"strconv"
	"strings"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Ref holds a cell reference.
type Ref struct {
	Row int
	Col int
}

// ParseRef parses a string into a Ref.
func ParseRef(str string) Ref {
	r := Ref{}
	state := 0
	for _, ch := range strings.ToUpper(str) {
		if state == 0 {
			if ch >= 'A' && ch <= 'Z' {
				r.Col *= 26
				r.Col += int(1 + ch - 'A')
			} else {
				state = 1
			}
		}
		if state == 1 {
			if ch >= '0' && ch <= '9' {
				r.Row *= 10
				r.Row += int(ch - '0')
			} else {
				break
			}
		}
	}
	if r.Col > 0 {
		r.Col--
	}
	if r.Row > 0 {
		r.Row--
	}
	return r
}

func (r Ref) String() string {
	var a [65]byte
	i := len(a)
	col := r.Col
	for col >= 26 {
		i--
		q := col / 26
		a[i] = letters[col-q*26]
		col = q - 1
	}
	i--
	a[i] = letters[col]
	return string(a[i:]) + strconv.Itoa(r.Row+1)
}
