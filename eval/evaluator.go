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
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

// VariableResolver is used to resolve variables in expressions into their values.
type VariableResolver interface {
	ResolveVariable(variableName string) string
}

type expressionOperand struct {
	value   string
	unaryOp *Operator
}

type expressionOperator struct {
	op      *Operator
	unaryOp *Operator
}

type expressionTree struct {
	evaluator *Evaluator
	left      any
	right     any
	op        *Operator
	unaryOp   *Operator
}

// Function provides a signature for a Function.
type Function func(evaluator *Evaluator, arguments string) (any, error)

type parsedFunction struct {
	function Function
	args     string
	unaryOp  *Operator
}

// Evaluator is used to evaluate an expression. If you do not have any variables that will be resolved, you can leave
// Resolver unset. StdOperators() and StdFunctions() can be used to populate the Operators and Functions fields.
type Evaluator struct {
	Resolver      VariableResolver
	Operators     []*Operator
	Functions     map[string]Function
	operandStack  []any
	operatorStack []*expressionOperator
}

// Evaluate an expression.
func (e *Evaluator) Evaluate(expression string) (any, error) {
	if err := e.parse(expression); err != nil {
		return nil, err
	}
	for len(e.operatorStack) != 0 {
		e.processTree()
	}
	if len(e.operandStack) == 0 {
		return "", nil
	}
	return e.evaluateOperand(e.operandStack[len(e.operandStack)-1])
}

// EvaluateNew reuses the Resolver, Operators, and Functions from this Evaluator to create a new Evaluator and then
// resolves an expression with it.
func (e *Evaluator) EvaluateNew(expression string) (any, error) {
	other := Evaluator{
		Resolver:  e.Resolver,
		Operators: e.Operators,
		Functions: e.Functions,
	}
	return other.Evaluate(expression)
}

func (e *Evaluator) parse(expression string) error {
	var unaryOp *Operator
	haveOperand := false
	haveOperator := false
	e.operandStack = nil
	e.operatorStack = nil
	i := 0
	for i < len(expression) {
		ch := expression[i]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			i++
			continue
		}
		opIndex, op := e.nextOperator(expression, i, nil)
		if opIndex > i || opIndex == -1 {
			var err error
			if i, err = e.processOperand(expression, i, opIndex, unaryOp); err != nil {
				return err
			}
			haveOperand = true
			haveOperator = false
			unaryOp = nil
		}
		if opIndex == i {
			if op != nil && op.EvaluateUnary != nil && (haveOperator || i == 0) {
				i = opIndex + len(op.Symbol)
				if unaryOp != nil {
					return errs.Newf("consecutive unary operators are not allowed at index %d", i)
				}
				unaryOp = op
			} else {
				var err error
				if i, err = e.processOperator(expression, opIndex, op, haveOperand, unaryOp); err != nil {
					return err
				}
				unaryOp = nil
			}
			if op == nil || op.Symbol != ")" {
				haveOperand = false
				haveOperator = true
			}
		}
	}
	return nil
}

func (e *Evaluator) nextOperator(expression string, start int, match *Operator) (int, *Operator) {
	for i := start; i < len(expression); i++ {
		if match != nil {
			if match.match(expression, i, len(expression)) {
				return i, match
			}
		} else {
			for _, op := range e.Operators {
				if op.match(expression, i, len(expression)) {
					return i, op
				}
			}
		}
	}
	return -1, nil
}

func (e *Evaluator) processOperand(expression string, start, opIndex int, unaryOp *Operator) (int, error) {
	if opIndex == -1 {
		text := strings.TrimSpace(expression[start:])
		if text == "" {
			return -1, errs.Newf("expression is invalid at index %d", start)
		}
		e.operandStack = append(e.operandStack, &expressionOperand{
			value:   text,
			unaryOp: unaryOp,
		})
		return len(expression), nil
	}
	text := strings.TrimSpace(expression[start:opIndex])
	if text == "" {
		return -1, errs.Newf("expression is invalid at index %d", start)
	}
	e.operandStack = append(e.operandStack, &expressionOperand{
		value:   text,
		unaryOp: unaryOp,
	})
	return opIndex, nil
}

