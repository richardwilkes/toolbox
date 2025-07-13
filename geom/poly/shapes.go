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
)

// FromRect returns a Polygon in the shape of the specified rectangle.
func FromRect(r geom.Rect) Polygon {
	right := r.Right() - 1
	bottom := r.Bottom() - 1
	return Polygon{Contour{r.Point, geom.NewPoint(r.X, bottom), geom.NewPoint(right, bottom), geom.NewPoint(right, r.Y)}}
}

// FromEllipse returns a Polygon that approximates an ellipse filling the given Rect. 'sections' indicates how many
// segments to break the ellipse contour into. Passing a value less than 4 for 'sections' will result in an automatic
// choice based on a call to EllipseSegmentCount, using half of the longest dimension for the 'r' parameter and 0.2 for
// the 'e' parameter.
func FromEllipse(r geom.Rect, sections int) Polygon {
	if sections < 4 {
		sections = EllipseSegmentCount(max(r.Width, r.Height)/2, 0.2)
	}
	halfWidth := r.Width / 2
	halfHeight := r.Height / 2
	inc := math.Pi * 2 / float32(sections)
	center := r.Center()
	contour := make(Contour, sections)
	var angle float32
	for i := range sections {
		contour[i] = geom.NewPoint(center.X+xmath.Cos(angle)*halfWidth, center.Y+xmath.Sin(angle)*halfHeight)
		angle += inc
	}
	return Polygon{contour}
}

// EllipseSegmentCount returns a suggested number of segments to use when generating an ellipse. 'r' is the largest
// radius of the ellipse. 'e' is the acceptable error, typically 1 or less.
func EllipseSegmentCount(r, e float32) int {
	d := 1 - e/r
	return max(int(xmath.Ceil(2*math.Pi/xmath.Acos(2*d*d-1))), 4)
}
