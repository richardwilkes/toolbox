// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval

import "github.com/richardwilkes/toolbox/v2/fixed"

// NewFixed64Evaluator creates a new evaluator whose number type is one of the fixed64 types.
func NewFixed64Evaluator[T fixed.Dx](resolver VariableResolver, divideByZeroReturnsZero bool) *Evaluator {
	return &Evaluator{
		Resolver:  resolver,
		Operators: Fixed64Operators[T](divideByZeroReturnsZero),
		Functions: Fixed64Functions[T](),
	}
}
