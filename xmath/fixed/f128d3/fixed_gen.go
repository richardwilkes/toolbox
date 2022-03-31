// Code created from "fixed128.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f128d3

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed/internal"
	"github.com/richardwilkes/toolbox/xmath/num"
	"gopkg.in/yaml.v3"
)

var (
	multiplierBigInt   = new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil)
	multiplierBigFloat = new(big.Float).SetPrec(128).SetInt(multiplierBigInt)
	multiplier         = num.Int128FromBigInt(multiplierBigInt)
	// Max holds the maximum value.
	Max = Int{data: num.MaxInt128}
	// Min holds the minimum value.
	Min = Int{data: num.MinInt128}
)

// Some commonly used values.
var (
	One     = FromInt(1)
	Half    = FromStringForced("0.5")
	NegHalf = Half.Neg()
)

// Int holds a fixed-point value that contains up to 3 decimal places. Values are truncated, not rounded.
type Int struct {
	data num.Int128
}

// FromFloat64 creates a new value from a float64.
func FromFloat64(value float64) Int {
	f, _ := FromString(new(big.Float).SetPrec(128).SetFloat64(value).Text('f', 4)) //nolint:errcheck // Failure means 0
	return f
}

// FromFloat32 creates a new value from a float32.
func FromFloat32(value float32) Int {
	return FromFloat64(float64(value))
}

// FromInt64 creates a new value from an int64.
func FromInt64(value int64) Int {
	return Int{data: num.Int128From64(value).Mul(multiplier)}
}

// FromInt creates a new value from an int.
func FromInt(value int) Int {
	return FromInt64(int64(value))
}

// FromString creates a new value from a string.
func FromString(str string) (Int, error) {
	if str == "" {
		return Int{}, errs.New("empty string is not valid")
	}
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return Int{}, err
		}
		return FromFloat64(f), nil
	}
	parts := strings.SplitN(str, ".", 2)
	var neg bool
	value := new(big.Int)
	fraction := new(big.Int)
	switch parts[0] {
	case "":
	case "-", "-0":
		neg = true
	default:
		if _, ok := value.SetString(parts[0], 10); !ok {
			return Int{}, errs.Newf("invalid value: %s", str)
		}
		if value.Sign() < 0 {
			neg = true
			value.Neg(value)
		}
		value.Mul(value, multiplierBigInt)
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 3+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 3+1 {
			frac = frac[:3+1]
		}
		if _, ok := fraction.SetString(frac, 10); !ok {
			return Int{}, errs.Newf("invalid value: %s", str)
		}
		value.Add(value, fraction).Sub(value, multiplierBigInt)
	}
	if neg {
		value.Neg(value)
	}
	return Int{data: num.Int128FromBigInt(value)}, nil
}

