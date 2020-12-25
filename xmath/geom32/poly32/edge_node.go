// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly32

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom32"
)

const float32Epsilon = float32(2.2204460492503131e-16)

// Epsilon controls how close a point has to be to be considered the same as
// another. If performance becomes an issue, you can increase this value,
// which will cause more points to merge together, reducing the amount of work
// that must be done. Don't raise it too high, though, as it can distort the
// resulting polygon.
var Epsilon = float32Epsilon

type horizontalEdgeStates int

const (
	noHorizontalEdge horizontalEdgeStates = iota
	bottomHorizontalEdge
	topHorizontalEdge
)

var nextHorizontalEdgeStates = [3][6]horizontalEdgeStates{
	{bottomHorizontalEdge, topHorizontalEdge, topHorizontalEdge, bottomHorizontalEdge, noHorizontalEdge, noHorizontalEdge},
	{noHorizontalEdge, noHorizontalEdge, noHorizontalEdge, noHorizontalEdge, topHorizontalEdge, topHorizontalEdge},
	{noHorizontalEdge, noHorizontalEdge, noHorizontalEdge, noHorizontalEdge, bottomHorizontalEdge, bottomHorizontalEdge},
}

type bundleState int

const (
	unbundled bundleState = iota
	bundleHead
	bundleTail
)

type edgeNode struct {
	vertex      geom32.Point
	bot         geom32.Point
	top         geom32.Point
	xb          float32
	xt          float32
	dx          float32
	outAbove    *polygonNode
	outBelow    *polygonNode
	prev        *edgeNode
	next        *edgeNode
	pred        *edgeNode
	successor   *edgeNode
	nextBound   *edgeNode
	aboveState  bundleState
	belowState  bundleState
	which       int
	bundleAbove [2]bool
	bundleBelow [2]bool
	subjectSide bool
	clipSide    bool
}

type sortedEdge struct {
	edge *edgeNode
	xb   float32
	xt   float32
	dx   float32
	prev *sortedEdge
}

func (e *edgeNode) insertInto(b **edgeNode) {
	switch {
	case *b == nil:
		*b = e
	case e.bot.X < (*b).bot.X || (e.bot.X == (*b).bot.X && e.dx < (*b).dx):
		e.nextBound = *b
		*b = e
	default:
		e.insertInto(&(*b).nextBound)
	}
}

func (e *edgeNode) addEdgeToActiveEdgeTable(aet, prev *edgeNode) *edgeNode {
	switch {
	case aet == nil:
		aet = e
		e.prev = prev
		e.next = nil
	case e.xb < aet.xb || (e.xb == aet.xb && e.dx < aet.dx):
		e.prev = prev
		e.next = aet
		aet.prev = e
		aet = e
	default:
		aet.next = e.addEdgeToActiveEdgeTable(aet.next, aet)
	}
	return aet
}

func (e *edgeNode) addLocalMin(p *polygonNode, pt geom32.Point) *polygonNode {
	v := &vertexNode{pt: pt}
	result := &polygonNode{
		left:   v,
		right:  v,
		next:   p,
		active: true,
	}
	result.proxy = result
	e.outAbove = result
	return result
}

func (e *edgeNode) buildIntersections(dy float32) *intersection {
	var se *sortedEdge
	var it *intersection
	for edge := e; edge != nil; edge = edge.next {
		if edge.aboveState == bundleHead || edge.bundleAbove[clipping] || edge.bundleAbove[subject] {
			edge.addToSortedEdgeTable(&se, &it, dy)
		}
	}
	return it
}

