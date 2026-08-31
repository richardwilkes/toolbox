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

// checkPolygon asserts the vertex count and enclosed area of a polygon, and that no two consecutive vertices,
// including the wrap-around pair, sit within the comparison tolerance of each other. The area tolerance scales with
// the expected area and the vertex separation tolerance with the polygon's extent, so both checks stay meaningful at
// every scene scale instead of being fixed in world units.
func checkPolygon(c check.Checker, polygon []geom.Point, wantVertices int, wantArea float64, name string) {
	c.Equal(wantVertices, len(polygon), "%s: %v", name, polygon)
	if len(polygon) == 0 {
		return
	}
	if area := polygonArea(polygon); math.Abs(area-wantArea) > math.Max(wantArea*1e-5, 1e-9) {
		c.Equal(wantArea, area, "%s: %v", name, polygon)
	}
	minPt := polygon[0]
	maxPt := polygon[0]
	for _, pt := range polygon {
		minPt.X = min(minPt.X, pt.X)
		minPt.Y = min(minPt.Y, pt.Y)
		maxPt.X = max(maxPt.X, pt.X)
		maxPt.Y = max(maxPt.Y, pt.Y)
	}
	// PolygonFrom promises consecutive vertices distinct to within its comparison tolerance, which it derives from the
	// scene's dimensions as max(minDim*1e-4, maxDim*1e-6, 1e-6). The polygon's span never exceeds the scene's, so the
	// same formula applied to the span is a conservative lower bound on that tolerance, and asserting separation
	// beyond it verifies the promised property rather than only its noise-floor term.
	spanX := maxPt.X - minPt.X
	spanY := maxPt.Y - minPt.Y
	eps := max(min(spanX, spanY)*1e-4, max(spanX, spanY)*1e-6, 1e-6)
	for i := range polygon {
		c.False(polygon[i].EqualWithin(polygon[(i+1)%len(polygon)], eps), "%s: zero-length edge at %d in %v", name, i,
			polygon)
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

	// An overlap boundary within the tolerance of one line's endpoint must merge with that endpoint for both lines,
	// not be dropped for one and kept for the other, which used to leave two pieces overlapping by half the scene.
	// This all-collinear input has a bounding box with no finer dimension, so its tolerance is the noise floor derived
	// from the 100-unit span, 1e-4, and the offset sits just inside it.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(50, 0)),
		geom.NewLine(geom.NewPoint(50, 0), geom.NewPoint(100, 0)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(100, 0)),
		geom.NewLine(geom.NewPoint(0.00005, 0), geom.NewPoint(50, 0)),
	})))

	// Cut points come verbatim from the endpoints of the lines in an overlap, not from interpolation along whichever
	// line happened to be examined first, so reversed input cannot produce a one-ULP-different duplicate of the
	// shared portion.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(3, 0)),
		geom.NewLine(geom.NewPoint(3, 0), geom.NewPoint(4, 0)),
		geom.NewLine(geom.NewPoint(4, 0), geom.NewPoint(9, 0)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(3, 0), geom.NewPoint(4, 0)),
		geom.NewLine(geom.NewPoint(9, 0), geom.NewPoint(0, 0)),
	})))
}

