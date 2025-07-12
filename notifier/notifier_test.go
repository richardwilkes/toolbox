// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package notifier_test

import (
	"sync"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/notifier"
)

// Mock target for testing
type mockTarget struct {
	notifications []notification
	mu            sync.Mutex
}

type notification struct {
	data     any
	producer any
	name     string
}

func (mt *mockTarget) HandleNotification(name string, data, producer any) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.notifications = append(mt.notifications, notification{
		name:     name,
		data:     data,
		producer: producer,
	})
}

func (mt *mockTarget) getNotifications() []notification {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	result := make([]notification, len(mt.notifications))
	copy(result, mt.notifications)
	return result
}

func (mt *mockTarget) reset() {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.notifications = nil
}

// Mock batch target for testing
type mockBatchTarget struct {
	batchStarts []bool
	mockTarget
}

func (mbt *mockBatchTarget) BatchMode(start bool) {
	mbt.mu.Lock()
	defer mbt.mu.Unlock()
	mbt.batchStarts = append(mbt.batchStarts, start)
}

func (mbt *mockBatchTarget) getBatchStarts() []bool {
	mbt.mu.Lock()
	defer mbt.mu.Unlock()
	result := make([]bool, len(mbt.batchStarts))
	copy(result, mbt.batchStarts)
	return result
}

// Test recovery handler
func testRecoveryHandler(_ error) {
	// Do nothing for tests
}

func TestNew(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	c.NotNil(n)
	c.True(n.Enabled())
	c.Equal(0, n.BatchLevel())
}

func TestNewWithNilRecoveryHandler(t *testing.T) {
	c := check.New(t)

	n := notifier.New(nil)
	c.NotNil(n)
	c.True(n.Enabled())
}

func TestRegisterAndNotify(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register target for "test" notifications
	n.Register(target, 0, "test")

	// Send notification
	n.Notify("test", "producer")

	notifications := target.getNotifications()
	c.Equal(1, len(notifications))
	c.Equal("test", notifications[0].name)
	c.Nil(notifications[0].data)
	c.Equal("producer", notifications[0].producer)
}

func TestRegisterWithData(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")
	n.NotifyWithData("test", "data", "producer")

	notifications := target.getNotifications()
	c.Equal(1, len(notifications))
	c.Equal("test", notifications[0].name)
	c.Equal("data", notifications[0].data)
	c.Equal("producer", notifications[0].producer)
}

func TestRegisterMultipleNames(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register for multiple names
	n.Register(target, 0, "foo", "bar")

	n.Notify("foo", "producer1")
	n.Notify("bar", "producer2")

	notifications := target.getNotifications()
	c.Equal(2, len(notifications))
	c.Equal("foo", notifications[0].name)
	c.Equal("bar", notifications[1].name)
}

func TestHierarchicalNames(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register for parent name
	n.Register(target, 0, "foo")

	// Notify with child names
	n.Notify("foo.bar", "producer1")
	n.Notify("foo.bar.baz", "producer2")
	n.Notify("foo.other", "producer3")

	notifications := target.getNotifications()
	c.Equal(3, len(notifications))
	c.Equal("foo.bar", notifications[0].name)
	c.Equal("foo.bar.baz", notifications[1].name)
	c.Equal("foo.other", notifications[2].name)
}

func TestPriority(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target1 := &mockTarget{}
	target2 := &mockTarget{}
	target3 := &mockTarget{}

	// Register with different priorities
	n.Register(target1, 1, "test")  // medium priority
	n.Register(target2, 10, "test") // high priority
	n.Register(target3, 0, "test")  // low priority

	n.Notify("test", "producer")

	// All should receive notification
	c.Equal(1, len(target1.getNotifications()))
	c.Equal(1, len(target2.getNotifications()))
	c.Equal(1, len(target3.getNotifications()))
}

func TestUnregister(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")
	n.Notify("test", "producer")
	c.Equal(1, len(target.getNotifications()))

	// Unregister and notify again
	n.Unregister(target)
	target.reset()
	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))
}

func TestUnregisterNonExistent(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Should not panic when unregistering non-existent target
	n.Unregister(target)
	c.True(true) // Just verify we got here without panicking
}

