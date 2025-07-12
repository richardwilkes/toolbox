// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"cmp"
	"slices"

	"github.com/richardwilkes/toolbox/v2/collection/quadtree"
	"github.com/richardwilkes/toolbox/v2/geom/poly"
)

var (
	epsilon    = poly.NumFromFloat(0.01)
	twoEpsilon = epsilon.Mul(poly.Two)
	tenEpsilon = epsilon.Mul(poly.Ten)
)

// Visibility holds state for computing a visibility polygon.
type Visibility struct {
	segments []Segment
	top      poly.Num
	left     poly.Num
	bottom   poly.Num
	right    poly.Num
}

// New creates a Visibility object. If the obstructions do not intersect each other, pass in true for hasNoIntersections
// to eliminate the costly pass to break the segments up into smaller parts.
func New(bounds poly.Rect, obstructions []Segment, hasNoIntersections bool) *Visibility {
	v := &Visibility{
		top:      bounds.Y,
		left:     bounds.X,
		bottom:   bounds.Bottom(),
		right:    bounds.Right(),
		segments: make([]Segment, len(obstructions)),
	}
	copy(v.segments, obstructions)
	if !hasNoIntersections {
		v.breakIntersections()
	}
	return v
}

// SetViewPoint sets a view point and generates a polygon with the unobstructed visible area.
func (v *Visibility) SetViewPoint(viewPt poly.Point) poly.Polygon {
	// If the view point is not within the bounding area, there is no visible area
	if !v.inViewport(viewPt) {
		return nil
	}

	// Generate a revised segment list by clipping the segments against the viewport and throwing out any that aren't
	// within the viewport.
	segments := make([]Segment, 0, len(v.segments)*2)
	viewport := []poly.Point{
		poly.NewPoint(v.left, v.top),
		poly.NewPoint(v.right, v.top),
		poly.NewPoint(v.right, v.bottom),
		poly.NewPoint(v.left, v.bottom),
	}
	for _, si := range v.segments {
		if (si.Start.X < v.left && si.End.X < v.left) ||
			(si.Start.Y < v.top && si.End.Y < v.top) ||
			(si.Start.X > v.right && si.End.X > v.right) ||
			(si.Start.Y > v.bottom && si.End.Y > v.bottom) {
			continue
		}
		intersections := make([]poly.Point, 0, len(viewport))
		for j := range viewport {
			k := (j + 1) % len(viewport)
			if hasIntersection(si.Start, si.End, viewport[j], viewport[k]) {
				pt, intersects := intersectLines(si.Start, si.End, viewport[j], viewport[k])
				if intersects && pt != si.Start && pt != si.End {
					intersections = append(intersections, pt)
				}
			}
		}
		segments = v.collectSegments(si, intersections, segments, true)
	}

	// Add the viewport bounds to the segment list
	topLeft := poly.Point{X: v.left - tenEpsilon, Y: v.top - tenEpsilon}
	topRight := poly.Point{X: v.right + tenEpsilon, Y: v.top - tenEpsilon}
	bottomLeft := poly.Point{X: v.left - tenEpsilon, Y: v.bottom + tenEpsilon}
	bottomRight := poly.Point{X: v.right + tenEpsilon, Y: v.bottom + tenEpsilon}
	segments = append(segments,
		Segment{Start: topLeft, End: topRight},
		Segment{Start: topRight, End: bottomRight},
		Segment{Start: bottomRight, End: bottomLeft},
		Segment{Start: bottomLeft, End: topLeft},
	)

	return v.computePolygon(viewPt, segments)
}

func (v *Visibility) computePolygon(viewPt poly.Point, segments []Segment) poly.Polygon {
	// Sweep through the points to generate the visibility contour
	sorted := sortPoints(viewPt, segments)
	mapper := &array{data: make([]int, len(segments))}
	for i := range mapper.data {
		mapper.data[i] = -1
	}
	heap := &array{}
	start := poly.Point{X: viewPt.X + poly.One, Y: viewPt.Y}
	for i := range segments {
		a1 := angle(segments[i].Start, viewPt)
		a2 := angle(segments[i].End, viewPt)
		if (a1 >= -poly.OneHundredEighty && a1 <= 0 && a2 <= poly.OneHundredEighty &&
			a2 >= 0 && a2-a1 > poly.OneHundredEighty) ||
			(a2 >= -poly.OneHundredEighty && a2 <= 0 && a1 <= poly.OneHundredEighty &&
				a1 >= 0 && a1-a2 > poly.OneHundredEighty) {
			insert(i, heap, mapper, segments, viewPt, start)
		}
	}
	contour := make(poly.Contour, 0, len(sorted)*2)
	i := 0
	for i < len(sorted) {
		extend := false
		shorten := false
		orig := i
		vertex := sorted[i].pt(segments)
		oldSeg := heap.elem(0)
		for {
			if mapper.elem(sorted[i].segmentIndex) != -1 {
				if sorted[i].segmentIndex == oldSeg {
					extend = true
					vertex = sorted[i].pt(segments)
				}
				remove(mapper.elem(sorted[i].segmentIndex), heap, mapper, segments, viewPt, vertex)
			} else {
				insert(sorted[i].segmentIndex, heap, mapper, segments, viewPt, vertex)
				if heap.elem(0) != oldSeg {
					shorten = true
				}
			}
			i++
			if i == len(sorted) || sorted[i].angle >= sorted[orig].angle+epsilon {
				break
			}
		}
		if extend {
			contour = append(contour, poly.Point{X: vertex.X, Y: vertex.Y})
			s := segments[heap.elem(0)]
			if cur, intersects := intersectLines(s.Start, s.End, viewPt, vertex); intersects && cur != vertex {
				contour = append(contour, poly.Point{X: cur.X, Y: cur.Y})
			}
		} else if shorten {
			s := segments[oldSeg]
			if cur, intersects := intersectLines(s.Start, s.End, viewPt, vertex); intersects {
				contour = append(contour, poly.Point{X: cur.X, Y: cur.Y})
			}
			s = segments[heap.elem(0)]
			if cur, intersects := intersectLines(s.Start, s.End, viewPt, vertex); intersects {
				contour = append(contour, poly.Point{X: cur.X, Y: cur.Y})
			}
		}
	}
	if len(contour) == 0 {
		return nil
	}
	return poly.Polygon{contour}
}

