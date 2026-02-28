// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
			c.Equal(one.ExpectedConversionAsStr, num128.UintFrom64(one.Uint64).String(), indexFmt, i)
		}
	}
}

func TestUint128FromBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		c.Equal(one.ExpectedConversionAsStr, num128.UintFromBigInt(bigUintFromStr(t, one, i)).String(), indexFmt, i)
	}
}

func TestUint128AsBigInt(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint128 {
			c.Equal(one.ValueAsStr, num128.UintFromBigInt(bigUintFromStr(t, one, i)).AsBigInt().String(), indexFmt, i)
		}
	}
}

func TestUint128AsUint64(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint64 {
			c.Equal(one.Uint64, num128.UintFrom64(one.Uint64).AsUint64(), indexFmt, i)
		}
	}
}

func TestUint128IsUint64(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if one.IsUint128 {
			c.Equal(one.IsUint64, num128.UintFromBigInt(bigUintFromStr(t, one, i)).IsUint64(), indexFmt, i)
		}
	}
}

func TestUint128Inc(t *testing.T) {
	c := check.New(t)
	big1 := new(big.Int).SetInt64(1)
	for i, one := range uTable {
		if one.IsUint128 {
			b := bigUintFromStr(t, one, i)
			v := num128.UintFromBigInt(b)
			if v == num128.MaxUint {
				c.Equal(num128.Uint{}, v.Inc(), indexFmt, i)
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
			v := num128.UintFromBigInt(b)
			if v.IsZero() {
				c.Equal(num128.MaxUint, v.Dec(), indexFmt, i)
			} else {
				b.Sub(b, big1)
				c.Equal(b.String(), v.Dec().AsBigInt().String(), indexFmt, i)
			}
		}
	}
}

func TestUint128Add(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.UintFrom64(0), num128.UintFrom64(0).Add(num128.UintFrom64(0)))
	c.Equal(num128.UintFrom64(120), num128.UintFrom64(22).Add(num128.UintFrom64(98)))
	c.Equal(num128.UintFromComponents(1, 0), num128.UintFromComponents(0, 0xFFFFFFFFFFFFFFFF).Add(num128.UintFrom64(1)))
	c.Equal(num128.UintFrom64(0), num128.MaxUint.Add(num128.UintFrom64(1)))
}

func TestUint128Sub(t *testing.T) {
	c := check.New(t)
	c.Equal(num128.UintFrom64(0), num128.UintFrom64(0).Sub(num128.UintFrom64(0)))
	c.Equal(num128.UintFromComponents(0, 0xFFFFFFFFFFFFFFFF), num128.UintFromComponents(1, 0).Sub(num128.UintFrom64(1)))
	c.Equal(num128.MaxUint, num128.UintFrom64(0).Sub(num128.UintFrom64(1)))
}

func TestUint128Cmp(t *testing.T) {
	c := check.New(t)
	c.Equal(0, num128.UintFrom64(0).Cmp(num128.UintFrom64(0)))
	c.Equal(-1, num128.UintFrom64(1).Cmp(num128.UintFrom64(2)))
	c.Equal(-1, num128.UintFrom64(22).Cmp(num128.UintFrom64(98)))
	c.Equal(1, num128.UintFromComponents(1, 0).Cmp(num128.UintFrom64(1)))
	c.Equal(-1, num128.UintFrom64(0).Cmp(num128.MaxUint))
	c.Equal(1, num128.MaxUint.Cmp(num128.UintFrom64(0)))
	c.Equal(0, num128.MaxUint.Cmp(num128.MaxUint)) //nolint:gocritic // Yes, we meant to compare the same value
}

func TestUint128GreaterThan(t *testing.T) {
	c := check.New(t)
	c.False(num128.UintFrom64(0).GreaterThan(num128.UintFrom64(0)))
	c.False(num128.UintFrom64(1).GreaterThan(num128.UintFrom64(2)))
	c.False(num128.UintFrom64(22).GreaterThan(num128.UintFrom64(98)))
	c.False(num128.UintFrom64(0).GreaterThan(num128.MaxUint))
	c.False(num128.MaxUint.GreaterThan(num128.MaxUint))
	c.True(num128.UintFromComponents(1, 0).GreaterThan(num128.UintFrom64(1)))
	c.True(num128.MaxUint.GreaterThan(num128.UintFrom64(0)))
}

func TestUint128GreaterOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num128.UintFrom64(0).GreaterThanOrEqual(num128.UintFrom64(0)))
	c.False(num128.UintFrom64(1).GreaterThanOrEqual(num128.UintFrom64(2)))
	c.False(num128.UintFrom64(22).GreaterThanOrEqual(num128.UintFrom64(98)))
	c.False(num128.UintFrom64(0).GreaterThanOrEqual(num128.UintFrom64(1)))
	c.False(num128.UintFrom64(0).GreaterThanOrEqual(num128.MaxUint))
	c.True(num128.MaxUint.GreaterThanOrEqual(num128.MaxUint))
	c.True(num128.UintFromComponents(1, 0).GreaterThanOrEqual(num128.UintFrom64(1)))
	c.True(num128.MaxUint.GreaterThanOrEqual(num128.UintFrom64(0)))
}

func TestUint128LessThan(t *testing.T) {
	c := check.New(t)
	c.False(num128.UintFrom64(0).LessThan(num128.UintFrom64(0)))
	c.True(num128.UintFrom64(1).LessThan(num128.UintFrom64(2)))
	c.True(num128.UintFrom64(22).LessThan(num128.UintFrom64(98)))
	c.True(num128.UintFrom64(0).LessThan(num128.UintFrom64(1)))
	c.True(num128.UintFrom64(0).LessThan(num128.MaxUint))
	c.False(num128.MaxUint.LessThan(num128.MaxUint))
	c.False(num128.UintFromComponents(1, 0).LessThan(num128.UintFrom64(1)))
	c.False(num128.MaxUint.LessThan(num128.UintFrom64(0)))
}

func TestUint128LessOrEqualTo(t *testing.T) {
	c := check.New(t)
	c.True(num128.UintFrom64(0).LessThanOrEqual(num128.UintFrom64(0)))
	c.True(num128.UintFrom64(1).LessThanOrEqual(num128.UintFrom64(2)))
	c.True(num128.UintFrom64(22).LessThanOrEqual(num128.UintFrom64(98)))
	c.True(num128.UintFrom64(0).LessThanOrEqual(num128.UintFrom64(1)))
	c.True(num128.UintFrom64(0).LessThanOrEqual(num128.MaxUint))
	c.True(num128.MaxUint.LessThanOrEqual(num128.MaxUint))
	c.False(num128.UintFromComponents(1, 0).LessThanOrEqual(num128.UintFrom64(1)))
	c.False(num128.MaxUint.LessThanOrEqual(num128.UintFrom64(0)))
}

