// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

func TestNewSize(t *testing.T) {
	c := check.New(t)
	s := geom.NewSize(10, 20)
	c.Equal(float32(10), s.Width)
	c.Equal(float32(20), s.Height)
	sf := geom.NewSize(10.5, 20.7)
	c.Equal(float32(10.5), sf.Width)
	c.Equal(float32(20.7), sf.Height)
}

func TestSizeAdd(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 20)
	s2 := geom.NewSize(5, 8)
	result := s1.Add(s2)

	c.Equal(float32(15), result.Width)
	c.Equal(float32(28), result.Height)

	// Original sizes should be unchanged
	c.Equal(float32(10), s1.Width)
	c.Equal(float32(20), s1.Height)
	c.Equal(float32(5), s2.Width)
	c.Equal(float32(8), s2.Height)
}

func TestSizeSub(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(15, 25)
	s2 := geom.NewSize(5, 8)
	result := s1.Sub(s2)

	c.Equal(float32(10), result.Width)
	c.Equal(float32(17), result.Height)

	// Original sizes should be unchanged
	c.Equal(float32(15), s1.Width)
	c.Equal(float32(25), s1.Height)
}

func TestSizeMul(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10, 20)
	result := s.Mul(2)

	c.Equal(float32(20), result.Width)
	c.Equal(float32(40), result.Height)

	// Original size should be unchanged
	c.Equal(float32(10), s.Width)
	c.Equal(float32(20), s.Height)
}

func TestSizeDiv(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(20, 40)
	result := s.Div(2)

	c.Equal(float32(10), result.Width)
	c.Equal(float32(20), result.Height)

	sf := geom.NewSize(15.0, 25.0)
	resultf := sf.Div(2.0)

	c.Equal(float32(7.5), resultf.Width)
	c.Equal(float32(12.5), resultf.Height)
}

func TestSizeFloor(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10.7, 20.2)
	result := s.Floor()

	c.Equal(float32(10), result.Width)
	c.Equal(float32(20), result.Height)

	// Test with negative values
	s2 := geom.NewSize(-10.7, -20.2)
	result2 := s2.Floor()

	c.Equal(float32(-11), result2.Width)
	c.Equal(float32(-21), result2.Height)
}

func TestSizeCeil(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10.7, 20.2)
	result := s.Ceil()

	c.Equal(float32(11), result.Width)
	c.Equal(float32(21), result.Height)

	// Test with negative values
	s2 := geom.NewSize(-10.7, -20.2)
	result2 := s2.Ceil()

	c.Equal(float32(-10), result2.Width)
	c.Equal(float32(-20), result2.Height)
}

func TestSizeMin(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 25)
	s2 := geom.NewSize(15, 20)
	result := s1.Min(s2)

	c.Equal(float32(10), result.Width)  // min(10, 15)
	c.Equal(float32(20), result.Height) // min(25, 20)

	// Test with equal values
	s3 := geom.NewSize(10, 20)
	s4 := geom.NewSize(10, 20)
	result2 := s3.Min(s4)

	c.Equal(float32(10), result2.Width)
	c.Equal(float32(20), result2.Height)
}

func TestSizeMax(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 25)
	s2 := geom.NewSize(15, 20)
	result := s1.Max(s2)

	c.Equal(float32(15), result.Width)  // max(10, 15)
	c.Equal(float32(25), result.Height) // max(25, 20)

	// Test with equal values
	s3 := geom.NewSize(10, 20)
	s4 := geom.NewSize(10, 20)
	result2 := s3.Max(s4)

	c.Equal(float32(10), result2.Width)
	c.Equal(float32(20), result2.Height)
}

func TestSizeConstrainForHint(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(100, 200)

	// Hint is larger, should not constrain
	hint1 := geom.NewSize(150, 250)
	result1 := s.ConstrainForHint(hint1)
	c.Equal(float32(100), result1.Width)
	c.Equal(float32(200), result1.Height)

	// Hint is smaller, should constrain
	hint2 := geom.NewSize(80, 180)
	result2 := s.ConstrainForHint(hint2)
	c.Equal(float32(80), result2.Width)
	c.Equal(float32(180), result2.Height)

	// Mixed constraints
	hint3 := geom.NewSize(80, 250)
	result3 := s.ConstrainForHint(hint3)
	c.Equal(float32(80), result3.Width)   // constrained
	c.Equal(float32(200), result3.Height) // not constrained

	// Hint values less than 1 should be ignored
	hint4 := geom.NewSize(0, 0)
	result4 := s.ConstrainForHint(hint4)
	c.Equal(float32(100), result4.Width)  // not constrained (hint < 1)
	c.Equal(float32(200), result4.Height) // not constrained (hint < 1)

	// Edge case: hint exactly 1
	hint5 := geom.NewSize(1, 1)
	result5 := s.ConstrainForHint(hint5)
	c.Equal(float32(1), result5.Width)  // constrained
	c.Equal(float32(1), result5.Height) // constrained
}

func TestSizeString(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10, 20)
	str := s.String()
	c.Equal("10,20", str)

	sf := geom.NewSize(10.5, 20.7)
	strf := sf.String()
	c.Equal("10.5,20.7", strf)
}
