// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/fixed"
)

// TestDecimalPlaces verifies, for every Dx type, that Places() reports the expected count and that Multiplier() equals
// 10^Places(). Computing the expected multiplier independently guards against a typo in the hand-written constants.
func TestDecimalPlaces(t *testing.T) {
	c := check.New(t)
	cases := []struct {
		dx     fixed.Dx
		places int
	}{
		{dx: fixed.D1(0), places: 1},
		{dx: fixed.D2(0), places: 2},
		{dx: fixed.D3(0), places: 3},
		{dx: fixed.D4(0), places: 4},
		{dx: fixed.D5(0), places: 5},
		{dx: fixed.D6(0), places: 6},
		{dx: fixed.D7(0), places: 7},
		{dx: fixed.D8(0), places: 8},
		{dx: fixed.D9(0), places: 9},
		{dx: fixed.D10(0), places: 10},
		{dx: fixed.D11(0), places: 11},
		{dx: fixed.D12(0), places: 12},
		{dx: fixed.D13(0), places: 13},
		{dx: fixed.D14(0), places: 14},
		{dx: fixed.D15(0), places: 15},
		{dx: fixed.D16(0), places: 16},
	}
	for _, tc := range cases {
		c.Equal(tc.places, tc.dx.Places())
		var want int64 = 1
		for range tc.places {
			want *= 10
		}
		c.Equal(want, tc.dx.Multiplier(), "Multiplier for %d places", tc.places)
	}
}
