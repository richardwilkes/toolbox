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

func TestNewRect(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	c.Equal(10, r.X)
	c.Equal(20, r.Y)
	c.Equal(30, r.Width)
	c.Equal(40, r.Height)
}

func TestConvertRect(t *testing.T) {
	c := check.New(t)

	// Convert from int to float
	intRect := geom.NewRect(10, 20, 30, 40)
	floatRect := geom.ConvertRect[float64](intRect)
	c.Equal(10.0, floatRect.X)
	c.Equal(20.0, floatRect.Y)
	c.Equal(30.0, floatRect.Width)
	c.Equal(40.0, floatRect.Height)

	// Convert from float to int
	floatRect2 := geom.NewRect(10.7, 20.9, 30.3, 40.8)
	intRect2 := geom.ConvertRect[int](floatRect2)
	c.Equal(10, intRect2.X)
	c.Equal(20, intRect2.Y)
	c.Equal(30, intRect2.Width)
	c.Equal(40, intRect2.Height)
}

func TestRectEmpty(t *testing.T) {
	c := check.New(t)

	// Non-empty rect
	r1 := geom.NewRect(10, 20, 30, 40)
	c.False(r1.Empty())

	// Zero width
	r2 := geom.NewRect(10, 20, 0, 40)
	c.True(r2.Empty())

	// Zero height
	r3 := geom.NewRect(10, 20, 30, 0)
	c.True(r3.Empty())

	// Negative width
	r4 := geom.NewRect(10, 20, -30, 40)
	c.True(r4.Empty())

	// Negative height
	r5 := geom.NewRect(10, 20, 30, -40)
	c.True(r5.Empty())
}

func TestRectCenter(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	center := r.Center()

	// Center X: 10 + 30/2 = 25
	// Center Y: 20 + 40/2 = 40
	c.Equal(25, center.X)
	c.Equal(40, center.Y)
}

func TestRectCenterX(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	centerX := r.CenterX()

	// Center X: 10 + 30/2 = 25
	c.Equal(25, centerX)
}

func TestRectCenterY(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	centerY := r.CenterY()

	// Center Y: 20 + 40/2 = 40
	c.Equal(40, centerY)
}

//nolint:gocritic // The "commented out code" is actually explanation
func TestRectRight(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	right := r.Right()

	// Right: 10 + 30 = 40
	c.Equal(40, right)
}

//nolint:gocritic // The "commented out code" is actually explanation
func TestRectBottom(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	bottom := r.Bottom()

	// Bottom: 20 + 40 = 60
	c.Equal(60, bottom)
}

func TestRectCornerPoints(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)

	// TopLeft should be (10, 20)
	topLeft := r.TopLeft()
	c.Equal(10, topLeft.X)
	c.Equal(20, topLeft.Y)

	// TopRight should be (40, 20)
	topRight := r.TopRight()
	c.Equal(40, topRight.X)
	c.Equal(20, topRight.Y)

	// BottomRight should be (40, 60)
	bottomRight := r.BottomRight()
	c.Equal(40, bottomRight.X)
	c.Equal(60, bottomRight.Y)

	// BottomLeft should be (10, 60)
	bottomLeft := r.BottomLeft()
	c.Equal(10, bottomLeft.X)
	c.Equal(60, bottomLeft.Y)
}

func TestRectContains(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)

	// Rectangle completely inside
	inner := geom.NewRect(15, 25, 10, 15)
	c.True(r.Contains(inner))

	// Rectangle touching edges from inside (should be contained)
	edge := geom.NewRect(10, 20, 15, 20)
	c.True(r.Contains(edge))

	// Rectangle extending outside
	outside := geom.NewRect(5, 15, 50, 60)
	c.False(r.Contains(outside))

	// Rectangle partially overlapping
	overlap := geom.NewRect(35, 55, 20, 20)
	c.False(r.Contains(overlap))

	// Empty rectangles
	empty1 := geom.NewRect(15, 25, 0, 15)
	c.False(r.Contains(empty1))

	emptyR := geom.NewRect(10, 20, 0, 40)
	normalInner := geom.NewRect(15, 25, 10, 15)
	c.False(emptyR.Contains(normalInner))
}