func TestBatchTarget(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockBatchTarget{}

	n.Register(target, 0, "test")

	// Start batch
	n.StartBatch()
	c.Equal(1, n.BatchLevel())

	batchStarts := target.getBatchStarts()
	c.Equal(1, len(batchStarts))
	c.True(batchStarts[0])

	// Send notifications during batch
	n.Notify("test", "producer1")
	n.Notify("test", "producer2")

	notifications := target.getNotifications()
	c.Equal(2, len(notifications))

	// End batch
	n.EndBatch()
	c.Equal(0, n.BatchLevel())

	batchStarts = target.getBatchStarts()
	c.Equal(2, len(batchStarts))
	c.True(batchStarts[0])
	c.False(batchStarts[1])
}

func TestNestedBatch(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockBatchTarget{}

	n.Register(target, 0, "test")

	// Start nested batches
	n.StartBatch()
	c.Equal(1, n.BatchLevel())
	n.StartBatch()
	c.Equal(2, n.BatchLevel())

	// Should only get one batch start notification
	batchStarts := target.getBatchStarts()
	c.Equal(1, len(batchStarts))
	c.True(batchStarts[0])

	// End first batch - should still be in batch mode
	n.EndBatch()
	c.Equal(1, n.BatchLevel())
	c.Equal(1, len(target.getBatchStarts())) // No new batch notifications

	// End second batch - should exit batch mode
	n.EndBatch()
	c.Equal(0, n.BatchLevel())

	batchStarts = target.getBatchStarts()
	c.Equal(2, len(batchStarts))
	c.False(batchStarts[1])
}

func TestEnabledDisabled(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	// Disable notifier
	n.SetEnabled(false)
	c.False(n.Enabled())

	// Should not receive notifications when disabled
	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))

	// Re-enable
	n.SetEnabled(true)
	c.True(n.Enabled())

	n.Notify("test", "producer")
	c.Equal(1, len(target.getNotifications()))
}

func TestBatchWhenDisabled(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockBatchTarget{}

	n.Register(target, 0, "test")
	n.SetEnabled(false)

	// Batch operations should not work when disabled
	n.StartBatch()
	c.Equal(0, n.BatchLevel())
	c.Equal(0, len(target.getBatchStarts()))

	n.EndBatch()
	c.Equal(0, n.BatchLevel())
	c.Equal(0, len(target.getBatchStarts()))
}

func TestRegisterFromNotifier(t *testing.T) {
	c := check.New(t)

	n1 := notifier.New(testRecoveryHandler)
	n2 := notifier.New(testRecoveryHandler)
	target1 := &mockTarget{}
	target2 := &mockBatchTarget{}

	// Register targets in first notifier
	n1.Register(target1, 5, "foo", "bar")
	n1.Register(target2, 10, "baz")

	// Copy registrations to second notifier
	n2.RegisterFromNotifier(n1)

	// Test that notifications work in second notifier
	n2.Notify("foo", "producer1")
	n2.Notify("baz", "producer2")

	c.Equal(1, len(target1.getNotifications()))
	c.Equal(1, len(target2.getNotifications()))
}

func TestRegisterFromNotifierMerge(t *testing.T) {
	c := check.New(t)

	n1 := notifier.New(testRecoveryHandler)
	n2 := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register target in both notifiers with different names
	n1.Register(target, 5, "foo")
	n2.Register(target, 10, "bar")

	// Copy registrations from n1 to n2 (should merge)
	n2.RegisterFromNotifier(n1)

	// Target should now receive notifications for both "foo" and "bar"
	n2.Notify("foo", "producer1")
	n2.Notify("bar", "producer2")

	notifications := target.getNotifications()
	c.Equal(2, len(notifications))

	// Should have received both notifications
	names := make(map[string]bool)
	for _, notif := range notifications {
		names[notif.name] = true
	}
	c.True(names["foo"])
	c.True(names["bar"])
}

func TestRegisterFromSameNotifier(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	// Should be safe to register from self
	n.RegisterFromNotifier(n)

	n.Notify("test", "producer")
	c.Equal(1, len(target.getNotifications()))
}

func TestReset(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}
	batchTarget := &mockBatchTarget{}

	n.Register(target, 0, "test")
	n.Register(batchTarget, 0, "test")
	n.StartBatch()

	// Reset should clear everything
	n.Reset()

	c.Equal(0, n.BatchLevel())

	// Notifications should not be received after reset
	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))
	c.Equal(0, len(batchTarget.getNotifications()))
}

