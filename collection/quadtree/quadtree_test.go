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

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/collection/quadtree"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

type node[T ~float32 | ~float64] struct {
	geom.Rect[T]
}

func newNode[T ~float32 | ~float64](x, y, width, height T) *node[T] {
	return &node[T]{Rect: geom.NewRect(x, y, width, height)}
}

func (n node[T]) Bounds() geom.Rect[T] {
	return n.Rect
}

func TestContainsPoint(t *testing.T) {
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	check.False(t, q.ContainsPoint(geom.Point[float64]{}))
	q.Insert(newNode[float64](5, 5, 5, 5))
	check.False(t, q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	check.True(t, q.ContainsPoint(geom.NewPoint[float64](5, 5)))
	check.True(t, q.ContainsPoint(geom.NewPoint(9.9, 9.9)))
	check.False(t, q.ContainsPoint(geom.NewPoint[float64](10, 10)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	check.True(t, q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	check.True(t, q.ContainsPoint(geom.Point[float64]{}))
	check.True(t, q.ContainsPoint(geom.NewPoint[float64](0, -5)))
	check.False(t, q.ContainsPoint(geom.NewPoint[float64](-1, 0)))
}

func TestContainsRect(t *testing.T) {
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	check.False(t, q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	q.Insert(newNode[float64](5, 5, 5, 5))
	check.False(t, q.ContainsRect(geom.NewRect[float64](4, 4, 10, 10)))
	check.True(t, q.ContainsRect(geom.NewRect[float64](5, 5, 2, 2)))
	check.True(t, q.ContainsRect(geom.NewRect(9.9, 9.9, .05, .05)))
	check.False(t, q.ContainsRect(geom.NewRect[float64](10, 10, 5, 5)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	check.True(t, q.ContainsRect(geom.NewRect[float64](6, 4, 1, 2)))
	for i := range 2 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	check.True(t, q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	check.True(t, q.ContainsRect(geom.NewRect[float64](0, -5, 4, 4)))
	check.False(t, q.ContainsRect(geom.NewRect[float64](-1, 0, 2, 2)))
}

func TestGeneral(t *testing.T) {
	q := &quadtree.QuadTree[float64, *node[float64]]{}
	r := rand.New(rand.NewPCG(22, 1967)) //nolint:gosec // Yes, it is ok to use a weak prng here
	mine := newNode[float64](22, 22, 22, 22)
	q.Insert(mine)
	for range 100 * quadtree.DefaultQuadTreeThreshold {
		q.Insert(newNode(float64(50000-r.IntN(100000)), float64(50000-r.IntN(100000)), float64(r.IntN(100000)), float64(r.IntN(100000))))
	}
	check.Equal(t, 1+100*quadtree.DefaultQuadTreeThreshold, q.Size())
	all := q.All()
	check.True(t, slices.Contains(all, mine))
	count := q.Size()
	for _, one := range all {
		if one != mine && r.IntN(10) == 1 {
			q.Remove(one)
			count--
			check.Equal(t, count, q.Size())
		}
	}
	check.Equal(t, count, q.Size())
	q.Reorganize()
	check.Equal(t, count, q.Size())
	check.True(t, slices.Contains(q.All(), mine))
	check.True(t, slices.Contains(q.FindContainedByRect(mine.Rect), mine))
}
