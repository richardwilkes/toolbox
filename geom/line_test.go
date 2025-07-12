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

func TestLineIntersection(t *testing.T) {
	c := check.New(t)

	// Test intersecting lines
	a1 := geom.NewPoint(0, 0)
	a2 := geom.NewPoint(10, 10)
	b1 := geom.NewPoint(0, 10)
	b2 := geom.NewPoint(10, 0)

	intersection := geom.LineIntersection(a1, a2, b1, b2)
	c.Equal(1, len(intersection))
	c.Equal(float32(5), intersection[0].X)
	c.Equal(float32(5), intersection[0].Y)

	// Test parallel lines (no intersection)
	c1 := geom.NewPoint(0, 0)
	c2 := geom.NewPoint(10, 0)
	d1 := geom.NewPoint(0, 5)
	d2 := geom.NewPoint(10, 5)

	noIntersection := geom.LineIntersection(c1, c2, d1, d2)
	c.Equal(0, len(noIntersection))

	// Test identical points
	e1 := geom.NewPoint(5, 5)
	e2 := geom.NewPoint(5, 5)
	f1 := geom.NewPoint(5, 5)
	f2 := geom.NewPoint(5, 5)

	identicalIntersection := geom.LineIntersection(e1, e2, f1, f2)
	c.Equal(1, len(identicalIntersection))
	c.Equal(float32(5), identicalIntersection[0].X)
	c.Equal(float32(5), identicalIntersection[0].Y)

	// Test point intersecting with line
	g1 := geom.NewPoint(5, 5)
	g2 := geom.NewPoint(5, 5) // Same point
	h1 := geom.NewPoint(0, 5)
	h2 := geom.NewPoint(10, 5)

	pointLineIntersection := geom.LineIntersection(g1, g2, h1, h2)
	c.Equal(1, len(pointLineIntersection))
	c.Equal(float32(5), pointLineIntersection[0].X)
	c.Equal(float32(5), pointLineIntersection[0].Y)

	// Test lines that don't intersect within their segments
	i1 := geom.NewPoint(0, 0)
	i2 := geom.NewPoint(1, 1)
	j1 := geom.NewPoint(2, 0)
	j2 := geom.NewPoint(3, 1)

	noSegmentIntersection := geom.LineIntersection(i1, i2, j1, j2)
	c.Equal(0, len(noSegmentIntersection))
}

func TestLineIntersectionOverlapping(t *testing.T) {
	c := check.New(t)

	// Test overlapping line segments
	a1 := geom.NewPoint(0, 0)
	a2 := geom.NewPoint(10, 0)
	b1 := geom.NewPoint(5, 0)
	b2 := geom.NewPoint(15, 0)

	overlapping := geom.LineIntersection(a1, a2, b1, b2)
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

func TestPointSegmentDistance(t *testing.T) {
	c := check.New(t)

	// Test point on the line segment
	s1 := geom.NewPoint(0, 0)
	s2 := geom.NewPoint(10, 0)
	p1 := geom.NewPoint(5, 0)

	distance1 := geom.PointSegmentDistance(s1, s2, p1)
	c.Equal(float32(0), distance1)

	// Test point perpendicular to the line segment
	p2 := geom.NewPoint(5, 3)
	distance2 := geom.PointSegmentDistance(s1, s2, p2)
	c.Equal(float32(3), distance2)

	// Test point closest to an endpoint
	p3 := geom.NewPoint(-2, 4)
	distance3 := geom.PointSegmentDistance(s1, s2, p3)
	expected3 := float32(4.47213595) // sqrt((-2)^2 + 4^2) = sqrt(20) ≈ 4.472
	c.True(distance3 > expected3-0.01 && distance3 < expected3+0.01)

	// Test with vertical line segment
	s3 := geom.NewPoint(0, 0)
	s4 := geom.NewPoint(0, 10)
	p4 := geom.NewPoint(3, 5)

	distance4 := geom.PointSegmentDistance(s3, s4, p4)
	c.Equal(float32(3), distance4)
}

func TestPointSegmentDistanceSquared(t *testing.T) {
	c := check.New(t)

	// Test point on the line segment
	s1 := geom.NewPoint(0, 0)
	s2 := geom.NewPoint(10, 0)
	p1 := geom.NewPoint(5, 0)

	distanceSquared1 := geom.PointSegmentDistanceSquared(s1, s2, p1)
	c.Equal(float32(0), distanceSquared1)

	// Test point perpendicular to the line segment
	p2 := geom.NewPoint(5, 3)
	distanceSquared2 := geom.PointSegmentDistanceSquared(s1, s2, p2)
	c.Equal(float32(9), distanceSquared2) // 3^2 = 9

	// Test point closest to an endpoint
	p3 := geom.NewPoint(-2, 4)
	distanceSquared3 := geom.PointSegmentDistanceSquared(s1, s2, p3)
	c.Equal(float32(20), distanceSquared3) // (-2)^2 + 4^2 = 4 + 16 = 20

	// Verify consistency with PointSegmentDistance
	s3 := geom.NewPoint(1, 2)
	s4 := geom.NewPoint(5, 6)
	p4 := geom.NewPoint(3, 1)

	distance := geom.PointSegmentDistance(s3, s4, p4)
	distanceSquared := geom.PointSegmentDistanceSquared(s3, s4, p4)

	// distance^2 should equal distanceSquared (within floating point tolerance)
	calculatedSquared := distance * distance
	c.True(calculatedSquared > distanceSquared-0.001 && calculatedSquared < distanceSquared+0.001)
}
