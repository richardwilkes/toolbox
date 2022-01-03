// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree

import "github.com/richardwilkes/toolbox/xmath/geom"

const (
	// DefaultQuadTreeThreshold is the default threshold that will be used if none is specified.
	DefaultQuadTreeThreshold = 64
	// MinQuadTreeThreshold is the minimum allowed threshold.
	MinQuadTreeThreshold = 4
)

// Node defines the methods an object that can be stored within the QuadTree must implement.
type Node interface {
	// Bounds returns the node's bounding rectangle.
	Bounds() geom.Rect
}

// Matcher is used to match nodes.
type Matcher interface {
	// Matches returns true if the node matches.
	Matches(n Node) bool
}

// QuadTree stores two-dimensional nodes for fast lookup.
type QuadTree struct {
	Threshold int
	count     int
	root      *node
	outside   []Node
}

// Size returns the number of nodes contained within the QuadTree.
func (q *QuadTree) Size() int {
	return q.count
}

func (q *QuadTree) threshold() int {
	if q.Threshold < MinQuadTreeThreshold {
		return DefaultQuadTreeThreshold
	}
	return q.Threshold
}

// Insert a node. NOTE: Once a node is inserted, the value it returns from a call to Bounds() MUST REMAIN THE SAME until
// the node is removed.
func (q *QuadTree) Insert(n Node) {
	rect := n.Bounds()
	if rect.IsEmpty() {
		return
	}
	q.count++
	if q.root != nil && q.root.rect.ContainsRect(rect) {
		q.root.insert(n)
	} else {
		q.outside = append(q.outside, n)
		if len(q.outside) > q.threshold() {
			q.Reorganize()
		}
	}
}

// Remove a node.
func (q *QuadTree) Remove(n Node) {
	for i, one := range q.outside {
		if one != n {
			continue
		}
		q.outside[i] = q.outside[len(q.outside)-1]
		q.outside[len(q.outside)-1] = nil
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
func (q *QuadTree) All() []Node {
	all := make([]Node, 0, q.count)
	all = append(all, q.outside...)
	if q.root != nil {
		all = q.root.all(all)
	}
	return all
}

// Reorganize the QuadTree to optimally fit its contents.
func (q *QuadTree) Reorganize() {
	all := q.All()
	var rect geom.Rect
	for _, one := range all {
		rect.Union(one.Bounds())
	}
	q.root = nil
	q.outside = make([]Node, 0)
	if len(all) > 0 {
		q.root = &node{
			rect:      rect,
			threshold: q.threshold(),
		}
		for _, one := range all {
			q.root.insert(one)
		}
	}
}

// Clear removes all nodes.
func (q *QuadTree) Clear() {
	q.count = 0
	q.root = nil
	q.outside = make([]Node, 0)
}

// ContainsPoint returns true if at least one node contains the point.
func (q *QuadTree) ContainsPoint(pt geom.Point) bool {
	if q.root != nil {
		if q.root.containsPoint(pt) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsPoint(pt) {
			return true
		}
	}
	return false
}

// FindContainsPoint returns the nodes that contain the point.
func (q *QuadTree) FindContainsPoint(pt geom.Point) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findContainsPoint(pt, result)
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsPoint(pt) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainsPoint returns true if at least one node that the matcher returns true for contains the point.
func (q *QuadTree) MatchedContainsPoint(matcher Matcher, pt geom.Point) bool {
	if q.root != nil {
		if q.root.matchedContainsPoint(matcher, pt) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsPoint(pt) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainsPoint returns the nodes that the matcher returns true for which contain the point.
func (q *QuadTree) FindMatchedContainsPoint(matcher Matcher, pt geom.Point) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findMatchedContainsPoint(matcher, pt, result)
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsPoint(pt) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}

// Intersects returns true if at least one node intersects the rect.
func (q *QuadTree) Intersects(rect geom.Rect) bool {
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
func (q *QuadTree) FindIntersects(rect geom.Rect) []Node {
	var result []Node
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
func (q *QuadTree) MatchedIntersects(matcher Matcher, rect geom.Rect) bool {
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
func (q *QuadTree) FindMatchedIntersects(matcher Matcher, rect geom.Rect) []Node {
	var result []Node
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
func (q *QuadTree) ContainsRect(rect geom.Rect) bool {
	if q.root != nil {
		if q.root.containsRect(rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsRect(rect) {
			return true
		}
	}
	return false
}

// FindContainsRect returns the nodes that contain the rect.
func (q *QuadTree) FindContainsRect(rect geom.Rect) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findContainsRect(rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsRect(rect) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainsRect returns true if at least one node that the matcher returns true for contains the rect.
func (q *QuadTree) MatchedContainsRect(matcher Matcher, rect geom.Rect) bool {
	if q.root != nil {
		if q.root.matchedContainsRect(matcher, rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsRect(rect) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainsRect returns the nodes that the matcher returns true for which contains the rect.
func (q *QuadTree) FindMatchedContainsRect(matcher Matcher, rect geom.Rect) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findMatchedContainsRect(matcher, rect, result)
	}
	for _, one := range q.outside {
		if one.Bounds().ContainsRect(rect) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}

// ContainedByRect returns true if at least one node is contained by the rect.
func (q *QuadTree) ContainedByRect(rect geom.Rect) bool {
	if q.root != nil {
		if q.root.containedByRect(rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if rect.ContainsRect(one.Bounds()) {
			return true
		}
	}
	return false
}

// FindContainedByRect returns the nodes that are contained by the rect.
func (q *QuadTree) FindContainedByRect(rect geom.Rect) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findContainedByRect(rect, result)
	}
	for _, one := range q.outside {
		if rect.ContainsRect(one.Bounds()) {
			result = append(result, one)
		}
	}
	return result
}

// MatchedContainedByRect returns true if at least one node that the matcher returns true for is contained by the rect.
func (q *QuadTree) MatchedContainedByRect(matcher Matcher, rect geom.Rect) bool {
	if q.root != nil {
		if q.root.matchedContainedByRect(matcher, rect) {
			return true
		}
	}
	for _, one := range q.outside {
		if rect.ContainsRect(one.Bounds()) && matcher.Matches(one) {
			return true
		}
	}
	return false
}

// FindMatchedContainedByRect returns the nodes that the matcher returns true for which are contained by the rect.
func (q *QuadTree) FindMatchedContainedByRect(matcher Matcher, rect geom.Rect) []Node {
	var result []Node
	if q.root != nil {
		result = q.root.findMatchedContainedByRect(matcher, rect, result)
	}
	for _, one := range q.outside {
		if rect.ContainsRect(one.Bounds()) && matcher.Matches(one) {
			result = append(result, one)
		}
	}
	return result
}
