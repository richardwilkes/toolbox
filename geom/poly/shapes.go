// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

// FromRect returns a Polygon in the shape of the specified rectangle.
func FromRect(r Rect) Polygon {
	right := r.Right()
	bottom := r.Bottom()
	return Polygon{Contour{
		r.Point,
		NewPoint(r.X, bottom),
		NewPoint(right, bottom),
		NewPoint(right, r.Y),
	}}
}

// FromEllipse returns a Polygon that approximates an ellipse filling the given Rect. 'sections' indicates how many
// segments to break the ellipse contour into. Passing a value less than 4 for 'sections' will result in an automatic
// choice based on a call to EllipseSegmentCount, using half of the longest dimension for the 'r' parameter and 0.2 for
// the 'e' parameter.
func FromEllipse(r Rect, sections Num) Polygon {
	if sections < Four {
		sections = EllipseSegmentCount(max(r.Width, r.Height).Div(Two), PointTwo)
	}
	halfWidth := r.Width.Div(Two)
	halfHeight := r.Height.Div(Two)
	inc := Pi.Mul(Two).Div(sections)
	center := r.Center()
	contour := make(Contour, NumAsInteger[int](sections))
	var angle Num
	for i := range sections {
		contour[i] = NewPoint(center.X+Cos(angle).Mul(halfWidth), center.Y+Sin(angle).Mul(halfHeight))
		angle += inc
	}
	return Polygon{contour}
}

// EllipseSegmentCount returns a suggested number of segments to use when generating an ellipse. 'r' is the largest
// radius of the ellipse. 'e' is the acceptable error, typically 1 or less.
func EllipseSegmentCount(r, e Num) Num {
	d := One - e.Div(r)
	return max(((Two.Mul(Pi).Div(Acos(Two.Mul(d).Mul(d) - One))).Ceil()), Four)
}
