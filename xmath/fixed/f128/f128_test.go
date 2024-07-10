// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f128_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128"
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
	check.Equal(t, "0.1", f128.From[T, float64](0.1).String())
	check.Equal(t, "0.2", f128.From[T, float64](0.2).String())
	check.Equal(t, "0.3", f128.FromStringForced[T]("0.3").String())
	check.Equal(t, "-0.1", f128.From[T, float64](-0.1).String())
	check.Equal(t, "-0.2", f128.From[T, float64](-0.2).String())
	check.Equal(t, "-0.3", f128.FromStringForced[T]("-0.3").String())
	threeFill := strings.Repeat("3", f128.MaxDecimalDigits[T]())
	check.Equal(t, "0."+threeFill, f128.FromStringForced[T]("0.33333333").String())
	check.Equal(t, "-0."+threeFill, f128.FromStringForced[T]("-0.33333333").String())
	sixFill := strings.Repeat("6", f128.MaxDecimalDigits[T]())
	check.Equal(t, "0."+sixFill, f128.FromStringForced[T]("0.66666666").String())
	check.Equal(t, "-0."+sixFill, f128.FromStringForced[T]("-0.66666666").String())
	check.Equal(t, "1", f128.From[T, float64](1.0000004).String())
	check.Equal(t, "1", f128.From[T, float64](1.00000049).String())
	check.Equal(t, "1", f128.From[T, float64](1.0000005).String())
	check.Equal(t, "1", f128.From[T, float64](1.0000009).String())
	check.Equal(t, "-1", f128.From[T, float64](-1.0000004).String())
	check.Equal(t, "-1", f128.From[T, float64](-1.00000049).String())
	check.Equal(t, "-1", f128.From[T, float64](-1.0000005).String())
	check.Equal(t, "-1", f128.From[T, float64](-1.0000009).String())
	zeroFill := strings.Repeat("0", f128.MaxDecimalDigits[T]()-1)
	check.Equal(t, "0."+zeroFill+"4", f128.FromStringForced[T]("0."+zeroFill+"405").String())
	check.Equal(t, "-0."+zeroFill+"4", f128.FromStringForced[T]("-0."+zeroFill+"405").String())

	v, err := f128.FromString[T]("33.0")
	check.NoError(t, err)
	check.Equal(t, v, f128.From[T, int](33))

	v, err = f128.FromString[T]("33.00000000000000000000")
	check.NoError(t, err)
	check.Equal(t, v, f128.From[T, int](33))
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
	oneThird := f128.FromStringForced[T]("0.333333")
	negTwoThirds := f128.FromStringForced[T]("-0.666666")
	one := f128.From[T, int](1)
	oneAndTwoThirds := f128.FromStringForced[T]("1.666666")
	nineThousandSix := f128.From[T, int](9006)
	two := f128.From[T, int](2)
	check.Equal(t, "0."+strings.Repeat("9", f128.MaxDecimalDigits[T]()), oneThird.Add(oneThird).Add(oneThird).String())
	check.Equal(t, "0."+strings.Repeat("6", f128.MaxDecimalDigits[T]()-1)+"7", one.Sub(oneThird).String())
	check.Equal(t, "-1."+strings.Repeat("6", f128.MaxDecimalDigits[T]()), negTwoThirds.Sub(one).String())
	check.Equal(t, "0", negTwoThirds.Sub(one).Add(oneAndTwoThirds).String())
	check.Equal(t, f128.From[T, int](10240), f128.From[T, int](1234).Add(nineThousandSix))
	check.Equal(t, "10240", f128.From[T, int](1234).Add(nineThousandSix).String())
	check.Equal(t, "-1.5", f128.From[T, float64](0.5).Sub(two).String())
	ninetyPointZeroSix := f128.FromStringForced[T]("90.06")
	twelvePointThirtyFour := f128.FromStringForced[T]("12.34")
	var answer string
	if f128.MaxDecimalDigits[T]() > 1 {
		answer = "102.4"
	} else {
		answer = "102.3"
	}
	check.Equal(t, f128.FromStringForced[T](answer), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	check.Equal(t, answer, twelvePointThirtyFour.Add(ninetyPointZeroSix).String())
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
	pointThree := f128.FromStringForced[T]("0.3")
	negativePointThree := f128.FromStringForced[T]("-0.3")
	threeFill := strings.Repeat("3", f128.MaxDecimalDigits[T]())
	check.Equal(t, "0."+threeFill, f128.From[T, int](1).Div(f128.From[T, int](3)).String())
	check.Equal(t, "-0."+threeFill, f128.From[T, int](1).Div(f128.From[T, int](-3)).String())
	check.Equal(t, "0.1", pointThree.Div(f128.From[T, int](3)).String())
	check.Equal(t, "0.9", pointThree.Mul(f128.From[T, int](3)).String())
	check.Equal(t, "-0.9", negativePointThree.Mul(f128.From[T, int](3)).String())
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
	check.Equal(t, f128.From[T, int](1), f128.From[T, int](3).Mod(f128.From[T, int](2)))
	check.Equal(t, f128.FromStringForced[T]("0.3"), f128.FromStringForced[T]("9.3").Mod(f128.From[T, int](3)))
	check.Equal(t, f128.FromStringForced[T]("0.1"), f128.FromStringForced[T]("3.1").Mod(f128.FromStringForced[T]("0.2")))
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
	check.Equal(t, f128.From[T, int](0), f128.FromStringForced[T]("0.3333").Trunc())
	check.Equal(t, f128.From[T, int](2), f128.FromStringForced[T]("2.6789").Trunc())
	check.Equal(t, f128.From[T, int](3), f128.From[T, int](3).Trunc())
	check.Equal(t, f128.From[T, int](0), f128.FromStringForced[T]("-0.3333").Trunc())
	check.Equal(t, f128.From[T, int](-2), f128.FromStringForced[T]("-2.6789").Trunc())
	check.Equal(t, f128.From[T, int](-3), f128.From[T, int](-3).Trunc())
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
	check.Equal(t, f128.From[T, int](1), f128.FromStringForced[T]("0.3333").Ceil())
	check.Equal(t, f128.From[T, int](3), f128.FromStringForced[T]("2.6789").Ceil())
	check.Equal(t, f128.From[T, int](3), f128.From[T, int](3).Ceil())
	check.Equal(t, f128.From[T, int](0), f128.FromStringForced[T]("-0.3333").Ceil())
	check.Equal(t, f128.From[T, int](-2), f128.FromStringForced[T]("-2.6789").Ceil())
	check.Equal(t, f128.From[T, int](-3), f128.From[T, int](-3).Ceil())
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
	check.Equal(t, f128.From[T, int](0), f128.FromStringForced[T]("0.3333").Round())
	check.Equal(t, f128.From[T, int](3), f128.FromStringForced[T]("2.6789").Round())
	check.Equal(t, f128.From[T, int](3), f128.From[T, int](3).Round())
	check.Equal(t, f128.From[T, int](0), f128.FromStringForced[T]("-0.3333").Round())
	check.Equal(t, f128.From[T, int](-3), f128.FromStringForced[T]("-2.6789").Round())
	check.Equal(t, f128.From[T, int](-3), f128.From[T, int](-3).Round())
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
	check.Equal(t, f128.FromStringForced[T]("0.3333"), f128.FromStringForced[T]("0.3333").Abs())
	check.Equal(t, f128.FromStringForced[T]("2.6789"), f128.FromStringForced[T]("2.6789").Abs())
	check.Equal(t, f128.From[T, int](3), f128.From[T, int](3).Abs())
	check.Equal(t, f128.FromStringForced[T]("0.3333"), f128.FromStringForced[T]("-0.3333").Abs())
	check.Equal(t, f128.FromStringForced[T]("2.6789"), f128.FromStringForced[T]("-2.6789").Abs())
	check.Equal(t, f128.From[T, int](3), f128.From[T, int](-3).Abs())
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
	check.Equal(t, f128.FromStringForced[T]("-0.3333"), f128.FromStringForced[T]("0.3333").Neg())
	check.Equal(t, f128.FromStringForced[T]("-2.6789"), f128.FromStringForced[T]("2.6789").Neg())
	check.Equal(t, f128.From[T](-3), f128.From[T](3).Neg())
	check.Equal(t, f128.FromStringForced[T]("0.3333"), f128.FromStringForced[T]("-0.3333").Neg())
	check.Equal(t, f128.FromStringForced[T]("2.6789"), f128.FromStringForced[T]("-2.6789").Neg())
	check.Equal(t, f128.From[T](3), f128.From[T](-3).Neg())
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
	check.Equal(t, 1, f128.FromStringForced[T]("0.3333").Cmp(f128.From[T](-3)))
	check.Equal(t, -1, f128.FromStringForced[T]("2.6789").Cmp(f128.From[T](3)))
	check.Equal(t, 0, f128.From[T](3).Cmp(f128.From[T](3)))
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
	check.Equal(t, false, f128.FromStringForced[T]("0.3333").Equal(f128.From[T](-3)))
	check.Equal(t, false, f128.FromStringForced[T]("2.6789").Equal(f128.From[T](3)))
	check.Equal(t, true, f128.From[T](3).Equal(f128.From[T](3)))
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
	check.Equal(t, true, f128.FromStringForced[T]("0.3333").GreaterThan(f128.From[T](-3)))
	check.Equal(t, false, f128.FromStringForced[T]("2.6789").GreaterThan(f128.From[T](3)))
	check.Equal(t, false, f128.From[T](3).GreaterThan(f128.From[T](3)))
	check.Equal(t, true, f128.From[T](4).GreaterThan(f128.From[T](3)))
	check.Equal(t, true, f128.FromStringForced[T]("2.6789").GreaterThan(f128.From[T](-1)))
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
	check.Equal(t, true, f128.FromStringForced[T]("0.3333").GreaterThanOrEqual(f128.From[T](-3)))
	check.Equal(t, false, f128.FromStringForced[T]("2.6789").GreaterThanOrEqual(f128.From[T](3)))
	check.Equal(t, true, f128.From[T](3).GreaterThanOrEqual(f128.From[T](3)))
	check.Equal(t, true, f128.From[T](4).GreaterThanOrEqual(f128.From[T](3)))
	check.Equal(t, true, f128.FromStringForced[T]("2.6789").GreaterThanOrEqual(f128.From[T](-1)))
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
	check.Equal(t, false, f128.FromStringForced[T]("0.3333").LessThan(f128.From[T](-3)))
	check.Equal(t, true, f128.FromStringForced[T]("2.6789").LessThan(f128.From[T](3)))
	check.Equal(t, false, f128.From[T](3).LessThan(f128.From[T](3)))
	check.Equal(t, false, f128.From[T](4).LessThan(f128.From[T](3)))
	check.Equal(t, false, f128.FromStringForced[T]("2.6789").LessThan(f128.From[T](-1)))
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
	check.Equal(t, false, f128.FromStringForced[T]("0.3333").LessThanOrEqual(f128.From[T](-3)))
	check.Equal(t, true, f128.FromStringForced[T]("2.6789").LessThanOrEqual(f128.From[T](3)))
	check.Equal(t, true, f128.From[T](3).LessThanOrEqual(f128.From[T](3)))
	check.Equal(t, false, f128.From[T](4).LessThanOrEqual(f128.From[T](3)))
	check.Equal(t, false, f128.FromStringForced[T]("2.6789").LessThanOrEqual(f128.From[T](-1)))
}

