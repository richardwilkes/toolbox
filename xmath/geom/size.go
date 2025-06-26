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

	"github.com/richardwilkes/toolbox/xmath"
)

// Size defines a width and height.
type Size[T xmath.Numeric] struct {
	Width  T `json:"w"`
	Height T `json:"h"`
}

// NewSize creates a new Size.
func NewSize[T xmath.Numeric](width, height T) Size[T] {
	return Size[T]{
		Width:  width,
		Height: height,
	}
}

// ConvertSize converts a Size of type F into one of type T.
func ConvertSize[T, F xmath.Numeric](s Size[F]) Size[T] {
	return NewSize(T(s.Width), T(s.Height))
}

// Add returns a new Size which is the result of adding this Size with the provided Size.
func (s Size[T]) Add(size Size[T]) Size[T] {
	return Size[T]{Width: s.Width + size.Width, Height: s.Height + size.Height}
}

// Sub returns a new Size which is the result of subtracting the provided Size from this Size.
func (s Size[T]) Sub(size Size[T]) Size[T] {
	return Size[T]{Width: s.Width - size.Width, Height: s.Height - size.Height}
}

// Mul returns a new Size which is the result of multiplying this Size by the value.
func (s Size[T]) Mul(value T) Size[T] {
	return Size[T]{Width: s.Width * value, Height: s.Height * value}
}

// Div returns a new Size which is the result of dividing this Size by the value.
func (s Size[T]) Div(value T) Size[T] {
	return Size[T]{Width: s.Width / value, Height: s.Height / value}
}

// Floor returns a new Size with its width and height floored.
func (s Size[T]) Floor() Size[T] {
	return Size[T]{Width: xmath.Floor(s.Width), Height: xmath.Floor(s.Height)}
}

// Ceil returns a new Size with its width and height ceiled.
func (s Size[T]) Ceil() Size[T] {
	return Size[T]{Width: xmath.Ceil(s.Width), Height: xmath.Ceil(s.Height)}
}

// Min returns the smallest Size between itself and 'other'.
func (s Size[T]) Min(other Size[T]) Size[T] {
	return Size[T]{Width: min(s.Width, other.Width), Height: min(s.Height, other.Height)}
}

// Max returns the largest Size between itself and 'other'.
func (s Size[T]) Max(other Size[T]) Size[T] {
	return Size[T]{Width: max(s.Width, other.Width), Height: max(s.Height, other.Height)}
}

// ConstrainForHint returns a size no larger than the hint value. Hint values less than one are ignored.
func (s Size[T]) ConstrainForHint(hint Size[T]) Size[T] {
	w := s.Width
	if hint.Width >= 1 && w > hint.Width {
		w = hint.Width
	}
	h := s.Height
	if hint.Height >= 1 && h > hint.Height {
		h = hint.Height
	}
	return Size[T]{Width: w, Height: h}
}

// String implements fmt.Stringer.
func (s Size[T]) String() string {
	return fmt.Sprintf("%#v,%#v", s.Width, s.Height)
}
