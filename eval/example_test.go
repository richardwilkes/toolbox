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
