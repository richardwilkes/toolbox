// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed128

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/num128"
	"github.com/richardwilkes/toolbox/v2/xstrings"
	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v3"
)

// Int holds a fixed-point value. Values are truncated, not rounded.
type Int[T fixed.Dx] struct {
	data num128.Int
}

// Maximum returns the maximum possible value the type can hold.
func Maximum[T fixed.Dx]() Int[T] {
	return Int[T]{data: num128.MaxInt}
}

// Minimum returns the minimum possible value the type can hold.
func Minimum[T fixed.Dx]() Int[T] {
	return Int[T]{data: num128.MinInt}
}

// MaxDecimalDigits returns the maximum number of digits after the decimal that will be used.
func MaxDecimalDigits[T fixed.Dx]() int {
	var t T
	return t.Places()
}

// Multiplier returns the multiplier used.
func Multiplier[T fixed.Dx]() Int[T] {
	return Int[T]{data: multiplier[T]()}
}

func multiplier[T fixed.Dx]() num128.Int {
	var t T
	return num128.IntFrom64(t.Multiplier())
}

// FromInteger creates a new value.
func FromInteger[T fixed.Dx, FROM constraints.Integer](value FROM) Int[T] {
	var t T
	return Int[T]{data: num128.IntFrom64(int64(value)).Mul(num128.IntFrom64(t.Multiplier()))}
}

// FromFloat creates a new value.
func FromFloat[T fixed.Dx, FROM constraints.Float](value FROM) Int[T] {
	f, _ := FromString[T](new(big.Float).SetPrec(128).SetFloat64(float64(value)).Text('f', MaxDecimalDigits[T]()+1)) //nolint:errcheck // Failure means 0
	return f
}

// FromString creates a new value from a string.
func FromString[T fixed.Dx](str string) (Int[T], error) {
	if str == "" {
		return Int[T]{}, errs.New("empty string is not valid")
	}
	str = strings.ReplaceAll(str, ",", "")
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return Int[T]{}, err
		}
		return FromFloat[T](f), nil
	}
	parts := strings.SplitN(str, ".", 2)
	var neg bool
	value := new(big.Int)
	fraction := new(big.Int)
	var t T
	switch parts[0] {
	case "":
	case "-", "-0":
		neg = true
	default:
		if _, ok := value.SetString(parts[0], 10); !ok {
			return Int[T]{}, errs.Newf("invalid value: %s", str)
		}
		if value.Sign() < 0 {
			neg = true
			value.Neg(value)
		}
		value.Mul(value, big.NewInt(t.Multiplier()))
	}
	if len(parts) > 1 {
		cutoff := 1 + MaxDecimalDigits[T]()
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < cutoff {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > cutoff {
			frac = frac[:cutoff]
		}
		if _, ok := fraction.SetString(frac, 10); !ok {
			return Int[T]{}, errs.Newf("invalid value: %s", str)
		}
		value.Add(value, fraction).Sub(value, big.NewInt(t.Multiplier()))
	}
	if neg {
		value.Neg(value)
	}
	return Int[T]{data: num128.IntFromBigInt(value)}, nil
}

// FromStringForced creates a new value from a string.
func FromStringForced[T fixed.Dx](str string) Int[T] {
	f, _ := FromString[T](str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// AsInteger returns the equivalent value in the destination type.
func AsInteger[T fixed.Dx, TO constraints.Integer](f Int[T]) TO {
	return TO(f.data.Div(multiplier[T]()).AsInt64())
}

// AsFloat returns the equivalent value in the destination type.
func AsFloat[T fixed.Dx, TO constraints.Float](f Int[T]) TO {
	var t T
	fixed64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(),
		new(big.Float).SetPrec(128).SetInt(big.NewInt(t.Multiplier()))).Float64()
	return TO(fixed64)
}

// AsIntegerChecked is the same as As(), except that it returns an error if the value cannot be represented exactly in
// the requested destination type.
func AsIntegerChecked[T fixed.Dx, TO constraints.Integer](f Int[T]) (TO, error) {
	n := TO(f.data.Div(multiplier[T]()).AsInt64())
	if FromInteger[T](n) != f {
		return 0, fixed.ErrDoesNotFitInRequestedType
	}
	return n, nil
}

// AsFloatChecked is the same as As(), except that it returns an error if the value cannot be represented exactly in the
// requested destination type.
func AsFloatChecked[T fixed.Dx, TO constraints.Float](f Int[T]) (TO, error) {
	var t T
	fixed64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(),
		new(big.Float).SetPrec(128).SetInt(big.NewInt(t.Multiplier()))).Float64()
	n := TO(fixed64)
	if strconv.FormatFloat(float64(n), 'g', -1, reflect.TypeOf(n).Bits()) != f.String() {
		return 0, fixed.ErrDoesNotFitInRequestedType
	}
	return n, nil
}

