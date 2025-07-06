// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed128_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed128"
)

func TestFraction(t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromStringForced[fixed.D4]("0.3333"), fixed128.NewFraction[fixed.D4]("1/3").Value())
	c.Equal(fixed128.FromStringForced[fixed.D4]("0.3333"), fixed128.NewFraction[fixed.D4]("1 / 3").Value())
	c.Equal(fixed128.FromStringForced[fixed.D4]("0.3333"), fixed128.NewFraction[fixed.D4]("-1/-3").Value())
	c.Equal(fixed128.From[fixed.D4](0), fixed128.NewFraction[fixed.D4]("5/0").Value())
	c.Equal(fixed128.From[fixed.D4](5), fixed128.NewFraction[fixed.D4]("5/1").Value())
	c.Equal(fixed128.From[fixed.D4](-5), fixed128.NewFraction[fixed.D4]("-5/1").Value())
	c.Equal(fixed128.From[fixed.D4](-5), fixed128.NewFraction[fixed.D4]("5/-1").Value())
}

func TestFractionFunctions(t *testing.T) {
	testFractionFunctions[fixed.D1](t)
	testFractionFunctions[fixed.D2](t)
	testFractionFunctions[fixed.D3](t)
	testFractionFunctions[fixed.D4](t)
	testFractionFunctions[fixed.D5](t)
	testFractionFunctions[fixed.D6](t)
}

func testFractionFunctions[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	// Test NewFraction
	frac := fixed128.NewFraction[T]("3/4") // 3/4
	// The precision depends on the decimal places of T
	result := frac.Value().String()
	c.True(strings.HasPrefix(result, "0.7")) // Should be 0.75 or 0.7 depending on precision

	// Test Normalize - only handles negative denominators and zero division
	frac2 := fixed128.NewFraction[T]("6/-8") // 6/-8 should become -6/8
	frac2.Normalize()
	c.Equal("-6/8", frac2.String()) // Normalize doesn't reduce fractions, just fixes sign

	// Test with positive denominator
	frac3 := fixed128.NewFraction[T]("6/8")
	frac3.Normalize()
	c.Equal("6/8", frac3.String()) // Should remain unchanged

	// Test StringWithSign for fractions
	posFrac := fixed128.NewFraction[T]("1/2")
	c.Equal("+1/2", posFrac.StringWithSign())

	negFrac := fixed128.NewFraction[T]("-1/2")
	c.Equal("-1/2", negFrac.StringWithSign())

	// Test JSON marshaling/unmarshaling for fractions
	data, err := json.Marshal(posFrac)
	c.NoError(err)
	c.Equal(`"1/2"`, string(data))

	var unmarshaled fixed128.Fraction[T]
	err = json.Unmarshal(data, &unmarshaled)
	c.NoError(err)
	c.Equal(posFrac, unmarshaled)
}

func testAdditionalFractionEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	// Test fraction with denominator = 1
	wholeFrac := fixed128.NewFraction[T]("5")
	c.Equal("5", wholeFrac.String())
	c.Equal("+5", wholeFrac.StringWithSign())

	// Test fraction JSON unmarshaling error
	var frac fixed128.Fraction[T]
	err := json.Unmarshal([]byte("invalid json"), &frac)
	c.HasError(err)
}
