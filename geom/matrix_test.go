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
)

func TestNewIdentityMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix[float64]()
	c.Equal(1.0, m.ScaleX)
	c.Equal(0.0, m.SkewX)
	c.Equal(0.0, m.TransX)
	c.Equal(0.0, m.SkewY)
	c.Equal(1.0, m.ScaleY)
	c.Equal(0.0, m.TransY)
}

func TestNewTranslationMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix[float64](10, 20)
	c.Equal(1.0, m.ScaleX)
	c.Equal(0.0, m.SkewX)
	c.Equal(10.0, m.TransX)
	c.Equal(0.0, m.SkewY)
	c.Equal(1.0, m.ScaleY)
	c.Equal(20.0, m.TransY)
}

func TestNewScaleMatrix(t *testing.T) {
	c := check.New(t)

	m := geom.NewScaleMatrix[float64](2, 3)
	c.Equal(2.0, m.ScaleX)
	c.Equal(0.0, m.SkewX)
	c.Equal(0.0, m.TransX)
	c.Equal(0.0, m.SkewY)
	c.Equal(3.0, m.ScaleY)
	c.Equal(0.0, m.TransY)
}

func TestNewRotationMatrix(t *testing.T) {
	c := check.New(t)

	// Test 90 degree rotation (π/2 radians)
	m := geom.NewRotationMatrix(math.Pi / 2)

	// cos(π/2) = 0, sin(π/2) = 1
	// For clockwise rotation: ScaleX = cos, SkewX = -sin, SkewY = sin, ScaleY = cos
	c.True(math.Abs(m.ScaleX) < 0.0001)       // Should be ~0
	c.True(math.Abs(m.SkewX-(-1.0)) < 0.0001) // Should be ~-1
	c.Equal(0.0, m.TransX)
	c.True(math.Abs(m.SkewY-1.0) < 0.0001) // Should be ~1
	c.True(math.Abs(m.ScaleY) < 0.0001)    // Should be ~0
	c.Equal(0.0, m.TransY)
}

func TestNewRotationByDegreesMatrix(t *testing.T) {
	c := check.New(t)

	// Test 90 degree rotation
	m := geom.NewRotationByDegreesMatrix[float64](90)

	c.True(math.Abs(m.ScaleX) < 0.0001)       // Should be ~0
	c.True(math.Abs(m.SkewX-(-1.0)) < 0.0001) // Should be ~-1
	c.Equal(0.0, m.TransX)
	c.True(math.Abs(m.SkewY-1.0) < 0.0001) // Should be ~1
	c.True(math.Abs(m.ScaleY) < 0.0001)    // Should be ~0
	c.Equal(0.0, m.TransY)
}

func TestMatrixTranslate(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix[float64]()
	translated := m.Translate(10, 20)

	c.Equal(1.0, translated.ScaleX)
	c.Equal(0.0, translated.SkewX)
	c.Equal(10.0, translated.TransX)
	c.Equal(0.0, translated.SkewY)
	c.Equal(1.0, translated.ScaleY)
	c.Equal(20.0, translated.TransY)

	// Original matrix should be unchanged
	c.Equal(0.0, m.TransX)
	c.Equal(0.0, m.TransY)

	// Test cumulative translation
	translated2 := translated.Translate(5, 7)
	c.Equal(15.0, translated2.TransX)
	c.Equal(27.0, translated2.TransY)
}

func TestMatrixScale(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix[float64](10, 20)
	scaled := m.Scale(2, 3)

	c.Equal(2.0, scaled.ScaleX)
	c.Equal(0.0, scaled.SkewX)
	c.Equal(20.0, scaled.TransX) // Translation is also scaled
	c.Equal(0.0, scaled.SkewY)
	c.Equal(3.0, scaled.ScaleY)
	c.Equal(60.0, scaled.TransY) // Translation is also scaled

	// Original matrix should be unchanged
	c.Equal(1.0, m.ScaleX)
	c.Equal(1.0, m.ScaleY)
	c.Equal(10.0, m.TransX)
	c.Equal(20.0, m.TransY)
}

