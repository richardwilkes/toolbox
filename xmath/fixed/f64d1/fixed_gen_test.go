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

package f64d1_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed/f64d1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

type embedded struct {
	Field f64d1.Int
}

func TestConversion(t *testing.T) {
	assert.Equal(t, "0.1", f64d1.FromFloat64(0.1).String())
	assert.Equal(t, "0.2", f64d1.FromFloat64(0.2).String())
	assert.Equal(t, "0.3", f64d1.FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", f64d1.FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", f64d1.FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", f64d1.FromStringForced("-0.3").String())
	assert.Equal(t, "0.3", f64d1.FromStringForced("0.333").String())
	assert.Equal(t, "-0.3", f64d1.FromStringForced("-0.333").String())
	assert.Equal(t, "0.6", f64d1.FromStringForced("0.666").String())
	assert.Equal(t, "-0.6", f64d1.FromStringForced("-0.666").String())
	assert.Equal(t, "1", f64d1.FromFloat64(1.04).String())
	assert.Equal(t, "1", f64d1.FromFloat64(1.049).String())
	assert.Equal(t, "1", f64d1.FromFloat64(1.05).String())
	assert.Equal(t, "1", f64d1.FromFloat64(1.09).String())
	assert.Equal(t, "-1", f64d1.FromFloat64(-1.04).String())
	assert.Equal(t, "-1", f64d1.FromFloat64(-1.049).String())
	assert.Equal(t, "-1", f64d1.FromFloat64(-1.05).String())
	assert.Equal(t, "-1", f64d1.FromFloat64(-1.09).String())
	assert.Equal(t, "0.4", f64d1.FromStringForced("0.405").String())
	assert.Equal(t, "-0.4", f64d1.FromStringForced("-0.405").String())

	v, err := f64d1.FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, f64d1.FromInt(33))

	v, err = f64d1.FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, f64d1.FromInt(33))
}

func TestAddSub(t *testing.T) {
	oneThird := f64d1.FromStringForced("0.3")
	negTwoThirds := f64d1.FromStringForced("-0.6")
	one := f64d1.FromInt(1)
	oneAndTwoThirds := f64d1.FromStringForced("1.6")
	nineThousandSix := f64d1.FromInt(9006)
	two := f64d1.FromInt(2)
	assert.Equal(t, "0.9", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.7", (one - oneThird).String())
	assert.Equal(t, "-1.6", (negTwoThirds - one).String())
	assert.Equal(t, "0", (negTwoThirds - one + oneAndTwoThirds).String())
	assert.Equal(t, f64d1.FromInt(10240), f64d1.FromInt(1234)+nineThousandSix)
	assert.Equal(t, "10240", (f64d1.FromInt(1234) + nineThousandSix).String())
	assert.Equal(t, "-1.5", (f64d1.FromFloat64(0.5) - two).String())
}

func TestMulDiv(t *testing.T) {
	pointThree := f64d1.FromStringForced("0.3")
	negativePointThree := f64d1.FromStringForced("-0.3")
	assert.Equal(t, "0.3", f64d1.FromInt(1).Div(f64d1.FromInt(3)).String())
	assert.Equal(t, "-0.3", f64d1.FromInt(1).Div(f64d1.FromInt(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(f64d1.FromInt(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(f64d1.FromInt(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(f64d1.FromInt(3)).String())
}

func TestMod(t *testing.T) {
	assert.Equal(t, f64d1.FromInt(1), f64d1.FromInt(3).Mod(f64d1.FromInt(2)))
	assert.Equal(t, f64d1.FromStringForced("0.3"), f64d1.FromStringForced("9.3").Mod(f64d1.FromInt(3)))
	assert.Equal(t, f64d1.FromStringForced("0.1"), f64d1.FromStringForced("3.1").Mod(f64d1.FromStringForced("0.2")))
}

func TestTrunc(t *testing.T) {
	assert.Equal(t, f64d1.FromInt(0), f64d1.FromStringForced("0.3333").Trunc())
	assert.Equal(t, f64d1.FromInt(2), f64d1.FromStringForced("2.6789").Trunc())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromInt(3).Trunc())
	assert.Equal(t, f64d1.FromInt(0), f64d1.FromStringForced("-0.3333").Trunc())
	assert.Equal(t, f64d1.FromInt(-2), f64d1.FromStringForced("-2.6789").Trunc())
	assert.Equal(t, f64d1.FromInt(-3), f64d1.FromInt(-3).Trunc())
}

func TestCeil(t *testing.T) {
	assert.Equal(t, f64d1.FromInt(1), f64d1.FromStringForced("0.3333").Ceil())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromStringForced("2.6789").Ceil())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromInt(3).Ceil())
	assert.Equal(t, f64d1.FromInt(0), f64d1.FromStringForced("-0.3333").Ceil())
	assert.Equal(t, f64d1.FromInt(-2), f64d1.FromStringForced("-2.6789").Ceil())
	assert.Equal(t, f64d1.FromInt(-3), f64d1.FromInt(-3).Ceil())
}

func TestRound(t *testing.T) {
	assert.Equal(t, f64d1.FromInt(0), f64d1.FromStringForced("0.3333").Round())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromStringForced("2.6789").Round())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromInt(3).Round())
	assert.Equal(t, f64d1.FromInt(0), f64d1.FromStringForced("-0.3333").Round())
	assert.Equal(t, f64d1.FromInt(-3), f64d1.FromStringForced("-2.6789").Round())
	assert.Equal(t, f64d1.FromInt(-3), f64d1.FromInt(-3).Round())
}

func TestAbs(t *testing.T) {
	assert.Equal(t, f64d1.FromStringForced("0.3333"), f64d1.FromStringForced("0.3333").Abs())
	assert.Equal(t, f64d1.FromStringForced("2.6789"), f64d1.FromStringForced("2.6789").Abs())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromInt(3).Abs())
	assert.Equal(t, f64d1.FromStringForced("0.3333"), f64d1.FromStringForced("-0.3333").Abs())
	assert.Equal(t, f64d1.FromStringForced("2.6789"), f64d1.FromStringForced("-2.6789").Abs())
	assert.Equal(t, f64d1.FromInt(3), f64d1.FromInt(-3).Abs())
}

func TestJSON(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		testJSON(t, f64d1.FromInt(i))
	}
	testJSON(t, f64d1.FromInt64(1844674407371259000))
}

func testJSON(t *testing.T, v f64d1.Int) {
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
		testYAML(t, f64d1.FromInt(i))
	}
	testYAML(t, f64d1.FromInt64(1844674407371259000))
}

func testYAML(t *testing.T, v f64d1.Int) {
	t.Helper()
	e1 := embedded{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