func (e *edgeNode) addToSortedEdgeTable(se **sortedEdge, it **intersection, dy float32) {
	if *se == nil {
		*se = &sortedEdge{
			edge: e,
			xb:   e.xb,
			xt:   e.xt,
			dx:   e.dx,
		}
	} else {
		den := ((*se).xt - (*se).xb) - (e.xt - e.xb)
		if e.xt >= (*se).xt || e.dx == (*se).dx || xmath.AbsFloat32(den) <= float32Epsilon {
			*se = &sortedEdge{
				edge: e,
				xb:   e.xb,
				xt:   e.xt,
				dx:   e.dx,
				prev: *se,
			}
		} else {
			r := (e.xb - (*se).xb) / den
			addIntersection(it, (*se).edge, e, geom32.Point{
				X: (*se).xb + r*((*se).xt-(*se).xb),
				Y: r * dy,
			})
			e.addToSortedEdgeTable(&(*se).prev, it, dy)
		}
	}
}

func addIntersection(it **intersection, edge0, edge1 *edgeNode, pt geom32.Point) {
	switch {
	case *it == nil:
		*it = &intersection{
			edge0: edge0,
			edge1: edge1,
			point: pt,
		}
	case (*it).point.Y > pt.Y:
		*it = &intersection{
			edge0: edge0,
			edge1: edge1,
			point: pt,
			next:  *it,
		}
	default:
		addIntersection(&(*it).next, edge0, edge1, pt)
	}
}

func (e *edgeNode) bundleFields(pt geom32.Point) {
	updated := e
	e.bundleAbove[e.which] = e.top.Y != pt.Y
	e.bundleAbove[1-e.which] = false
	e.aboveState = unbundled
	for nextEdge := e.next; nextEdge != nil; nextEdge = nextEdge.next {
		nextEdge.bundleAbove[nextEdge.which] = nextEdge.top.Y != pt.Y
		nextEdge.bundleAbove[1-nextEdge.which] = false
		nextEdge.aboveState = unbundled
		if nextEdge.bundleAbove[nextEdge.which] {
			if mostlyEqual(updated.xb, nextEdge.xb) && mostlyEqual(updated.dx, nextEdge.dx) && updated.top.Y != pt.Y {
				nextEdge.bundleAbove[nextEdge.which] = nextEdge.bundleAbove[nextEdge.which] != updated.bundleAbove[nextEdge.which]
				nextEdge.bundleAbove[1-nextEdge.which] = updated.bundleAbove[1-nextEdge.which]
				nextEdge.aboveState = bundleHead
				updated.bundleAbove[clipping] = false
				updated.bundleAbove[subject] = false
				updated.aboveState = bundleTail
			}
			updated = nextEdge
		}
	}
}

