// Package geom provides geometry primitives.
package geom

import (
	"fmt"
	"math"
)

// Rect defines a rectangle.
type Rect struct {
	Point
	Size
}

// NewRect creates a new Rect.
func NewRect(x, y, width, height float64) Rect {
	return Rect{
		Point: NewPoint(x, y),
		Size:  NewSize(width, height),
	}
}

// NewRectPtr creates a new *Rect.
func NewRectPtr(x, y, width, height float64) *Rect {
	r := NewRect(x, y, width, height)
	return &r
}

// Align modifies this rectangle to align with integer coordinates that would
// encompass the original rectangle.
func (r *Rect) Align() {
	x := math.Floor(r.X)
	r.Width = math.Ceil(r.X+r.Width) - x
	r.X = x
	y := math.Floor(r.Y)
	r.Height = math.Ceil(r.Y+r.Height) - y
	r.Y = y
}

// CopyAndZeroLocation creates a new copy of the Rect and sets the location of
// the copy to 0,0.
func (r *Rect) CopyAndZeroLocation() Rect {
	return Rect{Size: r.Size}
}

// CenterX returns the center x-coordinate of the rectangle.
func (r Rect) CenterX() float64 {
	return r.X + r.Width/2
}

// CenterY returns the center y-coordinate of the rectangle.
func (r Rect) CenterY() float64 {
	return r.Y + r.Height/2
}

// IsEmpty returns true if either the width or height is zero or less.
func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Intersects returns true if this rect and the other rect intersect.
func (r Rect) Intersects(other Rect) bool {
	if !r.IsEmpty() && !other.IsEmpty() {
		return math.Min(r.X+r.Width, other.X+other.Width)-math.Max(r.X, other.X) > 0 &&
			math.Min(r.Y+r.Height, other.Y+other.Height)-math.Max(r.Y, other.Y) > 0
	}
	return false
}

// Intersect this Rect with another Rect, storing the result into this Rect.
func (r *Rect) Intersect(other Rect) {
	if r.IsEmpty() || other.IsEmpty() {
		r.Width = 0
		r.Height = 0
	} else {
		x := math.Max(r.X, other.X)
		y := math.Max(r.Y, other.Y)
		w := math.Min(r.X+r.Width, other.X+other.Width) - x
		h := math.Min(r.Y+r.Height, other.Y+other.Height) - y
		if w > 0 && h > 0 {
			r.X = x
			r.Y = y
			r.Width = w
			r.Height = h
		} else {
			r.Width = 0
			r.Height = 0
		}
	}
}

// Union this Rect with another Rect, storing the result into this Rect.
func (r *Rect) Union(other Rect) {
	e1 := r.IsEmpty()
	e2 := other.IsEmpty()
	switch {
	case e1 && e2:
		r.Width = 0
		r.Height = 0
	case e1:
		*r = other
	case !e2:
		x := math.Min(r.X, other.X)
		y := math.Min(r.Y, other.Y)
		r.Width = math.Max(r.X+r.Width, other.X+other.Width) - x
		r.Height = math.Max(r.Y+r.Height, other.Y+other.Height) - y
		r.X = x
		r.Y = y
	}
}

// InsetUniform insets this Rect by the specified amount on all sides.
// Positive values make the Rect smaller, while negative values make it
// larger.
func (r *Rect) InsetUniform(amount float64) {
	r.X += amount
	r.Y += amount
	r.Width -= amount * 2
	if r.Width < 0 {
		r.Width = 0
		r.Height = 0
	} else {
		r.Height -= amount * 2
		if r.Height < 0 {
			r.Width = 0
			r.Height = 0
		}
	}
}

// Inset this Rect by the specified Insets.
func (r *Rect) Inset(insets Insets) {
	r.X += insets.Left
	r.Y += insets.Top
	r.Width -= insets.Left + insets.Right
	if r.Width <= 0 {
		r.Width = 0
	}
	r.Height -= insets.Top + insets.Bottom
	if r.Height < 0 {
		r.Height = 0
	}
}

// ContainsPoint returns true if the coordinates are within the Rect.
func (r Rect) ContainsPoint(pt Point) bool {
	if r.IsEmpty() {
		return false
	}
	return r.X <= pt.X && r.Y <= pt.Y && pt.X < r.X+r.Width && pt.Y < r.Y+r.Height
}

// ContainsRect returns true if this Rect fully contains the passed in Rect.
func (r Rect) ContainsRect(in Rect) bool {
	if r.IsEmpty() || in.IsEmpty() {
		return false
	}
	right := r.X + r.Width
	bottom := r.Y + r.Height
	inRight := in.X + in.Width - 1
	inBottom := in.Y + in.Height - 1
	return r.X <= in.X && r.Y <= in.Y && in.X < right && in.Y < bottom && r.X <= inRight && r.Y <= inBottom && inRight < right && inBottom < bottom
}

// String implements the fmt.Stringer interface.
func (r Rect) String() string {
	return fmt.Sprintf("%v, %v", r.Point, r.Size)
}
