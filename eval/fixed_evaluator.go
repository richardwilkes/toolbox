// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval

import "github.com/richardwilkes/toolbox/xmath/fixed"

// NewFixedEvaluator creates a new evaluator whose number type is one of the fixed types.
func NewFixedEvaluator[T fixed.Dx](resolver VariableResolver, divideByZeroReturnsZero bool) *Evaluator {
	return &Evaluator{
		Resolver:  resolver,
		Operators: FixedOperators[T](divideByZeroReturnsZero),
		Functions: FixedFunctions[T](),
	}
}
