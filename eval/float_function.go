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
	"github.com/richardwilkes/toolbox/v2/xmath"
	"github.com/richardwilkes/toolbox/v2/xstrings"
	"golang.org/x/exp/constraints"
)

// FloatFunctions returns standard functions that work with floats.
func FloatFunctions[T constraints.Float]() map[string]Function {
	return map[string]Function{
		"abs":   floatAbs[T],
		"cbrt":  floatCubeRoot[T],
		"ceil":  floatCeiling[T],
		"exp":   floatBaseEExponential[T],
		"exp2":  floatBase2Exponential[T],
		"floor": floatFloor[T],
		"if":    floatIf[T],
		"log":   floatNaturalLog[T],
		"log1p": floatNaturalLogSum1[T],
		"log10": floatDecimalLog[T],
		"max":   floatMaximum[T],
		"min":   floatMinimum[T],
		"round": floatRound[T],
		"sqrt":  floatSquareRoot[T],
	}
}

func floatAbs[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Abs[T])
}

func floatBase2Exponential[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Exp2[T])
}

func floatBaseEExponential[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Exp[T])
}

func floatCeiling[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Ceil[T])
}

func floatCubeRoot[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Cbrt[T])
}

func floatDecimalLog[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Log10[T])
}

func floatFloor[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Floor[T])
}

func floatIf[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	var arg string
	arg, arguments = NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value T
	if value, err = floatFrom[T](evaluated); err != nil {
		if s, ok := evaluated.(string); ok {
			if xstrings.IsTruthy(xstrings.Unquote(s)) {
				value = 1
			}
		} else {
			return nil, err
		}
	}
	if value == 0 {
		_, arguments = NextArg(arguments)
	}
	arg, _ = NextArg(arguments)
	return e.EvaluateNew(arg)
}

func floatMaximum[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	maxValue := xmath.MinValue[T]()
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFloat[T](e, arg)
		if err != nil {
			return nil, err
		}
		maxValue = max(value, maxValue)
	}
	return maxValue, nil
}

func floatMinimum[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	minValue := xmath.MaxValue[T]()
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFloat[T](e, arg)
		if err != nil {
			return nil, err
		}
		minValue = min(value, minValue)
	}
	return minValue, nil
}

func floatNaturalLog[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Log[T])
}

func floatNaturalLogSum1[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFloat[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return xmath.Log(value + 1), nil
}

func floatRound[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Round[T])
}

func floatSquareRoot[T constraints.Float](e *Evaluator, arguments string) (any, error) {
	return floatSingleNumberFunc(e, arguments, xmath.Sqrt[T])
}

func evalToFloat[T constraints.Float](e *Evaluator, arg string) (T, error) {
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return 0, err
	}
	return floatFrom[T](evaluated)
}

func floatSingleNumberFunc[T constraints.Float](e *Evaluator, arguments string, f func(T) T) (any, error) {
	value, err := evalToFloat[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return f(value), nil
}
