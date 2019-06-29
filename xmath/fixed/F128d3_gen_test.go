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

type embedded128d3 struct {
	Field fixed.F128d3
}

func TestConversion128d3(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F128d3FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F128d3FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F128d3FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F128d3FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F128d3FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F128d3FromStringForced("-0.3").String())
	assert.Equal(t, "0.333", fixed.F128d3FromStringForced("0.33333").String())
	assert.Equal(t, "-0.333", fixed.F128d3FromStringForced("-0.33333").String())
	assert.Equal(t, "0.666", fixed.F128d3FromStringForced("0.66666").String())
	assert.Equal(t, "-0.666", fixed.F128d3FromStringForced("-0.66666").String())
	assert.Equal(t, "1", fixed.F128d3FromFloat64(1.0004).String())
	assert.Equal(t, "1", fixed.F128d3FromFloat64(1.00049).String())
	assert.Equal(t, "1", fixed.F128d3FromFloat64(1.0005).String())
	assert.Equal(t, "1", fixed.F128d3FromFloat64(1.0009).String())
	assert.Equal(t, "-1", fixed.F128d3FromFloat64(-1.0004).String())
	assert.Equal(t, "-1", fixed.F128d3FromFloat64(-1.00049).String())
	assert.Equal(t, "-1", fixed.F128d3FromFloat64(-1.0005).String())
	assert.Equal(t, "-1", fixed.F128d3FromFloat64(-1.0009).String())
	assert.Equal(t, "0.004", fixed.F128d3FromStringForced("0.00405").String())
	assert.Equal(t, "-0.004", fixed.F128d3FromStringForced("-0.00405").String())

	v, err := fixed.F128d3FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d3FromInt64(33))

	v, err = fixed.F128d3FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d3FromInt64(33))
}

func TestAddSub128d3(t *testing.T) {
	oneThird := fixed.F128d3FromStringForced("0.333")
	negTwoThirds := fixed.F128d3FromStringForced("-0.666")
	one := fixed.F128d3FromInt64(1)
	oneAndTwoThirds := fixed.F128d3FromStringForced("1.666")
	nineThousandSix := fixed.F128d3FromInt64(9006)
	ninetyPointZeroSix := fixed.F128d3FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F128d3FromStringForced("12.34")
	two := fixed.F128d3FromInt64(2)
	assert.Equal(t, "0.999", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.667", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.666", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, fixed.F128d3FromInt64(10240), fixed.F128d3FromInt64(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (fixed.F128d3FromInt64(1234).Add(nineThousandSix)).String())
	assert.Equal(t, fixed.F128d3FromStringForced("102.4"), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	assert.Equal(t, "102.4", (twelvePointThirtyFour.Add(ninetyPointZeroSix)).String())
	assert.Equal(t, "-1.5", (fixed.F128d3FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv128d3(t *testing.T) {
	pointThree := fixed.F128d3FromStringForced("0.3")
	negativePointThree := fixed.F128d3FromStringForced("-0.3")
	assert.Equal(t, "0.333", fixed.F128d3FromInt64(1).Div(fixed.F128d3FromInt64(3)).String())
	assert.Equal(t, "-0.333", fixed.F128d3FromInt64(1).Div(fixed.F128d3FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F128d3FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F128d3FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F128d3FromInt64(3)).String())
}

func TestTrunc128d3(t *testing.T) {
	assert.Equal(t, fixed.F128d3FromInt64(0), fixed.F128d3FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F128d3FromInt64(2), fixed.F128d3FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F128d3FromInt64(3), fixed.F128d3FromInt64(3).Trunc())
	assert.Equal(t, fixed.F128d3FromInt64(0), fixed.F128d3FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F128d3FromInt64(-2), fixed.F128d3FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F128d3FromInt64(-3), fixed.F128d3FromInt64(-3).Trunc())
}

func TestYAML128d3(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded128d3{Field: fixed.F128d3FromInt64(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d3
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON128d3(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded128d3{Field: fixed.F128d3FromInt64(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d3
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
