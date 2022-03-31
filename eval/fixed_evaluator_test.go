// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/stretchr/testify/assert"
)

func TestFixedEvaluator(t *testing.T) {
	expected := []string{
		"2",
		"2.8",
		"2.8",
		"2.8001",
		"0.3333",
		"10.3333",
		"0.0769",
		"-0.0769",
		"0",
		"3",
		"1.4142",
		"3",
		"3",
		"24",
		"11",
		"2.01",
		"102",
	}
	e := eval.NewFixedEvaluator[f64d4.Int](resolver{}, true)
	for i, d := range testNumberResultExpressions {
		result, err := e.Evaluate(d)
		assert.NoError(t, err, "index %d", i)
		assert.Equal(t, f64d4.FromStringForced(expected[i]), result, "index %d", i)
	}
	for i, d := range testStringResultExpressions {
		result, err := e.Evaluate(d)
		assert.NoError(t, err, "index %d", i)
		assert.Equal(t, testStringResultExpected[i], result, "index %d", i)
	}

	result, err := e.Evaluate("2 > 1")
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	e = eval.NewFixedEvaluator[f64d4.Int](nil, false)
	_, err = e.Evaluate("1 / 0")
	assert.Error(t, err)
}
