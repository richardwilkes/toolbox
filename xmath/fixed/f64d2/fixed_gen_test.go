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

package f64d2_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed/f64d2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

type embedded struct {
	Field f64d2.Int
}

func TestConversion(t *testing.T) {
	assert.Equal(t, "0.1", f64d2.FromFloat64(0.1).String())
	assert.Equal(t, "0.2", f64d2.FromFloat64(0.2).String())
	assert.Equal(t, "0.3", f64d2.FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", f64d2.FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", f64d2.FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", f64d2.FromStringForced("-0.3").String())
	assert.Equal(t, "0.33", f64d2.FromStringForced("0.3333").String())
	assert.Equal(t, "-0.33", f64d2.FromStringForced("-0.3333").String())
	assert.Equal(t, "0.66", f64d2.FromStringForced("0.6666").String())
	assert.Equal(t, "-0.66", f64d2.FromStringForced("-0.6666").String())
	assert.Equal(t, "1", f64d2.FromFloat64(1.004).String())
	assert.Equal(t, "1", f64d2.FromFloat64(1.0049).String())
	assert.Equal(t, "1", f64d2.FromFloat64(1.005).String())
	assert.Equal(t, "1", f64d2.FromFloat64(1.009).String())
	assert.Equal(t, "-1", f64d2.FromFloat64(-1.004).String())
	assert.Equal(t, "-1", f64d2.FromFloat64(-1.0049).String())
	assert.Equal(t, "-1", f64d2.FromFloat64(-1.005).String())
	assert.Equal(t, "-1", f64d2.FromFloat64(-1.009).String())
	assert.Equal(t, "0.04", f64d2.FromStringForced("0.0405").String())
	assert.Equal(t, "-0.04", f64d2.FromStringForced("-0.0405").String())

	v, err := f64d2.FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, f64d2.FromInt(33))

	v, err = f64d2.FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, f64d2.FromInt(33))
}

func TestAddSub(t *testing.T) {
	oneThird := f64d2.FromStringForced("0.33")
	negTwoThirds := f64d2.FromStringForced("-0.66")
	one := f64d2.FromInt(1)
	oneAndTwoThirds := f64d2.FromStringForced("1.66")
	nineThousandSix := f64d2.FromInt(9006)
	two := f64d2.FromInt(2)
	assert.Equal(t, "0.99", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.67", (one - oneThird).String())
	assert.Equal(t, "-1.66", (negTwoThirds - one).String())
	assert.Equal(t, "0", (negTwoThirds - one + oneAndTwoThirds).String())
	assert.Equal(t, f64d2.FromInt(10240), f64d2.FromInt(1234)+nineThousandSix)
	assert.Equal(t, "10240", (f64d2.FromInt(1234) + nineThousandSix).String())
	assert.Equal(t, "-1.5", (f64d2.FromFloat64(0.5) - two).String())
	ninetyPointZeroSix := f64d2.FromStringForced("90.06")
	twelvePointThirtyFour := f64d2.FromStringForced("12.34")
	assert.Equal(t, f64d2.FromStringForced("102.4"), twelvePointThirtyFour+ninetyPointZeroSix)
	assert.Equal(t, "102.4", (twelvePointThirtyFour + ninetyPointZeroSix).String())
}

func TestMulDiv(t *testing.T) {
	pointThree := f64d2.FromStringForced("0.3")
	negativePointThree := f64d2.FromStringForced("-0.3")
	assert.Equal(t, "0.33", f64d2.FromInt(1).Div(f64d2.FromInt(3)).String())
	assert.Equal(t, "-0.33", f64d2.FromInt(1).Div(f64d2.FromInt(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(f64d2.FromInt(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(f64d2.FromInt(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(f64d2.FromInt(3)).String())
}

func TestMod(t *testing.T) {
	assert.Equal(t, f64d2.FromInt(1), f64d2.FromInt(3).Mod(f64d2.FromInt(2)))
	assert.Equal(t, f64d2.FromStringForced("0.3"), f64d2.FromStringForced("9.3").Mod(f64d2.FromInt(3)))
	assert.Equal(t, f64d2.FromStringForced("0.1"), f64d2.FromStringForced("3.1").Mod(f64d2.FromStringForced("0.2")))
}

func TestTrunc(t *testing.T) {
	assert.Equal(t, f64d2.FromInt(0), f64d2.FromStringForced("0.3333").Trunc())
	assert.Equal(t, f64d2.FromInt(2), f64d2.FromStringForced("2.6789").Trunc())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromInt(3).Trunc())
	assert.Equal(t, f64d2.FromInt(0), f64d2.FromStringForced("-0.3333").Trunc())
	assert.Equal(t, f64d2.FromInt(-2), f64d2.FromStringForced("-2.6789").Trunc())
	assert.Equal(t, f64d2.FromInt(-3), f64d2.FromInt(-3).Trunc())
}

func TestCeil(t *testing.T) {
	assert.Equal(t, f64d2.FromInt(1), f64d2.FromStringForced("0.3333").Ceil())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromStringForced("2.6789").Ceil())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromInt(3).Ceil())
	assert.Equal(t, f64d2.FromInt(0), f64d2.FromStringForced("-0.3333").Ceil())
	assert.Equal(t, f64d2.FromInt(-2), f64d2.FromStringForced("-2.6789").Ceil())
	assert.Equal(t, f64d2.FromInt(-3), f64d2.FromInt(-3).Ceil())
}

func TestRound(t *testing.T) {
	assert.Equal(t, f64d2.FromInt(0), f64d2.FromStringForced("0.3333").Round())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromStringForced("2.6789").Round())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromInt(3).Round())
	assert.Equal(t, f64d2.FromInt(0), f64d2.FromStringForced("-0.3333").Round())
	assert.Equal(t, f64d2.FromInt(-3), f64d2.FromStringForced("-2.6789").Round())
	assert.Equal(t, f64d2.FromInt(-3), f64d2.FromInt(-3).Round())
}

func TestAbs(t *testing.T) {
	assert.Equal(t, f64d2.FromStringForced("0.3333"), f64d2.FromStringForced("0.3333").Abs())
	assert.Equal(t, f64d2.FromStringForced("2.6789"), f64d2.FromStringForced("2.6789").Abs())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromInt(3).Abs())
	assert.Equal(t, f64d2.FromStringForced("0.3333"), f64d2.FromStringForced("-0.3333").Abs())
	assert.Equal(t, f64d2.FromStringForced("2.6789"), f64d2.FromStringForced("-2.6789").Abs())
	assert.Equal(t, f64d2.FromInt(3), f64d2.FromInt(-3).Abs())
}

func TestJSON(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testJSON(t, f64d2.FromInt(i))
	}
	testJSON(t, f64d2.FromInt64(1844674407371259000))
}

func testJSON(t *testing.T, v f64d2.Int) {
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
		testYAML(t, f64d2.FromInt(i))
	}
	testYAML(t, f64d2.FromInt64(1844674407371259000))
}

func testYAML(t *testing.T, v f64d2.Int) {
	t.Helper()
	e1 := embedded{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
