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

package fixed

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/num"
	"gopkg.in/yaml.v3"
)

var (
	// F128d2Max holds the maximum F128d2 value.
	F128d2Max = F128d2{data: num.MaxInt128}
	// F128d2Min holds the minimum F128d2 value.
	F128d2Min                = F128d2{data: num.MinInt128}
	multiplierF128d2BigInt   = new(big.Int).Exp(big.NewInt(10), big.NewInt(2), nil)
	multiplierF128d2BigFloat = new(big.Float).SetPrec(128).SetInt(multiplierF128d2BigInt)
	multiplierF128d2         = num.Int128FromBigInt(multiplierF128d2BigInt)
)

// Some commonly used values.
var (
	F128d2One     = F128d2FromInt(1)
	F128d2Half    = F128d2FromStringForced("0.5")
	F128d2NegHalf = F128d2Half.Neg()
)

// F128d2 holds a fixed-point value that contains up to 2 decimal places. Values are truncated, not rounded.
type F128d2 struct {
	data num.Int128
}

// F128d2FromFloat64 creates a new F128d2 value from a float64.
func F128d2FromFloat64(value float64) F128d2 {
	f, _ := F128d2FromString(new(big.Float).SetPrec(128).SetFloat64(value).Text('f', 3)) //nolint:errcheck // Failure means 0
	return f
}

// F128d2FromFloat32 creates a new F128d2 value from a float32.
func F128d2FromFloat32(value float32) F128d2 {
	return F128d2FromFloat64(float64(value))
}

// F128d2FromInt64 creates a new F128d2 value from an int64.
func F128d2FromInt64(value int64) F128d2 {
	return F128d2{data: num.Int128From64(value).Mul(multiplierF128d2)}
}

// F128d2FromInt creates a new F128d2 value from an int.
func F128d2FromInt(value int) F128d2 {
	return F128d2FromInt64(int64(value))
}

// F128d2FromString creates a new F128d2 value from a string.
func F128d2FromString(str string) (F128d2, error) {
	if str == "" {
		return F128d2{}, errs.New("empty string is not valid")
	}
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return F128d2{}, err
		}
		return F128d2FromFloat64(f), nil
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
			return F128d2{}, errs.Newf("invalid value: %s", str)
		}
		if value.Sign() < 0 {
			neg = true
			value.Neg(value)
		}
		value.Mul(value, multiplierF128d2BigInt)
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 2+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 2+1 {
			frac = frac[:2+1]
		}
		if _, ok := fraction.SetString(frac, 10); !ok {
			return F128d2{}, errs.Newf("invalid value: %s", str)
		}
		value.Add(value, fraction).Sub(value, multiplierF128d2BigInt)
	}
	if neg {
		value.Neg(value)
	}
	return F128d2{data: num.Int128FromBigInt(value)}, nil
}

