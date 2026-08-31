// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/xmath"
)

// LineBoundsEpsilon is the amount by which to expand the bounds of a line to account for floating-point imprecision.
var LineBoundsEpsilon float32 = 0.0001

// Line holds the start and end points of a line.
type Line struct {
	Start Point
	End   Point
}

// NewLine creates a new line from the specified start and end points.
func NewLine(start, end Point) Line {
	return Line{
		Start: start,
		End:   end,
	}
}

// Intersects returns true if this line intersects the other line.
func (l Line) Intersects(other Line) bool {
	return len(l.Intersection(other)) > 0
}

// Intersection determines the intersection of this line with the other line. A return of no points indicates no
// intersection. One point indicates intersection at a single point. Two points indicates an overlapping line segment.
func (l Line) Intersection(other Line) []Point {
	return LineIntersection(l.Start, l.End, other.Start, other.End)
}

// DistanceToPoint returns the distance from the point to the nearest point on this line segment, which is 0 if the
// point lies on the segment.
func (l Line) DistanceToPoint(pt Point) float32 {
	return PointSegmentDistance(l.Start, l.End, pt)
}

// DistanceToPointSquared returns the square of the distance from the point to the nearest point on this line segment,
// which is 0 if the point lies on the segment.
func (l Line) DistanceToPointSquared(pt Point) float32 {
	return PointSegmentDistanceSquared(l.Start, l.End, pt)
}

// lineBoundsEpsilonRatio scales the bounds expansion with the magnitude of the line's coordinates. A fixed expansion
// is rounded away entirely once it drops below half a float32 ULP of the coordinate it is added to (LineBoundsEpsilon
// vanishes at magnitudes of 4096 and up), which would leave an axis-aligned line with a degenerate zero-thickness
// bounds that intersects nothing, not even itself. The ratio is about eight float32 ULPs.
const lineBoundsEpsilonRatio = 1e-6

// Bounds returns the bounding rectangle of this Line, expanded on each side to compensate for floating-point
// imprecision. The expansion is LineBoundsEpsilon or lineBoundsEpsilonRatio times the largest coordinate magnitude,
// whichever is greater, so that it survives float32 rounding at every scale.
func (l Line) Bounds() Rect {
	minX := min(l.Start.X, l.End.X)
	minY := min(l.Start.Y, l.End.Y)
	maxX := max(l.Start.X, l.End.X)
	maxY := max(l.Start.Y, l.End.Y)
	eps := max(LineBoundsEpsilon, max(-minX, maxX, -minY, maxY)*lineBoundsEpsilonRatio)
	return NewRect(minX-eps, minY-eps, maxX-minX+eps*2, maxY-minY+eps*2)
}

