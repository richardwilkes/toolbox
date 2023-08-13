// Copyright ©2019-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

// Segment32 is an alias for the float32 version of Segment.
type Segment32 = Segment[float32]

// Segment64 is an alias for the float64 version of Segment.
type Segment64 = Segment[float64]

// Segment holds the start and end points of a line.
type Segment[T constraints.Float] struct {
	Start geom.Point[T]
	End   geom.Point[T]
}

// Bounds returns the bounding rectangle of this Segment. This includes a slight bit of expansion to compensate for
// floating-point imprecision.
func (s Segment[T]) Bounds() geom.Rect[T] {
	minX := min(s.Start.X, s.End.X)
	minY := min(s.Start.Y, s.End.Y)
	return geom.Rect[T]{
		Point: geom.Point[T]{
			X: minX - epsilon,
			Y: minY - epsilon,
		},
		Size: geom.Size[T]{
			Width:  max(s.Start.X, s.End.X) - minX + epsilon*2,
			Height: max(s.Start.Y, s.End.Y) - minY + epsilon*2,
		},
	}
}
