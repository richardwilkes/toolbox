// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/xmath"
)

// Insets defines margins on each side of a rectangle.
type Insets[T xmath.Numeric] struct {
	Top    T `json:"top"`
	Left   T `json:"left"`
	Bottom T `json:"bottom"`
	Right  T `json:"right"`
}

// NewInsets returns an Insets with the given values for its edges.
func NewInsets[T xmath.Numeric](top, left, bottom, right T) Insets[T] {
	return Insets[T]{Top: top, Left: left, Bottom: bottom, Right: right}
}

// NewUniformInsets returns an Insets whose edges all have the same value.
func NewUniformInsets[T xmath.Numeric](amount T) Insets[T] {
	return NewInsets(amount, amount, amount, amount)
}

// NewSymmetricInsets returns an Insets whose edges match their opposite edge.
func NewSymmetricInsets[T xmath.Numeric](h, v T) Insets[T] {
	return NewInsets(v, h, v, h)
}

// NewHorizontalInsets returns an Insets whose left and right edges have the specified value.
func NewHorizontalInsets[T xmath.Numeric](amount T) Insets[T] {
	return Insets[T]{Left: amount, Right: amount}
}

// NewVerticalInsets returns an Insets whose top and bottom edges have the specified value.
func NewVerticalInsets[T xmath.Numeric](amount T) Insets[T] {
	return Insets[T]{Top: amount, Bottom: amount}
}

// ConvertInsets converts a Insets of type F into one of type T.
func ConvertInsets[T, F xmath.Numeric](i Insets[F]) Insets[T] {
	return NewInsets(T(i.Top), T(i.Left), T(i.Bottom), T(i.Right))
}

// Add returns a new Insets which is the result of adding this Insets with the provided Insets.
func (i Insets[T]) Add(in Insets[T]) Insets[T] {
	return NewInsets(i.Top+in.Top, i.Left+in.Left, i.Bottom+in.Bottom, i.Right+in.Right)
}

// Sub returns a new Insets which is the result of subtracting the provided Insets from this Insets.
func (i Insets[T]) Sub(in Insets[T]) Insets[T] {
	return NewInsets(i.Top-in.Top, i.Left-in.Left, i.Bottom-in.Bottom, i.Right-in.Right)
}

// Mul returns a new Insets which is the result of multiplying the values of this Insets by the value.
func (i Insets[T]) Mul(value T) Insets[T] {
	return NewInsets(i.Top*value, i.Left*value, i.Bottom*value, i.Right*value)
}

// Div returns a new Insets which is the result of dividing the values of this Insets by the value.
func (i Insets[T]) Div(value T) Insets[T] {
	return NewInsets(i.Top/value, i.Left/value, i.Bottom/value, i.Right/value)
}

// Size returns the Size of the Insets.
func (i Insets[T]) Size() Size[T] {
	return NewSize(i.Width(), i.Height())
}

// Width returns the sum of the left and right insets.
func (i Insets[T]) Width() T {
	return i.Left + i.Right
}

// Height returns the sum of the top and bottom insets.
func (i Insets[T]) Height() T {
	return i.Top + i.Bottom
}

// String implements fmt.Stringer.
func (i Insets[T]) String() string {
	return fmt.Sprintf("%#v,%#v,%#v,%#v", i.Top, i.Left, i.Bottom, i.Right)
}
