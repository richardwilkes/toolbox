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

package f128d1_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed/f128d1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

type embedded struct {
	Field f128d1.Int
}

func TestConversion(t *testing.T) {
	assert.Equal(t, "0.1", f128d1.FromFloat64(0.1).String())
	assert.Equal(t, "0.2", f128d1.FromFloat64(0.2).String())
	assert.Equal(t, "0.3", f128d1.FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", f128d1.FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", f128d1.FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", f128d1.FromStringForced("-0.3").String())
	assert.Equal(t, "0.3", f128d1.FromStringForced("0.333").String())
	assert.Equal(t, "-0.3", f128d1.FromStringForced("-0.333").String())
	assert.Equal(t, "0.6", f128d1.FromStringForced("0.666").String())
	assert.Equal(t, "-0.6", f128d1.FromStringForced("-0.666").String())
	assert.Equal(t, "1", f128d1.FromFloat64(1.04).String())
	assert.Equal(t, "1", f128d1.FromFloat64(1.049).String())
	assert.Equal(t, "1", f128d1.FromFloat64(1.05).String())
	assert.Equal(t, "1", f128d1.FromFloat64(1.09).String())
	assert.Equal(t, "-1", f128d1.FromFloat64(-1.04).String())
	assert.Equal(t, "-1", f128d1.FromFloat64(-1.049).String())
	assert.Equal(t, "-1", f128d1.FromFloat64(-1.05).String())
	assert.Equal(t, "-1", f128d1.FromFloat64(-1.09).String())
	assert.Equal(t, "0.4", f128d1.FromStringForced("0.405").String())
	assert.Equal(t, "-0.4", f128d1.FromStringForced("-0.405").String())

	v, err := f128d1.FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, f128d1.FromInt(33))

	v, err = f128d1.FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, f128d1.FromInt(33))
}

