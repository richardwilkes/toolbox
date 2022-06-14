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
	"unicode"
	"unicode/utf8"
)

// OpFunc provides a signature for an Operator's Evaluate function.
type OpFunc func(left, right any) (any, error)

// UnaryOpFunc provides a signature for an Operator's EvaluateUnary function.
type UnaryOpFunc func(arg any) (any, error)

// Operator provides an operator implementation for the Evaluator.
type Operator struct {
	Symbol        string
	Precedence    int
	Evaluate      OpFunc
	EvaluateUnary UnaryOpFunc
}

func (o *Operator) match(expression string, start, max int) bool {
	if max-start < len(o.Symbol) {
		return false
	}
	matches := o.Symbol == expression[start:start+len(o.Symbol)]
	// Hack to allow negative exponents on floating point numbers (i.e. 1.2e-2)
	if matches && len(o.Symbol) == 1 && o.Symbol == "-" && start > 1 && expression[start-1:start] == "e" {
		ch, _ := utf8.DecodeRuneInString(expression[start-2 : start-1])
		if unicode.IsDigit(ch) {
			return false
		}
	}
	return matches
}

// OpenParen (
func OpenParen() *Operator {
	return &Operator{Symbol: "("}
}

// CloseParen )
func CloseParen() *Operator {
	return &Operator{Symbol: ")"}
}

// Not !
func Not(f UnaryOpFunc) *Operator {
	return &Operator{
		Symbol:        "!",
		EvaluateUnary: f,
	}
}

// LogicalOr ||
func LogicalOr(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "||",
		Precedence: 10,
		Evaluate:   f,
	}
}

// LogicalAnd &&
func LogicalAnd(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "&&",
		Precedence: 20,
		Evaluate:   f,
	}
}

// Equal ==
func Equal(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "==",
		Precedence: 30,
		Evaluate:   f,
	}
}

// NotEqual !=
func NotEqual(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "!=",
		Precedence: 30,
		Evaluate:   f,
	}
}

// GreaterThan >
func GreaterThan(f OpFunc) *Operator {
	return &Operator{
		Symbol:     ">",
		Precedence: 40,
		Evaluate:   f,
	}
}

// GreaterThanOrEqual >=
func GreaterThanOrEqual(f OpFunc) *Operator {
	return &Operator{
		Symbol:     ">=",
		Precedence: 40,
		Evaluate:   f,
	}
}

// LessThan <
func LessThan(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "<",
		Precedence: 40,
		Evaluate:   f,
	}
}

// LessThanOrEqual <=
func LessThanOrEqual(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "<=",
		Precedence: 40,
		Evaluate:   f,
	}
}

// Add +
func Add(f OpFunc, unary UnaryOpFunc) *Operator {
	return &Operator{
		Symbol:        "+",
		Precedence:    50,
		Evaluate:      f,
		EvaluateUnary: unary,
	}
}

// Subtract -
func Subtract(f OpFunc, unary UnaryOpFunc) *Operator {
	return &Operator{
		Symbol:        "-",
		Precedence:    50,
		Evaluate:      f,
		EvaluateUnary: unary,
	}
}

// Multiply *
func Multiply(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "*",
		Precedence: 60,
		Evaluate:   f,
	}
}

// Divide /
func Divide(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "/",
		Precedence: 60,
		Evaluate:   f,
	}
}

// Modulo %
func Modulo(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "%",
		Precedence: 60,
		Evaluate:   f,
	}
}

// Power ^
func Power(f OpFunc) *Operator {
	return &Operator{
		Symbol:     "^",
		Precedence: 70,
		Evaluate:   f,
	}
}
