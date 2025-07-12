// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

import (
	"fmt"

	"github.com/richardwilkes/toolbox/v2/geom"
)

// Size holds a fixed-point size.
type Size struct {
	Width  Num
	Height Num
}

// NewSize creates a new Size.
func NewSize(width, height Num) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

// SizeFrom converts a geom.Size into a Size.
func SizeFrom(s geom.Size) Size {
	return Size{
		Width:  NumFromFloat(s.Width),
		Height: NumFromFloat(s.Height),
	}
}

// Size converts this Size into a geom.Size.
func (s Size) Size() geom.Size {
	return geom.Size{
		Width:  NumAsFloat[float32](s.Width),
		Height: NumAsFloat[float32](s.Height),
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
func (s Size) Mul(value Num) Size {
	return Size{
		Width:  s.Width.Mul(value),
		Height: s.Height.Mul(value),
	}
}

// Div returns a new Size which is the result of dividing this Size by the value.
func (s Size) Div(value Num) Size {
	return Size{
		Width:  s.Width.Div(value),
		Height: s.Height.Div(value),
	}
}

// Floor returns a new Size with its width and height floored.
func (s Size) Floor() Size {
	return Size{
		Width:  s.Width.Floor(),
		Height: s.Height.Floor(),
	}
}

// Ceil returns a new Size with its width and height ceiled.
func (s Size) Ceil() Size {
	return Size{
		Width:  s.Width.Ceil(),
		Height: s.Height.Ceil(),
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

func (s Size) String() string {
	return fmt.Sprintf("%v,%v", s.Width, s.Height)
}
