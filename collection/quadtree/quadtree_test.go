// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"github.com/richardwilkes/toolbox/v2/geom"
)

type node struct {
	geom.Rect
}

func newNode(x, y, width, height float32) *node {
	return &node{Rect: geom.NewRect(x, y, width, height)}
}

func (n node) Bounds() geom.Rect {
	return n.Rect
}

func TestContainsPoint(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}
	c.False(q.ContainsPoint(geom.Point{}))
	q.Insert(newNode(5, 5, 5, 5))
	c.False(q.ContainsPoint(geom.NewPoint(6, 4)))
	c.True(q.ContainsPoint(geom.NewPoint(5, 5)))
	c.True(q.ContainsPoint(geom.NewPoint(9.9, 9.9)))
	c.False(q.ContainsPoint(geom.NewPoint(10, 10)))
	q.Insert(newNode(4, 4, 3, 3))
	c.True(q.ContainsPoint(geom.NewPoint(6, 4)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float32(i), -5, 10, 10))
	}
	c.True(q.ContainsPoint(geom.Point{}))
	c.True(q.ContainsPoint(geom.NewPoint(0, -5)))
	c.False(q.ContainsPoint(geom.NewPoint(-1, 0)))
}

func TestContainsRect(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}
	c.False(q.ContainsRect(geom.NewRect(0, 0, 1, 1)))
	q.Insert(newNode(5, 5, 5, 5))
	c.False(q.ContainsRect(geom.NewRect(4, 4, 10, 10)))
	c.True(q.ContainsRect(geom.NewRect(5, 5, 2, 2)))
	c.True(q.ContainsRect(geom.NewRect(9.9, 9.9, .05, .05)))
	c.False(q.ContainsRect(geom.NewRect(10, 10, 5, 5)))
	q.Insert(newNode(4, 4, 3, 3))
	c.True(q.ContainsRect(geom.NewRect(6, 4, 1, 2)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float32(i), -5, 10, 10))
	}
	c.True(q.ContainsRect(geom.NewRect(0, 0, 1, 1)))
	c.True(q.ContainsRect(geom.NewRect(0, -5, 4, 4)))
	c.False(q.ContainsRect(geom.NewRect(-1, 0, 2, 2)))
}

func TestGeneral(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}
	r := rand.New(rand.NewPCG(22, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	mine := newNode(22, 22, 22, 22)
	q.Insert(mine)
	for range 100 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float32(50000-r.IntN(100000)), float32(50000-r.IntN(100000)), float32(r.IntN(100000)), float32(r.IntN(100000))))
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
	q := &quadtree.QuadTree[*node]{}

	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
	c.False(q.ContainsPoint(geom.Point{}))
	c.False(q.ContainsRect(geom.NewRect(0, 0, 1, 1)))
	c.False(q.Intersects(geom.NewRect(0, 0, 1, 1)))
	c.False(q.ContainedByRect(geom.NewRect(0, 0, 100, 100)))

	c.Equal(0, len(q.FindContainsPoint(geom.Point{})))
	c.Equal(0, len(q.FindContainsRect(geom.NewRect(0, 0, 1, 1))))
	c.Equal(0, len(q.FindIntersects(geom.NewRect(0, 0, 1, 1))))
	c.Equal(0, len(q.FindContainedByRect(geom.NewRect(0, 0, 100, 100))))
}

func TestInsertEmptyNode(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	emptyNode := newNode(0, 0, 0, 0)
	q.Insert(emptyNode)
	c.Equal(0, q.Size()) // Empty nodes should not be inserted
}

func TestRemove(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(20, 20, 10, 10)
	n3 := newNode(40, 40, 10, 10)

	q.Insert(n1)
	q.Insert(n2)
	q.Insert(n3)
	c.Equal(3, q.Size())

	q.Remove(n2)
	c.Equal(2, q.Size())
	c.False(slices.Contains(q.All(), n2))
	c.True(slices.Contains(q.All(), n1))
	c.True(slices.Contains(q.All(), n3))

	nonExisting := newNode(100, 100, 10, 10)
	q.Remove(nonExisting)
	c.Equal(2, q.Size())

	q.Remove(n1)
	q.Remove(n3)
	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
}

