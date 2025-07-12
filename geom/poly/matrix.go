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

// Matrix provides a 2D matrix.
type Matrix struct {
	ScaleX Num
	SkewX  Num
	TransX Num
	SkewY  Num
	ScaleY Num
	TransY Num
}

// NewIdentityMatrix creates a new identity transformation Matrix.
func NewIdentityMatrix() Matrix {
	return Matrix{
		ScaleX: One,
		ScaleY: One,
	}
}

// NewTranslationMatrix creates a new Matrix that translates by 'tx' and 'ty'.
func NewTranslationMatrix(tx, ty Num) Matrix {
	return Matrix{
		ScaleX: One,
		ScaleY: One,
		TransX: tx,
		TransY: ty,
	}
}

// NewScaleMatrix creates a new Matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix(sx, sy Num) Matrix {
	return Matrix{
		ScaleX: sx,
		ScaleY: sy,
	}
}

// NewRotationMatrix creates a new Matrix that rotates by 'radians'. Positive values are clockwise.
func NewRotationMatrix(radians Num) Matrix {
	s := Sin(radians)
	c := Cos(radians)
	return Matrix{
		ScaleX: c,
		SkewX:  -s,
		SkewY:  s,
		ScaleY: c,
	}
}

// NewRotationByDegreesMatrix creates a new Matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationByDegreesMatrix(degrees Num) Matrix {
	return NewRotationMatrix(degrees.Mul(DegreesToRadians))
}

// MatrixFrom converts a geom.Matrix into a Matrix.
func MatrixFrom(m geom.Matrix) Matrix {
	return Matrix{
		ScaleX: NumFromFloat(m.ScaleX),
		SkewX:  NumFromFloat(m.SkewX),
		TransX: NumFromFloat(m.TransX),
		SkewY:  NumFromFloat(m.SkewY),
		ScaleY: NumFromFloat(m.ScaleY),
		TransY: NumFromFloat(m.TransY),
	}
}

// Matrix converts this Matrix into a geom.Matrix.
func (m Matrix) Matrix() geom.Matrix {
	return geom.Matrix{
		ScaleX: NumAsFloat[float32](m.ScaleX),
		SkewX:  NumAsFloat[float32](m.SkewX),
		TransX: NumAsFloat[float32](m.TransX),
		SkewY:  NumAsFloat[float32](m.SkewY),
		ScaleY: NumAsFloat[float32](m.ScaleY),
		TransY: NumAsFloat[float32](m.TransY),
	}
}

// Translate returns a new Matrix which is a copy of this Matrix translated by 'tx' and 'ty'.
func (m Matrix) Translate(tx, ty Num) Matrix {
	return Matrix{
		ScaleX: m.ScaleX,
		SkewX:  m.SkewX,
		TransX: m.TransX + tx,
		SkewY:  m.SkewY,
		ScaleY: m.ScaleY,
		TransY: m.TransY + ty,
	}
}

// Scale returns a new Matrix which is a copy of this Matrix scaled by 'sx' and 'sy'.
func (m Matrix) Scale(sx, sy Num) Matrix {
	return Matrix{
		ScaleX: m.ScaleX.Mul(sx),
		SkewX:  m.SkewX.Mul(sx),
		TransX: m.TransX.Mul(sx),
		SkewY:  m.SkewY.Mul(sy),
		ScaleY: m.ScaleY.Mul(sy),
		TransY: m.TransY.Mul(sy),
	}
}

// Rotate returns a new Matrix which is a copy of this Matrix rotated by 'radians'. Positive values are clockwise.
func (m Matrix) Rotate(radians Num) Matrix {
	s := Sin(radians)
	c := Cos(radians)
	return Matrix{
		ScaleX: m.ScaleX.Mul(c) - s.Mul(m.SkewY),
		SkewX:  m.SkewX.Mul(c) - s.Mul(m.ScaleY),
		TransX: m.TransX.Mul(c) - s.Mul(m.TransY),
		SkewY:  m.ScaleX.Mul(s) + m.SkewY.Mul(c),
		ScaleY: m.SkewX.Mul(s) + m.ScaleY.Mul(c),
		TransY: m.TransX.Mul(s) + m.TransY.Mul(c),
	}
}

// RotateByDegrees returns a new Matrix which is a copy of this Matrix rotated by 'degrees'. Positive values are
// clockwise.
func (m Matrix) RotateByDegrees(degrees Num) Matrix {
	return m.Rotate(degrees.Mul(DegreesToRadians))
}

// Multiply returns this Matrix multiplied by the other Matrix.
func (m Matrix) Multiply(other Matrix) Matrix {
	return Matrix{
		ScaleX: m.ScaleX.Mul(other.ScaleX) + m.SkewX.Mul(other.SkewY),
		SkewX:  m.ScaleX.Mul(other.SkewX) + m.SkewX.Mul(other.ScaleY),
		TransX: m.ScaleX.Mul(other.TransX) + m.SkewX.Mul(other.TransY) + m.TransX,
		SkewY:  m.SkewY.Mul(other.ScaleX) + m.ScaleY.Mul(other.SkewY),
		ScaleY: m.SkewY.Mul(other.SkewX) + m.ScaleY.Mul(other.ScaleY),
		TransY: m.SkewY.Mul(other.TransX) + m.ScaleY.Mul(other.TransY) + m.TransY,
	}
}

// TransformPoint returns the result of transforming the given Point by this Matrix.
func (m Matrix) TransformPoint(p Point) Point {
	return Point{
		X: m.ScaleX.Mul(p.X) + m.SkewX.Mul(p.Y) + m.TransX,
		Y: m.SkewY.Mul(p.X) + m.ScaleY.Mul(p.Y) + m.TransY,
	}
}

// String implements fmt.Stringer.
func (m Matrix) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v", m.ScaleX, m.SkewX, m.TransX, m.SkewY, m.ScaleY, m.TransY)
}
