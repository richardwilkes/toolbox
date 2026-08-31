// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility_test

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sync"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/geom/visibility"
)

// polygonArea returns the absolute area enclosed by the polygon, via the shoelace formula.
func polygonArea(pts []geom.Point) float64 {
	var area float64
	for i := range pts {
		j := (i + 1) % len(pts)
		area += float64(pts[i].X)*float64(pts[j].Y) - float64(pts[j].X)*float64(pts[i].Y)
	}
	return math.Abs(area) / 2
}

// checkPolygon asserts the vertex count and enclosed area of a polygon, and that none of its edges, including the
// implicit closing one, has zero length. The area is accumulated from float32 coordinates, so it is compared with a
// tolerance rather than exactly.
func checkPolygon(c check.Checker, polygon []geom.Point, wantVertices int, wantArea float64, name string) {
	c.Equal(wantVertices, len(polygon), "%s: %v", name, polygon)
	if len(polygon) == 0 {
		return
	}
	if area := polygonArea(polygon); math.Abs(area-wantArea) > 0.01 {
		c.Equal(wantArea, area, "%s: %v", name, polygon)
	}
	for i := range polygon {
		c.NotEqual(polygon[i], polygon[(i+1)%len(polygon)], "%s: zero-length edge at %d in %v", name, i, polygon)
	}
}

// normalizedLines returns the lines with each one's endpoints in a canonical order and the lines themselves sorted, so
// that a set of lines can be compared without depending on the order the algorithm happens to produce.
func normalizedLines(lines []geom.Line) []geom.Line {
	result := make([]geom.Line, len(lines))
	for i, one := range lines {
		if one.End.X < one.Start.X || (one.End.X == one.Start.X && one.End.Y < one.Start.Y) {
			one.Start, one.End = one.End, one.Start
		}
		result[i] = one
	}
	slices.SortFunc(result, func(a, b geom.Line) int {
		return cmp.Or(
			cmp.Compare(a.Start.X, b.Start.X),
			cmp.Compare(a.Start.Y, b.Start.Y),
			cmp.Compare(a.End.X, b.End.X),
			cmp.Compare(a.End.Y, b.End.Y),
		)
	})
	return result
}

func TestNew(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(40, 40), geom.NewPoint(60, 40)),
		geom.NewLine(geom.NewPoint(60, 40), geom.NewPoint(60, 60)),
		geom.NewLine(geom.NewPoint(60, 60), geom.NewPoint(40, 60)),
		geom.NewLine(geom.NewPoint(40, 60), geom.NewPoint(40, 40)),
	}

	v := visibility.New(bounds, obstructions)
	c.NotNil(v)
	before := v.PolygonFrom(geom.NewPoint(20, 20))
	checkPolygon(c, before, 8, 7200, "square obstruction seen from outside")

	// New copies the obstructions, so later edits to the caller's slice must not reach the Visibility.
	obstructions[0] = geom.NewLine(geom.NewPoint(0, 90), geom.NewPoint(100, 90))
	obstructions[1] = geom.NewLine(geom.NewPoint(0, 95), geom.NewPoint(100, 95))
	c.Equal(before, v.PolygonFrom(geom.NewPoint(20, 20)))
}

func TestNewWithEmptyObstructions(t *testing.T) {
	c := check.New(t)

	v := visibility.New(geom.NewRect(0, 0, 100, 100), nil)
	c.NotNil(v)

	// With no obstructions, the entire bounds is visible, so the polygon is exactly its four corners.
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 4, 10000, "no obstructions")
}

func TestBreakIntersections(t *testing.T) {
	c := check.New(t)

	// Two lines that never meet come back untouched.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 10)),
	}
	c.Equal(lines, visibility.BreakIntersections(lines))

	// Two that cross are each split at the crossing.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	})))
}

func TestBreakIntersectionsWithEmptySlice(t *testing.T) {
	c := check.New(t)

	c.Equal(0, len(visibility.BreakIntersections(nil)))
	c.Equal(0, len(visibility.BreakIntersections([]geom.Line{})))
}

func TestBreakIntersectionsWithParallelLines(t *testing.T) {
	c := check.New(t)

	// Parallel lines never meet, so neither one is cut.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(0, 1), geom.NewPoint(10, 1)),
	}
	c.Equal(lines, visibility.BreakIntersections(lines))
}

