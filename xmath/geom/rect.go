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
	"reflect"

	"github.com/richardwilkes/toolbox/xmath"
)

// Rect defines a rectangle.
type Rect[T xmath.Numeric] struct {
	Point[T] `json:",inline"`
	Size[T]  `json:",inline"`
}

// NewRect creates a new Rect.
func NewRect[T xmath.Numeric](x, y, width, height T) Rect[T] {
	return Rect[T]{
		Point: NewPoint[T](x, y),
		Size:  NewSize[T](width, height),
	}
}

// NewRectPtr creates a new *Rect.
func NewRectPtr[T xmath.Numeric](x, y, width, height T) *Rect[T] {
	r := NewRect[T](x, y, width, height)
	return &r
}

// CopyAndZeroLocation creates a new copy of the Rect and sets the location of the copy to 0,0.
func (r *Rect[T]) CopyAndZeroLocation() Rect[T] {
	return Rect[T]{Size: r.Size}
}

// Center returns the center of the rectangle.
func (r Rect[T]) Center() Point[T] {
	return Point[T]{
		X: r.CenterX(),
		Y: r.CenterY(),
	}
}

// CenterX returns the center x-coordinate of the rectangle.
func (r Rect[T]) CenterX() T {
	return r.X + r.Width/2
}

// CenterY returns the center y-coordinate of the rectangle.
func (r Rect[T]) CenterY() T {
	return r.Y + r.Height/2
}

// Right returns the right edge, or X + Width.
func (r Rect[T]) Right() T {
	return r.X + r.Width
}

// Bottom returns the bottom edge, or Y + Height.
func (r Rect[T]) Bottom() T {
	return r.Y + r.Height
}

// TopLeft returns the top-left point of the rectangle.
func (r Rect[T]) TopLeft() Point[T] {
	return r.Point
}

// TopRight returns the top-right point of the rectangle.
func (r Rect[T]) TopRight() Point[T] {
	return Point[T]{X: r.Right() - 1, Y: r.Y}
}

// BottomRight returns the bottom-right point of the rectangle.
func (r Rect[T]) BottomRight() Point[T] {
	return Point[T]{X: r.Right() - 1, Y: r.Bottom() - 1}
}

// BottomLeft returns the bottom-left point of the rectangle.
func (r Rect[T]) BottomLeft() Point[T] {
	return Point[T]{X: r.X, Y: r.Bottom() - 1}
}

// Max returns the bottom right corner of the rectangle.
func (r Rect[T]) Max() Point[T] {
	return Point[T]{X: r.Right(), Y: r.Bottom()}
}

// IsEmpty returns true if either the width or height is zero or less.
func (r Rect[T]) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Intersects returns true if this rect and the other rect intersect.
func (r Rect[T]) Intersects(other Rect[T]) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.X < other.Right() && r.Y < other.Bottom() && r.Right() > other.X && r.Bottom() > other.Y
}

// IntersectsLine returns true if this rect and the line described by start and end intersect.
func (r Rect[T]) IntersectsLine(start, end Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	if r.ContainsPoint(start) || r.ContainsPoint(end) {
		return true
	}
	s64 := start.toPoint64()
	e64 := end.toPoint64()
	rp64 := r.Point.toPoint64()
	rtrp64 := r.TopRight().toPoint64()
	if len(LineIntersection(s64, e64, rp64, rtrp64)) != 0 {
		return true
	}
	rblp64 := r.BottomLeft().toPoint64()
	if len(LineIntersection(s64, e64, rp64, rblp64)) != 0 {
		return true
	}
	rbrp64 := r.BottomRight().toPoint64()
	if len(LineIntersection(s64, e64, rtrp64, rbrp64)) != 0 {
		return true
	}
	if len(LineIntersection(s64, e64, rblp64, rbrp64)) != 0 {
		return true
	}
	return false
}

// ContainsPoint returns true if the coordinates are within the Rect.
func (r Rect[T]) ContainsPoint(pt Point[T]) bool {
	if r.IsEmpty() {
		return false
	}
	return r.X <= pt.X && r.Y <= pt.Y && pt.X < r.Right() && pt.Y < r.Bottom()
}

// ContainsRect returns true if this Rect fully contains the passed in Rect.
func (r Rect[T]) ContainsRect(in Rect[T]) bool {
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
func (r *Rect[T]) Intersect(other Rect[T]) *Rect[T] {
	if r.IsEmpty() || other.IsEmpty() {
		r.Width = 0
		r.Height = 0
	} else {
		x := xmath.Max(r.X, other.X)
		y := xmath.Max(r.Y, other.Y)
		w := xmath.Min(r.Right(), other.Right()) - x
		h := xmath.Min(r.Bottom(), other.Bottom()) - y
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
func (r *Rect[T]) Union(other Rect[T]) *Rect[T] {
	e1 := r.IsEmpty()
	e2 := other.IsEmpty()
	switch {
	case e1 && e2:
		r.Width = 0
		r.Height = 0
	case e1:
		*r = other
	case !e2:
		x := xmath.Min(r.X, other.X)
		y := xmath.Min(r.Y, other.Y)
		r.Width = xmath.Max(r.Right(), other.Right()) - x
		r.Height = xmath.Max(r.Bottom(), other.Bottom()) - y
		r.X = x
		r.Y = y
	}
	return r
}

// Align modifies this rectangle to align with integer coordinates that would encompass the original rectangle. Returns
// itself for easy chaining.
func (r *Rect[T]) Align() *Rect[T] {
	switch reflect.TypeOf(r.X).Kind() {
	case reflect.Float32, reflect.Float64:
		x := T(math.Floor(reflect.ValueOf(r.X).Float()))
		r.Width = T(math.Ceil(reflect.ValueOf(r.Right()).Float())) - x
		r.X = x
		y := T(math.Floor(reflect.ValueOf(r.Y).Float()))
		r.Height = T(math.Ceil(reflect.ValueOf(r.Bottom()).Float())) - y
		r.Y = y
	}
	return r
}

// InsetUniform insets this Rect by the specified amount on all sides. Positive values make the Rect smaller, while
// negative values make it larger. Returns itself for easy chaining.
func (r *Rect[T]) InsetUniform(amount T) *Rect[T] {
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
func (r *Rect[T]) Inset(insets Insets[T]) *Rect[T] {
	r.X += insets.Left
	r.Y += insets.Top
	r.Width -= insets.Width()
	if r.Width <= 0 {
		r.Width = 0
	}
	r.Height -= insets.Height()
	if r.Height < 0 {
		r.Height = 0
	}
	return r
}

// AddPoint adds a Point to this Rect. If the Rect has a negative width or height, then the Rect's upper-left corner
// will be set to the Point and its width and height will be set to 0. Returns itself for easy chaining.
func (r *Rect[T]) AddPoint(pt Point[T]) *Rect[T] {
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
func (r *Rect[T]) Bounds() Rect[T] {
	return *r
}

// String implements the fmt.Stringer interface.
func (r *Rect[T]) String() string {
	return fmt.Sprintf("%v,%v,%v,%v", r.X, r.Y, r.Width, r.Height)
}
