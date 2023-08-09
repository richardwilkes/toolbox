// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

type connector[T constraints.Float] struct {
	openPolys   []chain[T]
	closedPolys []chain[T]
}

func (c *connector[T]) add(s Segment[T]) {
	for i := range c.openPolys {
		one := &c.openPolys[i]
		if one.removeCoincidentSegment(s) {
			if one.empty() {
				c.openPolys = append(c.openPolys[0:i], c.openPolys[i+1:]...)
			}
			return
		}
	}
	for i := range c.openPolys {
		one := &c.openPolys[i]
		if !one.linkSegment(s) {
			continue
		}
		if one.closed {
			if len(one.points) == 2 {
				one.closed = false
				return
			}
			c.closedPolys = append(c.closedPolys, c.openPolys[i])
			c.openPolys = append(c.openPolys[:i], c.openPolys[i+1:]...)
			return
		}
		k := len(c.openPolys)
		for j := i + 1; j < k; j++ {
			if one.linkChain(&c.openPolys[j]) {
				c.openPolys = append(c.openPolys[:j], c.openPolys[j+1:]...)
				return
			}
		}
		return
	}
	c.openPolys = append(c.openPolys, chain[T]{points: []geom.Point[T]{s.Start, s.End}})
}

func (c *connector[T]) toPolygon() Polygon[T] {
	var poly Polygon[T]
	for _, one := range c.closedPolys {
		ct := make([]geom.Point[T], len(one.points))
		copy(ct, one.points)
		poly = append(poly, ct)
	}
	return poly
}

func (c *connector[T]) toPolyLine() Polygon[T] {
	var poly Polygon[T]
	for _, one := range c.openPolys {
		ct := make([]geom.Point[T], len(one.points))
		copy(ct, one.points)
		poly = append(poly, ct)
	}
	return poly
}
