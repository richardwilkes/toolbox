// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

// Subtract returns the result of removing the other polygon from this polygon. Does not handle self-intersecting
// polygons. Use Simplify() first, if necessary.
func (p Polygon[T]) Subtract(other Polygon[T]) Polygon[T] {
	if len(p) == 0 || len(other) == 0 {
		return p.Clone()
	}
	sb := p.Bounds()
	cb := other.Bounds()
	if !sb.Intersects(cb) {
		return p.Clone()
	}
	var q eventQueue[T]
	for _, one := range p {
		for i := range one {
			q.addProcessedSegment(one.segment(i), true)
		}
	}
	for _, one := range other {
		for i := range one {
			q.addProcessedSegment(one.segment(i), false)
		}
	}
	var c connector[T]
	var s sweepline[T]
	for !q.empty() {
		var prev, next *endpoint[T]
		e := q.dequeue()
		if e.pt.X > sb.Right() {
			return c.toPolygon()
		}
		if e.left {
			pos := s.insert(e)
			if pos > 0 {
				prev = s[pos-1]
			} else {
				prev = nil
			}
			if pos < len(s)-1 {
				next = s[pos+1]
			} else {
				next = nil
			}
			switch {
			case prev == nil:
				e.inside = false
				e.inout = false
			case prev.edgeType != edgeNormal:
				if pos-2 < 0 {
					e.inside = false
					e.inout = false
					if prev.subject != e.subject {
						e.inside = true
					} else {
						e.inout = true
					}
				} else if e.segment() == prev.segment() {
					if e.edgeType == edgeSameTransition || prev.edgeType == edgeSameTransition {
						e.inout = prev.inout
					} else {
						e.inout = !prev.inout
					}
				} else {
					prevTwo := s[pos-2]
					if prev.subject == e.subject {
						e.inout = !prev.inout
						e.inside = !prevTwo.inout
					} else {
						e.inout = !prevTwo.inout
						e.inside = !prev.inout
					}
				}
			case e.subject == prev.subject:
				e.inside = prev.inside
				e.inout = !prev.inout
			default:
				e.inside = !prev.inout
				e.inout = prev.inside
			}
			divided := make(map[*endpoint[T]]bool)
			if next != nil {
				for _, seg := range q.possibleIntersection(e, next) {
					if seg != nil {
						divided[seg] = true
					}
				}
			}
			if prev != nil {
				for _, seg := range q.possibleIntersection(prev, e) {
					if seg != nil {
						divided[seg] = true
					}
				}
			}
			if len(divided) > 0 && !divided[e] {
				s.remove(e)
				q.enqueue(e)
			}
		} else {
			otherPos := -1
			for i := range s {
				if s[i] == e.other {
					otherPos = i
					break
				}
			}
			if otherPos != -1 {
				if otherPos > 0 {
					prev = s[otherPos-1]
				} else {
					prev = nil
				}
				if otherPos < len(s)-1 {
					next = s[otherPos+1]
				} else {
					next = nil
				}
			}
			if e.edgeType == edgeDifferentTransition || e.subject == !e.other.inside {
				c.add(e.segment())
			}
			if otherPos != -1 {
				s.remove(s[otherPos])
			}
			if next != nil && prev != nil {
				q.possibleIntersection(next, prev)
			}
		}
	}
	return c.toPolygon()
}