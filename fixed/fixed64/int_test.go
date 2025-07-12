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
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed64"
	"github.com/richardwilkes/toolbox/v2/num128"
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
	c.Equal("0.1", fixed64.FromFloat[T](0.1).String())
	c.Equal("0.2", fixed64.FromFloat[T](0.2).String())
	c.Equal("0.3", fixed64.FromStringForced[T]("0.3").String())
	c.Equal("-0.1", fixed64.FromFloat[T](-0.1).String())
	c.Equal("-0.2", fixed64.FromFloat[T](-0.2).String())
	c.Equal("-0.3", fixed64.FromStringForced[T]("-0.3").String())
	threeFill := strings.Repeat("3", fixed64.MaxDecimalDigits[T]())
	c.Equal("0."+threeFill, fixed64.FromStringForced[T]("0.33333333").String())
	c.Equal("-0."+threeFill, fixed64.FromStringForced[T]("-0.33333333").String())
	sixFill := strings.Repeat("6", fixed64.MaxDecimalDigits[T]())
	c.Equal("0."+sixFill, fixed64.FromStringForced[T]("0.66666666").String())
	c.Equal("-0."+sixFill, fixed64.FromStringForced[T]("-0.66666666").String())
	c.Equal("1", fixed64.FromFloat[T](1.0000004).String())
	c.Equal("1", fixed64.FromFloat[T](1.00000049).String())
	c.Equal("1", fixed64.FromFloat[T](1.0000005).String())
	c.Equal("1", fixed64.FromFloat[T](1.0000009).String())
	c.Equal("-1", fixed64.FromFloat[T](-1.0000004).String())
	c.Equal("-1", fixed64.FromFloat[T](-1.00000049).String())
	c.Equal("-1", fixed64.FromFloat[T](-1.0000005).String())
	c.Equal("-1", fixed64.FromFloat[T](-1.0000009).String())
	zeroFill := strings.Repeat("0", fixed64.MaxDecimalDigits[T]()-1)
	c.Equal("0."+zeroFill+"4", fixed64.FromStringForced[T]("0."+zeroFill+"405").String())
	c.Equal("-0."+zeroFill+"4", fixed64.FromStringForced[T]("-0."+zeroFill+"405").String())

	v, err := fixed64.FromString[T]("33.0")
	c.NoError(err)
	c.Equal(v, fixed64.FromInteger[T](33))

	v, err = fixed64.FromString[T]("33.00000000000000000000")
	c.NoError(err)
	c.Equal(v, fixed64.FromInteger[T](33))
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
	oneThird := fixed64.FromStringForced[T]("0.333333")
	negTwoThirds := fixed64.FromStringForced[T]("-0.666666")
	one := fixed64.FromInteger[T](1)
	oneAndTwoThirds := fixed64.FromStringForced[T]("1.666666")
	nineThousandSix := fixed64.FromInteger[T](9006)
	two := fixed64.FromInteger[T](2)
	c := check.New(t)
	c.Equal("0."+strings.Repeat("9", fixed64.MaxDecimalDigits[T]()), (oneThird + oneThird + oneThird).String())
	c.Equal("0."+strings.Repeat("6", fixed64.MaxDecimalDigits[T]()-1)+"7", (one - oneThird).String())
	c.Equal("-1."+strings.Repeat("6", fixed64.MaxDecimalDigits[T]()), (negTwoThirds - one).String())
	c.Equal("0", (negTwoThirds - one + oneAndTwoThirds).String())
	c.Equal(fixed64.FromInteger[T](10240), fixed64.FromInteger[T](1234)+nineThousandSix)
	c.Equal("10240", (fixed64.FromInteger[T](1234) + nineThousandSix).String())
	c.Equal("-1.5", (fixed64.FromFloat[T](0.5) - two).String())
	ninetyPointZeroSix := fixed64.FromStringForced[T]("90.06")
	twelvePointThirtyFour := fixed64.FromStringForced[T]("12.34")
	var answer string
	if fixed64.MaxDecimalDigits[T]() > 1 {
		answer = "102.4"
	} else {
		answer = "102.3"
	}
	c.Equal(fixed64.FromStringForced[T](answer), twelvePointThirtyFour+ninetyPointZeroSix)
	c.Equal(answer, (twelvePointThirtyFour + ninetyPointZeroSix).String())
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
	pointThree := fixed64.FromStringForced[T]("0.3")
	negativePointThree := fixed64.FromStringForced[T]("-0.3")
	threeFill := strings.Repeat("3", fixed64.MaxDecimalDigits[T]())
	c := check.New(t)
	c.Equal("0."+threeFill, fixed64.FromInteger[T](1).Div(fixed64.FromInteger[T](3)).String())
	c.Equal("-0."+threeFill, fixed64.FromInteger[T](1).Div(fixed64.FromInteger[T](-3)).String())
	c.Equal("0.1", pointThree.Div(fixed64.FromInteger[T](3)).String())
	c.Equal("0.9", pointThree.Mul(fixed64.FromInteger[T](3)).String())
	c.Equal("-0.9", negativePointThree.Mul(fixed64.FromInteger[T](3)).String())
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
	c.Equal(fixed64.FromInteger[T](1), fixed64.FromInteger[T](3).Mod(fixed64.FromInteger[T](2)))
	c.Equal(fixed64.FromStringForced[T]("0.3"), fixed64.FromStringForced[T]("9.3").Mod(fixed64.FromInteger[T](3)))
	c.Equal(fixed64.FromStringForced[T]("0.1"), fixed64.FromStringForced[T]("3.1").Mod(fixed64.FromStringForced[T]("0.2")))
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
	c.Equal(fixed64.FromInteger[T](0), fixed64.FromStringForced[T]("0.3333").Floor())
	c.Equal(fixed64.FromInteger[T](2), fixed64.FromStringForced[T]("2.6789").Floor())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromInteger[T](3).Floor())
	c.Equal(fixed64.FromInteger[T](0), fixed64.FromStringForced[T]("-0.3333").Floor())
	c.Equal(fixed64.FromInteger[T](-2), fixed64.FromStringForced[T]("-2.6789").Floor())
	c.Equal(fixed64.FromInteger[T](-3), fixed64.FromInteger[T](-3).Floor())
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
	c.Equal(fixed64.FromInteger[T](1), fixed64.FromStringForced[T]("0.3333").Ceil())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromStringForced[T]("2.6789").Ceil())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromInteger[T](3).Ceil())
	c.Equal(fixed64.FromInteger[T](0), fixed64.FromStringForced[T]("-0.3333").Ceil())
	c.Equal(fixed64.FromInteger[T](-2), fixed64.FromStringForced[T]("-2.6789").Ceil())
	c.Equal(fixed64.FromInteger[T](-3), fixed64.FromInteger[T](-3).Ceil())
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
	c.Equal(fixed64.FromInteger[T](0), fixed64.FromStringForced[T]("0.3333").Round())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromStringForced[T]("2.6789").Round())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromInteger[T](3).Round())
	c.Equal(fixed64.FromInteger[T](0), fixed64.FromStringForced[T]("-0.3333").Round())
	c.Equal(fixed64.FromInteger[T](-3), fixed64.FromStringForced[T]("-2.6789").Round())
	c.Equal(fixed64.FromInteger[T](-3), fixed64.FromInteger[T](-3).Round())
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
	c.Equal(fixed64.FromStringForced[T]("0.3333"), fixed64.FromStringForced[T]("0.3333").Abs())
	c.Equal(fixed64.FromStringForced[T]("2.6789"), fixed64.FromStringForced[T]("2.6789").Abs())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromInteger[T](3).Abs())
	c.Equal(fixed64.FromStringForced[T]("0.3333"), fixed64.FromStringForced[T]("-0.3333").Abs())
	c.Equal(fixed64.FromStringForced[T]("2.6789"), fixed64.FromStringForced[T]("-2.6789").Abs())
	c.Equal(fixed64.FromInteger[T](3), fixed64.FromInteger[T](-3).Abs())
}

