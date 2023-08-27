// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
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
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	_ ErrorWrapper   = &Error{}
	_ StackError     = &Error{}
	_ fmt.Formatter  = &Error{}
	_ slog.LogValuer = &Error{}
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
	message string
	cause   error
	stack   []uintptr
	next    *Error
	wrapped bool
}

// Wrap an error and turn it into a detailed error. If error is already a detailed error or nil, it will be returned
// as-is.
func Wrap(cause error) error {
	if cause == nil {
		return nil
	}
	var errorPtr *Error
	if errors.As(cause, &errorPtr) {
		return cause
	}
	return &Error{
		message: cause.Error(),
		stack:   callStack(),
		cause:   cause,
		wrapped: true,
	}
}

// WrapTyped wraps an error and turns it into a detailed error. If error is already a detailed error or nil, it will be
// returned as-is. This method returns the error as an *Error. Use Wrap() to receive a generic error.
func WrapTyped(cause error) *Error {
	if cause == nil {
		return nil
	}
	// Intentionally not checking to see if there is a deeper wrapped *Error as the error must be wrapped again in order
	// to avoid losing information and still return a *Error
	//nolint:errorlint // See note above
	if err, ok := cause.(*Error); ok {
		return err
	}
	return &Error{
		message: cause.Error(),
		stack:   callStack(),
		cause:   cause,
		wrapped: true,
	}
}

// New creates a new detailed error with the 'message'.
func New(message string) *Error {
	return &Error{
		message: message,
		stack:   callStack(),
	}
}

// Newf creates a new detailed error using fmt.Sprintf() to format the message.
func Newf(format string, v ...any) *Error {
	return New(fmt.Sprintf(format, v...))
}

// NewWithCause creates a new detailed error with the 'message' and underlying 'cause'.
func NewWithCause(message string, cause error) *Error {
	return &Error{
		message: message,
		stack:   callStack(),
		cause:   cause,
	}
}

// NewWithCausef creates a new detailed error with an underlying 'cause' and using fmt.Sprintf() to format the message.
func NewWithCausef(cause error, format string, v ...any) *Error {
	return NewWithCause(fmt.Sprintf(format, v...), cause)
}

