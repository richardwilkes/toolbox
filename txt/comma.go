// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
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
	"strings"

	"golang.org/x/exp/constraints"
)

// Comma returns text version of the value that uses commas for every 3 orders of magnitude.
func Comma[T constraints.Integer | constraints.Float](value T) string {
	return CommaFromStringNum(fmt.Sprintf("%v", value))
}

// CommaFromStringNum returns a revised version of the numeric input string that uses commas for every 3 orders of
// magnitude.
func CommaFromStringNum(s string) string {
	var buffer strings.Builder
	if strings.HasPrefix(s, "-") {
		buffer.WriteByte('-')
		s = s[1:]
	}
	parts := strings.Split(s, ".")
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
		buffer.Write([]byte{'.'})
		buffer.WriteString(parts[1])
	}
	return buffer.String()
}
