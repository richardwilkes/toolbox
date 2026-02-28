// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
)

// Insets defines margins on each side of a rectangle.
type Insets struct {
	Top    float32 `json:"top"`
	Left   float32 `json:"left"`
	Bottom float32 `json:"bottom"`
	Right  float32 `json:"right"`
}

// NewInsets returns an Insets with the given values for its edges.
func NewInsets(top, left, bottom, right float32) Insets {
	return Insets{
		Top:    top,
		Left:   left,
		Bottom: bottom,
		Right:  right,
	}
}

// NewUniformInsets returns an Insets whose edges all have the same value.
func NewUniformInsets(amount float32) Insets {
	return Insets{
		Top:    amount,
		Left:   amount,
		Bottom: amount,
		Right:  amount,
	}
}

// NewSymmetricInsets returns an Insets whose edges match their opposite edge.
func NewSymmetricInsets(h, v float32) Insets {
	return Insets{
		Top:    v,
		Left:   h,
		Bottom: v,
		Right:  h,
	}
}

// NewHorizontalInsets returns an Insets whose left and right edges have the specified value.
func NewHorizontalInsets(amount float32) Insets {
	return Insets{
		Left:  amount,
		Right: amount,
	}
}

// NewVerticalInsets returns an Insets whose top and bottom edges have the specified value.
func NewVerticalInsets(amount float32) Insets {
	return Insets{
		Top:    amount,
		Bottom: amount,
	}
}

// Add returns a new Insets which is the result of adding this Insets with the provided Insets.
func (i Insets) Add(in Insets) Insets {
	return Insets{
		Top:    i.Top + in.Top,
		Left:   i.Left + in.Left,
		Bottom: i.Bottom + in.Bottom,
		Right:  i.Right + in.Right,
	}
}

// Sub returns a new Insets which is the result of subtracting the provided Insets from this Insets.
func (i Insets) Sub(in Insets) Insets {
	return Insets{
		Top:    i.Top - in.Top,
		Left:   i.Left - in.Left,
		Bottom: i.Bottom - in.Bottom,
		Right:  i.Right - in.Right,
	}
}

// Mul returns a new Insets which is the result of multiplying the values of this Insets by the value.
func (i Insets) Mul(value float32) Insets {
	return Insets{
		Top:    i.Top * value,
		Left:   i.Left * value,
		Bottom: i.Bottom * value,
		Right:  i.Right * value,
	}
}

// Div returns a new Insets which is the result of dividing the values of this Insets by the value.
func (i Insets) Div(value float32) Insets {
	return Insets{
		Top:    i.Top / value,
		Left:   i.Left / value,
		Bottom: i.Bottom / value,
		Right:  i.Right / value,
	}
}

// Size returns the Size of the Insets.
func (i Insets) Size() Size {
	return NewSize(i.Width(), i.Height())
}

// Width returns the sum of the left and right insets.
func (i Insets) Width() float32 {
	return i.Left + i.Right
}

// Height returns the sum of the top and bottom insets.
func (i Insets) Height() float32 {
	return i.Top + i.Bottom
}

// String implements fmt.Stringer.
func (i Insets) String() string {
	return fmt.Sprintf("%#v,%#v,%#v,%#v", i.Top, i.Left, i.Bottom, i.Right)
}
