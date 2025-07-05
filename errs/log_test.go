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
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/errs"
)

type recordCatcher struct {
	records []slog.Record
}

func (h *recordCatcher) Enabled(_ context.Context, level slog.Level) bool {
	return level > slog.LevelDebug
}

func (h *recordCatcher) Handle(_ context.Context, r slog.Record) error { //nolint:gocritic // Can't change the API
	h.records = append(h.records, r)
	return nil
}

func (h *recordCatcher) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *recordCatcher) WithGroup(_ string) slog.Handler {
	return h
}

func newLoggerWithCatcher() (*slog.Logger, *recordCatcher) {
	catcher := &recordCatcher{}
	logger := slog.New(catcher)
	return logger, catcher
}

func TestLogFunctions(t *testing.T) {
	c := check.New(t)
	logger, catcher := newLoggerWithCatcher()

	err := errs.New("test error")
	errs.LogTo(logger, err)
	c.Equal(1, len(catcher.records))
	c.Contains(catcher.records[0].Message, "test error")

	catcher.records = nil
	errs.LogContextTo(context.Background(), logger, err)
	c.Equal(1, len(catcher.records))

	catcher.records = nil
	errs.LogWithLevel(context.Background(), slog.LevelWarn, logger, err)
	c.Equal(1, len(catcher.records))
	c.Equal(slog.LevelWarn, catcher.records[0].Level)
	errs.LogWithLevel(context.Background(), slog.LevelDebug, logger, err)
	c.Equal(1, len(catcher.records))

	savedLogger := slog.Default()
	defer slog.SetDefault(savedLogger)
	slog.SetDefault(logger)
	catcher.records = nil
	errs.Log(err)
	c.Equal(1, len(catcher.records))

	catcher.records = nil
	errs.LogContext(context.Background(), err)
	c.Equal(1, len(catcher.records))
}

func TestLogAttrsFunctions(t *testing.T) {
	c := check.New(t)
	logger, catcher := newLoggerWithCatcher()

	err := errs.New("attr error")
	attr := slog.String("foo", "bar")
	errs.LogAttrsTo(logger, err, attr)
	c.Equal(1, len(catcher.records))
	found := false
	catcher.records[0].Attrs(func(a slog.Attr) bool {
		if a.Key == "foo" && a.Value.String() == "bar" {
			found = true
		}
		return true
	})
	c.True(found)

	catcher.records = nil
	errs.LogAttrsContextTo(context.Background(), logger, err, attr)
	c.Equal(1, len(catcher.records))

	catcher.records = nil
	errs.LogAttrsWithLevel(context.Background(), slog.LevelWarn, logger, err, attr)
	c.Equal(1, len(catcher.records))
	c.Equal(slog.LevelWarn, catcher.records[0].Level)
	errs.LogAttrsWithLevel(context.Background(), slog.LevelDebug, logger, err, attr)
	c.Equal(1, len(catcher.records))

	savedLogger := slog.Default()
	defer slog.SetDefault(savedLogger)
	slog.SetDefault(logger)
	catcher.records = nil
	errs.LogAttrs(err, attr)
	c.Equal(1, len(catcher.records))

	catcher.records = nil
	errs.LogAttrsContext(context.Background(), err, attr)
	c.Equal(1, len(catcher.records))
}

func TestLogNilError(t *testing.T) {
	c := check.New(t)
	logger, catcher := newLoggerWithCatcher()
	errs.LogTo(logger, nil)
	c.Equal(1, len(catcher.records))
	c.Equal("", catcher.records[0].Message)

	savedLogger := slog.Default()
	defer slog.SetDefault(savedLogger)
	slog.SetDefault(logger)
	catcher.records = nil
	errs.Log(nil)
	c.Equal(1, len(catcher.records))
	c.Equal("", catcher.records[0].Message)
}

func TestStackTraceLogging(t *testing.T) {
	c := check.New(t)
	err := errs.New("stack error")
	logger, catcher := newLoggerWithCatcher()
	errs.LogTo(logger, err)
	c.Equal(1, len(catcher.records))
	found := false
	catcher.records[0].Attrs(func(a slog.Attr) bool {
		if a.Key == errs.StackTraceKey {
			lines, ok := a.Value.Resolve().Any().([]string)
			if ok && len(lines) > 0 && strings.Contains(lines[0], "TestStackTraceLogging") {
				found = true
			}
		}
		return true
	})
	c.True(found)
}

func TestLogToNil(t *testing.T) {
	c := check.New(t)
	logger, catcher := newLoggerWithCatcher()
	savedLogger := slog.Default()
	defer slog.SetDefault(savedLogger)
	slog.SetDefault(logger)
	errs.LogTo(nil, nil)
	c.Equal(1, len(catcher.records))
	errs.LogContextTo(nil, nil, nil) //nolint:staticcheck // This is for testing only
	c.Equal(2, len(catcher.records))
	attr := slog.String("foo", "bar")
	errs.LogAttrsTo(nil, nil, attr)
	c.Equal(3, len(catcher.records))
	found := false
	catcher.records[2].Attrs(func(a slog.Attr) bool {
		if a.Key == "foo" && a.Value.String() == "bar" {
			found = true
		}
		return true
	})
	c.True(found)
	attr = slog.String("foo2", "bar2")
	errs.LogAttrsContextTo(nil, nil, nil, attr) //nolint:staticcheck // This is for testing only
	c.Equal(4, len(catcher.records))
	found = false
	catcher.records[3].Attrs(func(a slog.Attr) bool {
		if a.Key == "foo2" && a.Value.String() == "bar2" {
			found = true
		}
		return true
	})
	c.True(found)
}
