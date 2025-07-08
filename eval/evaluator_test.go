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
	"github.com/richardwilkes/toolbox/v2/fixed"
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
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)
	for i := 0; i < len(numExpr); i += 3 {
		result, err := e.Evaluate(numExpr[i])
		c.NoError(err, "%d: %s == %s", i, numExpr[i], numExpr[i+1])
		c.Equal(numExpr[i+1], fmt.Sprintf("%v", result), "%d: %s == %s", i, numExpr[i], numExpr[i+1])
	}
	for i := 0; i < len(strExpr); i += 2 {
		result, err := e.Evaluate(strExpr[i])
		c.NoError(err, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
		c.Equal(strExpr[i+1], result, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
	}

	result, err := e.Evaluate("2 >= 1")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 >= 2")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 >= 3")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("2 > 1")
	c.NoError(err)
	c.Equal(true, result)

	e = eval.NewFixed64Evaluator[fixed.D4](resolver{}, false)
	_, err = e.Evaluate("1 / 0")
	c.HasError(err)
}

func TestFloatEvaluator(t *testing.T) {
	c := check.New(t)
	e := eval.NewFloatEvaluator[float64](resolver{}, true)
	for i := 0; i < len(numExpr); i += 3 {
		result, err := e.Evaluate(numExpr[i])
		c.NoError(err, "%d: %s == %s", i, numExpr[i], numExpr[i+2])
		c.Equal(numExpr[i+2], fmt.Sprintf("%0.16f", result), "%d: %s == %s", i, numExpr[i], numExpr[i+2])
	}
	for i := 0; i < len(strExpr); i += 2 {
		result, err := e.Evaluate(strExpr[i])
		c.NoError(err, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
		c.Equal(strExpr[i+1], result, "%d: %s == %s", i, strExpr[i], strExpr[i+1])
	}

	result, err := e.Evaluate("2 > 1")
	c.NoError(err)
	c.Equal(true, result)

	e = eval.NewFloatEvaluator[float64](resolver{}, false)
	_, err = e.Evaluate("1 / 0")
	c.HasError(err)
}

func TestLogicalOperators(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test logical OR (||)
	result, err := e.Evaluate("1 || 0")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("0 || 1")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("0 || 0")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("2.5 || 0")
	c.NoError(err)
	c.Equal(true, result)

	// Test logical AND (&&)
	result, err = e.Evaluate("1 && 1")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("1 && 0")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("0 && 1")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("2.5 && 3.7")
	c.NoError(err)
	c.Equal(true, result)

	// Test logical NOT (!)
	result, err = e.Evaluate("!0")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("!1")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("!2.5")
	c.NoError(err)
	c.Equal(false, result)
}

func TestComparisonOperators(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test equality (==)
	result, err := e.Evaluate("2 == 2")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 == 3")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("2.5 == 2.5")
	c.NoError(err)
	c.Equal(true, result)

	// Test inequality (!=)
	result, err = e.Evaluate("2 != 3")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 != 2")
	c.NoError(err)
	c.Equal(false, result)

	// Test less than or equal (<=)
	result, err = e.Evaluate("2 <= 2")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 <= 3")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("3 <= 2")
	c.NoError(err)
	c.Equal(false, result)

	// Test greater than or equal (>=)
	result, err = e.Evaluate("3 >= 2")
	c.NoError(err)
	c.Equal(true, result)
}

func TestArithmeticOperators(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test modulo (%)
	result, err := e.Evaluate("10 % 3")
	c.NoError(err)
	c.Equal("1", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("10 % 2")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("7.5 % 2.5")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))

	// Test power (^)
	result, err = e.Evaluate("2 ^ 3")
	c.NoError(err)
	c.Equal("8", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("4 ^ 0.5")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("3 ^ 2")
	c.NoError(err)
	c.Equal("9", fmt.Sprintf("%v", result))

	// Test unary plus
	result, err = e.Evaluate("+5")
	c.NoError(err)
	c.Equal("5", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("+(2 + 3)")
	c.NoError(err)
	c.Equal("5", fmt.Sprintf("%v", result))

	// Test unary minus
	result, err = e.Evaluate("-5")
	c.NoError(err)
	c.Equal("-5", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("-(2 + 3)")
	c.NoError(err)
	c.Equal("-5", fmt.Sprintf("%v", result))
}

func TestMathFunctions(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test ceil function
	result, err := e.Evaluate("ceil(2.3)")
	c.NoError(err)
	c.Equal("3", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("ceil(2.0)")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("ceil(-2.3)")
	c.NoError(err)
	c.Equal("-2", fmt.Sprintf("%v", result))

	// Test floor function
	result, err = e.Evaluate("floor(2.7)")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("floor(2.0)")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("floor(-2.3)")
	c.NoError(err)
	c.Equal("-2", fmt.Sprintf("%v", result))

	// Test round function
	result, err = e.Evaluate("round(2.3)")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("round(2.7)")
	c.NoError(err)
	c.Equal("3", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("round(2.5)")
	c.NoError(err)
	c.Equal("3", fmt.Sprintf("%v", result))

	// Test exp function
	result, err = e.Evaluate("exp(0)")
	c.NoError(err)
	c.Equal("1", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("exp(1)")
	c.NoError(err)
	// e ≈ 2.718281828, but with D4 precision it's 2.7182
	c.Equal("2.7182", fmt.Sprintf("%v", result))

	// Test exp2 function
	result, err = e.Evaluate("exp2(3)")
	c.NoError(err)
	c.Equal("8", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("exp2(0)")
	c.NoError(err)
	c.Equal("1", fmt.Sprintf("%v", result))

	// Test log function
	result, err = e.Evaluate("log(1)")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))

	// Test log10 function
	result, err = e.Evaluate("log10(100)")
	c.NoError(err)
	c.Equal("2", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("log10(1)")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))

	// Test log1p function
	result, err = e.Evaluate("log1p(0)")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))
}

func TestIfFunction(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test if with numeric conditions
	result, err := e.Evaluate("if(1, yes, no)")
	c.NoError(err)
	c.Equal("yes", result)

	result, err = e.Evaluate("if(0, yes, no)")
	c.NoError(err)
	c.Equal("no", result)

	result, err = e.Evaluate("if(2.5, yes, no)")
	c.NoError(err)
	c.Equal("yes", result)

	// Test if with comparison expressions
	result, err = e.Evaluate("if(2 > 1, greater, smaller)")
	c.NoError(err)
	c.Equal("greater", result)

	result, err = e.Evaluate("if(1 > 2, greater, smaller)")
	c.NoError(err)
	c.Equal("smaller", result)

	// Test if with string conditions
	result, err = e.Evaluate("if(true, yes, no)")
	c.NoError(err)
	c.Equal("yes", result)

	result, err = e.Evaluate("if(\"true\", yes, no)")
	c.NoError(err)
	c.Equal("yes", result)

	result, err = e.Evaluate("if(false, yes, no)")
	c.NoError(err)
	c.Equal("no", result)

	result, err = e.Evaluate("if(\"\", yes, no)")
	c.NoError(err)
	c.Equal("no", result)

	result, err = e.Evaluate("if(\"hello\", yes, no)")
	c.NoError(err)
	c.Equal("no", result)

	result, err = e.Evaluate("if(\"false\", yes, no)")
	c.NoError(err)
	c.Equal("no", result)
}

func TestErrorHandling(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test empty expression
	result, err := e.Evaluate("")
	c.NoError(err)
	c.Equal("", result)

	// Test unmatched parentheses that actually cause errors
	_, err = e.Evaluate("1)")
	c.HasError(err)

	// Test undefined function
	_, err = e.Evaluate("undefined_func(1)")
	c.HasError(err)

	// Test function with unclosed parentheses
	_, err = e.Evaluate("sqrt(2")
	c.HasError(err)

	// Test evaluator without variable resolver but with variables
	eNoResolver := eval.NewFixed64Evaluator[fixed.D4](nil, true)
	_, err = eNoResolver.Evaluate("$foo")
	c.HasError(err)

	// Test invalid variable name
	_, err = e.Evaluate("$")
	c.HasError(err)

	// Test some other edge cases that should work
	result, err = e.Evaluate("1 +")
	c.NoError(err) // This actually parses as just "1"
	c.Equal("1", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("+ 1")
	c.NoError(err) // This parses as unary plus
	c.Equal("1", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("1 + + 2")
	c.NoError(err) // This parses as "1 + (+2)"
	c.Equal("3", fmt.Sprintf("%v", result))

	// Test variable that resolves to empty string
	eEmptyResolver := eval.NewFixed64Evaluator[fixed.D4](emptyResolver{}, true)
	_, err = eEmptyResolver.Evaluate("$empty")
	c.HasError(err)
}

func TestFloatEvaluatorSpecific(t *testing.T) {
	c := check.New(t)
	e := eval.NewFloatEvaluator[float64](resolver{}, true)

	// Test the same operators but with float evaluator
	result, err := e.Evaluate("!1")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("1 || 0")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("1 && 0")
	c.NoError(err)
	c.Equal(false, result)

	result, err = e.Evaluate("2 == 2")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 != 3")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("2 <= 3")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("3 >= 2")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("10 % 3")
	c.NoError(err)
	c.Equal(1.0, result)

	result, err = e.Evaluate("2 ^ 3")
	c.NoError(err)
	c.Equal(8.0, result)

	result, err = e.Evaluate("+5")
	c.NoError(err)
	c.Equal(5.0, result)

	// Test float-specific functions
	result, err = e.Evaluate("ceil(2.3)")
	c.NoError(err)
	c.Equal(3.0, result)

	result, err = e.Evaluate("floor(2.7)")
	c.NoError(err)
	c.Equal(2.0, result)

	result, err = e.Evaluate("round(2.7)")
	c.NoError(err)
	c.Equal(3.0, result)
}

func TestDivideByZeroErrors(t *testing.T) {
	c := check.New(t)

	// Test fixed evaluator with divideByZero=false
	eFixed := eval.NewFixed64Evaluator[fixed.D4](resolver{}, false)
	_, err := eFixed.Evaluate("1 / 0")
	c.HasError(err)

	_, err = eFixed.Evaluate("5 % 0")
	c.HasError(err)

	// Test float evaluator with divideByZero=false
	eFloat := eval.NewFloatEvaluator[float64](resolver{}, false)
	_, err = eFloat.Evaluate("1 / 0")
	c.HasError(err)

	_, err = eFloat.Evaluate("5 % 0")
	c.HasError(err)
}

func TestComplexExpressions(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test complex logical expressions
	result, err := e.Evaluate("(1 > 0) && (2 < 3)")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("!(1 > 2) || (3 == 3)")
	c.NoError(err)
	c.Equal(true, result)

	// Test complex arithmetic with functions
	result, err = e.Evaluate("sqrt(max(4, 9)) + min(1, 2)")
	c.NoError(err)
	c.Equal("4", fmt.Sprintf("%v", result))

	// Test nested function calls
	result, err = e.Evaluate("abs(sqrt(9) - 4)")
	c.NoError(err)
	c.Equal("1", fmt.Sprintf("%v", result))

	// Test mixed operators
	result, err = e.Evaluate("2 ^ 3 % 5")
	c.NoError(err)
	c.Equal("3", fmt.Sprintf("%v", result))

	// Test conditional with complex expressions
	result, err = e.Evaluate("if(sqrt(16) > 3, big, small)")
	c.NoError(err)
	c.Equal("big", result)
}

func TestNextArgFunction(t *testing.T) {
	c := check.New(t)

	// Test basic argument extraction
	arg, remaining := eval.NextArg("first,second,third")
	c.Equal("first", arg)
	c.Equal("second,third", remaining)

	// Test with spaces
	arg, remaining = eval.NextArg("  first  ,  second  ")
	c.Equal("  first  ", arg)
	c.Equal("  second  ", remaining)

	// Test with nested parentheses
	arg, remaining = eval.NextArg("func(a,b),next")
	c.Equal("func(a,b)", arg)
	c.Equal("next", remaining)

	// Test with multiple levels of nesting
	arg, remaining = eval.NextArg("func(a,func2(x,y)),next")
	c.Equal("func(a,func2(x,y))", arg)
	c.Equal("next", remaining)

	// Test single argument
	arg, remaining = eval.NextArg("onlyarg")
	c.Equal("onlyarg", arg)
	c.Equal("", remaining)

	// Test empty string
	arg, remaining = eval.NextArg("")
	c.Equal("", arg)
	c.Equal("", remaining)
}

func TestEvaluateNew(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test that EvaluateNew creates a new evaluator but reuses configuration
	result1, err1 := e.Evaluate("1 + 1")
	c.NoError(err1)

	result2, err2 := e.EvaluateNew("2 + 2")
	c.NoError(err2)

	c.Equal("2", fmt.Sprintf("%v", result1))
	c.Equal("4", fmt.Sprintf("%v", result2))

	// Test that original evaluator state is preserved
	result3, err3 := e.Evaluate("3 + 3")
	c.NoError(err3)
	c.Equal("6", fmt.Sprintf("%v", result3))
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

type emptyResolver struct{}

func (r emptyResolver) ResolveVariable(_ string) string {
	return "" // Always return empty string
}

func TestFloatMathFunctions(t *testing.T) {
	c := check.New(t)
	e := eval.NewFloatEvaluator[float64](resolver{}, true)

	// Test float-specific math functions that weren't covered
	result, err := e.Evaluate("exp2(3)")
	c.NoError(err)
	c.Equal(8.0, result)

	result, err = e.Evaluate("exp(1)")
	c.NoError(err)
	// e ≈ 2.718281828
	f, ok := result.(float64)
	c.True(ok)
	c.True(f > 2.7 && f < 2.8)

	result, err = e.Evaluate("log(1)")
	c.NoError(err)
	c.Equal(0.0, result)

	result, err = e.Evaluate("log10(100)")
	c.NoError(err)
	c.Equal(2.0, result)

	result, err = e.Evaluate("log1p(0)")
	c.NoError(err)
	c.Equal(0.0, result)
}

func TestBooleanOperatorEdgeCases(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test comparison operators with strings
	result, err := e.Evaluate("\"hello\" == \"hello\"")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("\"hello\" != \"world\"")
	c.NoError(err)
	c.Equal(true, result)

	// Test OR and AND with string comparisons
	result, err = e.Evaluate("\"a\" == \"a\" || \"b\" == \"c\"")
	c.NoError(err)
	c.Equal(true, result)

	result, err = e.Evaluate("\"a\" == \"a\" && \"b\" == \"b\"")
	c.NoError(err)
	c.Equal(true, result)

	// Test that bare string literals like 'true' and 'false' are treated as strings, not booleans
	// and cause errors when used with numeric operators
	_, err = e.Evaluate("!true")
	c.HasError(err) // 'true' is a string that can't be parsed as a number

	_, err = e.Evaluate("!false")
	c.HasError(err) // 'false' is a string that can't be parsed as a number
}

func TestErrorPathsAndEdgeCases(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test some error paths in functions that actually cause errors
	_, err := e.Evaluate("abs(invalid)")
	c.HasError(err)

	_, err = e.Evaluate("sqrt(invalid)")
	c.HasError(err)

	_, err = e.Evaluate("ceil(invalid)")
	c.HasError(err)

	_, err = e.Evaluate("floor(invalid)")
	c.HasError(err)

	_, err = e.Evaluate("round(invalid)")
	c.HasError(err)

	// Test unary operators with invalid operands that cause parsing errors
	_, err = e.Evaluate("-invalid")
	c.HasError(err)

	_, err = e.Evaluate("+invalid")
	c.HasError(err)

	// Note: Many operations like "1 + invalid" actually concatenate instead of erroring
	// This is because the parser treats them as string concatenation when parsing fails

	// Test that string concatenation works (not an error)
	result, err := e.Evaluate("1 + invalid")
	c.NoError(err)
	c.Equal("1invalid", result)

	// Test some operations that should actually error
	_, err = e.Evaluate("invalid * 2")
	c.HasError(err)

	_, err = e.Evaluate("invalid / 2")
	c.HasError(err)

	_, err = e.Evaluate("invalid % 2")
	c.HasError(err)

	_, err = e.Evaluate("invalid ^ 2")
	c.HasError(err)
}

func TestAdditionalFunctionTests(t *testing.T) {
	c := check.New(t)
	e := eval.NewFixed64Evaluator[fixed.D4](resolver{}, true)

	// Test if function with more complex conditions
	result, err := e.Evaluate("if(max(1, 2) > 1, big, small)")
	c.NoError(err)
	c.Equal("big", result)

	// Test max/min with more arguments through nested calls
	result, err = e.Evaluate("max(max(1, 2), 3)")
	c.NoError(err)
	c.Equal("3", fmt.Sprintf("%v", result))

	result, err = e.Evaluate("min(min(1, 2), 0)")
	c.NoError(err)
	c.Equal("0", fmt.Sprintf("%v", result))

	// Test functions with no arguments
	_, err = e.Evaluate("abs()")
	c.HasError(err) // Empty string is invalid for abs

	_, err = e.Evaluate("sqrt()")
	c.HasError(err) // Empty string is invalid for sqrt

	// Note: max() and min() don't error with empty args, they return boundary values
	_, err = e.Evaluate("max()")
	c.NoError(err) // Returns a very negative number

	_, err = e.Evaluate("min()")
	c.NoError(err) // Returns a very positive number

	// Test max/min with invalid arguments do cause errors
	_, err = e.Evaluate("max(1, invalid)")
	c.HasError(err)

	_, err = e.Evaluate("min(1, invalid)")
	c.HasError(err)
}