func TestBreakIntersectionsResultsNeverIntersect(t *testing.T) {
	c := check.New(t)

	// Random scenes on a small integer grid, dense enough in collinear overlaps and shared endpoints to exercise the
	// overlap grouping hard, checked against the documented postcondition itself: no two result lines intersect
	// anywhere but at their endpoints, and no two overlap.
	rng := rand.New(rand.NewPCG(1, 5150)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for scene := range 2500 {
		lines := make([]geom.Line, 2+rng.IntN(4))
		for i := range lines {
			lines[i] = geom.NewLine(
				geom.NewPoint(float32(rng.IntN(21)), float32(rng.IntN(5))),
				geom.NewPoint(float32(rng.IntN(21)), float32(rng.IntN(5))),
			)
		}
		result := visibility.BreakIntersections(lines)
		for a := range result {
			for b := a + 1; b < len(result); b++ {
				c.False(segmentsViolateNonIntersection(result[a], result[b], 0.01),
					"scene %d: %v and %v intersect: input=%v result=%v", scene, result[a], result[b], lines, result)
			}
		}
	}
}

func TestBreakIntersectionsResultsNeverIntersectNonInteger(t *testing.T) {
	c := check.New(t)

	// Scenes dense in exactly collinear overlaps whose coordinates are not exactly representable in the cross-product
	// arithmetic: points of the form (t,t) lie exactly on y=x for any float32 t, but the products in the collinearity
	// tests are inexact, which is precisely where fused multiply-subtract used to defeat the exact-zero comparisons.
	// An integer grid never leaves the exactly-representable region, so the grid-based test above cannot catch that.
	// A few unconstrained lines are mixed in to exercise the single-crossing paths at the same time.
	rng := rand.New(rand.NewPCG(2, 5150)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for scene := range 2500 {
		lines := make([]geom.Line, 0, 6)
		for range 2 + rng.IntN(3) {
			t1 := rng.Float32() * 20
			t2 := rng.Float32() * 20
			lines = append(lines, geom.NewLine(geom.NewPoint(t1, t1), geom.NewPoint(t2, t2)))
		}
		for range 1 + rng.IntN(2) {
			lines = append(lines, geom.NewLine(
				geom.NewPoint(rng.Float32()*20, rng.Float32()*5),
				geom.NewPoint(rng.Float32()*20, rng.Float32()*5),
			))
		}
		result := visibility.BreakIntersections(lines)
		for a := range result {
			for b := a + 1; b < len(result); b++ {
				c.False(segmentsViolateNonIntersection(result[a], result[b], 0.01),
					"scene %d: %v and %v intersect: input=%v result=%v", scene, result[a], result[b], lines, result)
			}
		}
	}
}

// segmentsViolateNonIntersection reports whether the two segments either properly cross -- each having its endpoints
// strictly on opposite sides of the other's line, by more than tol -- or lie on a common line and overlap by more
// than tol. Touching at endpoints is fine. The predicates are evaluated in float64, since a naive segment
// intersection is itself unreliable for the nearly collinear fragments splitting produces.
func segmentsViolateNonIntersection(a, b geom.Line, tol float64) bool {
	d1 := signedLineDistance(b, a.Start)
	d2 := signedLineDistance(b, a.End)
	d3 := signedLineDistance(a, b.Start)
	d4 := signedLineDistance(a, b.End)
	straddles := func(p, q float64) bool { return (p > tol && q < -tol) || (p < -tol && q > tol) }
	if straddles(d1, d2) && straddles(d3, d4) {
		return true
	}
	if math.Abs(d1) > tol || math.Abs(d2) > tol || math.Abs(d3) > tol || math.Abs(d4) > tol {
		return false
	}
	// Only near-parallel segments can genuinely lie along a common line: short fragments fanning out of a cluster of
	// nearly concurrent crossings pass within tol of each other's lines while heading in clearly different
	// directions, and that is legitimate.
	dx := float64(a.End.X) - float64(a.Start.X)
	dy := float64(a.End.Y) - float64(a.Start.Y)
	bdx := float64(b.End.X) - float64(b.Start.X)
	bdy := float64(b.End.Y) - float64(b.Start.Y)
	if math.Abs(dx*bdy-dy*bdx) > 0.1*math.Hypot(dx, dy)*math.Hypot(bdx, bdy) {
		return false
	}
	// Collinear to within tol: the projections of the two segments onto a's direction must not overlap by more than
	// tol.
	length := math.Hypot(dx, dy)
	dx /= length
	dy /= length
	t1 := (float64(b.Start.X)-float64(a.Start.X))*dx + (float64(b.Start.Y)-float64(a.Start.Y))*dy
	t2 := (float64(b.End.X)-float64(a.Start.X))*dx + (float64(b.End.Y)-float64(a.Start.Y))*dy
	return math.Min(length, math.Max(t1, t2))-math.Max(0, math.Min(t1, t2)) > tol
}

// signedLineDistance returns the signed distance from p to the infinite line through l, evaluated in float64.
func signedLineDistance(l geom.Line, p geom.Point) float64 {
	dx := float64(l.End.X) - float64(l.Start.X)
	dy := float64(l.End.Y) - float64(l.Start.Y)
	return (dx*(float64(p.Y)-float64(l.Start.Y)) - dy*(float64(p.X)-float64(l.Start.X))) / math.Hypot(dx, dy)
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

	// Identical duplicate lines reduce to a single copy, which must still be cut normally by everything it crosses.
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

	// Beyond vertex distinctness, the view-point-on-obstruction cases pin the polygon itself: a zero-thickness wall
	// through or touching the eye occludes only a measure-zero set of rays, so the whole bounds stays visible.
	bounds := geom.NewRect(0, 0, 100, 100)
	for _, tc := range []struct {
		name         string
		obstruction  geom.Line
		viewPt       geom.Point
		wantVertices int
		wantArea     float64
	}{
		{
			"wall across the middle",
			geom.NewLine(geom.NewPoint(0, 50), geom.NewPoint(100, 50)),
			geom.NewPoint(50, 25),
			4,
			5000,
		},
		{
			"view point on a wall endpoint",
			geom.NewLine(geom.NewPoint(0, 50), geom.NewPoint(100, 50)),
			geom.NewPoint(0, 50),
			4,
			10000,
		},
		{
			"view point on a wall",
			geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(100, 100)),
			geom.NewPoint(50, 50),
			4,
			10000,
		},
		{
			"wall along a bounds edge",
			geom.NewLine(geom.NewPoint(20, 0), geom.NewPoint(80, 0)),
			geom.NewPoint(50, 50),
			6,
			10000,
		},
	} {
		result := visibility.New(bounds, []geom.Line{tc.obstruction}).PolygonFrom(tc.viewPt)
		checkPolygon(c, result, tc.wantVertices, tc.wantArea, tc.name)
	}
}

func TestPolygonFromWithViewPointOnObstruction(t *testing.T) {
	c := check.New(t)

	// A zero-thickness wall through the eye occludes only the measure-zero set of rays along its own line, so every
	// rotation of it must leave the entire bounds visible; anything else breaks rotational symmetry.
	bounds := geom.NewRect(0, 0, 100, 100)
	viewPt := geom.NewPoint(50, 50)
	for i := range 24 {
		theta := float64(i) * 15 * math.Pi / 180
		dir := geom.NewPoint(float32(20*math.Cos(theta)), float32(20*math.Sin(theta)))
		wall := geom.NewLine(viewPt.Sub(dir), viewPt.Add(dir))
		polygon := visibility.New(bounds, []geom.Line{wall}).PolygonFrom(viewPt)
		checkPolygon(c, polygon, 4, 10000, fmt.Sprintf("wall rotated %d degrees", i*15))
	}

	// The same holds when the eye sits partway along a chord rather than at the wall's center.
	v := visibility.New(bounds, []geom.Line{geom.NewLine(geom.NewPoint(10, 80), geom.NewPoint(90, 20))})
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(50, 50)), 4, 10000, "view point on a chord")

	// And when the wall merely touches the eye with one endpoint.
	v = visibility.New(bounds, []geom.Line{geom.NewLine(viewPt, geom.NewPoint(90, 90))})
	checkPolygon(c, v.PolygonFrom(viewPt), 4, 10000, "wall endpoint at the view point")
}

