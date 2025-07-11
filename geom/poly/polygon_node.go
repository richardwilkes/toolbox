// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

type vertexNode struct {
	next *vertexNode
	pt   Point
}

type polygonNode struct {
	left   *vertexNode
	right  *vertexNode
	next   *polygonNode
	proxy  *polygonNode
	active bool
}

func (p *polygonNode) addLeft(pt Point) {
	p.proxy.left = &vertexNode{
		pt:   pt,
		next: p.proxy.left,
	}
}

func (p *polygonNode) addRight(pt Point) {
	v := &vertexNode{pt: pt}
	if p.proxy == nil {
		p.proxy = p
	}
	if p.proxy.right != nil {
		p.proxy.right.next = v
	}
	p.proxy.right = v
}

func (p *polygonNode) mergeLeft(other, list *polygonNode) {
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

func (p *polygonNode) mergeRight(other, list *polygonNode) {
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

func (p *polygonNode) generate() Polygon {
	contourCount := 0
	ptCounts := make([]int, 0, 32)

	// Count the points of each contour and disable any that don't have enough points.
	for poly := p; poly != nil; poly = poly.next {
		if poly.active {
			var prev *vertexNode
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
		return Polygon{}
	}

	// Create the polygon
	result := make([]Contour, contourCount)
	ci := 0
	for poly := p; poly != nil; poly = poly.next {
		if !poly.active {
			continue
		}
		var prev *vertexNode
		result[ci] = make([]Point, ptCounts[ci])
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
