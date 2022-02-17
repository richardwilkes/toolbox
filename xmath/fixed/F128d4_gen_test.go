// Code created from "fixed_test.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

type embedded128d4 struct {
	Field fixed.F128d4
}

func TestConversion128d4(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F128d4FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F128d4FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F128d4FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F128d4FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F128d4FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F128d4FromStringForced("-0.3").String())
	assert.Equal(t, "0.3333", fixed.F128d4FromStringForced("0.333333").String())
	assert.Equal(t, "-0.3333", fixed.F128d4FromStringForced("-0.333333").String())
	assert.Equal(t, "0.6666", fixed.F128d4FromStringForced("0.666666").String())
	assert.Equal(t, "-0.6666", fixed.F128d4FromStringForced("-0.666666").String())
	assert.Equal(t, "1", fixed.F128d4FromFloat64(1.00004).String())
	assert.Equal(t, "1", fixed.F128d4FromFloat64(1.000049).String())
	assert.Equal(t, "1", fixed.F128d4FromFloat64(1.00005).String())
	assert.Equal(t, "1", fixed.F128d4FromFloat64(1.00009).String())
	assert.Equal(t, "-1", fixed.F128d4FromFloat64(-1.00004).String())
	assert.Equal(t, "-1", fixed.F128d4FromFloat64(-1.000049).String())
	assert.Equal(t, "-1", fixed.F128d4FromFloat64(-1.00005).String())
	assert.Equal(t, "-1", fixed.F128d4FromFloat64(-1.00009).String())
	assert.Equal(t, "0.0004", fixed.F128d4FromStringForced("0.000405").String())
	assert.Equal(t, "-0.0004", fixed.F128d4FromStringForced("-0.000405").String())

	v, err := fixed.F128d4FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d4FromInt(33))

	v, err = fixed.F128d4FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d4FromInt(33))
}

func TestAddSub128d4(t *testing.T) {
	oneThird := fixed.F128d4FromStringForced("0.3333")
	negTwoThirds := fixed.F128d4FromStringForced("-0.6666")
	one := fixed.F128d4FromInt(1)
	oneAndTwoThirds := fixed.F128d4FromStringForced("1.6666")
	nineThousandSix := fixed.F128d4FromInt(9006)
	ninetyPointZeroSix := fixed.F128d4FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F128d4FromStringForced("12.34")
	two := fixed.F128d4FromInt(2)
	assert.Equal(t, "0.9999", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.6667", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.6666", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, fixed.F128d4FromInt(10240), fixed.F128d4FromInt(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (fixed.F128d4FromInt(1234).Add(nineThousandSix)).String())
	assert.Equal(t, fixed.F128d4FromStringForced("102.4"), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	assert.Equal(t, "102.4", (twelvePointThirtyFour.Add(ninetyPointZeroSix)).String())
	assert.Equal(t, "-1.5", (fixed.F128d4FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv128d4(t *testing.T) {
	pointThree := fixed.F128d4FromStringForced("0.3")
	negativePointThree := fixed.F128d4FromStringForced("-0.3")
	assert.Equal(t, "0.3333", fixed.F128d4FromInt(1).Div(fixed.F128d4FromInt(3)).String())
	assert.Equal(t, "-0.3333", fixed.F128d4FromInt(1).Div(fixed.F128d4FromInt(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F128d4FromInt(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F128d4FromInt(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F128d4FromInt(3)).String())
}

func TestMod128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromInt(1), fixed.F128d4FromInt(3).Mod(fixed.F128d4FromInt(2)))
	assert.Equal(t, fixed.F128d4FromStringForced("0.3"), fixed.F128d4FromStringForced("9.3").Mod(fixed.F128d4FromInt(3)))
	assert.Equal(t, fixed.F128d4FromStringForced("0.1"), fixed.F128d4FromStringForced("3.1").Mod(fixed.F128d4FromStringForced("0.2")))
}

func TestTrunc128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromInt(0), fixed.F128d4FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F128d4FromInt(2), fixed.F128d4FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(3).Trunc())
	assert.Equal(t, fixed.F128d4FromInt(0), fixed.F128d4FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F128d4FromInt(-2), fixed.F128d4FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F128d4FromInt(-3), fixed.F128d4FromInt(-3).Trunc())
}

func TestCeil128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromInt(1), fixed.F128d4FromStringForced("0.3333").Ceil())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromStringForced("2.6789").Ceil())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(3).Ceil())
	assert.Equal(t, fixed.F128d4FromInt(0), fixed.F128d4FromStringForced("-0.3333").Ceil())
	assert.Equal(t, fixed.F128d4FromInt(-2), fixed.F128d4FromStringForced("-2.6789").Ceil())
	assert.Equal(t, fixed.F128d4FromInt(-3), fixed.F128d4FromInt(-3).Ceil())
}

func TestRound128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromInt(0), fixed.F128d4FromStringForced("0.3333").Round())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromStringForced("2.6789").Round())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(3).Round())
	assert.Equal(t, fixed.F128d4FromInt(0), fixed.F128d4FromStringForced("-0.3333").Round())
	assert.Equal(t, fixed.F128d4FromInt(-3), fixed.F128d4FromStringForced("-2.6789").Round())
	assert.Equal(t, fixed.F128d4FromInt(-3), fixed.F128d4FromInt(-3).Round())
}

