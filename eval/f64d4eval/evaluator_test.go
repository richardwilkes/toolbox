// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f64d4eval_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/eval/f64d4eval"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	data := []string{
		"1 + 1", "2",
		"1.3 + 1.5", "2.8",
		"1.30015 + 1.5", "2.8001",
		"1 / 3", "0.3333",
		"1 / 3 + 10", "10.3333",
		"1 / (3 + 10)", "0.0769",
		"-1 / (3 + 10)", "-0.0769",
		"1 / 0", "0",
		"sqrt(9)", "3",
		"sqrt(2)", "1.4142",
		"sqrt(cbrt(8)+7)", "3",
		"  sqrt	(  cbrt    ( 8 ) +     7.0000 )    ", "3",
		"$foo + $bar", "24",
		"$foo / $bar", "11",
		"2 + 1e-2", "2.01",
		"2 + 1e2", "102",
	}
	e := f64d4eval.NewEvaluator(resolver{}, true)
	for i := 0; i < len(data); i += 2 {
		result, err := e.Evaluate(data[i])
		assert.NoError(t, err, "index %d", i)
		assert.Equal(t, fixed.F64d4FromStringForced(data[i+1]), result, "index %d", i)
	}

	data = []string{
		"foo + bar", "foobar",
		"foo +               \n    bar", "foobar",
		"$other", "22 + 2",
		"if($foo > $bar, yes, no)", "yes",
		"if($foo < $bar, yes, no)", "no",
	}
	for i := 0; i < len(data); i += 2 {
		result, err := e.Evaluate(data[i])
		assert.NoError(t, err, "index %d", i)
		assert.Equal(t, data[i+1], result, "index %d", i)
	}

	result, err := e.Evaluate("2 > 1")
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	e = f64d4eval.NewEvaluator(nil, false)
	_, err = e.Evaluate("1 / 0")
	assert.Error(t, err)
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
