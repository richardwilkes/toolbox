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
	"strings"

	"github.com/richardwilkes/toolbox/errs"
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

// Add modifies this Insets by adding the supplied Insets. Returns itself for
// easy chaining.
func (i *Insets) Add(insets Insets) *Insets {
	i.Top += insets.Top
	i.Left += insets.Left
	i.Bottom += insets.Bottom
	i.Right += insets.Right
	return i
}

// Subtract modifies this Insets by subtracting the supplied Insets. Returns
// itself for easy chaining.
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

// MarshalText implements the encoding.TextMarshaler interface.
func (i Insets) MarshalText() (text []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (i *Insets) UnmarshalText(text []byte) error {
	txt := string(text)
	parts := strings.SplitN(strings.TrimSpace(txt), ",", 4)
	if len(parts) != 4 {
		return errs.Newf("unable to parse '%s'", txt)
	}
	var err error
	if i.Top, err = parseFloat(parts[0]); err != nil {
		return err
	}
	if i.Left, err = parseFloat(parts[1]); err != nil {
		return err
	}
	if i.Bottom, err = parseFloat(parts[2]); err != nil {
		return err
	}
	if i.Right, err = parseFloat(parts[3]); err != nil {
		return err
	}
	return nil
}
