// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"math"
	"slices"

	"github.com/richardwilkes/toolbox/v2/collection/quadtree"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/xmath"
)

const (
	// pointEpsilonRatio is the fraction of a scene's overall magnitude used as the tolerance when comparing points for
	// equality. Deriving the tolerance from the data rather than fixing it at some number of world units keeps the
	// results consistent as a scene is scaled up or down.
	pointEpsilonRatio = 0.0001
	// minPointEpsilon is the floor for the point comparison tolerance, used when a scene has no measurable extent.
	minPointEpsilon = 1e-6
	// angleEpsilon is the tolerance, in degrees, for treating two endpoints as lying at the same angle from the view
	// point. This is an angular tolerance and is deliberately kept separate from the point comparison tolerance, which
	// is in world units.
	angleEpsilon = 0.01
)

// Visibility holds state for computing a visibility polygon.
type Visibility struct {
	lines   []geom.Line
	bounds  geom.Rect
	epsilon float32
}

// New creates a Visibility object. The obstructions must not intersect each other. If they do, call
// BreakIntersections() and pass the result instead.
func New(bounds geom.Rect, obstructions []geom.Line) *Visibility {
	magnitude := max(xmath.Abs(bounds.X), xmath.Abs(bounds.Y), xmath.Abs(bounds.Right()), xmath.Abs(bounds.Bottom()))
	v := &Visibility{
		lines:   make([]geom.Line, len(obstructions)),
		bounds:  bounds,
		epsilon: pointEpsilon(magnitude),
	}
	copy(v.lines, obstructions)
	return v
}

// BreakIntersections breaks the lines at their intersections, returning a new slice of lines that do not intersect.
// Collinear lines that overlap each other are split at the ends of the shared portion, and that shared portion is
// returned only once.
func BreakIntersections(lines []geom.Line) []geom.Line {
	var qt quadtree.QuadTree[geom.Line]
	var magnitude float32
	for _, line := range lines {
		qt.Insert(line)
		magnitude = max(magnitude, xmath.Abs(line.Start.X), xmath.Abs(line.Start.Y), xmath.Abs(line.End.X),
			xmath.Abs(line.End.Y))
	}
	eps := pointEpsilon(magnitude)
	revised := make([]geom.Line, 0, len(lines)*2)
	for _, line := range lines {
		var intersections []geom.Point
		for _, one := range qt.FindIntersects(line.Bounds()) {
			if line == one {
				continue
			}
			// geom.LineIntersection is used rather than the local intersectLines because it returns both endpoints of
			// the shared portion when the two segments are collinear and overlap, a case that cannot be expressed as a
			// single infinite-line intersection point.
			for _, pt := range geom.LineIntersection(line.Start, line.End, one.Start, one.End) {
				if !pt.EqualWithin(line.Start, eps) && !pt.EqualWithin(line.End, eps) {
					intersections = append(intersections, pt)
				}
			}
		}
		revised = collectLines(line, intersections, revised, nil, eps)
	}
	// A collinear overlap is broken out of each of the lines that contained it, so the shared portion appears more than
	// once. Drop those exact duplicates, since coincident lines would otherwise still intersect each other.
	seen := make(map[geom.Line]struct{}, len(revised))
	deduped := revised[:0]
	for _, line := range revised {
		key := line
		if key.End.X < key.Start.X || (key.End.X == key.Start.X && key.End.Y < key.Start.Y) {
			key.Start, key.End = key.End, key.Start
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, line)
	}
	return slices.Clip(deduped)
}

func collectLines(line geom.Line, intersections []geom.Point, lines []geom.Line, viewPort *geom.Rect, eps float32,
) []geom.Line {
	start := line.Start
	for len(intersections) > 0 {
		endIndex := 0
		endDis := distSqrd(start, intersections[0])
		for i := 1; i < len(intersections); i++ {
			if dis := distSqrd(start, intersections[i]); dis < endDis {
				endDis = dis
				endIndex = i
			}
		}
		end := intersections[endIndex]
		intersections = slices.Delete(intersections, endIndex, endIndex+1)
		// Three or more lines meeting at a single point yield that point once per line, so skip any repeat of the
		// current start rather than emitting a zero-length segment.
		if end.EqualWithin(start, eps) {
			continue
		}
		if viewPort == nil || viewPort.IntersectsLine(start, end) {
			lines = append(lines, geom.NewLine(start, end))
		}
		start = end
	}
	if !line.End.EqualWithin(start, eps) && (viewPort == nil || viewPort.IntersectsLine(start, line.End)) {
		lines = append(lines, geom.NewLine(start, line.End))
	}
	return lines
}

