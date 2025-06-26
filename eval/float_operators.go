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
	"reflect"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath"
)

// FloatOperators returns standard operators that work with floating point values.
func FloatOperators[T ~float32 | ~float64](divideByZeroReturnsZero bool) []*Operator {
	var divide, modulo OpFunc
	if divideByZeroReturnsZero {
		divide = floatDivideAllowDivideByZero[T]
		modulo = floatModuloAllowDivideByZero[T]
	} else {
		divide = floatDivide[T]
		modulo = floatModulo[T]
	}
	return []*Operator{
		OpenParen(),
		CloseParen(),
		LogicalOr(floatLogicalOr[T]),
		LogicalAnd(floatLogicalAnd[T]),
		NotEqual(floatNotEqual[T]),
		Not(floatNot[T]),
		Equal(floatEqual[T]),
		GreaterThanOrEqual(floatGreaterThanOrEqual[T]),
		GreaterThan(floatGreaterThan[T]),
		LessThanOrEqual(floatLessThanOrEqual[T]),
		LessThan(floatLessThan[T]),
		Add(floatAdd[T], floatAddUnary[T]),
		Subtract(floatSubtract[T], floatSubtractUnary[T]),
		Multiply(floatMultiply[T]),
		Divide(divide),
		Modulo(modulo),
		Power(floatPower[T]),
	}
}

func floatNot[T ~float32 | ~float64](arg any) (any, error) {
	if b, ok := arg.(bool); ok {
		return !b, nil
	}
	v, err := floatFrom[T](arg)
	if err != nil {
		return nil, err
	}
	if v == 0 {
		return true, nil
	}
	return false, nil
}

func floatLogicalOr[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l != 0 {
		return true, nil
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func floatLogicalAnd[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return false, nil
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func floatEqual[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	}
	return l == r, nil
}

func floatNotEqual[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	}
	return l != r, nil
}

func floatGreaterThan[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right), nil
	}
	return l > r, nil
}

func floatGreaterThanOrEqual[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) >= fmt.Sprintf("%v", right), nil
	}
	return l >= r, nil
}

func floatLessThan[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right), nil
	}
	return l < r, nil
}

func floatLessThanOrEqual[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) <= fmt.Sprintf("%v", right), nil
	}
	return l <= r, nil
}

func floatAdd[T ~float32 | ~float64](left, right any) (any, error) {
	var r T
	l, err := floatFrom[T](left)
	if err == nil {
		r, err = floatFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return l + r, nil
}

func floatAddUnary[T ~float32 | ~float64](arg any) (any, error) {
	return floatFrom[T](arg)
}

func floatSubtract[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l - r, nil
}

func floatSubtractUnary[T ~float32 | ~float64](arg any) (any, error) {
	v, err := floatFrom[T](arg)
	if err != nil {
		return nil, err
	}
	return -v, nil
}

func floatMultiply[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l * r, nil
}

func floatDivide[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l / r, nil
}

func floatDivideAllowDivideByZero[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return l / r, nil
}

func floatModulo[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return xmath.Mod(l, r), nil
}

func floatModuloAllowDivideByZero[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return xmath.Mod(l, r), nil
}

func floatPower[T ~float32 | ~float64](left, right any) (any, error) {
	l, err := floatFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = floatFrom[T](right)
	if err != nil {
		return nil, err
	}
	return xmath.Pow(l, r), nil
}

func floatFrom[T ~float32 | ~float64](arg any) (T, error) {
	switch a := arg.(type) {
	case bool:
		if a {
			return 1, nil
		}
		return 0, nil
	case T:
		return a, nil
	case string:
		var t T
		f, err := strconv.ParseFloat(a, reflect.TypeOf(t).Bits())
		if err != nil {
			return 0, errs.Wrap(err)
		}
		return T(f), nil
	default:
		return 0, errs.Newf("not a number: %v", arg)
	}
}