func TestRectIntersects(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)

	// Rectangle completely inside
	inner := geom.NewRect(15, 25, 10, 15)
	c.True(r.Intersects(inner))

	// Rectangle partially overlapping
	overlap := geom.NewRect(35, 55, 20, 20)
	c.True(r.Intersects(overlap))

	// Rectangle touching edge
	touching := geom.NewRect(40, 30, 10, 20)
	c.False(r.Intersects(touching)) // Just touching, not intersecting

	// Rectangle completely outside
	outside := geom.NewRect(50, 70, 10, 20)
	c.False(r.Intersects(outside))

	// Empty rectangles
	empty := geom.NewRect(15, 25, 0, 15)
	c.False(r.Intersects(empty))

	emptyR := geom.NewRect(10, 20, 0, 40)
	normal := geom.NewRect(15, 25, 10, 15)
	c.False(emptyR.Intersects(normal))
}

func TestRectIntersect(t *testing.T) {
	c := check.New(t)

	r1 := geom.NewRect(10, 20, 30, 40)
	r2 := geom.NewRect(25, 35, 30, 40)

	intersection := r1.Intersect(r2)

	// Intersection should be (25, 35, 15, 25)
	c.Equal(25, intersection.X)
	c.Equal(35, intersection.Y)
	c.Equal(15, intersection.Width)
	c.Equal(25, intersection.Height)

	// Non-intersecting rectangles
	r3 := geom.NewRect(50, 70, 10, 20)
	noIntersection := r1.Intersect(r3)
	c.True(noIntersection.Empty())

	// Empty rectangle intersection
	emptyR := geom.NewRect(10, 20, 0, 40)
	emptyIntersection := r1.Intersect(emptyR)
	c.True(emptyIntersection.Empty())
}

func TestRectUnion(t *testing.T) {
	c := check.New(t)

	r1 := geom.NewRect(10, 20, 30, 40)
	r2 := geom.NewRect(25, 35, 30, 40)

	union := r1.Union(r2)

	// Union should encompass both rectangles
	// Left: min(10, 25) = 10
	// Top: min(20, 35) = 20
	// Right: max(40, 55) = 55
	// Bottom: max(60, 75) = 75
	// So: (10, 20, 45, 55)
	c.Equal(10, union.X)
	c.Equal(20, union.Y)
	c.Equal(45, union.Width)
	c.Equal(55, union.Height)

	// Union with empty rectangle
	emptyR := geom.NewRect(10, 20, 0, 40)
	unionWithEmpty := r1.Union(emptyR)
	c.Equal(r1.X, unionWithEmpty.X)
	c.Equal(r1.Y, unionWithEmpty.Y)
	c.Equal(r1.Width, unionWithEmpty.Width)
	c.Equal(r1.Height, unionWithEmpty.Height)

	// Union of two empty rectangles
	empty1 := geom.NewRect(10, 20, 0, 40)
	empty2 := geom.NewRect(15, 25, -10, 20)
	emptyUnion := empty1.Union(empty2)
	c.True(emptyUnion.Empty())
}

func TestRectAlign(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10.3, 20.7, 30.2, 40.8)
	aligned := r.Align()

	// Point should be floored, size should be ceiled
	c.Equal(10.0, aligned.X)
	c.Equal(20.0, aligned.Y)
	c.Equal(31.0, aligned.Width)
	c.Equal(41.0, aligned.Height)
}

