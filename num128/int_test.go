// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//nolint:goconst // I'd rather have the strings inline than extracted out into a constant for the tests.
package num128_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/num128"
)

const (
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
	c.True(ok, "index %d", index)
	c.Equal(one.ValueAsStr, b.String(), "index %d", index)
	return b
}

func TestInt128FromInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt64 {
			c.Equal(one.ExpectedConversionAsStr, num128.IntFrom64(one.Int64).String(), "index %d", i)
		}
	}
}

func TestInt128FromBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		c.Equal(one.ExpectedConversionAsStr, num128.IntFromBigInt(bigIntFromStr(t, one, i)).String(), "index %d", i)
	}
}

func TestInt128AsBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt128 {
			c.Equal(one.ValueAsStr, num128.IntFromBigInt(bigIntFromStr(t, one, i)).AsBigInt().String(), "index %d", i)
		}
	}
}

func TestInt128AsInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt64 {
			c.Equal(one.Int64, num128.IntFrom64(one.Int64).AsInt64(), "index %d", i)
		}
	}
}

func TestInt128IsInt64(t *testing.T) {
	c := check.New(t)
	for i, one := range table {
		if one.IsInt128 {
			c.Equal(one.IsInt64, num128.IntFromBigInt(bigIntFromStr(t, one, i)).IsInt64(), "index %d", i)
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
			c.Equal(sign, num128.IntFromBigInt(bigIntFromStr(t, one, i)).Sign(), "index %d", i)
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
				c.Equal(num128.MinInt, v.Inc(), "index %d", i)
			} else {
				b.Add(b, big1)
				c.Equal(b.String(), v.Inc().AsBigInt().String(), "index %d", i)
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
				c.Equal(num128.MaxInt, v.Dec(), "index %d", i)
			} else {
				b.Sub(b, big1)
				c.Equal(b.String(), v.Dec().AsBigInt().String(), "index %d", i)
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
		c.NoError(err, "index %d", i)
		var out num128.Int
		c.NoError(json.Unmarshal(data, &out), "index %d", i)
		c.Equal(in, out, "index %d", i)
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
		c.NoError(err, "index %d", i)
		var out num128.Int
		c.NoError(yaml.Unmarshal(data, &out), "index %d", i)
		c.Equal(in, out, "index %d", i)
	}
}

func TestIntFromUint64(t *testing.T) {
	c := check.New(t)

	testCases := []struct { //nolint:govet // Don't care about optimal pointer bytes in tests
		input    uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{math.MaxUint64, "18446744073709551615"},
	}

	for i, tc := range testCases {
		result := num128.IntFromUint64(tc.input)
		c.Equal(tc.expected, result.String(), "index %d", i)
		c.True(result.IsUint64())
		c.Equal(tc.input, result.AsUint64())
	}
}

func TestIntFromFloat64(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.Int{}, num128.IntFromFloat64(0))
	c.Equal(num128.Int{}, num128.IntFromFloat64(math.NaN()))

	c.Equal(num128.IntFrom64(1), num128.IntFromFloat64(1.0))
	c.Equal(num128.IntFrom64(42), num128.IntFromFloat64(42.5)) // truncation

	c.Equal(num128.IntFrom64(-1), num128.IntFromFloat64(-1.0))
	c.Equal(num128.IntFrom64(-42), num128.IntFromFloat64(-42.5)) // truncation
	c.Equal(num128.IntFrom64(-5), num128.IntFromFloat64(-5.0))
	c.Equal(num128.IntFrom64(math.MinInt64), num128.IntFromFloat64(-9223372036854775808)) // -2^63

	c.Equal(num128.IntFromComponents(1, 0), num128.IntFromFloat64(18446744073709551616))        // 2^64
	c.Equal(num128.IntFromComponents(1, 0).Neg(), num128.IntFromFloat64(-18446744073709551616)) // -2^64
	c.Equal(num128.IntFromComponents(2, 0), num128.IntFromFloat64(36893488147419103232))        // 2^65
	c.Equal(num128.IntFromComponents(2, 0).Neg(), num128.IntFromFloat64(-36893488147419103232)) // -2^65
	c.Equal(num128.IntFromStringNoCheck("-1000000000000000019884624838656"), num128.IntFromFloat64(-1e30))

	c.Equal(num128.MaxInt, num128.IntFromFloat64(170141183460469231731687303715884105728))  // 2^127 clamps
	c.Equal(num128.MinInt, num128.IntFromFloat64(-170141183460469231731687303715884105728)) // -2^127 is MinInt
	c.Equal(num128.MaxInt, num128.IntFromFloat64(math.MaxFloat64))
	c.Equal(num128.MinInt, num128.IntFromFloat64(-math.MaxFloat64))
}

