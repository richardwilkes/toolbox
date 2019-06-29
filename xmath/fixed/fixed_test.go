package fixed_test

import (
	"encoding/json"
	"testing"

	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type embedded struct {
	Field fixed.Fixed
}

func TestConversion(t *testing.T) {
	assert.Equal(t, "0.1", fixed.FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.FromFloat64(0.3).String())
	assert.Equal(t, "-0.1", fixed.FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.FromFloat64(-0.3).String())
	assert.Equal(t, "0.3333", fixed.FromFloat64(1.0/3.0).String())
	assert.Equal(t, "-0.3333", fixed.FromFloat64(-1.0/3.0).String())
	assert.Equal(t, "0.6666", fixed.FromFloat64(2.0/3.0).String())
	assert.Equal(t, "-0.6666", fixed.FromFloat64(-2.0/3.0).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00004).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.000049).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00005).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00009).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00004).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.000049).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00005).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00009).String())
	assert.Equal(t, "0.0004", fixed.FromFloat64(0.000405).String())
	assert.Equal(t, "-0.0004", fixed.FromFloat64(-0.000405).String())

	v, err := fixed.Parse("33.0")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.FromInt(33))

	v, err = fixed.Parse("33.00000000000000000000")
	assert.NoError(t, err)
	assert.Equal(t, v, fixed.FromInt(33))
}

