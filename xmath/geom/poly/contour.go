package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

type segment struct {
	start geom.Point
	end   geom.Point
}

// Contour is a sequence of vertices connected by line segments, forming a
// closed shape.
type Contour []geom.Point

// Clone returns a copy of this contour.
func (c Contour) Clone() Contour {
	return append([]geom.Point{}, c...)
}

// Bounds returns the bounding rectangle of a contour.
func (c Contour) Bounds() geom.Rect {
	minX := math.Inf(1)
	minY := minX
	maxX := math.Inf(-1)
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
			Width:  maxX - minX,
			Height: maxY - minY,
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

func (c Contour) segment(index int) segment {
	right := 0
	if index != len(c)-1 {
		right = index + 1
	}
	return segment{c[index], c[right]}
}
