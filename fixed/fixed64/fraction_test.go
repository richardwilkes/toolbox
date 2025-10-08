// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed64_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed64"
)

func TestFraction(t *testing.T) {
	c := check.New(t)
	c.Equal(fixed64.FromStringForced[fixed.D4]("0.3333"), fixed64.NewFraction[fixed.D4]("1/3").Value())
	c.Equal(fixed64.FromStringForced[fixed.D4]("0.3333"), fixed64.NewFraction[fixed.D4]("1 / 3").Value())
	c.Equal(fixed64.FromStringForced[fixed.D4]("0.3333"), fixed64.NewFraction[fixed.D4]("-1/-3").Value())
	c.Equal(fixed64.FromInteger[fixed.D4](0), fixed64.NewFraction[fixed.D4]("5/0").Value())
	c.Equal(fixed64.FromInteger[fixed.D4](5), fixed64.NewFraction[fixed.D4]("5/1").Value())
	c.Equal(fixed64.FromInteger[fixed.D4](-5), fixed64.NewFraction[fixed.D4]("-5/1").Value())
	c.Equal(fixed64.FromInteger[fixed.D4](-5), fixed64.NewFraction[fixed.D4]("5/-1").Value())
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
	frac := fixed64.NewFraction[T]("3/4") // 3/4
	// The precision depends on the decimal places of T
	result := frac.Value().String()
	c.True(strings.HasPrefix(result, "0.7")) // Should be 0.75 or 0.7 depending on precision

	// Test Normalize - only handles negative denominators and zero division
	frac2 := fixed64.NewFraction[T]("6/-8") // 6/-8 should become -6/8
	frac2.Normalize()
	c.Equal("-6/8", frac2.String()) // Normalize doesn't reduce fractions, just fixes sign

	// Test with positive denominator
	frac3 := fixed64.NewFraction[T]("6/8")
	frac3.Normalize()
	c.Equal("6/8", frac3.String()) // Should remain unchanged

	// Test StringWithSign for fractions
	posFrac := fixed64.NewFraction[T]("1/2")
	c.Equal("+1/2", posFrac.StringWithSign())

	negFrac := fixed64.NewFraction[T]("-1/2")
	c.Equal("-1/2", negFrac.StringWithSign())

	// Test JSON marshaling/unmarshaling for fractions
	data, err := json.Marshal(posFrac)
	c.NoError(err)
	c.Equal(`"1/2"`, string(data))

	var unmarshaled fixed64.Fraction[T]
	err = json.Unmarshal(data, &unmarshaled)
	c.NoError(err)
	c.Equal(posFrac, unmarshaled)
}

func TestAdditionalFactionEdgeCases(t *testing.T) {
	testAdditionalFractionEdgeCases[fixed.D1](t)
	testAdditionalFractionEdgeCases[fixed.D2](t)
	testAdditionalFractionEdgeCases[fixed.D3](t)
	testAdditionalFractionEdgeCases[fixed.D4](t)
	testAdditionalFractionEdgeCases[fixed.D5](t)
	testAdditionalFractionEdgeCases[fixed.D6](t)
}

func testAdditionalFractionEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	// Test fraction with denominator = 1
	wholeFrac := fixed64.NewFraction[T]("5")
	c.Equal("5", wholeFrac.String())
	c.Equal("+5", wholeFrac.StringWithSign())

	// Test fraction JSON unmarshaling error
	var frac fixed64.Fraction[T]
	err := json.Unmarshal([]byte("invalid json"), &frac)
	c.HasError(err)
}

func TestFractionArithmeticAndSimplify(t *testing.T) {
	testFractionArithmeticAndSimplify[fixed.D1](t)
	testFractionArithmeticAndSimplify[fixed.D2](t)
	testFractionArithmeticAndSimplify[fixed.D3](t)
	testFractionArithmeticAndSimplify[fixed.D4](t)
	testFractionArithmeticAndSimplify[fixed.D5](t)
	testFractionArithmeticAndSimplify[fixed.D6](t)
}

// nolint:gocritic // The comments aren't "commented out code"
func testFractionArithmeticAndSimplify[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	frac1 := fixed64.NewFraction[T]("1/2")
	frac2 := fixed64.NewFraction[T]("1/3")

	// Add: 1/2 + 1/3 = (3+2)/6 = 5/6
	sum := frac1.Add(frac2)
	c.Equal("5/6", sum.Simplify().String())

	// Sub: 1/2 - 1/3 = (3-2)/6 = 1/6
	diff := frac1.Sub(frac2)
	c.Equal("1/6", diff.Simplify().String())

	// Mul: 1/2 * 1/3 = 1/6
	prod := frac1.Mul(frac2)
	c.Equal("1/6", prod.Simplify().String())

	// Div: (1/2) / (1/3) = (1*3)/(2*1) = 3/2
	quot := frac1.Div(frac2)
	c.Equal("3/2", quot.Simplify().String())

	// Simplify: 2/4 = 1/2
	simple := fixed64.NewFraction[T]("2/4").Simplify()
	c.Equal("1/2", simple.String())

	// Simplify: 10/100 = 1/10
	simple2 := fixed64.NewFraction[T]("10/100").Simplify()
	c.Equal("1/10", simple2.String())

	// Simplify: 7/13 (should remain 7/13)
	simple3 := fixed64.NewFraction[T]("7/13").Simplify()
	c.Equal("7/13", simple3.String())
}
