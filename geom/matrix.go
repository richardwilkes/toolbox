// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"fmt"

	"github.com/richardwilkes/toolbox/v2/xmath"
	"golang.org/x/exp/constraints"
)

// Matrix provides a 2D matrix.
type Matrix[T constraints.Float] struct {
	ScaleX T `json:"scale_x"`
	SkewX  T `json:"skew_x"`
	TransX T `json:"trans_x"`
	SkewY  T `json:"skew_y"`
	ScaleY T `json:"scale_y"`
	TransY T `json:"trans_y"`
}

// NewIdentityMatrix creates a new identity transformation Matrix.
func NewIdentityMatrix[T constraints.Float]() Matrix[T] {
	return Matrix[T]{
		ScaleX: 1,
		ScaleY: 1,
	}
}

// NewTranslationMatrix creates a new Matrix that translates by 'tx' and 'ty'.
func NewTranslationMatrix[T constraints.Float](tx, ty T) Matrix[T] {
	return Matrix[T]{
		ScaleX: 1,
		ScaleY: 1,
		TransX: tx,
		TransY: ty,
	}
}

// NewScaleMatrix creates a new Matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix[T constraints.Float](sx, sy T) Matrix[T] {
	return Matrix[T]{
		ScaleX: sx,
		ScaleY: sy,
	}
}

// NewRotationMatrix creates a new Matrix that rotates by 'radians'. Positive values are clockwise.
func NewRotationMatrix[T constraints.Float](radians T) Matrix[T] {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return Matrix[T]{
		ScaleX: c,
		SkewX:  -s,
		SkewY:  s,
		ScaleY: c,
	}
}

// NewRotationByDegreesMatrix creates a new Matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationByDegreesMatrix[T constraints.Float](degrees T) Matrix[T] {
	return NewRotationMatrix(degrees * xmath.DegreesToRadians)
}

// Translate returns a new Matrix which is a copy of this Matrix translated by 'tx' and 'ty'.
func (m Matrix[T]) Translate(tx, ty T) Matrix[T] {
	return Matrix[T]{
		ScaleX: m.ScaleX,
		SkewX:  m.SkewX,
		TransX: m.TransX + tx,
		SkewY:  m.SkewY,
		ScaleY: m.ScaleY,
		TransY: m.TransY + ty,
	}
}

// Scale returns a new Matrix which is a copy of this Matrix scaled by 'sx' and 'sy'.
func (m Matrix[T]) Scale(sx, sy T) Matrix[T] {
	return Matrix[T]{
		ScaleX: m.ScaleX * sx,
		SkewX:  m.SkewX * sx,
		TransX: m.TransX * sx,
		SkewY:  m.SkewY * sy,
		ScaleY: m.ScaleY * sy,
		TransY: m.TransY * sy,
	}
}

// Rotate returns a new Matrix which is a copy of this Matrix rotated by 'radians'. Positive values are clockwise.
func (m Matrix[T]) Rotate(radians T) Matrix[T] {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return Matrix[T]{
		ScaleX: m.ScaleX*c - s*m.SkewY,
		SkewX:  m.SkewX*c - s*m.ScaleY,
		TransX: m.TransX*c - s*m.TransY,
		SkewY:  m.ScaleX*s + m.SkewY*c,
		ScaleY: m.SkewX*s + m.ScaleY*c,
		TransY: m.TransX*s + m.TransY*c,
	}
}

// RotateByDegrees returns a new Matrix which is a copy of this Matrix rotated by 'degrees'. Positive values are clockwise.
func (m Matrix[T]) RotateByDegrees(degrees T) Matrix[T] {
	return m.Rotate(degrees * xmath.DegreesToRadians)
}

// Multiply returns this Matrix multiplied by the other Matrix.
func (m Matrix[T]) Multiply(other Matrix[T]) Matrix[T] {
	return Matrix[T]{
		ScaleX: m.ScaleX*other.ScaleX + m.SkewX*other.SkewY,
		SkewX:  m.ScaleX*other.SkewX + m.SkewX*other.ScaleY,
		TransX: m.ScaleX*other.TransX + m.SkewX*other.TransY + m.TransX,
		SkewY:  m.SkewY*other.ScaleX + m.ScaleY*other.SkewY,
		ScaleY: m.SkewY*other.SkewX + m.ScaleY*other.ScaleY,
		TransY: m.SkewY*other.TransX + m.ScaleY*other.TransY + m.TransY,
	}
}

// TransformPoint returns the result of transforming the given Point by this Matrix.
func (m Matrix[T]) TransformPoint(p Point[T]) Point[T] {
	return Point[T]{
		X: m.ScaleX*p.X + m.SkewX*p.Y + m.TransX,
		Y: m.SkewY*p.X + m.ScaleY*p.Y + m.TransY,
	}
}

// String implements fmt.Stringer.
func (m Matrix[T]) String() string {
	return fmt.Sprintf("%#v,%#v,%#v,%#v,%#v,%#v", m.ScaleX, m.SkewX, m.TransX, m.SkewY, m.ScaleY, m.TransY)
}
