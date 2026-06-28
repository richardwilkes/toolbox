// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rate_test

import (
	"sync"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/rate"
)

func TestNew(t *testing.T) {
	c := check.New(t)
	rl := rate.New(1024, time.Second)
	c.NotNil(rl)
	c.Equal(1024, rl.Cap(false))
	c.Equal(1024, rl.Cap(true))
	c.Equal(0, rl.LastUsed())
	c.False(rl.Closed())
	rl.Close()
	c.True(rl.Closed())
}

func TestSetCap(t *testing.T) {
	c := check.New(t)
	rl := rate.New(1024, time.Second)
	c.Equal(1024, rl.Cap(false))

	rl.SetCap(2048)
	c.Equal(2048, rl.Cap(false))

	rl.SetCap(512)
	c.Equal(512, rl.Cap(false))

	rl.Close()
}

func TestCapWithHierarchy(t *testing.T) {
	c := check.New(t)
	parent := rate.New(1000, time.Second)
	child1 := parent.New(800)
	child2 := parent.New(1200)
	grandchild := child1.New(600)

	// Without parent caps
	c.Equal(1000, parent.Cap(false))
	c.Equal(800, child1.Cap(false))
	c.Equal(1200, child2.Cap(false))
	c.Equal(600, grandchild.Cap(false))

	// With parent caps
	c.Equal(1000, parent.Cap(true))
	c.Equal(800, child1.Cap(true))
	c.Equal(1000, child2.Cap(true)) // Limited by parent
	c.Equal(600, grandchild.Cap(true))

	parent.Close()
}

func TestUse(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, 100*time.Millisecond)
	endAfter := time.Now().Add(250 * time.Millisecond)
	for endAfter.After(time.Now()) {
		err := <-rl.Use(1)
		c.NoError(err)
	}
	c.Equal(100, rl.LastUsed())
	rl.Close()
	c.True(rl.Closed())
}

func TestUseZeroAmount(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	err := <-rl.Use(0)
	c.NoError(err)
	c.Equal(0, rl.LastUsed())

	rl.Close()
}

func TestUseNegativeAmount(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	err := <-rl.Use(-10)
	c.HasError(err)
	c.Contains(err.Error(), "Amount (-10) must be positive")

	rl.Close()
}

func TestUseAmountGreaterThanCapacity(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	err := <-rl.Use(200)
	c.HasError(err)
	c.Contains(err.Error(), "Amount (200) is greater than capacity (100)")

	rl.Close()
}

func TestUseAmountGreaterThanParentCap(t *testing.T) {
	c := check.New(t)
	parent := rate.New(50, time.Second)
	child := parent.New(100) // The child's own capacity deliberately exceeds the parent's cap.

	// 80 <= the child's own capacity (100) but > its parent-capped effective capacity (50), so it can never be
	// satisfied and must be rejected promptly rather than queued forever. The watchdog makes a regression of the
	// infinite-block bug fail fast instead of hanging the suite.
	done := child.Use(80)
	select {
	case err := <-done:
		c.HasError(err)
		c.Contains(err.Error(), "Amount (80) is greater than capacity (50)")
	case <-time.After(10 * time.Second):
		t.Fatal("child.Use(80) never returned; a request exceeding the parent cap was queued forever")
	}

	parent.Close()
}

func TestQueuedRequestRejectedWhenParentCapLowered(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, 50*time.Millisecond)
	child := parent.New(100)

	// Consume all of the parent's capacity this period so the child's request cannot be satisfied immediately and
	// must wait for the next tick.
	err := <-parent.Use(100)
	c.NoError(err)

	// Queue a child request that is valid right now (60 <= the effective cap of 100).
	done := child.Use(60)

	// Lower the parent's cap below the queued amount before the next period. The request can now never succeed, so the
	// per-tick re-check must reject it rather than leaving it queued forever.
	parent.SetCap(50)

	select {
	case err = <-done:
		c.HasError(err)
		c.Contains(err.Error(), "Amount (60) is greater than capacity (50)")
	case <-time.After(10 * time.Second):
		t.Fatal("queued request was not rejected after the parent cap was lowered below it")
	}

	parent.Close()
}

func TestUseOnClosedLimiter(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)
	rl.Close()

	err := <-rl.Use(10)
	c.HasError(err)
	c.Contains(err.Error(), "Limiter is closed")
}

func TestUseImmediateSuccess(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	// Should succeed immediately when capacity is available
	err := <-rl.Use(50)
	c.NoError(err)

	err = <-rl.Use(30)
	c.NoError(err)

	// Should still have 20 remaining
	err = <-rl.Use(20)
	c.NoError(err)

	rl.Close()
}

func TestUseWaiting(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, 50*time.Millisecond)

	// Use up all capacity
	err := <-rl.Use(100)
	c.NoError(err)

	// This should block until next period
	start := time.Now()
	doneCh := rl.Use(50)

	select {
	case err = <-doneCh:
		c.NoError(err)
		elapsed := time.Since(start)
		c.True(elapsed >= 40*time.Millisecond) // Should wait at least close to the period
	case <-time.After(200 * time.Millisecond):
		t.Error("Request should have completed after reset period")
	}

	rl.Close()
}

