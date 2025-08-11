// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed64

import (
	"fmt"
	"math"
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

// Int holds a fixed-point value. Values are truncated, not rounded. Values can be added and subtracted directly. For
// multiplication and division, the provided Mul() and Div() methods should be used.
type Int[T fixed.Dx] int64

// Maximum returns the maximum possible value the type can hold.
func Maximum[T fixed.Dx]() Int[T] {
	return Int[T](math.MaxInt64)
}

// Minimum returns the minimum possible value the type can hold.
func Minimum[T fixed.Dx]() Int[T] {
	return Int[T](math.MinInt64)
}

// MaxDecimalDigits returns the maximum number of digits after the decimal that will be used.
func MaxDecimalDigits[T fixed.Dx]() int {
	var t T
	return t.Places()
}

// Multiplier returns the multiplier used.
func Multiplier[T fixed.Dx]() int64 {
	var t T
	return t.Multiplier()
}

// FromInteger creates a new value.
func FromInteger[T fixed.Dx, FROM constraints.Integer](value FROM) Int[T] {
	return Int[T](value * FROM(Multiplier[T]()))
}

// FromFloat creates a new value.
func FromFloat[T fixed.Dx, FROM constraints.Float](value FROM) Int[T] {
	return Int[T](value * FROM(Multiplier[T]()))
}

// FromString creates a new value from a string.
func FromString[T fixed.Dx](str string) (Int[T], error) {
	if str == "" {
		return 0, errs.New("empty string is not valid")
	}
	str = strings.ReplaceAll(str, ",", "")
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, err
		}
		return FromFloat[T](f), nil
	}
	mult := Multiplier[T]()
	parts := strings.SplitN(str, ".", 2)
	var value, fraction int64
	var neg bool
	var err error
	switch parts[0] {
	case "":
	case "-", "-0":
		neg = true
	default:
		if value, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		if value < 0 {
			neg = true
			value = -value
		}
		value *= mult
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
		if fraction, err = strconv.ParseInt(frac, 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		value += fraction - mult
	}
	if neg {
		value = -value
	}
	return Int[T](value), nil
}

// FromStringForced creates a new value from a string.
func FromStringForced[T fixed.Dx](str string) Int[T] {
	f, _ := FromString[T](str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// AsInteger returns the equivalent value in the destination type.
func AsInteger[T fixed.Dx, TO constraints.Integer](f Int[T]) TO {
	return TO(int64(f) / Multiplier[T]())
}

// AsFloat returns the equivalent value in the destination type.
func AsFloat[T fixed.Dx, TO constraints.Float](f Int[T]) TO {
	return TO(float64(f) / float64(Multiplier[T]()))
}

// AsIntegerChecked is the same as AsInteger(), except that it returns an error if the value cannot be represented
// exactly in the requested destination type.
func AsIntegerChecked[T fixed.Dx, TO constraints.Integer](f Int[T]) (TO, error) {
	n := TO(int64(f) / Multiplier[T]())
	if FromInteger[T](n) != f {
		return 0, fixed.ErrDoesNotFitInRequestedType
	}
	return n, nil
}

// AsFloatChecked is the same as AsFloat(), except that it returns an error if the value cannot be represented exactly
// in the requested destination type.
func AsFloatChecked[T fixed.Dx, TO constraints.Float](f Int[T]) (TO, error) {
	n := TO(float64(f) / float64(Multiplier[T]()))
	if strconv.FormatFloat(float64(n), 'g', -1, reflect.TypeOf(n).Bits()) != f.String() {
		return 0, fixed.ErrDoesNotFitInRequestedType
	}
	return n, nil
}

// Add adds this value to the passed-in value, returning a new value. Note that this method is only provided to make
// text templates easier to use with these objects, since you can just add two Int[T] values together like they were
// primitive types.
func (f Int[T]) Add(value Int[T]) Int[T] {
	return f + value
}

// Sub subtracts the passed-in value from this value, returning a new value. Note that this method is only provided to
// make text templates easier to use with these objects, since you can just subtract two Int[T] values together like
// they were primitive types.
func (f Int[T]) Sub(value Int[T]) Int[T] {
	return f - value
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f Int[T]) Mul(value Int[T]) Int[T] {
	return f.mul64(value, Int[T](Multiplier[T]()))
}

func (f Int[T]) mul64(value, div Int[T]) Int[T] {
	if f == 0 || value == 0 {
		return 0
	}
	result := f * value
	if f != math.MinInt64 && value != math.MinInt64 && result/value == f {
		return result / div
	}
	return Int[T](num128.IntFrom64(int64(f)).
		Mul(num128.IntFrom64(int64(value))).
		Div(num128.IntFrom64(int64(div))).
		AsInt64())
}

// Div divides this value by the passed-in value, returning a new value.
func (f Int[T]) Div(value Int[T]) Int[T] {
	return f.mul64(Int[T](Multiplier[T]()), value)
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f Int[T]) Mod(value Int[T]) Int[T] {
	return f - (value.Mul(f.Div(value).Floor()))
}

// Abs returns the absolute value of this value.
func (f Int[T]) Abs() Int[T] {
	if f < 0 {
		return -f
	}
	return f
}

// Floor returns the value rounded down to the nearest whole number.
func (f Int[T]) Floor() Int[T] {
	mult := Int[T](Multiplier[T]())
	return f / mult * mult
}

// Ceil returns the value rounded up to the nearest whole number.
func (f Int[T]) Ceil() Int[T] {
	v := f.Floor()
	if f > 0 && f != v {
		v += Int[T](Multiplier[T]())
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f Int[T]) Round() Int[T] {
	one := Int[T](Multiplier[T]())
	value := f.Floor()
	rem := f - value
	if rem >= one/2 {
		value += one
	} else if rem < -one/2 {
		value -= one
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f Int[T]) Min(value Int[T]) Int[T] {
	if f < value {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f Int[T]) Max(value Int[T]) Int[T] {
	if f > value {
		return f
	}
	return value
}

// Inc returns the value incremented by 1.
func (f Int[T]) Inc() Int[T] {
	return f + Int[T](Multiplier[T]())
}

// Dec returns the value decremented by 1.
func (f Int[T]) Dec() Int[T] {
	return f - Int[T](Multiplier[T]())
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive
func (f Int[T]) CommaWithSign() string {
	if f >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f Int[T]) Comma() string {
	return xstrings.CommaFromStringNum(f.String())
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive
func (f Int[T]) StringWithSign() string {
	if f >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f Int[T]) String() string {
	mult := Int[T](Multiplier[T]())
	integer := f / mult
	fraction := f % mult
	if fraction == 0 {
		return strconv.FormatInt(int64(integer), 10)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += mult
	fStr := strconv.FormatInt(int64(fraction), 10)
	for i := len(fStr) - 1; i > 0; i-- {
		if fStr[i] != '0' {
			fStr = fStr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && f < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%d.%s", neg, integer, fStr)
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
