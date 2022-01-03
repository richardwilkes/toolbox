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

type embedded128d2 struct {
	Field fixed.F128d2
}

func TestConversion128d2(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F128d2FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F128d2FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F128d2FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F128d2FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F128d2FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F128d2FromStringForced("-0.3").String())
	assert.Equal(t, "0.33", fixed.F128d2FromStringForced("0.3333").String())
	assert.Equal(t, "-0.33", fixed.F128d2FromStringForced("-0.3333").String())
	assert.Equal(t, "0.66", fixed.F128d2FromStringForced("0.6666").String())
	assert.Equal(t, "-0.66", fixed.F128d2FromStringForced("-0.6666").String())
	assert.Equal(t, "1", fixed.F128d2FromFloat64(1.004).String())
	assert.Equal(t, "1", fixed.F128d2FromFloat64(1.0049).String())
	assert.Equal(t, "1", fixed.F128d2FromFloat64(1.005).String())
	assert.Equal(t, "1", fixed.F128d2FromFloat64(1.009).String())
	assert.Equal(t, "-1", fixed.F128d2FromFloat64(-1.004).String())
	assert.Equal(t, "-1", fixed.F128d2FromFloat64(-1.0049).String())
	assert.Equal(t, "-1", fixed.F128d2FromFloat64(-1.005).String())
	assert.Equal(t, "-1", fixed.F128d2FromFloat64(-1.009).String())
	assert.Equal(t, "0.04", fixed.F128d2FromStringForced("0.0405").String())
	assert.Equal(t, "-0.04", fixed.F128d2FromStringForced("-0.0405").String())

	v, err := fixed.F128d2FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d2FromInt64(33))

	v, err = fixed.F128d2FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d2FromInt64(33))
}

func TestAddSub128d2(t *testing.T) {
	oneThird := fixed.F128d2FromStringForced("0.33")
	negTwoThirds := fixed.F128d2FromStringForced("-0.66")
	one := fixed.F128d2FromInt64(1)
	oneAndTwoThirds := fixed.F128d2FromStringForced("1.66")
	nineThousandSix := fixed.F128d2FromInt64(9006)
	ninetyPointZeroSix := fixed.F128d2FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F128d2FromStringForced("12.34")
	two := fixed.F128d2FromInt64(2)
	assert.Equal(t, "0.99", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.67", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.66", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, fixed.F128d2FromInt64(10240), fixed.F128d2FromInt64(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (fixed.F128d2FromInt64(1234).Add(nineThousandSix)).String())
	assert.Equal(t, fixed.F128d2FromStringForced("102.4"), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	assert.Equal(t, "102.4", (twelvePointThirtyFour.Add(ninetyPointZeroSix)).String())
	assert.Equal(t, "-1.5", (fixed.F128d2FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv128d2(t *testing.T) {
	pointThree := fixed.F128d2FromStringForced("0.3")
	negativePointThree := fixed.F128d2FromStringForced("-0.3")
	assert.Equal(t, "0.33", fixed.F128d2FromInt64(1).Div(fixed.F128d2FromInt64(3)).String())
	assert.Equal(t, "-0.33", fixed.F128d2FromInt64(1).Div(fixed.F128d2FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F128d2FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F128d2FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F128d2FromInt64(3)).String())
}

func TestTrunc128d2(t *testing.T) {
	assert.Equal(t, fixed.F128d2FromInt64(0), fixed.F128d2FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F128d2FromInt64(2), fixed.F128d2FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F128d2FromInt64(3), fixed.F128d2FromInt64(3).Trunc())
	assert.Equal(t, fixed.F128d2FromInt64(0), fixed.F128d2FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F128d2FromInt64(-2), fixed.F128d2FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F128d2FromInt64(-3), fixed.F128d2FromInt64(-3).Trunc())
}

func TestJSON128d2(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		testJSON128d2(t, fixed.F128d2FromInt64(i))
	}
	testJSON128d2(t, fixed.F128d2FromFloat64(18446744073712590000))
}

func testJSON128d2(t *testing.T, v fixed.F128d2) {
	t.Helper()
	e1 := embedded128d2{Field: v}
	data, err := json.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded128d2
	err = json.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}

func TestYAML128d2(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		testYAML128d2(t, fixed.F128d2FromInt64(i))
	}
	testYAML128d2(t, fixed.F128d2FromFloat64(18446744073712590000))
}

func testYAML128d2(t *testing.T, v fixed.F128d2) {
	t.Helper()
	e1 := embedded128d2{Field: v}
	data, err := yaml.Marshal(&e1)
	assert.NoError(t, err)
	var e2 embedded128d2
	err = yaml.Unmarshal(data, &e2)
	assert.NoError(t, err)
	require.Equal(t, e1, e2)
}
