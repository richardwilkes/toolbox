// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestToRoman(t *testing.T) {
	c := check.New(t)

	// Test basic single digits (1-9)
	c.Equal("I", xstrings.RomanNumerals(1))
	c.Equal("II", xstrings.RomanNumerals(2))
	c.Equal("III", xstrings.RomanNumerals(3))
	c.Equal("IV", xstrings.RomanNumerals(4))
	c.Equal("V", xstrings.RomanNumerals(5))
	c.Equal("VI", xstrings.RomanNumerals(6))
	c.Equal("VII", xstrings.RomanNumerals(7))
	c.Equal("VIII", xstrings.RomanNumerals(8))
	c.Equal("IX", xstrings.RomanNumerals(9))

	// Test tens
	c.Equal("X", xstrings.RomanNumerals(10))
	c.Equal("XI", xstrings.RomanNumerals(11))
	c.Equal("XIV", xstrings.RomanNumerals(14))

	// Test numbers in the 30s and 40s
	c.Equal("XXXIX", xstrings.RomanNumerals(39))
	c.Equal("XL", xstrings.RomanNumerals(40))
	c.Equal("XLI", xstrings.RomanNumerals(41))
	c.Equal("XLIX", xstrings.RomanNumerals(49))

	// Test 50s
	c.Equal("L", xstrings.RomanNumerals(50))
	c.Equal("LI", xstrings.RomanNumerals(51))

	// Test numbers in the 80s and 90s
	c.Equal("LXXXIX", xstrings.RomanNumerals(89))
	c.Equal("XC", xstrings.RomanNumerals(90))
	c.Equal("XCIX", xstrings.RomanNumerals(99))

	// Test hundreds
	c.Equal("C", xstrings.RomanNumerals(100))
	c.Equal("CCCXCIX", xstrings.RomanNumerals(399))
	c.Equal("CD", xstrings.RomanNumerals(400))
	c.Equal("CDXCIX", xstrings.RomanNumerals(499))

	// Test 500s
	c.Equal("D", xstrings.RomanNumerals(500))
	c.Equal("DCCCXCIX", xstrings.RomanNumerals(899))
	c.Equal("CM", xstrings.RomanNumerals(900))

	// Test complex historical dates
	c.Equal("MCMLXVII", xstrings.RomanNumerals(1967))
	c.Equal("MMXXI", xstrings.RomanNumerals(2021))
}
