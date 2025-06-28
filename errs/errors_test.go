// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package errs_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/errs"
)

func ExampleError() {
	var bad *int
	func() int {
		defer errs.Recovery(func(err error) { fmt.Println(err) })
		return *bad // trigger a panic due to a nil pointer dereference
	}()
	// Output: recovered from panic
	//     [github.com/richardwilkes/toolbox/v2/errs_test.ExampleError.func1] errors_test.go:27
	//     [github.com/richardwilkes/toolbox/v2/errs_test.ExampleError] errors_test.go:28
	//   Caused by: runtime error: invalid memory address or nil pointer dereference
}

func TestAppendError(t *testing.T) {
	original := errs.New("foo")
	result := errs.Append(original, errs.New("bar"))
	check.Equal(t, 2, result.Count())
}

func TestAppendToEmptyError(t *testing.T) {
	original := &errs.Error{}
	result := errs.Append(original, errs.New("bar"))
	check.Equal(t, 1, result.Count())
}

func TestAppendFlattening(t *testing.T) {
	original := errs.New("foo")
	result := errs.Append(original, errs.Append(nil, errs.New("foo"), errs.New("bar")))
	check.Equal(t, 3, result.Count())
}

func TestAppendTypedNil(t *testing.T) {
	var e *errs.Error
	result := errs.Append(e, errs.New("bar"))
	check.Equal(t, 1, result.Count())
}

func TestAppendNilError(t *testing.T) {
	var err error
	result := errs.Append(err, errs.New("bar"))
	check.Equal(t, 1, result.Count())
}

func TestAppendNilErrorArg(t *testing.T) {
	var err error
	var nilErr *errs.Error
	result := errs.Append(err, nilErr)
	check.Equal(t, 0, result.Count())
}

func TestAppendNilErrorInterfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := errs.Append(err, nilErr)
	check.Equal(t, 0, result.Count())
}

func TestAppendNonErrorError(t *testing.T) {
	original := errors.New("foo")
	result := errs.Append(original, errs.New("bar"))
	check.Equal(t, 2, result.Count())
}

func TestAppendNonErrorErrorWithAppend(t *testing.T) {
	original := errors.New("foo")
	result := errs.Append(original, errs.Append(nil, errors.New("bar")))
	check.Equal(t, 2, result.Count())
}

func TestErrorOrNil(t *testing.T) {
	var err errs.Error
	check.Nil(t, err.ErrorOrNil())
}

func TestErrorOrNilPointer(t *testing.T) {
	var err *errs.Error
	check.Nil(t, err.ErrorOrNil())
}

func TestWrap(t *testing.T) {
	notError := errors.New("foo")
	result := errs.Wrap(notError)
	check.NotNil(t, result)
	check.Equal(t, 1, strings.Count(result.Error(), "\n"))
}

func TestDoubleWrap(t *testing.T) {
	errError := error(errs.New("foo"))

	// Verify *errs.Error doesn't get wrapped again
	check.Equal(t, errError, errs.Wrap(errError))

	// Wrap the error using the standard library
	wrappedErr := fmt.Errorf("bar: %w", errError)
	check.Equal(t, errError, errors.Unwrap(wrappedErr))

	// Verify that an error with an embedded *errs.Error cause doesn't get wrapped again
	check.Equal(t, wrappedErr, errs.Wrap(wrappedErr))
}

func TestDoubleWrapTyped(t *testing.T) {
	errError := errs.New("foo")

	// Verify *errs.Error doesn't get wrapped again
	check.Equal(t, errError, errs.WrapTyped(errError))

	// Wrap the error using the standard library
	wrappedErr := fmt.Errorf("bar: %w", errError)
	check.Equal(t, error(errError), errors.Unwrap(wrappedErr))

	// It seems the best thing to do here is to wrap again
	rewrappedError := errs.WrapTyped(wrappedErr)
	check.Equal(t, wrappedErr, errors.Unwrap(rewrappedError))
}

func TestIs(t *testing.T) {
	err := errs.Wrap(os.ErrNotExist)
	check.NotNil(t, err)
	check.True(t, errors.Is(err, os.ErrNotExist))
	check.False(t, errors.Is(err, os.ErrClosed))
}

type customErr struct {
	value string
}

func (e *customErr) Error() string {
	return e.value
}

func TestAs(t *testing.T) {
	original := &customErr{value: "err"}
	wrapped := errs.Wrap(original)
	check.NotNil(t, wrapped)
	var target *customErr
	check.True(t, errors.As(wrapped, &target))
	check.Equal(t, original, target)
}

func TestNew(t *testing.T) {
	result := errs.New("foo")
	check.Equal(t, 1, strings.Count(result.Error(), "\n"))
}

func TestNewWithCause(t *testing.T) {
	cause := errs.New("bar")
	result := errs.NewWithCause("foo", cause)
	check.Equal(t, 3, strings.Count(result.Error(), "\n"))
}

func TestFormat(t *testing.T) {
	err := errs.New("test")
	check.Equal(t, "test", fmt.Sprintf("%s", err))
	check.Equal(t, `"test"`, fmt.Sprintf("%q", err))
	result := fmt.Sprintf("%v", err)
	check.Contains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestFormat]")
	check.NotContains(t, result, "[runtime.goexit]")
	result = fmt.Sprintf("%+v", err)
	check.Contains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestFormat]")
	check.Contains(t, result, "[runtime.goexit]")

	wrappedErrors := err.WrappedErrors()
	check.Equal(t, 1, len(wrappedErrors))
	check.Equal(t, "test", fmt.Sprintf("%s", wrappedErrors[0])) //nolint:gocritic // Testing %s, so necessary
	check.Equal(t, `"test"`, fmt.Sprintf("%q", wrappedErrors[0]))
	result = fmt.Sprintf("%v", wrappedErrors[0]) //nolint:gocritic // Testing %v, so necessary
	check.Contains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestFormat]")
	check.NotContains(t, result, "[runtime.goexit]")
	result = fmt.Sprintf("%+v", wrappedErrors[0])
	check.Contains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestFormat]")
	check.Contains(t, result, "[runtime.goexit]")
}

func TestWrappedErrors(t *testing.T) {
	foo := errs.New("foo")
	bar := errs.Append(foo, errs.New("bar"))
	foo2 := errs.New("foo2")
	bar2 := errs.Append(foo2, errs.New("bar2"))
	result := errs.Append(bar, bar2)
	list := result.WrappedErrors()
	check.Equal(t, 4, len(list))
	check.Equal(t, "foo", strings.SplitN(list[0].Error(), "\n", 2)[0])
	check.Equal(t, "bar", strings.SplitN(list[1].Error(), "\n", 2)[0])
	check.Equal(t, "foo2", strings.SplitN(list[2].Error(), "\n", 2)[0])
	check.Equal(t, "bar2", strings.SplitN(list[3].Error(), "\n", 2)[0])
}

func TestAlteredFilter(t *testing.T) {
	err := errs.New("test")
	result := fmt.Sprintf("%v", err)
	check.Contains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestAlteredFilter]")
	saved := errs.RuntimePrefixesToFilter
	errs.RuntimePrefixesToFilter = []string{"github.com/richardwilkes/toolbox/v2/errs_test.TestAlteredFilter"}
	result = fmt.Sprintf("%v", err)
	check.NotContains(t, result, "[github.com/richardwilkes/toolbox/v2/errs_test.TestAlteredFilter]")
	errs.RuntimePrefixesToFilter = saved
}
