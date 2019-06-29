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

type embedded64d3 struct {
	Field fixed.F64d3
}

func TestConversion64d3(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F64d3FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F64d3FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F64d3FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F64d3FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F64d3FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F64d3FromStringForced("-0.3").String())
	assert.Equal(t, "0.333", fixed.F64d3FromStringForced("0.33333").String())
	assert.Equal(t, "-0.333", fixed.F64d3FromStringForced("-0.33333").String())
	assert.Equal(t, "0.666", fixed.F64d3FromStringForced("0.66666").String())
	assert.Equal(t, "-0.666", fixed.F64d3FromStringForced("-0.66666").String())
	assert.Equal(t, "1", fixed.F64d3FromFloat64(1.0004).String())
	assert.Equal(t, "1", fixed.F64d3FromFloat64(1.00049).String())
	assert.Equal(t, "1", fixed.F64d3FromFloat64(1.0005).String())
	assert.Equal(t, "1", fixed.F64d3FromFloat64(1.0009).String())
	assert.Equal(t, "-1", fixed.F64d3FromFloat64(-1.0004).String())
	assert.Equal(t, "-1", fixed.F64d3FromFloat64(-1.00049).String())
	assert.Equal(t, "-1", fixed.F64d3FromFloat64(-1.0005).String())
	assert.Equal(t, "-1", fixed.F64d3FromFloat64(-1.0009).String())
	assert.Equal(t, "0.004", fixed.F64d3FromStringForced("0.00405").String())
	assert.Equal(t, "-0.004", fixed.F64d3FromStringForced("-0.00405").String())

	v, err := fixed.F64d3FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d3FromInt64(33))

	v, err = fixed.F64d3FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d3FromInt64(33))
}

func TestAddSub64d3(t *testing.T) {
	oneThird := fixed.F64d3FromStringForced("0.333")
	negTwoThirds := fixed.F64d3FromStringForced("-0.666")
	one := fixed.F64d3FromInt64(1)
	oneAndTwoThirds := fixed.F64d3FromStringForced("1.666")
	nineThousandSix := fixed.F64d3FromInt64(9006)
	ninetyPointZeroSix := fixed.F64d3FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F64d3FromStringForced("12.34")
	two := fixed.F64d3FromInt64(2)
	assert.Equal(t, "0.999", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.667", (one - oneThird).String())
	assert.Equal(t, "-1.666", (negTwoThirds - one).String())
	assert.Equal(t, "0", (negTwoThirds - one + oneAndTwoThirds).String())
	assert.Equal(t, fixed.F64d3FromInt64(10240), fixed.F64d3FromInt64(1234)+nineThousandSix)
	assert.Equal(t, "10240", (fixed.F64d3FromInt64(1234) + nineThousandSix).String())
	assert.Equal(t, fixed.F64d3FromStringForced("102.4"), twelvePointThirtyFour+ninetyPointZeroSix)
	assert.Equal(t, "102.4", (twelvePointThirtyFour + ninetyPointZeroSix).String())
	assert.Equal(t, "-1.5", (fixed.F64d3FromFloat64(0.5) - two).String())
}

func TestMulDiv64d3(t *testing.T) {
	pointThree := fixed.F64d3FromStringForced("0.3")
	negativePointThree := fixed.F64d3FromStringForced("-0.3")
	assert.Equal(t, "0.333", fixed.F64d3FromInt64(1).Div(fixed.F64d3FromInt64(3)).String())
	assert.Equal(t, "-0.333", fixed.F64d3FromInt64(1).Div(fixed.F64d3FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F64d3FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F64d3FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F64d3FromInt64(3)).String())
}

func TestTrunc64d3(t *testing.T) {
	assert.Equal(t, fixed.F64d3FromInt64(0), fixed.F64d3FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F64d3FromInt64(2), fixed.F64d3FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F64d3FromInt64(3), fixed.F64d3FromInt64(3).Trunc())
	assert.Equal(t, fixed.F64d3FromInt64(0), fixed.F64d3FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F64d3FromInt64(-2), fixed.F64d3FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F64d3FromInt64(-3), fixed.F64d3FromInt64(-3).Trunc())
}

func TestYAML64d3(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded64d3{Field: fixed.F64d3FromInt64(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded64d3
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON64d3(t *testing.T) {
	for i := int64(-25000); i < 25001; i += 13 {
		e1 := embedded64d3{Field: fixed.F64d3FromInt64(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded64d3
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
