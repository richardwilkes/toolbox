// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xmath provides math-related utilities.
package xmath

import (
	"math"
)

// Mathematical constants. Mostly just re-exported from the math package for convenience.
const (
	E   = math.E
	Pi  = math.Pi
	Phi = math.Phi

	Sqrt2   = math.Sqrt2
	SqrtE   = math.SqrtE
	SqrtPi  = math.SqrtPi
	SqrtPhi = math.SqrtPhi

	Ln2    = math.Ln2
	Log2E  = 1 / Ln2
	Ln10   = math.Ln10
	Log10E = 1 / Ln10
)

// Floating-point limit values. Mostly just re-exported from the math package for convenience. Max is the largest finite
// value representable by the type. SmallestNonzero is the smallest positive, non-zero value representable by the type.
const (
	MaxFloat32             = math.MaxFloat32
	SmallestNonzeroFloat32 = math.SmallestNonzeroFloat32
	MaxFloat64             = math.MaxFloat64
	SmallestNonzeroFloat64 = math.SmallestNonzeroFloat64
)

// Integer limit values. Mostly just re-exported from the math package for convenience.
const (
	MaxInt    = math.MaxInt
	MinInt    = math.MinInt
	MaxInt8   = math.MaxInt8
	MinInt8   = math.MinInt8
	MaxInt16  = math.MaxInt16
	MinInt16  = math.MinInt16
	MaxInt32  = math.MaxInt32
	MinInt32  = math.MinInt32
	MaxInt64  = math.MaxInt64
	MinInt64  = math.MinInt64
	MaxUint   = math.MaxUint
	MaxUint8  = math.MaxUint8
	MaxUint16 = math.MaxUint16
	MaxUint32 = math.MaxUint32
	MaxUint64 = math.MaxUint64
)

const (
	// DegreesToRadians converts a value in degrees to radians when multiplied with the value.
	DegreesToRadians = math.Pi / 180
	// RadiansToDegrees converts a value in radians to degrees when multiplied with the value.
	RadiansToDegrees = 180 / math.Pi
)
