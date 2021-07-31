// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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
