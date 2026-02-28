// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/geom"
)

func TestNewInsets(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(1, 2, 3, 4) // top, left, bottom, right
	c.Equal(float32(1), insets.Top)
	c.Equal(float32(2), insets.Left)
	c.Equal(float32(3), insets.Bottom)
	c.Equal(float32(4), insets.Right)
}

func TestNewUniformInsets(t *testing.T) {
	c := check.New(t)

	insets := geom.NewUniformInsets(5)
	c.Equal(float32(5), insets.Top)
	c.Equal(float32(5), insets.Left)
	c.Equal(float32(5), insets.Bottom)
	c.Equal(float32(5), insets.Right)
}

func TestNewSymmetricInsets(t *testing.T) {
	c := check.New(t)

	insets := geom.NewSymmetricInsets(3, 7) // horizontal, vertical
	c.Equal(float32(7), insets.Top)         // vertical
	c.Equal(float32(3), insets.Left)        // horizontal
	c.Equal(float32(7), insets.Bottom)      // vertical
	c.Equal(float32(3), insets.Right)       // horizontal
}

func TestNewHorizontalInsets(t *testing.T) {
	c := check.New(t)

	insets := geom.NewHorizontalInsets(4)
	c.Equal(float32(0), insets.Top)
	c.Equal(float32(4), insets.Left)
	c.Equal(float32(0), insets.Bottom)
	c.Equal(float32(4), insets.Right)
}

func TestNewVerticalInsets(t *testing.T) {
	c := check.New(t)

	insets := geom.NewVerticalInsets(6)
	c.Equal(float32(6), insets.Top)
	c.Equal(float32(0), insets.Left)
	c.Equal(float32(6), insets.Bottom)
	c.Equal(float32(0), insets.Right)
}

func TestInsetsAdd(t *testing.T) {
	c := check.New(t)

	i1 := geom.NewInsets(1, 2, 3, 4)
	i2 := geom.NewInsets(5, 6, 7, 8)
	result := i1.Add(i2)

	c.Equal(float32(6), result.Top)     // 1 + 5
	c.Equal(float32(8), result.Left)    // 2 + 6
	c.Equal(float32(10), result.Bottom) // 3 + 7
	c.Equal(float32(12), result.Right)  // 4 + 8

	// Original insets should be unchanged
	c.Equal(float32(1), i1.Top)
	c.Equal(float32(2), i1.Left)
	c.Equal(float32(3), i1.Bottom)
	c.Equal(float32(4), i1.Right)
}

func TestInsetsSub(t *testing.T) {
	c := check.New(t)

	i1 := geom.NewInsets(10, 12, 13, 14)
	i2 := geom.NewInsets(1, 2, 3, 4)
	result := i1.Sub(i2)

	c.Equal(float32(9), result.Top)     // 10 - 1
	c.Equal(float32(10), result.Left)   // 12 - 2
	c.Equal(float32(10), result.Bottom) // 13 - 3
	c.Equal(float32(10), result.Right)  // 14 - 4

	// Original insets should be unchanged
	c.Equal(float32(10), i1.Top)
	c.Equal(float32(12), i1.Left)
	c.Equal(float32(13), i1.Bottom)
	c.Equal(float32(14), i1.Right)
}

func TestInsetsMul(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(2, 3, 4, 5)
	result := insets.Mul(3)

	c.Equal(float32(6), result.Top)     // 2 * 3
	c.Equal(float32(9), result.Left)    // 3 * 3
	c.Equal(float32(12), result.Bottom) // 4 * 3
	c.Equal(float32(15), result.Right)  // 5 * 3

	// Original insets should be unchanged
	c.Equal(float32(2), insets.Top)
	c.Equal(float32(3), insets.Left)
	c.Equal(float32(4), insets.Bottom)
	c.Equal(float32(5), insets.Right)
}

func TestInsetsDiv(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(6, 9, 12, 15)
	result := insets.Div(3)

	c.Equal(float32(2), result.Top)    // 6 / 3
	c.Equal(float32(3), result.Left)   // 9 / 3
	c.Equal(float32(4), result.Bottom) // 12 / 3
	c.Equal(float32(5), result.Right)  // 15 / 3

	// Test with float
	insetsF := geom.NewInsets(7.0, 9.0, 11.0, 13.0)
	resultF := insetsF.Div(2.0)

	c.Equal(float32(3.5), resultF.Top)
	c.Equal(float32(4.5), resultF.Left)
	c.Equal(float32(5.5), resultF.Bottom)
	c.Equal(float32(6.5), resultF.Right)
}

func TestInsetsSize(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(2, 3, 4, 5)
	size := insets.Size()

	// Width: left + right = 3 + 5 = 8
	// Height: top + bottom = 2 + 4 = 6
	c.Equal(float32(8), size.Width)
	c.Equal(float32(6), size.Height)
}

func TestInsetsWidth(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(2, 3, 4, 5)
	width := insets.Width()

	// Width: left + right = 3 + 5 = 8
	c.Equal(float32(8), width)
}

func TestInsetsHeight(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(2, 3, 4, 5)
	height := insets.Height()

	// Height: top + bottom = 2 + 4 = 6
	c.Equal(float32(6), height)
}

func TestInsetsString(t *testing.T) {
	c := check.New(t)

	insets := geom.NewInsets(1, 2, 3, 4)
	str := insets.String()
	c.Equal("1,2,3,4", str)

	// Test with float
	insetsF := geom.NewInsets(1.5, 2.7, 3.2, 4.8)
	strF := insetsF.String()
	c.Equal("1.5,2.7,3.2,4.8", strF)
}
