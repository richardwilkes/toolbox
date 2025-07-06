// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num128_test

import (
	"encoding/json"
	"math"
	"math/big"
	"strconv"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/num128"
)

const (
	indexFmt               = "index %d"
	maxInt64PlusOneAsStr   = "9223372036854775808"
	minInt64MinusOneAsStr  = "-9223372036854775809"
	maxInt128AsStr         = "170141183460469231731687303715884105727"
	maxInt128PlusOneAsStr  = "170141183460469231731687303715884105728"
	minInt128AsStr         = "-170141183460469231731687303715884105728"
	minInt128MinusOneAsStr = "-170141183460469231731687303715884105729"
)

var table = []*info{
	{
		IsInt64:  true,
		IsInt128: true,
	},
	{
		Int64:    -1,
		IsInt64:  true,
		IsInt128: true,
	},
	{
		Int64:    1,
		IsInt64:  true,
		IsInt128: true,
	},
	{
		ValueAsStr: "18446744073712590000",
		IsInt128:   true,
	},
	{
		ValueAsStr: "-18446744073712590000",
		IsInt128:   true,
	},
	{
		Int64:    math.MaxInt64,
		IsInt64:  true,
		IsInt128: true,
	},
	{
		Int64:    math.MinInt64,
		IsInt64:  true,
		IsInt128: true,
	},
	{
		ValueAsStr: maxInt64PlusOneAsStr,
		IsInt64:    false,
		IsInt128:   true,
	},
	{
		ValueAsStr: minInt64MinusOneAsStr,
		IsInt64:    false,
		IsInt128:   true,
	},
	{
		ValueAsStr: maxInt128AsStr,
		IsInt64:    false,
		IsInt128:   true,
	},
	{
		ValueAsStr: minInt128AsStr,
		IsInt64:    false,
		IsInt128:   true,
	},
	{
		ValueAsStr:              maxInt128PlusOneAsStr,
		ExpectedConversionAsStr: maxInt128AsStr,
		IsInt64:                 false,
		IsInt128:                false,
	},
	{
		ValueAsStr:              minInt128MinusOneAsStr,
		ExpectedConversionAsStr: minInt128AsStr,
		IsInt64:                 false,
		IsInt128:                false,
	},
}

type info struct {
	ValueAsStr              string
	ExpectedConversionAsStr string
	Int64                   int64
	IsInt64                 bool
	IsInt128                bool
}

func init() {
	for _, d := range table {
		if d.IsInt64 {
			d.ValueAsStr = strconv.FormatInt(d.Int64, 10)
		}
		if d.ExpectedConversionAsStr == "" {
			d.ExpectedConversionAsStr = d.ValueAsStr
		}
	}
}

func bigIntFromStr(t *testing.T, one *info, index int) *big.Int {
	t.Helper()
	b, ok := new(big.Int).SetString(one.ValueAsStr, 10)
	c := check.New(t)
	c.True(ok, indexFmt, index)
	c.Equal(one.ValueAsStr, b.String(), indexFmt, index)
	return b
}

func TestInt128FromInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt64 {
			c.Equal(one.ExpectedConversionAsStr, num128.IntFrom64(one.Int64).String(), indexFmt, i)
		}
	}
}

func TestInt128FromBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.ExpectedConversionAsStr, num128.IntFromBigInt(bigIntFromStr(t, one, i)).String(), indexFmt, i)
	}
}

func TestInt128AsBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt128 {
			c.Equal(one.ValueAsStr, num128.IntFromBigInt(bigIntFromStr(t, one, i)).AsBigInt().String(), indexFmt, i)
		}
	}
}

func TestInt128AsInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt64 {
			c.Equal(one.Int64, num128.IntFrom64(one.Int64).AsInt64(), indexFmt, i)
		}
	}
}

func TestInt128IsInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt128 {
			c.Equal(one.IsInt64, num128.IntFromBigInt(bigIntFromStr(t, one, i)).IsInt64(), indexFmt, i)
		}
	}
}