func TestAbs128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromStringForced("0.3333"), fixed.F128d4FromStringForced("0.3333").Abs())
	assert.Equal(t, fixed.F128d4FromStringForced("2.6789"), fixed.F128d4FromStringForced("2.6789").Abs())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(3).Abs())
	assert.Equal(t, fixed.F128d4FromStringForced("0.3333"), fixed.F128d4FromStringForced("-0.3333").Abs())
	assert.Equal(t, fixed.F128d4FromStringForced("2.6789"), fixed.F128d4FromStringForced("-2.6789").Abs())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(-3).Abs())
}

func TestNeg128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromStringForced("-0.3333"), fixed.F128d4FromStringForced("0.3333").Neg())
	assert.Equal(t, fixed.F128d4FromStringForced("-2.6789"), fixed.F128d4FromStringForced("2.6789").Neg())
	assert.Equal(t, fixed.F128d4FromInt(-3), fixed.F128d4FromInt(3).Neg())
	assert.Equal(t, fixed.F128d4FromStringForced("0.3333"), fixed.F128d4FromStringForced("-0.3333").Neg())
	assert.Equal(t, fixed.F128d4FromStringForced("2.6789"), fixed.F128d4FromStringForced("-2.6789").Neg())
	assert.Equal(t, fixed.F128d4FromInt(3), fixed.F128d4FromInt(-3).Neg())
}

func TestCmp128d4(t *testing.T) {
	assert.Equal(t, 1, fixed.F128d4FromStringForced("0.3333").Cmp(fixed.F128d4FromInt(-3)))
	assert.Equal(t, -1, fixed.F128d4FromStringForced("2.6789").Cmp(fixed.F128d4FromInt(3)))
	assert.Equal(t, 0, fixed.F128d4FromInt(3).Cmp(fixed.F128d4FromInt(3)))
}

func TestEqual128d4(t *testing.T) {
	assert.Equal(t, false, fixed.F128d4FromStringForced("0.3333").Equal(fixed.F128d4FromInt(-3)))
	assert.Equal(t, false, fixed.F128d4FromStringForced("2.6789").Equal(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromInt(3).Equal(fixed.F128d4FromInt(3)))
}

func TestGreaterThan128d4(t *testing.T) {
	assert.Equal(t, true, fixed.F128d4FromStringForced("0.3333").GreaterThan(fixed.F128d4FromInt(-3)))
	assert.Equal(t, false, fixed.F128d4FromStringForced("2.6789").GreaterThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromInt(3).GreaterThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromInt(4).GreaterThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromStringForced("2.6789").GreaterThan(fixed.F128d4FromInt(-1)))
}

func TestGreaterThanOrEqual128d4(t *testing.T) {
	assert.Equal(t, true, fixed.F128d4FromStringForced("0.3333").GreaterThanOrEqual(fixed.F128d4FromInt(-3)))
	assert.Equal(t, false, fixed.F128d4FromStringForced("2.6789").GreaterThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromInt(3).GreaterThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromInt(4).GreaterThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromStringForced("2.6789").GreaterThanOrEqual(fixed.F128d4FromInt(-1)))
}

func TestLessThan128d4(t *testing.T) {
	assert.Equal(t, false, fixed.F128d4FromStringForced("0.3333").LessThan(fixed.F128d4FromInt(-3)))
	assert.Equal(t, true, fixed.F128d4FromStringForced("2.6789").LessThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromInt(3).LessThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromInt(4).LessThan(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromStringForced("2.6789").LessThan(fixed.F128d4FromInt(-1)))
}

func TestLessThanOrEqual128d4(t *testing.T) {
	assert.Equal(t, false, fixed.F128d4FromStringForced("0.3333").LessThanOrEqual(fixed.F128d4FromInt(-3)))
	assert.Equal(t, true, fixed.F128d4FromStringForced("2.6789").LessThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, true, fixed.F128d4FromInt(3).LessThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromInt(4).LessThanOrEqual(fixed.F128d4FromInt(3)))
	assert.Equal(t, false, fixed.F128d4FromStringForced("2.6789").LessThanOrEqual(fixed.F128d4FromInt(-1)))
}

func TestJSON128d4(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testJSON128d4(t, fixed.F128d4FromInt(i))
	}
	testJSON128d4(t, fixed.F128d4FromFloat64(18446744073712590000))
}

func testJSON128d4(t *testing.T, v fixed.F128d4) {
	t.Helper()
	e1 := embedded128d4{Field: v}
	data, err := json.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded128d4
	err = json.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}

func TestYAML128d4(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testYAML128d4(t, fixed.F128d4FromInt(i))
	}
	testYAML128d4(t, fixed.F128d4FromFloat64(18446744073712590000))
}

func testYAML128d4(t *testing.T, v fixed.F128d4) {
	t.Helper()
	e1 := embedded128d4{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded128d4
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
