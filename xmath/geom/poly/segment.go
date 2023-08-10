// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

// Segment32 is an alias for the float32 version of Segment.
type Segment32 = Segment[float32]

// Segment64 is an alias for the float64 version of Segment.
type Segment64 = Segment[float64]

// Segment holds the start and end point of a single segment within a contour of a polygon.
type Segment[T constraints.Float] struct {
	Start geom.Point[T]
	End   geom.Point[T]
}

// CoincidesWith returns true if the segment matches the given start and end in either direction.
func (s Segment[T]) CoincidesWith(start, end geom.Point[T]) bool {
	return s.Start == start && s.End == end || s.Start == end && s.End == start
}

// FindIntersection finds the points where the two segments intersect, if any.
func (s Segment[T]) FindIntersection(other Segment[T], tryBothDirections bool) (count int, pi0, pi1 geom.Point[T]) {
	p0 := s.Start
	d0 := geom.Point[T]{
		X: s.End.X - p0.X,
		Y: s.End.Y - p0.Y,
	}
	p1 := other.Start
	d1 := geom.Point[T]{
		X: other.End.X - p1.X,
		Y: other.End.Y - p1.Y,
	}
	const sqrEpsilon = 1e-21
	e := geom.Point[T]{
		X: p1.X - p0.X,
		Y: p1.Y - p0.Y,
	}
	cross := d0.X*d1.Y - d0.Y*d1.X
	sqrCross := cross * cross
	sqrLen0 := xmath.Sqrt(d0.X*d0.X + d0.Y*d0.Y)
	sqrLen1 := xmath.Sqrt(d1.X*d1.X + d1.Y*d1.Y)
	if sqrCross > sqrEpsilon*sqrLen0*sqrLen1 {
		s0 := (e.X*d1.Y - e.Y*d1.X) / cross
		if s0 < 0 || s0 > 1 {
			return 0, geom.Point[T]{}, geom.Point[T]{}
		}
		t := (e.X*d0.Y - e.Y*d0.X) / cross
		if t < 0 || t > 1 {
			return 0, geom.Point[T]{}, geom.Point[T]{}
		}
		pi0.X = p0.X + s0*d0.X
		pi0.Y = p0.Y + s0*d0.Y
		return 1, pi0, pi1
	}
	sqrLenE := xmath.Sqrt(e.X*e.X + e.Y*e.Y)
	cross = e.X*d0.Y - e.Y*d0.X
	sqrCross = cross * cross
	if sqrCross > sqrEpsilon*sqrLen0*sqrLenE {
		return 0, pi0, pi1
	}
	s0 := (d0.X*e.X + d0.Y*e.Y) / sqrLen0
	s1 := s0 + (d0.X*d1.X+d0.Y*d1.Y)/sqrLen0
	smin := min(s0, s1)
	smax := max(s0, s1)
	w := make([]T, 0, 2)
	switch {
	case smin > 1 || smax < 0:
	case smin == 1:
		w = append(w, 1)
	case smax == 0:
		w = append(w, 0)
	default:
		if smin > 0 {
			w = append(w, smin)
		} else {
			w = append(w, 0)
		}
		if smax < 1 {
			w = append(w, smax)
		} else {
			w = append(w, 1)
		}
	}
	if len(w) > 0 {
		pi0.X = p0.X + w[0]*d0.X
		pi0.Y = p0.Y + w[0]*d0.Y
	}
	if len(w) > 1 {
		pi1.X = p0.X + w[1]*d0.X
		pi1.Y = p0.Y + w[1]*d0.Y
	} else if tryBothDirections {
		if imax, otherPi0, otherPi1 := other.FindIntersection(s, false); imax > len(w) {
			return imax, otherPi0, otherPi1
		}
	}
	return len(w), pi0, pi1
}
