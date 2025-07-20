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

// NewTranslationMatrixPt creates a new Matrix that translates by translation.X and translation.Y.
func NewTranslationMatrixPt(translation Point) Matrix {
	return NewTranslationMatrix(translation.X, translation.Y)
}

// NewScaleMatrix creates a new Matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix(sx, sy float32) Matrix {
	return Matrix{
		ScaleX: sx,
		ScaleY: sy,
	}
}

// NewScaleMatrixPt creates a new Matrix that scales by scale.X and scale.Y.
func NewScaleMatrixPt(scale Point) Matrix {
	return Matrix{
		ScaleX: scale.X,
		ScaleY: scale.Y,
	}
}

// NewRotationMatrix creates a new Matrix that rotates by 'radians'. Positive values are clockwise.
func NewRotationMatrix(radians float32) Matrix {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return Matrix{
		ScaleX: c,
		SkewX:  -s,
		SkewY:  s,
		ScaleY: c,
	}
}

// NewRotationByDegreesMatrix creates a new Matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationByDegreesMatrix(degrees float32) Matrix {
	return NewRotationMatrix(degrees * xmath.DegreesToRadians)
}

// Translate returns a new Matrix which is a copy of this Matrix translated by 'tx' and 'ty'.
func (m Matrix) Translate(tx, ty float32) Matrix {
	return Matrix{
		ScaleX: m.ScaleX,
		SkewX:  m.SkewX,
		TransX: m.TransX + tx,
		SkewY:  m.SkewY,
		ScaleY: m.ScaleY,
		TransY: m.TransY + ty,
	}
}

// TranslatePt returns a new Matrix which is a copy of this Matrix translated by translate.X and translate.Y.
func (m Matrix) TranslatePt(translate Point) Matrix {
	return m.Translate(translate.X, translate.Y)
}

// Scale returns a new Matrix which is a copy of this Matrix scaled by 'sx' and 'sy'.
func (m Matrix) Scale(sx, sy float32) Matrix {
	return Matrix{
		ScaleX: m.ScaleX * sx,
		SkewX:  m.SkewX * sx,
		TransX: m.TransX * sx,
		SkewY:  m.SkewY * sy,
		ScaleY: m.ScaleY * sy,
		TransY: m.TransY * sy,
	}
}

// ScalePt returns a new Matrix which is a copy of this Matrix scaled by scale.X and scale.Y.
func (m Matrix) ScalePt(scale Point) Matrix {
	return m.Scale(scale.X, scale.Y)
}

// Rotate returns a new Matrix which is a copy of this Matrix rotated by 'radians'. Positive values are clockwise.
func (m Matrix) Rotate(radians float32) Matrix {
	s := xmath.Sin(radians)
	c := xmath.Cos(radians)
	return Matrix{
		ScaleX: m.ScaleX*c - s*m.SkewY,
		SkewX:  m.SkewX*c - s*m.ScaleY,
		TransX: m.TransX*c - s*m.TransY,
		SkewY:  m.ScaleX*s + m.SkewY*c,
		ScaleY: m.SkewX*s + m.ScaleY*c,
		TransY: m.TransX*s + m.TransY*c,
	}
}

// RotateByDegrees returns a new Matrix which is a copy of this Matrix rotated by 'degrees'. Positive values are clockwise.
func (m Matrix) RotateByDegrees(degrees float32) Matrix {
	return m.Rotate(degrees * xmath.DegreesToRadians)
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

// String implements fmt.Stringer.
func (m Matrix) String() string {
	return fmt.Sprintf("%#v,%#v,%#v,%#v,%#v,%#v", m.ScaleX, m.SkewX, m.TransX, m.SkewY, m.ScaleY, m.TransY)
}
