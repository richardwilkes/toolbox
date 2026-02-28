// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree

import (
	"github.com/richardwilkes/toolbox/v2/geom"
)

type node[N Node] struct {
	children  [4]*node[N]
	contents  []N
	threshold int
	rect      geom.Rect
}

func (n *node[N]) Bounds() geom.Rect {
	return n.rect
}

func (n *node[N]) all(result []N) []N {
	result = append(result, n.contents...)
	if !n.isLeaf() {
		for _, child := range n.children {
			result = child.all(result)
		}
	}
	return result
}

func (n *node[N]) isLeaf() bool {
	return n.children[0] == nil
}

func (n *node[N]) insert(obj N) {
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

func (n *node[N]) remove(obj N) bool {
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

func (n *node[N]) splitIfNeeded() {
	if n.isLeaf() {
		if len(n.contents) >= n.threshold {
			hw := n.rect.Width / 2
			hh := n.rect.Height / 2
			n.children[0] = &node[N]{
				rect: geom.Rect{
					Point: n.rect.Point,
					Size:  geom.NewSize(hw, hw),
				},
				threshold: n.threshold,
			}
			n.children[1] = &node[N]{
				rect:      geom.NewRect(n.rect.X+hw, n.rect.Y, n.rect.Width-hw, hh),
				threshold: n.threshold,
			}
			n.children[2] = &node[N]{
				rect:      geom.NewRect(n.rect.X, n.rect.Y+hh, hw, n.rect.Height-hh),
				threshold: n.threshold,
			}
			n.children[3] = &node[N]{
				rect:      geom.NewRect(n.rect.X+hw, n.rect.Y+hh, n.rect.Width-hw, n.rect.Height-hh),
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

func (n *node[N]) containsPoint(pt geom.Point) bool {
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

func (n *node[N]) findContainsPoint(pt geom.Point, result []N) []N {
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

func (n *node[N]) matchedContainsPoint(matcher Matcher[N], pt geom.Point) bool {
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

func (n *node[N]) findMatchedContainsPoint(matcher Matcher[N], pt geom.Point, result []N) []N {
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

func (n *node[N]) intersects(rect geom.Rect) bool {
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

func (n *node[N]) findIntersects(rect geom.Rect, result []N) []N {
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

func (n *node[N]) matchedIntersects(matcher Matcher[N], rect geom.Rect) bool {
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

func (n *node[N]) findMatchedIntersects(matcher Matcher[N], rect geom.Rect, result []N) []N {
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

func (n *node[N]) containsRect(rect geom.Rect) bool {
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

func (n *node[N]) findContainsRect(rect geom.Rect, result []N) []N {
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

func (n *node[N]) matchedContainsRect(matcher Matcher[N], rect geom.Rect) bool {
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

func (n *node[N]) findMatchedContainsRect(matcher Matcher[N], rect geom.Rect, result []N) []N {
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

func (n *node[N]) containedByRect(rect geom.Rect) bool {
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

func (n *node[N]) findContainedByRect(rect geom.Rect, result []N) []N {
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

func (n *node[N]) matchedContainedByRect(matcher Matcher[N], rect geom.Rect) bool {
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

func (n *node[N]) findMatchedContainedByRect(matcher Matcher[N], rect geom.Rect, result []N) []N {
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