// F128d2FromStringForced creates a new F128d2 value from a string.
func F128d2FromStringForced(str string) F128d2 {
	f, _ := F128d2FromString(str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// Add adds this value to the passed-in value, returning a new value.
func (f F128d2) Add(value F128d2) F128d2 {
	return F128d2{data: f.data.Add(value.data)}
}

// Sub subtracts the passed-in value from this value, returning a new value.
func (f F128d2) Sub(value F128d2) F128d2 {
	return F128d2{data: f.data.Sub(value.data)}
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F128d2) Mul(value F128d2) F128d2 {
	return F128d2{data: f.data.Mul(value.data).Div(multiplierF128d2)}
}

// Div divides this value by the passed-in value, returning a new value.
func (f F128d2) Div(value F128d2) F128d2 {
	return F128d2{data: f.data.Mul(multiplierF128d2).Div(value.data)}
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f F128d2) Mod(value F128d2) F128d2 {
	return f.Sub(value.Mul(f.Div(value).Trunc()))
}

// Neg negates this value, returning a new value.
func (f F128d2) Neg() F128d2 {
	return F128d2{data: f.data.Neg()}
}

// Abs returns the absolute value of this value.
func (f F128d2) Abs() F128d2 {
	return F128d2{data: f.data.Abs()}
}

// Cmp returns 1 if i > n, 0 if i == n, and -1 if i < n.
func (f F128d2) Cmp(n F128d2) int {
	return f.data.Cmp(n.data)
}

// GreaterThan returns true if i > n.
func (f F128d2) GreaterThan(n F128d2) bool {
	return f.data.GreaterThan(n.data)
}

// GreaterThanOrEqual returns true if i >= n.
func (f F128d2) GreaterThanOrEqual(n F128d2) bool {
	return f.data.GreaterThanOrEqual(n.data)
}

// Equal returns true if i == n.
func (f F128d2) Equal(n F128d2) bool {
	return f.data.Equal(n.data)
}

// LessThan returns true if i < n.
func (f F128d2) LessThan(n F128d2) bool {
	return f.data.LessThan(n.data)
}

// LessThanOrEqual returns true if i <= n.
func (f F128d2) LessThanOrEqual(n F128d2) bool {
	return f.data.LessThanOrEqual(n.data)
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f F128d2) Trunc() F128d2 {
	return F128d2{data: f.data.Div(multiplierF128d2).Mul(multiplierF128d2)}
}

// Ceil returns the value rounded up to the nearest whole number.
func (f F128d2) Ceil() F128d2 {
	v := f.Trunc()
	if f.GreaterThan(F128d2{}) && f != v {
		v = v.Add(F128d2One)
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f F128d2) Round() F128d2 {
	value := f.Trunc()
	rem := f.Sub(value)
	if rem.GreaterThanOrEqual(F128d2Half) {
		value = value.Add(F128d2One)
	} else if rem.LessThan(F128d2NegHalf) {
		value = value.Sub(F128d2One)
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f F128d2) Min(value F128d2) F128d2 {
	if f.data.LessThan(value.data) {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f F128d2) Max(value F128d2) F128d2 {
	if f.data.GreaterThan(value.data) {
		return f
	}
	return value
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f F128d2) AsInt64() int64 {
	return f.data.Div(multiplierF128d2).AsInt64()
}

// Int64 is the same as AsInt64(), except that it returns an error if the value cannot be represented exactly with an
// int64.
func (f F128d2) Int64() (int64, error) {
	n := f.AsInt64()
	if F128d2FromInt64(n) != f {
		return 0, errDoesNotFitInInt64
	}
	return n, nil
}

// AsInt returns the truncated equivalent integer to this value.
func (f F128d2) AsInt() int {
	return int(f.data.Div(multiplierF128d2).AsInt64())
}

// Int is the same as AsInt(), except that it returns an error if the value cannot be represented exactly with an int.
func (f F128d2) Int() (int, error) {
	n := f.AsInt()
	if F128d2FromInt(n) != f {
		return 0, errDoesNotFitInInt
	}
	return n, nil
}

// AsFloat64 returns the floating-point equivalent to this value.
func (f F128d2) AsFloat64() float64 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierF128d2BigFloat).Float64()
	return f64
}

// Float64 is the same as AsFloat64(), except that it returns an error if the value cannot be represented exactly with a
// float64.
func (f F128d2) Float64() (float64, error) {
	n := f.AsFloat64()
	if strconv.FormatFloat(n, 'g', -1, 64) != f.String() {
		return 0, errDoesNotFitInFloat64
	}
	return n, nil
}

// AsFloat32 returns the floating-point equivalent to this value.
func (f F128d2) AsFloat32() float32 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierF128d2BigFloat).Float32()
	return f64
}

// Float32 is the same as AsFloat32(), except that it returns an error if the value cannot be represented exactly with a
// float32.
func (f F128d2) Float32() (float32, error) {
	n := f.AsFloat32()
	if strconv.FormatFloat(float64(n), 'g', -1, 32) != f.String() {
		return 0, errDoesNotFitInFloat32
	}
	return n, nil
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive
func (f F128d2) CommaWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f F128d2) Comma() string {
	var iStr string
	integer := f.data.Div(multiplierF128d2)
	if integer.IsInt64() {
		iStr = humanize.Comma(integer.AsInt64())
	} else {
		iStr = humanize.BigComma(integer.AsBigInt())
	}
	fraction := f.data.Sub(integer.Mul(multiplierF128d2))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplierF128d2).String()
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

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive
func (f F128d2) StringWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f F128d2) String() string {
	integer := f.data.Div(multiplierF128d2)
	iStr := integer.String()
	fraction := f.data.Sub(integer.Mul(multiplierF128d2))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplierF128d2).String()
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
func (f F128d2) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F128d2) UnmarshalText(text []byte) error {
	f1, err := F128d2FromString(unquote(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f F128d2) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *F128d2) UnmarshalJSON(in []byte) error {
	v, err := F128d2FromString(unquote(in))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f F128d2) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: f.String(),
	}, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *F128d2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := F128d2FromString(str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}
