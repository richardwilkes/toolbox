// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f64

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"gopkg.in/yaml.v3"
)

const (
	// Max holds the maximum value.
	Max = math.MaxInt64
	// Min holds the minimum value.
	Min = math.MinInt64
)

// Int holds a fixed-point value. Values are truncated, not rounded. Values can be added and subtracted directly. For
// multiplication and division, the provided Mul() and Div() methods should be used.
type Int[T fixed.Dx] int64

// MaxSafeMultiply returns the maximum value that can be safely multiplied without overflow.
func MaxSafeMultiply[T fixed.Dx]() Int[T] {
	return Int[T](Max / Multiplier[T]())
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

// From creates a new value.
func From[T fixed.Dx, FROM xmath.Numeric](value FROM) Int[T] {
	return Int[T](value * FROM(Multiplier[T]()))
}

// FromString creates a new value from a string.
func FromString[T fixed.Dx](str string) (Int[T], error) {
	if str == "" {
		return 0, errs.New("empty string is not valid")
	}
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, err
		}
		return From[T](f), nil
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

// Mul multiplies this value by the passed-in value, returning a new value.
func (f Int[T]) Mul(value Int[T]) Int[T] {
	return f * value / Int[T](Multiplier[T]())
}

// Div divides this value by the passed-in value, returning a new value.
func (f Int[T]) Div(value Int[T]) Int[T] {
	return f * Int[T](Multiplier[T]()) / value
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f Int[T]) Mod(value Int[T]) Int[T] {
	return f - (value.Mul(f.Div(value).Trunc()))
}

// Abs returns the absolute value of this value.
func (f Int[T]) Abs() Int[T] {
	if f < 0 {
		return -f
	}
	return f
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f Int[T]) Trunc() Int[T] {
	mult := Int[T](Multiplier[T]())
	return f / mult * mult
}

// Ceil returns the value rounded up to the nearest whole number.
func (f Int[T]) Ceil() Int[T] {
	v := f.Trunc()
	if f > 0 && f != v {
		v += Int[T](Multiplier[T]())
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f Int[T]) Round() Int[T] {
	one := Int[T](Multiplier[T]())
	value := f.Trunc()
	rem := f - value //nolint:ifshort // don't want to embed this in the if
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

// As returns the equivalent value in the destination type.
func As[T fixed.Dx, TO xmath.Numeric](f Int[T]) TO {
	var n TO
	switch reflect.TypeOf(n).Kind() {
	case reflect.Float32, reflect.Float64:
		return TO(float64(f) / float64(Multiplier[T]()))
	default:
		return TO(int64(f) / Multiplier[T]())
	}
}

// CheckedAs is the same as As(), except that it returns an error if the value cannot be represented exactly in the
// requested destination type.
func CheckedAs[T fixed.Dx, TO xmath.Numeric](f Int[T]) (TO, error) {
	var n TO
	switch reflect.TypeOf(n).Kind() {
	case reflect.Float32, reflect.Float64:
		n = TO(float64(f) / float64(Multiplier[T]()))
		if strconv.FormatFloat(float64(n), 'g', -1, reflect.TypeOf(n).Bits()) != f.String() {
			return 0, fixed.ErrDoesNotFitInRequestedType
		}
	default:
		n = TO(int64(f) / Multiplier[T]())
		if From[T](n) != f {
			return 0, fixed.ErrDoesNotFitInRequestedType
		}
	}
	return n, nil
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
	mult := Int[T](Multiplier[T]())
	integer := f / mult
	fraction := f % mult
	if fraction == 0 {
		return humanize.Comma(int64(integer))
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
	return fmt.Sprintf("%s%s.%s", neg, humanize.Comma(int64(integer)), fStr)
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
	f1, err := FromString[T](txt.Unquote(string(text)))
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
	v, err := FromString[T](txt.Unquote(string(in)))
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
