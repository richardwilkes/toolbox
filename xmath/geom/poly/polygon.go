package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

type op int

const (
	unionOp op = iota
	intersectOp
	subtractOp
	xorOp
	clipLineOp
)

// Polygon holds one or more contour lines. The polygon may contain holes and
// may be self-intersecting. Note that the methods Union, Intersect, Subtract,
// Xor, and ClipLine all require that the polygons not have have
// self-intersections. Use of the Simplify() method will ensure a polygon can
// be used with those methods.
type Polygon []Contour

// ApproximateEllipse creates a polygon that approximates an ellipse.
// 'sections' indicates how many segments to break the ellipse contour into.
func ApproximateEllipse(bounds geom.Rect, sections int) Polygon {
	halfWidth := bounds.Width / 2
	halfHeight := bounds.Height / 2
	inc := math.Pi * 2 / float64(sections)
	center := bounds.Center()
	contour := make(Contour, sections)
	var angle float64
	for i := 0; i < sections; i++ {
		contour[i] = geom.Point{
			X: center.X + math.Cos(angle)*halfWidth,
			Y: center.Y + math.Sin(angle)*halfHeight,
		}
		angle += inc
	}
	return Polygon{contour}
}

// Rect creates a new polygon in the shape of a rectangle.
func Rect(bounds geom.Rect) Polygon {
	return Polygon{Contour{
		bounds.Point,
		geom.Point{X: bounds.X, Y: bounds.Bottom()},
		geom.Point{X: bounds.Right(), Y: bounds.Bottom()},
		geom.Point{X: bounds.Right(), Y: bounds.Y},
	}}
}

// Add a contour to a polygon.
func (p *Polygon) Add(c Contour) {
	*p = append(*p, c)
}

// Clone returns a duplicate of this polygon.
func (p Polygon) Clone() Polygon {
	clone := Polygon(make([]Contour, len(p)))
	for i := range p {
		clone[i] = p[i].Clone()
	}
	return clone
}

// Bounds returns the bounding rectangle of this polygon.
func (p Polygon) Bounds() geom.Rect {
	if len(p) == 0 {
		return geom.Rect{}
	}
	bb := p[0].Bounds()
	for _, c := range p[1:] {
		bb.Union(c.Bounds())
	}
	return bb
}

// Contains returns true if the point is contained by the polygon.
func (p Polygon) Contains(pt geom.Point) bool {
	for i := range p {
		if p[i].Contains(pt) {
			return true
		}
	}
	return false
}

// ContainsEvenOdd returns true if the point is contained by the polygon using
// the even-odd rule. https://en.wikipedia.org/wiki/Even-odd_rule
func (p Polygon) ContainsEvenOdd(pt geom.Point) bool {
	var count int
	for i := range p {
		if p[i].Contains(pt) {
			count++
		}
	}
	return count%2 == 1
}

// Union returns the union of both polygons.
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) Union(other Polygon) Polygon {
	return p.construct(unionOp, other)
}

// Intersect returns the intersection of both polygons.
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) Intersect(other Polygon) Polygon {
	return p.construct(intersectOp, other)
}

// Subtract returns the result of removing the other polygon from this
// polygon.
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) Subtract(other Polygon) Polygon {
	return p.construct(subtractOp, other)
}

// Xor returns the result of xor'ing this polygon with the other polygon.
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) Xor(other Polygon) Polygon {
	return p.construct(xorOp, other)
}

// ClipLine returns the result of removing the other polygon from this
// polygon. Assumes this polygon is actually a line.
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) ClipLine(other Polygon) Polygon {
	return p.construct(clipLineOp, other)
}

