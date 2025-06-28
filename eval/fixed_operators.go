// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval

import (
	"fmt"
	"math"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xmath/fixed"
	"github.com/richardwilkes/toolbox/v2/xmath/fixed/f64"
)

// FixedOperators returns standard operators that work with 64-bit fixed-point values.
func FixedOperators[T fixed.Dx](divideByZeroReturnsZero bool) []*Operator {
	var divide, modulo OpFunc
	if divideByZeroReturnsZero {
		divide = fixedDivideAllowDivideByZero[T]
		modulo = fixedModuloAllowDivideByZero[T]
	} else {
		divide = fixedDivide[T]
		modulo = fixedModulo[T]
	}
	return []*Operator{
		OpenParen(),
		CloseParen(),
		LogicalOr(fixedLogicalOr[T]),
		LogicalAnd(fixedLogicalAnd[T]),
		NotEqual(fixedNotEqual[T]),
		Not(fixedNot[T]),
		Equal(fixedEqual[T]),
		GreaterThanOrEqual(fixedGreaterThanOrEqual[T]),
		GreaterThan(fixedGreaterThan[T]),
		LessThanOrEqual(fixedLessThanOrEqual[T]),
		LessThan(fixedLessThan[T]),
		Add(fixedAdd[T], fixedAddUnary[T]),
		Subtract(fixedSubtract[T], fixedSubtractUnary[T]),
		Multiply(fixedMultiply[T]),
		Divide(divide),
		Modulo(modulo),
		Power(fixedPower[T]),
	}
}

func fixedNot[T fixed.Dx](arg any) (any, error) {
	if b, ok := arg.(bool); ok {
		return !b, nil
	}
	v, err := FixedFrom[T](arg)
	if err != nil {
		return nil, err
	}
	if v == 0 {
		return true, nil
	}
	return false, nil
}

func fixedLogicalOr[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l != 0 {
		return true, nil
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func fixedLogicalAnd[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return false, nil
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func fixedEqual[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	}
	return l == r, nil
}

func fixedNotEqual[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	}
	return l != r, nil
}

func fixedGreaterThan[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right), nil
	}
	return l > r, nil
}

func fixedGreaterThanOrEqual[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) >= fmt.Sprintf("%v", right), nil
	}
	return l >= r, nil
}

func fixedLessThan[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right), nil
	}
	return l < r, nil
}

func fixedLessThanOrEqual[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) <= fmt.Sprintf("%v", right), nil
	}
	return l <= r, nil
}

func fixedAdd[T fixed.Dx](left, right any) (any, error) {
	var r f64.Int[T]
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return l + r, nil
}

func fixedAddUnary[T fixed.Dx](arg any) (any, error) {
	return FixedFrom[T](arg)
}

func fixedSubtract[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l - r, nil
}

func fixedSubtractUnary[T fixed.Dx](arg any) (any, error) {
	v, err := FixedFrom[T](arg)
	if err != nil {
		return nil, err
	}
	return -v, nil
}

func fixedMultiply[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l.Mul(r), nil
}

func fixedDivide[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l.Div(r), nil
}

func fixedDivideAllowDivideByZero[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return l.Div(r), nil
}

func fixedModulo[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l.Mod(r), nil
}

func fixedModuloAllowDivideByZero[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return l.Mod(r), nil
}

func fixedPower[T fixed.Dx](left, right any) (any, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r f64.Int[T]
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return f64.From[T](math.Pow(f64.As[T, float64](l), f64.As[T, float64](r))), nil
}

// FixedFrom attempts to convert the arg into one of the fixed.F64 types.
func FixedFrom[T fixed.Dx](arg any) (f64.Int[T], error) {
	switch a := arg.(type) {
	case bool:
		if a {
			return f64.From[T](1), nil
		}
		return 0, nil
	case f64.Int[T]:
		return a, nil
	case string:
		return f64.FromString[T](a)
	default:
		return 0, errs.Newf("not a number: %v", arg)
	}
}
