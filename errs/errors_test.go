// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	"github.com/richardwilkes/toolbox/v2/xos"
)

func ExampleError() {
	var bad *int
	func() int {
		defer xos.PanicRecovery(func(err error) { fmt.Println(err) })
		return *bad // trigger a panic due to a nil pointer dereference
	}()
	// Output: recovered from panic
	//     [github.com/richardwilkes/toolbox/v2/errs_test.ExampleError.func1] errs/errors_test.go:28
	//     [github.com/richardwilkes/toolbox/v2/errs_test.ExampleError] errs/errors_test.go:29
	//   Caused by: runtime error: invalid memory address or nil pointer dereference
}

func TestDetail(t *testing.T) {
	c := check.New(t)
	detail := errs.New("").Detail()
	c.HasPrefix(detail, "    [github.com/richardwilkes/toolbox/v2/errs_test.TestDetail] errs/errors_test.go:")
	detail = (&errs.Error{}).Detail()
	c.Equal("<no detail>", detail)
}

func TestAppendError(t *testing.T) {
	c := check.New(t)
	original := errs.New("foo")
	result := errs.Append(original, errs.New("bar"))
	c.Equal(2, result.Count())
}

func TestAppendToEmptyError(t *testing.T) {
	c := check.New(t)
	original := &errs.Error{}
	result := errs.Append(original, errs.New("bar"))
	c.Equal(1, result.Count())
}

func TestAppendFlattening(t *testing.T) {
	c := check.New(t)
	original := errs.Newf("foo%d", 1)
	result := errs.Append(original, errs.Append(nil, errs.New("foo"), errs.New("bar")))
	c.Equal(3, result.Count())
}

func TestAppendTypedNil(t *testing.T) {
	c := check.New(t)
	var e *errs.Error
	result := errs.Append(e, errs.New("bar"))
	c.Equal(1, result.Count())
}

func TestAppendNilError(t *testing.T) {
	c := check.New(t)
	var err error
	result := errs.Append(err, errs.New("bar"))
	c.Equal(1, result.Count())
}

func TestAppendNilErrorArg(t *testing.T) {
	c := check.New(t)
	var err error
	var nilErr *errs.Error
	result := errs.Append(err, nilErr)
	c.Equal(0, result.Count())
}

func TestAppendNilErrorInterfaceArg(t *testing.T) {
	c := check.New(t)
	var err error
	var nilErr error
	result := errs.Append(err, nilErr)
	c.Equal(0, result.Count())
}

func TestAppendNonErrorError(t *testing.T) {
	c := check.New(t)
	original := errors.New("foo")
	result := errs.Append(original, errs.New("bar"))
	c.Equal(2, result.Count())
}

func TestAppendNonErrorErrorWithAppend(t *testing.T) {
	c := check.New(t)
	original := errors.New("foo")
	result := errs.Append(original, errs.Append(nil, errors.New("bar"), errors.New("baz")))
	c.Equal(3, result.Count())
}

func TestErrorOrNil(t *testing.T) {
	c := check.New(t)
	var err errs.Error
	c.Nil(err.ErrorOrNil())
}

func TestErrorOrNilPointer(t *testing.T) {
	c := check.New(t)
	var err *errs.Error
	c.Nil(err.ErrorOrNil())
}

func TestWrap(t *testing.T) {
	c := check.New(t)
	notError := errors.New("foo")
	result := errs.Wrap(notError)
	c.NotNil(result)
	c.Equal(1, strings.Count(result.Error(), "\n"))
}

func TestNilWrap(t *testing.T) {
	c := check.New(t)
	c.Nil(errs.Wrap(nil))
	c.Nil(errs.WrapTyped(nil))
}

func TestDoubleWrap(t *testing.T) {
	c := check.New(t)
	errError := error(errs.New("foo"))

	// Verify *errs.Error doesn't get wrapped again
	c.Equal(errError, errs.Wrap(errError))

	// Wrap the error using the standard library
	wrappedErr := fmt.Errorf("bar: %w", errError)
	c.Equal(errError, errors.Unwrap(wrappedErr))

	// Verify that an error with an embedded *errs.Error cause doesn't get wrapped again
	c.Equal(wrappedErr, errs.Wrap(wrappedErr))
}

func TestDoubleWrapTyped(t *testing.T) {
	c := check.New(t)
	errError := errs.New("foo")

	// Verify *errs.Error doesn't get wrapped again
	c.Equal(errError, errs.WrapTyped(errError))

	// Wrap the error using the standard library
	wrappedErr := fmt.Errorf("bar: %w", errError)
	c.Equal(error(errError), errors.Unwrap(wrappedErr))

	// It seems the best thing to do here is to wrap again
	rewrittenError := errs.WrapTyped(wrappedErr)
	c.Equal(wrappedErr, errors.Unwrap(rewrittenError))
}

