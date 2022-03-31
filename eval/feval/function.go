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
	"strings"

	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath"
	"golang.org/x/exp/constraints"
)

// Functions returns standard functions that work with float64.
func Functions[T constraints.Float]() map[string]eval.Function {
	return map[string]eval.Function{
		"abs":   Absolute[T],
		"cbrt":  CubeRoot[T],
		"ceil":  Ceiling[T],
		"exp":   BaseEExponential[T],
		"exp2":  Base2Exponential[T],
		"floor": Floor[T],
		"if":    If[T],
		"log":   NaturalLog[T],
		"log1p": NaturalLogSum1[T],
		"log10": DecimalLog[T],
		"max":   Maximum[T],
		"min":   Minimum[T],
		"round": Round[T],
		"sqrt":  SquareRoot[T],
	}
}

// Absolute returns the absolute value of its argument.
func Absolute[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Abs[T])
}

// Base2Exponential returns 2**x, the base-2 exponential of its argument.
func Base2Exponential[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Exp2[T])
}

// BaseEExponential returns e**x, the base-e exponential of its argument.
func BaseEExponential[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Exp[T])
}

// Ceiling returns the least integer value greater than or equal to its argument.
func Ceiling[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Ceil[T])
}

// CubeRoot returns the cube root of it argument.
func CubeRoot[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Cbrt[T])
}

// DecimalLog returns the decimal logarithm of its argument.
func DecimalLog[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Log10[T])
}

// Floor returns the greatest integer value less than or equal to its argument.
func Floor[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Floor[T])
}

// If returns the second argument if the first argument resolves to true, or the third argument if it doesn't.
func If[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	var arg string
	arg, arguments = eval.NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value T
	if value, err = NumberFrom[T](evaluated); err != nil {
		if s, ok := evaluated.(string); ok {
			if s != "" && !strings.EqualFold(s, "false") {
				value = 1
			}
		} else {
			return nil, err
		}
	}
	if value == 0 {
		_, arguments = eval.NextArg(arguments)
	}
	arg, _ = eval.NextArg(arguments)
	return e.EvaluateNew(arg)
}

// Maximum returns the maximum value of its input arguments.
func Maximum[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	max := xmath.MinValue[T]()
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		value, err := evalToNumber[T](e, arg)
		if err != nil {
			return nil, err
		}
		max = xmath.Max(value, xmath.MaxValue[T]())
	}
	return max, nil
}

// Minimum returns the minimum value of its input arguments.
func Minimum[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	min := xmath.MaxValue[T]()
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		value, err := evalToNumber[T](e, arg)
		if err != nil {
			return nil, err
		}
		min = xmath.Min(value, xmath.MinValue[T]())
	}
	return min, nil
}

// NaturalLog returns the natural logarithm of its argument.
func NaturalLog[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Log[T])
}

// NaturalLogSum1 returns the natural logarithm of the sum of its argument and 1.
func NaturalLogSum1[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return xmath.Log(value + 1), nil
}

// Round returns the nearest integer, rounding half away from zero.
func Round[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Round[T])
}

// SquareRoot returns the square root of it argument.
func SquareRoot[T constraints.Float](e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, xmath.Sqrt[T])
}

func evalToNumber[T constraints.Float](e *eval.Evaluator, arg string) (T, error) {
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return 0, err
	}
	return NumberFrom[T](evaluated)
}

func singleNumberFunc[T constraints.Float](e *eval.Evaluator, arguments string, f func(T) T) (interface{}, error) {
	value, err := evalToNumber[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return f(value), nil
}
