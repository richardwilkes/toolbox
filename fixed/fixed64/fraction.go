// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed64

import (
	"encoding/json"
	"strings"

	"github.com/richardwilkes/toolbox/v2/fixed"
	"github.com/richardwilkes/toolbox/v2/xmath"
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
		Denominator: FromInteger[T](1),
	}
	if len(parts) > 1 {
		f.Denominator = FromStringForced[T](strings.TrimSpace(parts[1]))
	}
	return f
}

// Normalize the fraction, eliminating any division by zero and ensuring a positive denominator.
func (f *Fraction[T]) Normalize() {
	if f.Denominator == 0 {
		f.Numerator = 0
		f.Denominator = FromInteger[T](1)
	} else if f.Denominator < 0 {
		negOne := FromInteger[T](-1)
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

// Add two fractions together and return the result.
func (f Fraction[T]) Add(other Fraction[T]) Fraction[T] {
	n := f
	n.Normalize()
	o := other
	o.Normalize()
	return Fraction[T]{
		Numerator:   n.Numerator.Mul(o.Denominator).Add(o.Numerator.Mul(n.Denominator)),
		Denominator: n.Denominator.Mul(o.Denominator),
	}
}

// Sub subtracts other from f and return the result.
func (f Fraction[T]) Sub(other Fraction[T]) Fraction[T] {
	n := f
	n.Normalize()
	o := other
	o.Normalize()
	return Fraction[T]{
		Numerator:   n.Numerator.Mul(o.Denominator).Sub(o.Numerator.Mul(n.Denominator)),
		Denominator: n.Denominator.Mul(o.Denominator),
	}
}

// Mul multiplies two fractions together and return the result.
func (f Fraction[T]) Mul(other Fraction[T]) Fraction[T] {
	n := f
	n.Normalize()
	o := other
	o.Normalize()
	return Fraction[T]{
		Numerator:   n.Numerator.Mul(o.Numerator),
		Denominator: n.Denominator.Mul(o.Denominator),
	}
}

// Div divides f by other and return the result.
func (f Fraction[T]) Div(other Fraction[T]) Fraction[T] {
	n := f
	n.Normalize()
	o := other
	o.Normalize()
	return Fraction[T]{
		Numerator:   n.Numerator.Mul(o.Denominator),
		Denominator: n.Denominator.Mul(o.Numerator),
	}
}

// Simplify the fraction, returning a new Fraction. If the numerator or denominator cannot be represented as an int, the
// original fraction is returned.
func (f Fraction[T]) Simplify() Fraction[T] {
	n := f
	n.Normalize()
	numerator := AsInteger[T, int](n.Numerator)
	if FromInteger[T](numerator) != n.Numerator {
		return n
	}
	denominator := AsInteger[T, int](n.Denominator)
	if FromInteger[T](denominator) != n.Denominator {
		return n
	}
	gcd := xmath.GCD(numerator, denominator)
	if gcd > 1 {
		g := FromInteger[T](gcd)
		n.Numerator = n.Numerator.Div(g)
		n.Denominator = n.Denominator.Div(g)
	}
	return n
}

// StringWithSign returns the same as String(), but prefixes the value with a '+' if it is positive.
func (f Fraction[T]) StringWithSign() string {
	n := f
	n.Normalize()
	s := n.Numerator.StringWithSign()
	if n.Denominator == FromInteger[T](1) {
		return s
	}
	return s + "/" + n.Denominator.String()
}

func (f Fraction[T]) String() string {
	n := f
	n.Normalize()
	s := n.Numerator.String()
	if n.Denominator == FromInteger[T](1) {
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
