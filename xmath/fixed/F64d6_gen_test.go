// Code created from "fixed_test.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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

type embedded64d6 struct {
	Field fixed.F64d6
}

func TestConversion64d6(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F64d6FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F64d6FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F64d6FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F64d6FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F64d6FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F64d6FromStringForced("-0.3").String())
	assert.Equal(t, "0.333333", fixed.F64d6FromStringForced("0.33333333").String())
	assert.Equal(t, "-0.333333", fixed.F64d6FromStringForced("-0.33333333").String())
	assert.Equal(t, "0.666666", fixed.F64d6FromStringForced("0.66666666").String())
	assert.Equal(t, "-0.666666", fixed.F64d6FromStringForced("-0.66666666").String())
	assert.Equal(t, "1", fixed.F64d6FromFloat64(1.0000004).String())
	assert.Equal(t, "1", fixed.F64d6FromFloat64(1.00000049).String())
	assert.Equal(t, "1", fixed.F64d6FromFloat64(1.0000005).String())
	assert.Equal(t, "1", fixed.F64d6FromFloat64(1.0000009).String())
	assert.Equal(t, "-1", fixed.F64d6FromFloat64(-1.0000004).String())
	assert.Equal(t, "-1", fixed.F64d6FromFloat64(-1.00000049).String())
	assert.Equal(t, "-1", fixed.F64d6FromFloat64(-1.0000005).String())
	assert.Equal(t, "-1", fixed.F64d6FromFloat64(-1.0000009).String())
	assert.Equal(t, "0.000004", fixed.F64d6FromStringForced("0.00000405").String())
	assert.Equal(t, "-0.000004", fixed.F64d6FromStringForced("-0.00000405").String())

	v, err := fixed.F64d6FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d6FromInt64(33))

	v, err = fixed.F64d6FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d6FromInt64(33))
}

func TestAddSub64d6(t *testing.T) {
	oneThird := fixed.F64d6FromStringForced("0.333333")
	negTwoThirds := fixed.F64d6FromStringForced("-0.666666")
	one := fixed.F64d6FromInt64(1)
	oneAndTwoThirds := fixed.F64d6FromStringForced("1.666666")
	nineThousandSix := fixed.F64d6FromInt64(9006)
	ninetyPointZeroSix := fixed.F64d6FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F64d6FromStringForced("12.34")
	two := fixed.F64d6FromInt64(2)
	assert.Equal(t, "0.999999", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.666667", (one - oneThird).String())
	assert.Equal(t, "-1.666666", (negTwoThirds - one).String())
	assert.Equal(t, "0", (negTwoThirds - one + oneAndTwoThirds).String())
	assert.Equal(t, fixed.F64d6FromInt64(10240), fixed.F64d6FromInt64(1234)+nineThousandSix)
	assert.Equal(t, "10240", (fixed.F64d6FromInt64(1234) + nineThousandSix).String())
	assert.Equal(t, fixed.F64d6FromStringForced("102.4"), twelvePointThirtyFour+ninetyPointZeroSix)
	assert.Equal(t, "102.4", (twelvePointThirtyFour + ninetyPointZeroSix).String())
	assert.Equal(t, "-1.5", (fixed.F64d6FromFloat64(0.5) - two).String())
}

func TestMulDiv64d6(t *testing.T) {
	pointThree := fixed.F64d6FromStringForced("0.3")
	negativePointThree := fixed.F64d6FromStringForced("-0.3")
	assert.Equal(t, "0.333333", fixed.F64d6FromInt64(1).Div(fixed.F64d6FromInt64(3)).String())
	assert.Equal(t, "-0.333333", fixed.F64d6FromInt64(1).Div(fixed.F64d6FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F64d6FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F64d6FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F64d6FromInt64(3)).String())
}

func TestTrunc64d6(t *testing.T) {
	assert.Equal(t, fixed.F64d6FromInt64(0), fixed.F64d6FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F64d6FromInt64(2), fixed.F64d6FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F64d6FromInt64(3), fixed.F64d6FromInt64(3).Trunc())
	assert.Equal(t, fixed.F64d6FromInt64(0), fixed.F64d6FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F64d6FromInt64(-2), fixed.F64d6FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F64d6FromInt64(-3), fixed.F64d6FromInt64(-3).Trunc())
}

func TestJSON64d6(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		testJSON64d6(t, fixed.F64d6FromInt64(i))
	}
	testJSON64d6(t, fixed.F64d6FromInt64(1844674407371259000))
}

func testJSON64d6(t *testing.T, v fixed.F64d6) {
	t.Helper()
	e1 := embedded64d6{Field: v}
	data, err := json.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded64d6
	err = json.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}

func TestYAML64d6(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		testYAML64d6(t, fixed.F64d6FromInt64(i))
	}
	testYAML64d6(t, fixed.F64d6FromInt64(1844674407371259000))
}

func testYAML64d6(t *testing.T, v fixed.F64d6) {
	t.Helper()
	e1 := embedded64d6{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded64d6
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
