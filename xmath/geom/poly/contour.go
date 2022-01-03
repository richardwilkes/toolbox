// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

// Contour is a sequence of vertices connected by line segments, forming a closed shape.
type Contour []geom.Point

// Clone returns a copy of this contour.
func (c Contour) Clone() Contour {
	return append([]geom.Point{}, c...)
}

// Bounds returns the bounding rectangle of a contour.
func (c Contour) Bounds() geom.Rect {
	if len(c) == 0 {
		return geom.Rect{}
	}
	minX := math.MaxFloat64
	minY := minX
	maxX := -math.MaxFloat64
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
	return geom.Rect{
		Point: geom.Point{
			X: minX,
			Y: minY,
		},
		Size: geom.Size{
			Width:  1 + maxX - minX,
			Height: 1 + maxY - minY,
		},
	}
}

// Contains returns true if the point is contained by the contour.
func (c Contour) Contains(pt geom.Point) bool {
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
		if pt.Y >= bottom.Y && pt.Y < top.Y && pt.X < math.Max(cur.X, next.X) && next.Y != cur.Y &&
			(cur.X == next.X || pt.X <= (pt.Y-cur.Y)*(next.X-cur.X)/(next.Y-cur.Y)+cur.X) {
			count++
		}
	}
	return count%2 == 1
}