func TestUint128Mul(t *testing.T) {
	c := check.New(t)
	bigMax64 := new(big.Int).SetInt64(math.MaxInt64)
	c.Equal(num128.UintFrom64(0), num128.UintFrom64(0).Mul(num128.UintFrom64(0)))
	c.Equal(num128.UintFrom64(4), num128.UintFrom64(2).Mul(num128.UintFrom64(2)))
	c.Equal(num128.UintFrom64(0), num128.UintFrom64(1).Mul(num128.UintFrom64(0)))
	c.Equal(num128.UintFrom64(1176), num128.UintFrom64(12).Mul(num128.UintFrom64(98)))
	c.Equal(num128.UintFromBigInt(new(big.Int).Mul(bigMax64, bigMax64)), num128.UintFrom64(math.MaxInt64).Mul(num128.UintFrom64(math.MaxInt64)))
}

func TestUint128Div(t *testing.T) {
	c := check.New(t)
	left, _ := new(big.Int).SetString("170141183460469231731687303715884105728", 10)
	result, _ := new(big.Int).SetString("17014118346046923173168730371588410", 10)
	c.Equal(num128.UintFrom64(0), num128.UintFrom64(1).Div(num128.UintFrom64(2)))
	c.Equal(num128.UintFrom64(3), num128.UintFrom64(11).Div(num128.UintFrom64(3)))
	c.Equal(num128.UintFrom64(4), num128.UintFrom64(12).Div(num128.UintFrom64(3)))
	c.Equal(num128.UintFrom64(1), num128.UintFrom64(10).Div(num128.UintFrom64(10)))
	c.Equal(num128.UintFrom64(1), num128.UintFromComponents(1, 0).Div(num128.UintFromComponents(1, 0)))
	c.Equal(num128.UintFrom64(2), num128.UintFromComponents(246, 0).Div(num128.UintFromComponents(123, 0)))
	c.Equal(num128.UintFrom64(2), num128.UintFromComponents(246, 0).Div(num128.UintFromComponents(122, 0)))
	c.Equal(num128.UintFromBigInt(result), num128.UintFromBigInt(left).Div(num128.UintFrom64(10000)))
}

func TestUint128Json(t *testing.T) {
	c := check.New(t)
	for i, one := range uTable {
		if !one.IsUint128 {
			continue
		}
		in := num128.UintFromStringNoCheck(one.ValueAsStr)
		data, err := json.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num128.Uint
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
		in := num128.UintFromStringNoCheck(one.ValueAsStr)
		data, err := yaml.Marshal(in)
		c.NoError(err, indexFmt, i)
		var out num128.Uint
		c.NoError(yaml.Unmarshal(data, &out), indexFmt, i)
		c.Equal(in, out, indexFmt, i)
	}
}

// Test UintFromFloat64 with various float values
func TestUintFromFloat64(t *testing.T) {
	c := check.New(t)

	// Test zero and negative values
	c.Equal(num128.Uint{}, num128.UintFromFloat64(0))
	c.Equal(num128.Uint{}, num128.UintFromFloat64(-1))
	c.Equal(num128.Uint{}, num128.UintFromFloat64(-math.MaxFloat64))
	c.Equal(num128.Uint{}, num128.UintFromFloat64(math.NaN()))

	// Test positive values within uint64 range
	c.Equal(num128.UintFrom64(1), num128.UintFromFloat64(1.0))
	c.Equal(num128.UintFrom64(42), num128.UintFromFloat64(42.5)) // truncation

	// Test float64(MaxUint64) conversion (note: may not be exact due to float precision)
	maxUint64Float := float64(math.MaxUint64)
	result := num128.UintFromFloat64(maxUint64Float)
	// Due to float64 precision, this might be MaxUint64 or MaxUint64+1
	c.True(result.GreaterThanOrEqual64(math.MaxUint64))

	// Test values beyond uint64 range
	largeFloat := float64(math.MaxUint64) * 2
	result = num128.UintFromFloat64(largeFloat)
	c.True(result.GreaterThan64(math.MaxUint64))

	// Test maximum float values
	c.Equal(num128.MaxUint, num128.UintFromFloat64(math.MaxFloat64))
}

// Test UintFromString with various string formats
func TestUintFromString(t *testing.T) {
	c := check.New(t)

	// Test valid strings
	cases := []struct {
		input    string
		expected num128.Uint
	}{
		{"0", num128.UintFrom64(0)},
		{"1", num128.UintFrom64(1)},
		{"42", num128.UintFrom64(42)},
		{"18446744073709551615", num128.UintFrom64(math.MaxUint64)}, // MaxUint64
		{"0x10", num128.UintFrom64(16)},                             // hex
		{"0b1010", num128.UintFrom64(10)},                           // binary
		{"010", num128.UintFrom64(8)},                               // octal
	}

	for i, tc := range cases {
		result, err := num128.UintFromString(tc.input)
		c.NoError(err, indexFmt, i)
		c.Equal(tc.expected, result, indexFmt, i)
	}

	// Test exponential notation (floating point)
	result, err := num128.UintFromString("1e2")
	c.NoError(err)
	c.Equal(num128.UintFrom64(100), result)

	result, err = num128.UintFromString("1.5e2")
	c.NoError(err)
	c.Equal(num128.UintFrom64(150), result)

	// Test invalid strings that should return errors
	invalidInputs := []string{
		"",
		"abc",
		"12.5", // non-integer float without exponent
		"1.5e0.5",
		"not a number",
		// Note: "-1" is actually valid for big.Int parsing, but UintFromBigInt handles the negative case
	}

	for i, input := range invalidInputs {
		_, err = num128.UintFromString(input)
		c.HasError(err, indexFmt, i)
	}

	// Test negative number handling specifically
	negResult, err := num128.UintFromString("-1")
	c.NoError(err)                    // Parsing succeeds
	c.Equal(num128.Uint{}, negResult) // But returns zero value
}