// construct uses the algorithm described at
// http://www.cs.ucr.edu/~vbz/cs230papers/martinez_boolean.pdf
//
// Note: This function is not designed to handle self-intersecting polygons.
// Remove self-intersections first using the Simplify function.
func (p Polygon) construct(operation op, other Polygon) Polygon {
	// If one is empty, short-circuit the result
	if len(p) == 0 || len(other) == 0 {
		switch operation {
		case subtractOp:
			return p.Clone()
		case unionOp:
			if len(p) == 0 {
				return other.Clone()
			}
			return p.Clone()
		}
		return Polygon{}
	}

	// If they do not intersect, short-circuit the result
	bounds := p.Bounds()
	otherBounds := other.Bounds()
	if !bounds.Intersects(otherBounds) {
		switch operation {
		case subtractOp:
			return p.Clone()
		case unionOp:
			result := p.Clone()
			for _, cont := range other {
				result.Add(cont.Clone())
			}
			return result
		}
		return Polygon{}
	}

	// Add each segment to the event queue, sorted from left to right
	q := &queue{}
	for _, cont := range p {
		for i := range cont {
			if operation != clipLineOp || i != len(cont)-1 {
				addSegmentToQueue(cont.segment(i), true, q)
			}
		}
	}
	for _, cont := range other {
		for i := range cont {
			addSegmentToQueue(cont.segment(i), false, q)
		}
	}

	// Process all events
	s := sweep{}
	minOfMaxX := math.Min(bounds.Right(), otherBounds.Right())
	conn := connector{op: operation}
	for q.more() {
		e := q.dequeue()
		if ((operation == intersectOp || operation == clipLineOp) && e.pt.X > minOfMaxX) || (operation == subtractOp && e.pt.X > bounds.Right()) {
			return conn.toPolygon()
		}

		var prev, next *edge
		if e.left {
			// Line segment must be inserted
			pos := s.insert(e)
			if pos > 0 {
				prev = s[pos-1]
			}
			if pos < len(s)-1 {
				next = s[pos+1]
			}

			// Compute the inside and inOut flags
			switch {
			case prev == nil:
				e.inside = false
				e.inOut = false
			case prev.edgeType != normalEdge:
				if pos-2 < 0 {
					if prev.subject != e.subject {
						e.inside = true
						e.inOut = false
					} else {
						e.inside = false
						e.inOut = true
					}
				} else {
					prevTwo := s[pos-2]
					if prev.subject == e.subject {
						e.inOut = !prev.inOut
						e.inside = !prevTwo.inOut
					} else {
						e.inOut = !prevTwo.inOut
						e.inside = !prev.inOut
					}
				}
			case e.subject == prev.subject:
				e.inside = prev.inside
				e.inOut = !prev.inOut
			default:
				e.inside = !prev.inOut
				e.inOut = prev.inside
			}

			// Process a possible intersections
			if next != nil {
				checkForIntersection(e, next, q)
			}
			if prev != nil {
				if divided := checkForIntersection(prev, e, q); len(divided) == 1 && divided[0] == prev {
					s.remove(e)
					q.enqueue(e)
				}
			}
		} else {
			// Line segment must be removed
			otherPos := -1
			for i := range s {
				if s[i].equals(e.other) {
					otherPos = i
					if otherPos > 0 {
						prev = s[otherPos-1]
					}
					if otherPos < len(s)-1 {
						next = s[otherPos+1]
					}
					break
				}
			}

			if operation == clipLineOp && e.other.inside && e.subject {
				conn.add(e.segment())
			}
			switch e.edgeType {
			case normalEdge:
				switch operation {
				case intersectOp:
					if e.other.inside {
						conn.add(e.segment())
					}
				case unionOp:
					if !e.other.inside {
						conn.add(e.segment())
					}
				case subtractOp:
					if (e.subject && !e.other.inside) || (!e.subject && e.other.inside) {
						conn.add(e.segment())
					}
				case xorOp:
					conn.add(e.segment())
				}
			case sameTransitionEdge:
				if operation == intersectOp || operation == unionOp {
					conn.add(e.segment())
				}
			case differentTransitionEdge:
				if operation == subtractOp {
					conn.add(e.segment())
				}
			}

			if otherPos != -1 {
				s.remove(s[otherPos])
			}
			if next != nil && prev != nil {
				checkForIntersection(next, prev, q)
			}
		}
	}
	return conn.toPolygon()
}