// SetViewPoint returns the polygon of the unobstructed area visible from viewPt, or nil if viewPt is outside the bounds
// or no area is visible.
func (v *Visibility) SetViewPoint(viewPt geom.Point) []geom.Point {
	if !viewPt.In(v.bounds) {
		return nil
	}

	// Generate a revised line list by clipping the lines against the viewport and throwing out any that aren't within
	// the viewport.
	lines := make([]geom.Line, 0, len(v.lines)*2)
	viewport := []geom.Point{
		v.bounds.Point,
		v.bounds.TopRight(),
		v.bounds.BottomRight(),
		v.bounds.BottomLeft(),
	}
	for _, line := range v.lines {
		if (line.Start.X < v.bounds.X && line.End.X < v.bounds.X) ||
			(line.Start.Y < v.bounds.Y && line.End.Y < v.bounds.Y) ||
			(line.Start.X > v.bounds.Right() && line.End.X > v.bounds.Right()) ||
			(line.Start.Y > v.bounds.Bottom() && line.End.Y > v.bounds.Bottom()) {
			continue
		}
		intersections := make([]geom.Point, 0, len(viewport))
		for j := range viewport {
			k := (j + 1) % len(viewport)
			if hasIntersection(line.Start, line.End, viewport[j], viewport[k]) {
				pt, intersects := intersectLines(line.Start, line.End, viewport[j], viewport[k])
				if intersects && !pt.EqualWithin(line.Start, v.epsilon) && !pt.EqualWithin(line.End, v.epsilon) {
					intersections = append(intersections, pt)
				}
			}
		}
		lines = collectLines(line, intersections, lines, &v.bounds, v.epsilon)
	}

	lines = append(lines,
		geom.NewLine(v.bounds.Point, v.bounds.TopRight()),
		geom.NewLine(v.bounds.TopRight(), v.bounds.BottomRight()),
		geom.NewLine(v.bounds.BottomRight(), v.bounds.BottomLeft()),
		geom.NewLine(v.bounds.BottomLeft(), v.bounds.Point),
	)

	return v.computePolygon(viewPt, lines)
}

func (v *Visibility) computePolygon(viewPt geom.Point, lines []geom.Line) []geom.Point {
	// Sweep through the points to generate the visibility polygon
	sorted := sortLines(viewPt, lines)
	mapper := &array{data: make([]int, len(lines))}
	for i := range mapper.data {
		mapper.data[i] = -1
	}
	heap := &array{}
	start := geom.Point{X: viewPt.X + 1, Y: viewPt.Y}
	for i := range lines {
		a1 := angle(lines[i].Start, viewPt)
		a2 := angle(lines[i].End, viewPt)
		if (a1 >= -180 && a1 <= 0 && a2 <= 180 && a2 >= 0 && a2-a1 > 180) ||
			(a2 >= -180 && a2 <= 0 && a1 <= 180 && a1 >= 0 && a1-a2 > 180) {
			v.insert(i, heap, mapper, lines, viewPt, start)
		}
	}
	polygon := make([]geom.Point, 0, len(sorted)*2)
	i := 0
	for i < len(sorted) {
		extend := false
		shorten := false
		orig := i
		vertex := sorted[i].pt(lines)
		oldLine := heap.top()
		for {
			if mapper.elem(sorted[i].lineIndex) != -1 {
				if sorted[i].lineIndex == oldLine {
					extend = true
					vertex = sorted[i].pt(lines)
				}
				v.remove(mapper.elem(sorted[i].lineIndex), heap, mapper, lines, viewPt, vertex)
			} else {
				v.insert(sorted[i].lineIndex, heap, mapper, lines, viewPt, vertex)
				if heap.top() != oldLine {
					shorten = true
				}
			}
			i++
			if i == len(sorted) || sorted[i].angle >= sorted[orig].angle+angleEpsilon {
				break
			}
		}
		// The heap can be emptied by the removals above, so every nearest-occluder lookup has to tolerate the empty
		// case rather than indexing blindly.
		if extend {
			polygon = append(polygon, geom.Point{X: vertex.X, Y: vertex.Y})
			if nearest := heap.top(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects &&
					!cur.EqualWithin(vertex, v.epsilon) {
					polygon = append(polygon, geom.Point{X: cur.X, Y: cur.Y})
				}
			}
		} else if shorten {
			if oldLine != -1 {
				line := lines[oldLine]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = append(polygon, geom.Point{X: cur.X, Y: cur.Y})
				}
			}
			if nearest := heap.top(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = append(polygon, geom.Point{X: cur.X, Y: cur.Y})
				}
			}
		}
	}
	if len(polygon) == 0 {
		return nil
	}
	return polygon
}

