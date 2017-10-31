package geom

import (
	"fmt"
)

// Insets defines margins on each side of a rectangle.
type Insets struct {
	Top    float64
	Left   float64
	Bottom float64
	Right  float64
}

// NewUniformInsets creates a new Insets whose edges all have the same value.
func NewUniformInsets(amount float64) Insets {
	return Insets{Top: amount, Left: amount, Bottom: amount, Right: amount}
}

// NewHorizontalInsets creates a new Insets whose left and right edges have
// the specified value.
func NewHorizontalInsets(amount float64) Insets {
	return Insets{Left: amount, Right: amount}
}

// NewVerticalInsets creates a new Insets whose top and bottom edges have the
// specified value.
func NewVerticalInsets(amount float64) Insets {
	return Insets{Top: amount, Bottom: amount}
}

// Add modifies this Insets by adding the supplied Insets.
func (i *Insets) Add(insets Insets) {
	i.Top += insets.Top
	i.Left += insets.Left
	i.Bottom += insets.Bottom
	i.Right += insets.Right
}

// String implements the fmt.Stringer interface.
func (i Insets) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", i.Top, i.Left, i.Bottom, i.Right)
}
