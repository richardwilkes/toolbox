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
	"slices"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

const snapTolerance = 8e-14

type eventQueue[T constraints.Float] struct {
	elements []*endpoint[T]
	sorted   bool
}

func (q *eventQueue[T]) enqueue(e *endpoint[T]) {
	if !q.sorted {
		q.elements = append(q.elements, e)
		return
	}
	if len(q.elements) == 0 {
		q.elements = append(q.elements, e)
		return
	}
	q.elements = append(q.elements, nil)
	i := len(q.elements) - 2
	for i >= 0 && endpointCmp(e, q.elements[i]) < 0 {
		q.elements[i+1] = q.elements[i]
		i--
	}
	q.elements[i+1] = e
}

func (q *eventQueue[T]) dequeue() *endpoint[T] {
	if !q.sorted {
		slices.SortFunc(q.elements, endpointCmp[T])
		q.sorted = true
	}
	e := q.elements[len(q.elements)-1]
	q.elements = q.elements[:len(q.elements)-1]
	return e
}

func (q *eventQueue[T]) addProcessedSegment(segment Segment[T], subject bool) {
	if segment.Start == segment.End {
		return
	}
	e1 := &endpoint[T]{
		pt:      segment.Start,
		subject: subject,
	}
	e2 := &endpoint[T]{
		pt:      segment.End,
		subject: subject,
		other:   e1,
	}
	e1.other = e2
	switch {
	case e1.pt.X < e2.pt.X:
		e1.left = true
	case e1.pt.X > e2.pt.X:
		e2.left = true
	case e1.pt.Y < e2.pt.Y:
		e1.left = true
	default:
		e2.left = true
	}
	q.enqueue(e1)
	q.enqueue(e2)
}

func (q *eventQueue[T]) divideSegment(e *endpoint[T], pt geom.Point[T]) *endpoint[T] {
	r := &endpoint[T]{
		pt:       pt,
		subject:  e.subject,
		edgeType: e.edgeType,
		other:    e,
	}
	l := &endpoint[T]{
		pt:       pt,
		subject:  e.subject,
		edgeType: e.other.edgeType,
		other:    e.other,
		left:     true,
	}
	if !l.isValidDirection() || !r.isValidDirection() {
		return nil
	}
	if endpointCmp(l, e.other) < 0 {
		e.other.left = true
		e.left = false
	}
	e.other.other = l
	e.other = r
	q.enqueue(l)
	q.enqueue(r)
	return e
}