func TestBreakIntersectionsWithDisjointCollinearLines(t *testing.T) {
	c := check.New(t)

	// Collinear but with a gap between them, so there is no shared portion to split out.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 0)),
		geom.NewLine(geom.NewPoint(6, 0), geom.NewPoint(10, 0)),
	}
	c.Equal(lines, visibility.BreakIntersections(lines))
}

func TestBreakIntersectionsWithCrossingLines(t *testing.T) {
	c := check.New(t)

	// A diagonal crossed by a vertical: both are cut at the single shared point.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5, 10)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 10)),
	})))
}

func TestBreakIntersectionsWithCollinearLines(t *testing.T) {
	c := check.New(t)

	// Three collinear segments covering 0-10, 5-15 and 12-20 of the same line. Every endpoint that falls inside
	// another segment is a cut, so the result is the five spans between consecutive endpoints, each appearing once.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(10, 0), geom.NewPoint(12, 0)),
		geom.NewLine(geom.NewPoint(12, 0), geom.NewPoint(15, 0)),
		geom.NewLine(geom.NewPoint(15, 0), geom.NewPoint(20, 0)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(15, 0)),
		geom.NewLine(geom.NewPoint(12, 0), geom.NewPoint(20, 0)),
	})))
}

func TestBreakIntersectionsSplitsCollinearOverlaps(t *testing.T) {
	c := check.New(t)

	// Partial overlap: each line must be split at the point where the other one starts or ends, and the shared portion
	// must appear only once.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(10, 0), geom.NewPoint(15, 0)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(15, 0)),
	})))

	// Full containment: the enclosing line must be split at both ends of the enclosed one.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(2, 0)),
		geom.NewLine(geom.NewPoint(2, 0), geom.NewPoint(8, 0)),
		geom.NewLine(geom.NewPoint(8, 0), geom.NewPoint(10, 0)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(2, 0), geom.NewPoint(8, 0)),
	})))

	// Diagonal overlap, to be sure the split is not limited to axis-aligned input.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(20, 20)),
	})))
}

func TestBreakIntersectionsWithConcurrentLines(t *testing.T) {
	c := check.New(t)

	// Three lines meeting at (5,5) yield that same intersection point more than once per line, which must not turn into
	// zero-length segments.
	result := visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 10)),
	})
	for _, one := range result {
		c.NotEqual(one.Start, one.End, "zero-length segment %v", one)
	}
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5, 10)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
	}, normalizedLines(result))
}

func TestBreakIntersectionsWithTouchingLines(t *testing.T) {
	c := check.New(t)

	// Meeting end to end is not an overlap, so neither line is cut.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
	}
	c.Equal(lines, visibility.BreakIntersections(lines))
}

func TestBreakIntersectionsWithDuplicateLines(t *testing.T) {
	c := check.New(t)

	// Skipping the self-comparison by value equality would make each of the two identical lines skip the other as
	// well, leaving both of them uncut by anything they overlap.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	})))
}

func TestPolygonFromOutsideBounds(t *testing.T) {
	c := check.New(t)

	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20)),
	})

	c.Nil(v.PolygonFrom(geom.NewPoint(-10, -10)))
	c.Nil(v.PolygonFrom(geom.NewPoint(110, 110)))
	c.Nil(v.PolygonFrom(geom.NewPoint(50, 110)))
	c.Nil(v.PolygonFrom(geom.NewPoint(float32(math.NaN()), float32(math.NaN()))))
}

func TestPolygonFromInsideBounds(t *testing.T) {
	c := check.New(t)

	// A wall spanning x 10-90 at y=10, seen from (50,50). The rays grazing its ends reach the bottom corners of the
	// bounds exactly, so the shadow is the quadrilateral (10,10)-(90,10)-(100,0)-(0,0), whose area is 900.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 6, 10000-900, "wall at y=10")
}

func TestPolygonFromOnBounds(t *testing.T) {
	c := check.New(t)

	// Point.In treats the bounds as half-open, so the top and left edges are inside and the bottom and right are not.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(90, 10)),
	})
	c.NotNil(v.PolygonFrom(geom.NewPoint(0, 0)))
	c.NotNil(v.PolygonFrom(geom.NewPoint(50, 0)))
	c.Nil(v.PolygonFrom(geom.NewPoint(100, 100)))
}

