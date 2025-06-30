// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f128

import (
	"encoding/json"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xmath/fixed"
)

// Fraction holds a fractional value.
type Fraction[T fixed.Dx] struct {
	Numerator   Int[T]
	Denominator Int[T]
}

// NewFraction creates a new fractional value from a string.
func NewFraction[T fixed.Dx](s string) Fraction[T] {
	parts := strings.SplitN(s, "/", 2)
	f := Fraction[T]{
		Numerator:   FromStringForced[T](strings.TrimSpace(parts[0])),
		Denominator: From[T, int](1),
	}
	if len(parts) > 1 {
		f.Denominator = FromStringForced[T](strings.TrimSpace(parts[1]))
	}
	return f
}

// Normalize the fraction, eliminating any division by zero and ensuring a positive denominator.
func (f *Fraction[T]) Normalize() {
	var zero Int[T]
	if f.Denominator == zero {
		f.Numerator = Int[T]{}
		f.Denominator = From[T, int](1)
	} else if f.Denominator.LessThan(zero) {
		negOne := From[T, int](-1)
		f.Numerator = f.Numerator.Mul(negOne)
		f.Denominator = f.Denominator.Mul(negOne)
	}
}

// Value returns the computed value.
func (f Fraction[T]) Value() Int[T] {
	n := f
	n.Normalize()
	return n.Numerator.Div(n.Denominator)
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive.
func (f Fraction[T]) StringWithSign() string {
	n := f
	n.Normalize()
	s := n.Numerator.StringWithSign()
	if n.Denominator == From[T, int](1) {
		return s
	}
	return s + "/" + n.Denominator.String()
}

func (f Fraction[T]) String() string {
	n := f
	n.Normalize()
	s := n.Numerator.String()
	if n.Denominator == From[T, int](1) {
		return s
	}
	return s + "/" + n.Denominator.String()
}

// MarshalJSON implements json.Marshaler.
func (f Fraction[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Fraction[T]) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}
	*f = NewFraction[T](s)
	f.Normalize()
	return nil
}
