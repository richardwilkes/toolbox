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

	// Test int size
	s := geom.NewSize(10, 20)
	c.Equal(10, s.Width)
	c.Equal(20, s.Height)

	// Test float size
	sf := geom.NewSize(10.5, 20.7)
	c.Equal(10.5, sf.Width)
	c.Equal(20.7, sf.Height)
}

func TestConvertSize(t *testing.T) {
	c := check.New(t)

	// Convert from int to float
	intSize := geom.NewSize(10, 20)
	floatSize := geom.ConvertSize[float64](intSize)
	c.Equal(10.0, floatSize.Width)
	c.Equal(20.0, floatSize.Height)

	// Convert from float to int
	floatSize2 := geom.NewSize(10.7, 20.9)
	intSize2 := geom.ConvertSize[int](floatSize2)
	c.Equal(10, intSize2.Width)
	c.Equal(20, intSize2.Height)
}

func TestSizeAdd(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 20)
	s2 := geom.NewSize(5, 8)
	result := s1.Add(s2)

	c.Equal(15, result.Width)
	c.Equal(28, result.Height)

	// Original sizes should be unchanged
	c.Equal(10, s1.Width)
	c.Equal(20, s1.Height)
	c.Equal(5, s2.Width)
	c.Equal(8, s2.Height)
}

func TestSizeSub(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(15, 25)
	s2 := geom.NewSize(5, 8)
	result := s1.Sub(s2)

	c.Equal(10, result.Width)
	c.Equal(17, result.Height)

	// Original sizes should be unchanged
	c.Equal(15, s1.Width)
	c.Equal(25, s1.Height)
}

func TestSizeMul(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10, 20)
	result := s.Mul(2)

	c.Equal(20, result.Width)
	c.Equal(40, result.Height)

	// Original size should be unchanged
	c.Equal(10, s.Width)
	c.Equal(20, s.Height)
}

func TestSizeDiv(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(20, 40)
	result := s.Div(2)

	c.Equal(10, result.Width)
	c.Equal(20, result.Height)

	// Test with float
	sf := geom.NewSize(15.0, 25.0)
	resultf := sf.Div(2.0)

	c.Equal(7.5, resultf.Width)
	c.Equal(12.5, resultf.Height)
}

func TestSizeFloor(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10.7, 20.2)
	result := s.Floor()

	c.Equal(10.0, result.Width)
	c.Equal(20.0, result.Height)

	// Test with negative values
	s2 := geom.NewSize(-10.7, -20.2)
	result2 := s2.Floor()

	c.Equal(-11.0, result2.Width)
	c.Equal(-21.0, result2.Height)
}

func TestSizeCeil(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10.7, 20.2)
	result := s.Ceil()

	c.Equal(11.0, result.Width)
	c.Equal(21.0, result.Height)

	// Test with negative values
	s2 := geom.NewSize(-10.7, -20.2)
	result2 := s2.Ceil()

	c.Equal(-10.0, result2.Width)
	c.Equal(-20.0, result2.Height)
}

func TestSizeMin(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 25)
	s2 := geom.NewSize(15, 20)
	result := s1.Min(s2)

	c.Equal(10, result.Width)  // min(10, 15)
	c.Equal(20, result.Height) // min(25, 20)

	// Test with equal values
	s3 := geom.NewSize(10, 20)
	s4 := geom.NewSize(10, 20)
	result2 := s3.Min(s4)

	c.Equal(10, result2.Width)
	c.Equal(20, result2.Height)
}

func TestSizeMax(t *testing.T) {
	c := check.New(t)

	s1 := geom.NewSize(10, 25)
	s2 := geom.NewSize(15, 20)
	result := s1.Max(s2)

	c.Equal(15, result.Width)  // max(10, 15)
	c.Equal(25, result.Height) // max(25, 20)

	// Test with equal values
	s3 := geom.NewSize(10, 20)
	s4 := geom.NewSize(10, 20)
	result2 := s3.Max(s4)

	c.Equal(10, result2.Width)
	c.Equal(20, result2.Height)
}

func TestSizeConstrainForHint(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(100, 200)

	// Hint is larger, should not constrain
	hint1 := geom.NewSize(150, 250)
	result1 := s.ConstrainForHint(hint1)
	c.Equal(100, result1.Width)
	c.Equal(200, result1.Height)

	// Hint is smaller, should constrain
	hint2 := geom.NewSize(80, 180)
	result2 := s.ConstrainForHint(hint2)
	c.Equal(80, result2.Width)
	c.Equal(180, result2.Height)

	// Mixed constraints
	hint3 := geom.NewSize(80, 250)
	result3 := s.ConstrainForHint(hint3)
	c.Equal(80, result3.Width)   // constrained
	c.Equal(200, result3.Height) // not constrained

	// Hint values less than 1 should be ignored
	hint4 := geom.NewSize(0, 0)
	result4 := s.ConstrainForHint(hint4)
	c.Equal(100, result4.Width)  // not constrained (hint < 1)
	c.Equal(200, result4.Height) // not constrained (hint < 1)

	// Edge case: hint exactly 1
	hint5 := geom.NewSize(1, 1)
	result5 := s.ConstrainForHint(hint5)
	c.Equal(1, result5.Width)  // constrained
	c.Equal(1, result5.Height) // constrained
}

func TestSizeString(t *testing.T) {
	c := check.New(t)

	s := geom.NewSize(10, 20)
	str := s.String()
	c.Equal("10,20", str)

	// Test with float
	sf := geom.NewSize(10.5, 20.7)
	strf := sf.String()
	c.Equal("10.5,20.7", strf)
}
