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

	c.Equal(1000, parent.Cap(false))
	c.Equal(800, child1.Cap(false))
	c.Equal(1200, child2.Cap(false))
	c.Equal(600, grandchild.Cap(false))

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

	// 80 exceeds the parent-capped effective capacity (50), so it can never be satisfied and must be rejected promptly
	// rather than queued forever.
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

	// Consume all of the parent's capacity so the child's request must wait for the next tick.
	err := <-parent.Use(100)
	c.NoError(err)

	// Queue a child request that is valid right now (60 <= the effective cap of 100).
	done := child.Use(60)

	// Lower the parent's cap below the queued amount; the per-tick re-check must now reject the request rather than
	// leave it queued forever.
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

	err := <-rl.Use(50)
	c.NoError(err)

	err = <-rl.Use(30)
	c.NoError(err)

	err = <-rl.Use(20)
	c.NoError(err)

	rl.Close()
}

func TestUseWaiting(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, 50*time.Millisecond)

	err := <-rl.Use(100)
	c.NoError(err)

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

func TestUseFIFOOrderingNoFastPathJump(t *testing.T) {
	c := check.New(t)
	rl := rate.New(10, 25*time.Millisecond)
	defer rl.Close()

	// Use part of the period so a full-capacity request must wait for the next tick while spare capacity remains.
	c.NoError(<-rl.Use(5))

	big := rl.Use(10)

	// Even though 5 units are free, FIFO admission must keep these behind the queued large request.
	small1 := rl.Use(1)
	small2 := rl.Use(1)

	// This runs synchronously right after submission, well within one period, so no tick can have intervened.
	assertPending := func(ch <-chan error, name string) {
		select {
		case <-ch:
			t.Fatalf("%s slipped through the fast path ahead of an earlier queued request", name)
		default:
		}
	}
	assertPending(small1, "small1")
	assertPending(small2, "small2")

	// The large request must eventually be served rather than starved by later requests.
	select {
	case err := <-big:
		c.NoError(err)
	case <-time.After(2 * time.Second):
		t.Fatal("large request was starved and never completed")
	}

	c.NoError(<-small1)
	c.NoError(<-small2)
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

	time.Sleep(60 * time.Millisecond)

	// After the reset, LastUsed reflects the previous period's usage.
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

	err := <-child1.Use(40)
	c.NoError(err)

	err = <-child2.Use(30)
	c.NoError(err)

	// The parent has used 70, so child1 can use 20 more to reach its own limit of 60.
	err = <-child1.Use(20)
	c.NoError(err)

	// child1 is now at its capacity, so this cannot succeed immediately.
	done := child1.Use(10)
	select {
	case err = <-done:
		c.HasError(err)
		c.Contains(err.Error(), "capacity")
	case <-time.After(10 * time.Millisecond):
		// Queued until the next period, which is also acceptable.
	}

	// The parent has only 10 left, so child2 can use exactly 10 more.
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

	// New on a closed limiter returns a non-nil, closed limiter rather than nil.
	child := rl.New(50)
	c.NotNil(child)
	c.True(child.Closed())

	// A chained call must yield a closed error rather than panic on a nil interface.
	err := <-rl.New(50).Use(10)
	c.HasError(err)
}

func TestConcurrentUse(t *testing.T) {
	c := check.New(t)
	rl := rate.New(1000, 100*time.Millisecond)

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for range 100 {
		wg.Go(func() {
			err := <-rl.Use(10)
			errors <- err
		})
	}

	wg.Wait()
	close(errors)

	successCount := 0
	for err := range errors {
		if err == nil {
			successCount++
		}
	}

	// 100 * 10 = 1000 fits the capacity, so all requests succeed.
	c.Equal(100, successCount)

	rl.Close()
}

func TestWaitingRequestsClearedOnClose(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, time.Second)

	err := <-rl.Use(100)
	c.NoError(err)

	doneCh := rl.Use(50)

	go func() {
		time.Sleep(10 * time.Millisecond)
		rl.Close()
	}()

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
	rl := rate.New(100, 100*time.Millisecond)

	start := time.Now()

	err := <-rl.Use(100)
	c.NoError(err)

	// The second 100 must wait for the period to reset.
	err = <-rl.Use(100)
	c.NoError(err)

	elapsed := time.Since(start)
	// Both requests together take at least one period.
	c.True(elapsed >= 80*time.Millisecond) // Allow some tolerance for timing

	rl.Close()
}

func TestHierarchicalConstraints(t *testing.T) {
	c := check.New(t)
	parent := rate.New(100, time.Second)
	child := parent.New(80)

	err := <-child.Use(80)
	c.NoError(err)

	// The child is at capacity even though the parent has 20 left, so this cannot succeed immediately.
	done := child.Use(10)

	select {
	case err = <-done:
		if err != nil {
			c.Contains(err.Error(), "capacity")
		}
	case <-time.After(10 * time.Millisecond):
		// Queued until the next period, which is also acceptable.
	}

	parent.Close()
}