func TestPolygonFromWithObstructionNearBoundsCorner(t *testing.T) {
	c := check.New(t)

	// A wall endpoint that subtends a small fraction of a degree from a bounds corner is a distinct event from the
	// corner itself. The old fixed 0.01-degree batching tolerance merged the two, silently deleting the wall's shadow
	// (or, mirrored, leaking area through the gap), so the expected area is computed independently here.
	viewPt := geom.NewPoint(50, 50)
	const offDeg = 0.008 // Degrees off the exact eye-to-corner ray, well inside the old fixed tolerance.
	theta := (225 + offDeg) * math.Pi / 180
	e := geom.NewPoint(float32(50+63.64*math.Cos(theta)), float32(50+63.64*math.Sin(theta)))
	w := geom.NewPoint(90, 5)
	v := visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{geom.NewLine(e, w)})
	polygon := v.PolygonFrom(viewPt)
	c.True(len(polygon) >= 5, "%v", polygon)
	// The wall's shadow is the quadrilateral between the wall and the bottom edge, bounded by the extensions of the
	// rays through its endpoints, both of which land on the bottom edge for this geometry.
	extend := func(pt geom.Point) (x, y float64) {
		t := 50 / (50 - float64(pt.Y))
		return 50 + t*(float64(pt.X)-50), 0
	}
	h1x, h1y := extend(e)
	h2x, h2y := extend(w)
	shadow := [][2]float64{{float64(e.X), float64(e.Y)}, {float64(w.X), float64(w.Y)}, {h2x, h2y}, {h1x, h1y}}
	var area float64
	for i := range shadow {
		j := (i + 1) % len(shadow)
		area += shadow[i][0]*shadow[j][1] - shadow[j][0]*shadow[i][1]
	}
	want := 10000 - math.Abs(area)/2
	got := polygonArea(polygon)
	c.True(math.Abs(got-want) < 2, "area %v, want %v (polygon %v)", got, want, polygon)
}

