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
	ctx := context.Background()
	start := time.Now()
	duration := 50 * time.Millisecond

	err := xtime.Sleep(ctx, duration)

	elapsed := time.Since(start)
	c.NoError(err)
	c.True(elapsed >= duration, "Expected elapsed time to be at least %v, got %v", duration, elapsed)
	c.True(elapsed < duration+20*time.Millisecond, "Expected elapsed time to be less than %v, got %v", duration+20*time.Millisecond, elapsed)
}

func TestSleepInterruptedByContextCancellation(t *testing.T) {
	c := check.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	start := time.Now()
	duration := 100 * time.Millisecond

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := xtime.Sleep(ctx, duration)

	elapsed := time.Since(start)
	c.Equal(context.Canceled, err)
	c.True(elapsed < duration, "Expected elapsed time to be less than %v, got %v", duration, elapsed)
	c.True(elapsed >= 8*time.Millisecond, "Expected elapsed time to be at least 8ms, got %v", elapsed)
}

func TestSleepInterruptedByContextTimeout(t *testing.T) {
	c := check.New(t)
	timeout := 25 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	duration := 100 * time.Millisecond

	err := xtime.Sleep(ctx, duration)

	elapsed := time.Since(start)
	c.Equal(context.DeadlineExceeded, err)
	c.True(elapsed >= timeout, "Expected elapsed time to be at least %v, got %v", timeout, elapsed)
	c.True(elapsed < duration, "Expected elapsed time to be less than %v, got %v", duration, elapsed)
}

func TestSleepWithZeroDuration(t *testing.T) {
	c := check.New(t)
	ctx := context.Background()
	start := time.Now()

	err := xtime.Sleep(ctx, 0)

	elapsed := time.Since(start)
	c.NoError(err)
	c.True(elapsed < 10*time.Millisecond, "Expected elapsed time to be less than 10ms, got %v", elapsed)
}

func TestSleepWithNegativeDuration(t *testing.T) {
	c := check.New(t)
	ctx := context.Background()
	start := time.Now()

	err := xtime.Sleep(ctx, -10*time.Millisecond)

	elapsed := time.Since(start)
	c.NoError(err)
	c.True(elapsed < 10*time.Millisecond, "Expected elapsed time to be less than 10ms, got %v", elapsed)
}

func TestSleepWithAlreadyCancelledContext(t *testing.T) {
	c := check.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	start := time.Now()
	duration := 50 * time.Millisecond

	err := xtime.Sleep(ctx, duration)

	elapsed := time.Since(start)
	c.Equal(context.Canceled, err)
	c.True(elapsed < 10*time.Millisecond, "Expected elapsed time to be less than 10ms, got %v", elapsed)
}

func TestSleepWithAlreadyExpiredContext(t *testing.T) {
	c := check.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), -1*time.Millisecond)
	defer cancel()

	start := time.Now()
	duration := 50 * time.Millisecond

	err := xtime.Sleep(ctx, duration)

	elapsed := time.Since(start)
	c.Equal(context.DeadlineExceeded, err)
	c.True(elapsed < 10*time.Millisecond, "Expected elapsed time to be less than 10ms, got %v", elapsed)
}

func TestSleepConcurrentOperations(t *testing.T) {
	c := check.New(t)
	ctx := context.Background()
	duration := 30 * time.Millisecond

	start := time.Now()
	errChan := make(chan error, 3)

	for range 3 {
		go func() {
			errChan <- xtime.Sleep(ctx, duration)
		}()
	}

	for range 3 {
		err := <-errChan
		c.NoError(err)
	}

	elapsed := time.Since(start)
	c.True(elapsed >= duration, "Expected elapsed time to be at least %v, got %v", duration, elapsed)
	c.True(elapsed < duration+20*time.Millisecond, "Expected elapsed time to be less than %v, got %v", duration+20*time.Millisecond, elapsed)
}
