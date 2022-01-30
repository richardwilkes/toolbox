// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package geom provides geometry primitives.
package geom

import (
	"fmt"
	"math"
)

// Rect defines a rectangle.
type Rect struct {
	Point `json:",inline"`
	Size  `json:",inline"`
}

// NewRect creates a new Rect.
func NewRect(x, y, width, height float64) Rect {
	return Rect{
		Point: NewPoint(x, y),
		Size:  NewSize(width, height),
	}
}

// NewRectPtr creates a new *Rect.
func NewRectPtr(x, y, width, height float64) *Rect {
	r := NewRect(x, y, width, height)
	return &r
}

// CopyAndZeroLocation creates a new copy of the Rect and sets the location of the copy to 0,0.
func (r *Rect) CopyAndZeroLocation() Rect {
	return Rect{Size: r.Size}
}

// Center returns the center of the rectangle.
func (r Rect) Center() Point {
	return Point{
		X: r.CenterX(),
		Y: r.CenterY(),
	}
}

// CenterX returns the center x-coordinate of the rectangle.
func (r Rect) CenterX() float64 {
	return r.X + r.Width/2
}

// CenterY returns the center y-coordinate of the rectangle.
func (r Rect) CenterY() float64 {
	return r.Y + r.Height/2
}

// Right returns the right edge, or X + Width.
func (r Rect) Right() float64 {
	return r.X + r.Width
}

// Bottom returns the bottom edge, or Y + Height.
func (r Rect) Bottom() float64 {
	return r.Y + r.Height
}

// TopLeft returns the top-left point of the rectangle.
func (r Rect) TopLeft() Point {
	return r.Point
}

// TopRight returns the top-right point of the rectangle.
func (r Rect) TopRight() Point {
	return Point{X: r.Right() - 1, Y: r.Y}
}

// BottomRight returns the bottom-right point of the rectangle.
func (r Rect) BottomRight() Point {
	return Point{X: r.Right() - 1, Y: r.Bottom() - 1}
}

// BottomLeft returns the bottom-left point of the rectangle.
func (r Rect) BottomLeft() Point {
	return Point{X: r.X, Y: r.Bottom() - 1}
}

// Max returns the bottom right corner of the rectangle.
func (r Rect) Max() Point {
	return Point{X: r.Right(), Y: r.Bottom()}
}

// IsEmpty returns true if either the width or height is zero or less.
func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Intersects returns true if this rect and the other rect intersect.
func (r Rect) Intersects(other Rect) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.X < other.Right() && r.Y < other.Bottom() && r.Right() > other.X && r.Bottom() > other.Y
}

// IntersectsLine returns true if this rect and the line described by start and end intersect.
func (r Rect) IntersectsLine(start, end Point) bool {
	if r.IsEmpty() {
		return false
	}
	if r.ContainsPoint(start) || r.ContainsPoint(end) {
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

// ContainsPoint returns true if the coordinates are within the Rect.
func (r Rect) ContainsPoint(pt Point) bool {
	if r.IsEmpty() {
		return false
	}
	return r.X <= pt.X && r.Y <= pt.Y && pt.X < r.Right() && pt.Y < r.Bottom()
}

// ContainsRect returns true if this Rect fully contains the passed in Rect.
func (r Rect) ContainsRect(in Rect) bool {
	if r.IsEmpty() || in.IsEmpty() {
		return false
	}
	right := r.X + r.Width
	bottom := r.Y + r.Height
	inRight := in.Right() - 1
	inBottom := in.Bottom() - 1
	return r.X <= in.X && r.Y <= in.Y && in.X < right && in.Y < bottom && r.X <= inRight && r.Y <= inBottom && inRight < right && inBottom < bottom
}

// Intersect this Rect with another Rect, storing the result into this Rect. Returns itself for easy chaining.
func (r *Rect) Intersect(other Rect) *Rect {
	if r.IsEmpty() || other.IsEmpty() {
		r.Width = 0
		r.Height = 0
	} else {
		x := math.Max(r.X, other.X)
		y := math.Max(r.Y, other.Y)
		w := math.Min(r.Right(), other.Right()) - x
		h := math.Min(r.Bottom(), other.Bottom()) - y
		if w > 0 && h > 0 {
			r.X = x
			r.Y = y
			r.Width = w
			r.Height = h
		} else {
			r.Width = 0
			r.Height = 0
		}
	}
	return r
}

// Union this Rect with another Rect, storing the result into this Rect. Returns itself for easy chaining.
func (r *Rect) Union(other Rect) *Rect {
	e1 := r.IsEmpty()
	e2 := other.IsEmpty()
	switch {
	case e1 && e2:
		r.Width = 0
		r.Height = 0
	case e1:
		*r = other
	case !e2:
		x := math.Min(r.X, other.X)
		y := math.Min(r.Y, other.Y)
		r.Width = math.Max(r.Right(), other.Right()) - x
		r.Height = math.Max(r.Bottom(), other.Bottom()) - y
		r.X = x
		r.Y = y
	}
	return r
}

// Align modifies this rectangle to align with integer coordinates that would encompass the original rectangle. Returns
// itself for easy chaining.
func (r *Rect) Align() *Rect {
	x := math.Floor(r.X)
	r.Width = math.Ceil(r.Right()) - x
	r.X = x
	y := math.Floor(r.Y)
	r.Height = math.Ceil(r.Bottom()) - y
	r.Y = y
	return r
}

// InsetUniform insets this Rect by the specified amount on all sides. Positive values make the Rect smaller, while
// negative values make it larger. Returns itself for easy chaining.
func (r *Rect) InsetUniform(amount float64) *Rect {
	r.X += amount
	r.Y += amount
	r.Width -= amount * 2
	if r.Width < 0 {
		r.Width = 0
		r.Height = 0
	} else {
		r.Height -= amount * 2
		if r.Height < 0 {
			r.Width = 0
			r.Height = 0
		}
	}
	return r
}

// Inset this Rect by the specified Insets. Returns itself for easy chaining.
func (r *Rect) Inset(insets Insets) *Rect {
	r.X += insets.Left
	r.Y += insets.Top
	r.Width -= insets.Left + insets.Right
	if r.Width <= 0 {
		r.Width = 0
	}
	r.Height -= insets.Top + insets.Bottom
	if r.Height < 0 {
		r.Height = 0
	}
	return r
}

// AddPoint adds a Point to this Rect. If the Rect has a negative width or height, then the Rect's upper-left corner
// will be set to the Point and its width and height will be set to 0. Returns itself for easy chaining.
func (r *Rect) AddPoint(pt Point) *Rect {
	if r.Width < 0 || r.Height < 0 {
		r.Point = pt
		r.Width = 0
		r.Height = 0
		return r
	}
	x2 := r.Right()
	y2 := r.Bottom()
	if r.X > pt.X {
		r.X = pt.X
	}
	if r.Y > pt.Y {
		r.Y = pt.Y
	}
	if x2 < pt.X {
		x2 = pt.X
	}
	if y2 < pt.Y {
		y2 = pt.Y
	}
	r.Width = x2 - r.X
	r.Height = y2 - r.Y
	return r
}

// Bounds merely returns this rectangle.
func (r Rect) Bounds() Rect {
	return r
}

// String implements the fmt.Stringer interface.
func (r Rect) String() string {
	return fmt.Sprintf("%v,%v,%v,%v", r.X, r.Y, r.Width, r.Height)
}
