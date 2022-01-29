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
	// F128d16Max holds the maximum F128d16 value.
	F128d16Max = F128d16{data: num.MaxInt128}
	// F128d16Min holds the minimum F128d16 value.
	F128d16Min                = F128d16{data: num.MinInt128}
	multiplierF128d16BigInt   = new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil)
	multiplierF128d16BigFloat = new(big.Float).SetPrec(128).SetInt(multiplierF128d16BigInt)
	multiplierF128d16         = num.Int128FromBigInt(multiplierF128d16BigInt)
)

// F128d16 holds a fixed-point value that contains up to 16 decimal places. Values are truncated, not rounded.
type F128d16 struct {
	data num.Int128
}

// F128d16FromFloat64 creates a new F128d16 value from a float64.
func F128d16FromFloat64(value float64) F128d16 {
	f, _ := F128d16FromString(new(big.Float).SetPrec(128).SetFloat64(value).Text('f', 17)) //nolint:errcheck // Failure means 0
	return f
}

// F128d16FromInt64 creates a new F128d16 value from an int64.
func F128d16FromInt64(value int64) F128d16 {
	return F128d16{data: num.Int128From64(value).Mul(multiplierF128d16)}
}

// F128d16FromString creates a new F128d16 value from a string.
func F128d16FromString(str string) (F128d16, error) {
	if str == "" {
		return F128d16{}, errs.New("empty string is not valid")
	}
	if strings.ContainsAny(str, "Ee") {
		// Given a floating-point value with an exponent, which technically
		// isn't valid input, but we'll try to convert it anyway.
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return F128d16{}, err
		}
		return F128d16FromFloat64(f), nil
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
			return F128d16{}, errs.Newf("invalid value: %s", str)
		}
		if value.Sign() < 0 {
			neg = true
			value.Neg(value)
		}
		value.Mul(value, multiplierF128d16BigInt)
	}
	if len(parts) > 1 {
		var buffer strings.Builder
		buffer.WriteString("1")
		buffer.WriteString(parts[1])
		for buffer.Len() < 16+1 {
			buffer.WriteString("0")
		}
		frac := buffer.String()
		if len(frac) > 16+1 {
			frac = frac[:16+1]
		}
		if _, ok := fraction.SetString(frac, 10); !ok {
			return F128d16{}, errs.Newf("invalid value: %s", str)
		}
		value.Add(value, fraction).Sub(value, multiplierF128d16BigInt)
	}
	if neg {
		value.Neg(value)
	}
	return F128d16{data: num.Int128FromBigInt(value)}, nil
}

// F128d16FromStringForced creates a new F128d16 value from a string.
func F128d16FromStringForced(str string) F128d16 {
	f, _ := F128d16FromString(str) //nolint:errcheck // failure results in 0, which is acceptable here
	return f
}

// Add adds this value to the passed-in value, returning a new value.
func (f F128d16) Add(value F128d16) F128d16 {
	return F128d16{data: f.data.Add(value.data)}
}

// Sub subtracts the passed-in value from this value, returning a new value.
func (f F128d16) Sub(value F128d16) F128d16 {
	return F128d16{data: f.data.Sub(value.data)}
}

// Mul multiplies this value by the passed-in value, returning a new value.
func (f F128d16) Mul(value F128d16) F128d16 {
	return F128d16{data: f.data.Mul(value.data).Div(multiplierF128d16)}
}

// Div divides this value by the passed-in value, returning a new value.
func (f F128d16) Div(value F128d16) F128d16 {
	return F128d16{data: f.data.Mul(multiplierF128d16).Div(value.data)}
}

// Trunc returns a new value which has everything to the right of the decimal
// place truncated.
func (f F128d16) Trunc() F128d16 {
	return F128d16{data: f.data.Div(multiplierF128d16).Mul(multiplierF128d16)}
}

// AsInt64 returns the truncated equivalent integer to this value.
func (f F128d16) AsInt64() int64 {
	return f.data.Div(multiplierF128d16).AsInt64()
}

// AsFloat64 returns the floating-point equivalent to this value.
func (f F128d16) AsFloat64() float64 {
	f64, _ := new(big.Float).SetPrec(128).Quo(f.data.AsBigFloat(), multiplierF128d16BigFloat).Float64()
	return f64
}

// CommaWithSign returns the same as Comma(), but prefixes the value with a '+' if it is positive
func (f F128d16) CommaWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.Comma()
	}
	return f.Comma()
}

// Comma returns the same as String(), but with commas for values of 1000 and
// greater.
func (f F128d16) Comma() string {
	var iStr string
	integer := f.data.Div(multiplierF128d16)
	if integer.IsInt64() {
		iStr = humanize.Comma(integer.AsInt64())
	} else {
		iStr = humanize.BigComma(integer.AsBigInt())
	}
	fraction := f.data.Sub(integer.Mul(multiplierF128d16))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplierF128d16).String()
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
func (f F128d16) StringWithSign() string {
	if f.data.Sign() >= 0 {
		return "+" + f.String()
	}
	return f.String()
}

func (f F128d16) String() string {
	integer := f.data.Div(multiplierF128d16)
	iStr := integer.String()
	fraction := f.data.Sub(integer.Mul(multiplierF128d16))
	if fraction.IsZero() {
		return iStr
	}
	if fraction.Sign() < 0 {
		fraction = fraction.Neg()
	}
	fStr := fraction.Add(multiplierF128d16).String()
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
func (f F128d16) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (f *F128d16) UnmarshalText(text []byte) error {
	f1, err := F128d16FromString(string(text))
	if err != nil {
		return err
	}
	*f = f1
	return nil
}

// Float64 implements json.Number. Intentionally returns an error if the value
// cannot be represented exactly with a float64, as we never want to emit
// inexact floating point values into json for fixed-point values.
func (f F128d16) Float64() (float64, error) {
	n := f.AsFloat64()
	if strconv.FormatFloat(n, 'g', -1, 64) != f.String() {
		return 0, errDoesNotFitInFloat64
	}
	return n, nil
}

// Int64 implements json.Number. Intentionally returns an error if the value
// cannot be represented exactly with an int64, as we never want to emit
// inexact values into json for fixed-point values.
func (f F128d16) Int64() (int64, error) {
	n := f.AsInt64()
	if F128d16FromInt64(n) != f {
		return 0, errDoesNotFitInInt64
	}
	return n, nil
}

// MarshalJSON implements json.Marshaler.
func (f F128d16) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *F128d16) UnmarshalJSON(in []byte) error {
	v, err := F128d16FromString(string(in))
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f F128d16) MarshalYAML() (interface{}, error) {
	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: f.String(),
	}, nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *F128d16) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	v, err := F128d16FromString(str)
	if err != nil {
		return err
	}
	*f = v
	return nil
}