func TestRectExpand(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)

	// Point inside rectangle
	pt1 := geom.NewPoint(15, 25)
	expanded1 := r.Expand(pt1)
	c.Equal(r.X, expanded1.X)
	c.Equal(r.Y, expanded1.Y)
	c.Equal(r.Width, expanded1.Width)
	c.Equal(r.Height, expanded1.Height)

	// Point outside rectangle (to the left and up)
	pt2 := geom.NewPoint(5, 15)
	expanded2 := r.Expand(pt2)
	c.Equal(5, expanded2.X)
	c.Equal(15, expanded2.Y)
	c.Equal(35, expanded2.Width)  // 40 - 5 = 35
	c.Equal(45, expanded2.Height) // 60 - 15 = 45

	// Point outside rectangle (to the right and down)
	pt3 := geom.NewPoint(50, 70)
	expanded3 := r.Expand(pt3)
	c.Equal(10, expanded3.X)
	c.Equal(20, expanded3.Y)
	c.Equal(40, expanded3.Width)  // 50 - 10 = 40
	c.Equal(50, expanded3.Height) // 70 - 20 = 50

	// Rectangle with negative width
	negativeR := geom.NewRect(10, 20, -30, 40)
	pt4 := geom.NewPoint(5, 15)
	expanded4 := negativeR.Expand(pt4)
	c.Equal(5, expanded4.X)
	c.Equal(15, expanded4.Y)
	c.Equal(0, expanded4.Width)
	c.Equal(0, expanded4.Height)
}

//nolint:gocritic // The "commented out code" is actually explanation
func TestRectInset(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	insets := geom.NewInsets(5, 3, 7, 2) // top, left, bottom, right

	insetRect := r.Inset(insets)

	// X should be moved right by left inset: 10 + 3 = 13
	// Y should be moved down by top inset: 20 + 5 = 25
	// Width should be reduced by left + right: 30 - 3 - 2 = 25
	// Height should be reduced by top + bottom: 40 - 5 - 7 = 28
	c.Equal(13, insetRect.X)
	c.Equal(25, insetRect.Y)
	c.Equal(25, insetRect.Width)
	c.Equal(28, insetRect.Height)

	// Test with insets larger than rectangle dimensions
	largeInsets := geom.NewInsets(50, 40, 50, 40)
	insetRect2 := r.Inset(largeInsets)

	c.Equal(50, insetRect2.X)     // 10 + 40
	c.Equal(70, insetRect2.Y)     // 20 + 50
	c.Equal(0, insetRect2.Width)  // max(30 - 80, 0) = 0
	c.Equal(0, insetRect2.Height) // max(40 - 100, 0) = 0
}

func TestRectString(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10, 20, 30, 40)
	str := r.String()
	c.Equal("10,20,30,40", str)

	// Test with float
	rf := geom.NewRect(10.5, 20.7, 30.2, 40.8)
	strf := rf.String()
	c.Equal("10.5,20.7,30.2,40.8", strf)
}

func TestRectIntersectsLine(t *testing.T) {
	c := check.New(t)

	r := geom.NewRect(10.0, 20.0, 30.0, 40.0)

	// Line completely inside the rectangle
	start1 := geom.NewPoint(15.0, 25.0)
	end1 := geom.NewPoint(25.0, 35.0)
	c.True(r.IntersectsLine(start1, end1))

	// Line with one endpoint inside, one outside
	start2 := geom.NewPoint(15.0, 25.0)
	end2 := geom.NewPoint(50.0, 70.0)
	c.True(r.IntersectsLine(start2, end2))

	// Line crossing through the rectangle
	start3 := geom.NewPoint(5.0, 30.0)
	end3 := geom.NewPoint(45.0, 30.0)
	c.True(r.IntersectsLine(start3, end3))

	// Line completely outside the rectangle
	start4 := geom.NewPoint(50.0, 70.0)
	end4 := geom.NewPoint(60.0, 80.0)
	c.False(r.IntersectsLine(start4, end4))

	// Line parallel to rectangle edge but outside
	start5 := geom.NewPoint(5.0, 5.0)
	end5 := geom.NewPoint(45.0, 5.0)
	c.False(r.IntersectsLine(start5, end5))

	// Line intersecting rectangle corner
	start6 := geom.NewPoint(5.0, 15.0)
	end6 := geom.NewPoint(15.0, 25.0)
	c.True(r.IntersectsLine(start6, end6))

	// Test with empty rectangle
	emptyR := geom.NewRect(10.0, 20.0, 0.0, 40.0)
	start7 := geom.NewPoint(5.0, 25.0)
	end7 := geom.NewPoint(15.0, 25.0)
	c.False(emptyR.IntersectsLine(start7, end7))
}