// LineIntersection determines the intersection of two lines, if any. A return of no points indicates no intersection.
// One point indicates intersection at a single point. Two points indicates an overlapping line segment.
func LineIntersection(a1, a2, b1, b2 Point) []Point {
	aIsPt := a1.X == a2.X && a1.Y == a2.Y
	bIsPt := b1.X == b2.X && b1.Y == b2.Y
	switch {
	case aIsPt && bIsPt:
		if a1.X == b1.X && a1.Y == b1.Y {
			return []Point{a1}
		}
	case aIsPt:
		if PointSegmentDistance(b1, b2, a1) == 0 {
			return []Point{a1}
		}
	case bIsPt:
		if PointSegmentDistance(a1, a2, b1) == 0 {
			return []Point{b1}
		}
	default:
		// The arithmetic is done in float64. The products of float32 differences carry at most 48 significant bits,
		// so they are exact in float64 and the exact-zero collinearity tests below become reliable; in float32, fused
		// multiply-subtract contraction of a*b - c*d (which the compiler emits on arm64) returns the rounding error of
		// the second product instead of zero for exactly collinear segments, so collinear overlaps went undetected. It
		// also keeps an interpolated crossing within a float32 rounding of the true crossing when one segment is far
		// longer than the other. The explicit float64 conversions around each product keep the compiler from fusing
		// even the float64 subtractions, which matters for differences that do not happen to be exact.
		abdx := float64(a1.X) - float64(b1.X)
		abdy := float64(a1.Y) - float64(b1.Y)
		bdx := float64(b2.X) - float64(b1.X)
		bdy := float64(b2.Y) - float64(b1.Y)
		uat := float64(bdx*abdy) - float64(bdy*abdx)
		adx := float64(a2.X) - float64(a1.X)
		ady := float64(a2.Y) - float64(a1.Y)
		ubt := float64(adx*abdy) - float64(ady*abdx)
		ub := float64(bdy*adx) - float64(bdx*ady)
		if ub != 0 {
			// Not parallel, so find intersection point
			a := uat / ub
			if a >= 0 && a <= 1 {
				b := ubt / ub
				if b >= 0 && b <= 1 {
					return []Point{
						{
							X: float32(float64(a1.X) + a*adx),
							Y: float32(float64(a1.Y) + a*ady),
						},
					}
				}
			}
		} else if uat == 0 && ubt == 0 {
			// Parallel, so check for overlap. Collinearity requires both cross-product numerators to be zero; in exact
			// arithmetic either being zero implies the other, but requiring both guards against a phantom overlap when
			// rounding zeroes only one of them for parallel-but-offset segments.
			var ub1, ub2 float64
			if math.Abs(adx) > math.Abs(ady) {
				ub1 = (float64(b1.X) - float64(a1.X)) / adx
				ub2 = (float64(b2.X) - float64(a1.X)) / adx
			} else {
				ub1 = (float64(b1.Y) - float64(a1.Y)) / ady
				ub2 = (float64(b2.Y) - float64(a1.Y)) / ady
			}
			left := max(0, min(ub1, ub2))
			right := min(1, max(ub1, ub2))
			if left > right {
				return nil
			}
			if left == right {
				return []Point{
					{
						X: float32(float64(a2.X)*left + float64(a1.X)*(1-left)),
						Y: float32(float64(a2.Y)*left + float64(a1.Y)*(1-left)),
					},
				}
			}
			return []Point{
				{
					X: float32(float64(a2.X)*left + float64(a1.X)*(1-left)),
					Y: float32(float64(a2.Y)*left + float64(a1.Y)*(1-left)),
				},
				{
					X: float32(float64(a2.X)*right + float64(a1.X)*(1-right)),
					Y: float32(float64(a2.Y)*right + float64(a1.Y)*(1-right)),
				},
			}
		}
	}
	return nil
}

// PointSegmentDistance returns the distance from p to the nearest point on the segment from s1 to s2, which is 0 if p
// lies on the segment.
func PointSegmentDistance(s1, s2, p Point) float32 {
	return xmath.Sqrt(PointSegmentDistanceSquared(s1, s2, p))
}

// PointSegmentDistanceSquared returns the square of the distance from p to the nearest point on the segment from s1 to
// s2, which is 0 if p lies on the segment. The math is done in float64 via the projection of p onto the segment:
// formulations that subtract two nearly equal float32 squared magnitudes cancel catastrophically for a point near a
// long segment, silently reporting zero for true distances up to roughly the segment length times 2.4e-4.
func PointSegmentDistanceSquared(s1, s2, p Point) float32 {
	vx := float64(s2.X) - float64(s1.X)
	vy := float64(s2.Y) - float64(s1.Y)
	px := float64(p.X) - float64(s1.X)
	py := float64(p.Y) - float64(s1.Y)
	t := 0.0
	if lenSqrd := vx*vx + vy*vy; lenSqrd > 0 {
		t = min(max((px*vx+py*vy)/lenSqrd, 0), 1)
	}
	dx := px - t*vx
	dy := py - t*vy
	return float32(dx*dx + dy*dy)
}