// FromStringForced creates a new value from a string.
func FromStringForced(str string) Int {
	f, _ := FromString(str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// Add adds this value to the passed-in value, returning a new value.
func (f Int) Add(value Int) Int {
	return Int{data: f.data.Add(value.data)}
}

// Sub subtracts the passed-in value from this value, returning a new value.
func (f Int) Sub(value Int) Int {
	return Int{data: f.data.Sub(value.data)}
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f Int) Mul(value Int) Int {
	return Int{data: f.data.Mul(value.data).Div(multiplier)}
}

// Div divides this value by the passed-in value, returning a new value.
func (f Int) Div(value Int) Int {
	return Int{data: f.data.Mul(multiplier).Div(value.data)}
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f Int) Mod(value Int) Int {
	return f.Sub(value.Mul(f.Div(value).Trunc()))
}

// Neg negates this value, returning a new value.
func (f Int) Neg() Int {
	return Int{data: f.data.Neg()}
}

// Abs returns the absolute value of this value.
func (f Int) Abs() Int {
	return Int{data: f.data.Abs()}
}

// Cmp returns 1 if i > n, 0 if i == n, and -1 if i < n.
func (f Int) Cmp(n Int) int {
	return f.data.Cmp(n.data)
}

// GreaterThan returns true if i > n.
func (f Int) GreaterThan(n Int) bool {
	return f.data.GreaterThan(n.data)
}

// GreaterThanOrEqual returns true if i >= n.
func (f Int) GreaterThanOrEqual(n Int) bool {
	return f.data.GreaterThanOrEqual(n.data)
}

// Equal returns true if i == n.
func (f Int) Equal(n Int) bool {
	return f.data.Equal(n.data)
}

// LessThan returns true if i < n.
func (f Int) LessThan(n Int) bool {
	return f.data.LessThan(n.data)
}

// LessThanOrEqual returns true if i <= n.
func (f Int) LessThanOrEqual(n Int) bool {
	return f.data.LessThanOrEqual(n.data)
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f Int) Trunc() Int {
	return Int{data: f.data.Div(multiplier).Mul(multiplier)}
}

// Ceil returns the value rounded up to the nearest whole number.
func (f Int) Ceil() Int {
	v := f.Trunc()
	if f.GreaterThan(Int{}) && f != v {
		v = v.Add(One)
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f Int) Round() Int {
	value := f.Trunc()
	rem := f.Sub(value)
	if rem.GreaterThanOrEqual(Half) {
		value = value.Add(One)
	} else if rem.LessThan(NegHalf) {
		value = value.Sub(One)
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f Int) Min(value Int) Int {
	if f.data.LessThan(value.data) {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f Int) Max(value Int) Int {
	if f.data.GreaterThan(value.data) {
		return f
	}
	return value
}

// Inc returns the value incremented by 1.
func (f Int) Inc() Int {
	return f.Add(One)
}

// Dec returns the value decremented by 1.
func (f Int) Dec() Int {
	return f.Sub(One)
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f Int) AsInt64() int64 {
	return f.data.Div(multiplier).AsInt64()
}

// Int64 is the same as AsInt64(), except that it returns an error if the value cannot be represented exactly with an
// int64.
func (f Int) Int64() (int64, error) {
	n := f.AsInt64()
	if FromInt64(n) != f {
		return 0, internal.ErrDoesNotFitInInt64
	}
	return n, nil
}

// AsInt returns the truncated equivalent integer to this value.
func (f Int) AsInt() int {
	return int(f.data.Div(multiplier).AsInt64())
}

// Int is the same as AsInt(), except that it returns an error if the value cannot be represented exactly with an int.
func (f Int) Int() (int, error) {
	n := f.AsInt()
	if FromInt(n) != f {
		return 0, internal.ErrDoesNotFitInInt
	}
	return n, nil
}

// AsFloat64 returns the floating-point equivalent to this value.
func (f Int) AsFloat64() float64 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierBigFloat).Float64()
	return f64
}

// Float64 is the same as AsFloat64(), except that it returns an error if the value cannot be represented exactly with a
// float64.
func (f Int) Float64() (float64, error) {
	n := f.AsFloat64()
	if strconv.FormatFloat(n, 'g', -1, 64) != f.String() {
		return 0, internal.ErrDoesNotFitInFloat64
	}
	return n, nil
}

// AsFloat32 returns the floating-point equivalent to this value.
func (f Int) AsFloat32() float32 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierBigFloat).Float32()
	return f64
}

// Float32 is the same as AsFloat32(), except that it returns an error if the value cannot be represented exactly with a
// float32.
func (f Int) Float32() (float32, error) {
	n := f.AsFloat32()
	if strconv.FormatFloat(float64(n), 'g', -1, 32) != f.String() {
		return 0, internal.ErrDoesNotFitInFloat32
	}
	return n, nil
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive.
func (f Int) CommaWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f Int) Comma() string {
	var iStr string
	integer := f.data.Div(multiplier)
	if integer.IsInt64() {
		iStr = humanize.Comma(integer.AsInt64())
	} else {
		iStr = humanize.BigComma(integer.AsBigInt())
	}
	fraction := f.data.Sub(integer.Mul(multiplier))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplier).String()
	for i := len(fStr) - 1; i > 0; i-- {
		if fStr[i] != '0' {
			fStr = fStr[1 : i+1]
			break
		}
	}
	var neg string
	if integer.IsZero() && f.data.Sign() < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, iStr, fStr)
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive.
func (f Int) StringWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f Int) String() string {
	integer := f.data.Div(multiplier)
	iStr := integer.String()
	fraction := f.data.Sub(integer.Mul(multiplier))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplier).String()
	for i := len(fStr) - 1; i > 0; i-- {
		if fStr[i] != '0' {
			fStr = fStr[1 : i+1]
			break
		}
	}
	var neg string
	if integer.IsZero() && f.data.Sign() < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, iStr, fStr)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (f Int) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *Int) UnmarshalText(text []byte) error {
	f1, err := FromString(internal.Unquote(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f Int) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Int) UnmarshalJSON(in []byte) error {
	v, err := FromString(internal.Unquote(in))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f Int) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: f.String(),
	}, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *Int) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := FromString(str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// FromFloat64 creates a new value from a float64. Here as a convenience for the eval package.
func (f *Int) FromFloat64(value float64) Int {
	return FromFloat64(value)
}

// FromFloat32 creates a new value from a float32. Here as a convenience for the eval package.
func (f *Int) FromFloat32(value float32) Int {
	return FromFloat32(value)
}

// FromInt64 creates a new value from an int64. Here as a convenience for the eval package.
func (f *Int) FromInt64(value int64) Int {
	return FromInt64(value)
}

// FromInt creates a new value from an int. Here as a convenience for the eval package.
func (f *Int) FromInt(value int) Int {
	return FromInt(value)
}

// FromString creates a new value from a string. Here as a convenience for the eval package.
func (f *Int) FromString(str string) (Int, error) {
	return FromString(str)
}
