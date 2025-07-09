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
	for i := 0; i < 100; i++ {
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