func TestPolygonFromWithElongatedBounds(t *testing.T) {
	c := check.New(t)

	// The comparison tolerance must track the short dimension of an elongated bounds. Derived from the long one, it
	// swallows the short dimension entirely and the sweep returns a bogus three-vertex polygon with roughly half the
	// true area even with no obstructions at all.
	v := visibility.New(geom.NewRect(0, 0, 1000, 0.05), nil)
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(500, 0.025)), 4, 50, "wide bounds")

	v = visibility.New(geom.NewRect(0, 0, 0.05, 1000), nil)
	checkPolygon(c, v.PolygonFrom(geom.NewPoint(0.025, 500)), 4, 50, "tall bounds")
}

func TestPolygonFromWithHugeBounds(t *testing.T) {
	c := check.New(t)

	// The sweep's starting ray direction is derived from the scene's extent. A fixed +1 offset vanishes below the
	// float32 resolution of scene-local coordinates at 2^24 and beyond, degrading the seeded heap to insertion order
	// and letting far walls show through near ones, so the same scene expressed at a huge scale has to produce the
	// same shape.
	obstructions := []geom.Line{
		geom.NewLine(geom.NewPoint(48, 24), geom.NewPoint(48, 40)),
		geom.NewLine(geom.NewPoint(40, 16), geom.NewPoint(40, 48)),
	}
	base := visibility.New(geom.NewRect(0, 0, 64, 64), obstructions).PolygonFrom(geom.NewPoint(56, 32))
	baseArea := polygonArea(base)
	c.True(baseArea > 0 && baseArea < 64*64, "base area %v", baseArea)

	const scale = 1 << 20 // Puts the scaled view point's scene-local coordinates beyond 2^24.
	scaled := make([]geom.Line, len(obstructions))
	for i, line := range obstructions {
		scaled[i] = geom.NewLine(line.Start.Mul(scale), line.End.Mul(scale))
	}
	huge := visibility.New(geom.NewRect(0, 0, 64*scale, 64*scale), scaled).PolygonFrom(geom.NewPoint(56*scale, 32*scale))
	ratio := polygonArea(huge) / (scale * scale)
	c.True(math.Abs(ratio-baseArea) <= baseArea*0.005, "scaled area ratio %v, want %v", ratio, baseArea)
}