func checkForIntersection(e1, e2 *edge, q *queue) []*edge {
	numIntersections, ip1, _ := findIntersection(e1.segment(), e2.segment(), true)
	switch numIntersections {
	case 0:
		return nil
	case 1:
		ip1 = snap(ip1, e1.pt, e2.pt, e1.other.pt, e2.other.pt)
		switch {
		case e1.pt == e2.pt || e1.other.pt == e2.other.pt || !validIntersection(e1, e2, ip1):
			return nil
		case e1.pt == ip1 || e1.other.pt == ip1:
			return []*edge{divideSegment(e2, ip1, q)}
		case e2.pt == ip1 || e2.other.pt == ip1:
			return []*edge{divideSegment(e1, ip1, q)}
		default:
			return []*edge{
				divideSegment(e1, ip1, q),
				divideSegment(e2, ip1, q),
			}
		}
	default:
		if e1.subject == e2.subject {
			return nil
		}
		sortedEvents := addSortedEvents(e1.other, e2.other, addSortedEvents(e1, e2, make([]*edge, 0, 4)))
		switch {
		case len(sortedEvents) == 2: // line segments are equal
			e1.edgeType = nonContributingEdge
			e1.other.edgeType = nonContributingEdge
			if e1.inOut == e2.inOut {
				e2.edgeType = sameTransitionEdge
				e2.other.edgeType = sameTransitionEdge
			} else {
				e2.edgeType = differentTransitionEdge
				e2.other.edgeType = differentTransitionEdge
			}
			return nil
		case len(sortedEvents) == 3: // line segments share an edge
			sortedEvents[1].edgeType = nonContributingEdge
			sortedEvents[1].other.edgeType = nonContributingEdge
			var which int
			if sortedEvents[0] != nil {
				which = 0
			} else {
				which = 2
			}
			if e1.inOut == e2.inOut {
				sortedEvents[which].other.edgeType = sameTransitionEdge
			} else {
				sortedEvents[which].other.edgeType = differentTransitionEdge
			}
			if which == 0 {
				return []*edge{divideSegment(sortedEvents[0], sortedEvents[1].pt, q)}
			}
			return []*edge{divideSegment(sortedEvents[2].other, sortedEvents[1].pt, q)}
		case sortedEvents[0] != sortedEvents[3].other: // no line segment completely includes the other one
			sortedEvents[1].edgeType = nonContributingEdge
			if e1.inOut == e2.inOut {
				sortedEvents[2].edgeType = sameTransitionEdge
			} else {
				sortedEvents[2].edgeType = differentTransitionEdge
			}
			return []*edge{
				divideSegment(sortedEvents[0], sortedEvents[1].pt, q),
				divideSegment(sortedEvents[1], sortedEvents[2].pt, q),
			}
		default: // one line segment includes the other one
			sortedEvents[1].edgeType = nonContributingEdge
			sortedEvents[1].other.edgeType = nonContributingEdge
			divideSegment(sortedEvents[0], sortedEvents[1].pt, q)
			if e1.inOut == e2.inOut {
				sortedEvents[3].other.edgeType = sameTransitionEdge
			} else {
				sortedEvents[3].other.edgeType = differentTransitionEdge
			}
			return []*edge{divideSegment(sortedEvents[3].other, sortedEvents[2].pt, q)}
		}
	}
}

func addSortedEvents(e1, e2 *edge, sortedEvents []*edge) []*edge {
	switch {
	case e1.pt == e2.pt:
		return append(sortedEvents, nil)
	case e1.less(e2):
		return append(sortedEvents, e2, e1)
	default:
		return append(sortedEvents, e1, e2)
	}
}

func validIntersection(e1, e2 *edge, pt geom.Point) bool {
	switch {
	case e1.pt.X == pt.X && e2.pt.X == pt.X:
		return (pt.Y-e1.pt.Y > 0) != (pt.Y-e2.pt.Y > 0)
	case e1.pt.Y == pt.Y && e2.pt.Y == pt.Y:
		return (pt.X-e1.pt.X > 0) != (pt.X-e2.pt.X > 0)
	case e1.other.pt.X == pt.X && e2.other.pt.X == pt.X:
		return (pt.Y-e1.other.pt.Y > 0) != (pt.Y-e2.other.pt.Y > 0)
	case e1.other.pt.Y == pt.Y && e2.other.pt.Y == pt.Y:
		return (pt.X-e1.other.pt.X > 0) != (pt.X-e2.other.pt.X > 0)
	}
	return true
}

