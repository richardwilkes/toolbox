// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xmath_test

import (
	"math"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xmath"
)

func TestAbsInt(t *testing.T) {
	c := check.New(t)

	// Ordinary values, including zero, across signed widths.
	c.Equal(5, xmath.AbsInt(5))
	c.Equal(5, xmath.AbsInt(-5))
	c.Equal(0, xmath.AbsInt(0))
	c.Equal(int8(7), xmath.AbsInt(int8(-7)))
	c.Equal(int16(7), xmath.AbsInt(int16(-7)))
	c.Equal(int32(7), xmath.AbsInt(int32(-7)))
	c.Equal(int64(7), xmath.AbsInt(int64(-7)))

	// The most-negative value of each type has no positive counterpart, so it saturates to the maximum rather than
	// overflowing back to a negative result.
	c.Equal(int8(math.MaxInt8), xmath.AbsInt(int8(math.MinInt8)))
	c.Equal(int16(math.MaxInt16), xmath.AbsInt(int16(math.MinInt16)))
	c.Equal(int32(math.MaxInt32), xmath.AbsInt(int32(math.MinInt32)))
	c.Equal(int64(math.MaxInt64), xmath.AbsInt(int64(math.MinInt64)))
	c.Equal(math.MaxInt, xmath.AbsInt(math.MinInt))

	// The value just above the minimum still negates normally.
	c.Equal(int8(math.MaxInt8), xmath.AbsInt(int8(math.MinInt8+1)))
}

func TestGCD(t *testing.T) {
	c := check.New(t)
	c.Equal(6, xmath.GCD(54, 24))
	c.Equal(1, xmath.GCD(17, 13))
	c.Equal(5, xmath.GCD(-5, 0))
	c.Equal(5, xmath.GCD(0, -5))
	c.Equal(4, xmath.GCD(-8, -12))
	c.Equal(0, xmath.GCD(0, 0))

	// The most-negative value must not overflow during negation. Its magnitude 2^63 shares the factor 2 with 6, so the
	// GCD is 2; previously the unfixed negation left a negative and yielded -2.
	c.Equal(2, xmath.GCD(math.MinInt, 6))
	c.Equal(2, xmath.GCD(6, math.MinInt))
	c.Equal(1, xmath.GCD(math.MinInt, 3))
	c.Equal(1, xmath.GCD(math.MinInt, math.MaxInt)) // GCD(2^63, 2^63-1): consecutive integers are coprime.

	// When the only non-zero operands are MinInt, the true GCD is MinInt's magnitude (2^63), which is not representable
	// as an int, so it saturates to the maximum int value rather than returning a negative result.
	c.Equal(math.MaxInt, xmath.GCD(math.MinInt, 0))
	c.Equal(math.MaxInt, xmath.GCD(0, math.MinInt))
	c.Equal(math.MaxInt, xmath.GCD(math.MinInt, math.MinInt))
}
