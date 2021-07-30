package f64d4eval

import "github.com/richardwilkes/toolbox/eval"

// NewEvaluator creates a new evaluator whose number type is fixed.F64d4.
func NewEvaluator(resolver eval.VariableResolver, divideByZeroReturnsZero bool) *eval.Evaluator {
	return &eval.Evaluator{
		Resolver:  resolver,
		Operators: Operators(divideByZeroReturnsZero),
		Functions: Functions(),
	}
}