// sightlineVisible reports whether p is visible from viewPt given the obstructions, along with whether the sample is
// far enough from every geometric boundary -- walls and their silhouette-casting endpoints -- for the answer to be
// trustworthy at the tolerance the sweep works to.
func sightlineVisible(viewPt, p geom.Point, obstructions []geom.Line, guard float32) (visible, ok bool) {
	for _, o := range obstructions {
		if geom.PointSegmentDistance(o.Start, o.End, p) <= guard ||
			geom.PointSegmentDistance(viewPt, p, o.Start) <= guard ||
			geom.PointSegmentDistance(viewPt, p, o.End) <= guard {
			return false, false
		}
	}
	for _, o := range obstructions {
		if len(geom.LineIntersection(viewPt, p, o.Start, o.End)) != 0 {
			return false, true
		}
	}
	return true, true
}

// pointInPolygon reports whether p lies inside the polygon, via the even-odd rule evaluated in float64.
func pointInPolygon(polygon []geom.Point, p geom.Point) bool {
	px, py := float64(p.X), float64(p.Y)
	in := false
	j := len(polygon) - 1
	for i := range polygon {
		xi, yi := float64(polygon[i].X), float64(polygon[i].Y)
		xj, yj := float64(polygon[j].X), float64(polygon[j].Y)
		if (yi > py) != (yj > py) && px < (xj-xi)*(py-yi)/(yj-yi)+xi {
			in = !in
		}
		j = i
	}
	return in
}

// nearPolygonBoundary reports whether p lies within guard of any edge of the polygon, including the implicit closing
// edge.
func nearPolygonBoundary(polygon []geom.Point, p geom.Point, guard float32) bool {
	j := len(polygon) - 1
	for i := range polygon {
		if geom.PointSegmentDistance(polygon[j], polygon[i], p) <= guard {
			return true
		}
		j = i
	}
	return false
}