func TestIs(t *testing.T) {
	c := check.New(t)
	err := errs.Wrap(os.ErrNotExist)
	c.NotNil(err)
	c.True(errors.Is(err, os.ErrNotExist))
	c.False(errors.Is(err, os.ErrClosed))
}

type customErr struct {
	value string
}

func (e *customErr) Error() string {
	return e.value
}

func TestAs(t *testing.T) {
	c := check.New(t)
	original := &customErr{value: "err"}
	wrapped := errs.Wrap(original)
	c.NotNil(wrapped)
	var target *customErr
	c.True(errors.As(wrapped, &target))
	c.Equal(original, target)
}

func TestNew(t *testing.T) {
	c := check.New(t)
	result := errs.New("foo")
	c.Equal(1, strings.Count(result.Error(), "\n"))
}

func TestNewWithCause(t *testing.T) {
	c := check.New(t)
	cause := errs.New("bar")
	result := errs.NewWithCause("foo", cause)
	c.Equal(3, strings.Count(result.Error(), "\n"))
	result = errs.NewWithCausef(cause, "foo%d", 1)
	c.Equal(3, strings.Count(result.Error(), "\n"))
}

func TestFormat(t *testing.T) {
	c := check.New(t)
	err := errs.New("test")
	c.Equal("test", fmt.Sprintf("%s", err))
	c.Equal(`"test"`, fmt.Sprintf("%q", err))
	result := fmt.Sprintf("%v", err)
	c.Contains(result, "/errs_test.TestFormat] ")

	wrappedErrors := err.WrappedErrors()
	c.Equal(1, len(wrappedErrors))
	c.Equal("test", fmt.Sprintf("%s", wrappedErrors[0])) //nolint:gocritic // Testing %s, so necessary
	c.Equal(`"test"`, fmt.Sprintf("%q", wrappedErrors[0]))
	result = fmt.Sprintf("%v", wrappedErrors[0]) //nolint:gocritic // Testing %v, so necessary
	c.Contains(result, "/errs_test.TestFormat] ")
}

func TestWrappedErrors(t *testing.T) {
	c := check.New(t)
	foo := errs.New("foo")
	bar := errs.Append(foo, errs.New("bar"))
	foo2 := errs.New("foo2")
	bar2 := errs.Append(foo2, errs.New("bar2"))
	result := errs.Append(bar, bar2)
	list := result.WrappedErrors()
	c.Equal(4, len(list))
	c.Equal("foo", strings.SplitN(list[0].Error(), "\n", 2)[0])
	c.Equal("bar", strings.SplitN(list[1].Error(), "\n", 2)[0])
	c.Equal("foo2", strings.SplitN(list[2].Error(), "\n", 2)[0])
	c.Equal("bar2", strings.SplitN(list[3].Error(), "\n", 2)[0])
}

func TestCloneWithPrefixMessage(t *testing.T) {
	c := check.New(t)
	err := errs.New("foo")
	cloned := err.CloneWithPrefixMessage("prefix: ")
	c.NotNil(cloned)
	c.Contains(cloned.Error(), "prefix: foo")
	c.NotEqual(err, cloned)
}

func TestRawStackTrace(t *testing.T) {
	c := check.New(t)
	err := errs.New("foo")
	raw := err.RawStackTrace()
	c.NotNil(raw)
	c.True(len(raw) > 0)
}

func TestEmptyMethod(t *testing.T) {
	c := check.New(t)
	var e *errs.Error
	c.True(e.ErrorOrNil() == nil)
	emptyErr := &errs.Error{}
	c.True(emptyErr.ErrorOrNil() == nil)
	nonEmptyErr := errs.New("foo")
	c.False(nonEmptyErr.ErrorOrNil() == nil)
}

func TestLogValue(t *testing.T) {
	c := check.New(t)
	err := errs.New("foo")
	val := err.LogValue()
	c.Equal("Group", val.Kind().String())
	// Check that the group contains the error message
	c.Contains(val.Group()[0].Value.String(), "foo")
}

func TestAppendWithNilMix(t *testing.T) {
	c := check.New(t)
	err1 := errs.New("foo")
	var nilErr *errs.Error
	result := errs.Append(err1, nilErr, errs.New("bar"), nil)
	c.Equal(2, result.Count())
}

func TestMultiErrorMessage(t *testing.T) {
	c := check.New(t)
	err := errs.Append(errs.New("foo"), errs.New("bar"))
	msg := err.Message()
	c.Equal(msg, "Multiple (2) errors occurred:\n- foo\n- bar")
}
