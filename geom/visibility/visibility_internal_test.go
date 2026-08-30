// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

func TestLessThanSortsMissedLinesLast(t *testing.T) {
	c := check.New(t)

	v := New(geom.NewRect(0, 0, 100, 100), nil)
	position := geom.NewPoint(0, 0)
	destination := geom.NewPoint(1, 0) // The ray runs along +X
	const (
		near = iota
		far
		missed
		alsoMissed
	)
	lines := []geom.Line{
		near:       geom.NewLine(geom.NewPoint(10, -5), geom.NewPoint(10, 5)),
		far:        geom.NewLine(geom.NewPoint(20, -5), geom.NewPoint(20, 5)),
		missed:     geom.NewLine(geom.NewPoint(5, 3), geom.NewPoint(15, 3)),   // Parallel to the ray
		alsoMissed: geom.NewLine(geom.NewPoint(5, -3), geom.NewPoint(15, -3)), // Parallel to the ray
	}

	c.True(v.lessThan(near, far, lines, position, destination))
	c.False(v.lessThan(far, near, lines, position, destination))

	// A line the ray misses has to sort as infinitely far away, so a line the ray actually hits always displaces it.
	// Reporting false in both directions would break the strict weak ordering the heap relies on and would leave the
	// missed line parked at the front of the heap as the reported nearest occluder.
	c.True(v.lessThan(near, missed, lines, position, destination))
	c.False(v.lessThan(missed, near, lines, position, destination))
	c.True(v.lessThan(far, missed, lines, position, destination))
	c.False(v.lessThan(missed, far, lines, position, destination))

	// Two missed lines are equivalent, which requires false in both directions.
	c.False(v.lessThan(missed, alsoMissed, lines, position, destination))
	c.False(v.lessThan(alsoMissed, missed, lines, position, destination))
}

func TestComputePolygonToleratesEmptyHeap(t *testing.T) {
	c := check.New(t)

	// SetViewPoint always appends the four bounds edges, so it is not known to be able to drain the heap. Calling
	// computePolygon directly with lines that do not enclose the view point does drain it, which is enough to pin the
	// guards on the nearest-occluder lookups.
	v := New(geom.NewRect(0, 0, 100, 100), nil)
	viewPt := geom.NewPoint(50, 50)
	for _, lines := range [][]geom.Line{
		{geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20))},
		{geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10))},
		{geom.NewLine(geom.NewPoint(40, 40), geom.NewPoint(60, 60))},
		{
			geom.NewLine(geom.NewPoint(10, 50), geom.NewPoint(90, 50)),
			geom.NewLine(geom.NewPoint(20, 20), geom.NewPoint(30, 30)),
		},
	} {
		c.NotPanics(func() { v.computePolygon(viewPt, lines) }, "lines: %v", lines)
	}
}

func TestPointEpsilonTracksMagnitude(t *testing.T) {
	c := check.New(t)

	c.Equal(float32(0.01), pointEpsilon(100))
	c.Equal(float32(0.0001), pointEpsilon(1))
	c.Equal(float32(minPointEpsilon), pointEpsilon(0))
	c.Equal(float32(minPointEpsilon), pointEpsilon(-1))
}
