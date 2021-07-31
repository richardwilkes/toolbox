// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package eval_test

import (
	"fmt"

	"github.com/richardwilkes/toolbox/eval/f64d4eval"
)

func Example() {
	e := f64d4eval.NewEvaluator(nil, true)
	result, err := e.Evaluate("1 + sqrt(2)")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// 2.4142
}
