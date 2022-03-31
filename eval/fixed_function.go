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
	"math"
	"strings"

	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// FixedFunctions returns standard functions that work with 64-bit fixed-point values.
func FixedFunctions[T fixed.F64]() map[string]Function {
	return map[string]Function{
		"abs":   fixedAbsolute[T],
		"cbrt":  fixedCubeRoot[T],
		"ceil":  fixedCeiling[T],
		"exp":   fixedBaseEExponential[T],
		"exp2":  fixedBase2Exponential[T],
		"floor": fixedFloor[T],
		"if":    fixedIf[T],
		"log":   fixedNaturalLog[T],
		"log1p": fixedNaturalLogSum1[T],
		"log10": fixedDecimalLog[T],
		"max":   fixedMaximum[T],
		"min":   fixedMinimum[T],
		"round": fixedRound[T],
		"sqrt":  fixedSquareRoot[T],
	}
}

func fixedAbsolute[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return interface{}(value).(interface{ Abs() T }).Abs(), nil
}

func fixedBase2Exponential[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Exp2)
}

func fixedBaseEExponential[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Exp)
}

func fixedCeiling[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return interface{}(value).(interface{ Ceil() T }).Ceil(), nil
}

func fixedCubeRoot[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Cbrt)
}

func fixedDecimalLog[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Log10)
}

func fixedFloor[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return interface{}(value).(interface{ Trunc() T }).Trunc(), nil
}

func fixedIf[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	var arg string
	arg, arguments = NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value T
	if value, err = fixedFrom[T](evaluated); err != nil {
		if s, ok := evaluated.(string); ok {
			if s != "" && !strings.EqualFold(s, "false") {
				value = interface{}(value).(interface{ Inc() T }).Inc()
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

func fixedMaximum[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	max := T(^(1<<63 - 1))
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFixed[T](e, arg)
		if err != nil {
			return nil, err
		}
		max = interface{}(max).(interface{ Max(T) T }).Max(value)
	}
	return max, nil
}

func fixedMinimum[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	min := T(1<<63 - 1)
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFixed[T](e, arg)
		if err != nil {
			return nil, err
		}
		min = interface{}(min).(interface{ Min(T) T }).Min(value)
	}
	return min, nil
}

func fixedNaturalLog[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Log)
}

func fixedNaturalLogSum1[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	value = interface{}(value).(interface{ Inc() T }).Inc()
	return fixedFromFloat64[T](math.Log(interface{}(value).(interface{ AsFloat64() float64 }).AsFloat64()))
}

func fixedRound[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return interface{}(value).(interface{ Round() T }).Round(), nil
}

func fixedSquareRoot[T fixed.F64](e *Evaluator, arguments string) (interface{}, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Sqrt)
}

func evalToFixed[T fixed.F64](e *Evaluator, arg string) (T, error) {
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return 0, err
	}
	return fixedFrom[T](evaluated)
}

func fixedSingleNumberFunc[T fixed.F64](e *Evaluator, arguments string, f func(float64) float64) (interface{}, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return fixedFromFloat64[T](f(interface{}(value).(interface{ AsFloat64() float64 }).AsFloat64()))
}
