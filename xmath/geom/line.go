// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import "math"

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
		abdx := a1.X - b1.X
		abdy := a1.Y - b1.Y
		bdx := b2.X - b1.X
		bdy := b2.Y - b1.Y
		uat := bdx*abdy - bdy*abdx
		adx := a2.X - a1.X
		ady := a2.Y - a1.Y
		ubt := adx*abdy - ady*abdx
		ub := bdy*adx - bdx*ady
		if ub != 0 {
			// Not parallel, so find intersection point
			a := uat / ub
			if a >= 0 && a <= 1 {
				b := ubt / ub
				if b >= 0 && b <= 1 {
					return []Point{{X: a1.X + a*adx, Y: a1.Y + a*ady}}
				}
			}
		} else if uat == 0 || ubt == 0 {
			// Parallel, so check for overlap
			var ub1, ub2 float64
			if math.Abs(adx) > math.Abs(ady) {
				ub1 = (b1.X - a1.X) / adx
				ub2 = (b2.X - a1.X) / adx
			} else {
				ub1 = (b1.Y - a1.Y) / ady
				ub2 = (b2.Y - a1.Y) / ady
			}
			left := math.Max(0, math.Min(ub1, ub2))
			right := math.Min(1, math.Max(ub1, ub2))
			if left < right {
				return []Point{
					{X: a2.X*left + a1.X*(1-left), Y: a2.Y*left + a1.Y*(1-left)},
					{X: a2.X*right + a1.X*(1-right), Y: a2.Y*right + a1.Y*(1-right)},
				}
			} else if left == right {
				return []Point{
					{X: a2.X*left + a1.X*(1-left), Y: a2.Y*left + a1.Y*(1-left)},
				}
			}
		}
	}
	return nil
}

// PointSegmentDistance returns the distance from a point to a line segment. The distance measured is the distance
// between the specified point and the closest point between the specified end points. If the specified point intersects
// the line segment in between the end points, this function returns 0.
func PointSegmentDistance(s1, s2, p Point) float64 {
	return math.Sqrt(PointSegmentDistanceSquared(s1, s2, p))
}

// PointSegmentDistanceSquared returns the square of the distance from a point to a line segment. The distance measured
// is the distance between the specified point and the closest point between the specified end points. If the specified
// point intersects the line segment in between the end points, this function returns 0.
func PointSegmentDistanceSquared(s1, s2, p Point) float64 {
	vx := s2.X - s1.X
	vy := s2.Y - s1.Y
	px := p.X - s1.X
	py := p.Y - s1.Y
	dp := px*vx + py*vy
	var projected float64
	if dp > 0 {
		px = vx - px
		py = vy - py
		dp = px*vx + py*vy
		if dp > 0 {
			projected = dp * dp / (vx*vx + vy*vy)
		}
	}
	lenSq := px*px + py*py - projected
	if lenSq < 0 {
		lenSq = 0
	}
	return lenSq
}