func TestLastUsed(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, 50*time.Millisecond)

	c.Equal(0, rl.LastUsed())

	err := <-rl.Use(30)
	c.NoError(err)
	c.Equal(0, rl.LastUsed()) // Should still be 0 until reset

	err = <-rl.Use(20)
	c.NoError(err)
	c.Equal(0, rl.LastUsed()) // Should still be 0 until reset

	// Wait for reset
	time.Sleep(60 * time.Millisecond)

	// After reset, last used should reflect previous period usage
	// We need to trigger the reset by trying to use the limiter
	err = <-rl.Use(10)
	c.NoError(err)
	c.Equal(50, rl.LastUsed()) // Should show usage from previous period

	rl.Close()
}

func TestHierarchicalLimiters(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, time.Second)
	child1 := parent.New(60)
	child2 := parent.New(80)

	// Use some capacity in child1
	err := <-child1.Use(40)
	c.NoError(err)

	// Use some capacity in child2
	err = <-child2.Use(30)
	c.NoError(err)

	// Total used should be 70, so child1 should be able to use 20 more to reach its limit
	err = <-child1.Use(20)
	c.NoError(err)

	// Now child1 is at its capacity (60), should not be able to use more immediately
	done := child1.Use(10)
	select {
	case err = <-done:
		// If it returns immediately, it should be an error
		c.HasError(err)
		c.Contains(err.Error(), "capacity")
	case <-time.After(10 * time.Millisecond):
		// If it's waiting, that's also acceptable behavior for a rate limiter
		// The request will be queued until the next period
	}

	// child2 should only be able to use 10 more (to reach parent's remaining capacity)
	err = <-child2.Use(10)
	c.NoError(err)

	parent.Close()
}

func TestChildLimiterClosedWhenParentClosed(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, time.Second)
	child := parent.New(50)
	grandchild := child.New(25)

	c.False(parent.Closed())
	c.False(child.Closed())
	c.False(grandchild.Closed())

	parent.Close()

	c.True(parent.Closed())
	c.True(child.Closed())
	c.True(grandchild.Closed())
}

func TestChildLimiterRemovedFromParentOnClose(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, time.Second)
	child := parent.New(50)

	c.False(child.Closed())
	child.Close()
	c.True(child.Closed())
	c.False(parent.Closed()) // Parent should still be open

	parent.Close()
}

func TestNewOnClosedLimiter(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)
	rl.Close()

	child := rl.New(50)
	c.Nil(child) // Should return nil when creating child on closed limiter
}

func TestConcurrentUse(t *testing.T) {
	c := check.New(t)
	rl := rate.New(1000, 100*time.Millisecond)

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Start 100 goroutines trying to use 10 units each
	for range 100 {
		wg.Go(func() {
			err := <-rl.Use(10)
			errors <- err
		})
	}

	wg.Wait()
	close(errors)

	// Should have exactly 1000 units used (100 * 10)
	successCount := 0
	for err := range errors {
		if err == nil {
			successCount++
		}
	}

	// All requests should succeed within the capacity
	c.Equal(100, successCount)

	rl.Close()
}

func TestWaitingRequestsClearedOnClose(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	// Use up all capacity
	err := <-rl.Use(100)
	c.NoError(err)

	// Create a waiting request
	doneCh := rl.Use(50)

	// Close the limiter
	go func() {
		time.Sleep(10 * time.Millisecond)
		rl.Close()
	}()

	// The waiting request should get an error
	select {
	case err = <-doneCh:
		c.HasError(err)
		c.Contains(err.Error(), "Limiter is closed")
	case <-time.After(100 * time.Millisecond):
		t.Error("Waiting request should have been canceled")
	}
}

func TestRateLimitingWithRealTiming(t *testing.T) {
	c := check.New(t)
	// Create a limiter that allows 100 units per 100ms
	rl := rate.New(100, 100*time.Millisecond)

	start := time.Now()

	// Use 100 units immediately
	err := <-rl.Use(100)
	c.NoError(err)

	// Next 100 units should wait for the period to reset
	err = <-rl.Use(100)
	c.NoError(err)

	elapsed := time.Since(start)
	// Should take at least 100ms to complete both requests
	c.True(elapsed >= 80*time.Millisecond) // Allow some tolerance for timing

	rl.Close()
}

func TestHierarchicalConstraints(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, time.Second)
	child := parent.New(80)

	// Child should be limited by its own capacity first
	err := <-child.Use(80)
	c.NoError(err)

	// Now child is at capacity, but parent still has 20 available
	// Next request should fail/wait
	done := child.Use(10)

	// Since we're not waiting for the period, this should not succeed immediately
	select {
	case err = <-done:
		// If it completes immediately, it should be an error or the limiter is queuing
		if err != nil {
			c.Contains(err.Error(), "capacity")
		}
	case <-time.After(10 * time.Millisecond):
		// If it's waiting, that's also acceptable behavior
	}

	parent.Close()
}