func (e *Evaluator) processOperator(expression string, index int, op *Operator, haveOperand bool, unaryOp *Operator) (int, error) {
	if haveOperand && op != nil && op.Symbol == "(" {
		var err error
		index, op, err = e.processFunction(expression, index)
		if err != nil {
			return -1, err
		}
		index += len(op.Symbol)
		var tmp int
		tmp, op = e.nextOperator(expression, index, nil)
		if op == nil {
			return index, nil
		}
		index = tmp
	}
	switch op.Symbol {
	case "(":
		e.operatorStack = append(e.operatorStack, &expressionOperator{
			op:      op,
			unaryOp: unaryOp,
		})
	case ")":
		var stackOp *expressionOperator
		if len(e.operatorStack) > 0 {
			stackOp = e.operatorStack[len(e.operatorStack)-1]
		}
		for stackOp != nil && stackOp.op.Symbol != "(" {
			e.processTree()
			if len(e.operatorStack) > 0 {
				stackOp = e.operatorStack[len(e.operatorStack)-1]
			} else {
				stackOp = nil
			}
		}
		if len(e.operatorStack) == 0 {
			return -1, errs.Newf("invalid expression at index %d", index)
		}
		stackOp = e.operatorStack[len(e.operatorStack)-1]
		if stackOp.op.Symbol != "(" {
			return -1, errs.Newf("invalid expression at index %d", index)
		}
		e.operatorStack = e.operatorStack[:len(e.operatorStack)-1]
		if stackOp.unaryOp != nil {
			left := e.operandStack[len(e.operandStack)-1]
			e.operandStack = e.operandStack[:len(e.operandStack)-1]
			e.operandStack = append(e.operandStack, &expressionTree{
				evaluator: e,
				left:      left,
				unaryOp:   stackOp.unaryOp,
			})
		}
	default:
		if len(e.operatorStack) > 0 {
			stackOp := e.operatorStack[len(e.operatorStack)-1]
			for stackOp != nil && stackOp.op.Precedence >= op.Precedence {
				e.processTree()
				if len(e.operatorStack) > 0 {
					stackOp = e.operatorStack[len(e.operatorStack)-1]
				} else {
					stackOp = nil
				}
			}
		}
		e.operatorStack = append(e.operatorStack, &expressionOperator{
			op:      op,
			unaryOp: unaryOp,
		})
	}
	return index + len(op.Symbol), nil
}

func (e *Evaluator) processFunction(expression string, opIndex int) (int, *Operator, error) {
	var op *Operator
	parens := 1
	next := opIndex
	for parens > 0 {
		if next, op = e.nextOperator(expression, next+1, nil); op == nil {
			return -1, nil, errs.Newf("function not closed at index %d", opIndex)
		}
		switch op.Symbol {
		case "(":
			parens++
		case ")":
			parens--
		default:
		}
	}
	if len(e.operandStack) == 0 {
		return -1, nil, errs.Newf("invalid stack at index %d", next)
	}
	operand, ok := e.operandStack[len(e.operandStack)-1].(*expressionOperand)
	if !ok {
		return -1, nil, errs.Newf("unexpected operand stack value at index %d", next)
	}
	e.operandStack = e.operandStack[:len(e.operandStack)-1]
	f, exists := e.Functions[operand.value]
	if !exists {
		return -1, nil, errs.Newf("function not defined: %s", operand.value)
	}
	e.operandStack = append(e.operandStack, &parsedFunction{
		function: f,
		args:     expression[opIndex+1 : next],
		unaryOp:  operand.unaryOp,
	})
	return next, op, nil
}

func (e *Evaluator) processTree() {
	var right any
	if len(e.operandStack) > 0 {
		right = e.operandStack[len(e.operandStack)-1]
		e.operandStack = e.operandStack[:len(e.operandStack)-1]
	}
	var left any
	if len(e.operandStack) > 0 {
		left = e.operandStack[len(e.operandStack)-1]
		e.operandStack = e.operandStack[:len(e.operandStack)-1]
	}
	op := e.operatorStack[len(e.operatorStack)-1]
	e.operatorStack = e.operatorStack[:len(e.operatorStack)-1]
	e.operandStack = append(e.operandStack, &expressionTree{
		evaluator: e,
		left:      left,
		right:     right,
		op:        op.op,
	})
}

