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
	"gopkg.in/yaml.v3"
)

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
	c.Equal("0.1", fixed128.From[T, float64](0.1).String())
	c.Equal("0.2", fixed128.From[T, float64](0.2).String())
	c.Equal("0.3", fixed128.FromStringForced[T]("0.3").String())
	c.Equal("-0.1", fixed128.From[T, float64](-0.1).String())
	c.Equal("-0.2", fixed128.From[T, float64](-0.2).String())
	c.Equal("-0.3", fixed128.FromStringForced[T]("-0.3").String())
	threeFill := strings.Repeat("3", fixed128.MaxDecimalDigits[T]())
	c.Equal("0."+threeFill, fixed128.FromStringForced[T]("0.33333333").String())
	c.Equal("-0."+threeFill, fixed128.FromStringForced[T]("-0.33333333").String())
	sixFill := strings.Repeat("6", fixed128.MaxDecimalDigits[T]())
	c.Equal("0."+sixFill, fixed128.FromStringForced[T]("0.66666666").String())
	c.Equal("-0."+sixFill, fixed128.FromStringForced[T]("-0.66666666").String())
	c.Equal("1", fixed128.From[T, float64](1.0000004).String())
	c.Equal("1", fixed128.From[T, float64](1.00000049).String())
	c.Equal("1", fixed128.From[T, float64](1.0000005).String())
	c.Equal("1", fixed128.From[T, float64](1.0000009).String())
	c.Equal("-1", fixed128.From[T, float64](-1.0000004).String())
	c.Equal("-1", fixed128.From[T, float64](-1.00000049).String())
	c.Equal("-1", fixed128.From[T, float64](-1.0000005).String())
	c.Equal("-1", fixed128.From[T, float64](-1.0000009).String())
	zeroFill := strings.Repeat("0", fixed128.MaxDecimalDigits[T]()-1)
	c.Equal("0."+zeroFill+"4", fixed128.FromStringForced[T]("0."+zeroFill+"405").String())
	c.Equal("-0."+zeroFill+"4", fixed128.FromStringForced[T]("-0."+zeroFill+"405").String())

	v, err := fixed128.FromString[T]("33.0")
	c.NoError(err)
	c.Equal(v, fixed128.From[T, int](33))

	v, err = fixed128.FromString[T]("33.00000000000000000000")
	c.NoError(err)
	c.Equal(v, fixed128.From[T, int](33))
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
	one := fixed128.From[T, int](1)
	oneAndTwoThirds := fixed128.FromStringForced[T]("1.666666")
	nineThousandSix := fixed128.From[T, int](9006)
	two := fixed128.From[T, int](2)
	c := check.New(t)
	c.Equal("0."+strings.Repeat("9", fixed128.MaxDecimalDigits[T]()), oneThird.Add(oneThird).Add(oneThird).String())
	c.Equal("0."+strings.Repeat("6", fixed128.MaxDecimalDigits[T]()-1)+"7", one.Sub(oneThird).String())
	c.Equal("-1."+strings.Repeat("6", fixed128.MaxDecimalDigits[T]()), negTwoThirds.Sub(one).String())
	c.Equal("0", negTwoThirds.Sub(one).Add(oneAndTwoThirds).String())
	c.Equal(fixed128.From[T, int](10240), fixed128.From[T, int](1234).Add(nineThousandSix))
	c.Equal("10240", fixed128.From[T, int](1234).Add(nineThousandSix).String())
	c.Equal("-1.5", fixed128.From[T, float64](0.5).Sub(two).String())
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
	c.Equal("0."+threeFill, fixed128.From[T, int](1).Div(fixed128.From[T, int](3)).String())
	c.Equal("-0."+threeFill, fixed128.From[T, int](1).Div(fixed128.From[T, int](-3)).String())
	c.Equal("0.1", pointThree.Div(fixed128.From[T, int](3)).String())
	c.Equal("0.9", pointThree.Mul(fixed128.From[T, int](3)).String())
	c.Equal("-0.9", negativePointThree.Mul(fixed128.From[T, int](3)).String())
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
	c.Equal(fixed128.From[T, int](1), fixed128.From[T, int](3).Mod(fixed128.From[T, int](2)))
	c.Equal(fixed128.FromStringForced[T]("0.3"), fixed128.FromStringForced[T]("9.3").Mod(fixed128.From[T, int](3)))
	c.Equal(fixed128.FromStringForced[T]("0.1"), fixed128.FromStringForced[T]("3.1").Mod(fixed128.FromStringForced[T]("0.2")))
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
	c.Equal(fixed128.From[T, int](0), fixed128.FromStringForced[T]("0.3333").Trunc())
	c.Equal(fixed128.From[T, int](2), fixed128.FromStringForced[T]("2.6789").Trunc())
	c.Equal(fixed128.From[T, int](3), fixed128.From[T, int](3).Trunc())
	c.Equal(fixed128.From[T, int](0), fixed128.FromStringForced[T]("-0.3333").Trunc())
	c.Equal(fixed128.From[T, int](-2), fixed128.FromStringForced[T]("-2.6789").Trunc())
	c.Equal(fixed128.From[T, int](-3), fixed128.From[T, int](-3).Trunc())
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
	c.Equal(fixed128.From[T, int](1), fixed128.FromStringForced[T]("0.3333").Ceil())
	c.Equal(fixed128.From[T, int](3), fixed128.FromStringForced[T]("2.6789").Ceil())
	c.Equal(fixed128.From[T, int](3), fixed128.From[T, int](3).Ceil())
	c.Equal(fixed128.From[T, int](0), fixed128.FromStringForced[T]("-0.3333").Ceil())
	c.Equal(fixed128.From[T, int](-2), fixed128.FromStringForced[T]("-2.6789").Ceil())
	c.Equal(fixed128.From[T, int](-3), fixed128.From[T, int](-3).Ceil())
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
	c.Equal(fixed128.From[T, int](0), fixed128.FromStringForced[T]("0.3333").Round())
	c.Equal(fixed128.From[T, int](3), fixed128.FromStringForced[T]("2.6789").Round())
	c.Equal(fixed128.From[T, int](3), fixed128.From[T, int](3).Round())
	c.Equal(fixed128.From[T, int](0), fixed128.FromStringForced[T]("-0.3333").Round())
	c.Equal(fixed128.From[T, int](-3), fixed128.FromStringForced[T]("-2.6789").Round())
	c.Equal(fixed128.From[T, int](-3), fixed128.From[T, int](-3).Round())
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
	c.Equal(fixed128.From[T, int](3), fixed128.From[T, int](3).Abs())
	c.Equal(fixed128.FromStringForced[T]("0.3333"), fixed128.FromStringForced[T]("-0.3333").Abs())
	c.Equal(fixed128.FromStringForced[T]("2.6789"), fixed128.FromStringForced[T]("-2.6789").Abs())
	c.Equal(fixed128.From[T, int](3), fixed128.From[T, int](-3).Abs())
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
	c.Equal(fixed128.From[T](-3), fixed128.From[T](3).Neg())
	c.Equal(fixed128.FromStringForced[T]("0.3333"), fixed128.FromStringForced[T]("-0.3333").Neg())
	c.Equal(fixed128.FromStringForced[T]("2.6789"), fixed128.FromStringForced[T]("-2.6789").Neg())
	c.Equal(fixed128.From[T](3), fixed128.From[T](-3).Neg())
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
	c.Equal(1, fixed128.FromStringForced[T]("0.3333").Cmp(fixed128.From[T](-3)))
	c.Equal(-1, fixed128.FromStringForced[T]("2.6789").Cmp(fixed128.From[T](3)))
	c.Equal(0, fixed128.From[T](3).Cmp(fixed128.From[T](3)))
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
	c.False(fixed128.FromStringForced[T]("0.3333").Equal(fixed128.From[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").Equal(fixed128.From[T](3)))
	c.True(fixed128.From[T](3).Equal(fixed128.From[T](3)))
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
	c.True(fixed128.FromStringForced[T]("0.3333").GreaterThan(fixed128.From[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").GreaterThan(fixed128.From[T](3)))
	c.False(fixed128.From[T](3).GreaterThan(fixed128.From[T](3)))
	c.True(fixed128.From[T](4).GreaterThan(fixed128.From[T](3)))
	c.True(fixed128.FromStringForced[T]("2.6789").GreaterThan(fixed128.From[T](-1)))
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
	c.True(fixed128.FromStringForced[T]("0.3333").GreaterThanOrEqual(fixed128.From[T](-3)))
	c.False(fixed128.FromStringForced[T]("2.6789").GreaterThanOrEqual(fixed128.From[T](3)))
	c.True(fixed128.From[T](3).GreaterThanOrEqual(fixed128.From[T](3)))
	c.True(fixed128.From[T](4).GreaterThanOrEqual(fixed128.From[T](3)))
	c.True(fixed128.FromStringForced[T]("2.6789").GreaterThanOrEqual(fixed128.From[T](-1)))
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
	c.False(fixed128.FromStringForced[T]("0.3333").LessThan(fixed128.From[T](-3)))
	c.True(fixed128.FromStringForced[T]("2.6789").LessThan(fixed128.From[T](3)))
	c.False(fixed128.From[T](3).LessThan(fixed128.From[T](3)))
	c.False(fixed128.From[T](4).LessThan(fixed128.From[T](3)))
	c.False(fixed128.FromStringForced[T]("2.6789").LessThan(fixed128.From[T](-1)))
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
	c.False(fixed128.FromStringForced[T]("0.3333").LessThanOrEqual(fixed128.From[T](-3)))
	c.True(fixed128.FromStringForced[T]("2.6789").LessThanOrEqual(fixed128.From[T](3)))
	c.True(fixed128.From[T](3).LessThanOrEqual(fixed128.From[T](3)))
	c.False(fixed128.From[T](4).LessThanOrEqual(fixed128.From[T](3)))
	c.False(fixed128.FromStringForced[T]("2.6789").LessThanOrEqual(fixed128.From[T](-1)))
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
		testJSONActual(t, fixed128.From[T](i))
	}
	testJSONActual(t, fixed128.From[T, int64](1844674407371259000))
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
		testYAMLActual(t, fixed128.From[T](i))
	}
	testYAMLActual(t, fixed128.From[T, int64](1844674407371259000))
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