func TestIntFromString(t *testing.T) {
	c := check.New(t)

	cases := []struct {
		input    string
		expected num128.Int
	}{
		{"0", num128.IntFrom64(0)},
		{"1", num128.IntFrom64(1)},
		{"-1", num128.IntFrom64(-1)},
		{"42", num128.IntFrom64(42)},
		{"-42", num128.IntFrom64(-42)},
		{"9223372036854775807", num128.IntFrom64(math.MaxInt64)},  // MaxInt64
		{"-9223372036854775808", num128.IntFrom64(math.MinInt64)}, // MinInt64
		{"0x10", num128.IntFrom64(16)},                            // hex
		{"-0x10", num128.IntFrom64(-16)},                          // negative hex
		{"0b1010", num128.IntFrom64(10)},                          // binary
		{"-0b1010", num128.IntFrom64(-10)},                        // negative binary
		{"010", num128.IntFrom64(8)},                              // octal
		{"-010", num128.IntFrom64(-8)},                            // negative octal
	}

	for i, tc := range cases {
		result, err := num128.IntFromString(tc.input)
		c.NoError(err, "index %d", i)
		c.Equal(tc.expected, result, "index %d", i)
	}

	result, err := num128.IntFromString("1e2")
	c.NoError(err)
	c.Equal(num128.IntFrom64(100), result)

	result, err = num128.IntFromString("-1.5e2")
	c.NoError(err)
	c.Equal(num128.IntFrom64(-150), result)

	invalidInputs := []string{
		"",
		"abc",
		"12.5", // non-integer float without exponent
		"1.5e0.5",
		"not a number",
	}

	for i, input := range invalidInputs {
		_, err = num128.IntFromString(input)
		c.HasError(err, "index %d", i)
	}
}

func TestIntFromStringNoCheck(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(42), num128.IntFromStringNoCheck("42"))
	c.Equal(num128.IntFrom64(-42), num128.IntFromStringNoCheck("-42"))

	c.Equal(num128.Int{}, num128.IntFromStringNoCheck("invalid"))
	c.Equal(num128.Int{}, num128.IntFromStringNoCheck(""))
}

func TestIntFromComponents(t *testing.T) {
	c := check.New(t)

	testCases := []struct { //nolint:govet // Don't care about optimal pointer bytes in tests
		hi, lo   uint64
		expected string
	}{
		{0, 0, "0"},
		{0, 1, "1"},
		{math.MaxUint64, math.MaxUint64, "-1"}, // Two's complement representation of -1
		{0x7FFFFFFFFFFFFFFF, math.MaxUint64, "170141183460469231731687303715884105727"}, // MaxInt
		{0x8000000000000000, 0, "-170141183460469231731687303715884105728"},             // MinInt
	}

	for i, tc := range testCases {
		result := num128.IntFromComponents(tc.hi, tc.lo)
		high, low := result.Components()
		c.Equal(tc.hi, high, "index %d", i)
		c.Equal(tc.lo, low, "index %d", i)
		c.Equal(tc.expected, result.String(), "index %d", i)
	}
}

func TestIntFromRand(t *testing.T) {
	c := check.New(t)

	source := rand.New(rand.NewSource(42)) //nolint:gosec // Fixed seed for reproducibility

	values := make(map[string]bool)
	for range 100 {
		result := num128.IntFromRand(source)
		str := result.String()
		c.False(values[str], "Random values should be unique: %s", str)
		values[str] = true
	}
}

func TestIntIsZero(t *testing.T) {
	c := check.New(t)

	c.True(num128.Int{}.IsZero())
	c.True(num128.IntFrom64(0).IsZero())
	c.False(num128.IntFrom64(1).IsZero())
	c.False(num128.IntFrom64(-1).IsZero())
	c.False(num128.IntFromComponents(0, 1).IsZero())
	c.False(num128.IntFromComponents(1, 0).IsZero())
	c.False(num128.MaxInt.IsZero())
	c.False(num128.MinInt.IsZero())
}

