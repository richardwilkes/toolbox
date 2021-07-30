package float64eval

import "github.com/richardwilkes/toolbox/eval"

// NewEvaluator creates a new evaluator whose number type is float64.
func NewEvaluator(resolver eval.VariableResolver, divideByZeroReturnsZero bool) *eval.Evaluator {
	return &eval.Evaluator{
		Resolver:  resolver,
		Operators: Operators(divideByZeroReturnsZero),
		Functions: Functions(),
	}
}
