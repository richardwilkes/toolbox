// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree_test

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/collection/quadtree"
	"github.com/richardwilkes/toolbox/v2/xmath/geom"
	"golang.org/x/exp/constraints"
)

type node[T constraints.Float] struct {
	geom.Rect[T]
}

func newNode[T constraints.Float](x, y, width, height T) *node[T] {
	return &node[T]{Rect: geom.NewRect(x, y, width, height)}
}

func (n node[T]) Bounds() geom.Rect[T] {
	return n.Rect
}

func TestContainsPoint(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	c.False(q.ContainsPoint(geom.Point[float64]{}))
	q.Insert(newNode[float64](5, 5, 5, 5))
	c.False(q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	c.True(q.ContainsPoint(geom.NewPoint[float64](5, 5)))
	c.True(q.ContainsPoint(geom.NewPoint(9.9, 9.9)))
	c.False(q.ContainsPoint(geom.NewPoint[float64](10, 10)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	c.True(q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	c.True(q.ContainsPoint(geom.Point[float64]{}))
	c.True(q.ContainsPoint(geom.NewPoint[float64](0, -5)))
	c.False(q.ContainsPoint(geom.NewPoint[float64](-1, 0)))
}

func TestContainsRect(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	c.False(q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	q.Insert(newNode[float64](5, 5, 5, 5))
	c.False(q.ContainsRect(geom.NewRect[float64](4, 4, 10, 10)))
	c.True(q.ContainsRect(geom.NewRect[float64](5, 5, 2, 2)))
	c.True(q.ContainsRect(geom.NewRect(9.9, 9.9, .05, .05)))
	c.False(q.ContainsRect(geom.NewRect[float64](10, 10, 5, 5)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	c.True(q.ContainsRect(geom.NewRect[float64](6, 4, 1, 2)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	c.True(q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	c.True(q.ContainsRect(geom.NewRect[float64](0, -5, 4, 4)))
	c.False(q.ContainsRect(geom.NewRect[float64](-1, 0, 2, 2)))
}

func TestGeneral(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	r := rand.New(rand.NewPCG(22, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	mine := newNode[float64](22, 22, 22, 22)
	q.Insert(mine)
	for range 100 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(50000-r.IntN(100000)), float64(50000-r.IntN(100000)), float64(r.IntN(100000)), float64(r.IntN(100000))))
	}
	c.Equal(1+100*quadtree.DefaultQuadTreeThreshold, q.Size())
	all := q.All()
	c.True(slices.Contains(all, mine))
	count := q.Size()
	for _, one := range all {
		if one != mine && r.IntN(10) == 1 {
			q.Remove(one)
			count--
			c.Equal(count, q.Size())
		}
	}
	c.Equal(count, q.Size())
	q.Reorganize()
	c.Equal(count, q.Size())
	c.True(slices.Contains(q.All(), mine))
	c.True(slices.Contains(q.FindContainedByRect(mine.Rect), mine))
}

func TestEmptyQuadTree(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Test empty quadtree behavior
	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
	c.False(q.ContainsPoint(geom.Point[float64]{}))
	c.False(q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	c.False(q.Intersects(geom.NewRect[float64](0, 0, 1, 1)))
	c.False(q.ContainedByRect(geom.NewRect[float64](0, 0, 100, 100)))

	// Test find methods on empty tree
	c.Equal(0, len(q.FindContainsPoint(geom.Point[float64]{})))
	c.Equal(0, len(q.FindContainsRect(geom.NewRect[float64](0, 0, 1, 1))))
	c.Equal(0, len(q.FindIntersects(geom.NewRect[float64](0, 0, 1, 1))))
	c.Equal(0, len(q.FindContainedByRect(geom.NewRect[float64](0, 0, 100, 100))))
}

func TestInsertEmptyNode(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Insert node with empty bounds
	emptyNode := newNode[float64](0, 0, 0, 0)
	q.Insert(emptyNode)
	c.Equal(0, q.Size()) // Empty nodes should not be inserted
}

func TestRemove(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	n1 := newNode[float64](0, 0, 10, 10)
	n2 := newNode[float64](20, 20, 10, 10)
	n3 := newNode[float64](40, 40, 10, 10)

	q.Insert(n1)
	q.Insert(n2)
	q.Insert(n3)
	c.Equal(3, q.Size())

	// Remove existing node
	q.Remove(n2)
	c.Equal(2, q.Size())
	c.False(slices.Contains(q.All(), n2))
	c.True(slices.Contains(q.All(), n1))
	c.True(slices.Contains(q.All(), n3))

	// Remove non-existing node (should not change size)
	nonExisting := newNode[float64](100, 100, 10, 10)
	q.Remove(nonExisting)
	c.Equal(2, q.Size())

	// Remove remaining nodes
	q.Remove(n1)
	q.Remove(n3)
	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
}

func TestClear(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Add some nodes
	for i := range 10 {
		q.Insert(newNode(float64(i*10), float64(i*10), 5, 5))
	}
	c.Equal(10, q.Size())

	// Clear and verify
	q.Clear()
	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
	c.False(q.ContainsPoint(geom.NewPoint[float64](25, 25)))
}

func TestThreshold(t *testing.T) {
	c := check.New(t)

	// Test default threshold (Threshold field starts at 0, internal threshold() method returns default)
	q1 := &quadtree.QuadTree[float64, *node[float64]]{}
	c.Equal(0, q1.Threshold) // Field is 0 by default

	// Test custom threshold
	q2 := &quadtree.QuadTree[float64, *node[float64]]{Threshold: 10}
	c.Equal(10, q2.Threshold)

	// Test threshold below minimum (should use default internally)
	q3 := &quadtree.QuadTree[float64, *node[float64]]{Threshold: 2}
	c.Equal(2, q3.Threshold) // Field value is preserved

	// Test that internal threshold logic works by triggering reorganization
	// Add enough nodes to trigger reorganization behavior
	for i := range quadtree.DefaultQuadTreeThreshold * 2 {
		q3.Insert(newNode(float64(i), 0, 1, 1))
	}
	// If internal threshold is working correctly, the tree should handle this without issues
	c.True(q3.Size() > 0)
}

func TestIntersects(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	n1 := newNode[float64](0, 0, 10, 10)
	n2 := newNode[float64](15, 15, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	// Test intersections
	c.True(q.Intersects(geom.NewRect[float64](5, 5, 10, 10)))      // Overlaps with n1
	c.True(q.Intersects(geom.NewRect[float64](20, 20, 10, 10)))    // Overlaps with n2
	c.True(q.Intersects(geom.NewRect[float64](0, 0, 30, 30)))      // Overlaps both
	c.False(q.Intersects(geom.NewRect[float64](100, 100, 10, 10))) // No overlap

	// Test FindIntersects
	intersects1 := q.FindIntersects(geom.NewRect[float64](5, 5, 10, 10))
	c.True(slices.Contains(intersects1, n1))

	intersects2 := q.FindIntersects(geom.NewRect[float64](0, 0, 30, 30))
	c.True(slices.Contains(intersects2, n1))
	c.True(slices.Contains(intersects2, n2))

	intersects3 := q.FindIntersects(geom.NewRect[float64](100, 100, 10, 10))
	c.Equal(0, len(intersects3))
}

func TestContainedByRect(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	n1 := newNode[float64](5, 5, 10, 10)   // Inside container
	n2 := newNode[float64](50, 50, 10, 10) // Outside container
	q.Insert(n1)
	q.Insert(n2)

	container := geom.NewRect[float64](0, 0, 20, 20)

	c.True(q.ContainedByRect(container))

	contained := q.FindContainedByRect(container)
	c.True(slices.Contains(contained, n1))
	c.False(slices.Contains(contained, n2))

	// Test with container that contains nothing
	smallContainer := geom.NewRect[float64](100, 100, 5, 5)
	c.False(q.ContainedByRect(smallContainer))
	c.Equal(0, len(q.FindContainedByRect(smallContainer)))
}

func TestReorganize(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Add nodes that will initially be outside
	nodes := make([]*node[float64], 0, 10)
	for i := range 10 {
		n := newNode(float64(i*100), float64(i*100), 10, 10)
		nodes = append(nodes, n)
		q.Insert(n)
	}

	initialSize := q.Size()
	c.Equal(10, initialSize)

	// Reorganize and verify all nodes are still present
	q.Reorganize()
	c.Equal(initialSize, q.Size())

	all := q.All()
	for _, n := range nodes {
		c.True(slices.Contains(all, n))
	}

	// Test reorganize with empty tree
	q.Clear()
	q.Reorganize()
	c.Equal(0, q.Size())
}

type testMatcher struct {
	target *node[float64]
}

func (m *testMatcher) Matches(n *node[float64]) bool {
	return n == m.target
}

func TestMatchedMethods(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	n1 := newNode[float64](0, 0, 10, 10)
	n2 := newNode[float64](20, 20, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	matcher := &testMatcher{target: n1}

	// Test MatchedContainsPoint
	c.True(q.MatchedContainsPoint(matcher, geom.NewPoint[float64](5, 5)))    // n1 contains point
	c.False(q.MatchedContainsPoint(matcher, geom.NewPoint[float64](25, 25))) // n2 contains point but doesn't match

	// Test FindMatchedContainsPoint
	matched := q.FindMatchedContainsPoint(matcher, geom.NewPoint[float64](5, 5))
	c.Equal(1, len(matched))
	c.Equal(n1, matched[0])

	// Test MatchedIntersects
	c.True(q.MatchedIntersects(matcher, geom.NewRect[float64](5, 5, 10, 10)))
	c.False(q.MatchedIntersects(matcher, geom.NewRect[float64](25, 25, 10, 10)))

	// Test FindMatchedIntersects
	matchedIntersects := q.FindMatchedIntersects(matcher, geom.NewRect[float64](0, 0, 30, 30))
	c.Equal(1, len(matchedIntersects))
	c.Equal(n1, matchedIntersects[0])

	// Test MatchedContainsRect
	c.True(q.MatchedContainsRect(matcher, geom.NewRect[float64](2, 2, 5, 5)))
	c.False(q.MatchedContainsRect(matcher, geom.NewRect[float64](22, 22, 5, 5)))

	// Test FindMatchedContainsRect
	matchedContains := q.FindMatchedContainsRect(matcher, geom.NewRect[float64](2, 2, 5, 5))
	c.Equal(1, len(matchedContains))
	c.Equal(n1, matchedContains[0])

	// Test MatchedContainedByRect
	c.True(q.MatchedContainedByRect(matcher, geom.NewRect[float64](0, 0, 50, 50)))
	c.False(q.MatchedContainedByRect(matcher, geom.NewRect[float64](5, 5, 3, 3)))

	// Test FindMatchedContainedByRect
	matchedContainedBy := q.FindMatchedContainedByRect(matcher, geom.NewRect[float64](0, 0, 50, 50))
	c.Equal(1, len(matchedContainedBy))
	c.Equal(n1, matchedContainedBy[0])
}

func TestLargeDataset(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Insert many nodes to test tree subdivision
	nodeCount := quadtree.DefaultQuadTreeThreshold * 10
	nodes := make([]*node[float64], 0, nodeCount)

	r := rand.New(rand.NewPCG(42, 2023)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for range nodeCount {
		x := float64(r.IntN(1000))
		y := float64(r.IntN(1000))
		w := float64(r.IntN(50) + 1)
		h := float64(r.IntN(50) + 1)
		n := newNode(x, y, w, h)
		nodes = append(nodes, n)
		q.Insert(n)
	}

	c.Equal(nodeCount, q.Size())

	// Test that all nodes are findable
	all := q.All()
	c.Equal(nodeCount, len(all))
	for _, n := range nodes {
		c.True(slices.Contains(all, n))
	}

	// Test spatial queries work correctly
	queryRect := geom.NewRect[float64](100, 100, 200, 200)
	intersecting := q.FindIntersects(queryRect)

	// Verify results by checking manually
	for _, n := range intersecting {
		c.True(n.Bounds().Intersects(queryRect))
	}

	// Verify no false negatives
	for _, n := range nodes {
		if n.Bounds().Intersects(queryRect) {
			c.True(slices.Contains(intersecting, n))
		}
	}
}

func TestEdgeCases(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float64, *node[float64]]{}

	// Test with nodes at exact boundaries
	n1 := newNode[float64](0, 0, 10, 10)
	n2 := newNode[float64](10, 10, 10, 10) // Touches corner of n1
	q.Insert(n1)
	q.Insert(n2)

	// Test point queries at boundaries
	c.True(q.ContainsPoint(geom.NewPoint[float64](0, 0)))    // Corner of n1
	c.True(q.ContainsPoint(geom.NewPoint[float64](10, 10)))  // Corner of both
	c.False(q.ContainsPoint(geom.NewPoint[float64](20, 20))) // Outside both

	// Test with very small nodes
	tiny := newNode(100, 100, 0.001, 0.001)
	q.Insert(tiny)
	c.True(q.ContainsPoint(geom.NewPoint[float64](100, 100)))
	c.False(q.ContainsPoint(geom.NewPoint(100.1, 100.1)))

	// Test with very large nodes
	huge := newNode[float64](-1000, -1000, 2000, 2000)
	q.Insert(huge)
	c.True(q.ContainsPoint(geom.NewPoint[float64](500, 500)))
	c.True(q.ContainsRect(geom.NewRect[float64](0, 0, 100, 100)))
}

func TestFloat32(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[float32, *node[float32]]{}

	n1 := newNode[float32](0, 0, 10, 10)
	n2 := newNode[float32](15, 15, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	c.Equal(2, q.Size())
	c.True(q.ContainsPoint(geom.NewPoint[float32](5, 5)))
	c.True(q.Intersects(geom.NewRect[float32](0, 0, 30, 30)))

	all := q.All()
	c.Equal(2, len(all))
	c.True(slices.Contains(all, n1))
	c.True(slices.Contains(all, n2))
}
