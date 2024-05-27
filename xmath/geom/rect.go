/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

// Package geom provides geometry primitives.
package geom

import (
	"github.com/richardwilkes/toolbox/xmath"
)

// Rect defines a rectangle.
type Rect[T xmath.Numeric] struct {
	Point[T] `json:",inline"`
	Size[T]  `json:",inline"`
}

// NewRect creates a new Rect.
func NewRect[T xmath.Numeric](x, y, width, height T) Rect[T] {
	return Rect[T]{Point: NewPoint[T](x, y), Size: NewSize[T](width, height)}
}

// ConvertRect converts a Rect of type F into one of type T.
func ConvertRect[T, F xmath.Numeric](r Rect[F]) Rect[T] {
	return NewRect(T(r.X), T(r.Y), T(r.Width), T(r.Height))
}

// Empty returns true if either the width or height is zero or less.
func (r Rect[T]) Empty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Center returns the center of the Rect.
func (r Rect[T]) Center() Point[T] {
	return NewPoint(r.CenterX(), r.CenterY())
}

// CenterX returns the center x-coordinate of the Rect.
func (r Rect[T]) CenterX() T {
	return r.X + r.Width/2
}

// CenterY returns the center y-coordinate of the Rect.
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

// TopLeft returns the top-left point of the Rect.
func (r Rect[T]) TopLeft() Point[T] {
	return r.Point
}

// TopRight returns the top-right point of the Rect.
func (r Rect[T]) TopRight() Point[T] {
	return NewPoint(r.Right(), r.Y)
}

// BottomRight returns the bottom-right point of the Rect.
func (r Rect[T]) BottomRight() Point[T] {
	return NewPoint(r.Right(), r.Bottom())
}

// BottomLeft returns the bottom-left point of the Rect.
func (r Rect[T]) BottomLeft() Point[T] {
	return NewPoint(r.X, r.Bottom())
}

// Contains returns true if this Rect fully contains the passed in Rect.
func (r Rect[T]) Contains(in Rect[T]) bool {
	if r.Empty() || in.Empty() {
		return false
	}
	right := r.Right()
	bottom := r.Bottom()
	inRight := in.Right() - 1
	inBottom := in.Bottom() - 1
	return r.X <= in.X && r.Y <= in.Y && in.X < right && in.Y < bottom && r.X <= inRight &&
		r.Y <= inBottom && inRight < right && inBottom < bottom
}

// IntersectsLine returns true if this rect and the line described by start and end intersect.
func (r Rect[T]) IntersectsLine(start, end Point[T]) bool {
	if r.Empty() {
		return false
	}
	if start.In(r) || end.In(r) {
		return true
	}
	if len(LineIntersection[T](start, end, r.Point, r.TopRight())) != 0 {
		return true
	}
	if len(LineIntersection[T](start, end, r.Point, r.BottomLeft())) != 0 {
		return true
	}
	if len(LineIntersection[T](start, end, r.TopRight(), r.BottomRight())) != 0 {
		return true
	}
	if len(LineIntersection[T](start, end, r.BottomLeft(), r.BottomRight())) != 0 {
		return true
	}
	return false
}

// Intersects returns true if this Rect and the other Rect intersect.
func (r Rect[T]) Intersects(other Rect[T]) bool {
	if r.Empty() || other.Empty() {
		return false
	}
	return r.X < other.Right() && r.Y < other.Bottom() && r.Right() > other.X && r.Bottom() > other.Y
}

// Intersect returns the result of intersecting this Rect with another Rect.
func (r Rect[T]) Intersect(other Rect[T]) Rect[T] {
	if r.Empty() || other.Empty() {
		return Rect[T]{}
	}
	x := max(r.X, other.X)
	y := max(r.Y, other.Y)
	w := min(r.Right(), other.Right()) - x
	h := min(r.Bottom(), other.Bottom()) - y
	if w <= 0 || h <= 0 {
		return Rect[T]{}
	}
	return NewRect(x, y, w, h)
}

// Union returns the result of unioning this Rect with another Rect.
func (r Rect[T]) Union(other Rect[T]) Rect[T] {
	e1 := r.Empty()
	e2 := other.Empty()
	switch {
	case e1 && e2:
		return Rect[T]{}
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
func (r Rect[T]) Align() Rect[T] {
	return Rect[T]{Point: r.Point.Floor(), Size: r.Size.Ceil()}
}

// Expand returns a new Rect that expands this Rect to encompass the provided Point. If the Rect has a negative width or
// height, then the Rect's upper-left corner will be set to the Point and its width and height will be set to 0.
func (r Rect[T]) Expand(pt Point[T]) Rect[T] {
	if r.Width < 0 || r.Height < 0 {
		return Rect[T]{Point: pt}
	}
	x := min(r.X, pt.X)
	y := min(r.Y, pt.Y)
	return NewRect(x, y, max(r.Right(), pt.X)-x, max(r.Bottom(), pt.Y)-y)
}

// Inset returns a new Rect which has been inset by the specified Insets.
func (r Rect[T]) Inset(insets Insets[T]) Rect[T] {
	return NewRect(r.X+insets.Left, r.Y+insets.Top, max(r.Width-insets.Width(), 0), max(r.Height-insets.Height(), 0))
}

// String implements fmt.Stringer.
func (r Rect[T]) String() string {
	return r.Point.String() + "," + r.Size.String()
}
