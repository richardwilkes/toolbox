// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/txt"
)

func TestToRoman(t *testing.T) {
	type data struct {
		v int
		e string
	}
	for _, one := range []data{
		{v: 1, e: "I"},
		{v: 2, e: "II"},
		{v: 3, e: "III"},
		{v: 4, e: "IV"},
		{v: 5, e: "V"},
		{v: 6, e: "VI"},
		{v: 7, e: "VII"},
		{v: 8, e: "VIII"},
		{v: 9, e: "IX"},
		{v: 10, e: "X"},
		{v: 11, e: "XI"},
		{v: 14, e: "XIV"},
		{v: 39, e: "XXXIX"},
		{v: 40, e: "XL"},
		{v: 41, e: "XLI"},
		{v: 49, e: "XLIX"},
		{v: 50, e: "L"},
		{v: 51, e: "LI"},
		{v: 89, e: "LXXXIX"},
		{v: 90, e: "XC"},
		{v: 99, e: "XCIX"},
		{v: 100, e: "C"},
		{v: 399, e: "CCCXCIX"},
		{v: 400, e: "CD"},
		{v: 499, e: "CDXCIX"},
		{v: 500, e: "D"},
		{v: 899, e: "DCCCXCIX"},
		{v: 900, e: "CM"},
		{v: 1967, e: "MCMLXVII"},
		{v: 2021, e: "MMXXI"},
	} {
		check.Equal(t, one.e, txt.RomanNumerals(one.v), "input: %d", one.v)
	}
}