func TestIntToBigIntAndAsBigInt(t *testing.T) {
	c := check.New(t)

	testCases := []num128.Int{
		num128.IntFrom64(0),
		num128.IntFrom64(1),
		num128.IntFrom64(-1),
		num128.IntFrom64(math.MaxInt64),
		num128.IntFrom64(math.MinInt64),
		num128.MaxInt,
		num128.MinInt,
	}

	for i, tc := range testCases {
		var b big.Int
		tc.ToBigInt(&b)
		c.Equal(tc.String(), b.String(), "index %d", i)

		b2 := tc.AsBigInt()
		c.Equal(tc.String(), b2.String(), "index %d", i)
	}
}

func TestIntAsBigFloat(t *testing.T) {
	c := check.New(t)

	testCases := []num128.Int{
		num128.IntFrom64(0),
		num128.IntFrom64(42),
		num128.IntFrom64(-42),
		num128.MaxInt,
		num128.MinInt,
	}

	for i, tc := range testCases {
		result := tc.AsBigFloat()
		expected := new(big.Float).SetInt(tc.AsBigInt())
		c.Equal(expected.String(), result.String(), "index %d", i)
	}
}

func TestIntAsFloat64(t *testing.T) {
	c := check.New(t)

	c.Equal(0.0, num128.IntFrom64(0).AsFloat64())

	c.Equal(1.0, num128.IntFrom64(1).AsFloat64())
	c.Equal(42.0, num128.IntFrom64(42).AsFloat64())

	c.Equal(-1.0, num128.IntFrom64(-1).AsFloat64())
	c.Equal(-42.0, num128.IntFrom64(-42).AsFloat64())

	negVal := num128.IntFromComponents(math.MaxUint64, math.MaxUint64-1) // -2
	c.Equal(-2.0, negVal.AsFloat64())

	// -2^64, where hi == MaxUint64 and lo == 0
	negVal = num128.IntFromComponents(math.MaxUint64, 0)
	c.Equal(-18446744073709551616.0, negVal.AsFloat64())

	posVal := num128.IntFromComponents(1, 0) // 2^64
	c.Equal(18446744073709551616.0, posVal.AsFloat64())

	negVal = num128.IntFromComponents(0x8000000000000001, 0) // -(2^127 - 2^64), which rounds to -2^127
	c.Equal(-170141183460469231731687303715884105728.0, negVal.AsFloat64())

	c.Equal(-170141183460469231731687303715884105728.0, num128.MinInt.AsFloat64())
}

func TestIntIsUintAndAsUint(t *testing.T) {
	c := check.New(t)

	c.True(num128.IntFrom64(0).IsUint())
	c.True(num128.IntFrom64(1).IsUint())
	c.True(num128.IntFrom64(math.MaxInt64).IsUint())

	c.False(num128.IntFrom64(-1).IsUint())
	c.False(num128.IntFrom64(math.MinInt64).IsUint())
	c.False(num128.MinInt.IsUint())

	val := num128.IntFrom64(42)
	uintVal := val.AsUint()
	c.Equal("42", uintVal.String())
}

func TestIntIsInt64AndAsInt64(t *testing.T) {
	c := check.New(t)

	c.True(num128.IntFrom64(0).IsInt64())
	c.True(num128.IntFrom64(1).IsInt64())
	c.True(num128.IntFrom64(-1).IsInt64())
	c.True(num128.IntFrom64(math.MaxInt64).IsInt64())
	c.True(num128.IntFrom64(math.MinInt64).IsInt64())

	c.False(num128.MaxInt.IsInt64())
	c.False(num128.MinInt.IsInt64())

	testCases := []int64{0, 1, -1, 42, -42, math.MaxInt64, math.MinInt64}
	for i, expected := range testCases {
		val := num128.IntFrom64(expected)
		c.Equal(expected, val.AsInt64(), "index %d", i)
	}
}

func TestIntIsUint64AndAsUint64(t *testing.T) {
	c := check.New(t)

	c.True(num128.IntFrom64(0).IsUint64())
	c.True(num128.IntFrom64(1).IsUint64())
	c.True(num128.IntFrom64(math.MaxInt64).IsUint64())

	c.False(num128.IntFrom64(-1).IsUint64())
	c.False(num128.IntFrom64(math.MinInt64).IsUint64())

	c.False(num128.IntFromComponents(1, 0).IsUint64())

	val := num128.IntFrom64(42)
	c.Equal(uint64(42), val.AsUint64())
}

