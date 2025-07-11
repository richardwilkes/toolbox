// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"strings"

	"github.com/richardwilkes/toolbox/v2/geom"
)

// Contour is a sequence of vertices connected by line segments, forming a closed shape.
type Contour []Point

// Points converts this Contour into a slice of geom.Point.
func (c Contour) Points() []geom.Point {
	points := make([]geom.Point, len(c))
	for i, p := range c {
		points[i] = p.Point()
	}
	return points
}

// Clone returns a copy of this contour.
func (c Contour) Clone() Contour {
	if len(c) == 0 {
		return nil
	}
	clone := make(Contour, len(c))
	copy(clone, c)
	return clone
}

// Bounds returns the bounding rectangle of the contour.
func (c Contour) Bounds() Rect {
	if len(c) == 0 {
		return Rect{}
	}
	minX := c[0].X
	minY := c[0].Y
	maxX := minX
	maxY := minY
	for _, p := range c[1:] {
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
	return NewRect(minX, minY, maxX-minX, maxY-minY)
}

// Contains returns true if the point is contained by the contour.
func (c Contour) Contains(pt Point) bool {
	if len(c) < 3 {
		return false // A contour needs at least 3 points to contain anything
	}
	var count int
	for i := range c {
		cur := c[i]
		bottom := cur
		next := c[(i+1)%len(c)]
		top := next
		if bottom.Y > top.Y {
			bottom, top = top, bottom
		}
		if pt.Y >= bottom.Y &&
			pt.Y < top.Y &&
			pt.X < max(cur.X, next.X) &&
			next.Y != cur.Y &&
			(cur.X == next.X ||
				pt.X <= (pt.Y-cur.Y).Mul(next.X-cur.X).Div(next.Y-cur.Y)+cur.X) {
			count++
		}
	}
	return count%2 == 1
}

func (c Contour) String() string {
	var buffer strings.Builder
	buffer.WriteByte('{')
	for j, pt := range c {
		if j != 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteByte('{')
		buffer.WriteString(pt.String())
		buffer.WriteByte('}')
	}
	buffer.WriteByte('}')
	return buffer.String()
}