func TestInt128Sign(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt128 {
			var sign int
			switch {
			case one.ValueAsStr[0] == '-':
				sign = -1
			case one.ValueAsStr == "0":
				sign = 0
			default:
				sign = 1
			}
			c.Equal(sign, num128.IntFromBigInt(bigIntFromStr(t, one, i)).Sign(), indexFmt, i)
		}
	}
}

func TestInt128Inc(t *testing.T) {
	c := check.New(t)
	big1 := new(big.Int).SetInt64(1)
	for i, one := range table {
		if one.IsInt128 {
			b := bigIntFromStr(t, one, i)
			v := num128.IntFromBigInt(b)
			if v == num128.MaxInt {
				c.Equal(num128.MinInt, v.Inc(), indexFmt, i)
			} else {
				b.Add(b, big1)
				c.Equal(b.String(), v.Inc().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestInt128Dec(t *testing.T) {
	c := check.New(t)
	big1 := new(big.Int).SetInt64(1)
	for i, one := range table {
		if one.IsInt128 {
			b := bigIntFromStr(t, one, i)
			v := num128.IntFromBigInt(b)
			if v == num128.MinInt {
				c.Equal(num128.MaxInt, v.Dec(), indexFmt, i)
			} else {
				b.Sub(b, big1)
				c.Equal(b.String(), v.Dec().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestInt128Add(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(0).Add(num128.IntFrom64(0)))
	c.Equal(num128.IntFrom64(-3), num128.IntFrom64(-2).Add(num128.IntFrom64(-1)))
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(1).Add(num128.IntFrom64(-1)))
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(-2).Add(num128.IntFrom64(1)))
	c.Equal(num128.IntFrom64(120), num128.IntFrom64(22).Add(num128.IntFrom64(98)))
	c.Equal(num128.IntFromComponents(1, 0), num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).Add(num128.IntFrom64(1)))
	c.Equal(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.IntFromComponents(1, 0).Add(num128.IntFrom64(-1)))
	c.Equal(num128.MinInt, num128.MaxInt.Add(num128.IntFrom64(1)))
}

func TestInt128Sub(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(0).Sub(num128.IntFrom64(0)))
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(-2).Sub(num128.IntFrom64(-1)))
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(1).Sub(num128.IntFrom64(2)))
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(-1).Sub(num128.IntFrom64(-2)))
	c.Equal(num128.IntFrom64(2), num128.IntFrom64(1).Sub(num128.IntFrom64(-1)))
	c.Equal(num128.IntFrom64(-2), num128.IntFrom64(-1).Sub(num128.IntFrom64(1)))
	c.Equal(num128.IntFrom64(-3), num128.IntFrom64(-2).Sub(num128.IntFrom64(1)))
	c.Equal(num128.IntFrom64(-76), num128.IntFrom64(22).Sub(num128.IntFrom64(98)))
	c.Equal(num128.IntFromComponents(1, 0), num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).Sub(num128.IntFrom64(-1)))
	c.Equal(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.IntFromComponents(1, 0).Sub(num128.IntFrom64(1)))
	c.Equal(num128.MaxInt, num128.MinInt.Sub(num128.IntFrom64(1)))
	c.Equal(num128.MinInt, num128.MaxInt.Sub(num128.IntFrom64(-1)))
	c.Equal(num128.IntFromComponents(0x8000000000000000, 1), num128.MinInt.Sub(num128.IntFrom64(-1)))
}

func TestInt128Neg(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(0).Neg())
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(-1).Neg())
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(1).Neg())
	c.Equal(num128.IntFromComponents(0xFFFFFFFFFFFFFFFF, 1), num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).Neg())
	c.Equal(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.IntFromComponents(0xFFFFFFFFFFFFFFFF, 1).Neg())
	c.Equal(num128.IntFromComponents(0x8000000000000000, 1), num128.MaxInt.Neg())
	c.Equal(num128.MinInt, num128.MinInt.Neg())
	c.Equal(num128.IntFromComponents(0xFFFFFFFFFFFFFFFF, 0), num128.IntFromComponents(1, 0).Neg())
}

func TestInt128Abs(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(0).Abs())
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(-1).Abs())
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(1).Abs())
	c.Equal(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).Abs())
	c.Equal(num128.IntFromComponents(1, 0), num128.IntFromComponents(0xFFFFFFFFFFFFFFFF, 0).Abs())
	c.Equal(num128.MaxInt, num128.MaxInt.Abs())
	c.Equal(num128.MinInt, num128.MinInt.Abs())
}