func (v *Visibility) inViewport(pt poly.Point) bool {
	return pt.X >= v.left-epsilon &&
		pt.Y >= v.top-epsilon &&
		pt.X <= v.right+epsilon &&
		pt.Y <= v.bottom+epsilon
}

func (v *Visibility) breakIntersections() {
	var qt quadtree.QuadTree[Segment]
	for _, si := range v.segments {
		qt.Insert(si)
	}
	segments := make([]Segment, 0, len(v.segments)*2)
	for _, si := range v.segments {
		var intersections []poly.Point
		for _, sj := range qt.FindIntersects(si.Bounds()) {
			if si == sj {
				continue
			}
			if hasIntersection(si.Start, si.End, sj.Start, sj.End) {
				pt, intersects := intersectLines(si.Start, si.End, sj.Start, sj.End)
				if intersects && pt != si.Start && pt != si.End {
					intersections = append(intersections, pt)
				}
			}
		}
		segments = v.collectSegments(si, intersections, segments, false)
	}
	v.segments = slices.Clip(segments)
}

func (v *Visibility) collectSegments(s Segment, intersections []poly.Point, segments []Segment, onlyInViewPort bool) []Segment {
	start := s.Start
	for len(intersections) > 0 {
		endIndex := 0
		endDis := distance(start, intersections[0])
		for i := 1; i < len(intersections); i++ {
			if dis := distance(start, intersections[i]); dis < endDis {
				endDis = dis
				endIndex = i
			}
		}
		if !onlyInViewPort || (v.inViewport(start) && v.inViewport(intersections[endIndex])) {
			segments = append(segments, Segment{Start: start, End: intersections[endIndex]})
		}
		start = intersections[endIndex]
		intersections = slices.Delete(intersections, endIndex, endIndex+1)
	}
	if !onlyInViewPort || (v.inViewport(start) && v.inViewport(s.End)) {
		segments = append(segments, Segment{Start: start, End: s.End})
	}
	return segments
}

func remove(index int, heap, mapper *array, segments []Segment, position, destination poly.Point) {
	mapper.set(heap.elem(index), -1)
	if index == heap.size()-1 {
		heap.pop()
		return
	}
	heap.set(index, heap.pop())
	mapper.set(heap.elem(index), index)
	cur := index
	parent := (cur - 1) / 2
	if cur != 0 && lessThan(heap.elem(cur), heap.elem(parent), segments, position, destination) {
		for cur > 0 {
			parent = (cur - 1) / 2
			parentElem := heap.elem(parent)
			curElem := heap.elem(cur)
			if !lessThan(curElem, parentElem, segments, position, destination) {
				break
			}
			mapper.set(parentElem, cur)
			mapper.set(curElem, parent)
			heap.set(cur, parentElem)
			heap.set(parent, curElem)
			cur = parent
		}
		return
	}
loop:
	for {
		left := 2*cur + 1
		right := left + 1
		switch {
		case left < heap.size() && lessThan(heap.elem(left), heap.elem(cur), segments, position, destination) &&
			(right == heap.size() || lessThan(heap.elem(left), heap.elem(right), segments, position, destination)):
			leftElem := heap.elem(left)
			curElem := heap.elem(cur)
			mapper.set(leftElem, cur)
			mapper.set(curElem, left)
			heap.set(left, curElem)
			heap.set(cur, leftElem)
			cur = left
		case right < heap.size() && lessThan(heap.elem(right), heap.elem(cur), segments, position, destination):
			rightElem := heap.elem(right)
			curElem := heap.elem(cur)
			mapper.set(rightElem, cur)
			mapper.set(curElem, right)
			heap.set(right, curElem)
			heap.set(cur, rightElem)
			cur = right
		default:
			break loop
		}
	}
}

