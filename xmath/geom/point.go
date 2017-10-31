package geom

import (
	"fmt"
)

// Point defines a location.
type Point struct {
	X, Y float64
}

// Add modifies this Point by adding the supplied coordinates.
func (p *Point) Add(pt Point) {
	p.X += pt.X
	p.Y += pt.Y
}

// Subtract modifies this Point by subtracting the supplied coordinates.
func (p *Point) Subtract(pt Point) {
	p.X -= pt.X
	p.Y -= pt.Y
}

// String implements the fmt.Stringer interface.
func (p Point) String() string {
	return fmt.Sprintf("%v, %v", p.X, p.Y)
}
