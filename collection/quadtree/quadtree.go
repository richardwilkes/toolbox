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

const (
	// DefaultQuadTreeThreshold is the threshold used when QuadTree.Threshold is below MinQuadTreeThreshold.
	DefaultQuadTreeThreshold = 64
	// MinQuadTreeThreshold is the minimum allowed threshold.
	MinQuadTreeThreshold = 4
)

// Node defines the methods an object stored in a QuadTree must implement.
type Node interface {
	// Bounds returns the node's bounding rectangle.
	Bounds() geom.Rect
	comparable
}

// Matcher is used to match nodes.
type Matcher[N Node] interface {
	// Matches returns true if the node matches.
	Matches(n N) bool
}

// QuadTree stores two-dimensional nodes for fast lookup.
//
// Each node is stored in the deepest cell that fully contains its Bounds(). A node whose bounds straddle a cell's
// center line cannot descend into any of its four children, so it stays in that cell. Input that concentrates many
// overlapping or center-crossing nodes in one region, or nodes that are large relative to the tree's bounds, therefore
// piles up in a single cell that cannot subdivide further, degrading queries over that region toward O(n). This is
// inherent to storing whole rectangles rather than points: keep node bounds small relative to the tree and spread out
// to preserve the expected logarithmic behavior.
type QuadTree[N Node] struct {
	root    *node[N]
	outside []N
	// Threshold is the number of nodes a cell may hold before it subdivides. Values below MinQuadTreeThreshold are
	// treated as DefaultQuadTreeThreshold.
	Threshold int
	count     int
}

// Size returns the number of nodes contained within the QuadTree.
func (q *QuadTree[N]) Size() int {
	return q.count
}

func (q *QuadTree[N]) threshold() int {
	if q.Threshold < MinQuadTreeThreshold {
		return DefaultQuadTreeThreshold
	}
	return q.Threshold
}

// Insert a node. Its Bounds() must not change until it is removed. Nodes with empty bounds are ignored.
func (q *QuadTree[N]) Insert(n N) {
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
func (q *QuadTree[N]) Remove(n N) {
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
func (q *QuadTree[N]) All() []N {
	all := make([]N, 0, q.count)
	all = append(all, q.outside...)
	if q.root != nil {
		all = q.root.all(all)
	}
	return all
}

// Reorganize rebuilds the QuadTree so that its root bounds fit its current contents.
func (q *QuadTree[N]) Reorganize() {
	all := q.All()
	var rect geom.Rect
	for _, one := range all {
		rect = rect.Union(one.Bounds())
	}
	q.root = nil
	q.outside = nil
	if len(all) > 0 {
		q.root = &node[N]{
			rect:      rect,
			threshold: q.threshold(),
		}
		for _, one := range all {
			q.root.insert(one)
		}
	}
}

// Clear removes all nodes.
func (q *QuadTree[N]) Clear() {
	q.count = 0
	q.root = nil
	q.outside = nil
}

// ContainsPoint returns true if at least one node contains the point.
func (q *QuadTree[N]) ContainsPoint(pt geom.Point) bool {
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
func (q *QuadTree[N]) FindContainsPoint(pt geom.Point) []N {
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

// MatchedContainsPoint returns true if at least one node accepted by the matcher contains the point.
func (q *QuadTree[N]) MatchedContainsPoint(matcher Matcher[N], pt geom.Point) bool {
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

// FindMatchedContainsPoint returns the nodes accepted by the matcher that contain the point.
func (q *QuadTree[N]) FindMatchedContainsPoint(matcher Matcher[N], pt geom.Point) []N {
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
func (q *QuadTree[N]) Intersects(rect geom.Rect) bool {
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
func (q *QuadTree[N]) FindIntersects(rect geom.Rect) []N {
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

// MatchedIntersects returns true if at least one node accepted by the matcher intersects the rect.
func (q *QuadTree[N]) MatchedIntersects(matcher Matcher[N], rect geom.Rect) bool {
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

// FindMatchedIntersects returns the nodes accepted by the matcher that intersect the rect.
func (q *QuadTree[N]) FindMatchedIntersects(matcher Matcher[N], rect geom.Rect) []N {
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
func (q *QuadTree[N]) ContainsRect(rect geom.Rect) bool {
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
func (q *QuadTree[N]) FindContainsRect(rect geom.Rect) []N {
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

// MatchedContainsRect returns true if at least one node accepted by the matcher contains the rect.
func (q *QuadTree[N]) MatchedContainsRect(matcher Matcher[N], rect geom.Rect) bool {
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

// FindMatchedContainsRect returns the nodes accepted by the matcher that contain the rect.
func (q *QuadTree[N]) FindMatchedContainsRect(matcher Matcher[N], rect geom.Rect) []N {
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
func (q *QuadTree[N]) ContainedByRect(rect geom.Rect) bool {
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
func (q *QuadTree[N]) FindContainedByRect(rect geom.Rect) []N {
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

// MatchedContainedByRect returns true if at least one node accepted by the matcher is contained by the rect.
func (q *QuadTree[N]) MatchedContainedByRect(matcher Matcher[N], rect geom.Rect) bool {
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

// FindMatchedContainedByRect returns the nodes accepted by the matcher that are contained by the rect.
func (q *QuadTree[N]) FindMatchedContainedByRect(matcher Matcher[N], rect geom.Rect) []N {
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
