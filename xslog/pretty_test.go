// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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
)

func TestPrettyHandlerBasic(t *testing.T) {
	var buf bytes.Buffer
	handler := xslog.NewPrettyHandler(&buf, nil)
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.Add("key1", "value1")
	c := check.New(t)
	c.NoError(handler.Handle(context.Background(), record))
	output := buf.String()
	c.Contains(output, "test message")
	c.Contains(output, "INF")
	c.Contains(output, `{"key1":"value1"}`)
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
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			record := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf("concurrent test %d", id), 0)
			c.NoError(handler.Handle(context.Background(), record))
		}(i)
	}
	wg.Wait()
	output := buf.String()
	for i := range numGoroutines {
		c.Contains(output, fmt.Sprintf("concurrent test %d", i))
	}
}
