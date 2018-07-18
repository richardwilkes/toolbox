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
	assert.Equal(t, "0.1", fixed.New(0.1).String())
	assert.Equal(t, "0.2", fixed.New(0.2).String())
	assert.Equal(t, "0.3", fixed.New(0.3).String())
	assert.Equal(t, "-0.1", fixed.New(-0.1).String())
	assert.Equal(t, "-0.2", fixed.New(-0.2).String())
	assert.Equal(t, "-0.3", fixed.New(-0.3).String())
	assert.Equal(t, "0.3333", fixed.New(1.0/3.0).String())
	assert.Equal(t, "-0.3333", fixed.New(-1.0/3.0).String())
	assert.Equal(t, "0.6666", fixed.New(2.0/3.0).String())
	assert.Equal(t, "-0.6666", fixed.New(-2.0/3.0).String())
	assert.Equal(t, "1", fixed.New(1.00004).String())
	assert.Equal(t, "1", fixed.New(1.000049).String())
	assert.Equal(t, "1", fixed.New(1.00005).String())
	assert.Equal(t, "1", fixed.New(1.00009).String())
	assert.Equal(t, "-1", fixed.New(-1.00004).String())
	assert.Equal(t, "-1", fixed.New(-1.000049).String())
	assert.Equal(t, "-1", fixed.New(-1.00005).String())
	assert.Equal(t, "-1", fixed.New(-1.00009).String())
	assert.Equal(t, "0.0004", fixed.New(0.000405).String())
	assert.Equal(t, "-0.0004", fixed.New(-0.000405).String())
}

func TestAddSub(t *testing.T) {
	oneThird := fixed.New(1.0 / 3.0)
	negTwoThirds := fixed.New(-2.0 / 3.0)
	assert.Equal(t, "0.9999", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.6667", (fixed.New(1) - oneThird).String())
	assert.Equal(t, "-1.6666", (negTwoThirds - fixed.New(1)).String())
	assert.Equal(t, "0", (negTwoThirds - fixed.New(1) + fixed.New(1.6666)).String())
	assert.Equal(t, fixed.New(10240), fixed.New(1234)+fixed.New(9006))
	assert.Equal(t, "10240", (fixed.New(1234) + fixed.New(9006)).String())
	assert.Equal(t, fixed.New(102.40), fixed.New(12.34)+fixed.New(90.06))
	assert.Equal(t, "102.4", (fixed.New(12.34) + fixed.New(90.06)).String())
	assert.Equal(t, "-1.5", (fixed.New(0.5) - fixed.New(2)).String())
}

func TestMulDiv(t *testing.T) {
	assert.Equal(t, "0.3333", fixed.New(1).Div(fixed.New(3)).String())
	assert.Equal(t, "-0.3333", fixed.New(1).Div(fixed.New(-3)).String())
	assert.Equal(t, "0.1", fixed.New(0.3).Div(fixed.New(3)).String())
	assert.Equal(t, "0.9", fixed.New(0.3).Mul(fixed.New(3)).String())
	assert.Equal(t, "-0.9", fixed.New(-0.3).Mul(fixed.New(3)).String())
}

func TestTrunc(t *testing.T) {
	assert.Equal(t, fixed.New(0), fixed.New(0.3333).Trunc())
	assert.Equal(t, fixed.New(2), fixed.New(2.6789).Trunc())
	assert.Equal(t, fixed.New(3), fixed.New(3).Trunc())
	assert.Equal(t, fixed.New(0), fixed.New(-0.3333).Trunc())
	assert.Equal(t, fixed.New(-2), fixed.New(-2.6789).Trunc())
	assert.Equal(t, fixed.New(-3), fixed.New(-3).Trunc())
}

func TestText(t *testing.T) {
	for i := -20000; i < 20001; i++ {
		f1 := fixed.Fixed(i)
		data, err := f1.MarshalText()
		assert.NoError(t, err)
		var f2 fixed.Fixed
		err = f2.UnmarshalText(data)
		assert.NoError(t, err)
		require.Equal(t, f1, f2)
	}
}

func TestYAML(t *testing.T) {
	for i := -20000; i < 20001; i++ {
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
	for i := -20000; i < 20001; i++ {
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

	assert.Equal(t, "0.1", fixed.New(0.1).String())
	assert.Equal(t, "0.2", fixed.New(0.2).String())
	assert.Equal(t, "0.3", fixed.New(0.3).String())
	assert.Equal(t, "-0.1", fixed.New(-0.1).String())
	assert.Equal(t, "-0.2", fixed.New(-0.2).String())
	assert.Equal(t, "-0.3", fixed.New(-0.3).String())
	assert.Equal(t, "0.33", fixed.New(1.0/3.0).String())
	assert.Equal(t, "-0.33", fixed.New(-1.0/3.0).String())
	assert.Equal(t, "0.66", fixed.New(2.0/3.0).String())
	assert.Equal(t, "-0.66", fixed.New(-2.0/3.0).String())
	assert.Equal(t, "1", fixed.New(1.00004).String())
	assert.Equal(t, "1", fixed.New(1.000049).String())
	assert.Equal(t, "1", fixed.New(1.00005).String())
	assert.Equal(t, "1", fixed.New(1.00009).String())
	assert.Equal(t, "-1", fixed.New(-1.00004).String())
	assert.Equal(t, "-1", fixed.New(-1.000049).String())
	assert.Equal(t, "-1", fixed.New(-1.00005).String())
	assert.Equal(t, "-1", fixed.New(-1.00009).String())
	assert.Equal(t, "0", fixed.New(0.000405).String())
	assert.Equal(t, "0", fixed.New(-0.000405).String())

	oneThird := fixed.New(1.0 / 3.0)
	negTwoThirds := fixed.New(-2.0 / 3.0)
	assert.Equal(t, "0.99", (oneThird + oneThird + oneThird).String())
	assert.Equal(t, "0.67", (fixed.New(1) - oneThird).String())
	assert.Equal(t, "-1.66", (negTwoThirds - fixed.New(1)).String())
	assert.Equal(t, "0", (negTwoThirds - fixed.New(1) + fixed.New(1.6666)).String())
	assert.Equal(t, fixed.New(10240), fixed.New(1234)+fixed.New(9006))
	assert.Equal(t, "10240", (fixed.New(1234) + fixed.New(9006)).String())
	assert.Equal(t, fixed.New(102.40), fixed.New(12.34)+fixed.New(90.06))
	assert.Equal(t, "102.4", (fixed.New(12.34) + fixed.New(90.06)).String())
	assert.Equal(t, "-1.5", (fixed.New(0.5) - fixed.New(2)).String())

	assert.Equal(t, "0.33", fixed.New(1).Div(fixed.New(3)).String())
	assert.Equal(t, "-0.33", fixed.New(1).Div(fixed.New(-3)).String())
	assert.Equal(t, "0.1", fixed.New(0.3).Div(fixed.New(3)).String())
	assert.Equal(t, "0.9", fixed.New(0.3).Mul(fixed.New(3)).String())
	assert.Equal(t, "-0.9", fixed.New(-0.3).Mul(fixed.New(3)).String())

	fixed.SetDigitsAfterDecimal(6)

	assert.Equal(t, "0.000405", fixed.New(0.000405).String())
	assert.Equal(t, "-0.000405", fixed.New(-0.000405).String())

	fixed.SetDigitsAfterDecimal(4)
}
