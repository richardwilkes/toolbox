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
	"fmt"

	"github.com/richardwilkes/toolbox/v2/geom"
)

// Point holds a fixed-point 2D coordinate.
type Point struct {
	X Num
	Y Num
}

// NewPoint returns a new Point.
func NewPoint(x, y Num) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// PointFrom converts a geom.Point into a Point.
func PointFrom(p geom.Point) Point {
	return Point{
		X: NumFromFloat(p.X),
		Y: NumFromFloat(p.Y),
	}
}

// Point converts this Point into a geom.Point.
func (p Point) Point() geom.Point {
	return geom.Point{
		X: NumAsFloat[float32](p.X),
		Y: NumAsFloat[float32](p.Y),
	}
}

// Add returns a new Point which is the result of adding this Point with the provided Point.
func (p Point) Add(pt Point) Point {
	return Point{
		X: p.X + pt.X,
		Y: p.Y + pt.Y,
	}
}

// Sub returns a new Point which is the result of subtracting the provided Point from this Point.
func (p Point) Sub(pt Point) Point {
	return Point{
		X: p.X - pt.X,
		Y: p.Y - pt.Y,
	}
}

// Mul returns a new Point which is the result of multiplying the coordinates of this point by the value.
func (p Point) Mul(value Num) Point {
	return Point{
		X: p.X * value,
		Y: p.Y * value,
	}
}

// Div returns a new Point which is the result of dividing the coordinates of this point by the value.
func (p Point) Div(value Num) Point {
	return Point{
		X: p.X / value,
		Y: p.Y / value,
	}
}

// Neg returns a new Point that holds the negated coordinates of this Point.
func (p Point) Neg() Point {
	return Point{
		X: -p.X,
		Y: -p.Y,
	}
}

// Floor returns a new Point which is aligned to integer coordinates by using Floor on them.
func (p Point) Floor() Point {
	return Point{
		X: p.X.Floor(),
		Y: p.Y.Floor(),
	}
}

// Ceil returns a new Point which is aligned to integer coordinates by using Ceil() on them.
func (p Point) Ceil() Point {
	return Point{
		X: p.X.Ceil(),
		Y: p.Y.Ceil(),
	}
}

// Dot returns the dot product of the two Points.
func (p Point) Dot(pt Point) Num {
	return p.X*pt.X + p.Y*pt.Y
}

// Cross returns the cross product of the two Points.
func (p Point) Cross(pt Point) Num {
	return p.X*pt.Y - p.Y*pt.X
}

// In returns true if this Point is within the Rect.
func (p Point) In(r Rect) bool {
	if r.Empty() {
		return false
	}
	return r.X <= p.X && r.Y <= p.Y && p.X < r.Right() && p.Y < r.Bottom()
}

func (p Point) String() string {
	return fmt.Sprintf("%v,%v", p.X, p.Y)
}
