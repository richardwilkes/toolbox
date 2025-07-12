// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom_test

import (
	"math"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/xmath"
)

func TestNewIdentityMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix()
	c.Equal(float32(1), m.ScaleX)
	c.Equal(float32(0), m.SkewX)
	c.Equal(float32(0), m.TransX)
	c.Equal(float32(0), m.SkewY)
	c.Equal(float32(1), m.ScaleY)
	c.Equal(float32(0), m.TransY)
}

func TestNewTranslationMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix(10, 20)
	c.Equal(float32(1), m.ScaleX)
	c.Equal(float32(0), m.SkewX)
	c.Equal(float32(10), m.TransX)
	c.Equal(float32(0), m.SkewY)
	c.Equal(float32(1), m.ScaleY)
	c.Equal(float32(20), m.TransY)
}

func TestNewScaleMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewScaleMatrix(2, 3)
	c.Equal(float32(2), m.ScaleX)
	c.Equal(float32(0), m.SkewX)
	c.Equal(float32(0), m.TransX)
	c.Equal(float32(0), m.SkewY)
	c.Equal(float32(3), m.ScaleY)
	c.Equal(float32(0), m.TransY)
}

func TestNewRotationMatrix(t *testing.T) {
	c := check.New(t)

	// Test 90 degree rotation (π/2 radians)
	m := geom.NewRotationMatrix(math.Pi / 2)

	// cos(π/2) = 0, sin(π/2) = 1
	// For clockwise rotation: ScaleX = cos, SkewX = -sin, SkewY = sin, ScaleY = cos
	c.True(xmath.Abs(m.ScaleX) < 0.0001)  // Should be ~0
	c.True(xmath.Abs(m.SkewX+1) < 0.0001) // Should be ~-1
	c.Equal(float32(0), m.TransX)
	c.True(xmath.Abs(m.SkewY-1) < 0.0001) // Should be ~1
	c.True(xmath.Abs(m.ScaleY) < 0.0001)  // Should be ~0
	c.Equal(float32(0), m.TransY)
}

func TestNewRotationByDegreesMatrix(t *testing.T) {
	c := check.New(t)

	// Test 90 degree rotation
	m := geom.NewRotationByDegreesMatrix(90)

	c.True(xmath.Abs(m.ScaleX) < 0.0001)  // Should be ~0
	c.True(xmath.Abs(m.SkewX+1) < 0.0001) // Should be ~-1
	c.Equal(float32(0), m.TransX)
	c.True(xmath.Abs(m.SkewY-1) < 0.0001) // Should be ~1
	c.True(xmath.Abs(m.ScaleY) < 0.0001)  // Should be ~0
	c.Equal(float32(0), m.TransY)
}

func TestMatrixTranslate(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix()
	translated := m.Translate(10, 20)

	c.Equal(float32(1), translated.ScaleX)
	c.Equal(float32(0), translated.SkewX)
	c.Equal(float32(10), translated.TransX)
	c.Equal(float32(0), translated.SkewY)
	c.Equal(float32(1), translated.ScaleY)
	c.Equal(float32(20), translated.TransY)

	// Original matrix should be unchanged
	c.Equal(float32(0), m.TransX)
	c.Equal(float32(0), m.TransY)

	// Test cumulative translation
	translated2 := translated.Translate(5, 7)
	c.Equal(float32(15), translated2.TransX)
	c.Equal(float32(27), translated2.TransY)
}

func TestMatrixScale(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix(10, 20)
	scaled := m.Scale(2, 3)

	c.Equal(float32(2), scaled.ScaleX)
	c.Equal(float32(0), scaled.SkewX)
	c.Equal(float32(20), scaled.TransX) // Translation is also scaled
	c.Equal(float32(0), scaled.SkewY)
	c.Equal(float32(3), scaled.ScaleY)
	c.Equal(float32(60), scaled.TransY) // Translation is also scaled

	// Original matrix should be unchanged
	c.Equal(float32(1), m.ScaleX)
	c.Equal(float32(1), m.ScaleY)
	c.Equal(float32(10), m.TransX)
	c.Equal(float32(20), m.TransY)
}