func TestPolygonFromHoldsItsInvariants(t *testing.T) {
	c := check.New(t)

	// Randomized scenes, checked against the properties that must hold for every one of them: the polygon's vertices
	// are finite and inside the bounds, it never encloses more than the bounds do, and sampled interior points agree
	// with an independent line-of-sight computation, which is what catches an under-reported polygon that the area
	// bound and clamped vertices cannot. Obstructions are allowed to run past the edges, which is what makes the
	// viewport clipping part of what is under test, and they are piped through BreakIntersections because New
	// documents that its obstructions must not intersect each other.
	rng := rand.New(rand.NewPCG(1967, 5150)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for i := range 20000 {
		x := float32(rng.IntN(200) - 100)
		y := float32(rng.IntN(200) - 100)
		width := float32(1 + rng.IntN(200))
		height := float32(1 + rng.IntN(200))
		bounds := geom.NewRect(x, y, width, height)
		raw := make([]geom.Line, rng.IntN(5))
		point := func() geom.Point {
			return geom.NewPoint(x+rng.Float32()*width*1.4-width*0.2, y+rng.Float32()*height*1.4-height*0.2)
		}
		for j := range raw {
			raw[j] = geom.NewLine(point(), point())
		}
		obstructions := visibility.BreakIntersections(raw)
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
		if len(polygon) < 3 {
			continue
		}
		// Cross-check sampled points against the independent line-of-sight computation, skipping samples that sit too
		// close to a wall, a silhouette edge, the polygon boundary, or a corner-grazing ray for a tolerance-based
		// sweep to classify them dependably. The guard scales with the smaller dimension so that elongated scenes are
		// still cross-checked rather than skipped wholesale.
		guard := max(min(width, height)*0.01, max(width, height)*0.001)
		if width <= 2*guard || height <= 2*guard {
			continue
		}
		// The view point only has to be clear of the sweep's own drop tolerance -- a wall within epsilon of the eye
		// is deliberately treated as not blocking anything -- rather than of the much larger sampling guard. Keeping
		// the eye guard near the actual drop threshold (the formula mirrors the sweep's epsilon derivation) is what
		// lets these scenes exercise eyes sitting close to walls and to the lines through them, where shared-point
		// event ordering is at its most delicate.
		eyeGuard := 4 * max(min(width, height)*1e-4, max(width, height)*1e-6, 1e-6)
		eyeClear := true
		for _, o := range obstructions {
			if geom.PointSegmentDistance(o.Start, o.End, viewPt) <= eyeGuard {
				eyeClear = false
				break
			}
		}
		if !eyeClear {
			continue
		}
		corners := [4]geom.Point{bounds.Point, bounds.TopRight(), bounds.BottomRight(), bounds.BottomLeft()}
		for range 12 {
			p := geom.NewPoint(x+guard+rng.Float32()*(width-2*guard), y+guard+rng.Float32()*(height-2*guard))
			visible, ok := sightlineVisible(viewPt, p, obstructions, guard)
			if !ok || nearPolygonBoundary(polygon, p, guard) {
				continue
			}
			grazesCorner := false
			for _, corner := range corners {
				if geom.PointSegmentDistance(viewPt, p, corner) <= guard {
					grazesCorner = true
					break
				}
			}
			if grazesCorner {
				continue
			}
			c.Equal(visible, pointInPolygon(polygon, p),
				"scene %d: sight line to %v disagrees with the polygon: bounds=%v viewPt=%v obstructions=%v polygon=%v",
				i, p, bounds, viewPt, obstructions, polygon)
		}
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
		wg.Go(func() {
			results[i] = v.PolygonFrom(viewPt)
		})
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
		{"offset 1e7, extent 100", geom.NewRect(1e7, 1e7, 100, 100)},
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

	// The same obstructed scene at an offset of 1e7, where a tolerance floor derived from the raw coordinate
	// magnitude reaches 10 world units and silently swallows the wall's shadow.
	v = visibility.New(geom.NewRect(1e7, 1e7, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(1e7+10, 1e7+10), geom.NewPoint(1e7+90, 1e7+10)),
	})
	polygon = v.PolygonFrom(geom.NewPoint(1e7+50, 1e7+50))
	c.Equal(6, len(polygon), "%v", polygon)
	area = polygonArea(polygon)
	c.True(math.Abs(area-9100) <= 9100*0.01, "area %v, want 9100", area)

	// A small scene very far from the origin keeps the polygon-or-nil contract too. World float32 coordinates at 5e7
	// are spaced 4 units apart, so the corners round and the area is only representation-accurate, but the polygon
	// must exist and be sane rather than collapse to nil.
	polygon = visibility.New(geom.NewRect(5e7, 5e7, 50, 50), nil).PolygonFrom(geom.NewPoint(5e7+25, 5e7+25))
	c.Equal(4, len(polygon), "%v", polygon)
	area = polygonArea(polygon)
	c.True(math.Abs(area-2500) <= 500, "area %v, want about 2500", area)
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

func TestBreakIntersectionsCollinearDiagonalNonInteger(t *testing.T) {
	c := check.New(t)

	// Exactly collinear diagonal segments with non-integer coordinates: every point of the form (t,t) lies exactly on
	// y=x, but the cross products of the collinearity tests are inexact, and fused multiply-subtract on arm64 used to
	// turn their mathematically zero results into rounding noise, sending these overlaps down the "not parallel"
	// branch and returning them uncut and overlapping. The scene spans 0-20 so that re-centering to (10,10) is exact
	// for every coordinate, letting the expected output be stated verbatim.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(5, 5)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(5.848, 5.848)),
		geom.NewLine(geom.NewPoint(5.848, 5.848), geom.NewPoint(9, 9)),
		geom.NewLine(geom.NewPoint(9, 9), geom.NewPoint(20, 20)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(20, 20)),
		geom.NewLine(geom.NewPoint(5.848, 5.848), geom.NewPoint(20, 20)),
		geom.NewLine(geom.NewPoint(5, 5), geom.NewPoint(9, 9)),
	})))
}

