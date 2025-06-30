// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package num_test

import (
	"encoding/json"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"gopkg.in/yaml.v3"

	"github.com/richardwilkes/toolbox/v2/xmath/num"
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
	check.True(t, ok, indexFmt, index)
	check.Equal(t, one.ValueAsStr, b.String(), indexFmt, index)
	return b
}

func TestInt128FromInt64(t *testing.T) {
	for i, one := range table {
		if one.IsInt64 {
			check.Equal(t, one.ExpectedConversionAsStr, num.Int128From64(one.Int64).String(), indexFmt, i)
		}
	}
}

func TestInt128FromBigInt(t *testing.T) {
	for i, one := range table {
		check.Equal(t, one.ExpectedConversionAsStr, num.Int128FromBigInt(bigIntFromStr(t, one, i)).String(), indexFmt, i)
	}
}

func TestInt128AsBigInt(t *testing.T) {
	for i, one := range table {
		if one.IsInt128 {
			check.Equal(t, one.ValueAsStr, num.Int128FromBigInt(bigIntFromStr(t, one, i)).AsBigInt().String(), indexFmt, i)
		}
	}
}

func TestInt128AsInt64(t *testing.T) {
	for i, one := range table {
		if one.IsInt64 {
			check.Equal(t, one.Int64, num.Int128From64(one.Int64).AsInt64(), indexFmt, i)
		}
	}
}

func TestInt128IsInt64(t *testing.T) {
	for i, one := range table {
		if one.IsInt128 {
			check.Equal(t, one.IsInt64, num.Int128FromBigInt(bigIntFromStr(t, one, i)).IsInt64(), indexFmt, i)
		}
	}
}

func TestInt128Sign(t *testing.T) {
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
			check.Equal(t, sign, num.Int128FromBigInt(bigIntFromStr(t, one, i)).Sign(), indexFmt, i)
		}
	}
}