func (e *Evaluator) evaluateOperand(operand any) (any, error) {
	switch op := operand.(type) {
	case *expressionTree:
		left, err := op.evaluator.evaluateOperand(op.left)
		if err != nil {
			return nil, err
		}
		var right any
		right, err = op.evaluator.evaluateOperand(op.right)
		if err != nil {
			return nil, err
		}
		if op.left != nil && op.right != nil {
			if op.op.Evaluate == nil {
				return nil, errs.New("operator does not have Evaluate function defined")
			}
			var v any
			v, err = op.op.Evaluate(left, right)
			if err != nil {
				return nil, err
			}
			if op.unaryOp != nil && op.unaryOp.EvaluateUnary != nil {
				return op.unaryOp.EvaluateUnary(v)
			}
			return v, nil
		}
		var v any
		if op.right == nil {
			v = left
		} else {
			v = right
		}
		if v != nil {
			if op.unaryOp != nil && op.unaryOp.EvaluateUnary != nil {
				v, err = op.unaryOp.EvaluateUnary(v)
			} else if op.op != nil && op.op.EvaluateUnary != nil {
				v, err = op.op.EvaluateUnary(v)
			}
			if err != nil {
				return nil, err
			}
		}
		if v == nil {
			return nil, errs.New("expression is invalid")
		}
		return v, nil
	case *expressionOperand:
		v, err := e.replaceVariables(op.value)
		if err != nil {
			return nil, err
		}
		if op.unaryOp != nil && op.unaryOp.EvaluateUnary != nil {
			return op.unaryOp.EvaluateUnary(v)
		}
		return v, nil
	case *parsedFunction:
		s, err := e.replaceVariables(op.args)
		if err != nil {
			return nil, err
		}
		var v any
		v, err = op.function(e, s)
		if err != nil {
			return nil, err
		}
		if op.unaryOp != nil && op.unaryOp.EvaluateUnary != nil {
			return op.unaryOp.EvaluateUnary(v)
		}
		return v, nil
	default:
		if op != nil {
			return nil, errs.New("invalid expression")
		}
		return nil, nil
	}
}

func (e *Evaluator) replaceVariables(expression string) (string, error) {
	dollar := strings.IndexRune(expression, '$')
	if dollar == -1 {
		return expression, nil
	}
	if e.Resolver == nil {
		return "", errs.Newf("no variable resolver, yet variables present at index %d", dollar)
	}
	for dollar >= 0 {
		last := dollar
		for i, ch := range expression[dollar+1:] {
			if ch == '_' || ch == '.' || ch == '#' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (i != 0 && ch >= '0' && ch <= '9') {
				last = dollar + 1 + i
			} else {
				break
			}
		}
		if dollar == last {
			return "", errs.Newf("invalid variable at index %d", dollar)
		}
		name := expression[dollar+1 : last+1]
		v := e.Resolver.ResolveVariable(name)
		if strings.TrimSpace(v) == "" {
			return "", errs.Newf("unable to resolve variable $%s", name)
		}
		var buffer strings.Builder
		if dollar > 0 {
			buffer.WriteString(expression[:dollar])
		}
		buffer.WriteString(v)
		if last+1 < len(expression) {
			buffer.WriteString(expression[last+1:])
		}
		expression = buffer.String()
		dollar = strings.IndexRune(expression, '$')
	}
	return expression, nil
}

// NextArg provides extraction of the next argument from an arguments string passed to a Function. An empty string will
// be returned if no argument remains.
func NextArg(args string) (arg, remaining string) {
	parens := 0
	for i, ch := range args {
		switch {
		case ch == '(':
			parens++
		case ch == ')':
			parens--
		case ch == ',' && parens == 0:
			return args[:i], args[i+1:]
		default:
		}
	}
	return args, ""
}
