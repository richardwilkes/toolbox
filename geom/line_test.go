// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom_test

import (
	"math"
	"slices"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

func TestLineIntersection(t *testing.T) {
	c := check.New(t)

	// Test intersecting lines
	l1 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10))
	l2 := geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0))
	intersection := l1.Intersection(l2)
	c.Equal(1, len(intersection))
	c.Equal(float32(5), intersection[0].X)
	c.Equal(float32(5), intersection[0].Y)

	// Test parallel lines (no intersection)
	l3 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0))
	l4 := geom.NewLine(geom.NewPoint(0, 5), geom.NewPoint(10, 5))
	c.Equal(0, len(l3.Intersection(l4)))

	// Test identical points
	l5 := geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5, 5))
	l6 := geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5, 5))
	identicalIntersection := l5.Intersection(l6)
	c.Equal(1, len(identicalIntersection))
	c.Equal(float32(5), identicalIntersection[0].X)
	c.Equal(float32(5), identicalIntersection[0].Y)

	// Test point intersecting with line
	pointLineIntersection := l5.Intersection(l4)
	c.Equal(1, len(pointLineIntersection))
	c.Equal(float32(5), pointLineIntersection[0].X)
	c.Equal(float32(5), pointLineIntersection[0].Y)

	// Test lines that don't intersect within their segments
	l7 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(1, 1))
	l8 := geom.NewLine(geom.NewPoint(2, 0), geom.NewPoint(3, 1))
	c.Equal(0, len(l7.Intersection(l8)))

	// Test overlapping line segments
	l9 := geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(15, 0))
	overlapping := l3.Intersection(l9)
	c.Equal(2, len(overlapping))

	// The overlapping segment should be from (5,0) to (10,0)
	// Order might vary, so check both points are present
	foundStart := false
	foundEnd := false
	for _, pt := range overlapping {
		if pt.X == 5 && pt.Y == 0 {
			foundStart = true
		}
		if pt.X == 10 && pt.Y == 0 {
			foundEnd = true
		}
	}
	c.True(foundStart)
	c.True(foundEnd)
}

// TestLineIntersectionParallel checks parallel segments, which intersect only when collinear and overlapping or
// touching. The diagonal rows have non-axis-aligned direction vectors, so neither collinearity numerator is trivially
// zero.
func TestLineIntersectionParallel(t *testing.T) {
	c := check.New(t)

	for _, tc := range []struct {
		name string
		want []geom.Point
		a, b geom.Line
	}{
		{"horizontal disjoint", nil, newLine(0, 0, 1, 0), newLine(2, 0, 3, 0)},
		{"horizontal disjoint, reversed", nil, newLine(2, 0, 3, 0), newLine(0, 0, 1, 0)},
		{"vertical disjoint", nil, newLine(0, 0, 0, 1), newLine(0, 2, 0, 3)},
		{"horizontal touching", []geom.Point{geom.NewPoint(1, 0)}, newLine(0, 0, 1, 0), newLine(1, 0, 2, 0)},
		{
			"horizontal overlap",
			[]geom.Point{geom.NewPoint(1, 0), geom.NewPoint(2, 0)},
			newLine(0, 0, 2, 0),
			newLine(1, 0, 3, 0),
		},
		{
			"diagonal overlap",
			[]geom.Point{geom.NewPoint(1, 1), geom.NewPoint(3, 3)},
			newLine(0, 0, 4, 4),
			newLine(1, 1, 3, 3),
		},
		{"diagonal touching", []geom.Point{geom.NewPoint(4, 4)}, newLine(0, 0, 4, 4), newLine(4, 4, 6, 6)},
		{"diagonal offset", nil, newLine(0, 0, 4, 4), newLine(1, 0, 5, 4)},
		{"diagonal offset, reversed", nil, newLine(1, 0, 5, 4), newLine(0, 0, 4, 4)},
		// (t,t) lies exactly on y=x, but the cross products of the collinearity tests are inexact in float32, where
		// arm64's fused multiply-subtract turned their zero into rounding noise. The power-of-two length keeps the
		// interpolated overlap endpoints exact.
		{
			"diagonal non-integer overlap",
			[]geom.Point{geom.NewPoint(5.848, 5.848), geom.NewPoint(20, 20)},
			newLine(0, 0, 8192, 8192),
			newLine(5.848, 5.848, 20, 20),
		},
	} {
		got := tc.a.Intersection(tc.b)
		c.Equal(len(tc.want), len(got), "%s: %v", tc.name, got)
		for _, pt := range tc.want {
			c.True(slices.Contains(got, pt), "%s: %v", tc.name, got)
		}
		c.Equal(len(tc.want) != 0, tc.a.Intersects(tc.b), tc.name)
	}

	// A rect must not report intersection with a line collinear with an edge but entirely beyond it.
	c.False(geom.NewRect(0, 0, 10, 10).IntersectsLine(geom.NewPoint(20, 10), geom.NewPoint(30, 10)))
}

