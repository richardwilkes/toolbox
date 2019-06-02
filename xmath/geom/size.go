package geom

import (
	"fmt"
	"math"
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

// Add modifies this Size by adding the supplied Size.
func (s *Size) Add(size Size) {
	s.Width += size.Width
	s.Height += size.Height
}

// AddInsets modifies this Size by expanding it to accommodate the specified
// insets.
func (s *Size) AddInsets(insets Insets) {
	s.Width += insets.Left + insets.Right
	s.Height += insets.Top + insets.Bottom
}

// Subtract modifies this Size by subtracting the supplied Size.
func (s *Size) Subtract(size Size) {
	s.Width -= size.Width
	s.Height -= size.Height
}

// SubtractInsets modifies this Size by reducing it to accommodate the
// specified insets.
func (s *Size) SubtractInsets(insets Insets) {
	s.Width -= insets.Left + insets.Right
	s.Height -= insets.Top + insets.Bottom
}

// GrowToInteger modifies this Size such that its width and height are both
// the smallest integers greater than or equal to their original values.
func (s *Size) GrowToInteger() {
	s.Width = math.Ceil(s.Width)
	s.Height = math.Ceil(s.Height)
}

// ConstrainForHint ensures this size is no larger than the hint. Hint values
// less than one are ignored.
func (s *Size) ConstrainForHint(hint Size) {
	if hint.Width >= 1 && s.Width > hint.Width {
		s.Width = hint.Width
	}
	if hint.Height >= 1 && s.Height > hint.Height {
		s.Height = hint.Height
	}
}

// Min modifies this Size to contain the smallest values between itself and
// 'other'.
func (s *Size) Min(other Size) {
	if s.Width > other.Width {
		s.Width = other.Width
	}
	if s.Height > other.Height {
		s.Height = other.Height
	}
}

// Max modifies this Size to contain the largest values between itself and
// 'other'.
func (s *Size) Max(other Size) {
	if s.Width < other.Width {
		s.Width = other.Width
	}
	if s.Height < other.Height {
		s.Height = other.Height
	}
}

// String implements the fmt.Stringer interface.
func (s Size) String() string {
	return fmt.Sprintf("%v, %v", s.Width, s.Height)
}
