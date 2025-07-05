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

	"gopkg.in/yaml.v3"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xmath/num"
)

const (
	maxUint64PlusOneAsStr  = "18446744073709551616"
	maxUint128AsStr        = "340282366920938463463374607431768211455"
	maxUint128PlusOneAsStr = "340282366920938463463374607431768211456"
)

var uTable = []*uInfo{
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

type uInfo struct {
	ValueAsStr              string
	ExpectedConversionAsStr string
	Uint64                  uint64
	IsUint64                bool
	IsUint128               bool
}

func init() {
	for _, d := range uTable {
		if d.IsUint64 {
			d.ValueAsStr = strconv.FormatUint(d.Uint64, 10)
		}
		if d.ExpectedConversionAsStr == "" {
			d.ExpectedConversionAsStr = d.ValueAsStr
		}
	}
}

func bigUintFromStr(t *testing.T, one *uInfo, index int) *big.Int {
	t.Helper()
	b, ok := new(big.Int).SetString(one.ValueAsStr, 10)
	c := check.New(t)
	c.True(ok, indexFmt, index)
	c.Equal(one.ValueAsStr, b.String(), indexFmt, index)
	return b
}

func TestUint128FromUint64(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint64 {
			c.Equal(one.ExpectedConversionAsStr, num.Uint128From64(one.Uint64).String(), indexFmt, i)
		}
	}
}

func TestUint128FromBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		c.Equal(one.ExpectedConversionAsStr, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).String(), indexFmt, i)
	}
}

func TestUint128AsBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint128 {
			c.Equal(one.ValueAsStr, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).AsBigInt().String(), indexFmt, i)
		}
	}
}

func TestUint128AsUint64(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint64 {
			c.Equal(one.Uint64, num.Uint128From64(one.Uint64).AsUint64(), indexFmt, i)
		}
	}
}

func TestUint128IsUint64(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint128 {
			c.Equal(one.IsUint64, num.Uint128FromBigInt(bigUintFromStr(t, one, i)).IsUint64(), indexFmt, i)
		}
	}
}

