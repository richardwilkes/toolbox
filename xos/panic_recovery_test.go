// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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
	}()
	c := check.New(t)
	c.False(called)
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
	c := check.New(t)
	c.Contains(msg, "recovered from panic")
	c.Contains(msg, "test panic")
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
	c := check.New(t)
	c.Contains(msg, "recovered from panic")
	c.Contains(msg, "original panic")
	c.Contains(msg, "handler panic")
}

func TestPanicRecovery_ErrorWrapping(t *testing.T) {
	originalErr := errors.New("original error")
	for _, tc := range []struct {
		value any
		name  string
		want  string
	}{
		{name: "error", value: originalErr, want: originalErr.Error()},
		{name: "string", value: "string panic", want: "string panic"},
		{name: "int", value: 42, want: "42"},
		{name: "nil", value: nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var capturedErr error
			func() {
				defer xos.PanicRecovery(func(err error) { capturedErr = err })
				panic(tc.value)
			}()
			c := check.New(t)
			c.NotNil(capturedErr)
			c.NotEqual(originalErr, capturedErr)
			c.Contains(capturedErr.Error(), "recovered from panic")
			c.Contains(capturedErr.Error(), tc.want)
		})
	}
}
