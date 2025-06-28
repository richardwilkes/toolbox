// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree

import (
	"github.com/richardwilkes/toolbox/v2/xmath"
	"github.com/richardwilkes/toolbox/v2/xmath/geom"
)

type node[T xmath.Numeric, N Node[T]] struct {
	rect      geom.Rect[T]
	children  [4]*node[T, N]
	contents  []N
	threshold int
}

func (n *node[T, N]) Bounds() geom.Rect[T] {
	return n.rect
}

func (n *node[T, N]) all(result []N) []N {
	result = append(result, n.contents...)
	if !n.isLeaf() {
		for _, child := range n.children {
			result = child.all(result)
		}
	}
	return result
}

func (n *node[T, N]) isLeaf() bool {
	return n.children[0] == nil
}

func (n *node[T, N]) insert(obj N) {
	n.splitIfNeeded()
	if !n.isLeaf() {
		rect := obj.Bounds()
		for _, child := range n.children {
			if child.rect.Contains(rect) {
				child.insert(obj)
				return
			}
		}
	}
	n.contents = append(n.contents, obj)
}

func (n *node[T, N]) remove(obj N) bool {
	for i, one := range n.contents {
		if one != obj {
			continue
		}
		n.contents[i] = n.contents[len(n.contents)-1]
		var zero N
		n.contents[len(n.contents)-1] = zero
		n.contents = n.contents[:len(n.contents)-1]
		return true
	}
	if !n.isLeaf() && n.rect.Contains(obj.Bounds()) {
		for _, child := range n.children {
			if child.remove(obj) {
				return true
			}
		}
	}
	return false
}

func (n *node[T, N]) splitIfNeeded() {
	if n.isLeaf() {
		if len(n.contents) >= n.threshold {
			hw := n.rect.Width / 2
			hh := n.rect.Height / 2
			n.children[0] = &node[T, N]{
				rect: geom.Rect[T]{
					Point: n.rect.Point,
					Size:  geom.NewSize[T](hw, hw),
				},
				threshold: n.threshold,
			}
			n.children[1] = &node[T, N]{
				rect:      geom.NewRect[T](n.rect.X+hw, n.rect.Y, n.rect.Width-hw, hh),
				threshold: n.threshold,
			}
			n.children[2] = &node[T, N]{
				rect:      geom.NewRect[T](n.rect.X, n.rect.Y+hh, hw, n.rect.Height-hh),
				threshold: n.threshold,
			}
			n.children[3] = &node[T, N]{
				rect:      geom.NewRect[T](n.rect.X+hw, n.rect.Y+hh, n.rect.Width-hw, n.rect.Height-hh),
				threshold: n.threshold,
			}
			contents := n.contents
			n.contents = nil
			for _, one := range contents {
				n.insert(one)
			}
		}
	}
}

func (n *node[T, N]) containsPoint(pt geom.Point[T]) bool {
	if pt.In(n.rect) {
		for _, one := range n.contents {
			if pt.In(one.Bounds()) {
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

func (n *node[T, N]) findContainsPoint(pt geom.Point[T], result []N) []N {
	if pt.In(n.rect) {
		for _, one := range n.contents {
			if pt.In(one.Bounds()) {
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

func (n *node[T, N]) matchedContainsPoint(matcher Matcher[T, N], pt geom.Point[T]) bool {
	if pt.In(n.rect) {
		for _, one := range n.contents {
			if pt.In(one.Bounds()) && matcher.Matches(one) {
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

func (n *node[T, N]) findMatchedContainsPoint(matcher Matcher[T, N], pt geom.Point[T], result []N) []N {
	if pt.In(n.rect) {
		for _, one := range n.contents {
			if pt.In(one.Bounds()) && matcher.Matches(one) {
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

func (n *node[T, N]) intersects(rect geom.Rect[T]) bool {
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

func (n *node[T, N]) findIntersects(rect geom.Rect[T], result []N) []N {
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

func (n *node[T, N]) matchedIntersects(matcher Matcher[T, N], rect geom.Rect[T]) bool {
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

func (n *node[T, N]) findMatchedIntersects(matcher Matcher[T, N], rect geom.Rect[T], result []N) []N {
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

func (n *node[T, N]) containsRect(rect geom.Rect[T]) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Contains(rect) {
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

func (n *node[T, N]) findContainsRect(rect geom.Rect[T], result []N) []N {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Contains(rect) {
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

func (n *node[T, N]) matchedContainsRect(matcher Matcher[T, N], rect geom.Rect[T]) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Contains(rect) && matcher.Matches(one) {
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

func (n *node[T, N]) findMatchedContainsRect(matcher Matcher[T, N], rect geom.Rect[T], result []N) []N {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if one.Bounds().Contains(rect) && matcher.Matches(one) {
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

func (n *node[T, N]) containedByRect(rect geom.Rect[T]) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.Contains(one.Bounds()) {
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

func (n *node[T, N]) findContainedByRect(rect geom.Rect[T], result []N) []N {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.Contains(one.Bounds()) {
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

func (n *node[T, N]) matchedContainedByRect(matcher Matcher[T, N], rect geom.Rect[T]) bool {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.Contains(one.Bounds()) && matcher.Matches(one) {
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

func (n *node[T, N]) findMatchedContainedByRect(matcher Matcher[T, N], rect geom.Rect[T], result []N) []N {
	if n.rect.Intersects(rect) {
		for _, one := range n.contents {
			if rect.Contains(one.Bounds()) && matcher.Matches(one) {
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
