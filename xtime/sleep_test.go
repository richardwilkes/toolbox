// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xtime"
)

func TestSleepCompletesNormally(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		d          time.Duration
		minElapsed time.Duration
		maxElapsed time.Duration
	}{
		{d: 50 * time.Millisecond, minElapsed: 50 * time.Millisecond, maxElapsed: 70 * time.Millisecond},
		{d: 0, maxElapsed: 10 * time.Millisecond},
		{d: -10 * time.Millisecond, maxElapsed: 10 * time.Millisecond},
	} {
		start := time.Now()
		err := xtime.Sleep(context.Background(), tc.d)
		elapsed := time.Since(start)
		c.NoError(err)
		c.True(elapsed >= tc.minElapsed && elapsed < tc.maxElapsed, "Sleep(%v) took %v, expected [%v, %v)", tc.d,
			elapsed, tc.minElapsed, tc.maxElapsed)
	}
}

func TestSleepInterruptedByContext(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		wantErr    error
		newCtx     func() (context.Context, context.CancelFunc)
		name       string
		d          time.Duration
		minElapsed time.Duration
		maxElapsed time.Duration
	}{
		{
			wantErr: context.Canceled,
			newCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				time.AfterFunc(10*time.Millisecond, cancel)
				return ctx, cancel
			},
			name:       "canceled while sleeping",
			d:          100 * time.Millisecond,
			minElapsed: 8 * time.Millisecond,
			maxElapsed: 100 * time.Millisecond,
		},
		{
			wantErr: context.DeadlineExceeded,
			newCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 25*time.Millisecond)
			},
			name:       "deadline passes while sleeping",
			d:          100 * time.Millisecond,
			minElapsed: 25 * time.Millisecond,
			maxElapsed: 100 * time.Millisecond,
		},
		{
			wantErr: context.Canceled,
			newCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			name:       "already canceled",
			d:          50 * time.Millisecond,
			maxElapsed: 10 * time.Millisecond,
		},
		{
			wantErr: context.DeadlineExceeded,
			newCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), -time.Millisecond)
			},
			name:       "already expired",
			d:          50 * time.Millisecond,
			maxElapsed: 10 * time.Millisecond,
		},
	} {
		start := time.Now()
		ctx, cancel := tc.newCtx()
		err := xtime.Sleep(ctx, tc.d)
		elapsed := time.Since(start)
		cancel()
		c.Equal(tc.wantErr, err, tc.name)
		c.True(elapsed >= tc.minElapsed && elapsed < tc.maxElapsed, "%s: took %v, expected [%v, %v)", tc.name,
			elapsed, tc.minElapsed, tc.maxElapsed)
	}
}
