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
	"fmt"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xmath"
)

// CommaInt formats the value with a comma inserted every three digits.
func CommaInt[T xmath.Integer](value T) string {
	return CommaFromStringNum(fmt.Sprintf("%d", value))
}

// CommaFloat formats the value with a comma inserted every three digits of the integer part.
func CommaFloat[T xmath.Float](value T) string {
	return CommaFromStringNum(strconv.FormatFloat(float64(value), 'f', -1, 64))
}

// CommaFromStringNum inserts a comma every three digits into the integer part of a numeric string. The input must
// consist of an optional leading sign, digits, and an optional decimal point followed by more digits.
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
