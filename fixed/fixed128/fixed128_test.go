// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed128"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed64"
	"gopkg.in/yaml.v3"
)

// TestFromStringLeadingPlus verifies that a leading "+" with an empty integer part (e.g. "+.5", "+") is accepted, as
// fixed64.FromString accepts it, rather than failing the underlying big.Int parse.
func TestFromStringLeadingPlus(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		in   string
		want string
	}{
		{in: "+.5", want: "0.5"},
		{in: "+", want: "0"},
		{in: "+0", want: "0"},
		{in: "+0.5", want: "0.5"},
		{in: "+5", want: "5"},
		{in: "+5.25", want: "5.25"},
	} {
		v, err := fixed128.FromString[fixed.D2](tc.in)
		c.NoError(err, tc.in)
		c.Equal(tc.want, v.String(), tc.in)

		// The 64- and 128-bit implementations must agree on these inputs.
		v64, err64 := fixed64.FromString[fixed.D2](tc.in)
		c.NoError(err64, tc.in)
		c.Equal(v64.String(), v.String(), tc.in)
	}
}

// TestMulDivSaturate verifies that Mul and Div saturate to Maximum/Minimum when the result overflows 128 bits, rather
// than silently wrapping the intermediate product.
func TestMulDivSaturate(t *testing.T) {
	c := check.New(t)
	maximum := fixed128.Maximum[fixed.D2]()
	minimum := fixed128.Minimum[fixed.D2]()
	two := fixed128.FromInteger[fixed.D2](2)
	negTwo := fixed128.FromInteger[fixed.D2](-2)
	half := fixed128.FromStringForced[fixed.D2]("0.5")
	negHalf := fixed128.FromStringForced[fixed.D2]("-0.5")

	// Mul overflow saturates to the correct extreme based on sign.
	c.Equal(maximum, maximum.Mul(two))
	c.Equal(minimum, minimum.Mul(two))
	c.Equal(minimum, maximum.Mul(negTwo))
	c.Equal(maximum, minimum.Mul(negTwo))

	// Div by a magnitude < 1 magnifies and can overflow; it must saturate too.
	c.Equal(maximum, maximum.Div(half))
	c.Equal(minimum, minimum.Div(half))
	c.Equal(minimum, maximum.Div(negHalf))

	// An operand too large for int64 whose product still fits takes the big.Int path and must not saturate.
	c.Equal("300000000000000000000",
		fixed128.FromStringForced[fixed.D2]("100000000000000000000").Mul(fixed128.FromInteger[fixed.D2](3)).String())
}

func TestConversion(t *testing.T) {
	testConversion[fixed.D1](t)
	testConversion[fixed.D2](t)
	testConversion[fixed.D3](t)
	testConversion[fixed.D4](t)
	testConversion[fixed.D5](t)
	testConversion[fixed.D6](t)
}

//nolint:goconst // Not helpful
func testConversion[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal("0.1", fixed128.FromFloat[T](0.1).String())
	c.Equal("0.2", fixed128.FromFloat[T](0.2).String())
	c.Equal("0.3", fixed128.FromStringForced[T]("0.3").String())
	c.Equal("-0.1", fixed128.FromFloat[T](-0.1).String())
	c.Equal("-0.2", fixed128.FromFloat[T](-0.2).String())
	c.Equal("-0.3", fixed128.FromStringForced[T]("-0.3").String())
	threeFill := strings.Repeat("3", fixed128.MaxDecimalDigits[T]())
	c.Equal("0."+threeFill, fixed128.FromStringForced[T]("0.33333333").String())
	c.Equal("-0."+threeFill, fixed128.FromStringForced[T]("-0.33333333").String())
	sixFill := strings.Repeat("6", fixed128.MaxDecimalDigits[T]())
	c.Equal("0."+sixFill, fixed128.FromStringForced[T]("0.66666666").String())
	c.Equal("1", fixed128.FromFloat[T](1.0000004).String())
	c.Equal("1", fixed128.FromFloat[T](1.00000049).String())
	c.Equal("1", fixed128.FromFloat[T](1.0000005).String())
	c.Equal("1", fixed128.FromFloat[T](1.0000009).String())
	c.Equal("-1", fixed128.FromFloat[T](-1.0000004).String())
	c.Equal("-1", fixed128.FromFloat[T](-1.00000049).String())
	c.Equal("-1", fixed128.FromFloat[T](-1.0000005).String())
	c.Equal("-1", fixed128.FromFloat[T](-1.0000009).String())
	zeroFill := strings.Repeat("0", fixed128.MaxDecimalDigits[T]()-1)
	c.Equal("0."+zeroFill+"4", fixed128.FromStringForced[T]("0."+zeroFill+"405").String())
	c.Equal("-0."+zeroFill+"4", fixed128.FromStringForced[T]("-0."+zeroFill+"405").String())

	v, err := fixed128.FromString[T]("33.0")
	c.NoError(err)
	c.Equal(v, fixed128.FromInteger[T](33))

	v, err = fixed128.FromString[T]("33.00000000000000000000")
	c.NoError(err)
	c.Equal(v, fixed128.FromInteger[T](33))
}

