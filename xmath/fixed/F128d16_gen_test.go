// Code created from "fixed_test.go.tmpl" - don't edit by hand

package fixed_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type embedded128d16 struct {
	Field fixed.F128d16
}

func TestConversion128d16(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F128d16FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F128d16FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F128d16FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F128d16FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F128d16FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F128d16FromStringForced("-0.3").String())
	assert.Equal(t, "0.3333333333333333", fixed.F128d16FromStringForced("0.333333333333333333").String())
	assert.Equal(t, "-0.3333333333333333", fixed.F128d16FromStringForced("-0.333333333333333333").String())
	assert.Equal(t, "0.6666666666666666", fixed.F128d16FromStringForced("0.666666666666666666").String())
	assert.Equal(t, "-0.6666666666666666", fixed.F128d16FromStringForced("-0.666666666666666666").String())
	assert.Equal(t, "1", fixed.F128d16FromFloat64(1.00000000000000004).String())
	assert.Equal(t, "1", fixed.F128d16FromFloat64(1.000000000000000049).String())
	assert.Equal(t, "1", fixed.F128d16FromFloat64(1.00000000000000005).String())
	assert.Equal(t, "1", fixed.F128d16FromFloat64(1.00000000000000009).String())
	assert.Equal(t, "-1", fixed.F128d16FromFloat64(-1.00000000000000004).String())
	assert.Equal(t, "-1", fixed.F128d16FromFloat64(-1.000000000000000049).String())
	assert.Equal(t, "-1", fixed.F128d16FromFloat64(-1.00000000000000005).String())
	assert.Equal(t, "-1", fixed.F128d16FromFloat64(-1.00000000000000009).String())
	assert.Equal(t, "0.0000000000000004", fixed.F128d16FromStringForced("0.000000000000000405").String())
	assert.Equal(t, "-0.0000000000000004", fixed.F128d16FromStringForced("-0.000000000000000405").String())

	v, err := fixed.F128d16FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d16FromInt64(33))

	v, err = fixed.F128d16FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d16FromInt64(33))
}

func TestAddSub128d16(t *testing.T) {
	oneThird := fixed.F128d16FromStringForced("0.3333333333333333")
	negTwoThirds := fixed.F128d16FromStringForced("-0.6666666666666666")
	one := fixed.F128d16FromInt64(1)
	oneAndTwoThirds := fixed.F128d16FromStringForced("1.6666666666666666")
	nineThousandSix := fixed.F128d16FromInt64(9006)
	ninetyPointZeroSix := fixed.F128d16FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F128d16FromStringForced("12.34")
	two := fixed.F128d16FromInt64(2)
	assert.Equal(t, "0.9999999999999999", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.6666666666666667", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.6666666666666666", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, fixed.F128d16FromInt64(10240), fixed.F128d16FromInt64(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (fixed.F128d16FromInt64(1234).Add(nineThousandSix)).String())
	assert.Equal(t, fixed.F128d16FromStringForced("102.4"), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	assert.Equal(t, "102.4", (twelvePointThirtyFour.Add(ninetyPointZeroSix)).String())
	assert.Equal(t, "-1.5", (fixed.F128d16FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv128d16(t *testing.T) {
	pointThree := fixed.F128d16FromStringForced("0.3")
	negativePointThree := fixed.F128d16FromStringForced("-0.3")
	assert.Equal(t, "0.3333333333333333", fixed.F128d16FromInt64(1).Div(fixed.F128d16FromInt64(3)).String())
	assert.Equal(t, "-0.3333333333333333", fixed.F128d16FromInt64(1).Div(fixed.F128d16FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F128d16FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F128d16FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F128d16FromInt64(3)).String())
}

func TestTrunc128d16(t *testing.T) {
	assert.Equal(t, fixed.F128d16FromInt64(0), fixed.F128d16FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F128d16FromInt64(2), fixed.F128d16FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F128d16FromInt64(3), fixed.F128d16FromInt64(3).Trunc())
	assert.Equal(t, fixed.F128d16FromInt64(0), fixed.F128d16FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F128d16FromInt64(-2), fixed.F128d16FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F128d16FromInt64(-3), fixed.F128d16FromInt64(-3).Trunc())
}

func TestYAML128d16(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded128d16{Field: fixed.F128d16FromInt64(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d16
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON128d16(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded128d16{Field: fixed.F128d16FromInt64(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d16
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
