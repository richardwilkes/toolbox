// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// FixedOperators returns standard operators that work with 64-bit fixed-point values.
func FixedOperators[T fixed.F64](divideByZeroReturnsZero bool) []*Operator {
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
		Not(fixedNot[T]),
		Equal(fixedEqual[T]),
		NotEqual(fixedNotEqual[T]),
		GreaterThan(fixedGreaterThan[T]),
		GreaterThanOrEqual(fixedGreaterThanOrEqual[T]),
		LessThan(fixedLessThan[T]),
		LessThanOrEqual(fixedLessThanOrEqual[T]),
		Add(fixedAdd[T], fixedAddUnary[T]),
		Subtract(fixedSubtract[T], fixedSubtractUnary[T]),
		Multiply(fixedMultiply[T]),
		Divide(divide),
		Modulo(modulo),
		Power(fixedPower[T]),
	}
}

func fixedNot[T fixed.F64](arg interface{}) (interface{}, error) {
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

func fixedLogicalOr[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l != 0 {
		return true, nil
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func fixedLogicalAnd[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return false, nil
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

func fixedEqual[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	}
	return l == r, nil
}

func fixedNotEqual[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	}
	return l != r, nil
}

func fixedGreaterThan[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right), nil
	}
	return l > r, nil
}

func fixedGreaterThanOrEqual[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) >= fmt.Sprintf("%v", right), nil
	}
	return l >= r, nil
}

func fixedLessThan[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right), nil
	}
	return l < r, nil
}

func fixedLessThanOrEqual[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) <= fmt.Sprintf("%v", right), nil
	}
	return l <= r, nil
}

func fixedAdd[T fixed.F64](left, right interface{}) (interface{}, error) {
	var r T
	l, err := FixedFrom[T](left)
	if err == nil {
		r, err = FixedFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return l + r, nil
}

func fixedAddUnary[T fixed.F64](arg interface{}) (interface{}, error) {
	return FixedFrom[T](arg)
}

func fixedSubtract[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l - r, nil
}

func fixedSubtractUnary[T fixed.F64](arg interface{}) (interface{}, error) {
	v, err := FixedFrom[T](arg)
	if err != nil {
		return nil, err
	}
	return -v, nil
}

func fixedMultiply[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return interface{}(&l).(interface{ Mul(T) T }).Mul(r), nil
}

func fixedDivide[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return interface{}(&l).(interface{ Div(T) T }), nil
}

func fixedDivideAllowDivideByZero[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return interface{}(&l).(interface{ Div(T) T }).Div(r), nil
}

func fixedModulo[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return interface{}(&l).(interface{ Mod(T) T }).Mod(r), nil
}

func fixedModuloAllowDivideByZero[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return interface{}(&l).(interface{ Mod(T) T }).Mod(r), nil
}

func fixedPower[T fixed.F64](left, right interface{}) (interface{}, error) {
	l, err := FixedFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = FixedFrom[T](right)
	if err != nil {
		return nil, err
	}
	return FixedFromFloat64[T](math.Pow(interface{}(&l).(interface{ AsFloat64() float64 }).AsFloat64(),
		interface{}(&r).(interface{ AsFloat64() float64 }).AsFloat64())), nil
}

// FixedFromFloat64 converts the value into one of the fixed.F64 types.
func FixedFromFloat64[T fixed.F64](value float64) T {
	var t T
	return interface{}(&t).(interface{ FromFloat64(float64) T }).FromFloat64(value)
}

// FixedFrom attempts to convert the arg into one of the fixed.F64 types.
func FixedFrom[T fixed.F64](arg interface{}) (T, error) {
	switch a := arg.(type) {
	case bool:
		if a {
			var t T
			return interface{}(t).(interface{ Inc() T }).Inc(), nil
		}
		return 0, nil
	case T:
		return a, nil
	case string:
		var t T
		return interface{}(&t).(interface{ FromString(string) (T, error) }).FromString(a)
	default:
		return 0, errs.Newf("not a number: %v", arg)
	}
}