func TestEmptyNames(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register with empty names
	n.Register(target, 0, "", ".", "...")

	// Should not receive any notifications
	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))
}

func TestNameNormalization(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "foo.bar")

	// These should all be normalized to "foo.bar"
	n.Notify("foo..bar", "producer1")
	n.Notify(".foo.bar.", "producer2")
	n.Notify("foo...bar", "producer3")

	notifications := target.getNotifications()
	c.Equal(3, len(notifications))
	for _, notif := range notifications {
		c.Equal("foo.bar", notif.name)
	}
}

func TestNotifyEmptyName(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	// Empty names should not trigger notifications
	n.Notify("", "producer")
	n.Notify(".", "producer")
	n.Notify("...", "producer")

	c.Equal(0, len(target.getNotifications()))
}

func TestConcurrentAccess(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	var wg sync.WaitGroup
	numGoroutines := 10
	numNotifications := 100

	// Concurrent notifications
	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range numNotifications {
				n.Notify("test", j)
			}
		}()
	}

	// Concurrent registrations/unregistrations
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 10 {
			newTarget := &mockTarget{}
			n.Register(newTarget, 0, "test")
			time.Sleep(time.Millisecond)
			n.Unregister(newTarget)
		}
	}()

	wg.Wait()

	// Should have received all notifications
	notifications := target.getNotifications()
	c.Equal(numGoroutines*numNotifications, len(notifications))
}

func TestPanicRecovery(t *testing.T) {
	c := check.New(t)

	recoveredErrors := make([]error, 0)
	recoveryHandler := func(err error) {
		recoveredErrors = append(recoveredErrors, err)
	}

	n := notifier.New(recoveryHandler)

	// Create a target that panics
	panicTarget := &panicTarget{}
	n.Register(panicTarget, 0, "test")

	// Should not panic when target panics
	n.Notify("test", "producer")

	// Recovery handler should have been called
	c.Equal(1, len(recoveredErrors))
}

func TestBatchTargetPanicRecovery(t *testing.T) {
	c := check.New(t)

	recoveredErrors := make([]error, 0)
	recoveryHandler := func(err error) {
		recoveredErrors = append(recoveredErrors, err)
	}

	n := notifier.New(recoveryHandler)

	// Create a batch target that panics
	panicBatchTarget := &panicBatchTarget{}
	n.Register(panicBatchTarget, 0, "test")

	// Should not panic when batch target panics
	n.StartBatch()
	n.EndBatch()

	// Recovery handler should have been called twice (start and end)
	c.Equal(2, len(recoveredErrors))
}

func TestHierarchicalNotMatching(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	// Register for "foo.bar"
	n.Register(target, 0, "foo.bar")

	// These should NOT match
	n.Notify("foo", "producer1")      // Parent doesn't match child registration
	n.Notify("foo.barn", "producer2") // Similar but different
	n.Notify("foobar", "producer3")   // No dot separator

	c.Equal(0, len(target.getNotifications()))
}

func TestMultipleTargetsPriority(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)

	// Create targets that record order of execution
	var executionOrder []string
	var mu sync.Mutex

	createTarget := func(name string) notifier.Target {
		return &orderTarget{
			name: name,
			onNotify: func() {
				mu.Lock()
				defer mu.Unlock()
				executionOrder = append(executionOrder, name)
			},
		}
	}

	target1 := createTarget("low")
	target2 := createTarget("high")
	target3 := createTarget("medium")

	// Register with different priorities (higher values get notified first)
	n.Register(target1, 1, "test")  // low priority
	n.Register(target2, 10, "test") // high priority
	n.Register(target3, 5, "test")  // medium priority

	n.Notify("test", "producer")

	// Should be executed in priority order: high, medium, low
	c.Equal(3, len(executionOrder))
	c.Equal("high", executionOrder[0])
	c.Equal("medium", executionOrder[1])
	c.Equal("low", executionOrder[2])
}

// Target that panics on notification
type panicTarget struct{}

func (pt *panicTarget) HandleNotification(_ string, _, _ any) {
	panic("test panic")
}

// Batch target that panics on batch mode
type panicBatchTarget struct {
	panicTarget
}

func (pbt *panicBatchTarget) BatchMode(_ bool) {
	panic("batch panic")
}

// Target that records execution order
type orderTarget struct {
	onNotify func()
	name     string
}

func (ot *orderTarget) HandleNotification(_ string, _, _ any) {
	ot.onNotify()
}