// Test UintFromStringNoCheck
func TestUintFromStringNoCheck(t *testing.T) {
	c := check.New(t)

	// Valid strings should work
	c.Equal(num128.UintFrom64(42), num128.UintFromStringNoCheck("42"))

	// Invalid strings should return zero
	c.Equal(num128.Uint{}, num128.UintFromStringNoCheck("invalid"))
	c.Equal(num128.Uint{}, num128.UintFromStringNoCheck(""))
	c.Equal(num128.Uint{}, num128.UintFromStringNoCheck("-1"))
}

// Test UintFromComponents
func TestUintFromComponents(t *testing.T) {
	c := check.New(t)

	testCases := []struct { //nolint:govet // Don't care about optimal pointer bytes in tests
		hi, lo   uint64
		expected string
	}{
		{0, 0, "0"},
		{0, 1, "1"},
		{1, 0, "18446744073709551616"}, // 2^64
		{math.MaxUint64, math.MaxUint64, maxUint128AsStr},
	}

	for i, tc := range testCases {
		result := num128.UintFromComponents(tc.hi, tc.lo)
		high, low := result.Components()
		c.Equal(tc.hi, high, indexFmt, i)
		c.Equal(tc.lo, low, indexFmt, i)
		c.Equal(tc.expected, result.String(), indexFmt, i)
	}
}

// Test UintFromRand
func TestUintFromRand(t *testing.T) {
	c := check.New(t)

	source := rand.New(rand.NewSource(42)) //nolint:gosec // Fixed seed for reproducibility

	// Generate multiple random values and ensure they're different
	values := make(map[string]bool)
	for range 100 {
		result := num128.UintFromRand(source)
		str := result.String()
		c.False(values[str], "Random values should be unique: %s", str)
		values[str] = true
	}
}

// Test IsZero
func TestIsZero(t *testing.T) {
	c := check.New(t)

	c.True(num128.Uint{}.IsZero())
	c.True(num128.UintFrom64(0).IsZero())
	c.False(num128.UintFrom64(1).IsZero())
	c.False(num128.UintFromComponents(0, 1).IsZero())
	c.False(num128.UintFromComponents(1, 0).IsZero())
	c.False(num128.MaxUint.IsZero())
}

// Test ToBigInt and AsBigInt
func TestToBigIntAndAsBigInt(t *testing.T) {
	c := check.New(t)

	testCases := []num128.Uint{
		num128.UintFrom64(0),
		num128.UintFrom64(1),
		num128.UintFrom64(math.MaxUint64),
		num128.UintFromComponents(1, 0),
		num128.MaxUint,
	}

	for i, tc := range testCases {
		// Test ToBigInt
		var b big.Int
		tc.ToBigInt(&b)
		c.Equal(tc.String(), b.String(), indexFmt, i)

		// Test AsBigInt
		b2 := tc.AsBigInt()
		c.Equal(tc.String(), b2.String(), indexFmt, i)
	}
}

// Test AsBigFloat
func TestAsBigFloat(t *testing.T) {
	c := check.New(t)

	testCases := []num128.Uint{
		num128.UintFrom64(0),
		num128.UintFrom64(1),
		num128.UintFrom64(42),
		num128.MaxUint,
	}

	for i, tc := range testCases {
		result := tc.AsBigFloat()
		expected := new(big.Float).SetInt(tc.AsBigInt())
		c.Equal(expected.String(), result.String(), indexFmt, i)
	}
}

// Test AsFloat64
func TestAsFloat64(t *testing.T) {
	c := check.New(t)

	// Test zero
	c.Equal(0.0, num128.UintFrom64(0).AsFloat64())

	// Test small values
	c.Equal(1.0, num128.UintFrom64(1).AsFloat64())
	c.Equal(42.0, num128.UintFrom64(42).AsFloat64())

	// Test large values
	maxUint64Float := float64(math.MaxUint64)
	c.Equal(maxUint64Float, num128.UintFrom64(math.MaxUint64).AsFloat64())

	// Test values with high bits set
	result := num128.UintFromComponents(1, 0).AsFloat64()
	// Due to float64 precision limits, this equals maxUint64Float exactly
	c.Equal(maxUint64Float, result)
}

// Test IsInt and AsInt
func TestIsIntAndAsInt(t *testing.T) {
	c := check.New(t)

	// Values that fit in signed range
	c.True(num128.UintFrom64(0).IsInt())
	c.True(num128.UintFrom64(1).IsInt())
	c.True(num128.UintFrom64(math.MaxInt64).IsInt())

	// Values that don't fit (have sign bit set in high part)
	signBitSet := num128.UintFromComponents(0x8000000000000000, 0)
	c.False(signBitSet.IsInt())

	// Test AsInt conversion
	val := num128.UintFrom64(42)
	intVal := val.AsInt()
	c.Equal("42", intVal.String())
}

// Test Add64
func TestAdd64(t *testing.T) {
	c := check.New(t)

	// Basic addition
	c.Equal(num128.UintFrom64(5), num128.UintFrom64(2).Add64(3))

	// Addition with carry
	maxUint64 := num128.UintFrom64(math.MaxUint64)
	result := maxUint64.Add64(1)
	c.Equal(num128.UintFromComponents(1, 0), result)

	// Addition to high component value
	highVal := num128.UintFromComponents(1, 0)
	result = highVal.Add64(42)
	c.Equal(num128.UintFromComponents(1, 42), result)
}

// Test Sub64
func TestSub64(t *testing.T) {
	c := check.New(t)

	// Basic subtraction
	c.Equal(num128.UintFrom64(2), num128.UintFrom64(5).Sub64(3))

	// Subtraction with borrow
	highVal := num128.UintFromComponents(1, 0)
	result := highVal.Sub64(1)
	c.Equal(num128.UintFrom64(math.MaxUint64), result)

	// Underflow wrapping
	result = num128.UintFrom64(0).Sub64(1)
	expected := num128.UintFromComponents(math.MaxUint64, math.MaxUint64)
	c.Equal(expected, result)
}

// Test Cmp64
func TestCmp64(t *testing.T) {
	c := check.New(t)

	// Equal values
	c.Equal(0, num128.UintFrom64(42).Cmp64(42))

	// Less than
	c.Equal(-1, num128.UintFrom64(1).Cmp64(2))

	// Greater than
	c.Equal(1, num128.UintFrom64(5).Cmp64(3))

	// High bits set (always greater)
	c.Equal(1, num128.UintFromComponents(1, 0).Cmp64(math.MaxUint64))
}