func TestMatrixRotate(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix[float64]()
	rotated := m.Rotate(math.Pi / 2) // 90 degrees

	c.True(math.Abs(rotated.ScaleX) < 0.0001)       // Should be ~0
	c.True(math.Abs(rotated.SkewX-(-1.0)) < 0.0001) // Should be ~-1
	c.Equal(0.0, rotated.TransX)
	c.True(math.Abs(rotated.SkewY-1.0) < 0.0001) // Should be ~1
	c.True(math.Abs(rotated.ScaleY) < 0.0001)    // Should be ~0
	c.Equal(0.0, rotated.TransY)

	// Original matrix should be unchanged
	c.Equal(1.0, m.ScaleX)
	c.Equal(1.0, m.ScaleY)
}

func TestMatrixRotateByDegrees(t *testing.T) {
	c := check.New(t)

	m := geom.NewIdentityMatrix[float64]()
	rotated := m.RotateByDegrees(90)

	c.True(math.Abs(rotated.ScaleX) < 0.0001)       // Should be ~0
	c.True(math.Abs(rotated.SkewX-(-1.0)) < 0.0001) // Should be ~-1
	c.Equal(0.0, rotated.TransX)
	c.True(math.Abs(rotated.SkewY-1.0) < 0.0001) // Should be ~1
	c.True(math.Abs(rotated.ScaleY) < 0.0001)    // Should be ~0
	c.Equal(0.0, rotated.TransY)
}

func TestMatrixMultiply(t *testing.T) {
	c := check.New(t)

	// Test identity multiplication
	identity := geom.NewIdentityMatrix[float64]()
	translation := geom.NewTranslationMatrix[float64](10, 20)

	result1 := identity.Multiply(translation)
	c.Equal(translation.ScaleX, result1.ScaleX)
	c.Equal(translation.SkewX, result1.SkewX)
	c.Equal(translation.TransX, result1.TransX)
	c.Equal(translation.SkewY, result1.SkewY)
	c.Equal(translation.ScaleY, result1.ScaleY)
	c.Equal(translation.TransY, result1.TransY)

	// Test scale and translation combination
	scale := geom.NewScaleMatrix[float64](2, 3)
	combined := scale.Multiply(translation)

	c.Equal(2.0, combined.ScaleX)
	c.Equal(0.0, combined.SkewX)
	c.Equal(20.0, combined.TransX) // Translation after scale
	c.Equal(0.0, combined.SkewY)
	c.Equal(3.0, combined.ScaleY)
	c.Equal(60.0, combined.TransY) // Translation after scale
}

//nolint:gocritic // The "commented out code" is actually explanation
func TestMatrixTransformPoint(t *testing.T) {
	c := check.New(t)

	// Test identity transformation
	identity := geom.NewIdentityMatrix[float64]()
	p := geom.NewPoint(5.0, 7.0)
	result := identity.TransformPoint(p)

	c.Equal(5.0, result.X)
	c.Equal(7.0, result.Y)

	// Test translation
	translation := geom.NewTranslationMatrix[float64](10, 20)
	result2 := translation.TransformPoint(p)

	c.Equal(15.0, result2.X) // 5 + 10
	c.Equal(27.0, result2.Y) // 7 + 20

	// Test scaling
	scale := geom.NewScaleMatrix[float64](2, 3)
	result3 := scale.TransformPoint(p)

	c.Equal(10.0, result3.X) // 5 * 2
	c.Equal(21.0, result3.Y) // 7 * 3

	// Test 90-degree rotation (point (1,0) should become (0,1))
	rotation := geom.NewRotationByDegreesMatrix[float64](90)
	p2 := geom.NewPoint(1.0, 0.0)
	result4 := rotation.TransformPoint(p2)

	c.True(math.Abs(result4.X) < 0.0001)     // Should be ~0
	c.True(math.Abs(result4.Y-1.0) < 0.0001) // Should be ~1

	// Test combined transformations
	combined := geom.NewTranslationMatrix[float64](5, 5).Scale(2, 2)
	p3 := geom.NewPoint(3.0, 4.0)
	result5 := combined.TransformPoint(p3)

	c.Equal(16.0, result5.X) // (3 + 5) * 2 = 16
	c.Equal(18.0, result5.Y) // (4 + 5) * 2 = 18
}

func TestMatrixString(t *testing.T) {
	c := check.New(t)

	m := geom.NewTranslationMatrix[float64](10, 20)
	str := m.String()
	c.Equal("1,0,10,0,1,20", str)

	// Test with decimals
	m2 := geom.NewScaleMatrix[float64](1.5, 2.5)
	str2 := m2.String()
	c.Equal("1.5,0,0,0,2.5,0", str2)
}
