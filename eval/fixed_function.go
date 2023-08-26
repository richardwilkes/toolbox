// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
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
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// FixedFunctions returns standard functions that work with 64-bit fixed-point values.
func FixedFunctions[T fixed.Dx]() map[string]Function {
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

func fixedAbsolute[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return any(&value).(interface{ Abs() f64.Int[T] }).Abs(), nil
}

func fixedBase2Exponential[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Exp2)
}

func fixedBaseEExponential[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Exp)
}

func fixedCeiling[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return any(&value).(interface{ Ceil() f64.Int[T] }).Ceil(), nil
}

func fixedCubeRoot[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Cbrt)
}

func fixedDecimalLog[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Log10)
}

func fixedFloor[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return any(&value).(interface{ Trunc() f64.Int[T] }).Trunc(), nil
}

func fixedIf[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	var arg string
	arg, arguments = NextArg(arguments)
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return nil, err
	}
	var value f64.Int[T]
	if value, err = FixedFrom[T](evaluated); err != nil {
		if s, ok := evaluated.(string); ok {
			if s != "" && !strings.EqualFold(s, "false") {
				value = any(&value).(interface{ Inc() f64.Int[T] }).Inc()
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

func fixedMaximum[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	maximum := f64.Int[T](f64.Min)
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFixed[T](e, arg)
		if err != nil {
			return nil, err
		}
		maximum = maximum.Max(value)
	}
	return maximum, nil
}

func fixedMinimum[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	minimum := f64.Int[T](f64.Max)
	for arguments != "" {
		var arg string
		arg, arguments = NextArg(arguments)
		value, err := evalToFixed[T](e, arg)
		if err != nil {
			return nil, err
		}
		minimum = minimum.Min(value)
	}
	return minimum, nil
}

func fixedNaturalLog[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Log)
}

func fixedNaturalLogSum1[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	value = any(&value).(interface{ Inc() f64.Int[T] }).Inc()
	return f64.From[T](math.Log(any(&value).(interface{ AsFloat64() float64 }).AsFloat64())), nil
}

func fixedRound[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return any(&value).(interface{ Round() f64.Int[T] }).Round(), nil
}

func fixedSquareRoot[T fixed.Dx](e *Evaluator, arguments string) (any, error) {
	return fixedSingleNumberFunc[T](e, arguments, math.Sqrt)
}

func evalToFixed[T fixed.Dx](e *Evaluator, arg string) (f64.Int[T], error) {
	evaluated, err := e.EvaluateNew(arg)
	if err != nil {
		return 0, err
	}
	return FixedFrom[T](evaluated)
}

func fixedSingleNumberFunc[T fixed.Dx](e *Evaluator, arguments string, f func(float64) float64) (any, error) {
	value, err := evalToFixed[T](e, arguments)
	if err != nil {
		return nil, err
	}
	return f64.From[T](f(f64.As[T, float64](value))), nil
}
