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
	"log/slog"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xslog"
)

func TestMultiHandler(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	h1 := slog.NewJSONHandler(&buf1, nil)
	h2 := slog.NewJSONHandler(&buf2, nil)
	multi := xslog.NewMultiHandler(h1, h2)
	ctx := context.Background()
	check.True(t, multi.Enabled(ctx, slog.LevelInfo))
	check.False(t, multi.Enabled(ctx, slog.LevelDebug))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "record 1", 0)
	record.Add("key1", "value1")
	check.NoError(t, multi.Handle(ctx, record))

	withAttrs := multi.WithAttrs(nil)
	check.Equal(t, multi, withAttrs)
	withAttrs = multi.WithAttrs([]slog.Attr{slog.String("extra", "data")})
	record = slog.NewRecord(time.Now(), slog.LevelInfo, "record 2", 0)
	record.Add("key2", "value2")
	check.NoError(t, withAttrs.Handle(ctx, record))

	withGroup := multi.WithGroup("")
	check.Equal(t, multi, withGroup)
	withGroup = multi.WithGroup("group1")
	record = slog.NewRecord(time.Now(), slog.LevelInfo, "record 3", 0)
	record.Add("key3", "value3")
	check.NoError(t, withGroup.Handle(ctx, record))

	check.Equal(t, buf1.String(), buf2.String())
}