func TestPolygonFromWithSquareObstruction(t *testing.T) {
	c := check.New(t)

	// The square's silhouette from (20,20) runs through (60,40) and (40,60); extending those rays reaches (100,60) and
	// (60,100), so the shadow is (60,40)-(100,60)-(100,100)-(60,100)-(40,60)-(40,40), whose area is 2800.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(40, 40), geom.NewPoint(60, 40)),
		geom.NewLine(geom.NewPoint(60, 40), geom.NewPoint(60, 60)),
		geom.NewLine(geom.NewPoint(60, 60), geom.NewPoint(40, 60)),
		geom.NewLine(geom.NewPoint(40, 60), geom.NewPoint(40, 40)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(20, 20)), 8, 10000-2800, "square obstruction")
}

func TestPolygonFromWithNoObstructions(t *testing.T) {
	c := check.New(t)

	v := visibility.New(geom.NewRect(0, 0, 100, 100), nil)
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 4, 10000, "no obstructions")
}

func TestPolygonFromWithOriginCenteredBounds(t *testing.T) {
	c := check.New(t)

	// The angles of the corners of a bounds centered on the view point cover all four quadrants, including the pair
	// that straddles the ±180° discontinuity the sweep starts from.
	v := visibility.New(geom.NewRect(-10, -10, 20, 20), nil)
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(0, 0)), 4, 400, "origin-centered bounds")
}

func TestPolygonFromBetweenTwoWalls(t *testing.T) {
	c := check.New(t)

	// Walls at y=3 and y=7 spanning x 3-7, with the view point between them at (5,5). Each wall casts a shadow of
	// area 21 onto the edge beyond it.
	v := visibility.New(geom.NewRect(0, 0, 10, 10), []geom.Line{
		geom.NewLine(geom.NewPoint(3, 3), geom.NewPoint(7, 3)),
		geom.NewLine(geom.NewPoint(3, 7), geom.NewPoint(7, 7)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(5, 5)), 8, 100-21-21, "between two walls")
}

func TestPolygonFromWithObstructionsOutsideBounds(t *testing.T) {
	c := check.New(t)

	v := visibility.New(geom.NewRect(10, 10, 80, 80), []geom.Line{
		geom.NewLine(geom.NewPoint(-10, -10), geom.NewPoint(0, 0)),
		geom.NewLine(geom.NewPoint(100, 100), geom.NewPoint(110, 110)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 4, 6400, "obstructions outside the bounds")
}

func TestPolygonFromWithObstructionsCrossingBounds(t *testing.T) {
	c := check.New(t)

	// A wall along y=50 running in from the left and one along x=50 running past both edges box the view point into
	// the 50x50 upper-left quadrant. The parts of each wall lying outside the bounds must be discarded rather than
	// merely clipped at their crossing, or the sweep casts rays out to their far endpoints.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(-10, 50), geom.NewPoint(50, 50)),
		geom.NewLine(geom.NewPoint(50, -10), geom.NewPoint(50, 110)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(25, 25)), 4, 2500, "obstructions crossing the bounds")
}

func TestPolygonFromWithZeroBounds(t *testing.T) {
	c := check.New(t)

	// No point lies inside an empty rectangle, so there is nothing to see from anywhere.
	v := visibility.New(geom.NewRect(0, 0, 0, 0), nil)
	c.Nil(v.PolygonFrom(geom.NewPoint(0, 0)))

	v = visibility.New(geom.NewRect(10, 10, -5, -5), nil)
	c.Nil(v.PolygonFrom(geom.NewPoint(10, 10)))
}

func TestPolygonFromWithNegativeBounds(t *testing.T) {
	c := check.New(t)

	// A full-width wall along y=0 leaves the view point at (0,-25) seeing only the upper half of the bounds.
	v := visibility.New(geom.NewRect(-50, -50, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(-50, 0), geom.NewPoint(50, 0)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(0, -25)), 4, 5000, "negative bounds")
}

func TestPolygonFromWithSmallObstruction(t *testing.T) {
	c := check.New(t)

	// A wall one hundredth of the bounds tall still casts its shadow: the rays grazing its ends reach (100,50) and
	// (100,53), so the shadow is (50,50)-(50,51)-(100,53)-(100,50), whose area is 100.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(50, 50), geom.NewPoint(50, 51)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(25, 50)), 8, 10000-100, "small obstruction")
}

func TestPolygonFromWithZeroLengthObstruction(t *testing.T) {
	c := check.New(t)

	// A degenerate obstruction has no extent to block anything, so the whole bounds stays visible.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(50, 20), geom.NewPoint(50, 20)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 4, 10000, "zero-length obstruction")
}

func TestPolygonFromWithObstructionAlongBoundsEdge(t *testing.T) {
	c := check.New(t)

	// An obstruction lying on top of one of the implicit bounds segments blocks nothing that the bounds did not
	// already block, so the whole area stays visible. Its endpoints do become vertices along that edge.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(20, 0), geom.NewPoint(80, 0)),
	})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 6, 10000, "obstruction along a bounds edge")
}