func (q *eventQueue[T]) possibleIntersection(e1, e2 *endpoint[T]) []*endpoint[T] {
	numIntersections, ip1, _ := e1.segment().FindIntersection(e2.segment(), true)
	if numIntersections == 0 || (numIntersections == 1 && (e1.pt == e2.pt || e1.other.pt == e2.other.pt)) {
		return nil
	}
	ip1 = snap(ip1, e1, e2)
	if numIntersections == 1 {
		switch {
		case e1.pt == ip1 || e1.other.pt == ip1:
			return []*endpoint[T]{q.divideSegment(e2, ip1)}
		case e2.pt == ip1 || e2.other.pt == ip1:
			return []*endpoint[T]{q.divideSegment(e1, ip1)}
		case !isValidSingleIntersection(e1, e2, ip1):
			return nil
		default:
			return []*endpoint[T]{
				q.divideSegment(e1, ip1),
				q.divideSegment(e2, ip1),
			}
		}
	}
	if numIntersections == 2 && e1.subject == e2.subject {
		return nil
	}
	sortedEvents := make([]*endpoint[T], 0, 4)
	switch {
	case e1.pt == e2.pt:
		sortedEvents = append(sortedEvents, nil)
	case endpointCmp(e1, e2) < 0:
		sortedEvents = append(sortedEvents, e2, e1)
	default:
		sortedEvents = append(sortedEvents, e1, e2)
	}
	switch {
	case e1.other.pt == e2.other.pt:
		sortedEvents = append(sortedEvents, nil)
	case endpointCmp(e1.other, e2.other) < 0:
		sortedEvents = append(sortedEvents, e2.other, e1.other)
	default:
		sortedEvents = append(sortedEvents, e1.other, e2.other)
	}
	if len(sortedEvents) == 2 {
		e1.edgeType, e1.other.edgeType = edgeNonContributing, edgeNonContributing
		if e1.inout == e2.inout {
			e2.edgeType, e2.other.edgeType = edgeSameTransition, edgeSameTransition
		} else {
			e2.edgeType, e2.other.edgeType = edgeDifferentTransition, edgeDifferentTransition
		}
		return nil
	}
	if len(sortedEvents) == 3 {
		sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = edgeNonContributing, edgeNonContributing
		var i int
		if sortedEvents[0] != nil {
			i = 0
		} else {
			i = 2
		}
		if e1.inout == e2.inout {
			sortedEvents[i].other.edgeType = edgeSameTransition
		} else {
			sortedEvents[i].other.edgeType = edgeDifferentTransition
		}
		if sortedEvents[0] != nil {
			return []*endpoint[T]{q.divideSegment(sortedEvents[0], sortedEvents[1].pt)}
		}
		return []*endpoint[T]{q.divideSegment(sortedEvents[2].other, sortedEvents[1].pt)}
	}
	if sortedEvents[0] != sortedEvents[3].other {
		sortedEvents[1].edgeType = edgeNonContributing
		if e1.inout == e2.inout {
			sortedEvents[2].edgeType = edgeSameTransition
		} else {
			sortedEvents[2].edgeType = edgeDifferentTransition
		}
		return []*endpoint[T]{
			q.divideSegment(sortedEvents[0], sortedEvents[1].pt),
			q.divideSegment(sortedEvents[1], sortedEvents[2].pt),
		}
	}
	sortedEvents[1].edgeType, sortedEvents[1].other.edgeType = edgeNonContributing, edgeNonContributing
	firstDivided := q.divideSegment(sortedEvents[0], sortedEvents[1].pt)
	if e1.inout == e2.inout {
		sortedEvents[3].other.edgeType = edgeSameTransition
	} else {
		sortedEvents[3].other.edgeType = edgeDifferentTransition
	}
	return []*endpoint[T]{
		firstDivided,
		q.divideSegment(sortedEvents[3].other, sortedEvents[2].pt),
	}
}

func (q *eventQueue[T]) empty() bool {
	return len(q.elements) == 0
}

func isValidSingleIntersection[T constraints.Float](e1, e2 *endpoint[T], ip geom.Point[T]) bool {
	switch {
	case e1.pt.X == ip.X && e2.pt.X == ip.X:
		return (ip.Y-e1.pt.Y > 0) != (ip.Y-e2.pt.Y > 0)
	case e1.pt.Y == ip.Y && e2.pt.Y == ip.Y:
		return (ip.X-e1.pt.X > 0) != (ip.X-e2.pt.X > 0)
	case e1.other.pt.X == ip.X && e2.other.pt.X == ip.X:
		return (ip.Y-e1.other.pt.Y > 0) != (ip.Y-e2.other.pt.Y > 0)
	case e1.other.pt.Y == ip.Y && e2.other.pt.Y == ip.Y:
		return (ip.X-e1.other.pt.X > 0) != (ip.X-e2.other.pt.X > 0)
	}
	return true
}

func snap[T constraints.Float](pt geom.Point[T], e1, e2 *endpoint[T]) geom.Point[T] {
	pts := []geom.Point[T]{e1.pt, e2.pt, e1.other.pt, e2.other.pt}
	for _, p := range pts {
		if pt == p {
			return p
		}
	}
	slices.SortFunc(pts, func(a, b geom.Point[T]) int {
		switch {
		case a == b:
			return 0
		case ptIsBefore(a, b):
			return -1
		default:
			return 1
		}
	})
	pts[0], pts[2] = pts[2], pts[0]
	for _, p := range pts {
		if ptEqualWithin(pt, p, snapTolerance) {
			return p
		}
	}
	return pt
}

func ptEqualWithin[T constraints.Float](p1, p2 geom.Point[T], tol T) bool {
	return xmath.EqualWithin(p1.X, p2.X, tol) && xmath.EqualWithin(p1.Y, p2.Y, tol)
}
