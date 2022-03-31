// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package feval

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath"
	"golang.org/x/exp/constraints"
)

// Operators returns standard operators that work with floating point values.
func Operators[T constraints.Float](divideByZeroReturnsZero bool) []*eval.Operator {
	var divide, modulo eval.OpFunc
	if divideByZeroReturnsZero {
		divide = DivideAllowDivideByZero[T]
		modulo = ModuloAllowDivideByZero[T]
	} else {
		divide = Divide[T]
		modulo = Modulo[T]
	}
	return []*eval.Operator{
		eval.OpenParen(),
		eval.CloseParen(),
		eval.Or(Or[T]),
		eval.And(And[T]),
		eval.Not(Not[T]),
		eval.Equal(Equal[T]),
		eval.NotEqual(NotEqual[T]),
		eval.GreaterThan(GreaterThan[T]),
		eval.GreaterThanOrEqual(GreaterThanOrEqual[T]),
		eval.LessThan(LessThan[T]),
		eval.LessThanOrEqual(LessThanOrEqual[T]),
		eval.Add(Add[T], AddUnary[T]),
		eval.Subtract(Subtract[T], SubtractUnary[T]),
		eval.Multiply(Multiply[T]),
		eval.Divide(divide),
		eval.Modulo(modulo),
		eval.Power(Power[T]),
	}
}

// Not !
func Not[T constraints.Float](arg interface{}) (interface{}, error) {
	if b, ok := arg.(bool); ok {
		return !b, nil
	}
	v, err := NumberFrom[T](arg)
	if err != nil {
		return nil, err
	}
	if v == 0 {
		return true, nil
	}
	return false, nil
}

// Or ||
func Or[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l != 0 {
		return true, nil
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

// And &&
func And[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return false, nil
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

// Equal ==
func Equal[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	}
	return l == r, nil
}

// NotEqual !=
func NotEqual[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	}
	return l != r, nil
}

// GreaterThan >
func GreaterThan[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right), nil
	}
	return l > r, nil
}

// GreaterThanOrEqual >=
func GreaterThanOrEqual[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) >= fmt.Sprintf("%v", right), nil
	}
	return l >= r, nil
}

// LessThan <
func LessThan[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right), nil
	}
	return l < r, nil
}

// LessThanOrEqual <=
func LessThanOrEqual[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) <= fmt.Sprintf("%v", right), nil
	}
	return l <= r, nil
}

// Add + (addition)
func Add[T constraints.Float](left, right interface{}) (interface{}, error) {
	var r T
	l, err := NumberFrom[T](left)
	if err == nil {
		r, err = NumberFrom[T](right)
	}
	if err != nil {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return l + r, nil
}

// AddUnary + (plus)
func AddUnary[T constraints.Float](arg interface{}) (interface{}, error) {
	return NumberFrom[T](arg)
}

// Subtract - (subtraction)
func Subtract[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l - r, nil
}

// SubtractUnary - (minus)
func SubtractUnary[T constraints.Float](arg interface{}) (interface{}, error) {
	v, err := NumberFrom[T](arg)
	if err != nil {
		return nil, err
	}
	return -v, nil
}

// Multiply *
func Multiply[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	return l * r, nil
}

// Divide /
func Divide[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l / r, nil
}

// DivideAllowDivideByZero / (returns 0 for division by 0)
func DivideAllowDivideByZero[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return 0.0, nil
	}
	return l / r, nil
}

// Modulo %
func Modulo[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return xmath.Mod(l, r), nil
}

// ModuloAllowDivideByZero % (returns 0 for modulo 0)
func ModuloAllowDivideByZero[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return xmath.Mod(l, r), nil
}

// Power ^
func Power[T constraints.Float](left, right interface{}) (interface{}, error) {
	l, err := NumberFrom[T](left)
	if err != nil {
		return nil, err
	}
	var r T
	r, err = NumberFrom[T](right)
	if err != nil {
		return nil, err
	}
	return xmath.Pow(l, r), nil
}

// NumberFrom attempts to extract a number from arg.
func NumberFrom[T constraints.Float](arg interface{}) (T, error) {
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