func TestUint128Inc(t *testing.T) {
	c := check.New(t)
	big1 := new(big.Int).SetInt64(1)
	for i, one := range uTable {
		if one.IsUint128 {
			b := bigUintFromStr(t, one, i)
			v := num.Uint128FromBigInt(b)
			if v == num.MaxUint128 {
				c.Equal(num.Uint128{}, v.Inc(), indexFmt, i)
			} else {
				b.Add(b, big1)
				c.Equal(b.String(), v.Inc().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestUint128Dec(t *testing.T) {
	c := check.New(t)
	big1 := new(big.Int).SetInt64(1)
	for i, one := range uTable {
		if one.IsUint128 {
			b := bigUintFromStr(t, one, i)
			v := num.Uint128FromBigInt(b)
			if v.IsZero() {
				c.Equal(num.MaxUint128, v.Dec(), indexFmt, i)
			} else {
				b.Sub(b, big1)
				c.Equal(b.String(), v.Dec().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestUint128Add(t *testing.T) {
	c := check.New(t)
	c.Equal(num.Uint128From64(0), num.Uint128From64(0).Add(num.Uint128From64(0)))
	c.Equal(num.Uint128From64(120), num.Uint128From64(22).Add(num.Uint128From64(98)))
	c.Equal(num.Uint128FromComponents(1, 0), num.Uint128FromComponents(0, 0xFFFFFFFFFFFFFFFF).Add(num.Uint128From64(1)))
	c.Equal(num.Uint128From64(0), num.MaxUint128.Add(num.Uint128From64(1)))
}

func TestUint128Sub(t *testing.T) {
	c := check.New(t)
	c.Equal(num.Uint128From64(0), num.Uint128From64(0).Sub(num.Uint128From64(0)))
	c.Equal(num.Uint128FromComponents(0, 0xFFFFFFFFFFFFFFFF), num.Uint128FromComponents(1, 0).Sub(num.Uint128From64(1)))
	c.Equal(num.MaxUint128, num.Uint128From64(0).Sub(num.Uint128From64(1)))
}

func TestUint128Cmp(t *testing.T) {
	c := check.New(t)
	c.Equal(0, num.Uint128From64(0).Cmp(num.Uint128From64(0)))
	c.Equal(-1, num.Uint128From64(1).Cmp(num.Uint128From64(2)))
	c.Equal(-1, num.Uint128From64(22).Cmp(num.Uint128From64(98)))
	c.Equal(1, num.Uint128FromComponents(1, 0).Cmp(num.Uint128From64(1)))
	c.Equal(-1, num.Uint128From64(0).Cmp(num.MaxUint128))
	c.Equal(1, num.MaxUint128.Cmp(num.Uint128From64(0)))
	c.Equal(0, num.MaxUint128.Cmp(num.MaxUint128)) //nolint:gocritic // Yes, we meant to compare the same value
}

func TestUint128GreaterThan(t *testing.T) {
	c := check.New(t)
	c.False(num.Uint128From64(0).GreaterThan(num.Uint128From64(0)))
	c.False(num.Uint128From64(1).GreaterThan(num.Uint128From64(2)))
	c.False(num.Uint128From64(22).GreaterThan(num.Uint128From64(98)))
	c.False(num.Uint128From64(0).GreaterThan(num.MaxUint128))
	c.False(num.MaxUint128.GreaterThan(num.MaxUint128))
	c.True(num.Uint128FromComponents(1, 0).GreaterThan(num.Uint128From64(1)))
	c.True(num.MaxUint128.GreaterThan(num.Uint128From64(0)))
}

func TestUint128GreaterOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num.Uint128From64(0).GreaterThanOrEqual(num.Uint128From64(0)))
	c.False(num.Uint128From64(1).GreaterThanOrEqual(num.Uint128From64(2)))
	c.False(num.Uint128From64(22).GreaterThanOrEqual(num.Uint128From64(98)))
	c.False(num.Uint128From64(0).GreaterThanOrEqual(num.Uint128From64(1)))
	c.False(num.Uint128From64(0).GreaterThanOrEqual(num.MaxUint128))
	c.True(num.MaxUint128.GreaterThanOrEqual(num.MaxUint128))
	c.True(num.Uint128FromComponents(1, 0).GreaterThanOrEqual(num.Uint128From64(1)))
	c.True(num.MaxUint128.GreaterThanOrEqual(num.Uint128From64(0)))
}

func TestUint128LessThan(t *testing.T) {
	c := check.New(t)
	c.False(num.Uint128From64(0).LessThan(num.Uint128From64(0)))
	c.True(num.Uint128From64(1).LessThan(num.Uint128From64(2)))
	c.True(num.Uint128From64(22).LessThan(num.Uint128From64(98)))
	c.True(num.Uint128From64(0).LessThan(num.Uint128From64(1)))
	c.True(num.Uint128From64(0).LessThan(num.MaxUint128))
	c.False(num.MaxUint128.LessThan(num.MaxUint128))
	c.False(num.Uint128FromComponents(1, 0).LessThan(num.Uint128From64(1)))
	c.False(num.MaxUint128.LessThan(num.Uint128From64(0)))
}

func TestUint128LessOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num.Uint128From64(0).LessThanOrEqual(num.Uint128From64(0)))
	c.True(num.Uint128From64(1).LessThanOrEqual(num.Uint128From64(2)))
	c.True(num.Uint128From64(22).LessThanOrEqual(num.Uint128From64(98)))
	c.True(num.Uint128From64(0).LessThanOrEqual(num.Uint128From64(1)))
	c.True(num.Uint128From64(0).LessThanOrEqual(num.MaxUint128))
	c.True(num.MaxUint128.LessThanOrEqual(num.MaxUint128))
	c.False(num.Uint128FromComponents(1, 0).LessThanOrEqual(num.Uint128From64(1)))
	c.False(num.MaxUint128.LessThanOrEqual(num.Uint128From64(0)))
}

func TestUint128Mul(t *testing.T) {
	c := check.New(t)
	bigMax64 := new(big.Int).SetInt64(math.MaxInt64)
	c.Equal(num.Uint128From64(0), num.Uint128From64(0).Mul(num.Uint128From64(0)))
	c.Equal(num.Uint128From64(4), num.Uint128From64(2).Mul(num.Uint128From64(2)))
	c.Equal(num.Uint128From64(0), num.Uint128From64(1).Mul(num.Uint128From64(0)))
	c.Equal(num.Uint128From64(1176), num.Uint128From64(12).Mul(num.Uint128From64(98)))
	c.Equal(num.Uint128FromBigInt(new(big.Int).Mul(bigMax64, bigMax64)), num.Uint128From64(math.MaxInt64).Mul(num.Uint128From64(math.MaxInt64)))
}

func TestUint128Div(t *testing.T) {
	c := check.New(t)
	left, _ := new(big.Int).SetString("170141183460469231731687303715884105728", 10)
	result, _ := new(big.Int).SetString("17014118346046923173168730371588410", 10)
	c.Equal(num.Uint128From64(0), num.Uint128From64(1).Div(num.Uint128From64(2)))
	c.Equal(num.Uint128From64(3), num.Uint128From64(11).Div(num.Uint128From64(3)))
	c.Equal(num.Uint128From64(4), num.Uint128From64(12).Div(num.Uint128From64(3)))
	c.Equal(num.Uint128From64(1), num.Uint128From64(10).Div(num.Uint128From64(10)))
	c.Equal(num.Uint128From64(1), num.Uint128FromComponents(1, 0).Div(num.Uint128FromComponents(1, 0)))
	c.Equal(num.Uint128From64(2), num.Uint128FromComponents(246, 0).Div(num.Uint128FromComponents(123, 0)))
	c.Equal(num.Uint128From64(2), num.Uint128FromComponents(246, 0).Div(num.Uint128FromComponents(122, 0)))
	c.Equal(num.Uint128FromBigInt(result), num.Uint128FromBigInt(left).Div(num.Uint128From64(10000)))
}

func TestUint128Json(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if !one.IsUint128 {
			continue
		}
		in := num.Uint128FromStringNoCheck(one.ValueAsStr)
		data, err := json.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num.Uint128
		c.NoError(json.Unmarshal(data, &out), indexFmt, i)
		c.Equal(in, out, indexFmt, i)
	}
}

func TestUint128Yaml(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if !one.IsUint128 {
			continue
		}
		in := num.Uint128FromStringNoCheck(one.ValueAsStr)
		data, err := yaml.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num.Uint128
		c.NoError(yaml.Unmarshal(data, &out), indexFmt, i)
		c.Equal(in, out, indexFmt, i)
	}
}