func insert(index int, heap, mapper *array, segments []Segment, position, destination poly.Point) {
	if _, intersects := intersectLines(segments[index].Start, segments[index].End, position, destination); !intersects {
		return
	}
	cur := heap.size()
	heap.push(index)
	mapper.set(index, cur)
	for cur > 0 {
		parent := (cur - 1) / 2
		parentElem := heap.elem(parent)
		curElem := heap.elem(cur)
		if !lessThan(curElem, parentElem, segments, position, destination) {
			break
		}
		mapper.set(parentElem, cur)
		mapper.set(curElem, parent)
		heap.set(cur, parentElem)
		heap.set(parent, curElem)
		cur = parent
	}
}

func lessThan(index1, index2 int, segments []Segment, position, destination poly.Point) bool {
	pt1, intersects1 := intersectLines(segments[index1].Start, segments[index1].End, position, destination)
	if !intersects1 {
		return false
	}
	pt2, intersects2 := intersectLines(segments[index2].Start, segments[index2].End, position, destination)
	if !intersects2 {
		return false
	}
	if pt1 != pt2 {
		d1 := distance(pt1, position)
		d2 := distance(pt2, position)
		return d1 < d2
	}
	var a1 poly.Num
	if pt1 == segments[index1].Start {
		a1 = angle2(segments[index1].End, pt1, position)
	} else {
		a1 = angle2(segments[index1].Start, pt1, position)
	}
	var a2 poly.Num
	if pt2 == segments[index2].Start {
		a2 = angle2(segments[index2].End, pt2, position)
	} else {
		a2 = angle2(segments[index2].Start, pt2, position)
	}
	if a1 < poly.OneHundredEighty {
		if a2 > poly.OneHundredEighty {
			return true
		}
		return a2 < a1
	}
	return a1 < a2
}

func sortPoints(position poly.Point, segments []Segment) []endPoint {
	points := make([]endPoint, len(segments)*2)
	pos := 0
	for i, s := range segments {
		points[pos].segmentIndex = i
		points[pos].angle = angle(s.Start, position)
		points[pos].start = true
		pos++
		points[pos].segmentIndex = i
		points[pos].angle = angle(s.End, position)
		points[pos].start = false
		pos++
	}
	slices.SortFunc(points, func(a, b endPoint) int {
		if result := cmp.Compare(a.angle, b.angle); result != 0 {
			return result
		}
		if result := cmp.Compare(distance(a.pt(segments), position), distance(b.pt(segments), position)); result != 0 {
			return result
		}
		if a.start == b.start {
			return 0
		}
		if a.start {
			return 1
		}
		return -1
	})
	return points
}

func angle2(a, b, c poly.Point) poly.Num {
	a1 := angle(a, b)
	a2 := angle(b, c)
	a3 := a1 - a2
	if a3 < 0 {
		a3 += poly.ThreeHundredSixty
	}
	if a3 > poly.ThreeHundredSixty {
		a3 -= poly.ThreeHundredSixty
	}
	return a3
}

func angle(a, b poly.Point) poly.Num {
	return poly.Atan2(b.Y-a.Y, b.X-a.X).Mul(poly.OneHundredEighty).Div(poly.Pi)
}

func distance(a, b poly.Point) poly.Num {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx.Mul(dx) + dy.Mul(dy)
}

func intersectLines(s1, e1, s2, e2 poly.Point) (poly.Point, bool) {
	dbx := e2.X - s2.X
	dby := e2.Y - s2.Y
	dax := e1.X - s1.X
	day := e1.Y - s1.Y
	ub := dby.Mul(dax) - dbx.Mul(day)
	if ub == 0 {
		return poly.Point{}, false
	}
	ua := (dbx.Mul(s1.Y-s2.Y) - dby.Mul(s1.X-s2.X)).Div(ub)
	return poly.Point{X: s1.X + ua.Mul(dax), Y: s1.Y + ua.Mul(day)}, true
}

func hasIntersection(s1, e1, s2, e2 poly.Point) bool {
	d1 := direction(s2, e2, s1)
	d2 := direction(s2, e2, e1)
	d3 := direction(s1, e1, s2)
	d4 := direction(s1, e1, e2)
	return (((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0))) ||
		(d1 == 0 && onSegment(s2, e2, s1)) ||
		(d2 == 0 && onSegment(s2, e2, e1)) ||
		(d3 == 0 && onSegment(s1, e1, s2)) ||
		(d4 == 0 && onSegment(s1, e1, e2))
}

func direction(a, b, c poly.Point) int {
	return cmp.Compare((c.X - a.X).Mul(b.Y-a.Y), (b.X - a.X).Mul(c.Y-a.Y))
}

func onSegment(a, b, c poly.Point) bool {
	return (a.X <= c.X || b.X <= c.X) &&
		(c.X <= a.X || c.X <= b.X) &&
		(a.Y <= c.Y || b.Y <= c.Y) &&
		(c.Y <= a.Y || c.Y <= b.Y)
}