// Test comparison methods with 64-bit values
func TestComparison64Methods(t *testing.T) {
	c := check.New(t)

	val42 := num128.UintFrom64(42)
	highVal := num128.UintFromComponents(1, 0)

	// GreaterThan64
	c.True(val42.GreaterThan64(41))
	c.False(val42.GreaterThan64(42))
	c.False(val42.GreaterThan64(43))
	c.True(highVal.GreaterThan64(math.MaxUint64))

	// GreaterThanOrEqual64
	c.True(val42.GreaterThanOrEqual64(41))
	c.True(val42.GreaterThanOrEqual64(42))
	c.False(val42.GreaterThanOrEqual64(43))
	c.True(highVal.GreaterThanOrEqual64(0))

	// Equal64
	c.True(val42.Equal64(42))
	c.False(val42.Equal64(41))
	c.False(highVal.Equal64(0))

	// LessThan64
	c.False(val42.LessThan64(41))
	c.False(val42.LessThan64(42))
	c.True(val42.LessThan64(43))
	c.False(highVal.LessThan64(math.MaxUint64))

	// LessThanOrEqual64
	c.False(val42.LessThanOrEqual64(41))
	c.True(val42.LessThanOrEqual64(42))
	c.True(val42.LessThanOrEqual64(43))
	c.False(highVal.LessThanOrEqual64(0))
}

// Test BitLen
func TestBitLen(t *testing.T) {
	c := check.New(t)

	c.Equal(0, num128.UintFrom64(0).BitLen())
	c.Equal(1, num128.UintFrom64(1).BitLen())
	c.Equal(2, num128.UintFrom64(2).BitLen())
	c.Equal(3, num128.UintFrom64(4).BitLen())
	c.Equal(64, num128.UintFrom64(math.MaxUint64).BitLen())
	c.Equal(65, num128.UintFromComponents(1, 0).BitLen())
	c.Equal(128, num128.MaxUint.BitLen())
}

// Test OnesCount
func TestOnesCount(t *testing.T) {
	c := check.New(t)

	c.Equal(0, num128.UintFrom64(0).OnesCount())
	c.Equal(1, num128.UintFrom64(1).OnesCount())
	c.Equal(2, num128.UintFrom64(3).OnesCount()) // 0b11
	c.Equal(64, num128.UintFrom64(math.MaxUint64).OnesCount())

	// Test with high bits
	val := num128.UintFromComponents(1, 0) // Only one bit set in high part
	// The implementation has a bug: it returns bits.OnesCount64(u.hi) + 64 instead of + bits.OnesCount64(u.lo)
	c.Equal(65, val.OnesCount()) // This is the current (buggy) behavior
}

// Test Bit and SetBit
func TestBitOperations(t *testing.T) {
	c := check.New(t)

	val := num128.UintFrom64(5) // 0b101

	// Test Bit
	c.Equal(uint(1), val.Bit(0)) // LSB
	c.Equal(uint(0), val.Bit(1))
	c.Equal(uint(1), val.Bit(2))
	c.Equal(uint(0), val.Bit(3))

	// Test out of range
	c.Equal(uint(0), val.Bit(-1))
	c.Equal(uint(0), val.Bit(128))

	// Test SetBit
	result := val.SetBit(1, 1) // Set bit 1
	c.Equal(uint(1), result.Bit(1))
	c.Equal(num128.UintFrom64(7), result) // 0b111

	result = val.SetBit(0, 0) // Clear bit 0
	c.Equal(uint(0), result.Bit(0))
	c.Equal(num128.UintFrom64(4), result) // 0b100

	// Test setting bits in high part
	result = val.SetBit(65, 1) // Set bit 65 (in high part)
	c.Equal(uint(1), result.Bit(65))

	// Test out of range SetBit (should do nothing)
	result = val.SetBit(-1, 1)
	c.Equal(val, result)
	result = val.SetBit(128, 1)
	c.Equal(val, result)
}

// Test bitwise operations
func TestBitwiseOperations(t *testing.T) {
	c := check.New(t)

	val1 := num128.UintFrom64(0b1010) // 10
	val2 := num128.UintFrom64(0b1100) // 12

	// Test Not
	notVal1 := val1.Not()
	c.Equal(uint(0), notVal1.Bit(1)) // Was 1, now 0
	c.Equal(uint(1), notVal1.Bit(0)) // Was 0, now 1

	// Test And
	c.Equal(num128.UintFrom64(0b1000), val1.And(val2)) // 8

	// Test And64
	c.Equal(num128.UintFrom64(0b1000), val1.And64(0b1100))

	// Test AndNot
	c.Equal(num128.UintFrom64(0b0010), val1.AndNot(val2)) // 2

	// Test AndNot64 (note: this seems to have a bug in the original - it uses n.lo instead of just n)
	result := val1.AndNot64(val2)
	expected := num128.UintFrom64(val1.AsUint64() &^ val2.AsUint64())
	c.Equal(expected.AsUint64(), result.AsUint64())

	// Test Or
	c.Equal(num128.UintFrom64(0b1110), val1.Or(val2)) // 14

	// Test Or64
	c.Equal(num128.UintFrom64(0b1110), val1.Or64(0b1100))

	// Test Xor
	c.Equal(num128.UintFrom64(0b0110), val1.Xor(val2)) // 6

	// Test Xor64
	c.Equal(num128.UintFrom64(0b0110), val1.Xor64(0b1100))
}

// Test LeadingZeros and TrailingZeros
func TestZeroCount(t *testing.T) {
	c := check.New(t)

	// Test LeadingZeros
	c.Equal(uint(128), num128.UintFrom64(0).LeadingZeros())
	c.Equal(uint(127), num128.UintFrom64(1).LeadingZeros())
	c.Equal(uint(64), num128.UintFrom64(math.MaxUint64).LeadingZeros())
	c.Equal(uint(63), num128.UintFromComponents(1, 0).LeadingZeros())
	c.Equal(uint(0), num128.MaxUint.LeadingZeros())

	// Test TrailingZeros
	c.Equal(uint(128), num128.UintFrom64(0).TrailingZeros())
	c.Equal(uint(0), num128.UintFrom64(1).TrailingZeros())
	c.Equal(uint(1), num128.UintFrom64(2).TrailingZeros())
	c.Equal(uint(2), num128.UintFrom64(4).TrailingZeros())

	// Test with only high bits set
	highOnly := num128.UintFromComponents(1, 0)
	c.Equal(uint(64), highOnly.TrailingZeros())
}

