// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/fixed/fixed64"
	"golang.org/x/exp/constraints"
)

// Common values that can be reused
var (
	Min               = Num(math.MinInt64)
	PointTwo          = NumFromFloat(0.2)
	One               = NumFromInteger(1)
	Two               = NumFromInteger(2)
	Four              = NumFromInteger(4)
	Ten               = NumFromInteger(10)
	OneHundredEighty  = NumFromInteger(180)
	ThreeHundredSixty = NumFromInteger(360)
	Pi                = NumFromFloat(math.Pi)
	DegreesToRadians  = NumFromFloat(math.Pi / 180)
	RadiansToDegrees  = NumFromFloat(180 / math.Pi)
)

// Num is an alias for the fixed-point type we are using.
type Num = fixed64.Int[fixed.D6]

// NumFromInteger creates a Num from a numeric value.
func NumFromInteger[T constraints.Integer](value T) Num {
	return fixed64.FromInteger[fixed.D6](value)
}

// NumFromFloat creates a Num from a numeric value.
func NumFromFloat[T constraints.Float](value T) Num {
	return fixed64.FromFloat[fixed.D6](value)
}

// NumFromString creates a Num from a string.
func NumFromString(value string) (Num, error) {
	return fixed64.FromString[fixed.D6](value)
}

// NumFromStringForced creates a Num from a string, ignoring any conversion inaccuracies.
func NumFromStringForced(value string) Num {
	return fixed64.FromStringForced[fixed.D6](value)
}

// NumAsInteger returns the equivalent value in the destination type.
func NumAsInteger[T constraints.Integer](value Num) T {
	return fixed64.AsInteger[fixed.D6, T](value)
}

// NumAsFloat returns the equivalent value in the destination type.
func NumAsFloat[T constraints.Float](value Num) T {
	return fixed64.AsFloat[fixed.D6, T](value)
}

// Sqrt returns the square root of the value.
func Sqrt(value Num) Num {
	return NumFromFloat(math.Sqrt(NumAsFloat[float64](value)))
}

// Sin returns the sine of the value.
func Sin(value Num) Num {
	return NumFromFloat(math.Sin(NumAsFloat[float64](value)))
}

// Cos returns the cosine of the value.
func Cos(value Num) Num {
	return NumFromFloat(math.Cos(NumAsFloat[float64](value)))
}

// Acos returns the arc cosine of the value.
func Acos(value Num) Num {
	return NumFromFloat(math.Acos(NumAsFloat[float64](value)))
}

// Atan2 returns the arc tangent of the quotient of its arguments.
func Atan2(y, x Num) Num {
	return NumFromFloat(math.Atan2(NumAsFloat[float64](y), NumAsFloat[float64](x)))
}
