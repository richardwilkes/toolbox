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

// Contour is a sequence of vertices connected by line segments, forming a closed shape.
type Contour[T constraints.Float] []geom.Point[T]

// Clone returns a copy of this contour.
func (c Contour[T]) Clone() Contour[T] {
	return append([]geom.Point[T]{}, c...)
}

// Bounds returns the bounding rectangle of a contour.
func (c Contour[T]) Bounds() geom.Rect[T] {
	if len(c) == 0 {
		return geom.Rect[T]{}
	}
	minX := xmath.MaxValue[T]()
	minY := minX
	maxX := xmath.MinValue[T]()
	maxY := maxX
	for _, p := range c {
		if p.X > maxX {
			maxX = p.X
		}
		if p.X < minX {
			minX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.Y < minY {
			minY = p.Y
		}
	}
	return geom.Rect[T]{
		Point: geom.Point[T]{
			X: minX,
			Y: minY,
		},
		Size: geom.Size[T]{
			Width:  1 + maxX - minX,
			Height: 1 + maxY - minY,
		},
	}
}

// Contains returns true if the point is contained by the contour.
func (c Contour[T]) Contains(pt geom.Point[T]) bool {
	var count int
	for i := range c {
		cur := c[i]
		bottom := cur
		n := i + 1
		if n == len(c) {
			n = 0
		}
		next := c[n]
		top := next
		if bottom.Y > top.Y {
			bottom, top = top, bottom
		}
		if pt.Y >= bottom.Y && pt.Y < top.Y && pt.X < xmath.Max(cur.X, next.X) && next.Y != cur.Y &&
			(cur.X == next.X || pt.X <= (pt.Y-cur.Y)*(next.X-cur.X)/(next.Y-cur.Y)+cur.X) {
			count++
		}
	}
	return count%2 == 1
}

func (c Contour[T]) segment(index int) Segment[T] {
	if index == len(c)-1 {
		return Segment[T]{c[len(c)-1], c[0]}
	}
	return Segment[T]{c[index], c[index+1]}
}
