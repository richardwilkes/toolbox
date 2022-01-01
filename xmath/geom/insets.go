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
)

// Insets defines margins on each side of a rectangle.
type Insets struct {
	Top    float64 `json:"top"`
	Left   float64 `json:"left"`
	Bottom float64 `json:"bottom"`
	Right  float64 `json:"right"`
}

// NewUniformInsets creates a new Insets whose edges all have the same value.
func NewUniformInsets(amount float64) Insets {
	return Insets{Top: amount, Left: amount, Bottom: amount, Right: amount}
}

// NewHorizontalInsets creates a new Insets whose left and right edges have the specified value.
func NewHorizontalInsets(amount float64) Insets {
	return Insets{Left: amount, Right: amount}
}

// NewVerticalInsets creates a new Insets whose top and bottom edges have the specified value.
func NewVerticalInsets(amount float64) Insets {
	return Insets{Top: amount, Bottom: amount}
}

// Add modifies this Insets by adding the supplied Insets. Returns itself for easy chaining.
func (i *Insets) Add(insets Insets) *Insets {
	i.Top += insets.Top
	i.Left += insets.Left
	i.Bottom += insets.Bottom
	i.Right += insets.Right
	return i
}

// Subtract modifies this Insets by subtracting the supplied Insets. Returns itself for easy chaining.
func (i *Insets) Subtract(insets Insets) *Insets {
	i.Top -= insets.Top
	i.Left -= insets.Left
	i.Bottom -= insets.Bottom
	i.Right -= insets.Right
	return i
}

// String implements the fmt.Stringer interface.
func (i Insets) String() string {
	return fmt.Sprintf("%f,%f,%f,%f", i.Top, i.Left, i.Bottom, i.Right)
}