func (e *edgeNode) process(op clipOp, pt geom32.Point, inPoly *polygonNode) (bPt geom32.Point, outPoly *polygonNode) {
	bPt = pt
	outPoly = inPoly
	var parityClipRight, paritySubjRight bool
	if op == subtractOp {
		parityClipRight = true
	}
	var horiz [2]horizontalEdgeStates
	var cf *polygonNode
	px := float32(-math.MaxFloat32)
	for edge := e; edge != nil; edge = edge.next {
		clipExistsState, clipExists := existsState(edge, clipping)
		subjExistsState, subjExists := existsState(edge, subject)
		if clipExists || subjExists {
			// Set bundle side
			edge.clipSide = parityClipRight
			edge.subjectSide = paritySubjRight

			// Determine contributing status and quadrant occupancies
			var br, bl, tr, tl, contributing bool
			pcb := parityClipRight != edge.bundleAbove[clipping]
			psb := paritySubjRight != edge.bundleAbove[subject]
			hc := horiz[clipping] != noHorizontalEdge
			hs := horiz[subject] != noHorizontalEdge
			phc := parityClipRight != hc
			phs := paritySubjRight != hs
			phcb := phc != edge.bundleBelow[clipping]
			phsb := phs != edge.bundleBelow[subject]
			switch op {
			case subtractOp, intersectOp:
				if contributing = (clipExists && (paritySubjRight || hs)) || (subjExists && (parityClipRight || hc)) || (clipExists && subjExists && parityClipRight == paritySubjRight); contributing {
					br = parityClipRight && paritySubjRight
					bl = pcb && psb
					tr = phc && phs
					tl = phcb && phsb
				}
			case xorOp:
				if contributing = clipExists || subjExists; contributing {
					br = parityClipRight != paritySubjRight
					bl = pcb != psb
					tr = phc != phs
					tl = phcb != phsb
				}
			case unionOp:
				if contributing = (clipExists && (!paritySubjRight || hs)) || (subjExists && (!parityClipRight || hc)) || (clipExists && subjExists && parityClipRight == paritySubjRight); contributing {
					br = parityClipRight || paritySubjRight
					bl = pcb || psb
					tr = phc || phs
					tl = phcb || phsb
				}
			default:
			}

			// Update parity
			parityClipRight = pcb
			paritySubjRight = psb

			// Update horizontal state
			if clipExists {
				horiz[clipping] = calcNextHState(clipExistsState, horiz[clipping], parityClipRight)
			}
			if subjExists {
				horiz[subject] = calcNextHState(subjExistsState, horiz[subject], paritySubjRight)
			}

			if contributing {
				bPt.X = edge.xb
				switch calcVertexType(tr, tl, br, bl) {
				case externalMinimum, internalMinimum:
					outPoly = edge.addLocalMin(outPoly, bPt)
					px = bPt.X
					cf = edge.outAbove
				case externalRightIntermediate:
					if cf != nil {
						if bPt.X != px {
							cf.addRight(bPt)
							px = bPt.X
						}
						edge.outAbove = cf
						cf = nil
					}
				case externalLeftIntermediate:
					edge.outBelow.addLeft(bPt)
					px = bPt.X
					cf = edge.outBelow
				case externalMaximum:
					if cf != nil {
						if bPt.X != px {
							cf.addLeft(bPt)
							px = bPt.X
						}
						cf.mergeRight(edge.outBelow, outPoly)
						cf = nil
					}
				case internalLeftIntermediate:
					if cf != nil {
						if bPt.X != px {
							cf.addLeft(bPt)
							px = bPt.X
						}
						edge.outAbove = cf
						cf = nil
					}
				case internalRightIntermediate:
					edge.outBelow.addRight(bPt)
					px = bPt.X
					cf = edge.outBelow
					edge.outBelow = nil
				case internalMaximum:
					if cf != nil {
						if bPt.X != px {
							cf.addRight(bPt)
							px = bPt.X
						}
						cf.mergeLeft(edge.outBelow, outPoly)
						cf = nil
						edge.outBelow = nil
					}
				case internalMaximumAndMinimum:
					if cf != nil {
						if bPt.X != px {
							cf.addRight(bPt)
							px = bPt.X
						}
						cf.mergeLeft(edge.outBelow, outPoly)
						edge.outBelow = nil
						outPoly = edge.addLocalMin(outPoly, bPt)
						cf = edge.outAbove
					}
				case externalMaximumAndMinimum:
					if cf != nil {
						if bPt.X != px {
							cf.addLeft(bPt)
							px = bPt.X
						}
						cf.mergeRight(edge.outBelow, outPoly)
						edge.outBelow = nil
						outPoly = edge.addLocalMin(outPoly, bPt)
						cf = edge.outAbove
					}
				case leftEdge:
					if edge.bot.Y == bPt.Y {
						edge.outBelow.addLeft(bPt)
					}
					edge.outAbove = edge.outBelow
					px = bPt.X
				case rightEdge:
					if edge.bot.Y == bPt.Y {
						edge.outBelow.addRight(bPt)
					}
					edge.outAbove = edge.outBelow
					px = bPt.X
				default:
				}
			}
		}
	}
	return
}