func TestBreakIntersectionsLargeExtentCollinear(t *testing.T) {
	c := check.New(t)

	// Two overlapping vertical lines at x=12000 in a batch whose other geometry keeps the re-centered coordinates near
	// ±12000. Line.Bounds() used to pad by a fixed epsilon that falls below the float32 ULP at that magnitude, giving
	// an axis-aligned line a degenerate zero-thickness bounds that intersects nothing -- not even itself -- so the
	// quadtree never offered the overlapping partner as a candidate and the pair came back uncut.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(-12000, 0), geom.NewPoint(-12000, 10)),
		geom.NewLine(geom.NewPoint(12000, 0), geom.NewPoint(12000, 2000)),
		geom.NewLine(geom.NewPoint(12000, 2000), geom.NewPoint(12000, 3000)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(12000, 3000), geom.NewPoint(12000, 0)),
		geom.NewLine(geom.NewPoint(12000, 3000), geom.NewPoint(12000, 2000)),
		geom.NewLine(geom.NewPoint(-12000, 0), geom.NewPoint(-12000, 10)),
	})))
}

func TestBreakIntersectionsElongatedScene(t *testing.T) {
	c := check.New(t)

	// A short crossing wall in a much longer scene: the preprocessing tolerance used to be derived from the overall
	// extent alone, so it dwarfed the wall, silently deleting it and leaving the crossing uncut, while PolygonFrom --
	// whose tolerance tracks the finer dimension -- would have resolved the same wall's shadow correctly.
	c.Equal([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0.025), geom.NewPoint(500, 0.025)),
		geom.NewLine(geom.NewPoint(500, 0), geom.NewPoint(500, 0.025)),
		geom.NewLine(geom.NewPoint(500, 0.025), geom.NewPoint(500, 0.05)),
		geom.NewLine(geom.NewPoint(500, 0.025), geom.NewPoint(1000, 0.025)),
	}, normalizedLines(visibility.BreakIntersections([]geom.Line{
		geom.NewLine(geom.NewPoint(0, 0.025), geom.NewPoint(1000, 0.025)),
		geom.NewLine(geom.NewPoint(500, 0), geom.NewPoint(500, 0.05)),
	})))

	// A pillar much shorter than the corridor is long, but well above the tolerance the sweep works to, must survive
	// the preprocessing untouched rather than vanish.
	lines := []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(1000, 0)),
		geom.NewLine(geom.NewPoint(0, 1), geom.NewPoint(1000, 1)),
		geom.NewLine(geom.NewPoint(500, 0.4), geom.NewPoint(500, 0.49)),
	}
	c.Equal(lines, visibility.BreakIntersections(lines))
}

func TestPolygonFromWithEyeNearLongWall(t *testing.T) {
	c := check.New(t)

	// The eye-proximity filter used to compute its distance with a formulation that cancels catastrophically for a
	// point near a long segment, silently discarding any wall whose true distance from the eye fell below roughly the
	// wall's length times 2.4e-4. A full-width wall in a long room, with the eye 100 tolerances away, then let the
	// eye see straight through: the whole 100000 came back instead of the 50000 above the wall.
	v := visibility.New(geom.NewRect(0, 0, 10000, 10), []geom.Line{
		geom.NewLine(geom.NewPoint(0, 5), geom.NewPoint(10000, 5)),
	})
	polygon := v.PolygonFrom(geom.NewPoint(5000, 6))
	c.True(len(polygon) >= 4, "%v", polygon)
	area := polygonArea(polygon)
	c.True(math.Abs(area-50000) <= 500, "area %v, want about 50000", area)
	for _, pt := range polygon {
		c.True(pt.Y >= 5-0.1, "vertex %v leaked below the wall", pt)
	}

	// Obstructions are explicitly allowed to run past the bounds, so the error must not scale with the caller's
	// segment length either: this wall is supplied as a two-million-unit segment, of which only the 100-unit span
	// inside the bounds matters, and the eye sits 45 units away from it.
	v = visibility.New(geom.NewRect(0, 0, 100, 100), []geom.Line{
		geom.NewLine(geom.NewPoint(-1e6, 50), geom.NewPoint(1e6+100, 50)),
	})
	polygon = v.PolygonFrom(geom.NewPoint(50, 95))
	c.True(len(polygon) >= 4, "%v", polygon)
	area = polygonArea(polygon)
	c.True(math.Abs(area-5000) <= 50, "area %v, want about 5000", area)
}