func TestIntAdd64(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(5), num128.IntFrom64(2).Add64(3))
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(2).Add64(-3))

	maxInt64 := num128.IntFrom64(math.MaxInt64)
	result := maxInt64.Add64(1)
	c.Equal("9223372036854775808", result.String())

	result = num128.IntFrom64(5).Add64(-10)
	c.Equal(num128.IntFrom64(-5), result)
}

func TestIntSub64(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(2), num128.IntFrom64(5).Sub64(3))
	c.Equal(num128.IntFrom64(8), num128.IntFrom64(5).Sub64(-3))

	minInt64 := num128.IntFrom64(math.MinInt64)
	result := minInt64.Sub64(1)
	c.Equal("-9223372036854775809", result.String())
}

func TestIntCmp64(t *testing.T) {
	c := check.New(t)

	c.Equal(0, num128.IntFrom64(42).Cmp64(42))
	c.Equal(0, num128.IntFrom64(-42).Cmp64(-42))

	c.Equal(-1, num128.IntFrom64(1).Cmp64(2))
	c.Equal(-1, num128.IntFrom64(-5).Cmp64(-2))

	c.Equal(1, num128.IntFrom64(5).Cmp64(3))
	c.Equal(1, num128.IntFrom64(-2).Cmp64(-5))

	large := num128.IntFromComponents(1, 0)
	c.Equal(1, large.Cmp64(math.MaxInt64))

	largeNeg := num128.IntFromComponents(0x8000000000000000, 0) // MinInt
	c.Equal(-1, largeNeg.Cmp64(math.MinInt64))
}

func TestIntComparison64Methods(t *testing.T) {
	c := check.New(t)

	val42 := num128.IntFrom64(42)
	valNeg42 := num128.IntFrom64(-42)

	c.True(val42.GreaterThan64(41))
	c.False(val42.GreaterThan64(42))
	c.False(val42.GreaterThan64(43))
	c.True(val42.GreaterThan64(-1))
	c.False(valNeg42.GreaterThan64(-41))

	c.True(val42.GreaterThanOrEqual64(41))
	c.True(val42.GreaterThanOrEqual64(42))
	c.False(val42.GreaterThanOrEqual64(43))
	c.True(valNeg42.GreaterThanOrEqual64(-42))

	c.True(val42.Equal64(42))
	c.False(val42.Equal64(41))
	c.True(valNeg42.Equal64(-42))

	c.False(val42.LessThan64(41))
	c.False(val42.LessThan64(42))
	c.True(val42.LessThan64(43))
	c.True(valNeg42.LessThan64(-41))

	c.False(val42.LessThanOrEqual64(41))
	c.True(val42.LessThanOrEqual64(42))
	c.True(val42.LessThanOrEqual64(43))
	c.True(valNeg42.LessThanOrEqual64(-42))
}

func TestIntMul64(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(15), num128.IntFrom64(3).Mul64(5))
	c.Equal(num128.IntFrom64(-15), num128.IntFrom64(3).Mul64(-5))
	c.Equal(num128.IntFrom64(-15), num128.IntFrom64(-3).Mul64(5))
	c.Equal(num128.IntFrom64(15), num128.IntFrom64(-3).Mul64(-5))

	large := num128.IntFrom64(math.MaxInt32)
	result := large.Mul64(math.MaxInt32)
	expected := int64(math.MaxInt32) * int64(math.MaxInt32)
	c.Equal(num128.IntFrom64(expected), result)
}

func TestIntDivisionByZeroPanics(t *testing.T) {
	c := check.New(t)

	val := num128.IntFrom64(42)

	c.Panics(func() {
		val.Div(num128.Int{})
	})

	c.Panics(func() {
		val.Div64(0)
	})

	c.Panics(func() {
		val.DivMod(num128.Int{})
	})

	c.Panics(func() {
		val.DivMod64(0)
	})

	c.Panics(func() {
		val.Mod(num128.Int{})
	})

	c.Panics(func() {
		val.Mod64(0)
	})
}

