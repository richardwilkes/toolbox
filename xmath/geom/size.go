// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
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
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

// Size defines a width and height.
type Size struct {
	Width, Height float64
}

// NewSize creates a new Size.
func NewSize(width, height float64) Size {
	return Size{
		Width:  width,
		Height: height,
	}
}

// NewSizePtr creates a new *Size.
func NewSizePtr(x, y float64) *Size {
	s := NewSize(x, y)
	return &s
}

// Add modifies this Size by adding the supplied Size. Returns itself for easy chaining.
func (s *Size) Add(size Size) *Size {
	s.Width += size.Width
	s.Height += size.Height
	return s
}

// AddInsets modifies this Size by expanding it to accommodate the specified insets. Returns itself for easy chaining.
func (s *Size) AddInsets(insets Insets) *Size {
	s.Width += insets.Left + insets.Right
	s.Height += insets.Top + insets.Bottom
	return s
}

// Subtract modifies this Size by subtracting the supplied Size. Returns itself for easy chaining.
func (s *Size) Subtract(size Size) *Size {
	s.Width -= size.Width
	s.Height -= size.Height
	return s
}

// SubtractInsets modifies this Size by reducing it to accommodate the specified insets. Returns itself for easy
// chaining.
func (s *Size) SubtractInsets(insets Insets) *Size {
	s.Width -= insets.Left + insets.Right
	s.Height -= insets.Top + insets.Bottom
	return s
}

// GrowToInteger modifies this Size such that its width and height are both the smallest integers greater than or equal
// to their original values. Returns itself for easy chaining.
func (s *Size) GrowToInteger() *Size {
	s.Width = math.Ceil(s.Width)
	s.Height = math.Ceil(s.Height)
	return s
}

// ConstrainForHint ensures this size is no larger than the hint. Hint values less than one are ignored. Returns itself
// for easy chaining.
func (s *Size) ConstrainForHint(hint Size) *Size {
	if hint.Width >= 1 && s.Width > hint.Width {
		s.Width = hint.Width
	}
	if hint.Height >= 1 && s.Height > hint.Height {
		s.Height = hint.Height
	}
	return s
}

// Min modifies this Size to contain the smallest values between itself and 'other'. Returns itself for easy chaining.
func (s *Size) Min(other Size) *Size {
	if s.Width > other.Width {
		s.Width = other.Width
	}
	if s.Height > other.Height {
		s.Height = other.Height
	}
	return s
}

// Max modifies this Size to contain the largest values between itself and 'other'. Returns itself for easy chaining.
func (s *Size) Max(other Size) *Size {
	if s.Width < other.Width {
		s.Width = other.Width
	}
	if s.Height < other.Height {
		s.Height = other.Height
	}
	return s
}

// String implements the fmt.Stringer interface.
func (s Size) String() string {
	return fmt.Sprintf("%f,%f", s.Width, s.Height)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s Size) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *Size) UnmarshalText(text []byte) error {
	txt := string(text)
	parts := strings.SplitN(strings.TrimSpace(txt), ",", 2)
	if len(parts) != 2 {
		return errs.Newf("unable to parse '%s'", txt)
	}
	var err error
	if s.Width, err = parseFloat(parts[0]); err != nil {
		return err
	}
	if s.Height, err = parseFloat(parts[1]); err != nil {
		return err
	}
	return nil
}
