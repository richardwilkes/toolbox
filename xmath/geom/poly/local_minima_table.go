// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type localMinimaNode struct {
	y          float64
	firstBound *edgeNode
	next       *localMinimaNode
}

func buildLocalMinimaTable(lmt *localMinimaNode, sbTree *scanBeamTree, p Polygon, nc []bool, which int, op clipOp) *localMinimaNode {
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
	edges := make([]edgeNode, count)
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
			for min := 0; min < count; min++ {
				if edges[previousIndex(min, count)].vertex.Y >= edges[min].vertex.Y && //nolint:gocritic
					edges[nextIndex(min, count)].vertex.Y > edges[min].vertex.Y {

					// Search for the next local maximum
					edgeCount := 1
					max := nextIndex(min, count)
					for edges[nextIndex(max, count)].vertex.Y > edges[max].vertex.Y {
						edgeCount++
						max = nextIndex(max, count)
					}

					// Build the next edge list
					e := &edges[edgeIndex]
					e.belowState = unbundled
					e.bundleBelow[clipping] = false
					e.bundleBelow[subject] = false
					vi := min
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
					lmt = lmt.insertBound(edges[min].vertex.Y, &edges[edgeIndex])
					edgeIndex += edgeCount
				}
			}

			// Do the contour reverse pass
			for min := 0; min < count; min++ {
				if edges[previousIndex(min, count)].vertex.Y > edges[min].vertex.Y && //nolint:gocritic
					edges[nextIndex(min, count)].vertex.Y >= edges[min].vertex.Y {

					// Search for the previous local maximum
					edgeCount := 1
					max := previousIndex(min, count)
					for edges[previousIndex(max, count)].vertex.Y > edges[max].vertex.Y {
						edgeCount++
						max = previousIndex(max, count)
					}

					// Build the previous edge list
					e := &edges[edgeIndex]
					e.belowState = unbundled
					e.bundleBelow[clipping] = false
					e.bundleBelow[subject] = false
					vi := min
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
					lmt = lmt.insertBound(edges[min].vertex.Y, &edges[edgeIndex])
					edgeIndex += edgeCount
				}
			}
		}
	}
	return lmt
}

func (n *localMinimaNode) insertBound(y float64, e *edgeNode) *localMinimaNode {
	lmn, en := n.boundList(y)
	e.insertInto(en)
	return lmn
}

func (n *localMinimaNode) boundList(y float64) (lmn *localMinimaNode, en **edgeNode) {
	switch {
	case n == nil:
		lmn = &localMinimaNode{y: y}
		return lmn, &lmn.firstBound
	case y < n.y:
		lmn = &localMinimaNode{
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

func optimal(v []geom.Point, i, n int) bool {
	return v[previousIndex(i, n)].Y != v[i].Y || v[nextIndex(i, n)].Y != v[i].Y
}

func previousIndex(i, n int) int {
	return (i - 1 + n) % n
}

func nextIndex(i, n int) int {
	return (i + 1) % n
}
