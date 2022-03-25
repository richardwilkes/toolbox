// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package quadtree_test

import (
	"math/rand"
	"testing"

	"github.com/richardwilkes/toolbox/collection/quadtree"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/stretchr/testify/assert"
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
	q := &quadtree.QuadTree[float64]{}
	assert.False(t, q.ContainsPoint(geom.Point[float64]{}))
	q.Insert(newNode[float64](5, 5, 5, 5))
	assert.False(t, q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	assert.True(t, q.ContainsPoint(geom.NewPoint[float64](5, 5)))
	assert.True(t, q.ContainsPoint(geom.NewPoint[float64](9.9, 9.9)))
	assert.False(t, q.ContainsPoint(geom.NewPoint[float64](10, 10)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	assert.True(t, q.ContainsPoint(geom.NewPoint[float64](6, 4)))
	for i := 0; i < 2*quadtree.DefaultQuadTreeThreshold; i++ {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	assert.True(t, q.ContainsPoint(geom.Point[float64]{}))
	assert.True(t, q.ContainsPoint(geom.NewPoint[float64](0, -5)))
	assert.False(t, q.ContainsPoint(geom.NewPoint[float64](-1, 0)))
}

func TestContainsRect(t *testing.T) {
	q := &quadtree.QuadTree[float64]{}
	assert.False(t, q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	q.Insert(newNode[float64](5, 5, 5, 5))
	assert.False(t, q.ContainsRect(geom.NewRect[float64](4, 4, 10, 10)))
	assert.True(t, q.ContainsRect(geom.NewRect[float64](5, 5, 2, 2)))
	assert.True(t, q.ContainsRect(geom.NewRect[float64](9.9, 9.9, .05, .05)))
	assert.False(t, q.ContainsRect(geom.NewRect[float64](10, 10, 5, 5)))
	q.Insert(newNode[float64](4, 4, 3, 3))
	assert.True(t, q.ContainsRect(geom.NewRect[float64](6, 4, 1, 2)))
	for i := 0; i < 2*quadtree.DefaultQuadTreeThreshold; i++ {
		q.Insert(newNode(float64(i), -5, 10, 10))
	}
	assert.True(t, q.ContainsRect(geom.NewRect[float64](0, 0, 1, 1)))
	assert.True(t, q.ContainsRect(geom.NewRect[float64](0, -5, 4, 4)))
	assert.False(t, q.ContainsRect(geom.NewRect[float64](-1, 0, 2, 2)))
}

func TestGeneral(t *testing.T) {
	q := &quadtree.QuadTree[float64]{}
	r := rand.New(rand.NewSource(22))
	mine := newNode[float64](22, 22, 22, 22)
	q.Insert(mine)
	for i := 0; i < 100*quadtree.DefaultQuadTreeThreshold; i++ {
		q.Insert(newNode(float64(50000-r.Intn(100000)), float64(50000-r.Intn(100000)), float64(r.Intn(100000)), float64(r.Intn(100000))))
	}
	assert.Equal(t, 1+100*quadtree.DefaultQuadTreeThreshold, q.Size())
	assert.Subset(t, q.All(), []quadtree.Node[float64]{mine})
	count := q.Size()
	for _, one := range q.All() {
		if one != mine && r.Intn(10) == 1 {
			q.Remove(one)
			count--
			assert.Equal(t, count, q.Size())
		}
	}
	assert.Equal(t, count, q.Size())
	q.Reorganize()
	assert.Equal(t, count, q.Size())
	assert.Subset(t, q.All(), []quadtree.Node[float64]{mine})
	assert.Subset(t, q.FindContainedByRect(mine.Rect), []quadtree.Node[float64]{mine})
}
