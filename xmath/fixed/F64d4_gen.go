// Code created from "fixed64.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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
)

const (
	// F64d4Max holds the maximum F64d4 value.
	F64d4Max = F64d4(1<<63 - 1)
	// F64d4Min holds the minimum F64d4 value.
	F64d4Min = F64d4(^(1<<63 - 1))
)

var multiplierF64d4 = int64(math.Pow(10, 4))

// F64d4 holds a fixed-point value that contains up to 4 decimal places.
// Values are truncated, not rounded. Values can be added and subtracted
// directly. For multiplication and division, the provided Mul() and Div()
// methods should be used.
type F64d4 int64

// F64d4FromFloat64 creates a new F64d4 value from a float64.
func F64d4FromFloat64(value float64) F64d4 {
	return F64d4(value * float64(multiplierF64d4))
}

// F64d4FromInt64 creates a new F64d4 value from an int64.
func F64d4FromInt64(value int64) F64d4 {
	return F64d4(value * multiplierF64d4)
}

// F64d4FromString creates a new F64d4 value from a string.
func F64d4FromString(str string) (F64d4, error) {
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
		return F64d4FromFloat64(f), nil
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
		value *= multiplierF64d4
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 4+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 4+1 {
			frac = frac[:4+1]
		}
		if fraction, err = strconv.ParseInt(frac, 10, 64); err != nil {
			return 0, errs.Wrap(err)
		}
		value += fraction - multiplierF64d4
	}
	if neg {
		value = -value
	}
	return F64d4(value), nil
}

// F64d4FromStringForced creates a new F64d4 value from a string.
func F64d4FromStringForced(str string) F64d4 {
	f, _ := F64d4FromString(str) //nolint:errcheck
	return f
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F64d4) Mul(value F64d4) F64d4 {
	return f * value / F64d4(multiplierF64d4)
}

// Div divides this value by the passed-in value, returning a new value.
func (f F64d4) Div(value F64d4) F64d4 {
	return f * F64d4(multiplierF64d4) / value
}

// Trunc returns a new value which has everything to the right of the decimal
// place truncated.
func (f F64d4) Trunc() F64d4 {
	return f / F64d4(multiplierF64d4) * F64d4(multiplierF64d4)
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f F64d4) AsInt64() int64 {
	return int64(f / F64d4(multiplierF64d4))
}

// AsFloat64 returns the floating-point equivalent to this value.
func (f F64d4) AsFloat64() float64 {
	return float64(f) / float64(multiplierF64d4)
}

// Comma returns the same as String(), but with commas for values of 1000 and
// greater.
func (f F64d4) Comma() string {
	integer := f / F64d4(multiplierF64d4)
	fraction := f % F64d4(multiplierF64d4)
	if fraction == 0 {
		return humanize.Comma(int64(integer))
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d4(multiplierF64d4)
	fstr := strconv.FormatInt(int64(fraction), 10)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && f < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%s.%s", neg, humanize.Comma(int64(integer)), fstr)
}

func (f F64d4) String() string {
	integer := f / F64d4(multiplierF64d4)
	fraction := f % F64d4(multiplierF64d4)
	if fraction == 0 {
		return strconv.FormatInt(int64(integer), 10)
	}
	if fraction < 0 {
		fraction = -fraction
	}
	fraction += F64d4(multiplierF64d4)
	fstr := strconv.FormatInt(int64(fraction), 10)
	for i := len(fstr) - 1; i > 0; i-- {
		if fstr[i] != '0' {
			fstr = fstr[1 : i+1]
			break
		}
	}
	var neg string
	if integer == 0 && f < 0 {
		neg = "-"
	} else {
		neg = ""
	}
	return fmt.Sprintf("%s%d.%s", neg, integer, fstr)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (f F64d4) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F64d4) UnmarshalText(text []byte) error {
	f1, err := F64d4FromString(string(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// Float64 implements json.Number. Intentionally returns an error if the value
// cannot be represented exactly with a float64, as we never want to emit
// inexact floating point values into json for fixed-point values.
func (f F64d4) Float64() (float64, error) {
	n := f.AsFloat64()
	if strconv.FormatFloat(n, 'g', -1, 64) != f.String() {
		return 0, errDoesNotFitInFloat64
	}
	return n, nil
}

// Int64 implements json.Number. Intentionally returns an error if the value
// cannot be represented exactly with an int64, as we never want to emit
// inexact values into json for fixed-point values.
func (f F64d4) Int64() (int64, error) {
	n := f.AsInt64()
	if F64d4FromInt64(n) != f {
		return 0, errDoesNotFitInInt64
	}
	return f.AsInt64(), nil
}

// MarshalJSON implements json.Marshaler.
func (f F64d4) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *F64d4) UnmarshalJSON(in []byte) error {
	v, err := F64d4FromString(string(in))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f F64d4) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *F64d4) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := F64d4FromString(str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}