func (e *edgeNode) deleteTerminatingEdges(pt geom32.Point, yt float32) *edgeNode {
	updated := e
	for edge := e; edge != nil; edge = edge.next {
		switch {
		case edge.top.Y == pt.Y:
			prevEdge := edge.prev
			nextEdge := edge.next
			if prevEdge != nil {
				prevEdge.next = nextEdge
			} else {
				updated = nextEdge
			}
			if nextEdge != nil {
				nextEdge.prev = prevEdge
			}
			if edge.belowState == bundleHead && prevEdge != nil && prevEdge.belowState == bundleTail {
				prevEdge.outBelow = edge.outBelow
				prevEdge.belowState = unbundled
				if prevEdge.prev != nil && prevEdge.prev.belowState == bundleTail {
					prevEdge.belowState = bundleHead
				}
			}
		case edge.top.Y == yt:
			edge.xt = edge.top.X
		default:
			edge.xt = edge.bot.X + edge.dx*(yt-edge.bot.Y)
		}
	}
	return updated
}

func (e *edgeNode) prepareForNextScanBeam(yt float32) *edgeNode {
	updated := e
	for edge := e; edge != nil; edge = edge.next {
		successorEdge := edge.successor
		if edge.top.Y == yt && successorEdge != nil {
			successorEdge.outBelow = edge.outAbove
			successorEdge.belowState = edge.aboveState
			successorEdge.bundleBelow[clipping] = edge.bundleAbove[clipping]
			successorEdge.bundleBelow[subject] = edge.bundleAbove[subject]
			prevEdge := edge.prev
			if prevEdge != nil {
				prevEdge.next = successorEdge
			} else {
				updated = successorEdge
			}
			if edge.next != nil {
				edge.next.prev = successorEdge
			}
			successorEdge.prev = prevEdge
			successorEdge.next = edge.next
		} else {
			edge.outBelow = edge.outAbove
			edge.belowState = edge.aboveState
			edge.bundleBelow[clipping] = edge.bundleAbove[clipping]
			edge.bundleBelow[subject] = edge.bundleAbove[subject]
			edge.xb = edge.xt
		}
		edge.outAbove = nil
	}
	return updated
}

func (e *edgeNode) swapIntersectingEdgeBundles(inter *intersection) *edgeNode {
	result := e
	e0 := inter.edge0
	e1 := inter.edge1
	e0t := e0
	e1t := e1
	e0n := e0.next
	e1n := e1.next

	e0p := e0.prev
	if e0.aboveState == bundleHead {
		for {
			e0t = e0p
			e0p = e0p.prev
			if e0p == nil || e0p.aboveState != bundleTail {
				break
			}
		}
	}

	e1p := e1.prev
	if e1.aboveState == bundleHead {
		for {
			e1t = e1p
			e1p = e1p.prev
			if e1p == nil || e1p.aboveState != bundleTail {
				break
			}
		}
	}

	if e0p != nil {
		if e1p != nil {
			if e0p != e1 {
				e0p.next = e1t
				e1t.prev = e0p
			}
			if e1p != e0 {
				e1p.next = e0t
				e0t.prev = e1p
			}
		} else {
			if e0p != e1 {
				e0p.next = e1t
				e1t.prev = e0p
			}
			result = e0t
			e0t.prev = nil
		}
	} else {
		if e1p != e0 {
			if e1p != nil {
				e1p.next = e0t
			}
			e0t.prev = e1p
		}
		result = e1t
		e1t.prev = nil
	}

	if e0p != e1 {
		e0.next = e1n
		if e1n != nil {
			e1n.prev = e0
		}
	} else {
		e0.next = e1t
		e1t.prev = e0
	}

	if e1p != e0 {
		e1.next = e0n
		if e0n != nil {
			e0n.prev = e1
		}
	} else {
		e1.next = e0t
		e0t.prev = e1
	}

	return result
}

func mostlyEqual(a, b float32) bool {
	return xmath.AbsFloat32(a-b) <= Epsilon
}

func calcNextHState(existsState int, current horizontalEdgeStates, parityRight bool) horizontalEdgeStates {
	i := (existsState - 1) << 1
	if parityRight {
		i++
	}
	return nextHorizontalEdgeStates[current][i]
}