func TestPolygonFromWithNonFiniteObstruction(t *testing.T) {
	c := check.New(t)

	// A coordinate that is not a number has neither an angle nor a distance, so such a line is ignored rather than
	// left to corrupt the ordering of everything it is compared against.
	nan := float32(math.NaN())
	inf := float32(math.Inf(1))
	for _, obstruction := range []geom.Line{
		geom.NewLine(geom.NewPoint(nan, nan), geom.NewPoint(50, 50)),
		geom.NewLine(geom.NewPoint(50, nan), geom.NewPoint(50, 50)),
		geom.NewLine(geom.NewPoint(inf, 20), geom.NewPoint(50, 20)),
		geom.NewLine(geom.NewPoint(-inf, 20), geom.NewPoint(50, 20)),
	} {
		v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{obstruction})
		checkPolygon(c, v.PolygonFrom(geom.NewPoint(25, 25)), 4, 10000, fmt.Sprint(obstruction))
	}
}

func TestPolygonFromHasNoDuplicateVertices(t *testing.T) {
	c := check.New(t)

	bounds := geom.NewRect(0, 0, 100, 100)
	for _, tc := range []struct {
		name        string
		obstruction geom.Line
		viewPt      geom.Point
	}{
		{"wall across the middle", geom.NewLine(geom.NewPoint(0, 50), geom.NewPoint(100, 50)), geom.NewPoint(50, 25)},
		{
			"view point on a wall endpoint",
			geom.NewLine(geom.NewPoint(0, 50), geom.NewPoint(100, 50)),
			geom.NewPoint(0, 50),
		},
		{"view point on a wall", geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(100, 100)), geom.NewPoint(50, 50)},
		{"wall along a bounds edge", geom.NewLine(geom.NewPoint(20, 0), geom.NewPoint(80, 0)), geom.NewPoint(50, 50)},
	} {
		result := visibility.New(bounds, []geom.Line{tc.obstruction}).PolygonFrom(tc.viewPt)
		c.True(len(result) > 2, "%s: only %d vertices", tc.name, len(result))
		for i := range result {
			// The closing edge is implicit, so the wrap-around pair has to be distinct too.
			j := (i + 1) % len(result)
			c.NotEqual(result[i], result[j], "%s: zero-length edge at %d in %v", tc.name, i, result)
		}
	}
}

