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

type embedded64d2 struct {
	Field fixed.F64d2
}

func TestConversion64d2(t *testing.T) {
	assert.Equal(t, "0.1", fixed.F64d2FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.F64d2FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.F64d2FromStringForced("0.3").String())
	assert.Equal(t, "-0.1", fixed.F64d2FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.F64d2FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.F64d2FromStringForced("-0.3").String())
	assert.Equal(t, "0.33", fixed.F64d2FromStringForced("0.3333").String())
	assert.Equal(t, "-0.33", fixed.F64d2FromStringForced("-0.3333").String())
	assert.Equal(t, "0.66", fixed.F64d2FromStringForced("0.6666").String())
	assert.Equal(t, "-0.66", fixed.F64d2FromStringForced("-0.6666").String())
	assert.Equal(t, "1", fixed.F64d2FromFloat64(1.004).String())
	assert.Equal(t, "1", fixed.F64d2FromFloat64(1.0049).String())
	assert.Equal(t, "1", fixed.F64d2FromFloat64(1.005).String())
	assert.Equal(t, "1", fixed.F64d2FromFloat64(1.009).String())
	assert.Equal(t, "-1", fixed.F64d2FromFloat64(-1.004).String())
	assert.Equal(t, "-1", fixed.F64d2FromFloat64(-1.0049).String())
	assert.Equal(t, "-1", fixed.F64d2FromFloat64(-1.005).String())
	assert.Equal(t, "-1", fixed.F64d2FromFloat64(-1.009).String())
	assert.Equal(t, "0.04", fixed.F64d2FromStringForced("0.0405").String())
	assert.Equal(t, "-0.04", fixed.F64d2FromStringForced("-0.0405").String())

	v, err := fixed.F64d2FromString("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d2FromInt64(33))

	v, err = fixed.F64d2FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F64d2FromInt64(33))
}

func TestAddSub64d2(t *testing.T) {
	oneThird := fixed.F64d2FromStringForced("0.33")
	negTwoThirds := fixed.F64d2FromStringForced("-0.66")
	one := fixed.F64d2FromInt64(1)
	oneAndTwoThirds := fixed.F64d2FromStringForced("1.66")
	nineThousandSix := fixed.F64d2FromInt64(9006)
	ninetyPointZeroSix := fixed.F64d2FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F64d2FromStringForced("12.34")
	two := fixed.F64d2FromInt64(2)
	assert.Equal(t, "0.99", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.67", (one - oneThird).String())
	assert.Equal(t, "-1.66", (negTwoThirds - one).String())
	assert.Equal(t, "0", (negTwoThirds - one + oneAndTwoThirds).String())
	assert.Equal(t, fixed.F64d2FromInt64(10240), fixed.F64d2FromInt64(1234)+nineThousandSix)
	assert.Equal(t, "10240", (fixed.F64d2FromInt64(1234) + nineThousandSix).String())
	assert.Equal(t, fixed.F64d2FromStringForced("102.4"), twelvePointThirtyFour+ninetyPointZeroSix)
	assert.Equal(t, "102.4", (twelvePointThirtyFour + ninetyPointZeroSix).String())
	assert.Equal(t, "-1.5", (fixed.F64d2FromFloat64(0.5) - two).String())
}

func TestMulDiv64d2(t *testing.T) {
	pointThree := fixed.F64d2FromStringForced("0.3")
	negativePointThree := fixed.F64d2FromStringForced("-0.3")
	assert.Equal(t, "0.33", fixed.F64d2FromInt64(1).Div(fixed.F64d2FromInt64(3)).String())
	assert.Equal(t, "-0.33", fixed.F64d2FromInt64(1).Div(fixed.F64d2FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F64d2FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F64d2FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F64d2FromInt64(3)).String())
}

func TestTrunc64d2(t *testing.T) {
	assert.Equal(t, fixed.F64d2FromInt64(0), fixed.F64d2FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F64d2FromInt64(2), fixed.F64d2FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F64d2FromInt64(3), fixed.F64d2FromInt64(3).Trunc())
	assert.Equal(t, fixed.F64d2FromInt64(0), fixed.F64d2FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F64d2FromInt64(-2), fixed.F64d2FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F64d2FromInt64(-3), fixed.F64d2FromInt64(-3).Trunc())
}

func TestText64d2(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		f1 := fixed.F64d2FromInt64(i)
		data, err := f1.MarshalText()
		assert.NoError(t, err)
		var f2 fixed.F64d2
		err = f2.UnmarshalText(data)
		assert.NoError(t, err)
		require.Equal(t, f1, f2)
	}
}

func TestYAML64d2(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		e1 := embedded64d2{Field: fixed.F64d2FromInt64(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded64d2
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON64d2(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		e1 := embedded64d2{Field: fixed.F64d2FromInt64(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded64d2
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}