func TestLineBoundsAtLargeCoordinates(t *testing.T) {
	c := check.New(t)

	// The padding must survive float32 rounding at every magnitude. A fixed pad falls below half a ULP at 4096 and
	// up, which used to leave an axis-aligned line with a degenerate zero-thickness bounds that intersects nothing,
	// not even itself.
	vertical := geom.NewLine(geom.NewPoint(12000, 0), geom.NewPoint(12000, 3000))
	bounds := vertical.Bounds()
	c.True(bounds.X < 12000, "%v", bounds)
	c.True(bounds.Right() > 12000, "%v", bounds)
	c.True(bounds.Intersects(bounds), "%v", bounds)

	horizontal := geom.NewLine(geom.NewPoint(-20000, -8000), geom.NewPoint(-10000, -8000))
	bounds = horizontal.Bounds()
	c.True(bounds.Y < -8000, "%v", bounds)
	c.True(bounds.Bottom() > -8000, "%v", bounds)
	c.True(bounds.Intersects(bounds), "%v", bounds)

	// Small coordinates keep the fixed floor.
	small := geom.NewLine(geom.NewPoint(1, 2), geom.NewPoint(3, 2))
	bounds = small.Bounds()
	c.True(bounds.Y < 2 && bounds.Bottom() > 2, "%v", bounds)
}

func TestPointSegmentDistance(t *testing.T) {
	c := check.New(t)

	for _, tc := range []struct {
		name      string
		line      geom.Line
		pt        geom.Point
		squared   float32
		distance  float32
		tolerance float64
	}{
		{"on the segment", newLine(0, 0, 10, 0), geom.NewPoint(5, 0), 0, 0, 0},
		{"perpendicular to the segment", newLine(0, 0, 10, 0), geom.NewPoint(5, 3), 9, 3, 0},
		{"nearest an endpoint", newLine(0, 0, 10, 0), geom.NewPoint(-2, 4), 20, float32(math.Sqrt(20)), 0},
		{"vertical segment", newLine(0, 0, 0, 10), geom.NewPoint(3, 5), 9, 3, 0},
		// A formulation that subtracts two nearly equal float32 squared magnitudes cancels catastrophically here.
		{"near a long segment", newLine(0, 5, 10000, 5), geom.NewPoint(5000, 6), 1, 1, 1e-4},
		{"near a very long segment", newLine(-1e6, 50, 1e6, 50), geom.NewPoint(50, 45), 25, 5, 1e-3},
		{"degenerate segment", newLine(3, 4, 3, 4), geom.NewPoint(0, 0), 25, 5, 0},
	} {
		sq := tc.line.DistanceToPointSquared(tc.pt)
		c.True(math.Abs(float64(sq-tc.squared)) <= tc.tolerance, "%s: squared distance %v", tc.name, sq)
		d := tc.line.DistanceToPoint(tc.pt)
		c.True(math.Abs(float64(d-tc.distance)) <= tc.tolerance, "%s: distance %v", tc.name, d)
	}

	// Verify consistency with PointSegmentDistance
	l3 := geom.NewLine(geom.NewPoint(1, 2), geom.NewPoint(5, 6))
	p5 := geom.NewPoint(3, 1)
	distance := l3.DistanceToPoint(p5)
	distanceSquared := l3.DistanceToPointSquared(p5)
	// distance^2 should equal distanceSquared (within floating point tolerance)
	calculatedSquared := distance * distance
	c.True(calculatedSquared > distanceSquared-0.001 && calculatedSquared < distanceSquared+0.001)
}

func newLine(x1, y1, x2, y2 float32) geom.Line {
	return geom.NewLine(geom.NewPoint(x1, y1), geom.NewPoint(x2, y2))
}