func TestPolygonFromWithEyeNearWallLine(t *testing.T) {
	c := check.New(t)

	// An eye sitting close to the infinite line through a wall or bounds edge used to flip the sweep's shared-point
	// tie-break: the tie between a wall and the bounds edge its clipped endpoint lies on was ordered through a
	// two-class angle key with a discontinuity exactly at the near-collinear angle, so whole scene-scale regions
	// leaked past the wall or were cut away. Both scenes are failures recorded while strengthening the randomized
	// invariants test.

	// Over-report: the wall's entire shadow triangle used to leak into the polygon. The eye is 0.55 from the wall's
	// line, and the wall is clipped against the bottom and left edges, so the visible area is the bounds minus the
	// triangle the wall cuts off: 118*36 - (93.4055*34.5872)/2.
	v := visibility.New(geom.NewRect(-30, 45, 118, 36), []geom.Line{
		geom.NewLine(geom.NewPoint(69.05682, 42.907578), geom.NewPoint(-48.18782, 86.32094)),
	})
	polygon := v.PolygonFrom(geom.NewPoint(35.22194, 56.027016))
	area := polygonArea(polygon)
	c.True(math.Abs(area-2632.7) <= 5, "area %v, want about 2632.7 (polygon %v)", area, polygon)
	for _, hidden := range []geom.Point{
		geom.NewPoint(13.169753, 59.869907),
		geom.NewPoint(2.5690413, 61.272205),
		geom.NewPoint(-12.879423, 70.65482),
		geom.NewPoint(-21.947197, 52.064766),
	} {
		c.False(pointInPolygon(polygon, hidden), "%v should be hidden (polygon %v)", hidden, polygon)
	}

	// Under-report: the eye is 0.23 from the top edge of a sliver bounds, and the visible region below the third wall
	// used to be cut away.
	v = visibility.New(geom.NewRect(-76, -54, 78, 5), []geom.Line{
		geom.NewLine(geom.NewPoint(-52.016544, -51.661396), geom.NewPoint(-71.98019, -51.785614)),
		geom.NewLine(geom.NewPoint(16.898869, -52.708485), geom.NewPoint(15.464588, -50.48328)),
		geom.NewLine(geom.NewPoint(-27.676712, -48.8627), geom.NewPoint(2.0002384, -51.862816)),
		geom.NewLine(geom.NewPoint(12.214455, -53.421143), geom.NewPoint(16.646137, -54.27759)),
	})
	polygon = v.PolygonFrom(geom.NewPoint(-67.283905, -49.233295))
	c.True(pointInPolygon(polygon, geom.NewPoint(-2.9682868, -52.259693)),
		"visible point missing from the polygon %v", polygon)
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

	// The same holds when the zero-length line lies on another line: it must be dropped without cutting that line.
	lines = []geom.Line{
		geom.NewLine(geom.NewPoint(0, 0), geom.NewPoint(10, 0)),
		geom.NewLine(geom.NewPoint(5, 0), geom.NewPoint(5, 0)),
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