func TestPolygonFromHoldsItsInvariants(t *testing.T) {
	c := check.New(t)

	// Randomized scenes, checked only against the properties that must hold for every one of them: the polygon's
	// vertices are finite and inside the bounds, and it never encloses more than the bounds do. Obstructions are
	// allowed to run past the edges, which is what makes the viewport clipping part of what is under test.
	rng := rand.New(rand.NewPCG(1967, 5150)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for i := range 20000 {
		x := float32(rng.IntN(200) - 100)
		y := float32(rng.IntN(200) - 100)
		width := float32(1 + rng.IntN(200))
		height := float32(1 + rng.IntN(200))
		bounds := geom.NewRect(x, y, width, height)
		obstructions := make([]geom.Line, rng.IntN(5))
		point := func() geom.Point {
			return geom.NewPoint(x+rng.Float32()*width*1.4-width*0.2, y+rng.Float32()*height*1.4-height*0.2)
		}
		for j := range obstructions {
			obstructions[j] = geom.NewLine(point(), point())
		}
		viewPt := geom.NewPoint(x+rng.Float32()*width, y+rng.Float32()*height)
		polygon := visibility.New(bounds, obstructions).PolygonFrom(viewPt)
		// The failure text is only built when something is actually wrong, so that describing 20000 healthy scenes
		// does not cost more than generating them.
		problem := ""
		switch {
		case len(polygon) == 0:
			problem = "no polygon at all"
		case polygonArea(polygon) > float64(width)*float64(height)+0.01:
			problem = fmt.Sprintf("area %v exceeds the bounds", polygonArea(polygon))
		default:
			for _, pt := range polygon {
				switch {
				case math.IsNaN(float64(pt.X)) || math.IsNaN(float64(pt.Y)) ||
					math.IsInf(float64(pt.X), 0) || math.IsInf(float64(pt.Y), 0):
					problem = fmt.Sprintf("non-finite vertex %v", pt)
				case pt.X < bounds.X || pt.X > bounds.Right() || pt.Y < bounds.Y || pt.Y > bounds.Bottom():
					problem = fmt.Sprintf("vertex %v lies outside the bounds", pt)
				}
				if problem != "" {
					break
				}
			}
		}
		c.Equal("", problem, "scene %d: bounds=%v viewPt=%v obstructions=%v polygon=%v", i, bounds, viewPt,
			obstructions, polygon)
	}
}

func TestPolygonFromIsConcurrencySafe(t *testing.T) {
	c := check.New(t)

	// A Visibility is documented as immutable once built, so this asserts the property the documentation promises.
	// The suite runs under -race, which is what makes the test meaningful.
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(40, 40), geom.NewPoint(60, 40)),
		geom.NewLine(geom.NewPoint(60, 40), geom.NewPoint(60, 60)),
		geom.NewLine(geom.NewPoint(60, 60), geom.NewPoint(40, 60)),
		geom.NewLine(geom.NewPoint(40, 60), geom.NewPoint(40, 40)),
	})
	viewPt := geom.NewPoint(20, 20)
	want := v.PolygonFrom(viewPt)
	results := make([][]geom.Point, 8)
	var wg sync.WaitGroup
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = v.PolygonFrom(viewPt)
		}()
	}
	wg.Wait()
	for i := range results {
		c.Equal(want, results[i], "goroutine %d disagreed", i)
	}
}

func TestPolygonFromIsScaleInvariant(t *testing.T) {
	c := check.New(t)

	// The same scene, expressed at wildly different scales, has to produce the same polygon shape. A tolerance fixed
	// in world units instead makes small scenes silently collapse: at a bounds size of 0.01 the square's shadow used
	// to disappear entirely, dropping the polygon from 8 vertices to 5.
	const wantRatio = 0.72
	for _, size := range []float32{1000, 100, 10, 1, 0.1, 0.01, 0.001} {
		bounds := geom.NewRect(0, 0, size, size)
		v := visibility.New(bounds, []geom.Line{
			geom.NewLine(geom.NewPoint(size*0.4, size*0.4), geom.NewPoint(size*0.6, size*0.4)),
			geom.NewLine(geom.NewPoint(size*0.6, size*0.4), geom.NewPoint(size*0.6, size*0.6)),
			geom.NewLine(geom.NewPoint(size*0.6, size*0.6), geom.NewPoint(size*0.4, size*0.6)),
			geom.NewLine(geom.NewPoint(size*0.4, size*0.6), geom.NewPoint(size*0.4, size*0.4)),
		})
		result := v.PolygonFrom(geom.NewPoint(size*0.2, size*0.2))
		c.Equal(8, len(result), "bounds size %v", size)
		ratio := polygonArea(result) / (float64(size) * float64(size))
		c.True(math.Abs(ratio-wantRatio) < 0.001, "bounds size %v: area ratio %v, want %v", size, ratio, wantRatio)
	}
}

