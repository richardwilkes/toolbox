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

// Size defines a width and height.
type Size struct {
	Width  float32 `json:"w"`
	Height float32 `json:"h"`
}

// NewSize creates a new Size.
func NewSize(width, height float32) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

// SizeFromPoint creates a new Size from a Point.
func SizeFromPoint(pt Point) Size {
	return Size{
		Width:  pt.X,
		Height: pt.Y,
	}
}

// Add returns a new Size which is the result of adding this Size with the provided Size.
func (s Size) Add(size Size) Size {
	return Size{
		Width:  s.Width + size.Width,
		Height: s.Height + size.Height,
	}
}

// Sub returns a new Size which is the result of subtracting the provided Size from this Size.
func (s Size) Sub(size Size) Size {
	return Size{
		Width:  s.Width - size.Width,
		Height: s.Height - size.Height,
	}
}

// Mul returns a new Size which is the result of multiplying this Size by the value.
func (s Size) Mul(value float32) Size {
	return Size{
		Width:  s.Width * value,
		Height: s.Height * value,
	}
}

// Div returns a new Size which is the result of dividing this Size by the value.
func (s Size) Div(value float32) Size {
	return Size{
		Width:  s.Width / value,
		Height: s.Height / value,
	}
}

// Floor returns a new Size with its width and height floored.
func (s Size) Floor() Size {
	return Size{
		Width:  xmath.Floor(s.Width),
		Height: xmath.Floor(s.Height),
	}
}

// Ceil returns a new Size with its width and height ceiled.
func (s Size) Ceil() Size {
	return Size{
		Width:  xmath.Ceil(s.Width),
		Height: xmath.Ceil(s.Height),
	}
}

// Min returns the smallest Size between itself and 'other'.
func (s Size) Min(other Size) Size {
	return Size{
		Width:  min(s.Width, other.Width),
		Height: min(s.Height, other.Height),
	}
}

// Max returns the largest Size between itself and 'other'.
func (s Size) Max(other Size) Size {
	return Size{
		Width:  max(s.Width, other.Width),
		Height: max(s.Height, other.Height),
	}
}

// ConstrainForHint returns a size no larger than the hint value. Hint values less than one are ignored.
func (s Size) ConstrainForHint(hint Size) Size {
	w := s.Width
	if hint.Width >= 1 && w > hint.Width {
		w = hint.Width
	}
	h := s.Height
	if hint.Height >= 1 && h > hint.Height {
		h = hint.Height
	}
	return Size{
		Width:  w,
		Height: h,
	}
}

// String implements fmt.Stringer.
func (s Size) String() string {
	return fmt.Sprintf("%#v,%#v", s.Width, s.Height)
}
