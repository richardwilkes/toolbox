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
	"slices"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

type chain[T constraints.Float] struct {
	points []geom.Point[T]
	closed bool
}

func (c *chain[T]) linkSegment(s Segment[T]) bool {
	front := c.points[0]
	back := c.points[len(c.points)-1]
	switch {
	case s.Start == front:
		if s.End == back {
			c.closed = true
		} else {
			c.points = append([]geom.Point[T]{s.End}, c.points...)
		}
		return true
	case s.End == back:
		if s.Start == front {
			c.closed = true
		} else {
			c.points = append(c.points, s.Start)
		}
		return true
	case s.End == front:
		if s.Start == back {
			c.closed = true
		} else {
			c.points = append([]geom.Point[T]{s.Start}, c.points...)
		}
		return true
	case s.Start == back:
		if s.End == front {
			c.closed = true
		} else {
			c.points = append(c.points, s.End)
		}
		return true
	}
	return false
}

func (c *chain[T]) linkChain(other *chain[T]) bool {
	front := c.points[0]
	back := c.points[len(c.points)-1]
	otherFront := other.points[0]
	otherBack := other.points[len(other.points)-1]
	switch {
	case otherFront == back:
		c.points = append(c.points, other.points[1:]...)
	case otherBack == front:
		c.points = append(other.points, c.points[1:]...)
	case otherFront == front:
		slices.Reverse(other.points)
		c.points = append(other.points, c.points[1:]...)
	case otherBack == back:
		slices.Reverse(other.points)
		c.points = append(c.points[:len(c.points)-1], other.points...)
	default:
		return false
	}
	other.points = nil
	return true
}

func (c *chain[T]) removeCoincidentSegment(s Segment[T]) bool {
	if s.CoincidesWith(c.points[0], c.points[1]) {
		c.points = c.points[1:]
		c.closed = false
		return true
	}
	last := len(c.points) - 1
	if s.CoincidesWith(c.points[last], c.points[last-1]) {
		c.points = c.points[:last]
		c.closed = false
		return true
	}
	return false
}

func (c *chain[T]) empty() bool {
	return len(c.points) < 2
}