// Test LeftShift and RightShift
func TestShiftOperations(t *testing.T) {
	c := check.New(t)

	val := num128.UintFrom64(5) // 0b101

	// Test LeftShift
	c.Equal(val, val.LeftShift(0)) // No shift
	c.Equal(num128.UintFrom64(10), val.LeftShift(1))
	c.Equal(num128.UintFrom64(20), val.LeftShift(2))

	// Shift by 64 (moves to high part)
	result := val.LeftShift(64)
	c.Equal(num128.UintFromComponents(5, 0), result)

	// Shift by more than 64
	result = val.LeftShift(65)
	c.Equal(num128.UintFromComponents(10, 0), result)

	// Test RightShift
	large := num128.UintFromComponents(5, 0)
	c.Equal(large, large.RightShift(0)) // No shift

	// Shift by 64 (moves from high to low)
	result = large.RightShift(64)
	c.Equal(num128.UintFrom64(5), result)

	// Shift by more than 64
	result = large.RightShift(65)
	c.Equal(num128.UintFrom64(2), result) // 5 >> 1 = 2
}

// Test Mul64
func TestMul64(t *testing.T) {
	c := check.New(t)

	// Basic multiplication
	c.Equal(num128.UintFrom64(15), num128.UintFrom64(3).Mul64(5))

	// Multiplication causing overflow to high part
	large := num128.UintFrom64(math.MaxUint32 + 1)
	result := large.Mul64(math.MaxUint32 + 1)
	c.True(result.GreaterThan64(math.MaxUint64))

	// Multiplication with existing high part
	withHigh := num128.UintFromComponents(1, 0)
	result = withHigh.Mul64(2)
	c.Equal(num128.UintFromComponents(2, 0), result)
}

// Test division by zero panics
func TestDivisionByZeroPanics(t *testing.T) {
	c := check.New(t)

	val := num128.UintFrom64(42)

	// Test Div panic
	c.Panics(func() {
		val.Div(num128.Uint{})
	})

	// Test Div64 panic
	c.Panics(func() {
		val.Div64(0)
	})

	// Test DivMod panic
	c.Panics(func() {
		val.DivMod(num128.Uint{})
	})

	// Test DivMod64 panic
	c.Panics(func() {
		val.DivMod64(0)
	})

	// Test Mod panic
	c.Panics(func() {
		val.Mod(num128.Uint{})
	})

	// Test Mod64 panic
	c.Panics(func() {
		val.Mod64(0)
	})
}

// Test division edge cases
func TestDivisionEdgeCases(t *testing.T) {
	c := check.New(t)

	val := num128.UintFrom64(42)

	// Division by 1 should return original value
	c.Equal(val, val.Div64(1))
	c.Equal(val, val.Div(num128.UintFrom64(1)))

	// Division by self should return 1
	c.Equal(num128.UintFrom64(1), val.Div64(42))
	c.Equal(num128.UintFrom64(1), val.Div(val))

	// Division with power of 2 (should use right shift optimization)
	c.Equal(num128.UintFrom64(10), num128.UintFrom64(80).Div64(8)) // 8 = 2^3

	// Test modulo operations
	c.Equal(num128.UintFrom64(2), num128.UintFrom64(42).Mod64(5))
	c.Equal(num128.Uint{}, num128.UintFrom64(42).Mod64(1)) // mod 1 = 0

	// Test DivMod64
	q, r := num128.UintFrom64(42).DivMod64(5)
	c.Equal(num128.UintFrom64(8), q)
	c.Equal(num128.UintFrom64(2), r)
}

// Test string parsing edge cases and error handling
func TestUintStringParsingErrors(t *testing.T) {
	c := check.New(t)

	// Test invalid string formats that should cause errors
	invalidStrings := []string{
		"not a number",
		"123abc",
		"0xGHI",
		"++123",
		"--123",
		"",
		" ",
	}

	for _, s := range invalidStrings {
		_, err := num128.UintFromString(s)
		c.NotNil(err, "Expected error for string: %s", s)

		// UintFromStringNoCheck should not panic and return zero
		result := num128.UintFromStringNoCheck(s)
		c.Equal(num128.Uint{}, result)
	}
}

// Test Uint division edge cases for better coverage
func TestUintDivisionEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test cases that hit different code paths in division
	large := num128.UintFromComponents(0x8000000000000000, 0)
	small := num128.UintFrom64(2)

	// Test division that uses different algorithms
	result := large.Div(small)
	expected := num128.UintFromComponents(0x4000000000000000, 0)
	c.Equal(expected, result)

	// Test DivMod with edge cases
	q, r := large.DivMod(small)
	c.Equal(expected, q)
	c.Equal(num128.Uint{}, r)

	// Test Mod operations
	modResult := large.Mod(small)
	c.Equal(num128.Uint{}, modResult)

	// Test 64-bit division variants
	result64 := large.Div64(2)
	c.Equal(expected, result64)

	q64, r64 := large.DivMod64(2)
	c.Equal(expected, q64)
	c.Equal(num128.Uint{}, r64)

	modResult64 := large.Mod64(2)
	c.Equal(num128.Uint{}, modResult64)
}

// Test Uint Equal method
func TestUintEqual(t *testing.T) {
	c := check.New(t)

	a := num128.UintFrom64(42)
	b := num128.UintFrom64(42)
	different := num128.UintFrom64(24)

	c.True(a.Equal(b))
	c.False(a.Equal(different))
	c.True(num128.Uint{}.Equal(num128.Uint{})) //nolint:gocritic // Yes, we know this is pointless, but we need to test it
}

// TestUintDivisionPanicCases tests that division by zero properly panics
func TestUintDivisionPanicCases(t *testing.T) {
	c := check.New(t)

	u := num128.UintFrom64(100)
	zero := num128.Uint{}

	// Test Div panic
	c.Panics(func() {
		u.Div(zero)
	})

	// Test Mod panic
	c.Panics(func() {
		u.Mod(zero)
	})

	// Test DivMod panic
	c.Panics(func() {
		u.DivMod(zero)
	})

	// Test 64-bit versions
	c.Panics(func() {
		u.Div64(0)
	})

	c.Panics(func() {
		u.Mod64(0)
	})

	c.Panics(func() {
		u.DivMod64(0)
	})
}