func TestInt128AbsUint128(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.UintFrom64(0), num128.IntFrom64(0).AbsUint())
	c.Equal(num128.UintFrom64(1), num128.IntFrom64(-1).AbsUint())
	c.Equal(num128.UintFrom64(1), num128.IntFrom64(1).AbsUint())
	c.Equal(num128.UintFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).AbsUint())
	c.Equal(num128.UintFromComponents(1, 0), num128.IntFromComponents(0xFFFFFFFFFFFFFFFF, 0).AbsUint())
	c.Equal(num128.UintFromComponents(0x7FFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF), num128.MaxInt.AbsUint())
	c.Equal(num128.UintFromComponents(0x8000000000000000, 0), num128.MinInt.AbsUint())
}

func TestInt128Cmp(t *testing.T) {
	c := check.New(t)
	c.Equal(0, num128.IntFrom64(0).Cmp(num128.IntFrom64(0)))
	c.Equal(-1, num128.IntFrom64(-2).Cmp(num128.IntFrom64(-1)))
	c.Equal(-1, num128.IntFrom64(1).Cmp(num128.IntFrom64(2)))
	c.Equal(1, num128.IntFrom64(-1).Cmp(num128.IntFrom64(-2)))
	c.Equal(1, num128.IntFrom64(1).Cmp(num128.IntFrom64(-1)))
	c.Equal(-1, num128.IntFrom64(-1).Cmp(num128.IntFrom64(1)))
	c.Equal(-1, num128.IntFrom64(-2).Cmp(num128.IntFrom64(1)))
	c.Equal(-1, num128.IntFrom64(22).Cmp(num128.IntFrom64(98)))
	c.Equal(1, num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).Cmp(num128.IntFrom64(-1)))
	c.Equal(1, num128.IntFromComponents(1, 0).Cmp(num128.IntFrom64(1)))
	c.Equal(-1, num128.MinInt.Cmp(num128.IntFrom64(1)))
	c.Equal(1, num128.MaxInt.Cmp(num128.IntFrom64(-1)))
	c.Equal(-1, num128.MinInt.Cmp(num128.MaxInt))
	c.Equal(1, num128.MaxInt.Cmp(num128.MinInt))
	c.Equal(0, num128.MaxInt.Cmp(num128.MaxInt)) //nolint:gocritic // Yes, we meant to compare the same value
	c.Equal(0, num128.MinInt.Cmp(num128.MinInt)) //nolint:gocritic // Yes, we meant to compare the same value
}

func TestInt128GreaterThan(t *testing.T) {
	c := check.New(t)
	c.False(num128.IntFrom64(0).GreaterThan(num128.IntFrom64(0)))
	c.False(num128.IntFrom64(-2).GreaterThan(num128.IntFrom64(-1)))
	c.False(num128.IntFrom64(1).GreaterThan(num128.IntFrom64(2)))
	c.False(num128.IntFrom64(-1).GreaterThan(num128.IntFrom64(1)))
	c.False(num128.IntFrom64(-2).GreaterThan(num128.IntFrom64(1)))
	c.False(num128.IntFrom64(22).GreaterThan(num128.IntFrom64(98)))
	c.False(num128.MinInt.GreaterThan(num128.IntFrom64(1)))
	c.False(num128.MinInt.GreaterThan(num128.MaxInt))
	c.False(num128.MaxInt.GreaterThan(num128.MaxInt))
	c.False(num128.MinInt.GreaterThan(num128.MinInt))
	c.True(num128.IntFrom64(-1).GreaterThan(num128.IntFrom64(-2)))
	c.True(num128.IntFrom64(1).GreaterThan(num128.IntFrom64(-1)))
	c.True(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).GreaterThan(num128.IntFrom64(-1)))
	c.True(num128.IntFromComponents(1, 0).GreaterThan(num128.IntFrom64(1)))
	c.True(num128.MaxInt.GreaterThan(num128.IntFrom64(-1)))
	c.True(num128.MaxInt.GreaterThan(num128.MinInt))
}