func (v *Visibility) remove(index int, heap, mapper *array, lines []geom.Line, position, destination geom.Point) {
	mapper.set(heap.elem(index), -1)
	if index == heap.size()-1 {
		heap.pop()
		return
	}
	heap.set(index, heap.pop())
	mapper.set(heap.elem(index), index)
	cur := index
	parent := (cur - 1) / 2
	if cur != 0 && v.lessThan(heap.elem(cur), heap.elem(parent), lines, position, destination) {
		for cur > 0 {
			parent = (cur - 1) / 2
			parentElem := heap.elem(parent)
			curElem := heap.elem(cur)
			if !v.lessThan(curElem, parentElem, lines, position, destination) {
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
		case left < heap.size() && v.lessThan(heap.elem(left), heap.elem(cur), lines, position, destination) &&
			(right == heap.size() || v.lessThan(heap.elem(left), heap.elem(right), lines, position, destination)):
			leftElem := heap.elem(left)
			curElem := heap.elem(cur)
			mapper.set(leftElem, cur)
			mapper.set(curElem, left)
			heap.set(left, curElem)
			heap.set(cur, leftElem)
			cur = left
		case right < heap.size() && v.lessThan(heap.elem(right), heap.elem(cur), lines, position, destination):
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

func (v *Visibility) insert(index int, heap, mapper *array, lines []geom.Line, position, destination geom.Point) {
	cur := heap.size()
	heap.push(index)
	mapper.set(index, cur)
	for cur > 0 {
		parent := (cur - 1) / 2
		parentElem := heap.elem(parent)
		curElem := heap.elem(cur)
		if !v.lessThan(curElem, parentElem, lines, position, destination) {
			break
		}
		mapper.set(parentElem, cur)
		mapper.set(curElem, parent)
		heap.set(cur, parentElem)
		heap.set(parent, curElem)
		cur = parent
	}
}

// lessThan reports whether the line at index1 is closer to position along the ray towards destination than the line at
// index2. A line that the ray misses entirely sorts as infinitely far away, so that it can always be displaced from the
// front of the heap by one that is actually hit. Returning false for both orderings instead would break the strict weak
// ordering the sift-up and sift-down code depends on and would let a missed line sit at the front of the heap forever.
func (v *Visibility) lessThan(index1, index2 int, lines []geom.Line, position, destination geom.Point) bool {
	pt1, intersects1 := intersectLines(lines[index1].Start, lines[index1].End, position, destination)
	pt2, intersects2 := intersectLines(lines[index2].Start, lines[index2].End, position, destination)
	if !intersects1 {
		return false
	}
	if !intersects2 {
		return true
	}
	if !pt1.EqualWithin(pt2, v.epsilon) {
		d1 := distSqrd(pt1, position)
		d2 := distSqrd(pt2, position)
		return d1 < d2
	}
	var a1 float32
	if pt1.EqualWithin(lines[index1].Start, v.epsilon) {
		a1 = angle2(lines[index1].End, pt1, position)
	} else {
		a1 = angle2(lines[index1].Start, pt1, position)
	}
	var a2 float32
	if pt2.EqualWithin(lines[index2].Start, v.epsilon) {
		a2 = angle2(lines[index2].End, pt2, position)
	} else {
		a2 = angle2(lines[index2].Start, pt2, position)
	}
	if a1 < 180 {
		if a2 > 180 {
			return true
		}
		return a2 < a1
	}
	return a1 < a2
}

func sortLines(position geom.Point, lines []geom.Line) []endPoint {
	points := make([]endPoint, len(lines)*2)
	pos := 0
	for i, line := range lines {
		points[pos].lineIndex = i
		points[pos].angle = angle(line.Start, position)
		points[pos].start = true
		pos++
		points[pos].lineIndex = i
		points[pos].angle = angle(line.End, position)
		points[pos].start = false
		pos++
	}
	slices.SortFunc(points, func(a, b endPoint) int {
		if result := cmp.Compare(a.angle, b.angle); result != 0 {
			return result
		}
		if result := cmp.Compare(distSqrd(a.pt(lines), position), distSqrd(b.pt(lines), position)); result != 0 {
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

func angle2(a, b, c geom.Point) float32 {
	a1 := angle(a, b)
	a2 := angle(b, c)
	a3 := a1 - a2
	if a3 < 0 {
		a3 += 360
	}
	if a3 > 360 {
		a3 -= 360
	}
	return a3
}

func angle(a, b geom.Point) float32 {
	return xmath.Atan2(b.Y-a.Y, b.X-a.X) * 180 / math.Pi
}

// pointEpsilon returns the tolerance to use when comparing points for equality in a scene whose coordinates reach the
// given magnitude. A tolerance proportional to the scene keeps behavior the same at every scale, while the floor keeps
// it from collapsing to zero for a scene with no measurable extent.
func pointEpsilon(magnitude float32) float32 {
	if eps := magnitude * pointEpsilonRatio; eps > minPointEpsilon {
		return eps
	}
	return minPointEpsilon
}

func distSqrd(a, b geom.Point) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func intersectLines(s1, e1, s2, e2 geom.Point) (geom.Point, bool) {
	dbx := e2.X - s2.X
	dby := e2.Y - s2.Y
	dax := e1.X - s1.X
	day := e1.Y - s1.Y
	ub := dby*dax - dbx*day
	if ub == 0 {
		return geom.Point{}, false
	}
	ua := (dbx*(s1.Y-s2.Y) - dby*(s1.X-s2.X)) / ub
	return geom.Point{X: s1.X + ua*dax, Y: s1.Y + ua*day}, true
}

func hasIntersection(s1, e1, s2, e2 geom.Point) bool {
	d1 := direction(s2, e2, s1)
	d2 := direction(s2, e2, e1)
	d3 := direction(s1, e1, s2)
	d4 := direction(s1, e1, e2)
	return (((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0))) ||
		(d1 == 0 && onLine(s2, e2, s1)) ||
		(d2 == 0 && onLine(s2, e2, e1)) ||
		(d3 == 0 && onLine(s1, e1, s2)) ||
		(d4 == 0 && onLine(s1, e1, e2))
}

func direction(a, b, c geom.Point) int {
	return cmp.Compare((c.X-a.X)*(b.Y-a.Y), (b.X-a.X)*(c.Y-a.Y))
}

func onLine(a, b, c geom.Point) bool {
	return (a.X <= c.X || b.X <= c.X) &&
		(c.X <= a.X || c.X <= b.X) &&
		(a.Y <= c.Y || b.Y <= c.Y) &&
		(c.Y <= a.Y || c.Y <= b.Y)
}

type endPoint struct {
	lineIndex int
	angle     float32
	start     bool
}

func (ep *endPoint) pt(lines []geom.Line) geom.Point {
	if ep.start {
		return lines[ep.lineIndex].Start
	}
	return lines[ep.lineIndex].End
}

type array struct {
	data []int
}

func (a *array) size() int {
	return len(a.data)
}

func (a *array) elem(index int) int {
	return a.data[index]
}

// top returns the element at the front of the heap, or -1 if the heap is empty.
func (a *array) top() int {
	if len(a.data) == 0 {
		return -1
	}
	return a.data[0]
}

func (a *array) set(index, value int) {
	a.data[index] = value
}

func (a *array) pop() int {
	v := a.data[len(a.data)-1]
	a.data = a.data[:len(a.data)-1]
	return v
}

func (a *array) push(v int) {
	a.data = append(a.data, v)
}