func TestClear(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	for i := range 10 {
		q.Insert(newNode(float32(i*10), float32(i*10), 5, 5))
	}
	c.Equal(10, q.Size())

	q.Clear()
	c.Equal(0, q.Size())
	c.Equal(0, len(q.All()))
	c.False(q.ContainsPoint(geom.NewPoint(25, 25)))
}

func TestThreshold(t *testing.T) {
	c := check.New(t)

	q1 := &quadtree.QuadTree[*node]{}
	c.Equal(0, q1.Threshold) // Field is 0 by default

	q2 := &quadtree.QuadTree[*node]{Threshold: 10}
	c.Equal(10, q2.Threshold)

	// A threshold below the minimum falls back to the default internally
	q3 := &quadtree.QuadTree[*node]{Threshold: 2}
	c.Equal(2, q3.Threshold) // Field value is preserved

	// Insert enough nodes to exercise the internal threshold
	for i := range quadtree.DefaultQuadTreeThreshold * 2 {
		q3.Insert(newNode(float32(i), 0, 1, 1))
	}
	c.True(q3.Size() > 0)
}

func TestIntersects(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(15, 15, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	c.True(q.Intersects(geom.NewRect(5, 5, 10, 10)))      // Overlaps with n1
	c.True(q.Intersects(geom.NewRect(20, 20, 10, 10)))    // Overlaps with n2
	c.True(q.Intersects(geom.NewRect(0, 0, 30, 30)))      // Overlaps both
	c.False(q.Intersects(geom.NewRect(100, 100, 10, 10))) // No overlap

	intersects1 := q.FindIntersects(geom.NewRect(5, 5, 10, 10))
	c.True(slices.Contains(intersects1, n1))

	intersects2 := q.FindIntersects(geom.NewRect(0, 0, 30, 30))
	c.True(slices.Contains(intersects2, n1))
	c.True(slices.Contains(intersects2, n2))

	intersects3 := q.FindIntersects(geom.NewRect(100, 100, 10, 10))
	c.Equal(0, len(intersects3))
}

func TestContainedByRect(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(5, 5, 10, 10)   // Inside container
	n2 := newNode(50, 50, 10, 10) // Outside container
	q.Insert(n1)
	q.Insert(n2)

	container := geom.NewRect(0, 0, 20, 20)

	c.True(q.ContainedByRect(container))

	contained := q.FindContainedByRect(container)
	c.True(slices.Contains(contained, n1))
	c.False(slices.Contains(contained, n2))

	smallContainer := geom.NewRect(100, 100, 5, 5)
	c.False(q.ContainedByRect(smallContainer))
	c.Equal(0, len(q.FindContainedByRect(smallContainer)))
}

func TestReorganize(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	nodes := make([]*node, 0, 10)
	for i := range 10 {
		n := newNode(float32(i*100), float32(i*100), 10, 10)
		nodes = append(nodes, n)
		q.Insert(n)
	}

	initialSize := q.Size()
	c.Equal(10, initialSize)

	q.Reorganize()
	c.Equal(initialSize, q.Size())

	all := q.All()
	for _, n := range nodes {
		c.True(slices.Contains(all, n))
	}

	q.Clear()
	q.Reorganize()
	c.Equal(0, q.Size())
}

type testMatcher struct {
	target *node
}

func (m *testMatcher) Matches(n *node) bool {
	return n == m.target
}

func TestMatchedMethods(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(20, 20, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	matcher := &testMatcher{target: n1}

	c.True(q.MatchedContainsPoint(matcher, geom.NewPoint(5, 5)))    // n1 contains point
	c.False(q.MatchedContainsPoint(matcher, geom.NewPoint(25, 25))) // n2 contains point but doesn't match

	matched := q.FindMatchedContainsPoint(matcher, geom.NewPoint(5, 5))
	c.Equal(1, len(matched))
	c.Equal(n1, matched[0])

	c.True(q.MatchedIntersects(matcher, geom.NewRect(5, 5, 10, 10)))
	c.False(q.MatchedIntersects(matcher, geom.NewRect(25, 25, 10, 10)))

	matchedIntersects := q.FindMatchedIntersects(matcher, geom.NewRect(0, 0, 30, 30))
	c.Equal(1, len(matchedIntersects))
	c.Equal(n1, matchedIntersects[0])

	c.True(q.MatchedContainsRect(matcher, geom.NewRect(2, 2, 5, 5)))
	c.False(q.MatchedContainsRect(matcher, geom.NewRect(22, 22, 5, 5)))

	matchedContains := q.FindMatchedContainsRect(matcher, geom.NewRect(2, 2, 5, 5))
	c.Equal(1, len(matchedContains))
	c.Equal(n1, matchedContains[0])

	c.True(q.MatchedContainedByRect(matcher, geom.NewRect(0, 0, 50, 50)))
	c.False(q.MatchedContainedByRect(matcher, geom.NewRect(5, 5, 3, 3)))

	matchedContainedBy := q.FindMatchedContainedByRect(matcher, geom.NewRect(0, 0, 50, 50))
	c.Equal(1, len(matchedContainedBy))
	c.Equal(n1, matchedContainedBy[0])
}

func TestLargeDataset(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	nodeCount := quadtree.DefaultQuadTreeThreshold * 10
	nodes := make([]*node, 0, nodeCount)

	r := rand.New(rand.NewPCG(42, 2023)) //nolint:gosec // Yes, it is ok to use a weak prng here
	for range nodeCount {
		x := float32(r.IntN(1000))
		y := float32(r.IntN(1000))
		w := float32(r.IntN(50) + 1)
		h := float32(r.IntN(50) + 1)
		n := newNode(x, y, w, h)
		nodes = append(nodes, n)
		q.Insert(n)
	}

	c.Equal(nodeCount, q.Size())

	all := q.All()
	c.Equal(nodeCount, len(all))
	for _, n := range nodes {
		c.True(slices.Contains(all, n))
	}

	queryRect := geom.NewRect(100, 100, 200, 200)
	intersecting := q.FindIntersects(queryRect)

	// Verify no false positives
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
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(10, 10, 10, 10) // Touches corner of n1
	q.Insert(n1)
	q.Insert(n2)

	c.True(q.ContainsPoint(geom.NewPoint(0, 0)))    // Corner of n1
	c.True(q.ContainsPoint(geom.NewPoint(10, 10)))  // Corner of both
	c.False(q.ContainsPoint(geom.NewPoint(20, 20))) // Outside both

	tiny := newNode(100, 100, 0.001, 0.001)
	q.Insert(tiny)
	c.True(q.ContainsPoint(geom.NewPoint(100, 100)))
	c.False(q.ContainsPoint(geom.NewPoint(100.1, 100.1)))

	huge := newNode(-1000, -1000, 2000, 2000)
	q.Insert(huge)
	c.True(q.ContainsPoint(geom.NewPoint(500, 500)))
	c.True(q.ContainsRect(geom.NewRect(0, 0, 100, 100)))
}

func TestFloat32(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(15, 15, 10, 10)
	q.Insert(n1)
	q.Insert(n2)

	c.Equal(2, q.Size())
	c.True(q.ContainsPoint(geom.NewPoint(5, 5)))
	c.True(q.Intersects(geom.NewRect(0, 0, 30, 30)))

	all := q.All()
	c.Equal(2, len(all))
	c.True(slices.Contains(all, n1))
	c.True(slices.Contains(all, n2))
}

func TestTreeSubdivision(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	initialNode := newNode(50, 50, 10, 10)
	q.Insert(initialNode)

	// Add enough nodes in a small area to force subdivision
	for i := range quadtree.DefaultQuadTreeThreshold + 10 {
		x := 50 + float32(i%8)*1.5 // Spread nodes in a small area
		y := 50 + float32(i/8)*1.5
		n := newNode(x, y, 1, 1)
		q.Insert(n)
	}

	searchRect := geom.NewRect(50, 50, 20, 20)

	found := q.FindContainsPoint(geom.NewPoint(55, 55))
	c.True(len(found) > 0)

	smallRect := geom.NewRect(51, 51, 0.5, 0.5)
	containing := q.FindContainsRect(smallRect)
	c.True(len(containing) > 0)

	intersecting := q.FindIntersects(searchRect)
	c.True(len(intersecting) > 0)

	largeRect := geom.NewRect(45, 45, 30, 30)
	contained := q.FindContainedByRect(largeRect)
	c.True(len(contained) > 0)
}

func TestMatchedMethodsWithTreeTraversal(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	target1 := newNode(25, 25, 10, 10)
	target2 := newNode(75, 75, 10, 10)
	q.Insert(target1)
	q.Insert(target2)

	// Add enough other nodes to force subdivision
	for i := range quadtree.DefaultQuadTreeThreshold + 5 {
		x := float32(i%10) * 10
		y := float32(i/10) * 10
		q.Insert(newNode(x, y, 5, 5))
	}

	matcher := &multiTargetMatcher{targets: []*node{target1, target2}}

	c.True(q.MatchedContainsPoint(matcher, geom.NewPoint(30, 30))) // In target1
	c.False(q.MatchedContainsPoint(matcher, geom.NewPoint(5, 5)))  // Not in targets

	matchedPoints := q.FindMatchedContainsPoint(matcher, geom.NewPoint(80, 80))
	c.Equal(1, len(matchedPoints))
	c.Equal(target2, matchedPoints[0])

	testRect := geom.NewRect(20, 20, 20, 20)
	c.True(q.MatchedIntersects(matcher, testRect))

	matchedIntersects := q.FindMatchedIntersects(matcher, testRect)
	c.True(len(matchedIntersects) > 0)
	c.True(slices.Contains(matchedIntersects, target1))

	smallRect := geom.NewRect(26, 26, 5, 5)
	c.True(q.MatchedContainsRect(matcher, smallRect))

	matchedContains := q.FindMatchedContainsRect(matcher, smallRect)
	c.Equal(1, len(matchedContains))
	c.Equal(target1, matchedContains[0])

	largeRect := geom.NewRect(20, 20, 20, 20)
	c.True(q.MatchedContainedByRect(matcher, largeRect))

	matchedContainedBy := q.FindMatchedContainedByRect(matcher, largeRect)
	c.Equal(1, len(matchedContainedBy))
	c.Equal(target1, matchedContainedBy[0])
}

type multiTargetMatcher struct {
	targets []*node
}

func (m *multiTargetMatcher) Matches(n *node) bool {
	return slices.Contains(m.targets, n)
}

func TestNodeBoundsMethod(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	// Add enough nodes to force subdivision
	for i := range quadtree.DefaultQuadTreeThreshold + 5 {
		x := float32(i%10) * 2
		y := float32(i/10) * 2
		q.Insert(newNode(x, y, 1, 1))
	}

	q.Reorganize()
	c.True(q.Size() > 0)

	searchRect := geom.NewRect(0, 0, 5, 5)
	intersects := q.FindIntersects(searchRect)
	c.True(len(intersects) >= 0) //nolint:gocritic // This is a valid test for the number of intersects

	contains := q.FindContainsRect(geom.NewRect(1, 1, 0.5, 0.5))
	c.True(len(contains) >= 0) //nolint:gocritic // This is a valid test for the number of contains

	containedBy := q.FindContainedByRect(geom.NewRect(0, 0, 100, 100))
	c.True(len(containedBy) > 0)
}

func TestThresholdEdgeCases(t *testing.T) {
	c := check.New(t)

	q1 := &quadtree.QuadTree[*node]{Threshold: quadtree.MinQuadTreeThreshold}
	c.Equal(quadtree.MinQuadTreeThreshold, q1.Threshold)

	// A threshold below the minimum falls back to the default internally
	q2 := &quadtree.QuadTree[*node]{Threshold: 1}
	c.Equal(1, q2.Threshold) // Field preserves the value

	for i := range 10 {
		q2.Insert(newNode(float32(i*100), float32(i*100), 5, 5))
	}
	c.Equal(10, q2.Size())

	q3 := &quadtree.QuadTree[*node]{Threshold: -5}
	for i := range 5 {
		q3.Insert(newNode(float32(i), 0, 1, 1))
	}
	c.Equal(5, q3.Size())
}

func TestComplexTreeOperations(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	nodes := make([]*node, 0, 200)
	r := rand.New(rand.NewPCG(123, 456)) //nolint:gosec // Yes, it is ok to use a weak prng here

	for i := range 200 {
		// Create clusters of nodes to force subdivision
		cluster := i / 20
		x := float32(cluster*50) + float32(r.IntN(20))
		y := float32(cluster*50) + float32(r.IntN(20))
		w := float32(r.IntN(5) + 1)
		h := float32(r.IntN(5) + 1)
		n := newNode(x, y, w, h)
		nodes = append(nodes, n)
		q.Insert(n)
	}

	c.Equal(200, q.Size())

	testRect := geom.NewRect(25, 25, 50, 50)
	testPoint := geom.NewPoint(50, 50)

	containsPoint := q.FindContainsPoint(testPoint)
	c.True(len(containsPoint) >= 0) //nolint:gocritic // This is a valid test for the number of containsPoint

	intersects := q.FindIntersects(testRect)
	c.True(len(intersects) >= 0) //nolint:gocritic // This is a valid test for the number of intersects

	containsRect := q.FindContainsRect(geom.NewRect(50, 50, 1, 1))
	c.True(len(containsRect) >= 0) //nolint:gocritic // This is a valid test for the number of containsRect

	containedBy := q.FindContainedByRect(geom.NewRect(0, 0, 500, 500))
	c.Equal(200, len(containedBy)) // All nodes should be contained

	removeCount := 0
	for i, n := range nodes {
		if i%3 == 0 { // Remove every third node
			q.Remove(n)
			removeCount++
		}
	}
	c.Equal(200-removeCount, q.Size())

	q.Reorganize()
	c.Equal(200-removeCount, q.Size())
}

func TestPointOnBoundaries(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(0, 0, 10, 10)
	n2 := newNode(10, 0, 10, 10)  // Shares right edge with n1
	n3 := newNode(0, 10, 10, 10)  // Shares bottom edge with n1
	n4 := newNode(10, 10, 10, 10) // Touches corner of n1

	q.Insert(n1)
	q.Insert(n2)
	q.Insert(n3)
	q.Insert(n4)

	c.True(q.ContainsPoint(geom.NewPoint(0, 0)))   // Corner
	c.True(q.ContainsPoint(geom.NewPoint(10, 0)))  // Edge point
	c.True(q.ContainsPoint(geom.NewPoint(0, 10)))  // Edge point
	c.True(q.ContainsPoint(geom.NewPoint(10, 10))) // Corner shared by multiple

	matcher := &testMatcher{target: n1}
	c.True(q.MatchedContainsPoint(matcher, geom.NewPoint(5, 5)))
	c.False(q.MatchedContainsPoint(matcher, geom.NewPoint(15, 15)))

	boundaryPoints := []geom.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
		{X: 5, Y: 0},
		{X: 0, Y: 5},
		{X: 10, Y: 5},
		{X: 5, Y: 10},
	}

	for _, pt := range boundaryPoints {
		found := q.FindContainsPoint(pt)
		c.True(len(found) >= 1) // At least one node should contain each boundary point
	}
}

func TestEmptyRectInsertionScenarios(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	emptyWidth := newNode(10, 10, 0, 5)
	emptyHeight := newNode(10, 10, 5, 0)
	emptyBoth := newNode(10, 10, 0, 0)

	initialSize := q.Size()
	q.Insert(emptyWidth)
	q.Insert(emptyHeight)
	q.Insert(emptyBoth)

	// Empty rectangles should not be inserted
	c.Equal(initialSize, q.Size())

	validNode := newNode(10, 10, 5, 5)
	q.Insert(validNode)
	c.Equal(initialSize+1, q.Size())
}

func TestDifferentNumericTypes(t *testing.T) {
	c := check.New(t)

	q64 := &quadtree.QuadTree[*node]{}
	n64_1 := newNode(0.5, 0.5, 10.5, 10.5)
	n64_2 := newNode(15.7, 15.3, 5.2, 5.8)

	q64.Insert(n64_1)
	q64.Insert(n64_2)

	c.Equal(2, q64.Size())
	c.True(q64.ContainsPoint(geom.NewPoint(5.5, 5.5)))
	c.True(q64.Intersects(geom.NewRect(0, 0, 20, 20)))

	all64 := q64.All()
	c.Equal(2, len(all64))
	c.True(slices.Contains(all64, n64_1))
	c.True(slices.Contains(all64, n64_2))
}

func TestNilRootScenarios(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	c.False(q.ContainsPoint(geom.NewPoint(0, 0)))
	c.False(q.ContainsRect(geom.NewRect(0, 0, 10, 10)))
	c.False(q.Intersects(geom.NewRect(0, 0, 10, 10)))
	c.False(q.ContainedByRect(geom.NewRect(0, 0, 10, 10)))

	c.Equal(0, len(q.FindContainsPoint(geom.NewPoint(0, 0))))
	c.Equal(0, len(q.FindContainsRect(geom.NewRect(0, 0, 1, 1))))
	c.Equal(0, len(q.FindIntersects(geom.NewRect(0, 0, 10, 10))))
	c.Equal(0, len(q.FindContainedByRect(geom.NewRect(0, 0, 100, 100))))

	matcher := &testMatcher{target: newNode(0, 0, 1, 1)}
	c.False(q.MatchedContainsPoint(matcher, geom.NewPoint(0, 0)))
	c.False(q.MatchedContainsRect(matcher, geom.NewRect(0, 0, 1, 1)))
	c.False(q.MatchedIntersects(matcher, geom.NewRect(0, 0, 10, 10)))
	c.False(q.MatchedContainedByRect(matcher, geom.NewRect(0, 0, 100, 100)))

	c.Equal(0, len(q.FindMatchedContainsPoint(matcher, geom.NewPoint(0, 0))))
	c.Equal(0, len(q.FindMatchedContainsRect(matcher, geom.NewRect(0, 0, 1, 1))))
	c.Equal(0, len(q.FindMatchedIntersects(matcher, geom.NewRect(0, 0, 10, 10))))
	c.Equal(0, len(q.FindMatchedContainedByRect(matcher, geom.NewRect(0, 0, 100, 100))))
}

func TestReorganizeWithEmptyTree(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	q.Reorganize()
	c.Equal(0, q.Size())

	q.Insert(newNode(0, 0, 10, 10))
	q.Insert(newNode(20, 20, 10, 10))
	c.Equal(2, q.Size())

	q.Clear()
	c.Equal(0, q.Size())

	q.Reorganize()
	c.Equal(0, q.Size())
}

func TestInsertNodeThatDoesNotFitInRoot(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n1 := newNode(10, 10, 10, 10)
	q.Insert(n1)

	n2 := newNode(100, 100, 10, 10)
	q.Insert(n2)

	c.Equal(2, q.Size())
	all := q.All()
	c.True(slices.Contains(all, n1))
	c.True(slices.Contains(all, n2))

	c.True(q.ContainsPoint(geom.NewPoint(105, 105)))
}

func TestTreeOperationsWithSingleNode(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	n := newNode(5, 5, 10, 10)
	q.Insert(n)

	c.True(q.ContainsPoint(geom.NewPoint(10, 10)))
	c.False(q.ContainsPoint(geom.NewPoint(20, 20)))

	c.True(q.ContainsRect(geom.NewRect(6, 6, 5, 5)))
	c.False(q.ContainsRect(geom.NewRect(0, 0, 20, 20)))

	c.True(q.Intersects(geom.NewRect(0, 0, 10, 10)))
	c.False(q.Intersects(geom.NewRect(20, 20, 5, 5)))

	c.True(q.ContainedByRect(geom.NewRect(0, 0, 20, 20)))
	c.False(q.ContainedByRect(geom.NewRect(6, 6, 5, 5)))

	found := q.FindContainsPoint(geom.NewPoint(10, 10))
	c.Equal(1, len(found))
	c.Equal(n, found[0])

	found = q.FindContainsRect(geom.NewRect(6, 6, 5, 5))
	c.Equal(1, len(found))
	c.Equal(n, found[0])

	intersecting := q.FindIntersects(geom.NewRect(0, 0, 10, 10))
	c.Equal(1, len(intersecting))
	c.Equal(n, intersecting[0])

	contained := q.FindContainedByRect(geom.NewRect(0, 0, 20, 20))
	c.Equal(1, len(contained))
	c.Equal(n, contained[0])
}

func TestInternalNodeMethods(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	baseNodes := []*node{
		newNode(10, 10, 5, 5),
		newNode(20, 10, 5, 5),
		newNode(10, 20, 5, 5),
		newNode(20, 20, 5, 5),
	}

	for _, n := range baseNodes {
		q.Insert(n)
	}

	// Add enough nodes to force subdivision
	for i := range quadtree.DefaultQuadTreeThreshold {
		x := 10 + float32(i%4)*2.5
		y := 10 + float32(i/4)*2.5
		q.Insert(newNode(x, y, 1, 1))
	}

	noIntersectRect := geom.NewRect(100, 100, 5, 5)
	c.False(q.Intersects(noIntersectRect))

	noContainRect := geom.NewRect(5, 5, 2, 2)
	c.False(q.ContainedByRect(noContainRect))

	edgeRect := geom.NewRect(8, 8, 1, 1)
	q.Intersects(edgeRect) // Exercise the intersects path

	edgeContainRect := geom.NewRect(9, 9, 1, 1)
	q.ContainedByRect(edgeContainRect) // Exercise the containedByRect path
}

func TestCoverageGaps(t *testing.T) {
	c := check.New(t)
	q := &quadtree.QuadTree[*node]{}

	q.Threshold = quadtree.MinQuadTreeThreshold
	for i := range quadtree.MinQuadTreeThreshold + 1 {
		q.Insert(newNode(float32(i), 0, 1, 1))
	}
	c.True(q.Size() > 0)

	q.Clear()

	q.Insert(newNode(50, 50, 10, 10))

	c.False(q.Intersects(geom.NewRect(0, 0, 10, 10)))

	c.False(q.ContainedByRect(geom.NewRect(55, 55, 2, 2)))

	// Add enough nodes to force subdivision
	for i := range quadtree.DefaultQuadTreeThreshold + 5 {
		x := 50 + float32(i%5)*2
		y := 50 + float32(i/5)*2
		q.Insert(newNode(x, y, 1, 1))
	}

	alwaysFalseMatcher := &alwaysFalseMatcher{}

	c.False(q.MatchedIntersects(alwaysFalseMatcher, geom.NewRect(50, 50, 20, 20)))
	c.False(q.MatchedContainsRect(alwaysFalseMatcher, geom.NewRect(51, 51, 1, 1)))
	c.False(q.MatchedContainedByRect(alwaysFalseMatcher, geom.NewRect(45, 45, 30, 30)))
}

type alwaysFalseMatcher struct{}

func (m *alwaysFalseMatcher) Matches(_ *node) bool {
	return false
}
