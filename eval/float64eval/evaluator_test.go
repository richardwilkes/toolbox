// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package float64eval_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/eval/float64eval"
	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	data := []string{
		"1 + 1",
		"1.3 + 1.5",
		"1.3+1.5",
		"1.30015 + 1.5",
		"1 / 3",
		"1 / 3 + 10",
		"1 / (3 + 10)",
		"-1 / (3 + 10)",
		"1 / 0",
		"sqrt(9)",
		"sqrt(2)",
		"sqrt(cbrt(8)+7)",
		"  sqrt	(  cbrt    ( 8 ) +     7.0000 )    ",
		"$foo + $bar",
		"$foo / $bar",
		"2 + 1e-2",
		"2 + 1e2",
	}
	expected := []float64{
		2,
		2.8,
		2.8,
		2.80015,
		0.3333333333333333,
		10.333333333333334,
		0.07692307692307693,
		-0.07692307692307693,
		0,
		3,
		1.4142135623730951,
		3,
		3,
		24,
		11,
		2.01,
		102,
	}
	e := float64eval.NewEvaluator(resolver{}, true)
	for i, d := range data {
		result, err := e.Evaluate(d)
		assert.NoError(t, err, "index %d", i)
		assert.Equal(t, expected[i], result, "index %d", i)
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

	e = float64eval.NewEvaluator(nil, false)
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
