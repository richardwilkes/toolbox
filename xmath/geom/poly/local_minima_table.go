/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package poly

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

type localMinimaNode[T constraints.Float] struct {
	y          T
	firstBound *edgeNode[T]
	next       *localMinimaNode[T]
}

func buildLocalMinimaTable[T constraints.Float](lmt *localMinimaNode[T], sbTree *scanBeamTree[T], p Polygon[T], nc []bool, which int, op clipOp) *localMinimaNode[T] {
	if len(p) == 0 {
		return lmt
	}
	count := 0
	for ci := range p {
		if !nc[ci] {
			for v := range p[ci] {
				if optimal(p[ci], v, len(p[ci])) {
					count++
				}
			}
		}
	}
	edges := make([]edgeNode[T], count)
	edgeIndex := 0
	for ci := range p {
		if !nc[ci] {
			// Perform contour optimization
			count = 0
			for v := range p[ci] {
				if optimal(p[ci], v, len(p[ci])) {
					edges[count].vertex = p[ci][v]
					sbTree.add(edges[count].vertex.Y)
					count++
				}
			}

			// Do the contour forward pass
			for minimum := 0; minimum < count; minimum++ {
				if edges[previousIndex(minimum, count)].vertex.Y < edges[minimum].vertex.Y ||
					edges[nextIndex(minimum, count)].vertex.Y <= edges[minimum].vertex.Y {
					continue
				}

				// Search for the next local maximum
				edgeCount := 1
				maximum := nextIndex(minimum, count)
				for edges[nextIndex(maximum, count)].vertex.Y > edges[maximum].vertex.Y {
					edgeCount++
					maximum = nextIndex(maximum, count)
				}

				// Build the next edge list
				e := &edges[edgeIndex]
				e.belowState = unbundled
				e.bundleBelow[clipping] = false
				e.bundleBelow[subject] = false
				vi := minimum
				for i := 0; i < edgeCount; i++ {
					e = &edges[edgeIndex+i]
					v := &edges[vi]
					e.xb = v.vertex.X
					e.bot = v.vertex
					vi = nextIndex(vi, count)
					v = &edges[vi]
					e.top = v.vertex
					e.dx = (v.vertex.X - e.bot.X) / (e.top.Y - e.bot.Y)
					e.which = which
					e.outAbove = nil
					e.outBelow = nil
					e.next = nil
					e.prev = nil
					if edgeCount > 1 && i < edgeCount-1 {
						e.successor = &edges[edgeIndex+i+1]
					} else {
						e.successor = nil
					}
					if edgeCount > 1 && i > 0 {
						e.pred = &edges[edgeIndex+i-1]
					} else {
						e.pred = nil
					}
					e.nextBound = nil
					e.clipSide = op == subtractOp
					e.subjectSide = false
				}
				lmt = lmt.insertBound(edges[minimum].vertex.Y, &edges[edgeIndex])
				edgeIndex += edgeCount
			}

			// Do the contour reverse pass
			for minimum := 0; minimum < count; minimum++ {
				if edges[previousIndex(minimum, count)].vertex.Y <= edges[minimum].vertex.Y ||
					edges[nextIndex(minimum, count)].vertex.Y < edges[minimum].vertex.Y {
					continue
				}
				// Search for the previous local maximum
				edgeCount := 1
				maximum := previousIndex(minimum, count)
				for edges[previousIndex(maximum, count)].vertex.Y > edges[maximum].vertex.Y {
					edgeCount++
					maximum = previousIndex(maximum, count)
				}

				// Build the previous edge list
				e := &edges[edgeIndex]
				e.belowState = unbundled
				e.bundleBelow[clipping] = false
				e.bundleBelow[subject] = false
				vi := minimum
				for i := 0; i < edgeCount; i++ {
					e = &edges[edgeIndex+i]
					v := &edges[vi]
					e.xb = v.vertex.X
					e.bot = v.vertex
					vi = previousIndex(vi, count)
					v = &edges[vi]
					e.top = v.vertex
					e.dx = (v.vertex.X - e.bot.X) / (e.top.Y - e.bot.Y)
					e.which = which
					e.outAbove = nil
					e.outBelow = nil
					e.next = nil
					e.prev = nil
					if edgeCount > 1 && i < edgeCount-1 {
						e.successor = &edges[edgeIndex+i+1]
					} else {
						e.successor = nil
					}
					if edgeCount > 1 && i > 0 {
						e.pred = &edges[edgeIndex+i-1]
					} else {
						e.pred = nil
					}
					e.nextBound = nil
					e.clipSide = op == subtractOp
					e.subjectSide = false
				}
				lmt = lmt.insertBound(edges[minimum].vertex.Y, &edges[edgeIndex])
				edgeIndex += edgeCount
			}
		}
	}
	return lmt
}

func (n *localMinimaNode[T]) insertBound(y T, e *edgeNode[T]) *localMinimaNode[T] {
	lmn, en := n.boundList(y)
	e.insertInto(en)
	return lmn
}

func (n *localMinimaNode[T]) boundList(y T) (lmn *localMinimaNode[T], en **edgeNode[T]) {
	switch {
	case n == nil:
		lmn = &localMinimaNode[T]{y: y}
		return lmn, &lmn.firstBound
	case y < n.y:
		lmn = &localMinimaNode[T]{
			y:    y,
			next: n,
		}
		return lmn, &lmn.firstBound
	case y > n.y:
		n.next, en = n.next.boundList(y)
		return n, en
	default:
		return n, &n.firstBound
	}
}

func optimal[T constraints.Float](v []geom.Point[T], i, n int) bool {
	return v[previousIndex(i, n)].Y != v[i].Y || v[nextIndex(i, n)].Y != v[i].Y
}

func previousIndex(i, n int) int {
	return (i - 1 + n) % n
}

func nextIndex(i, n int) int {
	return (i + 1) % n
}
