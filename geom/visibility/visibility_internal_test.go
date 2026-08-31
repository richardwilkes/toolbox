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

func TestLessThanIsStrictWeakOrderingAtEpsilonScale(t *testing.T) {
	c := check.New(t)

	// Walls whose ray crossings sit within a few tolerances of each other used to be ordered through a non-transitive
	// pairwise epsilon-equality gate, which admitted genuine cycles (a < b < c < a) that every insertion order turned
	// into a heap whose front was beaten by another element. The bucketed comparison has to be asymmetric and
	// transitive, with transitive equivalence, for every pair and triple.
	eps := pointEpsilon(100, 100)
	viewPt := geom.NewPoint(0, 0)
	destination := geom.NewPoint(120, 0)
	rng := rand.New(rand.NewPCG(5150, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	lines := make([]geom.Line, 16)
	for i := range lines {
		// Nearly vertical walls crossing the ray at 10 plus a handful of epsilons, with epsilon-scale slopes.
		x := 10 + (rng.Float32()-0.5)*8*eps
		lines[i] = geom.NewLine(
			geom.NewPoint(x+(rng.Float32()-0.5)*4*eps, -5),
			geom.NewPoint(x+(rng.Float32()-0.5)*4*eps, 5),
		)
	}
	h := newLineHeap(lines, viewPt, eps)
	less := func(a, b int) bool { return h.lessThan(a, b, destination) }
	for a := range lines {
		for b := range lines {
			if less(a, b) {
				c.False(less(b, a), "asymmetry violated for %d,%d", a, b)
			}
			for d := range lines {
				if less(a, b) && less(b, d) {
					c.True(less(a, d), "transitivity violated for %d,%d,%d", a, b, d)
				}
				if !less(a, b) && !less(b, a) && !less(b, d) && !less(d, b) {
					c.False(less(a, d) || less(d, a), "equivalence not transitive for %d,%d,%d", a, b, d)
				}
			}
		}
	}
}

func TestComputePolygonToleratesEmptyHeap(t *testing.T) {
	c := check.New(t)

	// PolygonFrom always appends the four bounds edges, so it is not known to be able to drain the heap. Calling
	// computePolygon directly with lines that do not enclose the view point does drain it, which is enough to pin the
	// guards on the nearest-occluder lookups. computePolygon works in scene-local coordinates, so the bounds is
	// centered on the origin to make the world and scene-local frames coincide, keeping these lines and the view point
	// in the same frame as v.bounds rather than exercising the guards only by geometric luck.
	v := New(geom.NewRect(-50, -50, 100, 100), nil)
	viewPt := geom.NewPoint(0, 0)
	for _, lines := range [][]geom.Line{
		{geom.NewLine(geom.NewPoint(-40, -40), geom.NewPoint(-30, -30))},
		{geom.NewLine(geom.NewPoint(-40, -40), geom.NewPoint(40, -40))},
		{geom.NewLine(geom.NewPoint(-10, -10), geom.NewPoint(10, 10))},
		{
			geom.NewLine(geom.NewPoint(-40, 0), geom.NewPoint(40, 0)),
			geom.NewLine(geom.NewPoint(-30, -30), geom.NewPoint(-20, -20)),
		},
	} {
		c.NotPanics(func() { v.computePolygon(viewPt, lines) }, "lines: %v", lines)
	}
}

func TestPointEpsilonTracksFeatureAndExtent(t *testing.T) {
	c := check.New(t)

	c.Equal(float32(0.01), pointEpsilon(100, 100))
	c.Equal(float32(0.0001), pointEpsilon(1, 1))
	// An elongated scene takes its tolerance from the finer dimension so that dimension is not swallowed, with the
	// noise floor still proportional to the overall extent, which is what bounds the rounding noise of scene-local
	// coordinates.
	c.Equal(float32(1000)*extentNoiseRatio, pointEpsilon(0.05, 1000))
	c.Equal(float32(minPointEpsilon), pointEpsilon(0, 0))
	c.Equal(float32(minPointEpsilon), pointEpsilon(-1, -1))
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

// checkLineHeapInvariants drives a heap over the given lines through a random insert/remove sequence, verifying the
// structural invariants after every operation.
func checkLineHeapInvariants(t *testing.T, lines []geom.Line, epsilon float32) {
	t.Helper()
	c := check.New(t)
	rng := rand.New(rand.NewPCG(1, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	viewPt := geom.NewPoint(50, 50)
	// One destination throughout gives a single consistent ordering, which is what these invariants are stated
	// against. The sweep varies the destination per operation, so it only ever relies on the ordering being correct
	// for the destination in hand.
	destination := geom.NewPoint(120, 63)
	h := newLineHeap(lines, viewPt, epsilon)
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

func TestLineHeapMaintainsInvariants(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	lines := make([]geom.Line, 24)
	for i := range lines {
		lines[i] = geom.NewLine(geom.NewPoint(rng.Float32()*100, rng.Float32()*100),
			geom.NewPoint(rng.Float32()*100, rng.Float32()*100))
	}
	checkLineHeapInvariants(t, lines, pointEpsilon(100, 100))
}

func TestLineHeapMaintainsInvariantsAtEpsilonScale(t *testing.T) {
	// Lines whose ray crossings cluster within a few tolerances of each other exercise the tie-handling paths of the
	// ordering, which lines spread far apart relative to the tolerance never reach.
	eps := pointEpsilon(100, 100)
	rng := rand.New(rand.NewPCG(2, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	lines := make([]geom.Line, 24)
	for i := range lines {
		x := 60 + (rng.Float32()-0.5)*8*eps
		lines[i] = geom.NewLine(
			geom.NewPoint(x+(rng.Float32()-0.5)*4*eps, 40+rng.Float32()*10),
			geom.NewPoint(x+(rng.Float32()-0.5)*4*eps, 55+rng.Float32()*10),
		)
	}
	checkLineHeapInvariants(t, lines, eps)
}
