package num_test

import (
	"encoding/json"
	"math"
	"math/big"
	"strconv"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/richardwilkes/toolbox/xmath/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	maxUint64PlusOneAsStr  = "18446744073709551616"
	maxUint128AsStr        = "340282366920938463463374607431768211455"
	maxUint128PlusOneAsStr = "340282366920938463463374607431768211456"
)

var (
	utable = []*uinfo{
		{
			IsUint64:  true,
			IsUint128: true,
		},
		{
			Uint64:    1,
			IsUint64:  true,
			IsUint128: true,
		},
		{
			ValueAsStr: "18446744073712590000",
			IsUint128:  true,
		},
		{
			Uint64:    math.MaxUint64,
			IsUint64:  true,
			IsUint128: true,
		},
		{
			ValueAsStr: maxUint64PlusOneAsStr,
			IsUint128:  true,
		},
		{
			ValueAsStr: maxUint128AsStr,
			IsUint128:  true,
		},
		{
			ValueAsStr:              maxUint128PlusOneAsStr,
			ExpectedConversionAsStr: maxUint128AsStr,
		},
	}
)

type uinfo struct {
	Uint64                  uint64
	ValueAsStr              string
	ExpectedConversionAsStr string
	IsUint64                bool
	IsUint128               bool
}

func init() {
	for _, d := range utable {
		if d.IsUint64 {
			d.ValueAsStr = strconv.FormatUint(d.Uint64, 10)
		}
		if d.ExpectedConversionAsStr == "" {
			d.ExpectedConversionAsStr = d.ValueAsStr
		}
	}
}

func bigUintFromStr(t *testing.T, one *uinfo, index int) *big.Int {
	t.Helper()
	b, ok := new(big.Int).SetString(one.ValueAsStr, 10)
	require.True(t, ok, indexFmt, index)
	require.Equal(t, one.ValueAsStr, b.String(), indexFmt, index)
	return b
}

func TestUint128FromUint64(t *testing.T) {
	for i, one := range utable {
		if one.IsUint64 {
			assert.Equal(t, one.ExpectedConversionAsStr, num.Uint128From64(one.Uint64).String(), indexFmt, i)
		}
	}
}

func TestUint128FromBigInt(t *testing.T) {
	for i, one := range utable {
		assert.Equal(t, one.ExpectedConversionAsStr, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).String(), indexFmt, i)
	}
}

func TestUint128AsBigInt(t *testing.T) {
	for i, one := range utable {
		if one.IsUint128 {
			assert.Equal(t, one.ValueAsStr, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).AsBigInt().String(), indexFmt, i)
		}
	}
}

func TestUint128AsUint64(t *testing.T) {
	for i, one := range utable {
		if one.IsUint64 {
			assert.Equal(t, one.Uint64, num.Uint128From64(one.Uint64).AsUint64(), indexFmt, i)
		}
	}
}

func TestUint128IsUint64(t *testing.T) {
	for i, one := range utable {
		if one.IsUint128 {
			assert.Equal(t, one.IsUint64, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).IsUint64(), indexFmt, i)
		}
	}
}