func TestNewWithNonFiniteBounds(t *testing.T) {
	c := check.New(t)

	// A bounds with a NaN or infinite coordinate cannot yield finite vertices, so it is treated as empty and every
	// polygon request returns nil rather than garbage such as [+Inf,+Inf NaN,NaN ...].
	nan := float32(math.NaN())
	inf := float32(math.Inf(1))
	obstructions := []geom.Line{geom.NewLine(geom.NewPoint(10, 10), geom.NewPoint(20, 20))}
	for _, bounds := range []geom.Rect{
		geom.NewRect(0, 0, inf, inf),
		geom.NewRect(nan, 0, 100, 100),
		geom.NewRect(0, nan, 100, 100),
		geom.NewRect(0, 0, 100, nan),
		geom.NewRect(float32(math.Inf(-1)), 0, 100, 100),
	} {
		c.Nil(visibility.New(bounds, obstructions).PolygonFrom(geom.NewPoint(50, 50)), "bounds %v", bounds)
	}
}

func TestPolygonFromWithBoundsFarFromOrigin(t *testing.T) {
	c := check.New(t)

	// The comparison tolerance is derived from the scene's extent, not the raw coordinate magnitude. A magnitude-based
	// tolerance exceeds the extent of a small scene far from the origin, collapsing the sweep to a single degenerate
	// vertex.
	for _, tc := range []struct {
		name   string
		bounds geom.Rect
	}{
		{"offset 10000, extent 1", geom.NewRect(10000, 10000, 1, 1)},
		{"offset 100000, extent 3", geom.NewRect(100000, 100000, 3, 3)},
		{"offset 1e6, extent 50", geom.NewRect(1e6, 1e6, 50, 50)},
	} {
		polygon := visibility.New(tc.bounds, nil).PolygonFrom(tc.bounds.Center())
		c.Equal(4, len(polygon), "%s: %v", tc.name, polygon)
		wantArea := float64(tc.bounds.Width) * float64(tc.bounds.Height)
		area := polygonArea(polygon)
		c.True(math.Abs(area-wantArea) <= wantArea*0.01, "%s: area %v, want %v", tc.name, area, wantArea)
	}

	// An obstruction in an offset scene still casts its shadow: this is the scene from TestPolygonFromInsideBounds
	// translated by (10000,10000), so the visible area is the same 9100 square units.
	v := visibility.New(geom.NewRect(10000, 10000, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(10010, 10010), geom.NewPoint(10090, 10010)),
	})
	polygon := v.PolygonFrom(geom.NewPoint(10050, 10050))
	c.Equal(6, len(polygon), "%v", polygon)
	area := polygonArea(polygon)
	c.True(math.Abs(area-9100) <= 9100*0.01, "area %v, want 9100", area)
}

func TestPolygonFromWithSliverBounds(t *testing.T) {
	c := check.New(t)

	// A bounds far thinner than it is long has an epsilon, derived from the long dimension, that swallows the short
	// one, so the sweep's output collapses to fewer than three distinct vertices. The contract is a real polygon or
	// nil, never a degenerate one- or two-vertex slice.
	v := visibility.New(geom.NewRect(0, 0, 100, 0.0001), nil)
	c.Nil(v.PolygonFrom(geom.NewPoint(50, 0.00005)))
}

func TestBreakIntersectionsFarFromOrigin(t *testing.T) {
	c := check.New(t)

	// These cross at (10000.5,10000), half a unit from an endpoint of each line. A tolerance derived from the raw
	// coordinate magnitude (~1) discards the crossing as coincident with those endpoints and returns both lines
	// uncut, violating the postcondition that the results do not intersect.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(10000, 10000), geom.NewPoint(10000.5, 10000)),
		geom.NewLine(geom.NewPoint(10000.5, 9999), geom.NewPoint(10000.5, 10000)),
		geom.NewLine(geom.NewPoint(10000.5, 10000), geom.NewPoint(10000.5, 10001)),
		geom.NewLine(geom.NewPoint(10000.5, 10000), geom.NewPoint(10010, 10000)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(10000, 10000), geom.NewPoint(10010, 10000)),
		geom.NewLine(geom.NewPoint(10000.5, 9999), geom.NewPoint(10000.5, 10001)),
	})))
}

