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

type edgeType int

const (
	edgeNormal edgeType = iota
	edgeNonContributing
	edgeSameTransition
	edgeDifferentTransition
)

type endpoint[T constraints.Float] struct {
	pt       geom.Point[T]
	edgeType edgeType
	other    *endpoint[T]
	subject  bool
	left     bool
	inout    bool
	inside   bool
}

func (e *endpoint[T]) segment() Segment[T] {
	return Segment[T]{
		Start: e.pt,
		End:   e.other.pt,
	}
}

func (e *endpoint[T]) above(pt geom.Point[T]) bool {
	return !e.below(pt)
}

func (e *endpoint[T]) below(pt geom.Point[T]) bool {
	if e.left {
		return signedArea(e.pt, e.other.pt, pt) > 0
	}
	return signedArea(e.other.pt, e.pt, pt) > 0
}

func (e *endpoint[T]) isValidDirection() bool {
	var left, right geom.Point[T]
	if e.left {
		left = e.pt
		right = e.other.pt
	} else {
		left = e.other.pt
		right = e.pt
	}
	return ptIsBefore(left, right)
}

func ptIsBefore[T constraints.Float](p1, p2 geom.Point[T]) bool {
	return p1.X < p2.X || p1.X == p2.X && p1.Y < p2.Y
}

func endpointCmp[T constraints.Float](a, b *endpoint[T]) int {
	// TODO: Use the cmp package once Go 1.21 ships
	if a.pt.X != b.pt.X {
		if a.pt.X > b.pt.X {
			return -1
		}
		return 1
	}
	if a.pt.Y != b.pt.Y {
		if a.pt.Y > b.pt.Y {
			return -1
		}
		return 1
	}
	if a.left != b.left {
		if a.left {
			return -1
		}
		return 1
	}
	if a.above(b.other.pt) {
		return -1
	}
	return 1
}

func signedArea[T constraints.Float](pt0, pt1, pt2 geom.Point[T]) T {
	return (pt0.X-pt2.X)*(pt1.Y-pt2.Y) - (pt1.X-pt2.X)*(pt0.Y-pt2.Y)
}
