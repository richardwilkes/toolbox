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
	"github.com/richardwilkes/toolbox/v2/geom"
)

type Rect struct {
	Point
	Size
}

// NewRect creates a new Rect.
func NewRect(x, y, width, height Num) Rect {
	return Rect{
		Point: NewPoint(x, y),
		Size:  NewSize(width, height),
	}
}

// RectFrom converts a geom.Rect into a Rect.
func RectFrom(r geom.Rect) Rect {
	return Rect{
		Point: PointFrom(r.Point),
		Size:  SizeFrom(r.Size),
	}
}

// Rect converts this Rect into a geom.Rect.
func (r Rect) Rect() geom.Rect {
	return geom.Rect{
		Point: r.Point.Point(),
		Size:  r.Size.Size(),
	}
}

// Empty returns true if either the width or height is zero or less.
func (r Rect) Empty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Center returns the center of the Rect.
func (r Rect) Center() Point {
	return NewPoint(r.CenterX(), r.CenterY())
}

// CenterX returns the center x-coordinate of the Rect.
func (r Rect) CenterX() Num {
	return r.X + r.Width.Div(Two)
}

// CenterY returns the center y-coordinate of the Rect.
func (r Rect) CenterY() Num {
	return r.Y + r.Height.Div(Two)
}

// Right returns the right edge, or X + Width.
func (r Rect) Right() Num {
	return r.X + r.Width
}

// Bottom returns the bottom edge, or Y + Height.
func (r Rect) Bottom() Num {
	return r.Y + r.Height
}

// TopLeft returns the top-left point of the Rect.
func (r Rect) TopLeft() Point {
	return r.Point
}

// TopRight returns the top-right point of the Rect.
func (r Rect) TopRight() Point {
	return NewPoint(r.Right(), r.Y)
}

// BottomRight returns the bottom-right point of the Rect.
func (r Rect) BottomRight() Point {
	return NewPoint(r.Right(), r.Bottom())
}

// BottomLeft returns the bottom-left point of the Rect.
func (r Rect) BottomLeft() Point {
	return NewPoint(r.X, r.Bottom())
}

// Contains returns true if this Rect fully contains the passed in Rect.
func (r Rect) Contains(in Rect) bool {
	if r.Empty() || in.Empty() {
		return false
	}
	right := r.Right()
	bottom := r.Bottom()
	inRight := in.Right()
	inBottom := in.Bottom()
	return r.X <= in.X && r.Y <= in.Y && in.X < right && in.Y < bottom && r.X < inRight &&
		r.Y < inBottom && inRight <= right && inBottom <= bottom
}

// IntersectsLine returns true if this rect and the line described by start and end intersect.
func (r Rect) IntersectsLine(start, end Point) bool {
	if r.Empty() {
		return false
	}
	if start.In(r) || end.In(r) {
		return true
	}
	if len(LineIntersection(start, end, r.Point, r.TopRight())) != 0 {
		return true
	}
	if len(LineIntersection(start, end, r.Point, r.BottomLeft())) != 0 {
		return true
	}
	if len(LineIntersection(start, end, r.TopRight(), r.BottomRight())) != 0 {
		return true
	}
	if len(LineIntersection(start, end, r.BottomLeft(), r.BottomRight())) != 0 {
		return true
	}
	return false
}

// Intersects returns true if this Rect and the other Rect intersect.
func (r Rect) Intersects(other Rect) bool {
	if r.Empty() || other.Empty() {
		return false
	}
	return r.X < other.Right() && r.Y < other.Bottom() && other.X < r.Right() && other.Y < r.Bottom()
}

// Intersect returns the result of intersecting this Rect with another Rect.
func (r Rect) Intersect(other Rect) Rect {
	if r.Empty() || other.Empty() {
		return Rect{}
	}
	x := max(r.X, other.X)
	y := max(r.Y, other.Y)
	w := min(r.Right(), other.Right()) - x
	h := min(r.Bottom(), other.Bottom()) - y
	if w <= 0 || h <= 0 {
		return Rect{}
	}
	return NewRect(x, y, w, h)
}

// Union returns the result of unioning this Rect with another Rect.
func (r Rect) Union(other Rect) Rect {
	e1 := r.Empty()
	e2 := other.Empty()
	switch {
	case e1 && e2:
		return Rect{}
	case e1:
		return other
	case e2:
		return r
	default:
		x := min(r.X, other.X)
		y := min(r.Y, other.Y)
		return NewRect(x, y, max(r.Right(), other.Right())-x, max(r.Bottom(), other.Bottom())-y)
	}
}

// Align returns a new Rect aligned with integer coordinates that would encompass the original rectangle.
func (r Rect) Align() Rect {
	return Rect{Point: r.Point.Floor(), Size: r.Size.Ceil()}
}

// Expand returns a new Rect that expands this Rect to encompass the provided Point. If the Rect has a negative width or
// height, then the Rect's upper-left corner will be set to the Point and its width and height will be set to 0.
func (r Rect) Expand(pt Point) Rect {
	if r.Width < 0 || r.Height < 0 {
		return Rect{Point: pt}
	}
	x := min(r.X, pt.X)
	y := min(r.Y, pt.Y)
	return NewRect(x, y, max(r.Right(), pt.X)-x, max(r.Bottom(), pt.Y)-y)
}

// String implements fmt.Stringer.
func (r Rect) String() string {
	return r.Point.String() + "," + r.Size.String()
}
