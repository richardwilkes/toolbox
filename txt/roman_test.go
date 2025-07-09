// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/txt"
)

func TestToRoman(t *testing.T) {
	c := check.New(t)

	// Test basic single digits (1-9)
	c.Equal("I", txt.RomanNumerals(1))
	c.Equal("II", txt.RomanNumerals(2))
	c.Equal("III", txt.RomanNumerals(3))
	c.Equal("IV", txt.RomanNumerals(4))
	c.Equal("V", txt.RomanNumerals(5))
	c.Equal("VI", txt.RomanNumerals(6))
	c.Equal("VII", txt.RomanNumerals(7))
	c.Equal("VIII", txt.RomanNumerals(8))
	c.Equal("IX", txt.RomanNumerals(9))

	// Test tens
	c.Equal("X", txt.RomanNumerals(10))
	c.Equal("XI", txt.RomanNumerals(11))
	c.Equal("XIV", txt.RomanNumerals(14))

	// Test numbers in the 30s and 40s
	c.Equal("XXXIX", txt.RomanNumerals(39))
	c.Equal("XL", txt.RomanNumerals(40))
	c.Equal("XLI", txt.RomanNumerals(41))
	c.Equal("XLIX", txt.RomanNumerals(49))

	// Test 50s
	c.Equal("L", txt.RomanNumerals(50))
	c.Equal("LI", txt.RomanNumerals(51))

	// Test numbers in the 80s and 90s
	c.Equal("LXXXIX", txt.RomanNumerals(89))
	c.Equal("XC", txt.RomanNumerals(90))
	c.Equal("XCIX", txt.RomanNumerals(99))

	// Test hundreds
	c.Equal("C", txt.RomanNumerals(100))
	c.Equal("CCCXCIX", txt.RomanNumerals(399))
	c.Equal("CD", txt.RomanNumerals(400))
	c.Equal("CDXCIX", txt.RomanNumerals(499))

	// Test 500s
	c.Equal("D", txt.RomanNumerals(500))
	c.Equal("DCCCXCIX", txt.RomanNumerals(899))
	c.Equal("CM", txt.RomanNumerals(900))

	// Test complex historical dates
	c.Equal("MCMLXVII", txt.RomanNumerals(1967))
	c.Equal("MMXXI", txt.RomanNumerals(2021))
}
