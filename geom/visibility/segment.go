// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package visibility

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/geom/poly"
)

// Segment holds the start and end points of a line.
type Segment struct {
	Start poly.Point
	End   poly.Point
}

// Bounds returns the bounding geom.Rect of this Segment for use in a quadtree. This includes a slight bit of expansion
// to compensate for floating-point imprecision.
func (s Segment) Bounds() geom.Rect {
	minX := min(s.Start.X, s.End.X)
	minY := min(s.Start.Y, s.End.Y)
	return poly.NewRect(
		minX-epsilon,
		minY-epsilon,
		max(s.Start.X, s.End.X)-minX+twoEpsilon,
		max(s.Start.Y, s.End.Y)-minY+twoEpsilon,
	).Rect()
}
