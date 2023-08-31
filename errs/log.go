// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package errs

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

// StackTraceKey is the key used for logging the stack trace.
const StackTraceKey = "stack_trace"

// Log an error with a stack trace.
func Log(err error, args ...any) {
	log(context.Background(), slog.LevelError, slog.Default(), WrapTyped(err), args...)
}

// LogContext logs an error with a stack trace.
func LogContext(ctx context.Context, err error, args ...any) {
	log(ctx, slog.LevelError, slog.Default(), WrapTyped(err), args...)
}

// LogTo logs an error with a stack trace.
func LogTo(logger *slog.Logger, err error, args ...any) {
	log(context.Background(), slog.LevelError, logger, WrapTyped(err), args...)
}

// LogContextTo logs an error with a stack trace.
func LogContextTo(ctx context.Context, logger *slog.Logger, err error, args ...any) {
	log(ctx, slog.LevelError, logger, WrapTyped(err), args...)
}

// LogWithLevel logs an error with a stack trace.
func LogWithLevel(ctx context.Context, level slog.Level, logger *slog.Logger, err error, args ...any) {
	log(ctx, level, logger, WrapTyped(err), args...)
}

func log(ctx context.Context, level slog.Level, logger *slog.Logger, err *Error, args ...any) {
	if logger == nil {
		logger = slog.Default()
	}
	if !logger.Enabled(ctx, level) {
		return
	}
	r := createRecord(level, err)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = logger.Handler().Handle(ctx, r) //nolint:errcheck // Since we are in the logger, nothing we can reasonably do to log this
}

func createRecord(level slog.Level, err *Error) slog.Record {
	var pc uintptr
	var msg string
	if err != nil {
		msg = err.Message()
		if len(err.stack) != 0 {
			pc = err.stack[0]
		}
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	if err != nil {
		r.AddAttrs(slog.Any(StackTraceKey, &stackValue{err: err}))
	}
	return r
}

// LogAttrs logs an error with a stack trace.
func LogAttrs(err error, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.LevelError, slog.Default(), WrapTyped(err), attrs...)
}

// LogAttrsContext logs an error with a stack trace.
func LogAttrsContext(ctx context.Context, err error, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelError, slog.Default(), WrapTyped(err), attrs...)
}

// LogAttrsTo logs an error with a stack trace.
func LogAttrsTo(logger *slog.Logger, err error, attrs ...slog.Attr) {
	logAttrs(context.Background(), slog.LevelError, logger, WrapTyped(err), attrs...)
}

// LogAttrsContextTo logs an error with a stack trace.
func LogAttrsContextTo(ctx context.Context, logger *slog.Logger, err error, attrs ...slog.Attr) {
	logAttrs(ctx, slog.LevelError, logger, WrapTyped(err), attrs...)
}

// LogAttrsWithLevel logs an error with a stack trace.
func LogAttrsWithLevel(ctx context.Context, level slog.Level, logger *slog.Logger, err error, attrs ...slog.Attr) {
	logAttrs(ctx, level, logger, WrapTyped(err), attrs...)
}

func logAttrs(ctx context.Context, level slog.Level, logger *slog.Logger, err *Error, attrs ...slog.Attr) {
	if logger == nil {
		logger = slog.Default()
	}
	if !logger.Enabled(ctx, level) {
		return
	}
	r := createRecord(level, err)
	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = logger.Handler().Handle(ctx, r) //nolint:errcheck // Since we are in the logger, nothing we can reasonably do to log this
}

type stackValue struct {
	err StackError
}

func (v *stackValue) StackError() StackError {
	return v.err
}

func (v *stackValue) LogValue() slog.Value {
	stack := strings.Split(v.err.StackTrace(true), "\n")
	for i := 0; i < len(stack); i++ {
		stack[i] = strings.TrimSpace(stack[i])
	}
	return slog.AnyValue(stack)
}