// Append one or more errors to an existing error. err may be nil.
func Append(err error, errs ...error) *Error {
	//nolint:errorlint // Explicitly only want to look at this exact error and not things wrapped inside it
	switch e := err.(type) {
	case *Error:
		var root *Error
		if !e.empty() {
			root = e
			for e.next != nil {
				e = e.next
			}
		} else {
			e = nil
		}
		for _, one := range errs {
			var next *Error
			//nolint:errorlint // Explicitly only want to look at this exact error and not things wrapped inside it
			switch typedErr := one.(type) {
			case *Error:
				if !typedErr.empty() {
					n := *typedErr
					localRoot := &n
					next = localRoot
					for next.next != nil {
						copied := *next.next
						next.next = &copied
						next = next.next
					}
					next = localRoot
				}
			default:
				if typedErr != nil {
					next = &Error{
						message: typedErr.Error(),
						stack:   callStack(),
						cause:   typedErr,
						wrapped: true,
					}
				}
			}
			if next != nil {
				if e == nil {
					root = next
				} else {
					e.next = next
				}
				e = next
			}
		}
		return root
	default:
		if e == nil {
			if len(errs) == 0 {
				return nil
			}
			return Append(errs[0], errs[1:]...)
		}
		return Append(Wrap(e), errs...)
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
func (e *Error) Count() int {
	count := 0
	err := e
	for err != nil {
		if !err.empty() {
			count++
		}
		err = err.next
	}
	return count
}

// Message returns the message attached to this error.
func (e *Error) Message() string {
	if e.next == nil {
		return e.message
	}
	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf("Multiple (%d) errors occurred:", e.Count()))
	err := e
	for err != nil {
		buffer.WriteString("\n- ")
		buffer.WriteString(err.message)
		err = err.next
	}
	return buffer.String()
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Detail(true)
}

// Detail returns the fully detailed error message, which includes the primary message, the call stack, and potentially
// one or more chained causes. Note that any included stack trace will be only for the first error in the case where
// multiple errors were accumulated into one via calls to .Append().
func (e *Error) Detail(trimRuntime bool) string {
	return strings.Join([]string{e.Message(), e.StackTrace(trimRuntime)}, "\n")
}

// StackTrace returns just the stack trace portion of the message.
func (e *Error) StackTrace(trimRuntime bool) string {
	var buffer strings.Builder
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		if frame.Function != "" {
			if trimRuntime && (strings.HasPrefix(frame.Function, "runtime.") ||
				strings.HasPrefix(frame.Function, "testing.") ||
				strings.HasPrefix(frame.Function, "github.com/richardwilkes/toolbox/errs.") ||
				(frame.Function == "main.main" && frame.File == "_testmain.go")) {
				continue
			}
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			buffer.WriteString("    [")
			buffer.WriteString(frame.Function)
			buffer.WriteString("] ")
			file := frame.File
			if i := strings.Index(file, "."); i != -1 {
				for i > 0 && file[i] != os.PathSeparator {
					i--
				}
				if i > 0 {
					file = file[i+1:]
				}
				if i = strings.LastIndexByte(file, os.PathSeparator); i != -1 {
					path := file[:i]
					offset := i + 1
					if i = strings.LastIndexByte(path, os.PathSeparator); i != -1 {
						if path[i+1:] == "_obj" {
							path = path[:i]
						}
					}
					if strings.HasPrefix(frame.Function, path) {
						file = file[offset:]
					}
				}
			}
			buffer.WriteString(file)
			buffer.WriteByte(':')
			buffer.WriteString(strconv.Itoa(frame.Line))
		}
		if !more {
			break
		}
	}
	if e.cause != nil && !e.wrapped {
		buffer.WriteString("\n  Caused by: ")
		//nolint:errorlint // Explicitly only want to look at this exact error and not things wrapped inside it
		if detailed, ok := e.cause.(*Error); ok {
			buffer.WriteString(detailed.Detail(trimRuntime))
		} else {
			buffer.WriteString(e.cause.Error())
		}
	}
	return buffer.String()
}

// RawStackTrace returns the raw call stack pointers for the first error within this error.
func (e *Error) RawStackTrace() []uintptr {
	return e.stack
}

// ErrorOrNil returns an error interface if this Error represents one or more errors, or nil if it is empty.
func (e *Error) ErrorOrNil() error {
	if e.empty() {
		return nil
	}
	return e
}

func (e *Error) empty() bool {
	return e == nil || (e.message == "" && e.stack == nil && e.cause == nil && e.next == nil)
}

// WrappedErrors returns the contained errors.
func (e *Error) WrappedErrors() []error {
	result := make([]error, 0, e.Count())
	err := e
	for err != nil {
		eCopy := *err
		eCopy.next = nil
		result = append(result, &eCopy)
		err = err.next
	}
	return result
}

// Unwrap implements errors.Unwrap and returns the underlying cause, if any.
func (e *Error) Unwrap() error {
	return e.cause
}

// Format implements the fmt.Formatter interface.
//
// Supported formats:
//   - "%s"  Just the message
//   - "%q"  Just the message, but quoted
//   - "%v"  The message plus a stack trace, trimmed of golang runtime calls
//   - "%+v" The message plus a stack trace
func (e *Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = state.Write([]byte(e.Detail(!state.Flag('+'))))
	case 's':
		_, _ = state.Write([]byte(e.Message()))
	case 'q':
		_, _ = fmt.Fprintf(state, "%q", e.Message())
	}
}

// LogValue implements the slog.LogValuer interface.
func (e *Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String(slog.MessageKey, e.Message()),
		slog.Any(StackLogKey, &stackValue{err: e}),
	)
}
