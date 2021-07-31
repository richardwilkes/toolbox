// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
)

const (
	clipping = iota
	subject
)

type clipOp int

const (
	subtractOp clipOp = iota
	intersectOp
	xorOp
	unionOp
)

type vertexType int

const (
	emptyNonIntersection vertexType = iota //nolint:deadcode,varcheck
	externalMaximum
	externalLeftIntermediate
	topEdge //nolint:deadcode,varcheck
	externalRightIntermediate
	rightEdge
	internalMaximumAndMinimum
	internalMinimum
	externalMinimum
	externalMaximumAndMinimum
	leftEdge
	internalLeftIntermediate
	bottomEdge //nolint:deadcode,varcheck
	internalRightIntermediate
	internalMaximum
	fullNonIntersection //nolint:deadcode,varcheck
)

// Polygon holds one or more contour lines. The polygon may contain holes and
// may be self-intersecting.
type Polygon []Contour

// CalcEllipseSegmentCount returns a suggested number of segments to use when
// generating an ellipse. 'r' is the largest radius of the ellipse. 'e' is the
// acceptable error, typically 1 or less.
func CalcEllipseSegmentCount(r, e float64) int {
	d := 1 - e/r
	n := int(math.Ceil(2 * math.Pi / math.Acos(2*d*d-1)))
	if n < 4 {
		n = 4
	}
	return n
}

// ApproximateEllipseAuto creates a polygon that approximates an ellipse,
// automatically choose the number of segments to break the ellipse contour
// into. This uses CalcEllipseSegmentCount() with an 'e' of 0.2.
func ApproximateEllipseAuto(bounds geom.Rect) Polygon {
	return ApproximateEllipse(bounds, CalcEllipseSegmentCount(math.Max(bounds.Width, bounds.Height)/2, 0.2))
}

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
		geom.Point{X: bounds.X, Y: bounds.Bottom() - 1},
		geom.Point{X: bounds.Right() - 1, Y: bounds.Bottom() - 1},
		geom.Point{X: bounds.Right() - 1, Y: bounds.Y},
	}}
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
	b := p[0].Bounds()
	for _, c := range p[1:] {
		b.Union(c.Bounds())
	}
	return b
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
func (p Polygon) Union(other Polygon) Polygon {
	return p.construct(unionOp, other)
}

// Intersect returns the intersection of both polygons.
func (p Polygon) Intersect(other Polygon) Polygon {
	return p.construct(intersectOp, other)
}

// Subtract returns the result of removing the other polygon from this
// polygon.
func (p Polygon) Subtract(other Polygon) Polygon {
	return p.construct(subtractOp, other)
}

// Xor returns the result of xor'ing this polygon with the other polygon.
func (p Polygon) Xor(other Polygon) Polygon {
	return p.construct(xorOp, other)
}

func (p Polygon) construct(op clipOp, other Polygon) Polygon {
	var result Polygon

	// Short-circuit the work if we can trivially determine the result is an
	// empty polygon.
	if (len(p) == 0 && len(other) == 0) ||
		(len(p) == 0 && (op == intersectOp || op == subtractOp)) ||
		(len(other) == 0 && op == intersectOp) {
		return result
	}

	// Build the local minima table and the scan beam table
	sbTree := &scanBeamTree{}
	subjNonContributing, clipNonContributing := p.identifyNonContributingContours(op, other)
	lmt := buildLocalMinimaTable(nil, sbTree, p, subjNonContributing, subject, op)
	if lmt = buildLocalMinimaTable(lmt, sbTree, other, clipNonContributing, clipping, op); lmt == nil {
		return result
	}
	sbt := sbTree.buildScanBeamTable()

	// Process each scan beam
	var aet *edgeNode
	var outPoly *polygonNode
	localMin := lmt
	i := 0
	for i < len(sbt) {

		// Set yb and yt to the bottom and top of the scanbeam
		var yt, dy float64
		var bPt geom.Point
		bPt.Y = sbt[i]
		i++
		if i < len(sbt) {
			yt = sbt[i]
			dy = yt - bPt.Y
		}

		// If LMT node corresponding to bPt.Y exists
		if localMin != nil && localMin.y == bPt.Y {
			// Add edges starting at this local minimum to the AET
			for edge := localMin.firstBound; edge != nil; edge = edge.nextBound {
				aet = edge.addEdgeToActiveEdgeTable(aet, nil)
			}
			localMin = localMin.next
		}
		if aet == nil {
			continue
		}

		aet.bundleFields(bPt)
		bPt, outPoly = aet.process(op, bPt, outPoly)
		aet = aet.deleteTerminatingEdges(bPt, yt)

		if i < len(sbt) {
			// Process each node in the intersection table
			for inter := aet.buildIntersections(dy); inter != nil; inter = inter.next {
				outPoly = inter.process(op, bPt, outPoly)
				aet = aet.swapIntersectingEdgeBundles(inter)
			}
			aet = aet.prepareForNextScanBeam(yt)
		}
	}

	// Generate the resulting polygon
	if outPoly != nil {
		return outPoly.generate()
	}
	return Polygon{}
}

func (p Polygon) identifyNonContributingContours(op clipOp, clip Polygon) (subjNonContributing, clipNonContributing []bool) {
	subjNonContributing = make([]bool, len(p))
	clipNonContributing = make([]bool, len(clip))
	if (op == intersectOp || op == subtractOp) && len(p) > 0 && len(clip) > 0 {

		// Check all subject contour bounding boxes against clip boxes
		overlaps := make([]bool, len(p)*len(clip))
		boxes := make([]geom.Rect, len(clip))
		for i, c := range clip {
			boxes[i] = c.Bounds()
		}
		for si := range p {
			box := p[si].Bounds()
			for ci := range clip {
				overlaps[ci*len(p)+si] = box.Intersects(boxes[ci])
			}
		}

		// For each clip contour, search for any subject contour overlaps
		for ci := range clip {
			clipNonContributing[ci] = true
			for si := range p {
				if overlaps[ci*len(p)+si] {
					clipNonContributing[ci] = false
					break
				}
			}
		}

		if op == intersectOp {
			// For each subject contour, search for any clip contour overlaps
			for si := range p {
				subjNonContributing[si] = true
				for ci := range clip {
					if overlaps[ci*len(p)+si] {
						subjNonContributing[si] = false
						break
					}
				}
			}
		}
	}
	return
}

func calcVertexType(tr, tl, br, bl bool) vertexType {
	var vt vertexType
	if tr {
		vt = 1
	}
	if tl {
		vt |= 2
	}
	if br {
		vt |= 4
	}
	if bl {
		vt |= 8
	}
	return vt
}

func existsState(edge *edgeNode, which int) (int, bool) {
	state := 0
	if edge.bundleAbove[which] {
		state = 1
	}
	if edge.bundleBelow[which] {
		state |= 2
	}
	return state, edge.bundleAbove[which] || edge.bundleBelow[which]
}
