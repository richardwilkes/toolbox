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

func TestAppendWrapsAndFlattens(t *testing.T) {
	for _, tc := range []struct {
		base  error
		arg   error
		name  string
		count int
	}{
		{name: "chain onto error", base: errs.Newf("foo%d", 1), arg: errs.Append(nil, errs.New("foo"), errs.New("bar")), count: 3},
		{name: "plain base", base: errors.New("foo"), arg: errs.New("bar"), count: 2},
		{name: "plain errors throughout", base: errors.New("foo"), arg: errs.Append(nil, errors.New("bar"), errors.New("baz")), count: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			check.New(t).Equal(tc.count, errs.Append(tc.base, tc.arg).Count())
		})
	}
}

func TestNilWrap(t *testing.T) {
	c := check.New(t)
	c.Nil(errs.Wrap(nil))
	c.Nil(errs.WrapTyped(nil))
}

func TestDoubleWrap(t *testing.T) {
	c := check.New(t)
	errError := error(errs.New("foo"))

	c.Equal(errError, errs.Wrap(errError))

	wrappedErr := fmt.Errorf("bar: %w", errError)
	c.Equal(errError, errors.Unwrap(wrappedErr))

	c.Equal(wrappedErr, errs.Wrap(wrappedErr))
}

func TestDoubleWrapTyped(t *testing.T) {
	c := check.New(t)
	errError := errs.New("foo")

	c.Equal(errError, errs.WrapTyped(errError))

	wrappedErr := fmt.Errorf("bar: %w", errError)
	c.Equal(error(errError), errors.Unwrap(wrappedErr))

	// Unlike Wrap, an error that merely wraps an *Error is wrapped again.
	rewrittenError := errs.WrapTyped(wrappedErr)
	c.Equal(wrappedErr, errors.Unwrap(rewrittenError))
}

type customErr struct {
	value string
}

func (e *customErr) Error() string {
	return e.value
}

func TestIsAndAsSeeWrappedCause(t *testing.T) {
	c := check.New(t)
	err := errs.Wrap(os.ErrNotExist)
	c.NotNil(err)
	c.True(errors.Is(err, os.ErrNotExist))
	c.False(errors.Is(err, os.ErrClosed))

	original := &customErr{value: "err"}
	var target *customErr
	c.True(errors.As(errs.Wrap(original), &target))
	c.Equal(original, target)
}

func TestCausedByOnlyForExplicitCause(t *testing.T) {
	cause := errs.New("bar")
	for _, tc := range []struct {
		err      error
		name     string
		newlines int
	}{
		{name: "New", err: errs.New("foo"), newlines: 1},
		{name: "Wrap", err: errs.Wrap(errors.New("foo")), newlines: 1},
		{name: "NewWithCause", err: errs.NewWithCause("foo", cause), newlines: 3},
		{name: "NewWithCausef", err: errs.NewWithCausef(cause, "foo%d", 1), newlines: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			check.New(t).Equal(tc.newlines, strings.Count(tc.err.Error(), "\n"))
		})
	}
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
	c.Contains(val.Group()[0].Value.String(), "foo")
}

func TestAppendIgnoresNilAndEmpty(t *testing.T) {
	var nilErr *errs.Error
	for _, tc := range []struct {
		name  string
		base  error
		args  []error
		count int
	}{
		{name: "nils between errors", base: errs.New("foo"), args: []error{nilErr, errs.New("bar"), nil}, count: 2},
		{name: "empty base", base: &errs.Error{}, args: []error{errs.New("bar")}, count: 1},
		{name: "typed nil base", base: nilErr, args: []error{errs.New("bar")}, count: 1},
		{name: "untyped nil base", base: nil, args: []error{errs.New("bar")}, count: 1},
		{name: "typed nil arg", base: nil, args: []error{nilErr}, count: 0},
		{name: "untyped nil arg", base: nil, args: []error{nil}, count: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			check.New(t).Equal(tc.count, errs.Append(tc.base, tc.args...).Count())
		})
	}
}

func TestMultiErrorMessage(t *testing.T) {
	c := check.New(t)
	err := errs.Append(errs.New("foo"), errs.New("bar"))
	msg := err.Message()
	c.Equal(msg, "Multiple (2) errors occurred:\n- foo\n- bar")
}
