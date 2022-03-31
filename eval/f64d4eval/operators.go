// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f64d4eval

import (
	"fmt"
	"math"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Operators returns standard operators that work with f64d4.Int.
func Operators(divideByZeroReturnsZero bool) []*eval.Operator {
	var divide, modulo eval.OpFunc
	if divideByZeroReturnsZero {
		divide = DivideAllowDivideByZero
		modulo = ModuloAllowDivideByZero
	} else {
		divide = Divide
		modulo = Modulo
	}
	return []*eval.Operator{
		eval.OpenParen(),
		eval.CloseParen(),
		eval.Or(Or),
		eval.And(And),
		eval.Not(Not),
		eval.Equal(Equal),
		eval.NotEqual(NotEqual),
		eval.GreaterThan(GreaterThan),
		eval.GreaterThanOrEqual(GreaterThanOrEqual),
		eval.LessThan(LessThan),
		eval.LessThanOrEqual(LessThanOrEqual),
		eval.Add(Add, AddUnary),
		eval.Subtract(Subtract, SubtractUnary),
		eval.Multiply(Multiply),
		eval.Divide(divide),
		eval.Modulo(modulo),
		eval.Power(Power),
	}
}

// Not !
func Not(arg interface{}) (interface{}, error) {
	if b, ok := arg.(bool); ok {
		return !b, nil
	}
	v, err := NumberFrom(arg)
	if err != nil {
		return nil, err
	}
	if v == 0 {
		return true, nil
	}
	return false, nil
}

// Or ||
func Or(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	if l != 0 {
		return true, nil
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

// And &&
func And(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return false, nil
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	return r != 0, nil
}

// Equal ==
func Equal(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	}
	return l == r, nil
}

// NotEqual !=
func NotEqual(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	}
	return l != r, nil
}

// GreaterThan >
func GreaterThan(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right), nil
	}
	return l > r, nil
}

// GreaterThanOrEqual >=
func GreaterThanOrEqual(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) >= fmt.Sprintf("%v", right), nil
	}
	return l >= r, nil
}

// LessThan <
func LessThan(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right), nil
	}
	return l < r, nil
}

// LessThanOrEqual <=
func LessThanOrEqual(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v", left) <= fmt.Sprintf("%v", right), nil
	}
	return l <= r, nil
}

// Add + (addition)
func Add(left, right interface{}) (interface{}, error) {
	var r f64d4.Int
	l, err := NumberFrom(left)
	if err == nil {
		r, err = NumberFrom(right)
	}
	if err != nil {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return l + r, nil
}

// AddUnary + (plus)
func AddUnary(arg interface{}) (interface{}, error) {
	return NumberFrom(arg)
}

// Subtract - (subtraction)
func Subtract(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	return l - r, nil
}

// SubtractUnary - (minus)
func SubtractUnary(arg interface{}) (interface{}, error) {
	v, err := NumberFrom(arg)
	if err != nil {
		return nil, err
	}
	return -v, nil
}

// Multiply *
func Multiply(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	return l.Mul(r), nil
}

// Divide /
func Divide(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l.Div(r), nil
}

// DivideAllowDivideByZero / (returns 0 for division by 0)
func DivideAllowDivideByZero(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return l.Div(r), nil
}

// Modulo %
func Modulo(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return nil, errs.New("divide by zero")
	}
	return l.Mod(r), nil
}

// ModuloAllowDivideByZero % (returns 0 for modulo 0)
func ModuloAllowDivideByZero(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	if r == 0 {
		return r, nil
	}
	return l.Mod(r), nil
}

// Power ^
func Power(left, right interface{}) (interface{}, error) {
	l, err := NumberFrom(left)
	if err != nil {
		return nil, err
	}
	var r f64d4.Int
	r, err = NumberFrom(right)
	if err != nil {
		return nil, err
	}
	return f64d4.FromFloat64(math.Pow(l.AsFloat64(), r.AsFloat64())), nil
}

// NumberFrom attempts to extract a number from arg.
func NumberFrom(arg interface{}) (f64d4.Int, error) {
	switch a := arg.(type) {
	case bool:
		if a {
			return f64d4.One, nil
		}
		return 0, nil
	case f64d4.Int:
		return a, nil
	case string:
		return f64d4.FromString(a)
	default:
		return 0, errs.Newf("not a number: %v", arg)
	}
}
