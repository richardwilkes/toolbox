// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree

import "github.com/richardwilkes/toolbox/xmath/geom"

type node struct {
	rect      geom.Rect
	contents  []Node
	children  [4]*node
	threshold int
}

func (n *node) Bounds() geom.Rect {
	return n.rect
}

func (n *node) all(result []Node) []Node {
	result = append(result, n.contents...)
	if !n.isLeaf() {
		for _, child := range n.children {
			result = child.all(result)
		}
	}
	return result
}

func (n *node) isLeaf() bool {
	return n.children[0] == nil
}

func (n *node) insert(obj Node) {
	n.splitIfNeeded()
	if !n.isLeaf() {
		rect := obj.Bounds()
		for _, child := range n.children {
			if child.rect.ContainsRect(rect) {
				child.insert(obj)
				return
			}
		}
	}
	n.contents = append(n.contents, obj)
}

func (n *node) remove(obj Node) bool {
	for i, one := range n.contents {
		if one == obj {
			n.contents[i] = n.contents[len(n.contents)-1]
			n.contents[len(n.contents)-1] = nil
			n.contents = n.contents[:len(n.contents)-1]
			return true
		}
	}
	if !n.isLeaf() && n.rect.ContainsRect(obj.Bounds()) {
		for _, child := range n.children {
			if child.remove(obj) {
				return true
			}
		}
	}
	return false
}

func (n *node) splitIfNeeded() {
	if n.isLeaf() {
		if len(n.contents) >= n.threshold {
			hw := n.rect.Width / 2
			hh := n.rect.Height / 2
			n.children[0] = &node{
				rect: geom.Rect{
					Point: n.rect.Point,
					Size: geom.Size{
						Width:  hw,
						Height: hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[1] = &node{
				rect: geom.Rect{
					Point: geom.Point{
						X: n.rect.X + hw,
						Y: n.rect.Y,
					},
					Size: geom.Size{
						Width:  n.rect.Width - hw,
						Height: hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[2] = &node{
				rect: geom.Rect{
					Point: geom.Point{
						X: n.rect.X,
						Y: n.rect.Y + hh,
					},
					Size: geom.Size{
						Width:  hw,
						Height: n.rect.Height - hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[3] = &node{
				rect: geom.Rect{
					Point: geom.Point{
						X: n.rect.X + hw,
						Y: n.rect.Y + hh,
					},
					Size: geom.Size{
						Width:  n.rect.Width - hw,
						Height: n.rect.Height - hh,
					},
				},
				threshold: n.threshold,
			}
			contents := n.contents
			n.contents = make([]Node, 0)
			for _, one := range contents {
				n.insert(one)
			}
		}
	}
}

func (n *node) containsPoint(pt geom.Point) bool {
	if n.rect.ContainsPoint(pt) {
		for _, one := range n.contents {
			if one.Bounds().ContainsPoint(pt) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.containsPoint(pt) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findContainsPoint(pt geom.Point, result []Node) []Node {
	if n.rect.ContainsPoint(pt) {
		for _, one := range n.contents {
			if one.Bounds().ContainsPoint(pt) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findContainsPoint(pt, result)
			}
		}
	}
	return result
}

func (n *node) matchedContainsPoint(matcher Matcher, pt geom.Point) bool {
	if n.rect.ContainsPoint(pt) {
		for _, one := range n.contents {
			if one.Bounds().ContainsPoint(pt) && matcher.Matches(one) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.matchedContainsPoint(matcher, pt) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findMatchedContainsPoint(matcher Matcher, pt geom.Point, result []Node) []Node {
	if n.rect.ContainsPoint(pt) {
		for _, one := range n.contents {
			if one.Bounds().ContainsPoint(pt) && matcher.Matches(one) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findMatchedContainsPoint(matcher, pt, result)
			}
		}
	}
	return result
}

func (n *node) intersects(rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Intersects(rect) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.intersects(rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findIntersects(rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Intersects(rect) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findIntersects(rect, result)
			}
		}
	}
	return result
}

func (n *node) matchedIntersects(matcher Matcher, rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Intersects(rect) && matcher.Matches(one) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.matchedIntersects(matcher, rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findMatchedIntersects(matcher Matcher, rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Intersects(rect) && matcher.Matches(one) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findMatchedIntersects(matcher, rect, result)
			}
		}
	}
	return result
}

func (n *node) containsRect(rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().ContainsRect(rect) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.containsRect(rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findContainsRect(rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().ContainsRect(rect) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findContainsRect(rect, result)
			}
		}
	}
	return result
}

func (n *node) matchedContainsRect(matcher Matcher, rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().ContainsRect(rect) && matcher.Matches(one) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.matchedContainsRect(matcher, rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findMatchedContainsRect(matcher Matcher, rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().ContainsRect(rect) && matcher.Matches(one) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findMatchedContainsRect(matcher, rect, result)
			}
		}
	}
	return result
}

func (n *node) containedByRect(rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.ContainsRect(one.Bounds()) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.containedByRect(rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findContainedByRect(rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.ContainsRect(one.Bounds()) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findContainedByRect(rect, result)
			}
		}
	}
	return result
}

func (n *node) matchedContainedByRect(matcher Matcher, rect geom.Rect) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.ContainsRect(one.Bounds()) && matcher.Matches(one) {
				return true
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				if one.matchedContainedByRect(matcher, rect) {
					return true
				}
			}
		}
	}
	return false
}

func (n *node) findMatchedContainedByRect(matcher Matcher, rect geom.Rect, result []Node) []Node {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.ContainsRect(one.Bounds()) && matcher.Matches(one) {
				result = append(result, one)
			}
		}
		if !n.isLeaf() {
			for _, one := range n.children {
				result = one.findMatchedContainedByRect(matcher, rect, result)
			}
		}
	}
	return result
}
