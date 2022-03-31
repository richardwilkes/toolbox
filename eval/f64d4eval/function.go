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
	"math"
	"strings"

	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Functions returns standard functions that work with f64d4.Int.
func Functions() map[string]eval.Function {
	return map[string]eval.Function{
		"abs":   Absolute,
		"cbrt":  CubeRoot,
		"ceil":  Ceiling,
		"exp":   BaseEExponential,
		"exp2":  Base2Exponential,
		"floor": Floor,
		"if":    If,
		"log":   NaturalLog,
		"log1p": NaturalLogSum1,
		"log10": DecimalLog,
		"max":   Maximum,
		"min":   Minimum,
		"round": Round,
		"sqrt":  SquareRoot,
	}
}

// Absolute returns the absolute value of its argument.
func Absolute(e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return value.Abs(), nil
}

// Base2Exponential returns 2**x, the base-2 exponential of its argument.
func Base2Exponential(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Exp2)
}

// BaseEExponential returns e**x, the base-e exponential of its argument.
func BaseEExponential(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Exp)
}

// Ceiling returns the least integer value greater than or equal to its argument.
func Ceiling(e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return value.Ceil(), nil
}

// CubeRoot returns the cube root of it argument.
func CubeRoot(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Cbrt)
}

// DecimalLog returns the decimal logarithm of its argument.
func DecimalLog(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Log10)
}

// Floor returns the greatest integer value less than or equal to its argument.
func Floor(e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return value.Trunc(), nil
}

// If returns the second argument if the first argument resolves to true, or the third argument if it doesn't.
func If(e *eval.Evaluator, arguments string) (interface{}, error) {
	var arg string
	arg, arguments = eval.NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value f64d4.Int
	if value, err = NumberFrom(evaluated); err != nil {
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
func Maximum(e *eval.Evaluator, arguments string) (interface{}, error) {
	max := f64d4.Min
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		value, err := evalToNumber(e, arg)
		if err != nil {
			return nil, err
		}
		max = max.Max(value)
	}
	return max, nil
}

// Minimum returns the minimum value of its input arguments.
func Minimum(e *eval.Evaluator, arguments string) (interface{}, error) {
	min := f64d4.Max
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		value, err := evalToNumber(e, arg)
		if err != nil {
			return nil, err
		}
		min = min.Min(value)
	}
	return min, nil
}

// NaturalLog returns the natural logarithm of its argument.
func NaturalLog(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Log)
}

// NaturalLogSum1 returns the natural logarithm of the sum of its argument and 1.
func NaturalLogSum1(e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return f64d4.FromFloat64(math.Log((value + f64d4.FromInt(1)).AsFloat64())), nil
}

// Round returns the nearest integer, rounding half away from zero.
func Round(e *eval.Evaluator, arguments string) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return value.Round(), nil
}

// SquareRoot returns the square root of it argument.
func SquareRoot(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Sqrt)
}

func evalToNumber(e *eval.Evaluator, arg string) (f64d4.Int, error) {
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return 0, err
	}
	return NumberFrom(evaluated)
}

func singleNumberFunc(e *eval.Evaluator, arguments string, f func(float64) float64) (interface{}, error) {
	value, err := evalToNumber(e, arguments)
	if err != nil {
		return nil, err
	}
	return f64d4.FromFloat64(f(value.AsFloat64())), nil
}