func TestComma(t *testing.T) {
	check.Equal(t, "0.12", f128.FromStringForced[fixed.D2]("0.12").Comma())
	check.Equal(t, "1,234,567,890.12", f128.FromStringForced[fixed.D2]("1234567890.12").Comma())
	check.Equal(t, "91,234,567,890.12", f128.FromStringForced[fixed.D2]("91234567890.12").Comma())
	check.Equal(t, "891,234,567,890.12", f128.FromStringForced[fixed.D2]("891234567890.12").Comma())
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
		testJSONActual(t, f128.From[T](i))
	}
	testJSONActual(t, f128.From[T, int64](1844674407371259000))
}

type embedded[T fixed.Dx] struct {
	Field f128.Int[T]
}

func testJSONActual[T fixed.Dx](t *testing.T, v f128.Int[T]) {
	t.Helper()
	e1 := embedded[T]{Field: v}
	data, err := json.Marshal(&e1)
	check.NoError(t, err)
	var e2 embedded[T]
	err = json.Unmarshal(data, &e2)
	check.NoError(t, err)
	check.Equal(t, e1, e2)
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
		testYAMLActual(t, f128.From[T](i))
	}
	testYAMLActual(t, f128.From[T, int64](1844674407371259000))
}

func testYAMLActual[T fixed.Dx](t *testing.T, v f128.Int[T]) {
	t.Helper()
	e1 := embedded[T]{Field: v}
	data, err := yaml.Marshal(&e1)
	check.NoError(t, err)
	var e2 embedded[T]
	err = yaml.Unmarshal(data, &e2)
	check.NoError(t, err)
	check.Equal(t, e1, e2)
}