func TestIntDivisionEdgeCases(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(8), num128.IntFrom64(42).Div64(5))
	c.Equal(num128.IntFrom64(2), num128.IntFrom64(42).Mod64(5))

	c.Equal(num128.IntFrom64(-8), num128.IntFrom64(-42).Div64(5))
	c.Equal(num128.IntFrom64(-2), num128.IntFrom64(-42).Mod64(5))

	c.Equal(num128.IntFrom64(-8), num128.IntFrom64(42).Div64(-5))
	c.Equal(num128.IntFrom64(2), num128.IntFrom64(42).Mod64(-5))

	c.Equal(num128.IntFrom64(8), num128.IntFrom64(-42).Div64(-5))
	c.Equal(num128.IntFrom64(-2), num128.IntFrom64(-42).Mod64(-5))

	q, r := num128.IntFrom64(-42).DivMod64(5)
	c.Equal(num128.IntFrom64(-8), q)
	c.Equal(num128.IntFrom64(-2), r)
}

func TestIntNeg(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.Int{}, num128.Int{}.Neg())

	c.Equal(num128.IntFrom64(-42), num128.IntFrom64(42).Neg())
	c.Equal(num128.IntFrom64(-1), num128.IntFrom64(1).Neg())

	c.Equal(num128.IntFrom64(42), num128.IntFrom64(-42).Neg())
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(-1).Neg())

	// MinInt remains MinInt due to overflow
	c.Equal(num128.MinInt, num128.MinInt.Neg())

	c.Equal(num128.MinInt.Inc(), num128.MaxInt.Neg())
}

func TestIntAbs(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.IntFrom64(42), num128.IntFrom64(42).Abs())
	c.Equal(num128.Int{}, num128.Int{}.Abs())

	c.Equal(num128.IntFrom64(42), num128.IntFrom64(-42).Abs())
	c.Equal(num128.IntFrom64(1), num128.IntFrom64(-1).Abs())

	c.Equal(num128.MinInt, num128.MinInt.Abs()) // Can't represent positive MinInt
}

func TestIntAbsUint(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.UintFrom64(42), num128.IntFrom64(42).AbsUint())
	c.Equal(num128.UintFrom64(0), num128.IntFrom64(0).AbsUint())

	c.Equal(num128.UintFrom64(42), num128.IntFrom64(-42).AbsUint())
	c.Equal(num128.UintFrom64(1), num128.IntFrom64(-1).AbsUint())

	c.Equal(num128.Uint(num128.MinInt), num128.MinInt.AbsUint())
}

func TestIntSign(t *testing.T) {
	c := check.New(t)

	c.Equal(0, num128.Int{}.Sign())
	c.Equal(0, num128.IntFrom64(0).Sign())

	c.Equal(1, num128.IntFrom64(1).Sign())
	c.Equal(1, num128.IntFrom64(42).Sign())
	c.Equal(1, num128.MaxInt.Sign())

	c.Equal(-1, num128.IntFrom64(-1).Sign())
	c.Equal(-1, num128.IntFrom64(-42).Sign())
	c.Equal(-1, num128.MinInt.Sign())
}

func TestIntStringFormatting(t *testing.T) {
	c := check.New(t)

	c.Equal("0", num128.IntFrom64(0).String())
	c.Equal("42", num128.IntFrom64(42).String())
	c.Equal("-42", num128.IntFrom64(-42).String())
	c.Equal("170141183460469231731687303715884105727", num128.MaxInt.String())
	c.Equal("-170141183460469231731687303715884105728", num128.MinInt.String())

	val := num128.IntFrom64(42)
	c.Equal("42", fmt.Sprintf("%d", val))
	c.Equal("2a", fmt.Sprintf("%x", val))
	c.Equal("52", fmt.Sprintf("%o", val))

	valNeg := num128.IntFrom64(-42)
	c.Equal("-42", fmt.Sprintf("%d", valNeg))
}

func TestIntScan(t *testing.T) {
	c := check.New(t)

	var val num128.Int
	n, err := fmt.Sscanf("42", "%v", &val)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.IntFrom64(42), val)

	n, err = fmt.Sscanf("-42", "%v", &val)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.IntFrom64(-42), val)

	_, err = fmt.Sscanf("invalid", "%v", &val)
	c.HasError(err)
}

