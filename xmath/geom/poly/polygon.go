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
	"strings"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
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

// Polygon holds one or more contour lines. The polygon may contain holes and may be self-intersecting.
type Polygon[T constraints.Float] []Contour[T]

// Clone returns a duplicate of this polygon.
func (p Polygon[T]) Clone() Polygon[T] {
	if len(p) == 0 {
		return nil
	}
	clone := Polygon[T](make([]Contour[T], len(p)))
	for i := range p {
		clone[i] = p[i].Clone()
	}
	return clone
}

// String implements fmt.Stringer.
func (p Polygon[T]) String() string {
	var buffer strings.Builder
	buffer.WriteByte('{')
	for i, c := range p {
		if i != 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(c.String())
	}
	buffer.WriteByte('}')
	return buffer.String()
}

// Bounds returns the bounding rectangle of this polygon.
func (p Polygon[T]) Bounds() geom.Rect[T] {
	if len(p) == 0 {
		return geom.Rect[T]{}
	}
	b := p[0].Bounds()
	for _, c := range p[1:] {
		b = b.Union(c.Bounds())
	}
	return b
}

// Contains returns true if the point is contained by this polygon.
func (p Polygon[T]) Contains(pt geom.Point[T]) bool {
	for i := range p {
		if p[i].Contains(pt) {
			return true
		}
	}
	return false
}

// ContainsEvenOdd returns true if the point is contained by the polygon using the even-odd rule.
// https://en.wikipedia.org/wiki/Even-odd_rule
func (p Polygon[T]) ContainsEvenOdd(pt geom.Point[T]) bool {
	var count int
	for i := range p {
		if p[i].Contains(pt) {
			count++
		}
	}
	return count%2 == 1
}

// Transform returns the result of transforming this Polygon by the Matrix.
func (p Polygon[T]) Transform(m geom.Matrix[T]) Polygon[T] {
	clone := p.Clone()
	for _, c := range clone {
		for i := range c {
			c[i] = m.TransformPoint(c[i])
		}
	}
	return clone
}

// Union returns a new Polygon holding the union of both Polygons.
func (p Polygon[T]) Union(other Polygon[T]) Polygon[T] {
	return p.construct(unionOp, other)
}

// Intersect returns a new Polygon holding the intersection of both Polygons.
func (p Polygon[T]) Intersect(other Polygon[T]) Polygon[T] {
	return p.construct(intersectOp, other)
}

// Sub returns a new Polygon holding the result of removing the other Polygon from this Polygon.
func (p Polygon[T]) Sub(other Polygon[T]) Polygon[T] {
	return p.construct(subtractOp, other)
}

// Xor returns a new Polygon holding the result of xor'ing this Polygon with the other Polygon.
func (p Polygon[T]) Xor(other Polygon[T]) Polygon[T] {
	return p.construct(xorOp, other)
}

func (p Polygon[T]) construct(op clipOp, other Polygon[T]) Polygon[T] {
	var result Polygon[T]

	// Short-circuit the work if we can trivially determine the result is an empty polygon.
	if (len(p) == 0 && len(other) == 0) ||
		(len(p) == 0 && (op == intersectOp || op == subtractOp)) ||
		(len(other) == 0 && op == intersectOp) {
		return result
	}

	// Build the local minima table and the scan beam table
	sbTree := &scanBeamTree[T]{}
	subjNonContributing, clipNonContributing := p.identifyNonContributingContours(op, other)
	lmt := buildLocalMinimaTable(nil, sbTree, p, subjNonContributing, subject, op)
	if lmt = buildLocalMinimaTable(lmt, sbTree, other, clipNonContributing, clipping, op); lmt == nil {
		return result
	}
	sbt := sbTree.buildScanBeamTable()

	// Process each scan beam
	var aet *edgeNode[T]
	var outPoly *polygonNode[T]
	localMin := lmt
	i := 0
	for i < len(sbt) {

		// Set yb and yt to the bottom and top of the scanbeam
		var yt, dy T
		var bPt geom.Point[T]
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
	return Polygon[T]{}
}

func (p Polygon[T]) identifyNonContributingContours(op clipOp, clip Polygon[T]) (subjNC, clipNC []bool) {
	subjNC = make([]bool, len(p))
	clipNC = make([]bool, len(clip))
	if (op == intersectOp || op == subtractOp) && len(p) > 0 && len(clip) > 0 {

		// Check all subject contour bounding boxes against clip boxes
		overlaps := make([]bool, len(p)*len(clip))
		boxes := make([]geom.Rect[T], len(clip))
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
			clipNC[ci] = true
			for si := range p {
				if overlaps[ci*len(p)+si] {
					clipNC[ci] = false
					break
				}
			}
		}

		if op == intersectOp {
			// For each subject contour, search for any clip contour overlaps
			for si := range p {
				subjNC[si] = true
				for ci := range clip {
					if overlaps[ci*len(p)+si] {
						subjNC[si] = false
						break
					}
				}
			}
		}
	}
	return
}
