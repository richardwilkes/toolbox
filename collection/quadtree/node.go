// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree

import (
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

type node[T xmath.Numeric] struct {
	rect      geom.Rect[T]
	contents  []Node[T]
	children  [4]*node[T]
	threshold int
}

func (n *node[T]) Bounds() geom.Rect[T] {
	return n.rect
}

func (n *node[T]) all(result []Node[T]) []Node[T] {
	result = append(result, n.contents...)
	if !n.isLeaf() {
		for _, child := range n.children {
			result = child.all(result)
		}
	}
	return result
}

func (n *node[T]) isLeaf() bool {
	return n.children[0] == nil
}

func (n *node[T]) insert(obj Node[T]) {
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

func (n *node[T]) remove(obj Node[T]) bool {
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

func (n *node[T]) splitIfNeeded() {
	if n.isLeaf() {
		if len(n.contents) >= n.threshold {
			hw := n.rect.Width / 2
			hh := n.rect.Height / 2
			n.children[0] = &node[T]{
				rect: geom.Rect[T]{
					Point: n.rect.Point,
					Size: geom.Size[T]{
						Width:  hw,
						Height: hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[1] = &node[T]{
				rect: geom.Rect[T]{
					Point: geom.Point[T]{
						X: n.rect.X + hw,
						Y: n.rect.Y,
					},
					Size: geom.Size[T]{
						Width:  n.rect.Width - hw,
						Height: hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[2] = &node[T]{
				rect: geom.Rect[T]{
					Point: geom.Point[T]{
						X: n.rect.X,
						Y: n.rect.Y + hh,
					},
					Size: geom.Size[T]{
						Width:  hw,
						Height: n.rect.Height - hh,
					},
				},
				threshold: n.threshold,
			}
			n.children[3] = &node[T]{
				rect: geom.Rect[T]{
					Point: geom.Point[T]{
						X: n.rect.X + hw,
						Y: n.rect.Y + hh,
					},
					Size: geom.Size[T]{
						Width:  n.rect.Width - hw,
						Height: n.rect.Height - hh,
					},
				},
				threshold: n.threshold,
			}
			contents := n.contents
			n.contents = make([]Node[T], 0)
			for _, one := range contents {
				n.insert(one)
			}
		}
	}
}

func (n *node[T]) containsPoint(pt geom.Point[T]) bool {
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

func (n *node[T]) findContainsPoint(pt geom.Point[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) matchedContainsPoint(matcher Matcher[T], pt geom.Point[T]) bool {
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

func (n *node[T]) findMatchedContainsPoint(matcher Matcher[T], pt geom.Point[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) intersects(rect geom.Rect[T]) bool {
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

func (n *node[T]) findIntersects(rect geom.Rect[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) matchedIntersects(matcher Matcher[T], rect geom.Rect[T]) bool {
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

func (n *node[T]) findMatchedIntersects(matcher Matcher[T], rect geom.Rect[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) containsRect(rect geom.Rect[T]) bool {
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

func (n *node[T]) findContainsRect(rect geom.Rect[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) matchedContainsRect(matcher Matcher[T], rect geom.Rect[T]) bool {
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

func (n *node[T]) findMatchedContainsRect(matcher Matcher[T], rect geom.Rect[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) containedByRect(rect geom.Rect[T]) bool {
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

func (n *node[T]) findContainedByRect(rect geom.Rect[T], result []Node[T]) []Node[T] {
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

func (n *node[T]) matchedContainedByRect(matcher Matcher[T], rect geom.Rect[T]) bool {
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

func (n *node[T]) findMatchedContainedByRect(matcher Matcher[T], rect geom.Rect[T], result []Node[T]) []Node[T] {
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