func TestBreakIntersectionsWithNonFiniteLines(t *testing.T) {
	c := check.New(t)

	// A single non-finite line must not poison the batch: an infinite coordinate used to drive the tolerance to +Inf,
	// making every point compare equal and dropping every segment, while a NaN collapsed the tolerance to its floor.
	// The bad line is dropped and the rest of the batch is still split normally.
	nan := float32(math.NaN())
	inf := float32(math.Inf(1))
	crossing := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}
	want := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(10, 10)),
	}
	for _, bad := range []geom.Line{
		geom.NewLine(geom.NewPoint(inf, 20), geom.NewPoint(30, 20)),
		geom.NewLine(geom.NewPoint(-inf, 20), geom.NewPoint(30, 20)),
		geom.NewLine(geom.NewPoint(nan, 20), geom.NewPoint(30, 20)),
		geom.NewLine(geom.NewPoint(20, 20), geom.NewPoint(20, nan)),
	} {
		c.Equal(want, normalizedLines(visibility.BreakIntersections(append(slices.Clone(crossing), bad))),
			"bad line %v", bad)
	}
}

func TestBreakIntersectionsDropsZeroLengthLines(t *testing.T) {
	c := check.New(t)

	// Dropping degenerate input is documented behavior: a zero-length line has nothing to split and nothing to
	// contribute as an obstruction.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5, 5)),
	}
	c.Equal([]geom.Line{lines[0]}, visibility.BreakIntersections(lines))
}

var (
	brokenSink  []geom.Line
	polygonSink []geom.Point
)

func BenchmarkBreakIntersections(b *testing.B) {
	lines := make([]geom.Line, 50)
	for i := range lines {
		angle := float64(i) * 2 * math.Pi / float64(len(lines))
		lines[i] = geom.NewLine(
			geom.NewPoint(0, 0),
			geom.NewPoint(float32(math.Cos(angle)*100), float32(math.Sin(angle)*100)),
		)
	}

	for b.Loop() {
		brokenSink = visibility.BreakIntersections(lines)
	}
}

func BenchmarkPolygonFrom(b *testing.B) {
	v := visibility.New(geom.NewRect(0, 0, 1000, 1000), diagonalObstructions(20))
	viewPoint := geom.NewPoint(500, 500)

	for b.Loop() {
		polygonSink = v.PolygonFrom(viewPoint)
	}
}

func BenchmarkPolygonFromWithManyObstructions(b *testing.B) {
	v := visibility.New(geom.NewRect(0, 0, 1000, 1000), diagonalObstructions(100))
	viewPoint := geom.NewPoint(500, 500)

	for b.Loop() {
		polygonSink = v.PolygonFrom(viewPoint)
	}
}

func diagonalObstructions(count int) []geom.Line {
	obstructions := make([]geom.Line, count)
	for i := range obstructions {
		obstructions[i] = geom.NewLine(
			geom.NewPoint(float32(i*40+50), float32(i*30+50)),
			geom.NewPoint(float32(i*40+100), float32(i*30+100)),
		)
	}
	return obstructions
}

func ExampleNew() {
	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(20, 20), geom.NewPoint(80, 20)),
		geom.NewLine(geom.NewPoint(80, 20), geom.NewPoint(80, 80)),
	}

	v := visibility.New(bounds, obstructions)
	polygon := v.PolygonFrom(geom.NewPoint(10, 10))

	// Use the visibility polygon for rendering, collision detection, etc.
	fmt.Println(len(polygon), "vertices")

	// Output: 7 vertices
}

func ExampleBreakIntersections() {
	// Lines that cross each other
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 10)),
		geom.NewLine(geom.NewPoint(0, 10), geom.NewPoint(10, 0)),
	}

	// Each one comes back split at the crossing, ready to pass to New
	for _, line := range visibility.BreakIntersections(lines) {
		fmt.Println(line)
	}

	// Output:
	// {0,0 5,5}
	// {5,5 10,10}
	// {0,10 5,5}
	// {5,5 10,0}
}

func ExampleVisibility_PolygonFrom() {
	bounds := geom.NewRect(0, 0, 100, 100)
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(30, 30), geom.NewPoint(70, 30)),
		geom.NewLine(geom.NewPoint(70, 30), geom.NewPoint(70, 70)),
	}

	v := visibility.New(bounds, obstructions)
	polygon := v.PolygonFrom(geom.NewPoint(15, 15))

	// polygon holds the vertices of the visible area, in sweep order
	fmt.Println(len(polygon), "vertices,", int(polygonArea(polygon)), "square units visible")

	// Output: 7 vertices, 7672 square units visible
}