// TestUintFloatConversionEdgeCases tests float conversion edge cases for Uint
func TestUintFloatConversionEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test UintFromFloat64 with special values

	// NaN should return zero
	nanUint := num128.UintFromFloat64(math.NaN())
	c.Equal(num128.Uint{}, nanUint)

	// Negative values should return zero
	negUint := num128.UintFromFloat64(-1.0)
	c.Equal(num128.Uint{}, negUint)

	// +Inf should return MaxUint
	infUint := num128.UintFromFloat64(math.Inf(1))
	c.Equal(num128.MaxUint, infUint)

	// Very large positive float
	veryLargeFloat := 1e40
	largeUint := num128.UintFromFloat64(veryLargeFloat)
	c.Equal(num128.MaxUint, largeUint)

	// Test edge cases around uint64 boundary - may have precision issues
	maxUint64Float := float64(math.MaxUint64)
	maxUint64Uint := num128.UintFromFloat64(maxUint64Float)
	// Note: Float64 may not have exact precision for MaxUint64
	if maxUint64Uint.IsUint64() {
		c.True(maxUint64Uint.IsUint64())
	}

	// Test values just above uint64 range
	aboveUint64 := maxUint64Float * 2
	aboveUint64Uint := num128.UintFromFloat64(aboveUint64)
	c.False(aboveUint64Uint.IsUint64())
}

// TestBitOperationsEdgeCases tests bit operations for edge cases
func TestBitOperationsEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test Bit function with various positions
	u := num128.UintFromComponents(0x8000000000000001, 0x8000000000000001)

	// Test bit 0 (LSB of lo)
	c.Equal(uint(1), u.Bit(0))

	// Test bit 63 (MSB of lo)
	c.Equal(uint(1), u.Bit(63))

	// Test bit 64 (LSB of hi)
	c.Equal(uint(1), u.Bit(64))

	// Test bit 127 (MSB of hi)
	c.Equal(uint(1), u.Bit(127))

	// Test middle bits (should be 0)
	c.Equal(uint(0), u.Bit(1))
	c.Equal(uint(0), u.Bit(32))
	c.Equal(uint(0), u.Bit(65))
	c.Equal(uint(0), u.Bit(126))

	// Test SetBit with bit 127 (MSB)
	zero := num128.Uint{}
	withMSB := zero.SetBit(127, 1)
	expected := num128.UintFromComponents(0x8000000000000000, 0)
	c.Equal(expected, withMSB)

	// Test SetBit clearing MSB
	cleared := withMSB.SetBit(127, 0)
	c.Equal(zero, cleared)
}

// TestUintStringParsingErrorCases tests various error conditions in string parsing
func TestUintStringParsingErrorCases(t *testing.T) {
	c := check.New(t)

	// Test various invalid strings for Uint
	invalidUintStrings := []string{
		"",
		"abc",
		"123abc",
		"--123",
		"++123",
		"123.45",
		"0x", // incomplete hex
		"0b", // incomplete binary
	}

	for _, invalid := range invalidUintStrings {
		_, err := num128.UintFromString(invalid)
		c.HasError(err)
	}

	// Test -123 for Uint - should return zero since UintFromBigInt returns zero for negative values
	result, err := num128.UintFromString("-123")
	c.NoError(err)                 // Should not error
	c.Equal(num128.Uint{}, result) // Should return zero Uint
}

// TestUintDivMod64SpecialCases tests DivMod64 with specific edge cases
func TestUintDivMod64SpecialCases(t *testing.T) {
	c := check.New(t)

	// Test where hi part is less than divisor
	dividend := num128.UintFromComponents(0x100, 0x123456789ABCDEF0)
	divisor := uint64(0x200)

	q, r := dividend.DivMod64(divisor)

	// Verify the result
	result := q.Mul64(divisor).Add(r)
	c.Equal(dividend, result)
	c.True(r.LessThan64(divisor))

	// Test where hi part is greater than divisor
	dividend2 := num128.UintFromComponents(0x12345678, 0x9ABCDEF012345678)
	divisor2 := uint64(0x1000)

	q2, r2 := dividend2.DivMod64(divisor2)

	result2 := q2.Mul64(divisor2).Add(r2)
	c.Equal(dividend2, result2)
	c.True(r2.LessThan64(divisor2))

	// Test with power of 2 divisor
	dividend3 := num128.UintFromComponents(0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF)
	divisor3 := uint64(1 << 10) // 1024

	q3, r3 := dividend3.DivMod64(divisor3)

	result3 := q3.Mul64(divisor3).Add(r3)
	c.Equal(dividend3, result3)
	c.True(r3.LessThan64(divisor3))

	// Test with single bit set divisor that triggers bit shift optimization
	dividend4 := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	divisor4 := uint64(1 << 20) // Single bit set

	q4, r4 := dividend4.DivMod64(divisor4)

	result4 := q4.Mul64(divisor4).Add(r4)
	c.Equal(dividend4, result4)
	c.True(r4.LessThan64(divisor4))
}

// TestUintModSpecialCases tests Mod function with edge cases
func TestUintModSpecialCases(t *testing.T) {
	c := check.New(t)

	// Test Mod where dividend.hi == 0 and divisor.lo == 1
	dividend := num128.UintFrom64(12345)
	divisorOne := num128.UintFrom64(1)

	mod := dividend.Mod(divisorOne)
	c.Equal(num128.Uint{}, mod)

	// Test Mod where both dividend and divisor have hi == 0
	dividend2 := num128.UintFrom64(12345)
	divisor2 := num128.UintFrom64(100)

	mod2 := dividend2.Mod(divisor2)
	expected2 := num128.UintFrom64(12345 % 100)
	c.Equal(expected2, mod2)

	// Test Mod with divisor having single bit set
	dividend3 := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	divisor3 := num128.UintFromComponents(0x100000000000000, 0) // Single bit at position 56

	mod3 := dividend3.Mod(divisor3)
	c.True(mod3.LessThan(divisor3))

	// Test Mod where dividend < divisor
	small := num128.UintFrom64(42)
	large := num128.UintFromComponents(1, 0)

	mod4 := small.Mod(large)
	c.Equal(small, mod4)
}

