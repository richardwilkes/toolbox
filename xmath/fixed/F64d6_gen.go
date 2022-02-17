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

package fixed

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/errs"
	"gopkg.in/yaml.v3"
)

const (
	// F64d6Max holds the maximum F64d6 value.
	F64d6Max = F64d6(1<<63 - 1)
	// F64d6Min holds the minimum F64d6 value.
	F64d6Min = F64d6(^(1<<63 - 1))
)

// Some commonly used values.
var (
	F64d6One     = F64d6FromInt(1)
	F64d6Half    = F64d6FromStringForced("0.5")
	F64d6NegHalf = -F64d6Half
)

var multiplierF64d6 = int64(math.Pow(10, 6))

// F64d6 holds a fixed-point value that contains up to 6 decimal places. Values are truncated, not rounded. Values can
// be added and subtracted directly. For multiplication and division, the provided Mul() and Div() methods should be
// used.
type F64d6 int64

// F64d6FromFloat64 creates a new F64d6 value from a float64.
func F64d6FromFloat64(value float64) F64d6 {
	return F64d6(value * float64(multiplierF64d6))
}

// F64d6FromFloat32 creates a new F64d6 value from a float32.
func F64d6FromFloat32(value float32) F64d6 {
	return F64d6(float64(value) * float64(multiplierF64d6))
}

// F64d6FromInt64 creates a new F64d6 value from an int64.
func F64d6FromInt64(value int64) F64d6 {
	return F64d6(value * multiplierF64d6)
}

// F64d6FromInt creates a new F64d6 value from an int.
func F64d6FromInt(value int) F64d6 {
	return F64d6(int64(value) * multiplierF64d6)
}

// F64d6FromString creates a new F64d6 value from a string.
func F64d6FromString(str string) (F64d6, error) {
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
		return F64d6FromFloat64(f), nil
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
		value *= multiplierF64d6
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
		value += fraction - multiplierF64d6
	}
	if neg {
		value = -value
	}
	return F64d6(value), nil
}

// F64d6FromStringForced creates a new F64d6 value from a string.
func F64d6FromStringForced(str string) F64d6 {
	f, _ := F64d6FromString(str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F64d6) Mul(value F64d6) F64d6 {
	return f * value / F64d6(multiplierF64d6)
}

// Div divides this value by the passed-in value, returning a new value.
func (f F64d6) Div(value F64d6) F64d6 {
	return f * F64d6(multiplierF64d6) / value
}

// Mod returns the remainder after subtracting all full multiples of the passed-in value.
func (f F64d6) Mod(value F64d6) F64d6 {
	return f - (value.Mul(f.Div(value).Trunc()))
}

// Abs returns the absolute value of this value.
func (f F64d6) Abs() F64d6 {
	if f < 0 {
		return -f
	}
	return f
}

// Trunc returns a new value which has everything to the right of the decimal place truncated.
func (f F64d6) Trunc() F64d6 {
	return f / F64d6(multiplierF64d6) * F64d6(multiplierF64d6)
}

// Ceil returns the value rounded up to the nearest whole number.
func (f F64d6) Ceil() F64d6 {
	v := f.Trunc()
	if f > 0 && f != v {
		v += F64d6One
	}
	return v
}

// Round returns the nearest integer, rounding half away from zero.
func (f F64d6) Round() F64d6 {
	value := f.Trunc()
	rem := f - value //nolint:ifshort // don't want to embed this in the if
	if rem >= F64d6Half {
		value += F64d6One
	} else if rem < F64d6NegHalf {
		value -= F64d6One
	}
	return value
}

// Min returns the minimum of this value or its argument.
func (f F64d6) Min(value F64d6) F64d6 {
	if f < value {
		return f
	}
	return value
}

// Max returns the maximum of this value or its argument.
func (f F64d6) Max(value F64d6) F64d6 {
	if f > value {
		return f
	}
	return value
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f F64d6) AsInt64() int64 {
	return int64(f / F64d6(multiplierF64d6))
}

// Int64 is the same as AsInt64(), except that it returns an error if the value cannot be represented exactly with an
// int64.
func (f F64d6) Int64() (int64, error) {
	n := f.AsInt64()
	if F64d6FromInt64(n) != f {
		return 0, errDoesNotFitInInt64
	}
	return n, nil
}

// AsInt returns the truncated equivalent integer to this value.
func (f F64d6) AsInt() int {
	return int(f / F64d6(multiplierF64d6))
}

// Int is the same as AsInt(), except that it returns an error if the value cannot be represented exactly with an int.
func (f F64d6) Int() (int, error) {
	n := f.AsInt()
	if F64d6FromInt(n) != f {
		return 0, errDoesNotFitInInt
	}
	return n, nil
}

// AsFloat64 returns the floating-point equivalent to this value.
func (f F64d6) AsFloat64() float64 {
	return float64(f) / float64(multiplierF64d6)
}

// Float64 is the same as AsFloat64(), except that it returns an error if the value cannot be represented exactly with a
// float64.
func (f F64d6) Float64() (float64, error) {
	n := f.AsFloat64()
	if strconv.FormatFloat(n, 'g', -1, 64) != f.String() {
		return 0, errDoesNotFitInFloat64
	}
	return n, nil
}

// AsFloat32 returns the floating-point equivalent to this value.
func (f F64d6) AsFloat32() float32 {
	return float32(f.AsFloat64())
}

// Float32 is the same as AsFloat32(), except that it returns an error if the value cannot be represented exactly with a
// float32.
func (f F64d6) Float32() (float32, error) {
	n := f.AsFloat32()
	if strconv.FormatFloat(float64(n), 'g', -1, 32) != f.String() {
		return 0, errDoesNotFitInFloat32
	}
	return n, nil
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive
func (f F64d6) CommaWithSign() string {
	if f >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and greater.
func (f F64d6) Comma() string {
	integer := f / F64d6(multiplierF64d6)
	fraction := f % F64d6(multiplierF64d6)
	if fraction == 0 {
		return humanize.Comma(int64(integer))
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d6(multiplierF64d6)
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
func (f F64d6) StringWithSign() string {
	if f >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f F64d6) String() string {
	integer := f / F64d6(multiplierF64d6)
	fraction := f % F64d6(multiplierF64d6)
	if fraction == 0 {
		return strconv.FormatInt(int64(integer), 10)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d6(multiplierF64d6)
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
func (f F64d6) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F64d6) UnmarshalText(text []byte) error {
	f1, err := F64d6FromString(unquote(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// MarshalJSON implements json.Marshaler.
func (f F64d6) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *F64d6) UnmarshalJSON(in []byte) error {
	v, err := F64d6FromString(unquote(in))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f F64d6) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: f.String(),
	}, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *F64d6) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := F64d6FromString(str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}