// Add adds this value to the passed-in value, returning a new value.
func (f Int[T]) Add(value Int[T]) Int[T] {
	return Int[T]{data: f.data.Add(value.data)}
}

// Sub subtracts the passed-in value from this value, returning a new value.
func (f Int[T]) Sub(value Int[T]) Int[T] {
	return Int[T]{data: f.data.Sub(value.data)}
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f Int[T]) Mul(value Int[T]) Int[T] {
	return Int[T]{data: f.data.Mul(value.data).Div(multiplier[T]())}
}

// Div divides this value by the passed-in value, returning a new value.
func (f Int[T]) Div(value Int[T]) Int[T] {
	return Int[T]{data: f.data.Mul(multiplier[T]()).Div(value.data)}
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f Int[T]) Mod(value Int[T]) Int[T] {
	return f.Sub(value.Mul(f.Div(value).Trunc()))
}

// Neg negates this value, returning a new value.
func (f Int[T]) Neg() Int[T] {
	return Int[T]{data: f.data.Neg()}
}

// Abs returns the absolute value of this value.
func (f Int[T]) Abs() Int[T] {
	return Int[T]{data: f.data.Abs()}
}

// Cmp returns 1 if i > n, 0 if i == n, and -1 if i < n.
func (f Int[T]) Cmp(n Int[T]) int {
	return f.data.Cmp(n.data)
}

// GreaterThan returns true if i > n.
func (f Int[T]) GreaterThan(n Int[T]) bool {
	return f.data.GreaterThan(n.data)
}

// GreaterThanOrEqual returns true if i >= n.
func (f Int[T]) GreaterThanOrEqual(n Int[T]) bool {
	return f.data.GreaterThanOrEqual(n.data)
}

// Equal returns true if i == n.
func (f Int[T]) Equal(n Int[T]) bool {
	return f.data.Equal(n.data)
}

// LessThan returns true if i < n.
func (f Int[T]) LessThan(n Int[T]) bool {
	return f.data.LessThan(n.data)
}

// LessThanOrEqual returns true if i <= n.
func (f Int[T]) LessThanOrEqual(n Int[T]) bool {
	return f.data.LessThanOrEqual(n.data)
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f Int[T]) Trunc() Int[T] {
	m := multiplier[T]()
	return Int[T]{data: f.data.Div(m).Mul(m)}
}

// Ceil returns the value rounded up to the nearest whole number.
func (f Int[T]) Ceil() Int[T] {
	v := f.Trunc()
	if f.GreaterThan(Int[T]{}) && f != v {
		v = v.Add(Multiplier[T]())
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f Int[T]) Round() Int[T] {
	one := Multiplier[T]()
	half := Int[T]{data: one.data.Div(num128.IntFrom64(2))}
	negHalf := half.Neg()
	value := f.Trunc()
	rem := f.Sub(value)
	if rem.GreaterThanOrEqual(half) {
		value = value.Add(one)
	} else if rem.LessThan(negHalf) {
		value = value.Sub(one)
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f Int[T]) Min(value Int[T]) Int[T] {
	if f.data.LessThan(value.data) {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f Int[T]) Max(value Int[T]) Int[T] {
	if f.data.GreaterThan(value.data) {
		return f
	}
	return value
}

// Inc returns the value incremented by 1.
func (f Int[T]) Inc() Int[T] {
	return f.Add(Multiplier[T]())
}

// Dec returns the value decremented by 1.
func (f Int[T]) Dec() Int[T] {
	return f.Sub(Multiplier[T]())
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive.
func (f Int[T]) CommaWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f Int[T]) Comma() string {
	return xstrings.CommaFromStringNum(f.String())
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive.
func (f Int[T]) StringWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f Int[T]) String() string {
	mult := multiplier[T]()
	integer := f.data.Div(mult)
	iStr := integer.String()
	fraction := f.data.Sub(integer.Mul(mult))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(mult).String()
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
func (f Int[T]) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *Int[T]) UnmarshalText(text []byte) error {
	f1, err := FromString[T](xstrings.Unquote(string(text)))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f Int[T]) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Int[T]) UnmarshalJSON(in []byte) error {
	v, err := FromString[T](xstrings.Unquote(string(in)))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f Int[T]) MarshalYAML() (any, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: f.String(),
	}, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *Int[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := FromString[T](str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}