func TestAddSub(t *testing.T) {
	testAddSub[fixed.D1](t)
	testAddSub[fixed.D2](t)
	testAddSub[fixed.D3](t)
	testAddSub[fixed.D4](t)
	testAddSub[fixed.D5](t)
	testAddSub[fixed.D6](t)
}

func testAddSub[T fixed.Dx](t *testing.T) {
	oneThird := fixed128.FromStringForced[T]("0.333333")
	negTwoThirds := fixed128.FromStringForced[T]("-0.666666")
	one := fixed128.FromInteger[T](1)
	oneAndTwoThirds := fixed128.FromStringForced[T]("1.666666")
	nineThousandSix := fixed128.FromInteger[T](9006)
	two := fixed128.FromInteger[T](2)
	c := check.New(t)
	c.Equal("0."+strings.Repeat("9", fixed128.MaxDecimalDigits[T]()), oneThird.Add(oneThird).Add(oneThird).String())
	c.Equal("0."+strings.Repeat("6", fixed128.MaxDecimalDigits[T]()-1)+"7", one.Sub(oneThird).String())
	c.Equal("-1."+strings.Repeat("6", fixed128.MaxDecimalDigits[T]()), negTwoThirds.Sub(one).String())
	c.Equal("0", negTwoThirds.Sub(one).Add(oneAndTwoThirds).String())
	c.Equal(fixed128.FromInteger[T](10240), fixed128.FromInteger[T](1234).Add(nineThousandSix))
	c.Equal("10240", fixed128.FromInteger[T](1234).Add(nineThousandSix).String())
	c.Equal("-1.5", fixed128.FromFloat[T](0.5).Sub(two).String())
	ninetyPointZeroSix := fixed128.FromStringForced[T]("90.06")
	twelvePointThirtyFour := fixed128.FromStringForced[T]("12.34")
	var answer string
	if fixed128.MaxDecimalDigits[T]() > 1 {
		answer = "102.4"
	} else {
		answer = "102.3"
	}
	c.Equal(fixed128.FromStringForced[T](answer), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	c.Equal(answer, twelvePointThirtyFour.Add(ninetyPointZeroSix).String())
}

func TestMulDiv(t *testing.T) {
	testMulDiv[fixed.D1](t)
	testMulDiv[fixed.D2](t)
	testMulDiv[fixed.D3](t)
	testMulDiv[fixed.D4](t)
	testMulDiv[fixed.D5](t)
	testMulDiv[fixed.D6](t)
}

func testMulDiv[T fixed.Dx](t *testing.T) {
	pointThree := fixed128.FromStringForced[T]("0.3")
	negativePointThree := fixed128.FromStringForced[T]("-0.3")
	threeFill := strings.Repeat("3", fixed128.MaxDecimalDigits[T]())
	c := check.New(t)
	c.Equal("0."+threeFill, fixed128.FromInteger[T](1).Div(fixed128.FromInteger[T](3)).String())
	c.Equal("-0."+threeFill, fixed128.FromInteger[T](1).Div(fixed128.FromInteger[T](-3)).String())
	c.Equal("0.1", pointThree.Div(fixed128.FromInteger[T](3)).String())
	c.Equal("0.9", pointThree.Mul(fixed128.FromInteger[T](3)).String())
	c.Equal("-0.9", negativePointThree.Mul(fixed128.FromInteger[T](3)).String())
}

func TestMod(t *testing.T) {
	testMod[fixed.D1](t)
	testMod[fixed.D2](t)
	testMod[fixed.D3](t)
	testMod[fixed.D4](t)
	testMod[fixed.D5](t)
	testMod[fixed.D6](t)
}

func testMod[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromInteger[T](1), fixed128.FromInteger[T](3).Mod(fixed128.FromInteger[T](2)))
	c.Equal(fixed128.FromStringForced[T]("0.3"), fixed128.FromStringForced[T]("9.3").Mod(fixed128.FromInteger[T](3)))
	c.Equal(fixed128.FromStringForced[T]("0.1"), fixed128.FromStringForced[T]("3.1").Mod(fixed128.FromStringForced[T]("0.2")))
}

func TestFloor(t *testing.T) {
	testFloor[fixed.D1](t)
	testFloor[fixed.D2](t)
	testFloor[fixed.D3](t)
	testFloor[fixed.D4](t)
	testFloor[fixed.D5](t)
	testFloor[fixed.D6](t)
}

