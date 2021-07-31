// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package float64eval

import (
	"math"

	"github.com/richardwilkes/toolbox/eval"
)

// Functions returns standard functions that work with float64.
func Functions() map[string]eval.Function {
	return map[string]eval.Function{
		"abs":   Absolute,
		"exp2":  Base2Exponential,
		"exp":   BaseEExponential,
		"ceil":  Ceiling,
		"cbrt":  CubeRoot,
		"log10": DecimalLog,
		"floor": Floor,
		"if":    If,
		"max":   Maximum,
		"min":   Minimum,
		"log":   NaturalLog,
		"round": Round,
		"sqrt":  SquareRoot,
	}
}

// Absolute returns the absolute value of its argument.
func Absolute(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, func(v float64) float64 {
		if v < 0 {
			return -v
		}
		return v
	})
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
	return singleNumberFunc(e, arguments, math.Ceil)
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
	return singleNumberFunc(e, arguments, math.Floor)
}

// If returns the second argument if the first argument resolves to true, or the third argument if it doesn't.
func If(e *eval.Evaluator, arguments string) (interface{}, error) {
	var arg string
	arg, arguments = eval.NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value float64
	if value, err = NumberFrom(evaluated); err != nil {
		return nil, err
	}
	if value == 0 {
		_, arguments = eval.NextArg(arguments)
	}
	arg, _ = eval.NextArg(arguments)
	return e.EvaluateNew(arg)
}

// Maximum returns the maximum value of its input arguments.
func Maximum(e *eval.Evaluator, arguments string) (interface{}, error) {
	max := -math.MaxFloat64
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		evaluated, err := e.EvaluateNew(arg)
		if err != nil {
			return nil, err
		}
		var value float64
		if value, err = NumberFrom(evaluated); err != nil {
			return nil, err
		}
		if max < value {
			max = value
		}
	}
	return max, nil
}

// Minimum returns the minimum value of its input arguments.
func Minimum(e *eval.Evaluator, arguments string) (interface{}, error) {
	min := math.MaxFloat64
	for arguments != "" {
		var arg string
		arg, arguments = eval.NextArg(arguments)
		evaluated, err := e.EvaluateNew(arg)
		if err != nil {
			return nil, err
		}
		var value float64
		if value, err = NumberFrom(evaluated); err != nil {
			return nil, err
		}
		if min > value {
			min = value
		}
	}
	return min, nil
}

// NaturalLog returns the natural logarithm of its argument.
func NaturalLog(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Log)
}

// Round returns the nearest integer, rounding half away from zero.
func Round(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Round)
}

// SquareRoot returns the square root of it argument.
func SquareRoot(e *eval.Evaluator, arguments string) (interface{}, error) {
	return singleNumberFunc(e, arguments, math.Sqrt)
}

func singleNumberFunc(e *eval.Evaluator, arguments string, f func(float64) float64) (interface{}, error) {
	arg, err := e.EvaluateNew(arguments)
	if err != nil {
		return nil, err
	}
	var value float64
	value, err = NumberFrom(arg)
	if err != nil {
		return nil, err
	}
	return f(value), nil
}
