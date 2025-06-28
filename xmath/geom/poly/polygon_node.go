// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"github.com/richardwilkes/toolbox/v2/xmath/geom"
)

type vertexNode[T ~float32 | ~float64] struct {
	pt   geom.Point[T]
	next *vertexNode[T]
}

type polygonNode[T ~float32 | ~float64] struct {
	left   *vertexNode[T]
	right  *vertexNode[T]
	next   *polygonNode[T]
	proxy  *polygonNode[T]
	active bool
}

func (p *polygonNode[T]) addLeft(pt geom.Point[T]) {
	p.proxy.left = &vertexNode[T]{
		pt:   pt,
		next: p.proxy.left,
	}
}

func (p *polygonNode[T]) addRight(pt geom.Point[T]) {
	v := &vertexNode[T]{pt: pt}
	if p.proxy == nil {
		p.proxy = p
	}
	if p.proxy.right != nil {
		p.proxy.right.next = v
	}
	p.proxy.right = v
}

func (p *polygonNode[T]) mergeLeft(other, list *polygonNode[T]) {
	if other != nil && p.proxy != other.proxy {
		p.proxy.right.next = other.proxy.left
		other.proxy.left = p.proxy.left
		for target := p.proxy; list != nil; list = list.next {
			if list.proxy == target {
				list.active = false
				list.proxy = other.proxy
			}
		}
	}
}

func (p *polygonNode[T]) mergeRight(other, list *polygonNode[T]) {
	if other != nil && p.proxy != other.proxy {
		other.proxy.right.next = p.proxy.left
		other.proxy.right = p.proxy.right
		for target := p.proxy; list != nil; list = list.next {
			if list.proxy == target {
				list.active = false
				list.proxy = other.proxy
			}
		}
	}
}

func (p *polygonNode[T]) generate() Polygon[T] {
	contourCount := 0
	ptCounts := make([]int, 0, 32)

	// Count the points of each contour and disable any that don't have enough points.
	for poly := p; poly != nil; poly = poly.next {
		if poly.active {
			var prev *vertexNode[T]
			ptCount := 0
			for v := poly.proxy.left; v != nil; v = v.next {
				if prev == nil || prev.pt != v.pt {
					ptCount++
				}
				prev = v
			}
			if ptCount > 2 {
				ptCounts = append(ptCounts, ptCount)
				contourCount++
			} else {
				poly.active = false
			}
		}
	}
	if contourCount == 0 {
		return Polygon[T]{}
	}

	// Create the polygon
	result := make([]Contour[T], contourCount)
	ci := 0
	for poly := p; poly != nil; poly = poly.next {
		if !poly.active {
			continue
		}
		var prev *vertexNode[T]
		result[ci] = make([]geom.Point[T], ptCounts[ci])
		v := len(result[ci]) - 1
		for vtx := poly.proxy.left; vtx != nil; vtx = vtx.next {
			if prev == nil || prev.pt != vtx.pt {
				result[ci][v] = vtx.pt
				v--
			}
			prev = vtx
		}
		ci++
	}
	return result
}