func TestMatrixRotate(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix()
	rotated := m.Rotate(math.Pi / 2) // 90 degrees

	c.True(xmath.Abs(rotated.ScaleX) < 0.0001)  // Should be ~0
	c.True(xmath.Abs(rotated.SkewX+1) < 0.0001) // Should be ~-1
	c.Equal(float32(0), rotated.TransX)
	c.True(xmath.Abs(rotated.SkewY-1) < 0.0001) // Should be ~1
	c.True(xmath.Abs(rotated.ScaleY) < 0.0001)  // Should be ~0
	c.Equal(float32(0), rotated.TransY)

	// Original matrix should be unchanged
	c.Equal(float32(1), m.ScaleX)
	c.Equal(float32(1), m.ScaleY)
}

func TestMatrixRotateByDegrees(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix()
	rotated := m.RotateByDegrees(90)

	c.True(xmath.Abs(rotated.ScaleX) < 0.0001)  // Should be ~0
	c.True(xmath.Abs(rotated.SkewX+1) < 0.0001) // Should be ~-1
	c.Equal(float32(0), rotated.TransX)
	c.True(xmath.Abs(rotated.SkewY-1) < 0.0001) // Should be ~1
	c.True(xmath.Abs(rotated.ScaleY) < 0.0001)  // Should be ~0
	c.Equal(float32(0), rotated.TransY)
}

func TestMatrixMultiply(t *testing.T) {
	c := check.New(t)

	// Test identity multiplication
	identity := geom.NewIdentityMatrix()
	translation := geom.NewTranslationMatrix(10, 20)

	result1 := identity.Multiply(translation)
	c.Equal(translation.ScaleX, result1.ScaleX)
	c.Equal(translation.SkewX, result1.SkewX)
	c.Equal(translation.TransX, result1.TransX)
	c.Equal(translation.SkewY, result1.SkewY)
	c.Equal(translation.ScaleY, result1.ScaleY)
	c.Equal(translation.TransY, result1.TransY)

	// Test scale and translation combination
	scale := geom.NewScaleMatrix(2, 3)
	combined := scale.Multiply(translation)

	c.Equal(float32(2), combined.ScaleX)
	c.Equal(float32(0), combined.SkewX)
	c.Equal(float32(20), combined.TransX) // Translation after scale
	c.Equal(float32(0), combined.SkewY)
	c.Equal(float32(3), combined.ScaleY)
	c.Equal(float32(60), combined.TransY) // Translation after scale
}

//nolint:gocritic // The "commented out code" is actually explanation
func TestMatrixTransformPoint(t *testing.T) {
	c := check.New(t)

	// Test identity transformation
	identity := geom.NewIdentityMatrix()
	p := geom.NewPoint(5.0, 7.0)
	result := identity.TransformPoint(p)

	c.Equal(float32(5), result.X)
	c.Equal(float32(7), result.Y)

	// Test translation
	translation := geom.NewTranslationMatrix(10, 20)
	result2 := translation.TransformPoint(p)

	c.Equal(float32(15), result2.X) // 5 + 10
	c.Equal(float32(27), result2.Y) // 7 + 20

	// Test scaling
	scale := geom.NewScaleMatrix(2, 3)
	result3 := scale.TransformPoint(p)

	c.Equal(float32(10), result3.X) // 5 * 2
	c.Equal(float32(21), result3.Y) // 7 * 3

	// Test 90-degree rotation (point (1,0) should become (0,1))
	rotation := geom.NewRotationByDegreesMatrix(90)
	p2 := geom.NewPoint(1.0, 0.0)
	result4 := rotation.TransformPoint(p2)

	c.True(xmath.Abs(result4.X) < 0.0001)     // Should be ~0
	c.True(xmath.Abs(result4.Y-1.0) < 0.0001) // Should be ~1

	// Test combined transformations
	combined := geom.NewTranslationMatrix(5, 5).Scale(2, 2)
	p3 := geom.NewPoint(3.0, 4.0)
	result5 := combined.TransformPoint(p3)

	c.Equal(float32(16), result5.X) // (3 + 5) * 2 = 16
	c.Equal(float32(18), result5.Y) // (4 + 5) * 2 = 18
}

func TestMatrixString(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix(10, 20)
	str := m.String()
	c.Equal("1,0,10,0,1,20", str)

	// Test with decimals
	m2 := geom.NewScaleMatrix(1.5, 2.5)
	str2 := m2.String()
	c.Equal("1.5,0,0,0,2.5,0", str2)
}
