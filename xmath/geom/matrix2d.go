// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"github.com/richardwilkes/toolbox/xmath"
	"golang.org/x/exp/constraints"
)

// Matrix2D provides a 2D matrix.
type Matrix2D[T constraints.Float] struct {
	ScaleX T
	SkewX  T
	TransX T
	SkewY  T
	ScaleY T
	TransY T
}

// NewIdentityMatrix2D creates a new identity transformation 2D matrix.
func NewIdentityMatrix2D[T constraints.Float]() *Matrix2D[T] {
	return &Matrix2D[T]{ScaleX: 1, ScaleY: 1}
}

// NewTranslationMatrix2D creates a new 2D matrix that translates by 'tx' and 'ty'.
func NewTranslationMatrix2D[T constraints.Float](tx, ty T) *Matrix2D[T] {
	return &Matrix2D[T]{ScaleX: 1, ScaleY: 1, TransX: tx, TransY: ty}
}

// NewScaleMatrix2D creates a new 2D matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix2D[T constraints.Float](sx, sy T) *Matrix2D[T] {
	return &Matrix2D[T]{ScaleX: sx, ScaleY: sy}
}

// NewRotationMatrix2D creates a new 2D matrix that rotates by 'radians'. Positive values are clockwise.
func NewRotationMatrix2D[T constraints.Float](radians T) *Matrix2D[T] {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return &Matrix2D[T]{ScaleX: c, SkewX: -s, SkewY: s, ScaleY: c}
}

// NewRotationByDegreesMatrix2D creates a new 2D matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationByDegreesMatrix2D[T constraints.Float](degrees T) *Matrix2D[T] {
	return NewRotationMatrix2D[T](degrees * xmath.DegreesToRadians)
}

// Translate this matrix by 'tx' and 'ty'.
func (m *Matrix2D[T]) Translate(tx, ty T) {
	m.TransX += tx
	m.TransY += ty
}

// Scale this matrix by 'sx' and 'sy'.
func (m *Matrix2D[T]) Scale(sx, sy T) {
	m.ScaleX *= sx
	m.SkewX *= sx
	m.TransX *= sx
	m.SkewY *= sy
	m.ScaleY *= sy
	m.TransY *= sy
}

// Rotate this matrix by 'radians'. Positive values are clockwise.
func (m *Matrix2D[T]) Rotate(radians T) {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	x := m.ScaleX*c - s*m.SkewY
	m.SkewY = m.ScaleX*s + m.SkewY*c
	m.ScaleX = x
	x = m.SkewX*c - s*m.ScaleY
	m.ScaleY = m.SkewX*s + m.ScaleY*c
	m.SkewX = x
	x = m.TransX*c - s*m.TransY
	m.TransY = m.TransX*s + m.TransY*c
	m.TransX = x
}

// Multiply this matrix by 'other'.
func (m *Matrix2D[T]) Multiply(other *Matrix2D[T]) {
	x := m.ScaleX*other.ScaleX + m.SkewY*other.SkewX
	m.SkewY = m.ScaleX*other.SkewY + m.SkewY*other.ScaleY
	m.ScaleX = x
	x = m.SkewX*other.ScaleX + m.ScaleY*other.SkewX
	m.ScaleY = m.SkewX*other.SkewY + m.ScaleY*other.ScaleY
	m.SkewX = x
	x = m.TransX*other.ScaleX + m.TransY*other.SkewX + other.TransX
	m.TransY = m.TransX*other.SkewY + m.TransY*other.ScaleY + other.TransY
	m.TransX = x
}

// TransformDistance returns the result of transforming the distance vector by this matrix. This is similar to
// TransformPoint(), except that the translation components of the transformation are ignored.
func (m *Matrix2D[T]) TransformDistance(distance Size[T]) Size[T] {
	x := m.ScaleX*distance.Width + m.SkewX*distance.Height
	distance.Height = m.SkewY*distance.Width + m.ScaleY*distance.Height
	distance.Width = x
	return distance
}

// TransformPoint returns the result of transforming the point by this matrix.
func (m *Matrix2D[T]) TransformPoint(where Point[T]) Point[T] {
	x := m.ScaleX*where.X + m.SkewX*where.Y + m.TransX
	where.Y = m.SkewY*where.X + m.ScaleY*where.Y + m.TransY
	where.X = x
	return where
}