func TestUint128Inc(t *testing.T) {
	big1 := new(big.Int).SetInt64(1)
	for i, one := range utable {
		if one.IsUint128 {
			b := bigUintFromStr(t, one, i)
			v := num.Uint128FromBigInt(b)
			if v == num.MaxUint128 {
				assert.Equal(t, num.Uint128{}, v.Inc(), indexFmt, i)
			} else {
				b.Add(b, big1)
				assert.Equal(t, b.String(), v.Inc().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestUint128Dec(t *testing.T) {
	big1 := new(big.Int).SetInt64(1)
	for i, one := range utable {
		if one.IsUint128 {
			b := bigUintFromStr(t, one, i)
			v := num.Uint128FromBigInt(b)
			if v.IsZero() {
				assert.Equal(t, num.MaxUint128, v.Dec(), indexFmt, i)
			} else {
				b.Sub(b, big1)
				assert.Equal(t, b.String(), v.Dec().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestUint128Add(t *testing.T) {
	assert.Equal(t, num.Uint128From64(0), num.Uint128From64(0).Add(num.Uint128From64(0)))
	assert.Equal(t, num.Uint128From64(120), num.Uint128From64(22).Add(num.Uint128From64(98)))
	assert.Equal(t, num.Uint128FromComponents(1, 0), num.Uint128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Add(num.Uint128From64(1)))
	assert.Equal(t, num.Uint128From64(0), num.MaxUint128.Add(num.Uint128From64(1)))
}

func TestUint128Sub(t *testing.T) {
	assert.Equal(t, num.Uint128From64(0), num.Uint128From64(0).Sub(num.Uint128From64(0)))
	assert.Equal(t, num.Uint128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Uint128FromComponents(1, 0).Sub(num.Uint128From64(1)))
	assert.Equal(t, num.MaxUint128, num.Uint128From64(0).Sub(num.Uint128From64(1)))
}

func TestUint128Cmp(t *testing.T) {
	assert.Equal(t, 0, num.Uint128From64(0).Cmp(num.Uint128From64(0)))
	assert.Equal(t, -1, num.Uint128From64(1).Cmp(num.Uint128From64(2)))
	assert.Equal(t, -1, num.Uint128From64(22).Cmp(num.Uint128From64(98)))
	assert.Equal(t, 1, num.Uint128FromComponents(1, 0).Cmp(num.Uint128From64(1)))
	assert.Equal(t, -1, num.Uint128From64(0).Cmp(num.MaxUint128))
	assert.Equal(t, 1, num.MaxUint128.Cmp(num.Uint128From64(0)))
	assert.Equal(t, 0, num.MaxUint128.Cmp(num.MaxUint128))
}

func TestUint128GreaterThan(t *testing.T) {
	assert.Equal(t, false, num.Uint128From64(0).GreaterThan(num.Uint128From64(0)))
	assert.Equal(t, false, num.Uint128From64(1).GreaterThan(num.Uint128From64(2)))
	assert.Equal(t, false, num.Uint128From64(22).GreaterThan(num.Uint128From64(98)))
	assert.Equal(t, false, num.Uint128From64(0).GreaterThan(num.MaxUint128))
	assert.Equal(t, false, num.MaxUint128.GreaterThan(num.MaxUint128))
	assert.Equal(t, true, num.Uint128FromComponents(1, 0).GreaterThan(num.Uint128From64(1)))
	assert.Equal(t, true, num.MaxUint128.GreaterThan(num.Uint128From64(0)))
}

func TestUint128GreaterOrEqualTo(t *testing.T) {
	assert.Equal(t, true, num.Uint128From64(0).GreaterOrEqualTo(num.Uint128From64(0)))
	assert.Equal(t, false, num.Uint128From64(1).GreaterOrEqualTo(num.Uint128From64(2)))
	assert.Equal(t, false, num.Uint128From64(22).GreaterOrEqualTo(num.Uint128From64(98)))
	assert.Equal(t, false, num.Uint128From64(0).GreaterOrEqualTo(num.Uint128From64(1)))
	assert.Equal(t, false, num.Uint128From64(0).GreaterOrEqualTo(num.MaxUint128))
	assert.Equal(t, true, num.MaxUint128.GreaterOrEqualTo(num.MaxUint128))
	assert.Equal(t, true, num.Uint128FromComponents(1, 0).GreaterOrEqualTo(num.Uint128From64(1)))
	assert.Equal(t, true, num.MaxUint128.GreaterOrEqualTo(num.Uint128From64(0)))
}

func TestUint128LessThan(t *testing.T) {
	assert.Equal(t, false, num.Uint128From64(0).LessThan(num.Uint128From64(0)))
	assert.Equal(t, true, num.Uint128From64(1).LessThan(num.Uint128From64(2)))
	assert.Equal(t, true, num.Uint128From64(22).LessThan(num.Uint128From64(98)))
	assert.Equal(t, true, num.Uint128From64(0).LessThan(num.Uint128From64(1)))
	assert.Equal(t, true, num.Uint128From64(0).LessThan(num.MaxUint128))
	assert.Equal(t, false, num.MaxUint128.LessThan(num.MaxUint128))
	assert.Equal(t, false, num.Uint128FromComponents(1, 0).LessThan(num.Uint128From64(1)))
	assert.Equal(t, false, num.MaxUint128.LessThan(num.Uint128From64(0)))
}

func TestUint128LessOrEqualTo(t *testing.T) {
	assert.Equal(t, true, num.Uint128From64(0).LessOrEqualTo(num.Uint128From64(0)))
	assert.Equal(t, true, num.Uint128From64(1).LessOrEqualTo(num.Uint128From64(2)))
	assert.Equal(t, true, num.Uint128From64(22).LessOrEqualTo(num.Uint128From64(98)))
	assert.Equal(t, true, num.Uint128From64(0).LessOrEqualTo(num.Uint128From64(1)))
	assert.Equal(t, true, num.Uint128From64(0).LessOrEqualTo(num.MaxUint128))
	assert.Equal(t, true, num.MaxUint128.LessOrEqualTo(num.MaxUint128))
	assert.Equal(t, false, num.Uint128FromComponents(1, 0).LessOrEqualTo(num.Uint128From64(1)))
	assert.Equal(t, false, num.MaxUint128.LessOrEqualTo(num.Uint128From64(0)))
}

func TestUint128Mul(t *testing.T) {
	bigMax64 := new(big.Int).SetInt64(math.MaxInt64)
	assert.Equal(t, num.Uint128From64(0), num.Uint128From64(0).Mul(num.Uint128From64(0)))
	assert.Equal(t, num.Uint128From64(4), num.Uint128From64(2).Mul(num.Uint128From64(2)))
	assert.Equal(t, num.Uint128From64(0), num.Uint128From64(1).Mul(num.Uint128From64(0)))
	assert.Equal(t, num.Uint128From64(1176), num.Uint128From64(12).Mul(num.Uint128From64(98)))
	assert.Equal(t, num.Uint128FromBigInt(new(big.Int).Mul(bigMax64, bigMax64)), num.Uint128From64(math.MaxInt64).Mul(num.Uint128From64(math.MaxInt64)))
}

func TestUint128Div(t *testing.T) {
	left, _ := new(big.Int).SetString("170141183460469231731687303715884105728", 10)
	result, _ := new(big.Int).SetString("17014118346046923173168730371588410", 10)
	assert.Equal(t, num.Uint128From64(0), num.Uint128From64(1).Div(num.Uint128From64(2)))
	assert.Equal(t, num.Uint128From64(3), num.Uint128From64(11).Div(num.Uint128From64(3)))
	assert.Equal(t, num.Uint128From64(4), num.Uint128From64(12).Div(num.Uint128From64(3)))
	assert.Equal(t, num.Uint128From64(1), num.Uint128From64(10).Div(num.Uint128From64(10)))
	assert.Equal(t, num.Uint128From64(1), num.Uint128FromComponents(1, 0).Div(num.Uint128FromComponents(1, 0)))
	assert.Equal(t, num.Uint128From64(2), num.Uint128FromComponents(246, 0).Div(num.Uint128FromComponents(123, 0)))
	assert.Equal(t, num.Uint128From64(2), num.Uint128FromComponents(246, 0).Div(num.Uint128FromComponents(122, 0)))
	assert.Equal(t, num.Uint128FromBigInt(result), num.Uint128FromBigInt(left).Div(num.Uint128From64(10000)))
}

func TestUint128Json(t *testing.T) {
	for i, one := range utable {
		if !one.IsUint128 {
			continue
		}
		in := num.Uint128FromStringNoCheck(one.ValueAsStr)
		data, err := json.Marshal(in)
		assert.NoError(t, err, indexFmt, i)
		var out num.Uint128
		assert.NoError(t, json.Unmarshal(data, &out), indexFmt, i)
		assert.Equal(t, in, out, indexFmt, i)
	}
}

func TestUint128Yaml(t *testing.T) {
	for i, one := range utable {
		if !one.IsUint128 {
			continue
		}
		in := num.Uint128FromStringNoCheck(one.ValueAsStr)
		data, err := yaml.Marshal(in)
		assert.NoError(t, err, indexFmt, i)
		var out num.Uint128
		assert.NoError(t, yaml.Unmarshal(data, &out), indexFmt, i)
		assert.Equal(t, in, out, indexFmt, i)
	}
}
