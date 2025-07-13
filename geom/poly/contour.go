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
	"math"
	"strings"

	"github.com/richardwilkes/toolbox/v2/geom"
)

// Contour is a sequence of vertices connected by line segments, forming a closed shape.
type Contour []geom.Point

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
func (c Contour) Bounds() geom.Rect {
	if len(c) == 0 {
		return geom.Rect{}
	}
	minX := float32(math.MaxFloat32)
	minY := minX
	maxX := float32(-math.MaxFloat32)
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
	return geom.NewRect(minX, minY, 1+maxX-minX, 1+maxY-minY)
}

// Contains returns true if the point is contained by the contour.
func (c Contour) Contains(pt geom.Point) bool {
	var count int
	for i := range c {
		cur := c[i]
		bottom := cur
		next := c[(i+1)%len(c)]
		top := next
		if bottom.Y > top.Y {
			bottom, top = top, bottom
		}
		if pt.Y >= bottom.Y && pt.Y < top.Y && pt.X < max(cur.X, next.X) && next.Y != cur.Y &&
			(cur.X == next.X || pt.X <= (pt.Y-cur.Y)*(next.X-cur.X)/(next.Y-cur.Y)+cur.X) {
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