func TestIntTextMarshaling(t *testing.T) {
	c := check.New(t)

	val := num128.IntFrom64(42)
	valNeg := num128.IntFrom64(-42)

	text, err := val.MarshalText()
	c.NoError(err)
	c.Equal("42", string(text))

	textNeg, err := valNeg.MarshalText()
	c.NoError(err)
	c.Equal("-42", string(textNeg))

	var val2 num128.Int
	err = val2.UnmarshalText([]byte("42"))
	c.NoError(err)
	c.Equal(val, val2)

	err = val2.UnmarshalText([]byte("-42"))
	c.NoError(err)
	c.Equal(valNeg, val2)

	err = val2.UnmarshalText([]byte("invalid"))
	c.HasError(err)
}

func TestIntJSONNumberInterface(t *testing.T) {
	c := check.New(t)

	val := num128.IntFrom64(42)

	_, err := val.Float64()
	c.HasError(err)

	i64, err := val.Int64()
	c.NoError(err)
	c.Equal(int64(42), i64)

	largeVal := num128.MaxInt
	_, err = largeVal.Int64()
	c.HasError(err)
}

func TestIntJSONMarshalingErrors(t *testing.T) {
	c := check.New(t)

	var val num128.Int

	err := val.UnmarshalJSON([]byte(`"invalid"`))
	c.HasError(err)
}

func TestIntYAMLMarshalingErrors(t *testing.T) {
	c := check.New(t)

	var val num128.Int

	err := val.UnmarshalYAML(func(v any) error {
		*v.(*string) = "invalid" //nolint:errcheck // We are simulating an error
		return nil
	})
	c.HasError(err)
}

func TestIntEqual(t *testing.T) {
	c := check.New(t)

	a := num128.IntFrom64(42)
	b := num128.IntFrom64(42)
	different := num128.IntFrom64(24)

	c.True(a.Equal(b))
	c.False(a.Equal(different))
	c.True(num128.Int{}.Equal(num128.Int{})) //nolint:gocritic // Yes, we know this is pointless, but we need to test it
}