func testFloor[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("0.3333").Floor())
	c.Equal(fixed128.FromInteger[T](2), fixed128.FromStringForced[T]("2.6789").Floor())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](3).Floor())
	c.Equal(fixed128.FromInteger[T](-1), fixed128.FromStringForced[T]("-0.3333").Floor())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromStringForced[T]("-2.6789").Floor())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromInteger[T](-3).Floor())
}

func TestTrunc(t *testing.T) {
	testTrunc[fixed.D1](t)
	testTrunc[fixed.D2](t)
	testTrunc[fixed.D3](t)
	testTrunc[fixed.D4](t)
	testTrunc[fixed.D5](t)
	testTrunc[fixed.D6](t)
}

func testTrunc[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("0.3333").Trunc())
	c.Equal(fixed128.FromInteger[T](2), fixed128.FromStringForced[T]("2.6789").Trunc())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](3).Trunc())
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("-0.3333").Trunc())
	c.Equal(fixed128.FromInteger[T](-2), fixed128.FromStringForced[T]("-2.6789").Trunc())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromInteger[T](-3).Trunc())
}

func TestCeil(t *testing.T) {
	testCeil[fixed.D1](t)
	testCeil[fixed.D2](t)
	testCeil[fixed.D3](t)
	testCeil[fixed.D4](t)
	testCeil[fixed.D5](t)
	testCeil[fixed.D6](t)
}

func testCeil[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromInteger[T](1), fixed128.FromStringForced[T]("0.3333").Ceil())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromStringForced[T]("2.6789").Ceil())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](3).Ceil())
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("-0.3333").Ceil())
	c.Equal(fixed128.FromInteger[T](-2), fixed128.FromStringForced[T]("-2.6789").Ceil())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromInteger[T](-3).Ceil())
}

func TestRound(t *testing.T) {
	testRound[fixed.D1](t)
	testRound[fixed.D2](t)
	testRound[fixed.D3](t)
	testRound[fixed.D4](t)
	testRound[fixed.D5](t)
	testRound[fixed.D6](t)
}

func testRound[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("0.3333").Round())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromStringForced[T]("2.6789").Round())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](3).Round())
	c.Equal(fixed128.FromInteger[T](0), fixed128.FromStringForced[T]("-0.3333").Round())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromStringForced[T]("-2.6789").Round())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromInteger[T](-3).Round())
	c.Equal(fixed128.FromInteger[T](1), fixed128.FromStringForced[T]("0.5").Round())
	c.Equal(fixed128.FromInteger[T](-1), fixed128.FromStringForced[T]("-0.5").Round())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromStringForced[T]("2.5").Round())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromStringForced[T]("-2.5").Round())
}

func TestAbs(t *testing.T) {
	testAbs[fixed.D1](t)
	testAbs[fixed.D2](t)
	testAbs[fixed.D3](t)
	testAbs[fixed.D4](t)
	testAbs[fixed.D5](t)
	testAbs[fixed.D6](t)
}

func testAbs[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromStringForced[T]("0.3333"), fixed128.FromStringForced[T]("0.3333").Abs())
	c.Equal(fixed128.FromStringForced[T]("2.6789"), fixed128.FromStringForced[T]("2.6789").Abs())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](3).Abs())
	c.Equal(fixed128.FromStringForced[T]("0.3333"), fixed128.FromStringForced[T]("-0.3333").Abs())
	c.Equal(fixed128.FromStringForced[T]("2.6789"), fixed128.FromStringForced[T]("-2.6789").Abs())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](-3).Abs())
}

func TestNeg(t *testing.T) {
	testNeg[fixed.D1](t)
	testNeg[fixed.D2](t)
	testNeg[fixed.D3](t)
	testNeg[fixed.D4](t)
	testNeg[fixed.D5](t)
	testNeg[fixed.D6](t)
}

func testNeg[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(fixed128.FromStringForced[T]("-0.3333"), fixed128.FromStringForced[T]("0.3333").Neg())
	c.Equal(fixed128.FromStringForced[T]("-2.6789"), fixed128.FromStringForced[T]("2.6789").Neg())
	c.Equal(fixed128.FromInteger[T](-3), fixed128.FromInteger[T](3).Neg())
	c.Equal(fixed128.FromStringForced[T]("0.3333"), fixed128.FromStringForced[T]("-0.3333").Neg())
	c.Equal(fixed128.FromStringForced[T]("2.6789"), fixed128.FromStringForced[T]("-2.6789").Neg())
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromInteger[T](-3).Neg())
}

func TestCmp(t *testing.T) {
	testCmp[fixed.D1](t)
	testCmp[fixed.D2](t)
	testCmp[fixed.D3](t)
	testCmp[fixed.D4](t)
	testCmp[fixed.D5](t)
	testCmp[fixed.D6](t)
}

func testCmp[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.Equal(1, fixed128.FromStringForced[T]("0.3333").Cmp(fixed128.FromInteger[T](-3)))
	c.Equal(-1, fixed128.FromStringForced[T]("2.6789").Cmp(fixed128.FromInteger[T](3)))
	c.Equal(0, fixed128.FromInteger[T](3).Cmp(fixed128.FromInteger[T](3)))
}