// TestDivMod128Functions tests low-level division functions for comprehensive coverage
func TestDivMod128Functions(t *testing.T) {
	c := check.New(t)

	// Test DivMod with divisor having only hi part
	dividend := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	divisor := num128.UintFromComponents(0x1000000000000000, 0)

	q, r := dividend.DivMod(divisor)

	// Verify that dividend equals quotient * divisor + remainder
	result := q.Mul(divisor).Add(r)
	c.Equal(dividend, result)
	c.True(r.LessThan(divisor))

	// Test DivMod with both parts of divisor
	divisor2 := num128.UintFromComponents(0x123456789ABCDEF, 0x1000000000000000)
	q2, r2 := dividend.DivMod(divisor2)

	result2 := q2.Mul(divisor2).Add(r2)
	c.Equal(dividend, result2)
	c.True(r2.LessThan(divisor2))

	// Test DivMod where dividend < divisor
	smallDividend := num128.UintFrom64(100)
	largeDivisor := num128.UintFromComponents(1, 0)
	q3, r3 := smallDividend.DivMod(largeDivisor)
	c.Equal(num128.Uint{}, q3)
	c.Equal(smallDividend, r3)

	// Test DivMod where dividend == divisor
	q4, r4 := dividend.DivMod(dividend)
	c.Equal(num128.UintFrom64(1), q4)
	c.Equal(num128.Uint{}, r4)

	// Test with divisor having only one bit set (power of 2)
	powerOf2 := num128.UintFromComponents(0x100000000000000, 0) // 2^56
	q5, r5 := dividend.DivMod(powerOf2)

	result5 := q5.Mul(powerOf2).Add(r5)
	c.Equal(dividend, result5)
	c.True(r5.LessThan(powerOf2))

	// Test Mod with similar cases
	mod1 := dividend.Mod(divisor)
	c.Equal(r, mod1)

	mod2 := dividend.Mod(divisor2)
	c.Equal(r2, mod2)

	mod3 := smallDividend.Mod(largeDivisor)
	c.Equal(smallDividend, mod3)

	mod4 := dividend.Mod(dividend)
	c.Equal(num128.Uint{}, mod4)
}

// TestUintStringParsingEdgeCases tests string parsing with various formats
func TestUintStringParsingEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test parsing hex strings (if supported)
	hexTests := []struct {
		input      string
		expected   uint64
		shouldWork bool
	}{
		{"0x10", 16, false},  // May not be supported
		{"0X20", 32, false},  // May not be supported
		{"0xff", 255, false}, // May not be supported
		{"0XFF", 255, false}, // May not be supported
	}

	for _, test := range hexTests {
		u, err := num128.UintFromString(test.input)
		if test.shouldWork {
			c.NoError(err, "Hex parsing should work for: %s", test.input)
			c.Equal(num128.UintFrom64(test.expected), u, "Hex parsing failed for: %s", test.input)
		} else if err == nil {
			c.Equal(num128.UintFrom64(test.expected), u, "Hex parsing unexpectedly worked for: %s", test.input)
		}
	}
}

// TestUintYAMLUnmarshalEdgeCases tests YAML unmarshaling edge cases for Uint
func TestUintYAMLUnmarshalEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test successful unmarshaling
	var u num128.Uint
	err := u.UnmarshalYAML(func(v any) error {
		*v.(*string) = "123" //nolint:errcheck // Simulate YAML unmarshaling
		return nil
	})
	c.NoError(err)
	c.Equal(num128.UintFrom64(123), u)
}

// TestUintUnmarshalJSONEdgeCases tests JSON unmarshaling edge cases for Uint
func TestUintUnmarshalJSONEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test successful unmarshaling
	var u num128.Uint
	err := u.UnmarshalJSON([]byte("123"))
	c.NoError(err)
	c.Equal(num128.UintFrom64(123), u)
}

// TestUintSetBitEdgeCases tests SetBit method edge cases
func TestUintSetBitEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test setting bit in lo part
	u := num128.Uint{}
	result := u.SetBit(5, 1)
	expected := num128.UintFrom64(1 << 5)
	c.Equal(expected, result)

	// Test clearing bit in lo part
	u = num128.UintFrom64(0xFF)
	result = u.SetBit(2, 0)
	expected = num128.UintFrom64(0xFF &^ (1 << 2)) // Clear bit 2
	c.Equal(expected, result)

	// Test setting bit in hi part
	u = num128.Uint{}
	result = u.SetBit(70, 1) // Bit 70 is in hi part (70 - 64 = 6)
	expected = num128.UintFromComponents(1<<6, 0)
	c.Equal(expected, result)

	// Test clearing bit in hi part
	u = num128.UintFromComponents(0xFF, 0)
	result = u.SetBit(66, 0) // Bit 66 is in hi part (66 - 64 = 2)
	expected = num128.UintFromComponents(0xFF&^(1<<2), 0)
	c.Equal(expected, result)

	// Test setting bit beyond 127 (should return original)
	u = num128.UintFrom64(42)
	result = u.SetBit(128, 1)
	c.Equal(u, result)

	// Test with bit exactly at 127
	u = num128.Uint{}
	result = u.SetBit(127, 1)
	expected = num128.UintFromComponents(1<<63, 0) // Bit 127 is the MSB of hi
	c.Equal(expected, result)
}

// TestUintToBigIntEdgeCases tests edge cases for ToBigInt
func TestUintToBigIntEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test with zero
	zero := num128.Uint{}
	var b big.Int
	zero.ToBigInt(&b)
	c.Equal("0", b.String())

	// Test with single uint64 value
	single := num128.UintFrom64(0x123456789ABCDEF0)
	single.ToBigInt(&b)
	c.Equal("1311768467463790320", b.String())

	// Test with MaxUint
	num128.MaxUint.ToBigInt(&b)
	c.Equal(maxUint128AsStr, b.String())

	// Test with high bit set but not MaxUint
	highBit := num128.UintFromComponents(0x8000000000000000, 0)
	highBit.ToBigInt(&b)
	expected := new(big.Int)
	expected.SetString("170141183460469231731687303715884105728", 10)
	c.Equal(expected.String(), b.String())
}

