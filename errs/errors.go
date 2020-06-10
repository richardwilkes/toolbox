// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package errs implements a detailed error object that provides stack traces
// with source locations, along with nested causes, if any.
package errs

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	_ ErrorWrapper = &Error{}
	_ StackError   = &Error{}
)

// ErrorWrapper contains methods for interacting with the wrapped errors.
type ErrorWrapper interface {
	error
	Count() int
	WrappedErrors() []error
}

// StackError contains methods with the stack trace and message.
type StackError interface {
	error
	Message() string
	Detail(trimRuntime bool) string
	StackTrace(trimRuntime bool) string
}

// Error holds the detailed error message.
type Error struct {
	errors []detail
}

// Wrap an error and turn it into a detailed error. If error is already a
// detailed error or nil, it will be returned as-is.
func Wrap(cause error) error {
	// Note: even though this function body is identical to WrapTyped(), do
	// not make a call to it here, otherwise the code that trims the stack
	// will be incorrect.
	if cause == nil {
		return nil
	}
	if err, ok := cause.(*Error); ok {
		return err
	}
	return &Error{
		errors: []detail{
			{
				message: cause.Error(),
				stack:   callStack(),
				cause: cause,
			},
		},
	}
}

// WrapTyped wraps an error and turns it into a detailed error. If error is
// already a detailed error or nil, it will be returned as-is. This method
// returns the error as an *Error. Use Wrap() to receive a generic error.
func WrapTyped(cause error) *Error {
	// Note: even though this function body is identical to Wrap(), do not
	// make a call to it here, otherwise the code that trims the stack will
	// be incorrect.
	if cause == nil {
		return nil
	}
	if err, ok := cause.(*Error); ok {
		return err
	}
	return &Error{
		errors: []detail{
			{
				message: cause.Error(),
				stack:   callStack(),
				cause: cause,
			},
		},
	}
}

// New creates a new detailed error with the 'message'.
func New(message string) *Error {
	return &Error{
		errors: []detail{
			{
				message: message,
				stack:   callStack(),
			},
		},
	}
}

// Newf creates a new detailed error using fmt.Sprintf() to format the
// message.
func Newf(format string, v ...interface{}) *Error {
	return New(fmt.Sprintf(format, v...))
}

// NewWithCause creates a new detailed error with the 'message' and underlying
// 'cause'.
func NewWithCause(message string, cause error) *Error {
	return &Error{
		errors: []detail{
			{
				message: message,
				stack:   callStack(),
				cause:   cause,
			},
		},
	}
}

// NewfWithCause creates a new detailed error with an underlying 'cause' and
// using fmt.Sprintf() to format the message.
//
// Deprecated: Use NewWithCausef instead.
func NewfWithCause(cause error, format string, v ...interface{}) *Error {
	return NewWithCause(fmt.Sprintf(format, v...), cause)
}

// NewWithCausef creates a new detailed error with an underlying 'cause' and
// using fmt.Sprintf() to format the message.
func NewWithCausef(cause error, format string, v ...interface{}) *Error {
	return NewWithCause(fmt.Sprintf(format, v...), cause)
}

// Append one or more errors to an existing error. err may be nil.
func Append(err error, errs ...error) *Error {
	switch err := err.(type) {
	case *Error:
		if err == nil {
			err = &Error{}
		}
		for _, one := range errs {
			switch typedErr := one.(type) {
			case *Error:
				if typedErr != nil {
					err.errors = append(err.errors, typedErr.errors...)
				}
			default:
				if typedErr != nil {
					err.errors = append(err.errors, detail{
						message: typedErr.Error(),
						stack:   callStack(),
					})
				}
			}
		}
		return err
	default:
		if err == nil {
			return Append(&Error{}, errs...)
		}
		all := make([]error, 0, len(errs)+1)
		all = append(all, err)
		all = append(all, errs...)
		return Append(&Error{}, all...)
	}
}

func callStack() []uintptr {
	var pcs [512]uintptr
	n := runtime.Callers(3, pcs[:])
	cs := make([]uintptr, n)
	copy(cs, pcs[:n])
	return cs
}

// Count returns the number of contained errors, not including causes.
func (d *Error) Count() int {
	return len(d.errors)
}

// Message returns the message attached to this error.
func (d *Error) Message() string {
	switch len(d.errors) {
	case 0:
		return ""
	case 1:
		return d.errors[0].message
	default:
		var buffer strings.Builder
		buffer.WriteString(fmt.Sprintf("Multiple (%d) errors occurred:", len(d.errors)))
		for _, one := range d.errors {
			buffer.WriteString("\n- ")
			buffer.WriteString(one.message)
		}
		return buffer.String()
	}
}

// Error implements the error interface.
func (d Error) Error() string {
	return d.Detail(true)
}

// Detail returns the fully detailed error message, which includes the primary
// message, the call stack, and potentially one or more chained causes.
func (d *Error) Detail(trimRuntime bool) string {
	switch len(d.errors) {
	case 0:
		return ""
	case 1:
		return d.errors[0].detail(true, trimRuntime)
	default:
		return d.Message() + d.errors[0].detail(false, trimRuntime)
	}
}

// StackTrace returns just the stack trace portion of the message.
func (d *Error) StackTrace(trimRuntime bool) string {
	if len(d.errors) == 0 {
		return ""
	}
	return d.errors[0].detail(false, trimRuntime)
}

// RawStackTrace returns the raw call stack pointers.
func (d *Error) RawStackTrace() []uintptr {
	if len(d.errors) == 0 {
		return nil
	}
	return d.errors[0].StackTrace()
}

// ErrorOrNil returns an error interface if this Error represents one or more
// errors, or nil if it is empty.
func (d *Error) ErrorOrNil() error {
	if d == nil || len(d.errors) == 0 {
		return nil
	}
	return d
}

// WrappedErrors returns the contained errors.
func (d *Error) WrappedErrors() []error {
	result := make([]error, len(d.errors))
	for i, one := range d.errors {
		result[i] = one
	}
	return result
}

// Unwrap implements errors.Unwrap and returns the underlying cause, if any.
func (d *Error) Unwrap() error {
	if len(d.errors) == 0 {
		return nil
	}
	return d.errors[0].Cause()
}

// Format implements the fmt.Formatter interface.
//
// Supported formats:
//   - "%s"  Just the message
//   - "%q"  Just the message, but quoted
//   - "%v"  The message plus a stack trace, trimmed of golang runtime calls
//   - "%+v" The message plus a stack trace
func (d *Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		// noinspection GoUnhandledErrorResult
		state.Write([]byte(d.Detail(!state.Flag('+'))))
	case 's':
		// noinspection GoUnhandledErrorResult
		state.Write([]byte(d.Message()))
	case 'q':
		fmt.Fprintf(state, "%q", d.Message())
	}
}
