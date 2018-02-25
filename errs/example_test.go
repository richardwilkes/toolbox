package errs_test

import (
	"errors"
	"fmt"

	"github.com/richardwilkes/toolbox/errs"
)

func ExampleNewWithCause() {
	err := errors.New("fake error")
	fmt.Println(errs.NewWithCause("This is a wrapped error", err))
	// Output:
	// This is a wrapped error
	//     [github.com/richardwilkes/toolbox/errs_test.ExampleNewWithCause] example_test.go:12
	//     [main.main] _testmain.go:74
	//   Caused by: fake error
}