func TestEqual(t *testing.T) {
	testEqual[fixed.D1](t)
	testEqual[fixed.D2](t)
	testEqual[fixed.D3](t)
	testEqual[fixed.D4](t)
	testEqual[fixed.D5](t)
	testEqual[fixed.D6](t)
}

func testEqual[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.False(fixed128.FromStringForced[T]("0.3333").Equal(fixed128.FromInteger[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").Equal(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromInteger[T](3).Equal(fixed128.FromInteger[T](3)))
}

func TestGreaterThan(t *testing.T) {
	testGreaterThan[fixed.D1](t)
	testGreaterThan[fixed.D2](t)
	testGreaterThan[fixed.D3](t)
	testGreaterThan[fixed.D4](t)
	testGreaterThan[fixed.D5](t)
	testGreaterThan[fixed.D6](t)
}

func testGreaterThan[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.True(fixed128.FromStringForced[T]("0.3333").GreaterThan(fixed128.FromInteger[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").GreaterThan(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromInteger[T](3).GreaterThan(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromInteger[T](4).GreaterThan(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromStringForced[T]("2.6789").GreaterThan(fixed128.FromInteger[T](-1)))
}

func TestGreaterThanOrEqual(t *testing.T) {
	testGreaterThanOrEqual[fixed.D1](t)
	testGreaterThanOrEqual[fixed.D2](t)
	testGreaterThanOrEqual[fixed.D3](t)
	testGreaterThanOrEqual[fixed.D4](t)
	testGreaterThanOrEqual[fixed.D5](t)
	testGreaterThanOrEqual[fixed.D6](t)
}

func testGreaterThanOrEqual[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.True(fixed128.FromStringForced[T]("0.3333").GreaterThanOrEqual(fixed128.FromInteger[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").GreaterThanOrEqual(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromInteger[T](3).GreaterThanOrEqual(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromInteger[T](4).GreaterThanOrEqual(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromStringForced[T]("2.6789").GreaterThanOrEqual(fixed128.FromInteger[T](-1)))
}

func TestLessThan(t *testing.T) {
	testLessThan[fixed.D1](t)
	testLessThan[fixed.D2](t)
	testLessThan[fixed.D3](t)
	testLessThan[fixed.D4](t)
	testLessThan[fixed.D5](t)
	testLessThan[fixed.D6](t)
}

func testLessThan[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.False(fixed128.FromStringForced[T]("0.3333").LessThan(fixed128.FromInteger[T](-3)))
	c.True(fixed128.FromStringForced[T]("2.6789").LessThan(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromInteger[T](3).LessThan(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromInteger[T](4).LessThan(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromStringForced[T]("2.6789").LessThan(fixed128.FromInteger[T](-1)))
}

func TestLessThanOrEqual(t *testing.T) {
	testLessThanOrEqual[fixed.D1](t)
	testLessThanOrEqual[fixed.D2](t)
	testLessThanOrEqual[fixed.D3](t)
	testLessThanOrEqual[fixed.D4](t)
	testLessThanOrEqual[fixed.D5](t)
	testLessThanOrEqual[fixed.D6](t)
}

func testLessThanOrEqual[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	c.False(fixed128.FromStringForced[T]("0.3333").LessThanOrEqual(fixed128.FromInteger[T](-3)))
	c.True(fixed128.FromStringForced[T]("2.6789").LessThanOrEqual(fixed128.FromInteger[T](3)))
	c.True(fixed128.FromInteger[T](3).LessThanOrEqual(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromInteger[T](4).LessThanOrEqual(fixed128.FromInteger[T](3)))
	c.False(fixed128.FromStringForced[T]("2.6789").LessThanOrEqual(fixed128.FromInteger[T](-1)))
}

func TestComma(t *testing.T) {
	c := check.New(t)
	c.Equal("0.12", fixed128.FromStringForced[fixed.D2]("0.12").Comma())
	c.Equal("1,234,567,890.12", fixed128.FromStringForced[fixed.D2]("1234567890.12").Comma())
	c.Equal("91,234,567,890.12", fixed128.FromStringForced[fixed.D2]("91234567890.12").Comma())
	c.Equal("891,234,567,890.12", fixed128.FromStringForced[fixed.D2]("891234567890.12").Comma())
}

func TestJSON(t *testing.T) {
	testJSON[fixed.D1](t)
	testJSON[fixed.D2](t)
	testJSON[fixed.D3](t)
	testJSON[fixed.D4](t)
	testJSON[fixed.D5](t)
	testJSON[fixed.D6](t)
}

func testJSON[T fixed.Dx](t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testJSONActual(t, fixed128.FromInteger[T](i))
	}
	testJSONActual(t, fixed128.FromInteger[T, int64](1844674407371259000))
}

type embedded[T fixed.Dx] struct {
	Field fixed128.Int[T]
}

func testJSONActual[T fixed.Dx](t *testing.T, v fixed128.Int[T]) {
	c := check.New(t)
	c.Helper()
	e1 := embedded[T]{Field: v}
	data, err := json.Marshal(&e1)
	c.NoError(err)
	var e2 embedded[T]
	err = json.Unmarshal(data, &e2)
	c.NoError(err)
	c.Equal(e1, e2)
}

func TestYAML(t *testing.T) {
	testYAML[fixed.D1](t)
	testYAML[fixed.D2](t)
	testYAML[fixed.D3](t)
	testYAML[fixed.D4](t)
	testYAML[fixed.D5](t)
	testYAML[fixed.D6](t)
}

func testYAML[T fixed.Dx](t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testYAMLActual(t, fixed128.FromInteger[T](i))
	}
	testYAMLActual(t, fixed128.FromInteger[T, int64](1844674407371259000))
}

func testYAMLActual[T fixed.Dx](t *testing.T, v fixed128.Int[T]) {
	c := check.New(t)
	c.Helper()
	e1 := embedded[T]{Field: v}
	data, err := yaml.Marshal(&e1)
	c.NoError(err)
	var e2 embedded[T]
	err = yaml.Unmarshal(data, &e2)
	c.NoError(err)
	c.Equal(e1, e2)
}

func TestBoundaryValues(t *testing.T) {
	testBoundaryValues[fixed.D1](t)
	testBoundaryValues[fixed.D2](t)
	testBoundaryValues[fixed.D3](t)
	testBoundaryValues[fixed.D4](t)
	testBoundaryValues[fixed.D5](t)
	testBoundaryValues[fixed.D6](t)
}

func testBoundaryValues[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	maxValue := fixed128.Maximum[T]()
	minValue := fixed128.Minimum[T]()
	c.True(maxValue.GreaterThan(minValue))
	c.True(minValue.LessThan(maxValue))
}

func TestMinMax(t *testing.T) {
	testMinMax[fixed.D1](t)
	testMinMax[fixed.D2](t)
	testMinMax[fixed.D3](t)
	testMinMax[fixed.D4](t)
	testMinMax[fixed.D5](t)
	testMinMax[fixed.D6](t)
}

func testMinMax[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	a := fixed128.FromInteger[T](5)
	b := fixed128.FromInteger[T](10)
	negativeA := fixed128.FromInteger[T](-5)

	c.Equal(a, a.Min(b))
	c.Equal(a, b.Min(a))
	c.Equal(negativeA, negativeA.Min(a))
	c.Equal(a, a.Min(a))

	c.Equal(b, a.Max(b))
	c.Equal(b, b.Max(a))
	c.Equal(a, negativeA.Max(a))
	c.Equal(a, a.Max(a))
}

func TestIncDec(t *testing.T) {
	testIncDec[fixed.D1](t)
	testIncDec[fixed.D2](t)
	testIncDec[fixed.D3](t)
	testIncDec[fixed.D4](t)
	testIncDec[fixed.D5](t)
	testIncDec[fixed.D6](t)
}

func testIncDec[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	zero := fixed128.FromInteger[T](0)
	one := fixed128.FromInteger[T](1)
	negativeOne := fixed128.FromInteger[T](-1)

	c.Equal(one, zero.Inc())
	c.Equal(fixed128.FromInteger[T](2), one.Inc())
	c.Equal(zero, negativeOne.Inc())

	c.Equal(negativeOne, zero.Dec())
	c.Equal(zero, one.Dec())
	c.Equal(fixed128.FromInteger[T](-2), negativeOne.Dec())
}

func TestAs(t *testing.T) {
	testAs[fixed.D1](t)
	testAs[fixed.D2](t)
	testAs[fixed.D3](t)
	testAs[fixed.D4](t)
	testAs[fixed.D5](t)
	testAs[fixed.D6](t)
}

func testAs[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	intVal := fixed128.FromInteger[T](42)
	c.Equal(int(42), intVal.AsInteger[int]())
	c.Equal(int8(42), intVal.AsInteger[int8]())
	c.Equal(int16(42), intVal.AsInteger[int16]())
	c.Equal(int32(42), intVal.AsInteger[int32]())
	c.Equal(int64(42), intVal.AsInteger[int64]())
	c.Equal(uint(42), intVal.AsInteger[uint]())
	c.Equal(uint8(42), intVal.AsInteger[uint8]())
	c.Equal(uint16(42), intVal.AsInteger[uint16]())
	c.Equal(uint32(42), intVal.AsInteger[uint32]())
	c.Equal(uint64(42), intVal.AsInteger[uint64]())

	floatVal := fixed128.FromStringForced[T]("3.1")
	f32Result := floatVal.AsFloat[float32]()
	f64Result := floatVal.AsFloat[float64]()
	c.True(f32Result > 3.0 && f32Result < 3.2)
	c.True(f64Result > 3.0 && f64Result < 3.2)

	negVal := fixed128.FromInteger[T](-10)
	c.Equal(int(-10), negVal.AsInteger[int]())
	c.Equal(float64(-10.0), negVal.AsFloat[float64]())
}

func TestCheckedAs(t *testing.T) {
	testCheckedAs[fixed.D1](t)
	testCheckedAs[fixed.D2](t)
	testCheckedAs[fixed.D3](t)
	testCheckedAs[fixed.D4](t)
	testCheckedAs[fixed.D5](t)
	testCheckedAs[fixed.D6](t)
}

func testCheckedAs[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	intVal := fixed128.FromInteger[T](42)
	result, err := intVal.AsIntegerChecked[int]()
	c.NoError(err)
	c.Equal(int(42), result)

	floatResult, err := intVal.AsFloatChecked[float64]()
	c.NoError(err)
	c.Equal(float64(42.0), floatResult)
	fracVal := fixed128.FromStringForced[T]("42.5")
	_, err = fracVal.AsIntegerChecked[int]()
	c.HasError(err)

	floatResult, err = fracVal.AsFloatChecked[float64]()
	c.NoError(err)
	c.True(floatResult > 42.4 && floatResult < 42.6)

	// Values whose shortest float representation uses an exponent should still convert
	million := fixed128.FromInteger[T](1000000)
	floatResult, err = million.AsFloatChecked[float64]()
	c.NoError(err)
	c.Equal(float64(1000000), floatResult)
}

func TestStringWithSign(t *testing.T) {
	testStringWithSign[fixed.D1](t)
	testStringWithSign[fixed.D2](t)
	testStringWithSign[fixed.D3](t)
	testStringWithSign[fixed.D4](t)
	testStringWithSign[fixed.D5](t)
	testStringWithSign[fixed.D6](t)
}

func testStringWithSign[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	positive := fixed128.FromInteger[T](42)
	negative := fixed128.FromInteger[T](-42)
	zero := fixed128.FromInteger[T](0)

	c.Equal("+42", positive.StringWithSign())
	c.Equal("-42", negative.StringWithSign())
	c.Equal("+0", zero.StringWithSign())

	positiveFrac := fixed128.FromStringForced[T]("3.1")
	negativeFrac := fixed128.FromStringForced[T]("-3.1")

	c.Equal("+3.1", positiveFrac.StringWithSign())
	c.Equal("-3.1", negativeFrac.StringWithSign())
}

func TestCommaWithSign(t *testing.T) {
	testCommaWithSign[fixed.D1](t)
	testCommaWithSign[fixed.D2](t)
	testCommaWithSign[fixed.D3](t)
	testCommaWithSign[fixed.D4](t)
	testCommaWithSign[fixed.D5](t)
	testCommaWithSign[fixed.D6](t)
}

func testCommaWithSign[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	large := fixed128.FromStringForced[T]("1234567.8")
	largenegative := fixed128.FromStringForced[T]("-1234567.8")
	zero := fixed128.FromInteger[T](0)

	c.Equal("+1,234,567.8", large.CommaWithSign())
	c.Equal("-1,234,567.8", largenegative.CommaWithSign())
	c.Equal("+0", zero.CommaWithSign())
}

func TestMarshalText(t *testing.T) {
	testMarshalText[fixed.D1](t)
	testMarshalText[fixed.D2](t)
	testMarshalText[fixed.D3](t)
	testMarshalText[fixed.D4](t)
	testMarshalText[fixed.D5](t)
	testMarshalText[fixed.D6](t)
}

func testMarshalText[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	val := fixed128.FromStringForced[T]("123.4")

	data, err := val.MarshalText()
	c.NoError(err)
	c.Equal("123.4", string(data))

	var unmarshaled fixed128.Int[T]
	err = unmarshaled.UnmarshalText(data)
	c.NoError(err)
	c.Equal(val, unmarshaled)

	err = unmarshaled.UnmarshalText([]byte(`"123.4"`))
	c.NoError(err)
	c.Equal(val, unmarshaled)
}

func TestUnmarshalErrors(t *testing.T) {
	testUnmarshalErrors[fixed.D1](t)
	testUnmarshalErrors[fixed.D2](t)
	testUnmarshalErrors[fixed.D3](t)
	testUnmarshalErrors[fixed.D4](t)
	testUnmarshalErrors[fixed.D5](t)
	testUnmarshalErrors[fixed.D6](t)
}

func testUnmarshalErrors[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	var val fixed128.Int[T]
	err := val.UnmarshalJSON([]byte("invalid"))
	c.HasError(err)

	err = val.UnmarshalText([]byte("invalid"))
	c.HasError(err)

	err = val.UnmarshalYAML(func(any) error {
		return fmt.Errorf("test error")
	})
	c.HasError(err)
}

func TestFromStringEdgeCases(t *testing.T) {
	testFromStringEdgeCases[fixed.D1](t)
	testFromStringEdgeCases[fixed.D2](t)
	testFromStringEdgeCases[fixed.D3](t)
	testFromStringEdgeCases[fixed.D4](t)
	testFromStringEdgeCases[fixed.D5](t)
	testFromStringEdgeCases[fixed.D6](t)
}

func testFromStringEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	_, err := fixed128.FromString[T]("")
	c.HasError(err)

	val, err := fixed128.FromString[T]("1,234.5")
	c.NoError(err)
	c.Equal("1234.5", val.String())

	val, err = fixed128.FromString[T]("1.23e2")
	c.NoError(err)
	c.Equal("123", val.String())

	val, err = fixed128.FromString[T]("-1.23E2")
	c.NoError(err)
	c.Equal("-123", val.String())

	_, err = fixed128.FromString[T]("1.23ez")
	c.HasError(err)

	_, err = fixed128.FromString[T]("abc.123")
	c.HasError(err)

	_, err = fixed128.FromString[T]("123.abc")
	c.HasError(err)

	val, err = fixed128.FromString[T]("-0")
	c.NoError(err)
	c.Equal("0", val.String())

	val, err = fixed128.FromString[T]("-0.000")
	c.NoError(err)
	c.Equal("0", val.String())

	val, err = fixed128.FromString[T](".5")
	c.NoError(err)
	c.Equal("0.5", val.String())

	val, err = fixed128.FromString[T]("-.5")
	c.NoError(err)
	c.Equal("-0.5", val.String())

	// Excess fractional digits are truncated, not rejected.
	val, err = fixed128.FromString[T]("0.123456789012345678901234567890")
	c.NoError(err)
	c.NotEqual("", val.String())
	c.HasPrefix(val.String(), "0.1")
}

func TestFloorCeilRoundExtremes(t *testing.T) {
	testFloorCeilRoundExtremes[fixed.D1](t)
	testFloorCeilRoundExtremes[fixed.D2](t)
	testFloorCeilRoundExtremes[fixed.D3](t)
	testFloorCeilRoundExtremes[fixed.D4](t)
	testFloorCeilRoundExtremes[fixed.D5](t)
	testFloorCeilRoundExtremes[fixed.D6](t)
}

func testFloorCeilRoundExtremes[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	zero := fixed128.FromInteger[T](0)
	// The floor of Minimum() and ceiling of Maximum() are not representable, so they saturate rather than wrapping.
	c.Equal(fixed128.Minimum[T](), fixed128.Minimum[T]().Floor())
	c.Equal(fixed128.Maximum[T](), fixed128.Maximum[T]().Ceil())
	// The extreme that already rounds toward zero stays representable.
	c.True(fixed128.Minimum[T]().Ceil().GreaterThan(fixed128.Minimum[T]()))
	c.True(fixed128.Maximum[T]().Floor().LessThan(fixed128.Maximum[T]()))
	// Whether Round() of an extreme saturates depends on the precision, but it must never wrap around.
	c.True(fixed128.Minimum[T]().Round().LessThan(zero))
	c.True(fixed128.Minimum[T]().Round().GreaterThanOrEqual(fixed128.Minimum[T]()))
	c.True(fixed128.Maximum[T]().Round().GreaterThan(zero))
	c.True(fixed128.Maximum[T]().Round().LessThanOrEqual(fixed128.Maximum[T]()))
}

func TestCeilEdgeCases(t *testing.T) {
	testCeilEdgeCases[fixed.D1](t)
	testCeilEdgeCases[fixed.D2](t)
	testCeilEdgeCases[fixed.D3](t)
	testCeilEdgeCases[fixed.D4](t)
	testCeilEdgeCases[fixed.D5](t)
	testCeilEdgeCases[fixed.D6](t)
}

func testCeilEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	negFrac := fixed128.FromStringForced[T]("-2.5")
	c.Equal(fixed128.FromInteger[T](-2), negFrac.Ceil())

	zero := fixed128.FromInteger[T](0)
	c.Equal(zero, zero.Ceil())

	negZero := fixed128.FromStringForced[T]("-0.0")
	c.Equal(fixed128.FromInteger[T](0), negZero.Ceil())
}

func TestStringEdgeCases(t *testing.T) {
	testStringEdgeCases[fixed.D1](t)
	testStringEdgeCases[fixed.D2](t)
	testStringEdgeCases[fixed.D3](t)
	testStringEdgeCases[fixed.D4](t)
	testStringEdgeCases[fixed.D5](t)
	testStringEdgeCases[fixed.D6](t)
}

func testStringEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	negFrac := fixed128.FromStringForced[T]("-0.5")
	c.Equal("-0.5", negFrac.String())

	posFrac := fixed128.FromStringForced[T]("0.5")
	c.Equal("0.5", posFrac.String())

	val := fixed128.FromStringForced[T]("1.2000")
	c.Equal("1.2", val.String())
}

func TestAdditionalEdgeCases(t *testing.T) {
	testAdditionalEdgeCases[fixed.D1](t)
	testAdditionalEdgeCases[fixed.D2](t)
	testAdditionalEdgeCases[fixed.D3](t)
	testAdditionalEdgeCases[fixed.D4](t)
	testAdditionalEdgeCases[fixed.D5](t)
	testAdditionalEdgeCases[fixed.D6](t)
}

func testAdditionalEdgeCases[T fixed.Dx](t *testing.T) {
	c := check.New(t)

	val := fixed128.FromStringForced[T]("999999999999999999999999999.9")
	_, _ = val.AsFloatChecked[float32]() //nolint:errcheck // This might succeed or fail depending on precision, but shouldn't panic. We'll just test that it doesn't panic
	c.NotNil(val)

	var intVal fixed128.Int[T]
	err := intVal.UnmarshalYAML(func(v any) error {
		*v.(*string) = "42" //nolint:errcheck // This is just a test, we know it will succeed
		return nil
	})
	c.NoError(err)

	err = intVal.UnmarshalYAML(func(any) error {
		return fmt.Errorf("unmarshal error")
	})
	c.HasError(err)
}

// TestOutOfRangeSaturates verifies that a value too large for the type converts to Maximum/Minimum rather than 0, and
// that an infinity or NaN does not panic. big.Float.SetFloat64 panics on NaN, and an infinity used to become 0.
func TestOutOfRangeSaturates(t *testing.T) {
	testOutOfRangeSaturates[fixed.D1](t)
	testOutOfRangeSaturates[fixed.D2](t)
	testOutOfRangeSaturates[fixed.D3](t)
	testOutOfRangeSaturates[fixed.D4](t)
	testOutOfRangeSaturates[fixed.D5](t)
	testOutOfRangeSaturates[fixed.D6](t)
}

func testOutOfRangeSaturates[T fixed.Dx](t *testing.T) {
	c := check.New(t)
	maximum := fixed128.Maximum[T]()
	minimum := fixed128.Minimum[T]()

	c.Equal(maximum, fixed128.FromFloat[T](math.MaxFloat64))
	c.Equal(minimum, fixed128.FromFloat[T](-math.MaxFloat64))
	c.Equal(maximum, fixed128.FromFloat[T](math.Inf(1)))
	c.Equal(minimum, fixed128.FromFloat[T](math.Inf(-1)))
	c.NotPanics(func() { c.Equal(fixed128.Int[T]{}, fixed128.FromFloat[T](math.NaN())) })
	c.Equal(fixed128.FromInteger[T](3), fixed128.FromFloat[T](3.0))
}

// TestFromStringOutOfRangeReportsAndSaturates verifies that FromString reports a value beyond what the type can hold.
// num128 saturates silently, so such a string used to come back as the bound with no error.
func TestFromStringOutOfRangeReportsAndSaturates(t *testing.T) {
	c := check.New(t)
	maximum := fixed128.Maximum[fixed.D4]()
	minimum := fixed128.Minimum[fixed.D4]()

	for _, tc := range []struct {
		str  string
		want fixed128.Int[fixed.D4]
	}{
		{str: "1e300", want: maximum},
		{str: "-1e300", want: minimum},
		{str: "1e400", want: maximum},
		{str: "-1e400", want: minimum},
		{str: "99999999999999999999999999999999999999999", want: maximum},
		{str: "-99999999999999999999999999999999999999999", want: minimum},
	} {
		v, err := fixed128.FromString[fixed.D4](tc.str)
		c.HasError(err, "%q should be reported as out of range", tc.str)
		c.Equal(tc.want, v, "%q", tc.str)
		c.Equal(tc.want, fixed128.FromStringForced[fixed.D4](tc.str), "%q forced", tc.str)
	}

	// A malformed string has no nearest value, so it stays 0.
	for _, str := range []string{"abc", ""} {
		v, err := fixed128.FromString[fixed.D4](str)
		c.HasError(err, "%q should be rejected", str)
		c.Equal(fixed128.Int[fixed.D4]{}, v, "%q", str)
	}

	// Values within range are unaffected, including the extremes themselves.
	for _, tc := range []struct {
		str  string
		want fixed128.Int[fixed.D4]
	}{
		{str: "1e5", want: fixed128.FromInteger[fixed.D4](100000)},
		{str: maximum.String(), want: maximum},
		{str: minimum.String(), want: minimum},
	} {
		v, err := fixed128.FromString[fixed.D4](tc.str)
		c.NoError(err, "%q", tc.str)
		c.Equal(tc.want, v, "%q", tc.str)
	}
}