func TestAddSub(t *testing.T) {
	oneThird := f128d1.FromStringForced("0.3")
	negTwoThirds := f128d1.FromStringForced("-0.6")
	one := f128d1.FromInt(1)
	oneAndTwoThirds := f128d1.FromStringForced("1.6")
	nineThousandSix := f128d1.FromInt(9006)
	two := f128d1.FromInt(2)
	assert.Equal(t, "0.9", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.7", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.6", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, f128d1.FromInt(10240), f128d1.FromInt(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (f128d1.FromInt(1234).Add(nineThousandSix)).String())
	assert.Equal(t, "-1.5", (f128d1.FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv(t *testing.T) {
	pointThree := f128d1.FromStringForced("0.3")
	negativePointThree := f128d1.FromStringForced("-0.3")
	assert.Equal(t, "0.3", f128d1.FromInt(1).Div(f128d1.FromInt(3)).String())
	assert.Equal(t, "-0.3", f128d1.FromInt(1).Div(f128d1.FromInt(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(f128d1.FromInt(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(f128d1.FromInt(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(f128d1.FromInt(3)).String())
}

func TestMod(t *testing.T) {
	assert.Equal(t, f128d1.FromInt(1), f128d1.FromInt(3).Mod(f128d1.FromInt(2)))
	assert.Equal(t, f128d1.FromStringForced("0.3"), f128d1.FromStringForced("9.3").Mod(f128d1.FromInt(3)))
	assert.Equal(t, f128d1.FromStringForced("0.1"), f128d1.FromStringForced("3.1").Mod(f128d1.FromStringForced("0.2")))
}

func TestTrunc(t *testing.T) {
	assert.Equal(t, f128d1.FromInt(0), f128d1.FromStringForced("0.3333").Trunc())
	assert.Equal(t, f128d1.FromInt(2), f128d1.FromStringForced("2.6789").Trunc())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(3).Trunc())
	assert.Equal(t, f128d1.FromInt(0), f128d1.FromStringForced("-0.3333").Trunc())
	assert.Equal(t, f128d1.FromInt(-2), f128d1.FromStringForced("-2.6789").Trunc())
	assert.Equal(t, f128d1.FromInt(-3), f128d1.FromInt(-3).Trunc())
}

func TestCeil(t *testing.T) {
	assert.Equal(t, f128d1.FromInt(1), f128d1.FromStringForced("0.3333").Ceil())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromStringForced("2.6789").Ceil())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(3).Ceil())
	assert.Equal(t, f128d1.FromInt(0), f128d1.FromStringForced("-0.3333").Ceil())
	assert.Equal(t, f128d1.FromInt(-2), f128d1.FromStringForced("-2.6789").Ceil())
	assert.Equal(t, f128d1.FromInt(-3), f128d1.FromInt(-3).Ceil())
}

func TestRound(t *testing.T) {
	assert.Equal(t, f128d1.FromInt(0), f128d1.FromStringForced("0.3333").Round())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromStringForced("2.6789").Round())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(3).Round())
	assert.Equal(t, f128d1.FromInt(0), f128d1.FromStringForced("-0.3333").Round())
	assert.Equal(t, f128d1.FromInt(-3), f128d1.FromStringForced("-2.6789").Round())
	assert.Equal(t, f128d1.FromInt(-3), f128d1.FromInt(-3).Round())
}

func TestAbs(t *testing.T) {
	assert.Equal(t, f128d1.FromStringForced("0.3333"), f128d1.FromStringForced("0.3333").Abs())
	assert.Equal(t, f128d1.FromStringForced("2.6789"), f128d1.FromStringForced("2.6789").Abs())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(3).Abs())
	assert.Equal(t, f128d1.FromStringForced("0.3333"), f128d1.FromStringForced("-0.3333").Abs())
	assert.Equal(t, f128d1.FromStringForced("2.6789"), f128d1.FromStringForced("-2.6789").Abs())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(-3).Abs())
}

func TestNeg(t *testing.T) {
	assert.Equal(t, f128d1.FromStringForced("-0.3333"), f128d1.FromStringForced("0.3333").Neg())
	assert.Equal(t, f128d1.FromStringForced("-2.6789"), f128d1.FromStringForced("2.6789").Neg())
	assert.Equal(t, f128d1.FromInt(-3), f128d1.FromInt(3).Neg())
	assert.Equal(t, f128d1.FromStringForced("0.3333"), f128d1.FromStringForced("-0.3333").Neg())
	assert.Equal(t, f128d1.FromStringForced("2.6789"), f128d1.FromStringForced("-2.6789").Neg())
	assert.Equal(t, f128d1.FromInt(3), f128d1.FromInt(-3).Neg())
}

func TestCmp(t *testing.T) {
	assert.Equal(t, 1, f128d1.FromStringForced("0.3333").Cmp(f128d1.FromInt(-3)))
	assert.Equal(t, -1, f128d1.FromStringForced("2.6789").Cmp(f128d1.FromInt(3)))
	assert.Equal(t, 0, f128d1.FromInt(3).Cmp(f128d1.FromInt(3)))
}

func TestEqual(t *testing.T) {
	assert.Equal(t, false, f128d1.FromStringForced("0.3333").Equal(f128d1.FromInt(-3)))
	assert.Equal(t, false, f128d1.FromStringForced("2.6789").Equal(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromInt(3).Equal(f128d1.FromInt(3)))
}

func TestGreaterThan(t *testing.T) {
	assert.Equal(t, true, f128d1.FromStringForced("0.3333").GreaterThan(f128d1.FromInt(-3)))
	assert.Equal(t, false, f128d1.FromStringForced("2.6789").GreaterThan(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromInt(3).GreaterThan(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromInt(4).GreaterThan(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromStringForced("2.6789").GreaterThan(f128d1.FromInt(-1)))
}

func TestGreaterThanOrEqual(t *testing.T) {
	assert.Equal(t, true, f128d1.FromStringForced("0.3333").GreaterThanOrEqual(f128d1.FromInt(-3)))
	assert.Equal(t, false, f128d1.FromStringForced("2.6789").GreaterThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromInt(3).GreaterThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromInt(4).GreaterThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromStringForced("2.6789").GreaterThanOrEqual(f128d1.FromInt(-1)))
}

func TestLessThan(t *testing.T) {
	assert.Equal(t, false, f128d1.FromStringForced("0.3333").LessThan(f128d1.FromInt(-3)))
	assert.Equal(t, true, f128d1.FromStringForced("2.6789").LessThan(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromInt(3).LessThan(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromInt(4).LessThan(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromStringForced("2.6789").LessThan(f128d1.FromInt(-1)))
}

func TestLessThanOrEqual(t *testing.T) {
	assert.Equal(t, false, f128d1.FromStringForced("0.3333").LessThanOrEqual(f128d1.FromInt(-3)))
	assert.Equal(t, true, f128d1.FromStringForced("2.6789").LessThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, true, f128d1.FromInt(3).LessThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromInt(4).LessThanOrEqual(f128d1.FromInt(3)))
	assert.Equal(t, false, f128d1.FromStringForced("2.6789").LessThanOrEqual(f128d1.FromInt(-1)))
}

func TestJSON(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testJSON(t, f128d1.FromInt(i))
	}
	testJSON(t, f128d1.FromFloat64(18446744073712590000))
}

func testJSON(t *testing.T, v f128d1.Int) {
	t.Helper()
	e1 := embedded{Field: v}
	data, err := json.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded
	err = json.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}

func TestYAML(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testYAML(t, f128d1.FromInt(i))
	}
	testYAML(t, f128d1.FromFloat64(18446744073712590000))
}

func testYAML(t *testing.T, v f128d1.Int) {
	t.Helper()
	e1 := embedded{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
