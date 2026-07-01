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

func TestLineIntersectionCollinear(t *testing.T) {
	c := check.New(t)

	// Collinear but non-overlapping segments must report no intersection. This previously returned a phantom 2-point
	// overlap because the empty clamped interval (left > right) was not detected.
	horiz := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(1, 0))
	beyond := geom.NewLine(geom.NewPoint(2, 0), geom.NewPoint(3, 0))
	c.Equal(0, len(horiz.Intersection(beyond)))
	c.Equal(0, len(beyond.Intersection(horiz))) // Order independent
	c.False(horiz.Intersects(beyond))

	// Same, but exercising the vertical (ady) branch of the overlap math.
	vert := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(0, 1))
	vertBeyond := geom.NewLine(geom.NewPoint(0, 2), geom.NewPoint(0, 3))
	c.Equal(0, len(vert.Intersection(vertBeyond)))

	// Collinear segments that touch at exactly one endpoint yield a single point.
	touching := geom.NewLine(geom.NewPoint(1, 0), geom.NewPoint(2, 0))
	touch := horiz.Intersection(touching)
	c.Equal(1, len(touch))
	c.Equal(geom.NewPoint(1, 0), touch[0])

	// Genuinely overlapping collinear segments still yield the two-point overlap.
	a := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(2, 0))
	b := geom.NewLine(geom.NewPoint(1, 0), geom.NewPoint(3, 0))
	c.Equal(2, len(a.Intersection(b)))

	// Propagation: a rect must not report intersection with a line collinear with an edge but entirely beyond it.
	c.False(geom.NewRect(0, 0, 10, 10).IntersectsLine(geom.NewPoint(20, 10), geom.NewPoint(30, 10)))
}

func TestLineIntersectionParallelCollinearity(t *testing.T) {
	c := check.New(t)

	// The parallel branch treats two segments as collinear only when both cross-product numerators are zero. These
	// diagonal cases exercise that math with non-axis-aligned direction vectors, so neither numerator is trivially zero.

	// Collinear overlapping diagonal segments: both numerators are zero, so the overlap is reported as two points.
	a := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(4, 4))
	b := geom.NewLine(geom.NewPoint(1, 1), geom.NewPoint(3, 3))
	overlap := a.Intersection(b)
	c.Equal(2, len(overlap))
	foundStart := false
	foundEnd := false
	for _, pt := range overlap {
		if pt == geom.NewPoint(1, 1) {
			foundStart = true
		}
		if pt == geom.NewPoint(3, 3) {
			foundEnd = true
		}
	}
	c.True(foundStart)
	c.True(foundEnd)

	// Collinear diagonal segments touching at a single endpoint yield exactly one point.
	touching := geom.NewLine(geom.NewPoint(4, 4), geom.NewPoint(6, 6))
	touch := a.Intersection(touching)
	c.Equal(1, len(touch))
	c.Equal(geom.NewPoint(4, 4), touch[0])

	// Parallel-but-offset diagonal segments (same slope, different line) are not collinear and must not intersect.
	offset := geom.NewLine(geom.NewPoint(1, 0), geom.NewPoint(5, 4))
	c.Equal(0, len(a.Intersection(offset)))
	c.Equal(0, len(offset.Intersection(a))) // Order independent
	c.False(a.Intersects(offset))
}

func TestPointSegmentDistance(t *testing.T) {
	c := check.New(t)

	// Test point on the line segment
	l1 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0))
	p1 := geom.NewPoint(5, 0)
	c.Equal(float32(0), l1.DistanceToPoint(p1))

	// Test point perpendicular to the line segment
	p2 := geom.NewPoint(5, 3)
	c.Equal(float32(3), l1.DistanceToPoint(p2))

	// Test point closest to an endpoint
	p3 := geom.NewPoint(-2, 4)
	distance3 := l1.DistanceToPoint(p3)
	expected3 := float32(4.47213595) // sqrt((-2)^2 + 4^2) = sqrt(20) ≈ 4.472
	c.True(distance3 > expected3-0.01 && distance3 < expected3+0.01)

	// Test with vertical line segment
	l2 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(0, 10))
	p4 := geom.NewPoint(3, 5)
	c.Equal(float32(3), l2.DistanceToPoint(p4))
}

func TestPointSegmentDistanceSquared(t *testing.T) {
	c := check.New(t)

	// Test point on the line segment
	l1 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0))
	p1 := geom.NewPoint(5, 0)
	c.Equal(float32(0), l1.DistanceToPointSquared(p1))

	// Test point perpendicular to the line segment
	p2 := geom.NewPoint(5, 3)
	c.Equal(float32(9), l1.DistanceToPointSquared(p2))

	// Test point closest to an endpoint
	p3 := geom.NewPoint(-2, 4)
	c.Equal(float32(20), l1.DistanceToPointSquared(p3)) // (-2)^2 + 4^2 = 4 + 16 = 20

	// Test with vertical line segment
	l2 := geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(0, 10))
	p4 := geom.NewPoint(3, 5)
	c.Equal(float32(9), l2.DistanceToPointSquared(p4))

	// Verify consistency with PointSegmentDistance
	l3 := geom.NewLine(geom.NewPoint(1, 2), geom.NewPoint(5, 6))
	p5 := geom.NewPoint(3, 1)
	distance := l3.DistanceToPoint(p5)
	distanceSquared := l3.DistanceToPointSquared(p5)
	// distance^2 should equal distanceSquared (within floating point tolerance)
	calculatedSquared := distance * distance
	c.True(calculatedSquared > distanceSquared-0.001 && calculatedSquared < distanceSquared+0.001)
}