func divideSegment(e *edge, pt geom.Point, q *queue) *edge {
	left := &edge{
		pt:       pt,
		left:     true,
		subject:  e.subject,
		other:    e.other,
		edgeType: e.other.edgeType,
	}
	if !left.isValidDirection() {
		return nil
	}
	right := &edge{
		pt:       pt,
		subject:  e.subject,
		other:    e,
		edgeType: e.edgeType,
	}
	if !right.isValidDirection() {
		return nil
	}
	if left.less(e.other) {
		e.other.left = true
		e.left = false
	}
	e.other.other = left
	e.other = right
	q.enqueue(left)
	q.enqueue(right)
	return e
}

func snap(pt geom.Point, toPts ...geom.Point) geom.Point {
	for _, p := range toPts {
		if equalWithin(pt.X, p.X) && equalWithin(pt.Y, p.Y) {
			return p
		}
	}
	return pt
}

func equalWithin(a, b float64) bool {
	if a == b {
		return true
	}
	delta := math.Abs(a - b)
	const tolerance = 3e-14
	if delta <= tolerance {
		return true
	}
	return delta/math.Max(math.Abs(a), math.Abs(b)) <= tolerance
}

func findIntersection(seg0, seg1 segment, tryBothDirections bool) (numIntersections int, intersectionPt1, intersectionPt2 geom.Point) {
	const epsilonSquared = 1e-15
	d0 := geom.Point{X: seg0.end.X - seg0.start.X, Y: seg0.end.Y - seg0.start.Y}
	d1 := geom.Point{X: seg1.end.X - seg1.start.X, Y: seg1.end.Y - seg1.start.Y}
	d2 := geom.Point{X: seg1.start.X - seg0.start.X, Y: seg1.start.Y - seg0.start.Y}
	d3 := d0.X*d1.Y - d0.Y*d1.X
	d3Squared := d3 * d3
	d0Dist := distanceToZero(d0)
	d1Dist := distanceToZero(d1)
	if d3Squared > epsilonSquared*d0Dist*d1Dist {
		if s := (d2.X*d1.Y - d2.Y*d1.X) / d3; s >= 0 && s <= 1 {
			if t := (d2.X*d0.Y - d2.Y*d0.X) / d3; t >= 0 && t <= 1 {
				intersectionPt1.X = seg0.start.X + s*d0.X
				intersectionPt1.Y = seg0.start.Y + s*d0.Y
				numIntersections = 1
			}
		}
		return
	}
	d2Dist := distanceToZero(d2)
	d3 = d2.X*d0.Y - d2.Y*d0.X
	d3Squared = d3 * d3
	if d3Squared > epsilonSquared*d0Dist*d2Dist {
		return
	}
	s0 := (d0.X*d2.X + d0.Y*d2.Y) / d0Dist
	s1 := s0 + (d0.X*d1.X+d0.Y*d1.Y)/d0Dist
	m1 := math.Min(s0, s1)
	m2 := math.Max(s0, s1)
	var w0, w1 float64
	switch {
	case m1 > 1 || m2 < 0:
	case m1 == 1:
		w0 = 1
		numIntersections = 1
	case m2 == 0:
		numIntersections = 1
	default:
		if m1 > 0 {
			w0 = m1
		}
		if m2 < 1 {
			w1 = m2
		} else {
			w1 = 1
		}
		numIntersections = 2
	}
	if numIntersections > 0 {
		intersectionPt1.X = seg0.start.X + w0*d0.X
		intersectionPt1.Y = seg0.start.Y + w0*d0.Y
	}
	if numIntersections > 1 {
		intersectionPt2.X = seg0.start.X + w1*d0.X
		intersectionPt2.Y = seg0.start.Y + w1*d0.Y
	} else if tryBothDirections {
		if num, pt1, pt2 := findIntersection(seg1, seg0, false); num > numIntersections {
			return num, pt1, pt2
		}
	}
	return
}