func TestInt128Inc(t *testing.T) {
	big1 := new(big.Int).SetInt64(1)
	for i, one := range table {
		if one.IsInt128 {
			b := bigIntFromStr(t, one, i)
			v := num.Int128FromBigInt(b)
			if v == num.MaxInt128 {
				check.Equal(t, num.MinInt128, v.Inc(), indexFmt, i)
			} else {
				b.Add(b, big1)
				check.Equal(t, b.String(), v.Inc().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestInt128Dec(t *testing.T) {
	big1 := new(big.Int).SetInt64(1)
	for i, one := range table {
		if one.IsInt128 {
			b := bigIntFromStr(t, one, i)
			v := num.Int128FromBigInt(b)
			if v == num.MinInt128 {
				check.Equal(t, num.MaxInt128, v.Dec(), indexFmt, i)
			} else {
				b.Sub(b, big1)
				check.Equal(t, b.String(), v.Dec().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestInt128Add(t *testing.T) {
	check.Equal(t, num.Int128From64(0), num.Int128From64(0).Add(num.Int128From64(0)))
	check.Equal(t, num.Int128From64(-3), num.Int128From64(-2).Add(num.Int128From64(-1)))
	check.Equal(t, num.Int128From64(0), num.Int128From64(1).Add(num.Int128From64(-1)))
	check.Equal(t, num.Int128From64(-1), num.Int128From64(-2).Add(num.Int128From64(1)))
	check.Equal(t, num.Int128From64(120), num.Int128From64(22).Add(num.Int128From64(98)))
	check.Equal(t, num.Int128FromComponents(1, 0), num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Add(num.Int128From64(1)))
	check.Equal(t, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Int128FromComponents(1, 0).Add(num.Int128From64(-1)))
	check.Equal(t, num.MinInt128, num.MaxInt128.Add(num.Int128From64(1)))
}

func TestInt128Sub(t *testing.T) {
	check.Equal(t, num.Int128From64(0), num.Int128From64(0).Sub(num.Int128From64(0)))
	check.Equal(t, num.Int128From64(-1), num.Int128From64(-2).Sub(num.Int128From64(-1)))
	check.Equal(t, num.Int128From64(-1), num.Int128From64(1).Sub(num.Int128From64(2)))
	check.Equal(t, num.Int128From64(1), num.Int128From64(-1).Sub(num.Int128From64(-2)))
	check.Equal(t, num.Int128From64(2), num.Int128From64(1).Sub(num.Int128From64(-1)))
	check.Equal(t, num.Int128From64(-2), num.Int128From64(-1).Sub(num.Int128From64(1)))
	check.Equal(t, num.Int128From64(-3), num.Int128From64(-2).Sub(num.Int128From64(1)))
	check.Equal(t, num.Int128From64(-76), num.Int128From64(22).Sub(num.Int128From64(98)))
	check.Equal(t, num.Int128FromComponents(1, 0), num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Sub(num.Int128From64(-1)))
	check.Equal(t, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Int128FromComponents(1, 0).Sub(num.Int128From64(1)))
	check.Equal(t, num.MaxInt128, num.MinInt128.Sub(num.Int128From64(1)))
	check.Equal(t, num.MinInt128, num.MaxInt128.Sub(num.Int128From64(-1)))
	check.Equal(t, num.Int128FromComponents(0x8000000000000000, 1), num.MinInt128.Sub(num.Int128From64(-1)))
}

func TestInt128Neg(t *testing.T) {
	check.Equal(t, num.Int128From64(0), num.Int128From64(0).Neg())
	check.Equal(t, num.Int128From64(1), num.Int128From64(-1).Neg())
	check.Equal(t, num.Int128From64(-1), num.Int128From64(1).Neg())
	check.Equal(t, num.Int128FromComponents(0xFFFFFFFFFFFFFFFF, 1), num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Neg())
	check.Equal(t, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Int128FromComponents(0xFFFFFFFFFFFFFFFF, 1).Neg())
	check.Equal(t, num.Int128FromComponents(0x8000000000000000, 1), num.MaxInt128.Neg())
	check.Equal(t, num.MinInt128, num.MinInt128.Neg())
	check.Equal(t, num.Int128FromComponents(0xFFFFFFFFFFFFFFFF, 0), num.Int128FromComponents(1, 0).Neg())
}

func TestInt128Abs(t *testing.T) {
	check.Equal(t, num.Int128From64(0), num.Int128From64(0).Abs())
	check.Equal(t, num.Int128From64(1), num.Int128From64(-1).Abs())
	check.Equal(t, num.Int128From64(1), num.Int128From64(1).Abs())
	check.Equal(t, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Abs())
	check.Equal(t, num.Int128FromComponents(1, 0), num.Int128FromComponents(0xFFFFFFFFFFFFFFFF, 0).Abs())
	check.Equal(t, num.MaxInt128, num.MaxInt128.Abs())
	check.Equal(t, num.MinInt128, num.MinInt128.Abs())
}

func TestInt128AbsUint128(t *testing.T) {
	check.Equal(t, num.Uint128From64(0), num.Int128From64(0).AbsUint128())
	check.Equal(t, num.Uint128From64(1), num.Int128From64(-1).AbsUint128())
	check.Equal(t, num.Uint128From64(1), num.Int128From64(1).AbsUint128())
	check.Equal(t, num.Uint128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).AbsUint128())
	check.Equal(t, num.Uint128FromComponents(1, 0), num.Int128FromComponents(0xFFFFFFFFFFFFFFFF, 0).AbsUint128())
	check.Equal(t, num.Uint128FromComponents(0x7FFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF), num.MaxInt128.AbsUint128())
	check.Equal(t, num.Uint128FromComponents(0x8000000000000000, 0), num.MinInt128.AbsUint128())
}

func TestInt128Cmp(t *testing.T) {
	check.Equal(t, 0, num.Int128From64(0).Cmp(num.Int128From64(0)))
	check.Equal(t, -1, num.Int128From64(-2).Cmp(num.Int128From64(-1)))
	check.Equal(t, -1, num.Int128From64(1).Cmp(num.Int128From64(2)))
	check.Equal(t, 1, num.Int128From64(-1).Cmp(num.Int128From64(-2)))
	check.Equal(t, 1, num.Int128From64(1).Cmp(num.Int128From64(-1)))
	check.Equal(t, -1, num.Int128From64(-1).Cmp(num.Int128From64(1)))
	check.Equal(t, -1, num.Int128From64(-2).Cmp(num.Int128From64(1)))
	check.Equal(t, -1, num.Int128From64(22).Cmp(num.Int128From64(98)))
	check.Equal(t, 1, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Cmp(num.Int128From64(-1)))
	check.Equal(t, 1, num.Int128FromComponents(1, 0).Cmp(num.Int128From64(1)))
	check.Equal(t, -1, num.MinInt128.Cmp(num.Int128From64(1)))
	check.Equal(t, 1, num.MaxInt128.Cmp(num.Int128From64(-1)))
	check.Equal(t, -1, num.MinInt128.Cmp(num.MaxInt128))
	check.Equal(t, 1, num.MaxInt128.Cmp(num.MinInt128))
	check.Equal(t, 0, num.MaxInt128.Cmp(num.MaxInt128)) //nolint:gocritic // Yes, we meant to compare the same value
	check.Equal(t, 0, num.MinInt128.Cmp(num.MinInt128)) //nolint:gocritic // Yes, we meant to compare the same value
}

func TestInt128GreaterThan(t *testing.T) {
	check.Equal(t, false, num.Int128From64(0).GreaterThan(num.Int128From64(0)))
	check.Equal(t, false, num.Int128From64(-2).GreaterThan(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128From64(1).GreaterThan(num.Int128From64(2)))
	check.Equal(t, false, num.Int128From64(-1).GreaterThan(num.Int128From64(1)))
	check.Equal(t, false, num.Int128From64(-2).GreaterThan(num.Int128From64(1)))
	check.Equal(t, false, num.Int128From64(22).GreaterThan(num.Int128From64(98)))
	check.Equal(t, false, num.MinInt128.GreaterThan(num.Int128From64(1)))
	check.Equal(t, false, num.MinInt128.GreaterThan(num.MaxInt128))
	check.Equal(t, false, num.MaxInt128.GreaterThan(num.MaxInt128))
	check.Equal(t, false, num.MinInt128.GreaterThan(num.MinInt128))
	check.Equal(t, true, num.Int128From64(-1).GreaterThan(num.Int128From64(-2)))
	check.Equal(t, true, num.Int128From64(1).GreaterThan(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).GreaterThan(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128FromComponents(1, 0).GreaterThan(num.Int128From64(1)))
	check.Equal(t, true, num.MaxInt128.GreaterThan(num.Int128From64(-1)))
	check.Equal(t, true, num.MaxInt128.GreaterThan(num.MinInt128))
}

func TestInt128GreaterOrEqualTo(t *testing.T) {
	check.Equal(t, true, num.Int128From64(0).GreaterThanOrEqual(num.Int128From64(0)))
	check.Equal(t, false, num.Int128From64(-2).GreaterThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128From64(1).GreaterThanOrEqual(num.Int128From64(2)))
	check.Equal(t, false, num.Int128From64(-1).GreaterThanOrEqual(num.Int128From64(1)))
	check.Equal(t, false, num.Int128From64(-2).GreaterThanOrEqual(num.Int128From64(1)))
	check.Equal(t, false, num.Int128From64(22).GreaterThanOrEqual(num.Int128From64(98)))
	check.Equal(t, false, num.MinInt128.GreaterThanOrEqual(num.Int128From64(1)))
	check.Equal(t, false, num.MinInt128.GreaterThanOrEqual(num.MaxInt128))
	check.Equal(t, true, num.MaxInt128.GreaterThanOrEqual(num.MaxInt128))
	check.Equal(t, true, num.MinInt128.GreaterThanOrEqual(num.MinInt128))
	check.Equal(t, true, num.Int128From64(-1).GreaterThanOrEqual(num.Int128From64(-2)))
	check.Equal(t, true, num.Int128From64(1).GreaterThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).GreaterThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128FromComponents(1, 0).GreaterThanOrEqual(num.Int128From64(1)))
	check.Equal(t, true, num.MaxInt128.GreaterThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, true, num.MaxInt128.GreaterThanOrEqual(num.MinInt128))
}

func TestInt128LessThan(t *testing.T) {
	check.Equal(t, false, num.Int128From64(0).LessThan(num.Int128From64(0)))
	check.Equal(t, true, num.Int128From64(-2).LessThan(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128From64(1).LessThan(num.Int128From64(2)))
	check.Equal(t, true, num.Int128From64(-1).LessThan(num.Int128From64(1)))
	check.Equal(t, true, num.Int128From64(-2).LessThan(num.Int128From64(1)))
	check.Equal(t, true, num.Int128From64(22).LessThan(num.Int128From64(98)))
	check.Equal(t, true, num.MinInt128.LessThan(num.Int128From64(1)))
	check.Equal(t, true, num.MinInt128.LessThan(num.MaxInt128))
	check.Equal(t, false, num.MaxInt128.LessThan(num.MaxInt128))
	check.Equal(t, false, num.MinInt128.LessThan(num.MinInt128))
	check.Equal(t, false, num.Int128From64(-1).LessThan(num.Int128From64(-2)))
	check.Equal(t, false, num.Int128From64(1).LessThan(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).LessThan(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128FromComponents(1, 0).LessThan(num.Int128From64(1)))
	check.Equal(t, false, num.MaxInt128.LessThan(num.Int128From64(-1)))
	check.Equal(t, false, num.MaxInt128.LessThan(num.MinInt128))
}

func TestInt128LessOrEqualTo(t *testing.T) {
	check.Equal(t, true, num.Int128From64(0).LessThanOrEqual(num.Int128From64(0)))
	check.Equal(t, true, num.Int128From64(-2).LessThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, true, num.Int128From64(1).LessThanOrEqual(num.Int128From64(2)))
	check.Equal(t, true, num.Int128From64(-1).LessThanOrEqual(num.Int128From64(1)))
	check.Equal(t, true, num.Int128From64(-2).LessThanOrEqual(num.Int128From64(1)))
	check.Equal(t, true, num.Int128From64(22).LessThanOrEqual(num.Int128From64(98)))
	check.Equal(t, true, num.MinInt128.LessThanOrEqual(num.Int128From64(1)))
	check.Equal(t, true, num.MinInt128.LessThanOrEqual(num.MaxInt128))
	check.Equal(t, true, num.MaxInt128.LessThanOrEqual(num.MaxInt128))
	check.Equal(t, true, num.MinInt128.LessThanOrEqual(num.MinInt128))
	check.Equal(t, false, num.Int128From64(-1).LessThanOrEqual(num.Int128From64(-2)))
	check.Equal(t, false, num.Int128From64(1).LessThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128FromComponents(0, 0xFFFFFFFFFFFFFFFF).LessThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, false, num.Int128FromComponents(1, 0).LessThanOrEqual(num.Int128From64(1)))
	check.Equal(t, false, num.MaxInt128.LessThanOrEqual(num.Int128From64(-1)))
	check.Equal(t, false, num.MaxInt128.LessThanOrEqual(num.MinInt128))
}

func TestInt128Mul(t *testing.T) {
	bigMax64 := new(big.Int).SetInt64(math.MaxInt64)
	bigMin64 := new(big.Int).SetInt64(math.MinInt64)
	check.Equal(t, num.Int128From64(0), num.Int128From64(0).Mul(num.Int128From64(0)))
	check.Equal(t, num.Int128From64(4), num.Int128From64(-2).Mul(num.Int128From64(-2)))
	check.Equal(t, num.Int128From64(-4), num.Int128From64(-2).Mul(num.Int128From64(2)))
	check.Equal(t, num.Int128From64(0), num.Int128From64(1).Mul(num.Int128From64(0)))
	check.Equal(t, num.Int128From64(1176), num.Int128From64(12).Mul(num.Int128From64(98)))
	check.Equal(t, num.Int128FromBigInt(new(big.Int).Mul(bigMax64, bigMax64)), num.Int128From64(math.MaxInt64).Mul(num.Int128From64(math.MaxInt64)))
	check.Equal(t, num.Int128FromBigInt(new(big.Int).Mul(bigMin64, bigMin64)), num.Int128From64(math.MinInt64).Mul(num.Int128From64(math.MinInt64)))
	check.Equal(t, num.Int128FromBigInt(new(big.Int).Mul(bigMin64, bigMax64)), num.Int128From64(math.MinInt64).Mul(num.Int128From64(math.MaxInt64)))
}

func TestInt128Div(t *testing.T) {
	left, _ := new(big.Int).SetString("-170141183460469231731687303715884105728", 10)
	result, _ := new(big.Int).SetString("-17014118346046923173168730371588410", 10)
	check.Equal(t, num.Int128From64(0), num.Int128From64(1).Div(num.Int128From64(2)))
	check.Equal(t, num.Int128From64(3), num.Int128From64(11).Div(num.Int128From64(3)))
	check.Equal(t, num.Int128From64(4), num.Int128From64(12).Div(num.Int128From64(3)))
	check.Equal(t, num.Int128From64(-3), num.Int128From64(11).Div(num.Int128From64(-3)))
	check.Equal(t, num.Int128From64(-4), num.Int128From64(12).Div(num.Int128From64(-3)))
	check.Equal(t, num.Int128From64(1), num.Int128From64(10).Div(num.Int128From64(10)))
	check.Equal(t, num.Int128From64(-1), num.Int128From64(10).Div(num.Int128From64(-10)))
	check.Equal(t, num.Int128From64(1), num.Int128FromComponents(1, 0).Div(num.Int128FromComponents(1, 0)))
	check.Equal(t, num.Int128From64(2), num.Int128FromComponents(246, 0).Div(num.Int128FromComponents(123, 0)))
	check.Equal(t, num.Int128From64(2), num.Int128FromComponents(246, 0).Div(num.Int128FromComponents(122, 0)))
	check.Equal(t, num.Int128FromBigInt(result), num.Int128FromBigInt(left).Div(num.Int128From64(10000)))
}

func TestInt128Json(t *testing.T) {
	for i, one := range table {
		if !one.IsInt128 {
			continue
		}
		in := num.Int128FromStringNoCheck(one.ValueAsStr)
		data, err := json.Marshal(in)
		check.NoError(t, err, indexFmt, i)
		var out num.Int128
		check.NoError(t, json.Unmarshal(data, &out), indexFmt, i)
		check.Equal(t, in, out, indexFmt, i)
	}
}

func TestInt128Yaml(t *testing.T) {
	for i, one := range table {
		if !one.IsInt128 {
			continue
		}
		in := num.Int128FromStringNoCheck(one.ValueAsStr)
		data, err := yaml.Marshal(in)
		check.NoError(t, err, indexFmt, i)
		var out num.Int128
		check.NoError(t, yaml.Unmarshal(data, &out), indexFmt, i)
		check.Equal(t, in, out, indexFmt, i)
	}
}
