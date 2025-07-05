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

const (
	// DefaultQuadTreeThreshold is the default threshold that will be used if none is specified.
	DefaultQuadTreeThreshold = 64
	// MinQuadTreeThreshold is the minimum allowed threshold.
	MinQuadTreeThreshold = 4
)

// Node defines the methods an object that can be stored within the QuadTree must implement.
type Node[T xmath.Numeric] interface {
	comparable
	// Bounds returns the node's bounding rectangle.
	Bounds() geom.Rect[T]
}

// Matcher is used to match nodes.
type Matcher[T xmath.Numeric, N Node[T]] interface {
	// Matches returns true if the node matches.
	Matches(n N) bool
}

// QuadTree stores two-dimensional nodes for fast lookup.
type QuadTree[T xmath.Numeric, N Node[T]] struct {
	root      *node[T, N]
	outside   []N
	Threshold int
	count     int
}

// Size returns the number of nodes contained within the QuadTree.
func (q *QuadTree[T, N]) Size() int {
	return q.count
}

func (q *QuadTree[T, N]) threshold() int {
	if q.Threshold < MinQuadTreeThreshold {
		return DefaultQuadTreeThreshold
	}
	return q.Threshold
}

// Insert a node. NOTE: Once a node is inserted, the value it returns from a call to Bounds() MUST REMAIN THE SAME until
// the node is removed.
func (q *QuadTree[T, N]) Insert(n N) {
	rect := n.Bounds()
	if rect.Empty() {
		return
	}
	q.count++
	if q.root != nil && q.root.rect.Contains(rect) {
		q.root.insert(n)
	} else {
		q.outside = append(q.outside, n)
		if len(q.outside) > q.threshold() {
			q.Reorganize()
		}
	}
}

// Remove a node.
func (q *QuadTree[T, N]) Remove(n N) {
	for i, one := range q.outside {
		if one != n {
			continue
		}
		q.outside[i] = q.outside[len(q.outside)-1]
		var zero N
		q.outside[len(q.outside)-1] = zero
		q.outside = q.outside[:len(q.outside)-1]
		q.count--
		return
	}
	if q.root != nil {
		if q.root.remove(n) {
			q.count--
		}
	}
}

// All returns all nodes.
func (q *QuadTree[T, N]) All() []N {
	all := make([]N, 0, q.count)
	all = append(all, q.outside...)
	if q.root != nil {
		all = q.root.all(all)
	}
	return all
}

// Reorganize the QuadTree to optimally fit its contents.
func (q *QuadTree[T, N]) Reorganize() {
	all := q.All()
	var rect geom.Rect[T]
	for _, one := range all {
		rect = rect.Union(one.Bounds())
	}
	q.root = nil
	q.outside = nil
	if len(all) > 0 {
		q.root = &node[T, N]{
			rect:      rect,
			threshold: q.threshold(),
		}
		for _, one := range all {
			q.root.insert(one)
		}
	}
}

// Clear removes all nodes.
func (q *QuadTree[T, N]) Clear() {
	q.count = 0
	q.root = nil
	q.outside = nil
}

// ContainsPoint returns true if at least one node contains the point.
func (q *QuadTree[T, N]) ContainsPoint(pt geom.Point[T]) bool {
	if q.root != nil {
		if q.root.containsPoint(pt) {
			return true
		}
	}
	for _, one := range q.outside {
		if pt.In(one.Bounds()) {
			return true
		}
	}
	return false
}

// FindContainsPoint returns the nodes that contain the point.
func (q *QuadTree[T, N]) FindContainsPoint(pt geom.Point[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findContainsPoint(pt, result)
	}
	for _, one := range q.outside {
		if pt.In(one.Bounds()) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainsPoint returns true if at least one node that the matcher returns true for contains the point.
func (q *QuadTree[T, N]) MatchedContainsPoint(matcher Matcher[T, N], pt geom.Point[T]) bool {
	if q.root != nil {
		if q.root.matchedContainsPoint(matcher, pt) {
			return true
		}
	}
	for _, one := range q.outside {
		if pt.In(one.Bounds()) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainsPoint returns the nodes that the matcher returns true for which contain the point.
func (q *QuadTree[T, N]) FindMatchedContainsPoint(matcher Matcher[T, N], pt geom.Point[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findMatchedContainsPoint(matcher, pt, result)
	}
	for _, one := range q.outside {
		if pt.In(one.Bounds()) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}

// Intersects returns true if at least one node intersects the rect.
func (q *QuadTree[T, N]) Intersects(rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.intersects(rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().Intersects(rect) {
			return true
		}
	}
	return false
}

// FindIntersects returns the nodes that intersect the rect.
func (q *QuadTree[T, N]) FindIntersects(rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findIntersects(rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().Intersects(rect) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedIntersects returns true if at least one node that the matcher returns true for intersects the rect.
func (q *QuadTree[T, N]) MatchedIntersects(matcher Matcher[T, N], rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.matchedIntersects(matcher, rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().Intersects(rect) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedIntersects returns the nodes that the matcher returns true for which intersect the rect.
func (q *QuadTree[T, N]) FindMatchedIntersects(matcher Matcher[T, N], rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findMatchedIntersects(matcher, rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().Intersects(rect) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}

// ContainsRect returns true if at least one node contains the rect.
func (q *QuadTree[T, N]) ContainsRect(rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.containsRect(rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().Contains(rect) {
			return true
		}
	}
	return false
}

// FindContainsRect returns the nodes that contain the rect.
func (q *QuadTree[T, N]) FindContainsRect(rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findContainsRect(rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().Contains(rect) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainsRect returns true if at least one node that the matcher returns true for contains the rect.
func (q *QuadTree[T, N]) MatchedContainsRect(matcher Matcher[T, N], rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.matchedContainsRect(matcher, rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().Contains(rect) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainsRect returns the nodes that the matcher returns true for which contains the rect.
func (q *QuadTree[T, N]) FindMatchedContainsRect(matcher Matcher[T, N], rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findMatchedContainsRect(matcher, rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().Contains(rect) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}

// ContainedByRect returns true if at least one node is contained by the rect.
func (q *QuadTree[T, N]) ContainedByRect(rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.containedByRect(rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if rect.Contains(one.Bounds()) {
			return true
		}
	}
	return false
}

// FindContainedByRect returns the nodes that are contained by the rect.
func (q *QuadTree[T, N]) FindContainedByRect(rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findContainedByRect(rect, result)
	}
	for _, one := range q.outside {
		if rect.Contains(one.Bounds()) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainedByRect returns true if at least one node that the matcher returns true for is contained by the rect.
func (q *QuadTree[T, N]) MatchedContainedByRect(matcher Matcher[T, N], rect geom.Rect[T]) bool {
	if q.root != nil {
		if q.root.matchedContainedByRect(matcher, rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if rect.Contains(one.Bounds()) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainedByRect returns the nodes that the matcher returns true for which are contained by the rect.
func (q *QuadTree[T, N]) FindMatchedContainedByRect(matcher Matcher[T, N], rect geom.Rect[T]) []N {
	var result []N
	if q.root != nil {
		result = q.root.findMatchedContainedByRect(matcher, rect, result)
	}
	for _, one := range q.outside {
		if rect.Contains(one.Bounds()) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}
