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
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

type boundsNode struct {
	geom.Rect
}

func (n boundsNode) Bounds() geom.Rect { return n.Rect }

// TestSplitQuadrants verifies that splitting a non-square node produces four correctly-sized
// quadrants. A square node would mask a width/height mix-up in the top-left child, so a
// deliberately wide rect is used here.
func TestSplitQuadrants(t *testing.T) {
	c := check.New(t)

	n := &node[boundsNode]{
		rect:      geom.NewRect(0, 0, 1000, 10),
		threshold: 1,
	}
	for i := range 2 {
		n.contents = append(n.contents, boundsNode{Rect: geom.NewRect(float32(i), 0, 1, 1)})
	}
	n.splitIfNeeded()

	c.False(n.isLeaf())
	hw := float32(500)
	hh := float32(5)
	c.Equal(geom.NewRect(0, 0, hw, hh), n.children[0].rect)   // top-left
	c.Equal(geom.NewRect(hw, 0, hw, hh), n.children[1].rect)  // top-right
	c.Equal(geom.NewRect(0, hh, hw, hh), n.children[2].rect)  // bottom-left
	c.Equal(geom.NewRect(hw, hh, hw, hh), n.children[3].rect) // bottom-right
}

// TestStraddlingContentsAccumulate documents the degenerate-input behavior described on splitIfNeeded and the QuadTree
// type: objects whose bounds straddle a cell's center line fit none of the children, so they pile into a single cell
// that cannot subdivide further. This verifies both that such input produces exactly that shape and that queries over
// the degenerate cell still return correct results despite scanning it linearly.
func TestStraddlingContentsAccumulate(t *testing.T) {
	c := check.New(t)

	const threshold = 4
	n := &node[boundsNode]{
		rect:      geom.NewRect(0, 0, 100, 100),
		threshold: threshold,
	}
	// Every object is centered on (50, 50) — the cell's center — so none can descend into any child, no matter how
	// many are inserted or how far past the threshold the cell grows.
	const count = 3 * threshold
	for i := range count {
		half := float32(i + 1)
		n.insert(boundsNode{Rect: geom.NewRect(50-half, 50-half, 2*half, 2*half)})
	}

	// The cell subdivided once (it is no longer a leaf), but every object stayed behind, so the children are empty
	// leaves and contents holds them all.
	c.False(n.isLeaf())
	c.Equal(count, len(n.contents))
	for _, child := range n.children {
		c.True(child.isLeaf())
		c.Equal(0, len(child.contents))
	}

	// Despite the degeneracy, a query over the cell still returns every matching object.
	c.Equal(count, len(n.findContainsPoint(geom.NewPoint(50, 50), nil)))
}