func TestInt128GreaterOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num128.IntFrom64(0).GreaterThanOrEqual(num128.IntFrom64(0)))
	c.False(num128.IntFrom64(-2).GreaterThanOrEqual(num128.IntFrom64(-1)))
	c.False(num128.IntFrom64(1).GreaterThanOrEqual(num128.IntFrom64(2)))
	c.False(num128.IntFrom64(-1).GreaterThanOrEqual(num128.IntFrom64(1)))
	c.False(num128.IntFrom64(-2).GreaterThanOrEqual(num128.IntFrom64(1)))
	c.False(num128.IntFrom64(22).GreaterThanOrEqual(num128.IntFrom64(98)))
	c.False(num128.MinInt.GreaterThanOrEqual(num128.IntFrom64(1)))
	c.False(num128.MinInt.GreaterThanOrEqual(num128.MaxInt))
	c.True(num128.MaxInt.GreaterThanOrEqual(num128.MaxInt))
	c.True(num128.MinInt.GreaterThanOrEqual(num128.MinInt))
	c.True(num128.IntFrom64(-1).GreaterThanOrEqual(num128.IntFrom64(-2)))
	c.True(num128.IntFrom64(1).GreaterThanOrEqual(num128.IntFrom64(-1)))
	c.True(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).GreaterThanOrEqual(num128.IntFrom64(-1)))
	c.True(num128.IntFromComponents(1, 0).GreaterThanOrEqual(num128.IntFrom64(1)))
	c.True(num128.MaxInt.GreaterThanOrEqual(num128.IntFrom64(-1)))
	c.True(num128.MaxInt.GreaterThanOrEqual(num128.MinInt))
}

func TestInt128LessThan(t *testing.T) {
	c := check.New(t)
	c.False(num128.IntFrom64(0).LessThan(num128.IntFrom64(0)))
	c.True(num128.IntFrom64(-2).LessThan(num128.IntFrom64(-1)))
	c.True(num128.IntFrom64(1).LessThan(num128.IntFrom64(2)))
	c.True(num128.IntFrom64(-1).LessThan(num128.IntFrom64(1)))
	c.True(num128.IntFrom64(-2).LessThan(num128.IntFrom64(1)))
	c.True(num128.IntFrom64(22).LessThan(num128.IntFrom64(98)))
	c.True(num128.MinInt.LessThan(num128.IntFrom64(1)))
	c.True(num128.MinInt.LessThan(num128.MaxInt))
	c.False(num128.MaxInt.LessThan(num128.MaxInt))
	c.False(num128.MinInt.LessThan(num128.MinInt))
	c.False(num128.IntFrom64(-1).LessThan(num128.IntFrom64(-2)))
	c.False(num128.IntFrom64(1).LessThan(num128.IntFrom64(-1)))
	c.False(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).LessThan(num128.IntFrom64(-1)))
	c.False(num128.IntFromComponents(1, 0).LessThan(num128.IntFrom64(1)))
	c.False(num128.MaxInt.LessThan(num128.IntFrom64(-1)))
	c.False(num128.MaxInt.LessThan(num128.MinInt))
}

func TestInt128LessOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num128.IntFrom64(0).LessThanOrEqual(num128.IntFrom64(0)))
	c.True(num128.IntFrom64(-2).LessThanOrEqual(num128.IntFrom64(-1)))
	c.True(num128.IntFrom64(1).LessThanOrEqual(num128.IntFrom64(2)))
	c.True(num128.IntFrom64(-1).LessThanOrEqual(num128.IntFrom64(1)))
	c.True(num128.IntFrom64(-2).LessThanOrEqual(num128.IntFrom64(1)))
	c.True(num128.IntFrom64(22).LessThanOrEqual(num128.IntFrom64(98)))
	c.True(num128.MinInt.LessThanOrEqual(num128.IntFrom64(1)))
	c.True(num128.MinInt.LessThanOrEqual(num128.MaxInt))
	c.True(num128.MaxInt.LessThanOrEqual(num128.MaxInt))
	c.True(num128.MinInt.LessThanOrEqual(num128.MinInt))
	c.False(num128.IntFrom64(-1).LessThanOrEqual(num128.IntFrom64(-2)))
	c.False(num128.IntFrom64(1).LessThanOrEqual(num128.IntFrom64(-1)))
	c.False(num128.IntFromComponents(0, 0xFFFFFFFFFFFFFFFF).LessThanOrEqual(num128.IntFrom64(-1)))
	c.False(num128.IntFromComponents(1, 0).LessThanOrEqual(num128.IntFrom64(1)))
	c.False(num128.MaxInt.LessThanOrEqual(num128.IntFrom64(-1)))
	c.False(num128.MaxInt.LessThanOrEqual(num128.MinInt))
}

