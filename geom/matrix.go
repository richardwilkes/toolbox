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
)

// Matrix provides a 2D matrix.
type Matrix struct {
	ScaleX float32 `json:"scale_x"`
	SkewX  float32 `json:"skew_x"`
	TransX float32 `json:"trans_x"`
	SkewY  float32 `json:"skew_y"`
	ScaleY float32 `json:"scale_y"`
	TransY float32 `json:"trans_y"`
}

// NewIdentityMatrix creates a new identity transformation Matrix.
func NewIdentityMatrix() Matrix {
	return Matrix{
		ScaleX: 1,
		ScaleY: 1,
	}
}

// NewTranslationMatrix creates a new Matrix that translates by 'tx' and 'ty'.
func NewTranslationMatrix(tx, ty float32) Matrix {
	return Matrix{
		ScaleX: 1,
		ScaleY: 1,
		TransX: tx,
		TransY: ty,
	}
}

// NewScaleMatrix creates a new Matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix(sx, sy float32) Matrix {
	return Matrix{
		ScaleX: sx,
		ScaleY: sy,
	}
}

// NewRotationMatrix creates a new Matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationMatrix(degrees float32) Matrix {
	radians := degrees * xmath.DegreesToRadians
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return Matrix{
		ScaleX: c,
		SkewX:  -s,
		SkewY:  s,
		ScaleY: c,
	}
}

// NewSkewMatrix creates a new Matrix that skews by 'sx' and 'sy' degrees.
func NewSkewMatrix(sx, sy float32) Matrix {
	return Matrix{
		ScaleX: 1,
		SkewX:  xmath.Tan(sx * xmath.DegreesToRadians),
		SkewY:  xmath.Tan(sy * xmath.DegreesToRadians),
		ScaleY: 1,
	}
}

// IsIdentity returns true if this is an identity matrix.
func (m Matrix) IsIdentity() bool {
	return m.ScaleX == 1 && m.SkewX == 0 && m.TransX == 0 && m.SkewY == 0 && m.ScaleY == 1 && m.TransY == 0
}

// Translate returns a new Matrix which is a copy of this Matrix translated by 'tx' and 'ty'.
func (m Matrix) Translate(tx, ty float32) Matrix {
	return NewTranslationMatrix(tx, ty).Multiply(m)
}

// Scale returns a new Matrix which is a copy of this Matrix scaled by 'sx' and 'sy'.
func (m Matrix) Scale(sx, sy float32) Matrix {
	return NewScaleMatrix(sx, sy).Multiply(m)
}

// Skew returns a new Matrix which is a copy of this Matrix skewed by 'sx' and 'sy' degrees.
func (m Matrix) Skew(sx, sy float32) Matrix {
	return NewSkewMatrix(sx, sy).Multiply(m)
}

// Rotate returns a new Matrix which is a copy of this Matrix rotated by 'degrees'. Positive values are clockwise.
func (m Matrix) Rotate(degrees float32) Matrix {
	return NewRotationMatrix(degrees).Multiply(m)
}

// RotateAround returns a new Matrix which is a copy of this Matrix rotated by 'degrees' around the point (cx, cy).
func (m Matrix) RotateAround(degrees, cx, cy float32) Matrix {
	return m.Translate(-cx, -cy).Rotate(degrees).Translate(cx, cy)
}

// Multiply returns this Matrix multiplied by the other Matrix.
func (m Matrix) Multiply(other Matrix) Matrix {
	return Matrix{
		ScaleX: m.ScaleX*other.ScaleX + m.SkewX*other.SkewY,
		SkewX:  m.ScaleX*other.SkewX + m.SkewX*other.ScaleY,
		TransX: m.ScaleX*other.TransX + m.SkewX*other.TransY + m.TransX,
		SkewY:  m.SkewY*other.ScaleX + m.ScaleY*other.SkewY,
		ScaleY: m.SkewY*other.SkewX + m.ScaleY*other.ScaleY,
		TransY: m.SkewY*other.TransX + m.ScaleY*other.TransY + m.TransY,
	}
}

// TransformPoint returns the result of transforming the given Point by this Matrix.
func (m Matrix) TransformPoint(p Point) Point {
	return Point{
		X: m.ScaleX*p.X + m.SkewX*p.Y + m.TransX,
		Y: m.SkewY*p.X + m.ScaleY*p.Y + m.TransY,
	}
}

// Invert returns the inverse of this Matrix. If the Matrix is non-invertible, an identity Matrix is returned.
func (m Matrix) Invert() Matrix {
	det := m.ScaleX*m.ScaleY - m.SkewX*m.SkewY
	if det == 0 {
		// Non-invertible matrix; return identity
		return NewIdentityMatrix()
	}
	invDet := 1 / det
	return Matrix{
		ScaleX: m.ScaleY * invDet,
		SkewX:  -m.SkewX * invDet,
		TransX: (m.SkewX*m.TransY - m.ScaleY*m.TransX) * invDet,
		SkewY:  -m.SkewY * invDet,
		ScaleY: m.ScaleX * invDet,
		TransY: (m.SkewY*m.TransX - m.ScaleX*m.TransY) * invDet,
	}
}

// String implements fmt.Stringer.
func (m Matrix) String() string {
	return fmt.Sprintf("%#v,%#v,%#v,%#v,%#v,%#v", m.ScaleX, m.SkewX, m.TransX, m.SkewY, m.ScaleY, m.TransY)
}
