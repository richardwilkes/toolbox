// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval_test

import (
	"fmt"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/eval"
	"github.com/richardwilkes/toolbox/v2/xmath/fixed"
)

var (
	numExpr = []string{
		"1 + 1", "2", "2.0000000000000000",
		"1.3 + 1.5", "2.8", "2.7999999999999998",
		"1.3+1.5", "2.8", "2.7999999999999998",
		"1.30015 + 1.5", "2.8001", "2.8001499999999999",
		"1 / 3", "0.3333", "0.3333333333333333",
		"1 / 3 + 10", "10.3333", "10.3333333333333339",
		"1 / (3 + 10)", "0.0769", "0.0769230769230769",
		"-1 / (3 + 10)", "-0.0769", "-0.0769230769230769",
		"1 / 0", "0", "0.0000000000000000",
		"sqrt(9)", "3", "3.0000000000000000",
		"sqrt(2)", "1.4142", "1.4142135623730951",
		"sqrt(cbrt(8)+7)", "3", "3.0000000000000000",
		"  sqrt\t(  cbrt    ( 8 ) +     7.0000 )    ", "3", "3.0000000000000000",
		"$foo + $bar", "24", "24.0000000000000000",
		"$foo / $bar", "11", "11.0000000000000000",
		"2 + 1e-2", "2.01", "2.0099999999999998",
		"2 + 1e2", "102", "102.0000000000000000",
		"(1 + 1) + 1", "3", "3.0000000000000000",
		"(1 + (1 + 1))", "3", "3.0000000000000000",
		"(1 + ((1) + 1))", "3", "3.0000000000000000",
		"(1 + (((1)) + 1))", "3", "3.0000000000000000",
		"1 + abs(1)", "2", "2.0000000000000000",
		"(1 + abs(1))", "2", "2.0000000000000000",
		"1 + (abs(1))", "2", "2.0000000000000000",
		"(abs(1)) + 1", "2", "2.0000000000000000",
		"(1 + abs(1)) + 1", "3", "3.0000000000000000",
		"max(0, 1)", "1", "1.0000000000000000",
		"(1 + max(0, 0)) - 10", "-9", "-9.0000000000000000",
		"abs(-12)", "12", "12.0000000000000000",
		"min(0, 1)", "0", "0.0000000000000000",
		"(1 + (2 * max(3, min(-4, 5) + 2) - ((14 - (13 - (12 - (11 - (10 - (9 - (8 - (7 + 6))))))))))) - 10", "-1", "-1.0000000000000000",
	}

	strExpr = []string{
		"foo + bar", "foobar",
		"foo +               \n    bar", "foobar",
		"$other", "22 + 2",
		"if($foo > $bar, yes, no)", "yes",
		"if($foo < $bar, yes, no)", "no",
	}
)

func TestFixedEvaluator(t *testing.T) {
	e := eval.NewFixedEvaluator[fixed.D4](resolver{}, true)
	for i := 0; i < len(numExpr); i += 3 {
		result, err := e.Evaluate(numExpr[i])
		check.NoError(t, err, "%d: %s == %s", i, numExpr[i], numExpr[i+1])
		check.Equal(t, numExpr[i+1], fmt.Sprintf("%v", result), "%d: %s == %s", i, numExpr[i], numExpr[i+1])
	}
	for i := 0; i < len(strExpr); i += 2 {
		result, err := e.Evaluate(strExpr[i])
		check.NoError(t, err, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
		check.Equal(t, strExpr[i+1], result, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
	}

	result, err := e.Evaluate("2 >= 1")
	check.NoError(t, err)
	check.Equal(t, true, result)

	result, err = e.Evaluate("2 >= 2")
	check.NoError(t, err)
	check.Equal(t, true, result)

	result, err = e.Evaluate("2 >= 3")
	check.NoError(t, err)
	check.Equal(t, false, result)

	result, err = e.Evaluate("2 > 1")
	check.NoError(t, err)
	check.Equal(t, true, result)

	e = eval.NewFixedEvaluator[fixed.D4](resolver{}, false)
	_, err = e.Evaluate("1 / 0")
	check.Error(t, err)
}

func TestFloatEvaluator(t *testing.T) {
	e := eval.NewFloatEvaluator[float64](resolver{}, true)
	for i := 0; i < len(numExpr); i += 3 {
		result, err := e.Evaluate(numExpr[i])
		check.NoError(t, err, "%d: %s == %s", i, numExpr[i], numExpr[i+2])
		check.Equal(t, numExpr[i+2], fmt.Sprintf("%0.16f", result), "%d: %s == %s", i, numExpr[i], numExpr[i+2])
	}
	for i := 0; i < len(strExpr); i += 2 {
		result, err := e.Evaluate(strExpr[i])
		check.NoError(t, err, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
		check.Equal(t, strExpr[i+1], result, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
	}

	result, err := e.Evaluate("2 > 1")
	check.NoError(t, err)
	check.Equal(t, true, result)

	e = eval.NewFloatEvaluator[float64](resolver{}, false)
	_, err = e.Evaluate("1 / 0")
	check.Error(t, err)
}

type resolver struct{}

func (r resolver) ResolveVariable(variableName string) string {
	switch variableName {
	case "foo":
		return "22"
	case "bar":
		return "2"
	default:
		return "$foo + $bar"
	}
}
