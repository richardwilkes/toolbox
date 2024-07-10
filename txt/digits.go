// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
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
	"unicode"

	"github.com/richardwilkes/toolbox/errs"
)

// DigitToValue converts a unicode digit into a numeric value.
func DigitToValue(ch rune) (int, error) {
	if ch < '\U00010000' {
		r16 := uint16(ch)
		for _, one := range unicode.Digit.R16 {
			if one.Lo <= r16 && one.Hi >= r16 {
				return int(r16 - one.Lo), nil
			}
		}
	} else {
		r32 := uint32(ch)
		for _, one := range unicode.Digit.R32 {
			if one.Lo <= r32 && one.Hi >= r32 {
				return int(r32 - one.Lo), nil
			}
		}
	}
	return 0, errs.New(fmt.Sprintf("Not a digit: '%v'", ch))
}
