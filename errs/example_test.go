package errs_test

import (
	"errors"
	"fmt"

	"github.com/richardwilkes/gokit/errs"
)

func ExampleNewWithCause() {
	err := errors.New("fake error")
	fmt.Println(errs.NewWithCause("This is a wrapped error", err))
	// Output:
	// This is a wrapped error
	//     [github.com/richardwilkes/gokit/errs_test.ExampleNewWithCause] example_test.go:12
	//     [main.main] github.com/richardwilkes/gokit/errs/_test/_testmain.go:76
	//   Caused by: fake error
}
