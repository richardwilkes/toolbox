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
	// parallelEpsilonSqrd is the square of the largest sine of the angle between two lines that still counts as
	// parallel. Squared because intersectLines can obtain the squared sine without a square root.
	parallelEpsilonSqrd = 1e-12
)

// Visibility holds state for computing a visibility polygon. A Visibility is immutable once created, so it is safe for
// concurrent use by multiple goroutines.
type Visibility struct {
	lines   []geom.Line
	bounds  geom.Rect
	epsilon float32
}

// New creates a Visibility object.
//
// bounds should not be empty. Nothing rejects one that is, but no point lies inside an empty rectangle, so every call
// to PolygonFrom on such a Visibility returns nil.
//
// The obstructions must not intersect each other, which is not verified. If they do, call BreakIntersections() and
// pass the result instead.
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
	var intersections []geom.Point
	for _, line := range lines {
		intersections = intersections[:0]
		for _, one := range qt.FindIntersects(line.Bounds()) {
			// geom.LineIntersection is used rather than the local intersectLines because it returns both endpoints of
			// the shared portion when the two segments are collinear and overlap, a case that cannot be expressed as a
			// single infinite-line intersection point. Every candidate is passed through it, including the line
			// itself: a line intersected with itself yields its own endpoints, which the filter below discards.
			// Skipping the self case by equality instead would also make each of two identical input lines skip its
			// twin, leaving them uncut.
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

// collectLines splits line at each of the given intersection points and appends the resulting segments to lines,
// returning the extended slice. Zero-length segments are dropped, as are segments that fall outside viewPort when one
// is supplied. The intersections slice is sorted in place.
func collectLines(line geom.Line, intersections []geom.Point, lines []geom.Line, viewPort *geom.Rect, eps float32,
) []geom.Line {
	// The intersection points all lie on line, so ordering them by distance from its start walks them from one end to
	// the other.
	slices.SortFunc(intersections, func(a, b geom.Point) int {
		return cmp.Compare(distSqrd(line.Start, a), distSqrd(line.Start, b))
	})
	start := line.Start
	for _, end := range intersections {
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

// PolygonFrom returns the polygon of the unobstructed area visible from viewPt, or nil if viewPt is outside the bounds
// or no area is visible. Consecutive vertices are always distinct, so the polygon has no zero-length edges.
func (v *Visibility) PolygonFrom(viewPt geom.Point) []geom.Point {
	if !viewPt.In(v.bounds) {
		return nil
	}

	// Generate a revised line list by clipping the lines against the viewport and throwing out any that aren't within
	// the viewport.
	lines := make([]geom.Line, 0, len(v.lines)*2+4)
	viewport := [4]geom.Point{
		v.bounds.Point,
		v.bounds.TopRight(),
		v.bounds.BottomRight(),
		v.bounds.BottomLeft(),
	}
	var intersections []geom.Point
	for _, line := range v.lines {
		if (line.Start.X < v.bounds.X && line.End.X < v.bounds.X) ||
			(line.Start.Y < v.bounds.Y && line.End.Y < v.bounds.Y) ||
			(line.Start.X > v.bounds.Right() && line.End.X > v.bounds.Right()) ||
			(line.Start.Y > v.bounds.Bottom() && line.End.Y > v.bounds.Bottom()) {
			continue
		}
		intersections = intersections[:0]
		for j := range viewport {
			k := (j + 1) % len(viewport)
			for _, pt := range geom.LineIntersection(line.Start, line.End, viewport[j], viewport[k]) {
				if !pt.EqualWithin(line.Start, v.epsilon) && !pt.EqualWithin(line.End, v.epsilon) {
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
	heap := newLineHeap(lines, viewPt, v.epsilon)
	start := geom.Point{X: viewPt.X + 1, Y: viewPt.Y}
	for i := range lines {
		a1 := angle(lines[i].Start, viewPt)
		a2 := angle(lines[i].End, viewPt)
		// Seed the heap with the lines that straddle the ±180° discontinuity the sweep starts from. angle returns an
		// atan2 result, so both angles already lie within [-180,180]; all that is left to test is that they fall on
		// opposite sides of zero and are more than 180° apart.
		if (a1 <= 0 && a2 >= 0 && a2-a1 > 180) || (a2 <= 0 && a1 >= 0 && a1-a2 > 180) {
			heap.insert(i, start)
		}
	}
	// The sweep emits at most two points per batch of endpoints that share an angle, but in practice lands a little
	// above one per endpoint, so size the polygon for that rather than for the worst case.
	polygon := make([]geom.Point, 0, len(sorted)+len(sorted)/2)
	i := 0
	for i < len(sorted) {
		extend := false
		shorten := false
		orig := i
		vertex := sorted[i].pt(lines)
		oldLine := heap.nearest()
		for {
			if heap.contains(sorted[i].lineIndex) {
				if sorted[i].lineIndex == oldLine {
					extend = true
					vertex = sorted[i].pt(lines)
				}
				heap.remove(sorted[i].lineIndex, vertex)
			} else {
				heap.insert(sorted[i].lineIndex, vertex)
				if heap.nearest() != oldLine {
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
			polygon = appendVertex(polygon, vertex, v.epsilon)
			if nearest := heap.nearest(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects &&
					!cur.EqualWithin(vertex, v.epsilon) {
					polygon = appendVertex(polygon, cur, v.epsilon)
				}
			}
		} else if shorten {
			if oldLine != -1 {
				line := lines[oldLine]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = appendVertex(polygon, cur, v.epsilon)
				}
			}
			if nearest := heap.nearest(); nearest != -1 {
				line := lines[nearest]
				if cur, intersects := intersectLines(line.Start, line.End, viewPt, vertex); intersects {
					polygon = appendVertex(polygon, cur, v.epsilon)
				}
			}
		}
	}
	// The sweep wraps all the way around, so it can finish on the vertex it started from. Drop that repeat as well,
	// since the closing edge is implicit.
	if len(polygon) > 1 && polygon[0].EqualWithin(polygon[len(polygon)-1], v.epsilon) {
		polygon = polygon[:len(polygon)-1]
	}
	if len(polygon) == 0 {
		return nil
	}
	return polygon
}

// appendVertex appends pt to the polygon unless it repeats the vertex already at the end. Emitting the repeat would
// leave a zero-length edge for every consumer to filter out, and would break polygon algorithms that assume each edge
// has a direction.
func appendVertex(polygon []geom.Point, pt geom.Point, eps float32) []geom.Point {
	if len(polygon) != 0 && polygon[len(polygon)-1].EqualWithin(pt, eps) {
		return polygon
	}
	return append(polygon, pt)
}

// lineHeap is a min-heap of indices into lines, ordered by how close each line is to viewPt along a ray towards a
// destination supplied per operation. order holds line indices keyed by heap position, and slot holds heap positions
// keyed by line index; the two are inverses of each other, which is what lets a line be located in the heap without
// searching.
type lineHeap struct {
	lines   []geom.Line
	order   []int
	slot    []int
	viewPt  geom.Point
	epsilon float32
}

// newLineHeap returns an empty heap over the given lines. Both index slices come from one allocation, since neither
// can ever hold more than one entry per line.
func newLineHeap(lines []geom.Line, viewPt geom.Point, epsilon float32) lineHeap {
	backing := make([]int, len(lines)*2)
	h := lineHeap{
		lines:   lines,
		order:   backing[:0:len(lines)],
		slot:    backing[len(lines):],
		viewPt:  viewPt,
		epsilon: epsilon,
	}
	for i := range h.slot {
		h.slot[i] = -1
	}
	return h
}

// contains reports whether the given line is currently in the heap.
func (h *lineHeap) contains(lineIndex int) bool {
	return h.slot[lineIndex] != -1
}

// nearest returns the index of the line closest to the view point, or -1 if the heap is empty.
func (h *lineHeap) nearest() int {
	if len(h.order) == 0 {
		return -1
	}
	return h.order[0]
}

// insert adds the given line to the heap, ordering it against the ray from the view point towards destination.
func (h *lineHeap) insert(lineIndex int, destination geom.Point) {
	cur := len(h.order)
	h.order = append(h.order, lineIndex)
	h.slot[lineIndex] = cur
	h.siftUp(cur, destination)
}

// remove takes the given line back out of the heap, reordering what is left against the ray from the view point
// towards destination.
func (h *lineHeap) remove(lineIndex int, destination geom.Point) {
	cur := h.slot[lineIndex]
	h.slot[lineIndex] = -1
	last := len(h.order) - 1
	if cur != last {
		h.order[cur] = h.order[last]
		h.slot[h.order[cur]] = cur
	}
	h.order = h.order[:last]
	if cur != last && !h.siftUp(cur, destination) {
		h.siftDown(cur, destination)
	}
}

// swap exchanges the entries at two heap positions, keeping slot in step with order.
func (h *lineHeap) swap(a, b int) {
	h.slot[h.order[a]], h.slot[h.order[b]] = b, a
	h.order[a], h.order[b] = h.order[b], h.order[a]
}

// siftUp moves the entry at the given heap position towards the front while it is closer than its parent, reporting
// whether it moved at all.
func (h *lineHeap) siftUp(cur int, destination geom.Point) bool {
	moved := false
	for cur > 0 {
		parent := (cur - 1) / 2
		if !h.lessThan(h.order[cur], h.order[parent], destination) {
			break
		}
		h.swap(cur, parent)
		cur = parent
		moved = true
	}
	return moved
}

// siftDown moves the entry at the given heap position towards the back while either child is closer than it is.
func (h *lineHeap) siftDown(cur int, destination geom.Point) {
	for {
		left := 2*cur + 1
		if left >= len(h.order) {
			return
		}
		closest := left
		if right := left + 1; right < len(h.order) && h.lessThan(h.order[right], h.order[left], destination) {
			closest = right
		}
		if !h.lessThan(h.order[closest], h.order[cur], destination) {
			return
		}
		h.swap(cur, closest)
		cur = closest
	}
}

// lessThan reports whether the line at index1 is closer to the view point along the ray towards destination than the
// line at index2. A line that the ray misses entirely sorts as infinitely far away, so that it can always be displaced
// from the front of the heap by one that is actually hit. Returning false for both orderings instead would break the
// strict weak ordering the sift-up and sift-down code depends on and would let a missed line sit at the front of the
// heap forever.
func (h *lineHeap) lessThan(index1, index2 int, destination geom.Point) bool {
	pt1, intersects1 := intersectLines(h.lines[index1].Start, h.lines[index1].End, h.viewPt, destination)
	pt2, intersects2 := intersectLines(h.lines[index2].Start, h.lines[index2].End, h.viewPt, destination)
	if !intersects1 {
		return false
	}
	if !intersects2 {
		return true
	}
	if !pt1.EqualWithin(pt2, h.epsilon) {
		return distSqrd(pt1, h.viewPt) < distSqrd(pt2, h.viewPt)
	}
	var a1 float32
	if pt1.EqualWithin(h.lines[index1].Start, h.epsilon) {
		a1 = angle2(h.lines[index1].End, pt1, h.viewPt)
	} else {
		a1 = angle2(h.lines[index1].Start, pt1, h.viewPt)
	}
	var a2 float32
	if pt2.EqualWithin(h.lines[index2].Start, h.epsilon) {
		a2 = angle2(h.lines[index2].End, pt2, h.viewPt)
	} else {
		a2 = angle2(h.lines[index2].Start, pt2, h.viewPt)
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
	// Both angles are atan2 results in [-180,180], so their difference is in [-360,360] and one adjustment is enough
	// to bring it into [0,360].
	a3 := angle(a, b) - angle(b, c)
	if a3 < 0 {
		a3 += 360
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

// intersectLines returns the point where the infinite lines through the two segments cross. It reports false when the
// lines are parallel, or so nearly parallel that the crossing point would be far enough away to swamp every distance
// it was compared against. ub is the cross product of the two direction vectors, so comparing its square against the
// product of their squared lengths tests the squared sine of the angle between them, which is independent of scale.
func intersectLines(s1, e1, s2, e2 geom.Point) (geom.Point, bool) {
	dbx := e2.X - s2.X
	dby := e2.Y - s2.Y
	dax := e1.X - s1.X
	day := e1.Y - s1.Y
	ub := dby*dax - dbx*day
	if ub*ub <= (parallelEpsilonSqrd*(dax*dax+day*day))*(dbx*dbx+dby*dby) {
		return geom.Point{}, false
	}
	ua := (dbx*(s1.Y-s2.Y) - dby*(s1.X-s2.X)) / ub
	return geom.Point{X: s1.X + ua*dax, Y: s1.Y + ua*day}, true
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
