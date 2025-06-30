// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xos"
)

func TestPanicRecovery_NoPanic(t *testing.T) {
	var called bool
	func() {
		defer xos.PanicRecovery(func(_ error) { called = true })
		// Normal execution, no panic
	}()
	check.False(t, called)
}

func TestPanicRecovery_PanicWithError(t *testing.T) {
	var capturedErr error
	func() {
		defer xos.PanicRecovery(func(err error) { capturedErr = err })
		panic(errors.New("original error"))
	}()
	check.NotNil(t, capturedErr)
	msg := capturedErr.Error()
	check.Contains(t, msg, "recovered from panic")
	check.Contains(t, msg, "original error")
}

func TestPanicRecovery_PanicWithString(t *testing.T) {
	var capturedErr error
	func() {
		defer xos.PanicRecovery(func(err error) { capturedErr = err })
		panic("string panic")
	}()
	check.NotNil(t, capturedErr)
	msg := capturedErr.Error()
	check.Contains(t, msg, "recovered from panic")
	check.Contains(t, msg, "string panic")
}

func TestPanicRecovery_PanicWithInt(t *testing.T) {
	var capturedErr error
	func() {
		defer xos.PanicRecovery(func(err error) { capturedErr = err })
		panic(42)
	}()
	check.NotNil(t, capturedErr)
	msg := capturedErr.Error()
	check.Contains(t, msg, "recovered from panic")
	check.Contains(t, msg, "42")
}

func TestPanicRecovery_NilHandler(t *testing.T) {
	oldLogger := slog.Default()
	defer func() { slog.SetDefault(oldLogger) }()
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	func() {
		defer xos.PanicRecovery(nil)
		panic("test panic")
	}()
	msg := buf.String()
	check.Contains(t, msg, "recovered from panic")
	check.Contains(t, msg, "test panic")
}

func TestPanicRecovery_HandlerPanics(t *testing.T) {
	oldLogger := slog.Default()
	defer func() { slog.SetDefault(oldLogger) }()
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	func() {
		defer xos.PanicRecovery(func(err error) { panic(errs.NewWithCause("handler panic", err)) })
		panic("original panic")
	}()
	msg := buf.String()
	check.Contains(t, msg, "recovered from panic")
	check.Contains(t, msg, "original panic")
	check.Contains(t, msg, "handler panic")
}

func TestPanicRecovery_PanicWithNil(t *testing.T) {
	var capturedErr error
	func() {
		defer xos.PanicRecovery(func(err error) { capturedErr = err })
		panic(nil) //nolint:govet // Intentionally panicking with nil
	}()
	check.NotNil(t, capturedErr)
	check.Contains(t, capturedErr.Error(), "recovered from panic")
}

func TestPanicRecovery_ErrorWrapping(t *testing.T) {
	originalErr := errors.New("original error")
	var capturedErr error
	func() {
		defer xos.PanicRecovery(func(err error) { capturedErr = err })
		panic(originalErr)
	}()
	check.NotNil(t, capturedErr)
	check.NotEqual(t, originalErr, capturedErr)
	check.Contains(t, capturedErr.Error(), "recovered from panic")
	check.Contains(t, capturedErr.Error(), originalErr.Error())
}
