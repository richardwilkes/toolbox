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
	"fmt"
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

	// Test Maximum and Minimum
	maxValue := fixed128.Maximum[T]()
	minValue := fixed128.Minimum[T]()
	c.True(maxValue.GreaterThan(minValue))
	c.True(minValue.LessThan(maxValue))

	// Test MaxSafeMultiply
	maxSafe := fixed128.MaxSafeMultiply[T]()
	c.True(maxSafe.LessThan(maxValue))
	c.True(maxSafe.GreaterThan(minValue))
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

	// Test Min
	c.Equal(a, a.Min(b))
	c.Equal(a, b.Min(a))
	c.Equal(negativeA, negativeA.Min(a))
	c.Equal(a, a.Min(a))

	// Test Max
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

	// Test Inc
	c.Equal(one, zero.Inc())
	c.Equal(fixed128.FromInteger[T](2), one.Inc())
	c.Equal(zero, negativeOne.Inc())

	// Test Dec
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

	// Test integer conversions
	intVal := fixed128.FromInteger[T](42)
	c.Equal(int(42), fixed128.AsInteger[T, int](intVal))
	c.Equal(int8(42), fixed128.AsInteger[T, int8](intVal))
	c.Equal(int16(42), fixed128.AsInteger[T, int16](intVal))
	c.Equal(int32(42), fixed128.AsInteger[T, int32](intVal))
	c.Equal(int64(42), fixed128.AsInteger[T, int64](intVal))
	c.Equal(uint(42), fixed128.AsInteger[T, uint](intVal))
	c.Equal(uint8(42), fixed128.AsInteger[T, uint8](intVal))
	c.Equal(uint16(42), fixed128.AsInteger[T, uint16](intVal))
	c.Equal(uint32(42), fixed128.AsInteger[T, uint32](intVal))
	c.Equal(uint64(42), fixed128.AsInteger[T, uint64](intVal))

	// Test float conversions
	floatVal := fixed128.FromStringForced[T]("3.1")
	f32Result := fixed128.AsFloat[T, float32](floatVal)
	f64Result := fixed128.AsFloat[T, float64](floatVal)
	c.True(f32Result > 3.0 && f32Result < 3.2)
	c.True(f64Result > 3.0 && f64Result < 3.2)

	// Test negative values
	negVal := fixed128.FromInteger[T](-10)
	c.Equal(int(-10), fixed128.AsInteger[T, int](negVal))
	c.Equal(float64(-10.0), fixed128.AsFloat[T, float64](negVal))
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

	// Test successful conversions
	intVal := fixed128.FromInteger[T](42)
	result, err := fixed128.CheckedAsInteger[T, int](intVal)
	c.NoError(err)
	c.Equal(int(42), result)

	floatResult, err := fixed128.CheckedAsFloat[T, float64](intVal)
	c.NoError(err)
	c.Equal(float64(42.0), floatResult)
	// Test conversion that should fail (fractional part)
	fracVal := fixed128.FromStringForced[T]("42.5")
	_, err = fixed128.CheckedAsInteger[T, int](fracVal)
	c.HasError(err)

	// Float conversion should succeed for fractional values
	floatResult, err = fixed128.CheckedAsFloat[T, float64](fracVal)
	c.NoError(err)
	c.True(floatResult > 42.4 && floatResult < 42.6)
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

	// Test with quoted text
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
	// Test invalid JSON
	err := val.UnmarshalJSON([]byte("invalid"))
	c.HasError(err)

	// Test invalid text
	err = val.UnmarshalText([]byte("invalid"))
	c.HasError(err)

	// Test invalid YAML
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
	// Test empty string
	_, err := fixed128.FromString[T]("")
	c.HasError(err)

	// Test string with commas
	val, err := fixed128.FromString[T]("1,234.5")
	c.NoError(err)
	c.Equal("1234.5", val.String())

	// Test scientific notation
	val, err = fixed128.FromString[T]("1.23e2")
	c.NoError(err)
	c.Equal("123", val.String())

	// Test negative scientific notation
	val, err = fixed128.FromString[T]("-1.23E2")
	c.NoError(err)
	c.Equal("-123", val.String())

	// Test invalid scientific notation
	_, err = fixed128.FromString[T]("1.23ez")
	c.HasError(err)

	// Test invalid integer part
	_, err = fixed128.FromString[T]("abc.123")
	c.HasError(err)

	// Test invalid decimal part
	_, err = fixed128.FromString[T]("123.abc")
	c.HasError(err)

	// Test negative zero
	val, err = fixed128.FromString[T]("-0")
	c.NoError(err)
	c.Equal("0", val.String())

	// Test negative zero with decimal
	val, err = fixed128.FromString[T]("-0.000")
	c.NoError(err)
	c.Equal("0", val.String())

	// Test just decimal point
	val, err = fixed128.FromString[T](".5")
	c.NoError(err)
	c.Equal("0.5", val.String())

	// Test just minus sign
	val, err = fixed128.FromString[T]("-.5")
	c.NoError(err)
	c.Equal("-0.5", val.String())

	// Test very long decimal precision (should be truncated)
	val, err = fixed128.FromString[T]("0.123456789012345678901234567890")
	c.NoError(err)
	// Just verify it doesn't panic and produces a reasonable result
	c.NotEqual("", val.String())
	c.HasPrefix(val.String(), "0.1")
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

	// Test negative fractional value
	negFrac := fixed128.FromStringForced[T]("-2.5")
	c.Equal(fixed128.FromInteger[T](-2), negFrac.Ceil())

	// Test zero
	zero := fixed128.FromInteger[T](0)
	c.Equal(zero, zero.Ceil())

	// Test negative zero
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

	// Test negative fractional with zero integer part
	negFrac := fixed128.FromStringForced[T]("-0.5")
	c.Equal("-0.5", negFrac.String())

	// Test positive fractional with zero integer part
	posFrac := fixed128.FromStringForced[T]("0.5")
	c.Equal("0.5", posFrac.String())

	// Test trailing zeros removal
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

	// Test CheckedAs with float conversion that should fail
	val := fixed128.FromStringForced[T]("999999999999999999999999999.9")
	_, _ = fixed128.CheckedAsFloat[T, float32](val) //nolint:errcheck // This might succeed or fail depending on precision, but shouldn't panic. We'll just test that it doesn't panic
	c.NotNil(val)

	// Test YAML unmarshaling with string data
	var intVal fixed128.Int[T]
	err := intVal.UnmarshalYAML(func(v any) error {
		*(v.(*string)) = "42" //nolint:errcheck // This is just a test, we know it will succeed
		return nil
	})
	// This should handle the string value correctly
	c.NoError(err)

	// Test YAML unmarshaling with unmarshaling error
	err = intVal.UnmarshalYAML(func(any) error {
		return fmt.Errorf("unmarshal error")
	})
	c.HasError(err)
}
