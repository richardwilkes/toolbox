// Code created from "fixed64.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f64d6

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed/internal"
	"gopkg.in/yaml.v3"
)

const (
	multiplier = 1000000
	// Max holds the maximum value.
	Max = Int(1<<63 - 1)
	// Min holds the minimum value.
	Min = Int(^(1<<63 - 1))
)

// Some commonly used values.
var (
	One     = Int(multiplier)
	Half    = Int(multiplier / 2)
	NegHalf = -Half
)

// Int holds a fixed-point value that contains up to 6 decimal places. Values are truncated, not rounded. Values can be
// added and subtracted directly. For multiplication and division, the provided Mul() and Div() methods should be used.
type Int int64

// FromFloat64 creates a new value from a float64.
func FromFloat64(value float64) Int {
	return Int(value * multiplier)
}

// FromFloat32 creates a new value from a float32.
func FromFloat32(value float32) Int {
	return Int(float64(value) * multiplier)
}

// FromInt64 creates a new value from an int64.
func FromInt64(value int64) Int {
	return Int(value * multiplier)
}

// FromInt creates a new value from an int.
func FromInt(value int) Int {
	return Int(int64(value) * multiplier)
}

// FromString creates a new value from a string.
func FromString(str string) (Int, error) {
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
		return FromFloat64(f), nil
	}
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
		value *= multiplier
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 6+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 6+1 {
			frac = frac[:6+1]
		}
		if fraction, err = strconv.ParseInt(frac, 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		value += fraction - multiplier
	}
	if neg {
		value = -value
	}
	return Int(value), nil
}

// FromStringForced creates a new value from a string.
func FromStringForced(str string) Int {
	f, _ := FromString(str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f Int) Mul(value Int) Int {
	return f * value / multiplier
}

// Div divides this value by the passed-in value, returning a new value.
func (f Int) Div(value Int) Int {
	return f * multiplier / value
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f Int) Mod(value Int) Int {
	return f - (value.Mul(f.Div(value).Trunc()))
}

// Abs returns the absolute value of this value.
func (f Int) Abs() Int {
	if f < 0 {
		return -f
	}
	return f
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f Int) Trunc() Int {
	return f / multiplier * multiplier
}

// Ceil returns the value rounded up to the nearest whole number.
func (f Int) Ceil() Int {
	v := f.Trunc()
	if f > 0 && f != v {
		v += One
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f Int) Round() Int {
	value := f.Trunc()
	rem := f - value //nolint:ifshort // don't want to embed this in the if
	if rem >= Half {
		value += One
	} else if rem < NegHalf {
		value -= One
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f Int) Min(value Int) Int {
	if f < value {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f Int) Max(value Int) Int {
	if f > value {
		return f
	}
	return value
}

// Inc returns the value incremented by 1.
func (f Int) Inc() Int {
	return f + One
}

// Dec returns the value decremented by 1.
func (f Int) Dec() Int {
	return f - One
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f Int) AsInt64() int64 {
	return int64(f / multiplier)
}

// Int64 is the same as AsInt64(), except that it returns an error if the value cannot be represented exactly with an
// int64
func (f Int) Int64() (int64, error) {
	n := f.AsInt64()
	if FromInt64(n) != f {
		return 0, internal.ErrDoesNotFitInInt64
	}
	return n, nil
}

// AsInt returns the truncated equivalent integer to this value.
func (f Int) AsInt() int {
	return int(f / multiplier)
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
	return float64(f) / multiplier
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
	return float32(f.AsFloat64())
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

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive
func (f Int) CommaWithSign() string {
	if f >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f Int) Comma() string {
	integer := f / multiplier
	fraction := f % multiplier
	if fraction == 0 {
		return humanize.Comma(int64(integer))
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += multiplier
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
func (f Int) StringWithSign() string {
	if f >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f Int) String() string {
	integer := f / multiplier
	fraction := f % multiplier
	if fraction == 0 {
		return strconv.FormatInt(int64(integer), 10)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += multiplier
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
