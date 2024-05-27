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

package geom

import (
	"fmt"

	"github.com/richardwilkes/toolbox/xmath"
)

// Point defines a location.
type Point[T xmath.Numeric] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// NewPoint returns a new Point.
func NewPoint[T xmath.Numeric](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

// ConvertPoint converts a Point of type F into one of type T.
func ConvertPoint[T, F xmath.Numeric](pt Point[F]) Point[T] {
	return NewPoint(T(pt.X), T(pt.Y))
}

// Add returns a new Point which is the result of adding this Point with the provided Point.
func (p Point[T]) Add(pt Point[T]) Point[T] {
	return Point[T]{X: p.X + pt.X, Y: p.Y + pt.Y}
}

// Sub returns a new Point which is the result of subtracting the provided Point from this Point.
func (p Point[T]) Sub(pt Point[T]) Point[T] {
	return Point[T]{X: p.X - pt.X, Y: p.Y - pt.Y}
}

// Mul returns a new Point which is the result of multiplying the coordinates of this point by the value.
func (p Point[T]) Mul(value T) Point[T] {
	return Point[T]{X: p.X * value, Y: p.Y * value}
}

// Div returns a new Point which is the result of dividing the coordinates of this point by the value.
func (p Point[T]) Div(value T) Point[T] {
	return Point[T]{X: p.X / value, Y: p.Y / value}
}

// Neg returns a new Point that holds the negated coordinates of this Point.
func (p Point[T]) Neg() Point[T] {
	return Point[T]{X: -p.X, Y: -p.Y}
}

// Floor returns a new Point which is aligned to integer coordinates by using Floor on them.
func (p Point[T]) Floor() Point[T] {
	return Point[T]{X: xmath.Floor(p.X), Y: xmath.Floor(p.Y)}
}

// Ceil returns a new Point which is aligned to integer coordinates by using Ceil() on them.
func (p Point[T]) Ceil() Point[T] {
	return Point[T]{X: xmath.Ceil(p.X), Y: xmath.Ceil(p.Y)}
}

// Dot returns the dot product of the two Points.
func (p Point[T]) Dot(pt Point[T]) T {
	return p.X*pt.X + p.Y*pt.Y
}

// Cross returns the cross product of the two Points.
func (p Point[T]) Cross(pt Point[T]) T {
	return p.X*pt.Y - p.Y*pt.X
}

// In returns true if this Point is within the Rect.
func (p Point[T]) In(r Rect[T]) bool {
	if r.Empty() {
		return false
	}
	return r.X <= p.X && r.Y <= p.Y && p.X < r.Right() && p.Y < r.Bottom()
}

// EqualWithin returns true if the two points are within the given tolerance of each other.
func (p Point[T]) EqualWithin(pt Point[T], tolerance T) bool {
	return xmath.EqualWithin(p.X, pt.X, tolerance) && xmath.EqualWithin(p.Y, pt.Y, tolerance)
}

// String implements the fmt.Stringer interface.
func (p Point[T]) String() string {
	return fmt.Sprintf("%#v,%#v", p.X, p.Y)
}
