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

// NewUniformInsets creates a new Insets whose edges all have the same value.
func NewUniformInsets[T xmath.Numeric](amount T) Insets[T] {
	return Insets[T]{Top: amount, Left: amount, Bottom: amount, Right: amount}
}

// NewHorizontalInsets creates a new Insets whose left and right edges have the specified value.
func NewHorizontalInsets[T xmath.Numeric](amount T) Insets[T] {
	return Insets[T]{Left: amount, Right: amount}
}

// NewVerticalInsets creates a new Insets whose top and bottom edges have the specified value.
func NewVerticalInsets[T xmath.Numeric](amount T) Insets[T] {
	return Insets[T]{Top: amount, Bottom: amount}
}

// Add modifies this Insets by adding the supplied Insets. Returns itself for easy chaining.
func (i *Insets[T]) Add(insets Insets[T]) *Insets[T] {
	i.Top += insets.Top
	i.Left += insets.Left
	i.Bottom += insets.Bottom
	i.Right += insets.Right
	return i
}

// Subtract modifies this Insets by subtracting the supplied Insets. Returns itself for easy chaining.
func (i *Insets[T]) Subtract(insets Insets[T]) *Insets[T] {
	i.Top -= insets.Top
	i.Left -= insets.Left
	i.Bottom -= insets.Bottom
	i.Right -= insets.Right
	return i
}

// String implements the fmt.Stringer interface.
func (i *Insets[T]) String() string {
	return fmt.Sprintf("%v,%v,%v,%v", i.Top, i.Left, i.Bottom, i.Right)
}
