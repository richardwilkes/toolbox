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
	"math/rand/v2"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

func TestLessThanSortsMissedLinesLast(t *testing.T) {
	c := check.New(t)

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

	h := newLineHeap(lines, position, New(geom.NewRect(0, 0, 100, 100), nil).epsilon)

	c.True(h.lessThan(near, far, destination))
	c.False(h.lessThan(far, near, destination))

	// A line the ray misses has to sort as infinitely far away, so a line the ray actually hits always displaces it.
	// Reporting false in both directions would break the strict weak ordering the heap relies on and would leave the
	// missed line parked at the front of the heap as the reported nearest occluder.
	c.True(h.lessThan(near, missed, destination))
	c.False(h.lessThan(missed, near, destination))
	c.True(h.lessThan(far, missed, destination))
	c.False(h.lessThan(missed, far, destination))

	// Two missed lines are equivalent, which requires false in both directions.
	c.False(h.lessThan(missed, alsoMissed, destination))
	c.False(h.lessThan(alsoMissed, missed, destination))
}

func TestComputePolygonToleratesEmptyHeap(t *testing.T) {
	c := check.New(t)

	// PolygonFrom always appends the four bounds edges, so it is not known to be able to drain the heap. Calling
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

func TestIntersectLinesRejectsNearlyParallel(t *testing.T) {
	c := check.New(t)

	// Exactly parallel.
	_, ok := intersectLines(geom.NewPoint(0, 0), geom.NewPoint(10, 0), geom.NewPoint(0, 5), geom.NewPoint(10, 5))
	c.False(ok)

	// Nearly parallel. Testing ub against zero exactly used to accept this and hand back a point tens of millions of
	// units away, a distance that dominates every comparison it then takes part in.
	_, ok = intersectLines(geom.NewPoint(0, 0), geom.NewPoint(10, 0), geom.NewPoint(0, 5), geom.NewPoint(10, 5.000001))
	c.False(ok)

	// A shallow but real crossing still has to be reported.
	pt, ok := intersectLines(geom.NewPoint(0, 0), geom.NewPoint(10, 0), geom.NewPoint(0, 5), geom.NewPoint(10, 4.99))
	c.True(ok)
	c.True(pt.X > 10, "expected the crossing beyond the segment, got %v", pt)

	// So does a square one.
	pt, ok = intersectLines(geom.NewPoint(0, 0), geom.NewPoint(10, 0), geom.NewPoint(5, -5), geom.NewPoint(5, 5))
	c.True(ok)
	c.Equal(geom.NewPoint(5, 0), pt)
}

func TestLineHeapMaintainsInvariants(t *testing.T) {
	c := check.New(t)

	rng := rand.New(rand.NewPCG(1, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	lines := make([]geom.Line, 24)
	for i := range lines {
		lines[i] = geom.NewLine(geom.NewPoint(rng.Float32()*100, rng.Float32()*100),
			geom.NewPoint(rng.Float32()*100, rng.Float32()*100))
	}
	viewPt := geom.NewPoint(50, 50)
	// One destination throughout gives a single consistent ordering, which is what these invariants are stated
	// against. The sweep varies the destination per operation, so it only ever relies on the ordering being correct
	// for the destination in hand.
	destination := geom.NewPoint(120, 63)
	h := newLineHeap(lines, viewPt, pointEpsilon(100))
	present := make(map[int]bool, len(lines))
	for range 4000 {
		lineIndex := rng.IntN(len(lines))
		if present[lineIndex] {
			h.remove(lineIndex, destination)
			delete(present, lineIndex)
		} else {
			h.insert(lineIndex, destination)
			present[lineIndex] = true
		}
		c.Equal(len(present), len(h.order))
		for i := range h.order {
			c.True(present[h.order[i]], "heap holds line %d that was removed", h.order[i])
			c.Equal(i, h.slot[h.order[i]])
			if i > 0 {
				c.False(h.lessThan(h.order[i], h.order[(i-1)/2], destination), "heap order broken at position %d", i)
			}
			// Ties are legitimate, so the front only has to be a minimum, not the minimum.
			c.False(h.lessThan(h.order[i], h.nearest(), destination), "position %d beats the front", i)
		}
		for lineIndex := range lines {
			c.Equal(present[lineIndex], h.contains(lineIndex))
		}
	}
	c.True(len(present) > 0, "the run left the heap empty, so it proved little")
}
