// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/xmath"
	"golang.org/x/exp/constraints"
)

// FromRect returns a Polygon in the shape of the specified rectangle.
func FromRect[T constraints.Float](r geom.Rect[T]) Polygon[T] {
	right := r.Right() - 1
	bottom := r.Bottom() - 1
	return Polygon[T]{Contour[T]{
		r.Point,
		geom.NewPoint(r.X, bottom),
		geom.NewPoint(right, bottom),
		geom.NewPoint(right, r.Y),
	}}
}

// FromEllipse returns a Polygon that approximates an ellipse filling the given Rect. 'sections' indicates how many
// segments to break the ellipse contour into. Passing a value less than 4 for 'sections' will result in an automatic
// choice based on a call to EllipseSegmentCount, using half of the longest dimension for the 'r' parameter and 0.2 for
// the 'e' parameter.
func FromEllipse[T constraints.Float](r geom.Rect[T], sections int) Polygon[T] {
	if sections < 4 {
		sections = EllipseSegmentCount(max(r.Width, r.Height)/2, 0.2)
	}
	halfWidth := r.Width / 2
	halfHeight := r.Height / 2
	inc := math.Pi * 2 / T(sections)
	center := r.Center()
	contour := make(Contour[T], sections)
	var angle T
	for i := range sections {
		contour[i] = geom.NewPoint(center.X+xmath.Cos(angle)*halfWidth, center.Y+xmath.Sin(angle)*halfHeight)
		angle += inc
	}
	return Polygon[T]{contour}
}

// EllipseSegmentCount returns a suggested number of segments to use when generating an ellipse. 'r' is the largest
// radius of the ellipse. 'e' is the acceptable error, typically 1 or less.
func EllipseSegmentCount[T constraints.Float](r, e T) int {
	d := 1 - e/r
	return max(int(xmath.Ceil(2*math.Pi/xmath.Acos(2*d*d-1))), 4)
}