func distanceToZero(pt geom.Point) float64 {
	return math.Sqrt(pt.X*pt.X + pt.Y*pt.Y)
}

func addSegmentToQueue(segment segment, isSubject bool, q *queue) {
	if segment.start == segment.end {
		return
	}
	e1 := &edge{
		pt:      segment.start,
		left:    true,
		subject: isSubject,
	}
	e2 := &edge{
		pt:      segment.end,
		left:    true,
		subject: isSubject,
		other:   e1,
	}
	e1.other = e2
	switch {
	case e1.pt.X < e2.pt.X:
		e2.left = false
	case e1.pt.X > e2.pt.X:
		e1.left = false
	case e1.pt.Y < e2.pt.Y:
		e2.left = false
	default:
		e1.left = false
	}
	q.enqueue(e1)
	q.enqueue(e2)
}

// Simplify returns a polygon with all self-intersections and repeated edges
// removed.
func (p Polygon) Simplify() Polygon {
	q := &queue{}
	var edgeCount int
	for _, cont := range p {
		for i := range cont {
			addSegmentToQueue(cont.segment(i), true, q)
			edgeCount++
		}
	}
	s := sweep{}
	edges := make([]*edge, 0, edgeCount)
	for q.more() {
		var prev, next *edge
		e := q.dequeue()
		if e.left {
			// Line segment must be inserted
			pos := s.insert(e)
			if pos > 0 {
				prev = s[pos-1]
			}
			if pos < len(s)-1 {
				next = s[pos+1]
			}
			if next != nil {
				simplify(e, next, q)
			}
			if prev != nil {
				if divided := simplify(prev, e, q); len(divided) == 1 && divided[0] == prev {
					s.remove(e)
					q.enqueue(e)
				}
			}
		} else {
			// Line segment must be removed
			otherPos := -1
			for i := range s {
				if s[i].equals(e.other) {
					otherPos = i
					break
				}
			}
			if otherPos != -1 {
				if otherPos > 0 {
					prev = s[otherPos-1]
				}
				if otherPos < len(s)-1 {
					next = s[otherPos+1]
				}
			}
			edges = append(edges, e)
			if otherPos != -1 {
				s.remove(s[otherPos])
			}
			if next != nil && prev != nil {
				simplify(next, prev, q)
			}
		}
	}
	conn := connector{op: unionOp}
	for i, e := range edges {
		if i == 0 || i == len(edges)-1 || (!(e.pt == edges[i+1].pt && e.other.pt == edges[i+1].other.pt) && !(e.pt == edges[i-1].pt && e.other.pt == edges[i-1].other.pt)) {
			conn.add(e.segment())
		}
	}
	return conn.toPolygon()
}

func simplify(e1, e2 *edge, q *queue) []*edge {
	numIntersections, pt1, pt2 := findIntersection(e1.segment(), e2.segment(), true)
	if numIntersections == 0 {
		return nil
	}
	e := make([]*edge, 0, 4)
	pt1 = snap(pt1, e1.pt, e2.pt, e1.other.pt, e2.other.pt)
	if numIntersections == 1 {
		if pt1 != e1.pt && pt1 != e1.other.pt {
			e = append(e, divideSegment(e1, pt1, q))
		}
		if pt1 != e2.pt && pt1 != e2.other.pt {
			e = append(e, divideSegment(e2, pt1, q))
		}
		return e
	}
	pt2 = snap(pt2, e1.pt, e2.pt, e1.other.pt, e2.other.pt)
	if pt1 != e1.pt && pt2 != e1.other.pt {
		e = append(e, divideSegment(e1, pt1, q))
	}
	if pt1 != e2.pt && pt2 != e2.other.pt {
		e = append(e, divideSegment(e2, pt1, q))
	}
	return e
}
