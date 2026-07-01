// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xslog"
	"github.com/richardwilkes/toolbox/v2/xterm"
)

// TestPrettyHandlerNoColorForNonTerminal verifies that a handler over a non-terminal writer (such as a log file)
// auto-detects Dumb and emits no ANSI escape sequences, while a handler with an explicit color ColorSupportOverride
// does emit them.
func TestPrettyHandlerNoColorForNonTerminal(t *testing.T) {
	c := check.New(t)

	var plain bytes.Buffer
	plainHandler := xslog.NewPrettyHandler(&plain, nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "no color please", 0)
	c.NoError(plainHandler.Handle(context.Background(), record))
	c.Contains(plain.String(), "no color please")
	c.NotContains(plain.String(), "\x1b[")

	var colored bytes.Buffer
	coloredHandler := xslog.NewPrettyHandler(&colored, &xslog.PrettyOptions{ColorSupportOverride: xterm.Color24})
	record = slog.NewRecord(time.Now(), slog.LevelInfo, "colorful", 0)
	c.NoError(coloredHandler.Handle(context.Background(), record))
	c.Contains(colored.String(), "colorful")
	c.Contains(colored.String(), "\x1b[")
}

func TestPrettyHandlerBasic(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.Add("key1", "value1")
	record.Add("elapsed", 5*time.Millisecond)
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	output := buf.String()
	c.Contains(output, "test message")
	c.Contains(output, "INF")
	c.Contains(output, `"key1":"value1"`)
	c.Contains(output, `"elapsed":"5ms"`)
}

func TestPrettyHandlerLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, &xslog.PrettyOptions{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug - 4},
	})
	for _, one := range []struct {
		prefix string
		level  slog.Level
	}{
		{prefix: "DEBUG ", level: slog.LevelDebug},
		{prefix: "INFO ", level: slog.LevelInfo},
		{prefix: "WARN ", level: slog.LevelWarn},
		{prefix: "ERROR ", level: slog.LevelError},
		{prefix: "WARN+2 ", level: slog.LevelWarn + 2},
		{prefix: "DEBUG-2 ", level: slog.LevelDebug - 2},
	} {
		t.Run(one.level.String(), func(t *testing.T) {
			record := slog.NewRecord(time.Now(), one.level, "level msg", 0)
			c := check.New(t)
			c.NoError(handler.Handle(context.Background(), record))
			c.True(strings.HasPrefix(buf.String(), one.prefix))
			buf.Reset()
		})
	}
}

func TestPrettyHandlerCallerInfo(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, &xslog.PrettyOptions{HandlerOptions: slog.HandlerOptions{AddSource: true}})
	pc, _, _, _ := runtime.Caller(0)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test caller", pc)
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	c.Contains(buf.String(), "pretty_test.go:")
}

func TestPrettyHandlerStackTrace(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	pc, _, _, _ := runtime.Caller(0)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test stack", pc)
	record.Add(errs.StackTraceKey, []string{
		"[main.main] play/main.go:32",
		"[runtime.main] runtime/proc.go:283",
	})
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	c.Contains(buf.String(), "\n    [main.main] play/main.go:32\n    [runtime.main] runtime/proc.go:283")
}

// TestPrettyHandlerDerivedStackTrace verifies that a handler derived via WithAttrs/WithGroup still emits the stack
// trace carried by its own records. The JSON handler's stack-capturing ReplaceAttr closure is bound to the original
// handler, so a derived clone must share the capture slot rather than reading its own empty copy.
func TestPrettyHandlerDerivedStackTrace(t *testing.T) {
	c := check.New(t)
	for _, derive := range []struct {
		fn   func(slog.Handler) slog.Handler
		name string
	}{
		{name: "WithAttrs", fn: func(h slog.Handler) slog.Handler {
			return h.WithAttrs([]slog.Attr{{Key: "k", Value: slog.StringValue("v")}})
		}},
		{name: "WithGroup", fn: func(h slog.Handler) slog.Handler { return h.WithGroup("g") }},
	} {
		t.Run(derive.name, func(_ *testing.T) {
			var buf bytes.Buffer
			derived := derive.fn(xslog.NewPrettyHandler(&buf, nil))
			record := slog.NewRecord(time.Now(), slog.LevelInfo, "derived stack", 0)
			record.Add(errs.StackTraceKey, []string{"[main.main] play/main.go:32"})
			c.NoError(derived.Handle(context.Background(), record))
			c.Contains(buf.String(), "\n    [main.main] play/main.go:32")
		})
	}
}

// TestPrettyHandlerDerivedStackTraceNoLeak verifies that a stack trace handled by a derived clone does not leak into a
// later, unrelated record handled by the base handler.
func TestPrettyHandlerDerivedStackTraceNoLeak(t *testing.T) {
	c := check.New(t)
	var buf bytes.Buffer
	base := xslog.NewPrettyHandler(&buf, nil)
	derived := base.WithAttrs([]slog.Attr{{Key: "k", Value: slog.StringValue("v")}})
	withStack := slog.NewRecord(time.Now(), slog.LevelInfo, "derived msg", 0)
	withStack.Add(errs.StackTraceKey, []string{"[leak.func] leak/leak.go:1"})
	c.NoError(derived.Handle(context.Background(), withStack))
	buf.Reset()
	c.NoError(base.Handle(context.Background(), slog.NewRecord(time.Now(), slog.LevelInfo, "base msg", 0)))
	c.NotContains(buf.String(), "leak/leak.go:1")
}

func TestPrettyHandlerEmptyMsg(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "", 0)
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	c.Equal(2, strings.Count(buf.String(), " | "))
}

func TestPrettyHandlerMultiLineMsg(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "first line\nsecond line", 0)
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	output := buf.String()
	c.Contains(output, "first line\n    second line")
}

func TestPrettyHandlerEnabled(t *testing.T) {
	handler := xslog.NewPrettyHandler(nil, nil)
	c := check.New(t)
	c.False(handler.Enabled(context.Background(), slog.LevelDebug))
	c.True(handler.Enabled(context.Background(), slog.LevelInfo))
}

func TestPrettyHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	h := handler.WithAttrs([]slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "key2", Value: slog.IntValue(42)},
	})
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "with attrs", 0)
	record.Add("key3", "value3")
	c := check.New(t)
	c.NoError(h.Handle(context.Background(), record))
	c.Contains(buf.String(), `{"key1":"value1","key2":42,"key3":"value3"}`)
}

func TestPrettyHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	h := handler.WithGroup("group")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "with group", 0)
	record.Add("key4", "value4")
	c := check.New(t)
	c.NoError(h.Handle(context.Background(), record))
	c.Contains(buf.String(), `{"group":{"key4":"value4"}}`)
}

func TestPrettyHandlerWithReplacer(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, &xslog.PrettyOptions{
		HandlerOptions: slog.HandlerOptions{
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == "r" {
					return slog.Attr{Key: "r", Value: slog.StringValue("replaced")}
				}
				return a
			},
		},
	})
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "replacer test", 0)
	record.Add("r", "foo")
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	c.Contains(buf.String(), `{"r":"replaced"}`)
}

func TestPrettyHandlerConcurrency(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	var wg sync.WaitGroup
	const numGoroutines = 20
	c := check.New(t)
	for i := range numGoroutines {
		wg.Go(func() {
			record := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf("concurrent test %d", i), 0)
			c.NoError(handler.Handle(context.Background(), record))
		})
	}
	wg.Wait()
	output := buf.String()
	for i := range numGoroutines {
		c.Contains(output, fmt.Sprintf("concurrent test %d", i))
	}
}
