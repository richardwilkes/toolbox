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

// Point defines a location.
type Point struct {
	X, Y float64
}

// NewPoint creates a new Point.
func NewPoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// NewPointPtr creates a new *Point.
func NewPointPtr(x, y float64) *Point {
	p := NewPoint(x, y)
	return &p
}

// Align modifies this Point to align with integer coordinates. Returns itself for easy chaining.
func (p *Point) Align() *Point {
	p.X = math.Floor(p.X)
	p.Y = math.Floor(p.Y)
	return p
}

// Add modifies this Point by adding the supplied coordinates. Returns itself for easy chaining.
func (p *Point) Add(pt Point) *Point {
	p.X += pt.X
	p.Y += pt.Y
	return p
}

// Subtract modifies this Point by subtracting the supplied coordinates. Returns itself for easy chaining.
func (p *Point) Subtract(pt Point) *Point {
	p.X -= pt.X
	p.Y -= pt.Y
	return p
}

// Negate modifies this Point by negating both the X and Y coordinates.
func (p *Point) Negate() *Point {
	p.X = -p.X
	p.Y = -p.Y
	return p
}

// String implements the fmt.Stringer interface.
func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.X, p.Y)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p Point) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *Point) UnmarshalText(text []byte) error {
	txt := string(text)
	parts := strings.SplitN(strings.TrimSpace(txt), ",", 2)
	if len(parts) != 2 {
		return errs.Newf("unable to parse '%s'", txt)
	}
	var err error
	if p.X, err = parseFloat(parts[0]); err != nil {
		return err
	}
	if p.Y, err = parseFloat(parts[1]); err != nil {
		return err
	}
	return nil
}
