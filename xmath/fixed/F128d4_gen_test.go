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
	assert.Equal(t, v, fixed.F128d4FromInt64(33))

	v, err = fixed.F128d4FromString("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.F128d4FromInt64(33))
}

func TestAddSub128d4(t *testing.T) {
	oneThird := fixed.F128d4FromStringForced("0.3333")
	negTwoThirds := fixed.F128d4FromStringForced("-0.6666")
	one := fixed.F128d4FromInt64(1)
	oneAndTwoThirds := fixed.F128d4FromStringForced("1.6666")
	nineThousandSix := fixed.F128d4FromInt64(9006)
	ninetyPointZeroSix := fixed.F128d4FromStringForced("90.06")
	twelvePointThirtyFour := fixed.F128d4FromStringForced("12.34")
	two := fixed.F128d4FromInt64(2)
	assert.Equal(t, "0.9999", (oneThird.Add(oneThird).Add(oneThird)).String())
	assert.Equal(t, "0.6667", (one.Sub(oneThird)).String())
	assert.Equal(t, "-1.6666", (negTwoThirds.Sub(one)).String())
	assert.Equal(t, "0", (negTwoThirds.Sub(one).Add(oneAndTwoThirds)).String())
	assert.Equal(t, fixed.F128d4FromInt64(10240), fixed.F128d4FromInt64(1234).Add(nineThousandSix))
	assert.Equal(t, "10240", (fixed.F128d4FromInt64(1234).Add(nineThousandSix)).String())
	assert.Equal(t, fixed.F128d4FromStringForced("102.4"), twelvePointThirtyFour.Add(ninetyPointZeroSix))
	assert.Equal(t, "102.4", (twelvePointThirtyFour.Add(ninetyPointZeroSix)).String())
	assert.Equal(t, "-1.5", (fixed.F128d4FromFloat64(0.5).Sub(two)).String())
}

func TestMulDiv128d4(t *testing.T) {
	pointThree := fixed.F128d4FromStringForced("0.3")
	negativePointThree := fixed.F128d4FromStringForced("-0.3")
	assert.Equal(t, "0.3333", fixed.F128d4FromInt64(1).Div(fixed.F128d4FromInt64(3)).String())
	assert.Equal(t, "-0.3333", fixed.F128d4FromInt64(1).Div(fixed.F128d4FromInt64(-3)).String())
	assert.Equal(t, "0.1", pointThree.Div(fixed.F128d4FromInt64(3)).String())
	assert.Equal(t, "0.9", pointThree.Mul(fixed.F128d4FromInt64(3)).String())
	assert.Equal(t, "-0.9", negativePointThree.Mul(fixed.F128d4FromInt64(3)).String())
}

func TestTrunc128d4(t *testing.T) {
	assert.Equal(t, fixed.F128d4FromInt64(0), fixed.F128d4FromStringForced("0.3333").Trunc())
	assert.Equal(t, fixed.F128d4FromInt64(2), fixed.F128d4FromStringForced("2.6789").Trunc())
	assert.Equal(t, fixed.F128d4FromInt64(3), fixed.F128d4FromInt64(3).Trunc())
	assert.Equal(t, fixed.F128d4FromInt64(0), fixed.F128d4FromStringForced("-0.3333").Trunc())
	assert.Equal(t, fixed.F128d4FromInt64(-2), fixed.F128d4FromStringForced("-2.6789").Trunc())
	assert.Equal(t, fixed.F128d4FromInt64(-3), fixed.F128d4FromInt64(-3).Trunc())
}

func TestText128d4(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		f1 := fixed.F128d4FromInt64(i)
		data, err := f1.MarshalText()
		assert.NoError(t, err)
		var f2 fixed.F128d4
		err = f2.UnmarshalText(data)
		assert.NoError(t, err)
		require.Equal(t, f1, f2)
	}
}

func TestYAML128d4(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		e1 := embedded128d4{Field: fixed.F128d4FromInt64(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d4
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON128d4(t *testing.T) {
	for i := int64(-20000); i < 20001; i++ {
		e1 := embedded128d4{Field: fixed.F128d4FromInt64(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded128d4
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}