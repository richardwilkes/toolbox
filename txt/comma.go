// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// CommaInt returns text version of the value that uses commas for every 3 orders of magnitude.
func CommaInt[T constraints.Integer](value T) string {
	return CommaFromStringNum(fmt.Sprintf("%d", value))
}

// CommaFloat returns text version of the value that uses commas for every 3 orders of magnitude.
func CommaFloat[T constraints.Float](value T) string {
	return CommaFromStringNum(strconv.FormatFloat(float64(value), 'f', -1, 64))
}

// CommaFromStringNum returns a revised version of the numeric input string that uses commas for every 3 orders of
// magnitude. Note that this function assumes the input is nothing more than an optional leading sign followed by
// digits.
func CommaFromStringNum(s string) string {
	if s == "" {
		return ""
	}
	var buffer strings.Builder
	if s[0] == '-' || s[0] == '+' {
		buffer.WriteByte(s[0])
		s = s[1:]
	}
	parts := strings.SplitN(s, ".", 2)
	i := 0
	needComma := false
	if len(parts[0])%3 != 0 {
		i += len(parts[0]) % 3
		buffer.WriteString(parts[0][:i])
		needComma = true
	}
	for ; i < len(parts[0]); i += 3 {
		if needComma {
			buffer.WriteByte(',')
		} else {
			needComma = true
		}
		buffer.WriteString(parts[0][i : i+3])
	}
	if len(parts) > 1 {
		buffer.WriteByte('.')
		buffer.WriteString(parts[1])
	}
	return buffer.String()
}
