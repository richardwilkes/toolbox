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
	"math"
	"reflect"

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

// NewSizePtr creates a new *Size.
func NewSizePtr[T xmath.Numeric](x, y T) *Size[T] {
	s := NewSize[T](x, y)
	return &s
}

// Add modifies this Size by adding the supplied Size. Returns itself for easy chaining.
func (s *Size[T]) Add(size Size[T]) *Size[T] {
	s.Width += size.Width
	s.Height += size.Height
	return s
}

// AddInsets modifies this Size by expanding it to accommodate the specified insets. Returns itself for easy chaining.
func (s *Size[T]) AddInsets(insets Insets[T]) *Size[T] {
	s.Width += insets.Left + insets.Right
	s.Height += insets.Top + insets.Bottom
	return s
}

// Subtract modifies this Size by subtracting the supplied Size. Returns itself for easy chaining.
func (s *Size[T]) Subtract(size Size[T]) *Size[T] {
	s.Width -= size.Width
	s.Height -= size.Height
	return s
}

// SubtractInsets modifies this Size by reducing it to accommodate the specified insets. Returns itself for easy
// chaining.
func (s *Size[T]) SubtractInsets(insets Insets[T]) *Size[T] {
	s.Width -= insets.Left + insets.Right
	s.Height -= insets.Top + insets.Bottom
	return s
}

// GrowToInteger modifies this Size such that its width and height are both the smallest integers greater than or equal
// to their original values. Returns itself for easy chaining.
func (s *Size[T]) GrowToInteger() *Size[T] {
	switch reflect.TypeOf(s.Width).Kind() {
	case reflect.Float32, reflect.Float64:
		s.Width = T(math.Ceil(reflect.ValueOf(s.Width).Float()))
		s.Height = T(math.Ceil(reflect.ValueOf(s.Height).Float()))
	}
	return s
}

// ConstrainForHint ensures this size is no larger than the hint. Hint values less than one are ignored. Returns itself
// for easy chaining.
func (s *Size[T]) ConstrainForHint(hint Size[T]) *Size[T] {
	if hint.Width >= 1 && s.Width > hint.Width {
		s.Width = hint.Width
	}
	if hint.Height >= 1 && s.Height > hint.Height {
		s.Height = hint.Height
	}
	return s
}

// Min modifies this Size to contain the smallest values between itself and 'other'. Returns itself for easy chaining.
func (s *Size[T]) Min(other Size[T]) *Size[T] {
	if s.Width > other.Width {
		s.Width = other.Width
	}
	if s.Height > other.Height {
		s.Height = other.Height
	}
	return s
}

// Max modifies this Size to contain the largest values between itself and 'other'. Returns itself for easy chaining.
func (s *Size[T]) Max(other Size[T]) *Size[T] {
	if s.Width < other.Width {
		s.Width = other.Width
	}
	if s.Height < other.Height {
		s.Height = other.Height
	}
	return s
}

// String implements the fmt.Stringer interface.
func (s *Size[T]) String() string {
	return fmt.Sprintf("%v,%v", s.Width, s.Height)
}