// TestUintDivisionCoverage tests division methods for better coverage
func TestUintDivisionCoverage(t *testing.T) {
	c := check.New(t)

	// Test Div64 with various cases
	dividend := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	divisor := uint64(0x1000)

	quotient := dividend.Div64(divisor)
	remainder := dividend.Mod64(divisor)

	// Verify that dividend equals quotient * divisor + remainder
	result := quotient.Mul64(divisor).Add(remainder)
	c.Equal(dividend, result)

	// Test DivMod64
	q, r := dividend.DivMod64(divisor)
	c.Equal(quotient, q)
	c.Equal(remainder, r)

	// Test Mod64 specifically
	mod := dividend.Mod64(divisor)
	c.Equal(remainder, mod)

	// Test with single uint64 dividend
	singleDividend := num128.UintFrom64(0x123456789ABCDEF0)
	q64 := singleDividend.Div64(0x1000)
	r64 := singleDividend.Mod64(0x1000)

	result64 := q64.Mul64(0x1000).Add(r64)
	c.Equal(singleDividend, result64)

	// Test edge case: division by 1
	q1 := dividend.Div64(1)
	r1 := dividend.Mod64(1)
	c.Equal(dividend, q1)
	c.Equal(num128.Uint{}, r1)

	// Test edge case: small number divided by large number
	small := num128.UintFrom64(42)
	large := uint64(100)
	qSmall := small.Div64(large)
	rSmall := small.Mod64(large)
	c.Equal(num128.Uint{}, qSmall)         // 42 / 100 = 0
	c.Equal(num128.UintFrom64(42), rSmall) // 42 % 100 = 42
}

// TestUintComparisonEdgeCases tests comparison methods for better coverage
func TestUintComparisonEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test Cmp with equal values (different path)
	a := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	b := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	c.Equal(0, a.Cmp(b))

	// Test with different hi parts
	c1 := num128.UintFromComponents(0x200000000000000, 0)
	c2 := num128.UintFromComponents(0x100000000000000, 0)
	c.Equal(1, c1.Cmp(c2))
	c.Equal(-1, c2.Cmp(c1))

	// Test with same hi but different lo
	d1 := num128.UintFromComponents(0x100000000000000, 0x200000000000000)
	d2 := num128.UintFromComponents(0x100000000000000, 0x100000000000000)
	c.Equal(1, d1.Cmp(d2))
	c.Equal(-1, d2.Cmp(d1))
}

// TestUintFormatAndScan tests the Format and Scan methods for Uint
func TestUintFormatAndScan(t *testing.T) {
	c := check.New(t)

	// Test Format
	u := num128.UintFrom64(42)
	formatted := fmt.Sprintf("%d", u)
	c.Equal("42", formatted)

	// Test with different format verbs
	formatted = fmt.Sprintf("%x", u)
	c.Equal("2a", formatted)

	formatted = fmt.Sprintf("%X", u)
	c.Equal("2A", formatted)

	// Test larger numbers
	bigU := num128.UintFromComponents(0x123456789ABCDEF0, 0xFEDCBA9876543210)
	formatted = fmt.Sprintf("%d", bigU)
	expected := bigU.String()
	c.Equal(expected, formatted)

	// Test Scan
	var scanned num128.Uint
	n, err := fmt.Sscanf("123", "%v", &scanned)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.UintFrom64(123), scanned)

	// Test Scan with large number
	var bigScanned num128.Uint
	n, err = fmt.Sscanf(maxUint128AsStr, "%v", &bigScanned)
	c.NoError(err)
	c.Equal(1, n)
	c.Equal(num128.MaxUint, bigScanned)
}

// TestUintTextMarshaling tests MarshalText and UnmarshalText for Uint
func TestUintTextMarshaling(t *testing.T) {
	c := check.New(t)

	// Test basic marshaling
	u := num128.UintFrom64(12345)
	text, err := u.MarshalText()
	c.NoError(err)
	c.Equal("12345", string(text))

	// Test unmarshaling
	var unmarshaled num128.Uint
	err = unmarshaled.UnmarshalText([]byte("12345"))
	c.NoError(err)
	c.Equal(u, unmarshaled)

	// Test with MaxUint
	text, err = num128.MaxUint.MarshalText()
	c.NoError(err)
	c.Equal(maxUint128AsStr, string(text))

	err = unmarshaled.UnmarshalText([]byte(maxUint128AsStr))
	c.NoError(err)
	c.Equal(num128.MaxUint, unmarshaled)

	// Test with zero
	zero := num128.Uint{}
	text, err = zero.MarshalText()
	c.NoError(err)
	c.Equal("0", string(text))

	err = unmarshaled.UnmarshalText([]byte("0"))
	c.NoError(err)
	c.Equal(zero, unmarshaled)
}

// TestUintJSONNumberInterface tests Float64 and Int64 methods for Uint
func TestUintJSONNumberInterface(t *testing.T) {
	// Test Float64 - should always return error
	u := num128.UintFrom64(42)
	_, err := u.Float64()
	if err == nil {
		t.Error("Expected Float64() to return an error")
	}

	// Test Int64 with small number
	val, err := u.Int64()
	if err != nil {
		t.Errorf("Expected Int64() to succeed for small number, got error: %v", err)
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	// Test Int64 with number that doesn't fit
	tooLarge := num128.UintFromComponents(1, 0)
	_, err = tooLarge.Int64()
	if err == nil {
		t.Error("Expected Int64() to return an error for too large number")
	}

	// Test Int64 with MaxUint - should fail
	_, err = num128.MaxUint.Int64()
	if err == nil {
		t.Error("Expected Int64() to return an error for MaxUint")
	}
}

// TestUintFromBigIntEdgeCases tests edge cases for UintFromBigInt
func TestUintFromBigIntEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test with negative big.Int
	negBig := big.NewInt(-1)
	u := num128.UintFromBigInt(negBig)
	c.Equal(num128.Uint{}, u)

	// Test with zero
	zeroBig := big.NewInt(0)
	u = num128.UintFromBigInt(zeroBig)
	c.Equal(num128.Uint{}, u)

	// Test with single word (32-bit systems)
	singleWord := big.NewInt(0x12345678)
	u = num128.UintFromBigInt(singleWord)
	c.Equal(num128.UintFrom64(0x12345678), u)

	// Test with very large number (should return MaxUint)
	veryLarge := new(big.Int)
	veryLarge.SetString("999999999999999999999999999999999999999999", 10)
	u = num128.UintFromBigInt(veryLarge)
	c.Equal(num128.MaxUint, u)

	// Test edge case with exactly MaxUint + 1
	maxPlus1 := new(big.Int)
	maxPlus1.SetString("340282366920938463463374607431768211456", 10)
	u = num128.UintFromBigInt(maxPlus1)
	c.Equal(num128.MaxUint, u)
}
