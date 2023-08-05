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
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

// Polygon holds one or more contour lines. The polygon may contain holes and may be self-intersecting.
type Polygon[T constraints.Float] []Contour[T]

// Clone returns a duplicate of this polygon.
func (p Polygon[T]) Clone() Polygon[T] {
	clone := Polygon[T](make([]Contour[T], len(p)))
	for i := range p {
		clone[i] = p[i].Clone()
	}
	return clone
}

// Bounds returns the bounding rectangle of this polygon.
func (p Polygon[T]) Bounds() geom.Rect[T] {
	if len(p) == 0 {
		return geom.Rect[T]{}
	}
	b := p[0].Bounds()
	for _, c := range p[1:] {
		b.Union(c.Bounds())
	}
	return b
}

// Contains returns true if the point is contained by the polygon.
func (p Polygon[T]) Contains(pt geom.Point[T]) bool {
	for i := range p {
		if p[i].Contains(pt) {
			return true
		}
	}
	return false
}

// ContainsEvenOdd returns true if the point is contained by the polygon using the even-odd rule.
// https://en.wikipedia.org/wiki/Even-odd_rule
func (p Polygon[T]) ContainsEvenOdd(pt geom.Point[T]) bool {
	var count int
	for i := range p {
		if p[i].Contains(pt) {
			count++
		}
	}
	return count%2 == 1
}
