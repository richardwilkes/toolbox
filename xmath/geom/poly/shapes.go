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
	"math"

	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"golang.org/x/exp/constraints"
)

// CalcEllipseSegmentCount returns a suggested number of segments to use when generating an ellipse. 'r' is the largest
// radius of the ellipse. 'e' is the acceptable error, typically 1 or less.
func CalcEllipseSegmentCount[T constraints.Float](r, e T) int {
	d := 1 - e/r
	n := int(xmath.Ceil(2 * math.Pi / xmath.Acos(2*d*d-1)))
	if n < 4 {
		n = 4
	}
	return n
}

// ApproximateEllipseAuto creates a polygon that approximates an ellipse, automatically choose the number of segments to
// break the ellipse contour into. This uses CalcEllipseSegmentCount() with an 'e' of 0.2.
func ApproximateEllipseAuto[T constraints.Float](bounds geom.Rect[T]) Polygon[T] {
	return ApproximateEllipse(bounds, CalcEllipseSegmentCount(max(bounds.Width, bounds.Height)/2, 0.2))
}

// ApproximateEllipse creates a polygon that approximates an ellipse. 'sections' indicates how many segments to break
// the ellipse contour into.
func ApproximateEllipse[T constraints.Float](bounds geom.Rect[T], sections int) Polygon[T] {
	halfWidth := bounds.Width / 2
	halfHeight := bounds.Height / 2
	inc := math.Pi * 2 / T(sections)
	center := bounds.Center()
	contour := make(Contour[T], sections)
	var angle T
	for i := 0; i < sections; i++ {
		contour[i] = geom.Point[T]{
			X: center.X + xmath.Cos(angle)*halfWidth,
			Y: center.Y + xmath.Sin(angle)*halfHeight,
		}
		angle += inc
	}
	return Polygon[T]{contour}
}

// Rect creates a new polygon in the shape of a rectangle.
func Rect[T constraints.Float](bounds geom.Rect[T]) Polygon[T] {
	right := bounds.Right() - 1
	bottom := bounds.Bottom() - 1
	return Polygon[T]{Contour[T]{
		bounds.Point,
		geom.Point[T]{X: bounds.X, Y: bottom},
		geom.Point[T]{X: right, Y: bottom},
		geom.Point[T]{X: right, Y: bounds.Y},
	}}
}
