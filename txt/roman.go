// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import (
	"strings"
)

var (
	romanValues = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romanText   = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
)

// RomanNumerals converts a number into roman numerals.
func RomanNumerals(value int) string {
	var buffer strings.Builder
	for value > 0 {
		for i, v := range romanValues {
			if value >= v {
				buffer.WriteString(romanText[i])
				value -= v
				break
			}
		}
	}
	return buffer.String()
}
