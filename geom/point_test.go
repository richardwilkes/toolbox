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

func TestNewPoint(t *testing.T) {
	c := check.New(t)

	// Test int point
	p := geom.NewPoint(3, 4)
	c.Equal(float32(3), p.X)
	c.Equal(float32(4), p.Y)

	// Test float point
	pf := geom.NewPoint(3.5, 4.7)
	c.Equal(float32(3.5), pf.X)
	c.Equal(float32(4.7), pf.Y)
}

func TestPointAdd(t *testing.T) {
	c := check.New(t)

	p1 := geom.NewPoint(3, 4)
	p2 := geom.NewPoint(1, 2)
	result := p1.Add(p2)

	c.Equal(float32(4), result.X)
	c.Equal(float32(6), result.Y)

	// Original points should be unchanged
	c.Equal(float32(3), p1.X)
	c.Equal(float32(4), p1.Y)
	c.Equal(float32(1), p2.X)
	c.Equal(float32(2), p2.Y)
}

func TestPointSub(t *testing.T) {
	c := check.New(t)

	p1 := geom.NewPoint(5, 7)
	p2 := geom.NewPoint(2, 3)
	result := p1.Sub(p2)

	c.Equal(float32(3), result.X)
	c.Equal(float32(4), result.Y)

	// Original points should be unchanged
	c.Equal(float32(5), p1.X)
	c.Equal(float32(7), p1.Y)
}

func TestPointMul(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(3, 4)
	result := p.Mul(2)

	c.Equal(float32(6), result.X)
	c.Equal(float32(8), result.Y)

	// Original point should be unchanged
	c.Equal(float32(3), p.X)
	c.Equal(float32(4), p.Y)
}

func TestPointDiv(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(6, 8)
	result := p.Div(2)

	c.Equal(float32(3), result.X)
	c.Equal(float32(4), result.Y)

	// Test with float
	pf := geom.NewPoint(7.0, 9.0)
	resultf := pf.Div(2.0)

	c.Equal(float32(3.5), resultf.X)
	c.Equal(float32(4.5), resultf.Y)
}

func TestPointNeg(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(3, -4)
	result := p.Neg()

	c.Equal(float32(-3), result.X)
	c.Equal(float32(4), result.Y)

	// Original point should be unchanged
	c.Equal(float32(3), p.X)
	c.Equal(float32(-4), p.Y)
}

func TestPointFloor(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(3.7, 4.2)
	result := p.Floor()

	c.Equal(float32(3), result.X)
	c.Equal(float32(4), result.Y)

	// Test with negative values
	p2 := geom.NewPoint(-3.7, -4.2)
	result2 := p2.Floor()

	c.Equal(float32(-4), result2.X)
	c.Equal(float32(-5), result2.Y)
}

func TestPointCeil(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(3.7, 4.2)
	result := p.Ceil()

	c.Equal(float32(4), result.X)
	c.Equal(float32(5), result.Y)

	// Test with negative values
	p2 := geom.NewPoint(-3.7, -4.2)
	result2 := p2.Ceil()

	c.Equal(float32(-3), result2.X)
	c.Equal(float32(-4), result2.Y)
}

func TestPointDot(t *testing.T) {
	c := check.New(t)

	p1 := geom.NewPoint(3, 4)
	p2 := geom.NewPoint(2, 1)

	// Dot product: (3*2) + (4*1) = 6 + 4 = 10
	result := p1.Dot(p2)
	c.Equal(float32(10), result)

	// Test commutativity
	result2 := p2.Dot(p1)
	c.Equal(float32(10), result2)

	// Test with perpendicular vectors (should be 0)
	p3 := geom.NewPoint(1, 0)
	p4 := geom.NewPoint(0, 1)
	result3 := p3.Dot(p4)
	c.Equal(float32(0), result3)
}

func TestPointCross(t *testing.T) {
	c := check.New(t)

	p1 := geom.NewPoint(3, 4)
	p2 := geom.NewPoint(2, 1)

	// Cross product: (3*1) - (4*2) = 3 - 8 = -5
	result := p1.Cross(p2)
	c.Equal(float32(-5), result)

	// Test anti-commutativity
	result2 := p2.Cross(p1)
	c.Equal(float32(5), result2)

	// Test with parallel vectors (should be 0)
	p3 := geom.NewPoint(2, 4)
	p4 := geom.NewPoint(1, 2)
	result3 := p3.Cross(p4)
	c.Equal(float32(0), result3)
}

func TestPointIn(t *testing.T) {
	c := check.New(t)

	rect := geom.NewRect(10, 20, 30, 40)

	// Point inside the rectangle
	p1 := geom.NewPoint(15, 25)
	c.True(p1.In(rect))

	// Point on the left edge (should be inside)
	p2 := geom.NewPoint(10, 25)
	c.True(p2.In(rect))

	// Point on the top edge (should be inside)
	p3 := geom.NewPoint(15, 20)
	c.True(p3.In(rect))

	// Point on the right edge (should be outside)
	p4 := geom.NewPoint(40, 25)
	c.False(p4.In(rect))

	// Point on the bottom edge (should be outside)
	p5 := geom.NewPoint(15, 60)
	c.False(p5.In(rect))

	// Point completely outside
	p6 := geom.NewPoint(5, 15)
	c.False(p6.In(rect))

	// Test with empty rectangle
	emptyRect := geom.NewRect(10, 20, 0, 40)
	p7 := geom.NewPoint(15, 25)
	c.False(p7.In(emptyRect))
}

func TestPointEqualWithin(t *testing.T) {
	c := check.New(t)

	p1 := geom.NewPoint(3.0, 4.0)
	p2 := geom.NewPoint(3.1, 4.05)

	// Within tolerance
	c.True(p1.EqualWithin(p2, 0.2))

	// Outside tolerance
	c.False(p1.EqualWithin(p2, 0.01))

	// Exact match
	p3 := geom.NewPoint(3.0, 4.0)
	c.True(p1.EqualWithin(p3, 0.0))

	// Test with negative coordinates
	p4 := geom.NewPoint(-3.0, -4.0)
	p5 := geom.NewPoint(-3.05, -4.02)
	c.True(p4.EqualWithin(p5, 0.1))
}

func TestPointString(t *testing.T) {
	c := check.New(t)

	p := geom.NewPoint(3, 4)
	str := p.String()
	c.Equal("3,4", str)

	// Test with float
	pf := geom.NewPoint(3.5, 4.7)
	strf := pf.String()
	c.Equal("3.5,4.7", strf)
}