func TestAddSub(t *testing.T) {
	oneThird := fixed.FromFloat64(1.0 / 3.0)
	negTwoThirds := fixed.FromFloat64(-2.0 / 3.0)
	assert.Equal(t, "0.9999", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.6667", (fixed.FromInt(1) - oneThird).String())
	assert.Equal(t, "-1.6666", (negTwoThirds - fixed.FromInt(1)).String())
	assert.Equal(t, "0", (negTwoThirds - fixed.FromInt(1) + fixed.FromFloat64(1.6666)).String())
	assert.Equal(t, fixed.FromInt(10240), fixed.FromInt(1234)+fixed.FromInt(9006))
	assert.Equal(t, "10240", (fixed.FromInt(1234) + fixed.FromInt(9006)).String())
	assert.Equal(t, fixed.FromFloat64(102.40), fixed.FromFloat64(12.34)+fixed.FromFloat64(90.06))
	assert.Equal(t, "102.4", (fixed.FromFloat64(12.34) + fixed.FromFloat64(90.06)).String())
	assert.Equal(t, "-1.5", (fixed.FromFloat64(0.5) - fixed.FromInt(2)).String())
}

func TestMulDiv(t *testing.T) {
	assert.Equal(t, "0.3333", fixed.FromInt(1).Div(fixed.FromInt(3)).String())
	assert.Equal(t, "-0.3333", fixed.FromInt(1).Div(fixed.FromInt(-3)).String())
	assert.Equal(t, "0.1", fixed.FromFloat64(0.3).Div(fixed.FromInt(3)).String())
	assert.Equal(t, "0.9", fixed.FromFloat64(0.3).Mul(fixed.FromInt(3)).String())
	assert.Equal(t, "-0.9", fixed.FromFloat64(-0.3).Mul(fixed.FromInt(3)).String())
}

func TestTrunc(t *testing.T) {
	assert.Equal(t, fixed.FromInt(0), fixed.FromFloat64(0.3333).Trunc())
	assert.Equal(t, fixed.FromInt(2), fixed.FromFloat64(2.6789).Trunc())
	assert.Equal(t, fixed.FromInt(3), fixed.FromInt(3).Trunc())
	assert.Equal(t, fixed.FromInt(0), fixed.FromFloat64(-0.3333).Trunc())
	assert.Equal(t, fixed.FromInt(-2), fixed.FromFloat64(-2.6789).Trunc())
	assert.Equal(t, fixed.FromInt(-3), fixed.FromInt(-3).Trunc())
}

func TestYAML(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		e1 := embedded{Field: fixed.Fixed(i)}
		data, err := yaml.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded
		err = yaml.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestJSON(t *testing.T) {
	for i := -25000; i < 25001; i += 13 {
		e1 := embedded{Field: fixed.Fixed(i)}
		data, err := json.Marshal(&e1)
		assert.NoError(t, err)
		var e2 embedded
		err = json.Unmarshal(data, &e2)
		assert.NoError(t, err)
		require.Equal(t, e1, e2)
	}
}

func TestPrecision(t *testing.T) {
	fixed.SetDigitsAfterDecimal(2)

	assert.Equal(t, "0.1", fixed.FromFloat64(0.1).String())
	assert.Equal(t, "0.2", fixed.FromFloat64(0.2).String())
	assert.Equal(t, "0.3", fixed.FromFloat64(0.3).String())
	assert.Equal(t, "-0.1", fixed.FromFloat64(-0.1).String())
	assert.Equal(t, "-0.2", fixed.FromFloat64(-0.2).String())
	assert.Equal(t, "-0.3", fixed.FromFloat64(-0.3).String())
	assert.Equal(t, "0.33", fixed.FromFloat64(1.0/3.0).String())
	assert.Equal(t, "-0.33", fixed.FromFloat64(-1.0/3.0).String())
	assert.Equal(t, "0.66", fixed.FromFloat64(2.0/3.0).String())
	assert.Equal(t, "-0.66", fixed.FromFloat64(-2.0/3.0).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00004).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.000049).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00005).String())
	assert.Equal(t, "1", fixed.FromFloat64(1.00009).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00004).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.000049).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00005).String())
	assert.Equal(t, "-1", fixed.FromFloat64(-1.00009).String())
	assert.Equal(t, "0", fixed.FromFloat64(0.000405).String())
	assert.Equal(t, "0", fixed.FromFloat64(-0.000405).String())

	oneThird := fixed.FromFloat64(1.0 / 3.0)
	negTwoThirds := fixed.FromFloat64(-2.0 / 3.0)
	assert.Equal(t, "0.99", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.67", (fixed.FromInt(1) - oneThird).String())
	assert.Equal(t, "-1.66", (negTwoThirds - fixed.FromInt(1)).String())
	assert.Equal(t, "0", (negTwoThirds - fixed.FromInt(1) + fixed.FromFloat64(1.6666)).String())
	assert.Equal(t, fixed.FromInt(10240), fixed.FromInt(1234)+fixed.FromInt(9006))
	assert.Equal(t, "10240", (fixed.FromInt(1234) + fixed.FromInt(9006)).String())
	assert.Equal(t, fixed.FromFloat64(102.40), fixed.FromFloat64(12.34)+fixed.FromFloat64(90.06))
	assert.Equal(t, "102.4", (fixed.FromFloat64(12.34) + fixed.FromFloat64(90.06)).String())
	assert.Equal(t, "-1.5", (fixed.FromFloat64(0.5) - fixed.FromInt(2)).String())

	assert.Equal(t, "0.33", fixed.FromInt(1).Div(fixed.FromInt(3)).String())
	assert.Equal(t, "-0.33", fixed.FromInt(1).Div(fixed.FromInt(-3)).String())
	assert.Equal(t, "0.1", fixed.FromFloat64(0.3).Div(fixed.FromInt(3)).String())
	assert.Equal(t, "0.9", fixed.FromFloat64(0.3).Mul(fixed.FromInt(3)).String())
	assert.Equal(t, "-0.9", fixed.FromFloat64(-0.3).Mul(fixed.FromInt(3)).String())

	fixed.SetDigitsAfterDecimal(6)

	assert.Equal(t, "0.000405", fixed.FromFloat64(0.000405).String())
	assert.Equal(t, "-0.000405", fixed.FromFloat64(-0.000405).String())

	fixed.SetDigitsAfterDecimal(4)
}