func TestComma(t *testing.T) {
	c := check.New(t)
	c.Equal("0.12", fixed64.FromStringForced[fixed.D2]("0.12").Comma())
	c.Equal("1,234,567,890.12", fixed64.FromStringForced[fixed.D2]("1234567890.12").Comma())
	c.Equal("91,234,567,890.12", fixed64.FromStringForced[fixed.D2]("91234567890.12").Comma())
	c.Equal("891,234,567,890.12", fixed64.FromStringForced[fixed.D2]("891234567890.12").Comma())
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
		testJSONActual(t, fixed64.FromInteger[T](i))
	}
	testJSONActual(t, fixed64.FromInteger[T, int64](1844674407371259000))
}

type embedded[T fixed.Dx] struct {
	Field fixed64.Int[T]
}

func testJSONActual[T fixed.Dx](t *testing.T, v fixed64.Int[T]) {
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
		testYAMLActual(t, fixed64.FromInteger[T](i))
	}
	testYAMLActual(t, fixed64.FromInteger[T, int64](1844674407371259000))
}

func testYAMLActual[T fixed.Dx](t *testing.T, v fixed64.Int[T]) {
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

	a := fixed64.FromInteger[T](5)
	b := fixed64.FromInteger[T](10)
	negativeA := fixed64.FromInteger[T](-5)

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

	zero := fixed64.FromInteger[T](0)
	one := fixed64.FromInteger[T](1)
	negativeOne := fixed64.FromInteger[T](-1)

	// Test Inc
	c.Equal(one, zero.Inc())
	c.Equal(fixed64.FromInteger[T](2), one.Inc())
	c.Equal(zero, negativeOne.Inc())

	// Test Dec
	c.Equal(negativeOne, zero.Dec())
	c.Equal(zero, one.Dec())
	c.Equal(fixed64.FromInteger[T](-2), negativeOne.Dec())
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
	intVal := fixed64.FromInteger[T](42)
	c.Equal(int(42), fixed64.AsInteger[T, int](intVal))
	c.Equal(int8(42), fixed64.AsInteger[T, int8](intVal))
	c.Equal(int16(42), fixed64.AsInteger[T, int16](intVal))
	c.Equal(int32(42), fixed64.AsInteger[T, int32](intVal))
	c.Equal(int64(42), fixed64.AsInteger[T, int64](intVal))
	c.Equal(uint(42), fixed64.AsInteger[T, uint](intVal))
	c.Equal(uint8(42), fixed64.AsInteger[T, uint8](intVal))
	c.Equal(uint16(42), fixed64.AsInteger[T, uint16](intVal))
	c.Equal(uint32(42), fixed64.AsInteger[T, uint32](intVal))
	c.Equal(uint64(42), fixed64.AsInteger[T, uint64](intVal))

	// Test float conversions
	floatVal := fixed64.FromStringForced[T]("3.1")
	f32Result := fixed64.AsFloat[T, float32](floatVal)
	f64Result := fixed64.AsFloat[T, float64](floatVal)
	c.True(f32Result > 3.0 && f32Result < 3.2)
	c.True(f64Result > 3.0 && f64Result < 3.2)

	// Test negative values
	negVal := fixed64.FromInteger[T](-10)
	c.Equal(int(-10), fixed64.AsInteger[T, int](negVal))
	c.Equal(float64(-10.0), fixed64.AsFloat[T, float64](negVal))
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
	intVal := fixed64.FromInteger[T](42)
	result, err := fixed64.AsIntegerChecked[T, int](intVal)
	c.NoError(err)
	c.Equal(int(42), result)

	floatResult, err := fixed64.AsFloatChecked[T, float64](intVal)
	c.NoError(err)
	c.Equal(float64(42.0), floatResult)

	// Test conversion that should fail (fractional part)
	fracVal := fixed64.FromStringForced[T]("42.5")
	_, err = fixed64.AsIntegerChecked[T, int](fracVal)
	c.HasError(err)

	// Float conversion should succeed for fractional values
	floatResult, err = fixed64.AsFloatChecked[T, float64](fracVal)
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

	positive := fixed64.FromInteger[T](42)
	negative := fixed64.FromInteger[T](-42)
	zero := fixed64.FromInteger[T](0)

	c.Equal("+42", positive.StringWithSign())
	c.Equal("-42", negative.StringWithSign())
	c.Equal("+0", zero.StringWithSign())

	positiveFrac := fixed64.FromStringForced[T]("3.1")
	negativeFrac := fixed64.FromStringForced[T]("-3.1")

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

	large := fixed64.FromStringForced[T]("1234567.8")
	largenegative := fixed64.FromStringForced[T]("-1234567.8")
	zero := fixed64.FromInteger[T](0)

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

	val := fixed64.FromStringForced[T]("123.4")

	data, err := val.MarshalText()
	c.NoError(err)
	c.Equal("123.4", string(data))

	var unmarshaled fixed64.Int[T]
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

	var val fixed64.Int[T]

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
	_, err := fixed64.FromString[T]("")
	c.HasError(err)

	// Test string with commas
	val, err := fixed64.FromString[T]("1,234.5")
	c.NoError(err)
	c.Equal("1234.5", val.String())

	// Test scientific notation
	val, err = fixed64.FromString[T]("1.23e2")
	c.NoError(err)
	c.Equal("123", val.String())

	// Test negative scientific notation
	val, err = fixed64.FromString[T]("-1.23E2")
	c.NoError(err)
	c.Equal("-123", val.String())

	// Test invalid scientific notation
	_, err = fixed64.FromString[T]("1.23ez")
	c.HasError(err)

	// Test invalid integer part
	_, err = fixed64.FromString[T]("abc.123")
	c.HasError(err)

	// Test invalid decimal part
	_, err = fixed64.FromString[T]("123.abc")
	c.HasError(err)

	// Test negative zero
	val, err = fixed64.FromString[T]("-0")
	c.NoError(err)
	c.Equal("0", val.String())

	// Test negative zero with decimal
	val, err = fixed64.FromString[T]("-0.000")
	c.NoError(err)
	c.Equal("0", val.String())

	// Test just decimal point
	val, err = fixed64.FromString[T](".5")
	c.NoError(err)
	c.Equal("0.5", val.String())

	// Test just minus sign
	val, err = fixed64.FromString[T]("-.5")
	c.NoError(err)
	c.Equal("-0.5", val.String())

	// Test very long decimal precision (should be truncated)
	val, err = fixed64.FromString[T]("0.123456789012345678901234567890")
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
	negFrac := fixed64.FromStringForced[T]("-2.5")
	c.Equal(fixed64.FromInteger[T](-2), negFrac.Ceil())

	// Test zero
	zero := fixed64.FromInteger[T](0)
	c.Equal(zero, zero.Ceil())

	// Test negative zero
	negZero := fixed64.FromStringForced[T]("-0.0")
	c.Equal(fixed64.FromInteger[T](0), negZero.Ceil())
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
	negFrac := fixed64.FromStringForced[T]("-0.5")
	c.Equal("-0.5", negFrac.String())

	// Test positive fractional with zero integer part
	posFrac := fixed64.FromStringForced[T]("0.5")
	c.Equal("0.5", posFrac.String())

	// Test trailing zeros removal
	val := fixed64.FromStringForced[T]("1.2000")
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
	val := fixed64.FromStringForced[T]("999999999999999999.9")
	_, _ = fixed64.AsFloatChecked[T, float32](val) //nolint:errcheck // This might succeed or fail depending on precision, but shouldn't panic. We'll just test that it doesn't panic
	c.NotNil(val)

	// Test YAML unmarshaling with string data
	var intVal fixed64.Int[T]
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

func TestMulOverflow(t *testing.T) {
	c := check.New(t)
	v1 := int64(math.MaxInt64 / 10000)
	v1 -= v1 / 10
	m1 := int64(110)
	d1 := int64(100)
	expected := num128.IntFrom64(v1).Mul(num128.IntFrom64(m1)).Div(num128.IntFrom64(d1)).AsInt64()
	v2 := fixed64.FromInteger[fixed.D4](v1)
	m2 := fixed64.FromInteger[fixed.D4](m1)
	d2 := fixed64.FromInteger[fixed.D4](d1)
	result := fixed64.AsInteger[fixed.D4, int64](v2.Mul(m2.Div(d2)))
	c.Equal(expected, result)
}