func TestInt128Mul(t *testing.T) {
	c := check.New(t)
	bigMax64 := new(big.Int).SetInt64(math.MaxInt64)
	bigMin64 := new(big.Int).SetInt64(math.MinInt64)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(0).Mul(num128.IntFrom64(0)))
	c.Equal(num128.IntFrom64(4), num128.IntFrom64(-2).Mul(num128.IntFrom64(-2)))
	c.Equal(num128.IntFrom64(-4), num128.IntFrom64(-2).Mul(num128.IntFrom64(2)))
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(1).Mul(num128.IntFrom64(0)))
	c.Equal(num128.IntFrom64(1176), num128.IntFrom64(12).Mul(num128.IntFrom64(98)))
	c.Equal(num128.IntFromBigInt(new(big.Int).Mul(bigMax64, bigMax64)), num128.IntFrom64(math.MaxInt64).Mul(num128.IntFrom64(math.MaxInt64)))
	c.Equal(num128.IntFromBigInt(new(big.Int).Mul(bigMin64, bigMin64)), num128.IntFrom64(math.MinInt64).Mul(num128.IntFrom64(math.MinInt64)))
	c.Equal(num128.IntFromBigInt(new(big.Int).Mul(bigMin64, bigMax64)), num128.IntFrom64(math.MinInt64).Mul(num128.IntFrom64(math.MaxInt64)))
}

func TestInt128Div(t *testing.T) {
	left, _ := new(big.Int).SetString("-170141183460469231731687303715884105728", 10)
	result, _ := new(big.Int).SetString("-17014118346046923173168730371588410", 10)
	c := check.New(t)
	c.Equal(num128.IntFrom64(0), num128.IntFrom64(1).Div(num128.IntFrom64(2)))
	c.Equal(num128.IntFrom64(3), num128.IntFrom64(11).Div(num128.IntFrom64(3)))
	c.Equal(num128.IntFrom64(4), num128.IntFrom64(12).Div(num128.IntFrom64(3)))
	c.Equal(num128.IntFrom64(-3), num128.IntFrom64(11).Div(num128.IntFrom64(-3)))
	c.Equal(num128.IntFrom64(-4), num128.IntFrom64(12).Div(num128.IntFrom64(-3)))
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(10).Div(num128.IntFrom64(10)))
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(10).Div(num128.IntFrom64(-10)))
	c.Equal(num128.IntFrom64(1), num128.IntFromComponents(1, 0).Div(num128.IntFromComponents(1, 0)))
	c.Equal(num128.IntFrom64(2), num128.IntFromComponents(246, 0).Div(num128.IntFromComponents(123, 0)))
	c.Equal(num128.IntFrom64(2), num128.IntFromComponents(246, 0).Div(num128.IntFromComponents(122, 0)))
	c.Equal(num128.IntFromBigInt(result), num128.IntFromBigInt(left).Div(num128.IntFrom64(10000)))
}

func TestInt128Json(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if !one.IsInt128 {
			continue
		}
		in := num128.IntFromStringNoCheck(one.ValueAsStr)
		data, err := json.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num128.Int
		c.NoError(json.Unmarshal(data, &out), indexFmt, i)
		c.Equal(in, out, indexFmt, i)
	}
}

func TestInt128Yaml(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if !one.IsInt128 {
			continue
		}
		in := num128.IntFromStringNoCheck(one.ValueAsStr)
		data, err := yaml.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num128.Int
		c.NoError(yaml.Unmarshal(data, &out), indexFmt, i)
		c.Equal(in, out, indexFmt, i)
	}
}