func TestIntFromBigIntEdgeCases(t *testing.T) {
	c := check.New(t)

	hugeBig := new(big.Int)
	hugeBig.SetString("999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
	result := num128.IntFromBigInt(hugeBig)
	c.Equal(num128.MaxInt, result)

	hugeNegBig := new(big.Int)
	hugeNegBig.SetString("-999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
	result = num128.IntFromBigInt(hugeNegBig)
	c.Equal(num128.MinInt, result)

	minIntBig := new(big.Int)
	minIntBig.SetString("-170141183460469231731687303715884105728", 10) // MinInt value
	result = num128.IntFromBigInt(minIntBig)
	c.Equal(num128.MinInt, result)
}

func TestIntNegEdgeCases(t *testing.T) {
	c := check.New(t)

	c.Equal(num128.MinInt, num128.MinInt.Neg())

	c.Equal(num128.Int{}, num128.Int{}.Neg())

	negMaxInt := num128.MaxInt.Neg()
	c.Equal(-1, negMaxInt.Sign()) // Should be negative

	doubleNeg := negMaxInt.Neg()
	c.Equal(1, doubleNeg.Sign()) // Should be positive
}

func TestIntDivisionPanicCases(t *testing.T) {
	c := check.New(t)

	i := num128.IntFrom64(100)
	zeroInt := num128.Int{}

	c.Panics(func() {
		i.Div(zeroInt)
	})

	c.Panics(func() {
		i.Mod(zeroInt)
	})

	c.Panics(func() {
		i.DivMod(zeroInt)
	})

	c.Panics(func() {
		i.Div64(0)
	})

	c.Panics(func() {
		i.Mod64(0)
	})

	c.Panics(func() {
		i.DivMod64(0)
	})
}

func TestIntFloatConversionEdgeCases(t *testing.T) {
	c := check.New(t)

	nanInt := num128.IntFromFloat64(math.NaN())
	c.Equal(num128.Int{}, nanInt)

	infInt := num128.IntFromFloat64(math.Inf(1))
	c.Equal(num128.MaxInt, infInt)

	negInfInt := num128.IntFromFloat64(math.Inf(-1))
	c.Equal(num128.MinInt, negInfInt)

	veryLargeFloat := 1e40
	largeInt := num128.IntFromFloat64(veryLargeFloat)
	c.Equal(num128.MaxInt, largeInt)

	veryLargeNegFloat := -1e40
	largeNegInt := num128.IntFromFloat64(veryLargeNegFloat)
	c.Equal(num128.MinInt, largeNegInt)

	// float64(math.MaxInt64) rounds up to 2^63, which no longer fits in an int64
	maxInt64Int := num128.IntFromFloat64(float64(math.MaxInt64))
	c.Equal(num128.IntFromComponents(0, 1<<63), maxInt64Int)
	c.False(maxInt64Int.IsInt64())

	// float64(math.MinInt64) is exactly -2^63
	minInt64Int := num128.IntFromFloat64(float64(math.MinInt64))
	c.Equal(num128.IntFrom64(math.MinInt64), minInt64Int)
	c.True(minInt64Int.IsInt64())
	c.Equal(int64(math.MinInt64), minInt64Int.AsInt64())
}

func TestIntStringParsingErrorCases(t *testing.T) {
	c := check.New(t)

	invalidIntStrings := []string{
		"",
		"abc",
		"123abc",
		"--123",
		"++123",
		"123.45",
		"0x", // incomplete hex
		"0b", // incomplete binary
	}

	for _, invalid := range invalidIntStrings {
		_, err := num128.IntFromString(invalid)
		c.HasError(err)
	}
}

func TestIntStringParsingEdgeCases(t *testing.T) {
	c := check.New(t)

	signTests := []struct {
		input      string
		expected   int64
		shouldWork bool
	}{
		{"+123", 123, true},
		{"-123", -123, true},
	}

	for _, test := range signTests {
		i, err := num128.IntFromString(test.input)
		if test.shouldWork {
			c.NoError(err, "Sign parsing should work for: %s", test.input)
			c.Equal(num128.IntFrom64(test.expected), i, "Sign parsing failed for: %s", test.input)
		}
	}
}

func TestIntNegAdditionalEdgeCases(t *testing.T) {
	c := check.New(t)

	zero := num128.Int{}
	negZero := zero.Neg()
	c.Equal(zero, negZero)

	negMinInt := num128.MinInt.Neg()
	c.Equal(num128.MinInt, negMinInt)

	pos := num128.IntFrom64(42)
	neg := pos.Neg()
	c.Equal(num128.IntFrom64(-42), neg)

	negVal := num128.IntFrom64(-42)
	posVal := negVal.Neg()
	c.Equal(num128.IntFrom64(42), posVal)

	special := num128.IntFromComponents(0, 1) // positive 1
	negSpecial := special.Neg()
	expectedHi := ^uint64(0)     // all 1s
	expectedLo := ^uint64(1) + 1 // ~1 + 1 = 0xFFFFFFFFFFFFFFFE + 1 = 0xFFFFFFFFFFFFFFFF
	expected := num128.IntFromComponents(expectedHi, expectedLo)
	c.Equal(expected, negSpecial)
}

func TestIntComparison64EdgeCases(t *testing.T) {
	c := check.New(t)

	negInt := num128.IntFrom64(-42)
	c.False(negInt.GreaterThan64(1))

	posInt := num128.IntFrom64(42)
	c.True(posInt.GreaterThan64(-1))

	c.True(posInt.GreaterThanOrEqual64(42))
	c.True(negInt.GreaterThanOrEqual64(-42))

	c.False(negInt.GreaterThanOrEqual64(1))

	c.True(posInt.GreaterThanOrEqual64(-1))

	c.False(posInt.LessThan64(-1))

	c.True(negInt.LessThan64(1))

	c.True(posInt.LessThanOrEqual64(42))
	c.True(negInt.LessThanOrEqual64(-42))

	c.False(posInt.LessThanOrEqual64(-1))

	c.True(negInt.LessThanOrEqual64(1))

	c.Equal(0, posInt.Cmp64(42))
	c.Equal(0, negInt.Cmp64(-42))
	c.Equal(1, posInt.Cmp64(-1))
	c.Equal(-1, negInt.Cmp64(1))
}

func TestIntScanEdgeCases(t *testing.T) {
	c := check.New(t)

	var i num128.Int
	n, err := fmt.Sscanf("123", "%v", &i)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.IntFrom64(123), i)

	var negI num128.Int
	n, err = fmt.Sscanf("-456", "%v", &negI)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.IntFrom64(-456), negI)
}

func TestIntYAMLUnmarshalEdgeCases(t *testing.T) {
	c := check.New(t)

	var i num128.Int
	err := i.UnmarshalYAML(func(v any) error {
		*v.(*string) = "123" //nolint:errcheck // Simulate YAML unmarshaling
		return nil
	})
	c.NoError(err)
	c.Equal(num128.IntFrom64(123), i)
}
