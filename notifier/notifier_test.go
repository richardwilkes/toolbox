// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

func testRecoveryHandler(_ error) {
}

func TestRegisterAndNotify(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	n.Notify("test", "producer")

	notifications := target.getNotifications()
	c.Equal(1, len(notifications))
	c.Equal("test", notifications[0].name)
	c.Nil(notifications[0].data)
	c.Equal("producer", notifications[0].producer)

	target.reset()
	n.NotifyWithData("test", "data", "producer")
	notifications = target.getNotifications()
	c.Equal(1, len(notifications))
	c.Equal("test", notifications[0].name)
	c.Equal("data", notifications[0].data)
	c.Equal("producer", notifications[0].producer)
}

func TestNameMatching(t *testing.T) {
	for _, tc := range []struct {
		name       string
		registered []string
		notified   []string
		delivered  []string
	}{
		{
			name:       "sub-names match a registered parent",
			registered: []string{"foo"},
			notified:   []string{"foo.bar", "foo.bar.baz", "foo.other"},
			delivered:  []string{"foo.bar", "foo.bar.baz", "foo.other"},
		},
		{
			name:       "each registered name matches",
			registered: []string{"foo", "bar"},
			notified:   []string{"foo", "bar"},
			delivered:  []string{"foo", "bar"},
		},
		{
			name:       "notified names are normalized",
			registered: []string{"foo.bar"},
			notified:   []string{"foo..bar", ".foo.bar.", "foo...bar"},
			delivered:  []string{"foo.bar", "foo.bar", "foo.bar"},
		},
		{
			name:       "parent and look-alike names do not match",
			registered: []string{"foo.bar"},
			notified:   []string{"foo", "foo.barn", "foobar"},
			delivered:  []string{},
		},
		{
			name:       "empty registered names are not wildcards",
			registered: []string{"", ".", "..."},
			notified:   []string{"test"},
			delivered:  []string{},
		},
		{
			name:       "empty notified names are dropped",
			registered: []string{"test"},
			notified:   []string{"", ".", "..."},
			delivered:  []string{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n := notifier.New(testRecoveryHandler)
			target := &mockTarget{}
			n.Register(target, 0, tc.registered...)
			for _, name := range tc.notified {
				n.Notify(name, "producer")
			}
			notifications := target.getNotifications()
			delivered := make([]string, 0, len(notifications))
			for _, one := range notifications {
				delivered = append(delivered, one.name)
			}
			check.New(t).Equal(tc.delivered, delivered)
		})
	}
}

func TestUnregister(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")
	n.Notify("test", "producer")
	c.Equal(1, len(target.getNotifications()))

	n.Unregister(target)
	target.reset()
	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))

	c.NotPanics(func() { n.Unregister(target) })
}

func TestBatchTarget(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockBatchTarget{}

	n.Register(target, 0, "test")

	n.StartBatch()
	c.Equal(1, n.BatchLevel())

	batchStarts := target.getBatchStarts()
	c.Equal(1, len(batchStarts))
	c.True(batchStarts[0])

	n.Notify("test", "producer1")
	n.Notify("test", "producer2")

	notifications := target.getNotifications()
	c.Equal(2, len(notifications))

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

	n.StartBatch()
	c.Equal(1, n.BatchLevel())
	n.StartBatch()
	c.Equal(2, n.BatchLevel())

	batchStarts := target.getBatchStarts()
	c.Equal(1, len(batchStarts))
	c.True(batchStarts[0])

	n.EndBatch()
	c.Equal(1, n.BatchLevel())
	c.Equal(1, len(target.getBatchStarts())) // No new batch notifications

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

	n.SetEnabled(false)
	c.False(n.Enabled())

	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))

	n.SetEnabled(true)
	c.True(n.Enabled())

	n.Notify("test", "producer")
	c.Equal(1, len(target.getNotifications()))
}

func TestBatchWhenDisabled(t *testing.T) {
	for _, tc := range []struct {
		run  func(c check.Checker, n *notifier.Notifier, target *mockBatchTarget)
		name string
	}{
		{
			name: "disabled before the batch",
			run: func(c check.Checker, n *notifier.Notifier, target *mockBatchTarget) {
				n.SetEnabled(false)

				// The batch lifecycle is independent of the enabled state: the level is still tracked and the
				// BatchMode brackets are still delivered, so StartBatch/EndBatch stay balanced.
				n.StartBatch()
				c.Equal(1, n.BatchLevel())
				c.Equal([]bool{true}, target.getBatchStarts())

				// Data notifications, however, remain suppressed while disabled.
				n.Notify("test", "producer")
				c.Equal(0, len(target.getNotifications()))

				n.EndBatch()
				c.Equal(0, n.BatchLevel())
				c.Equal([]bool{true, false}, target.getBatchStarts())
			},
		},
		{
			name: "disabled mid-batch",
			run: func(c check.Checker, n *notifier.Notifier, target *mockBatchTarget) {
				n.StartBatch()
				c.Equal(1, n.BatchLevel())
				c.Equal([]bool{true}, target.getBatchStarts())

				n.SetEnabled(false)
				n.EndBatch()
				c.Equal(0, n.BatchLevel(), "batch level must drop to 0 even though the notifier was disabled mid-batch")
				c.Equal([]bool{true, false}, target.getBatchStarts(), "matching BatchMode(false) must still be delivered")

				// Re-enabling must yield a fully balanced batch rather than a level stuck above 0.
				n.SetEnabled(true)
				n.StartBatch()
				c.Equal(1, n.BatchLevel())
				n.EndBatch()
				c.Equal(0, n.BatchLevel())
				c.Equal([]bool{true, false, true, false}, target.getBatchStarts())
			},
		},
		{
			name: "disabled mid nested batch",
			run: func(c check.Checker, n *notifier.Notifier, target *mockBatchTarget) {
				n.StartBatch()
				n.StartBatch()
				c.Equal(2, n.BatchLevel())
				c.Equal([]bool{true}, target.getBatchStarts())

				n.SetEnabled(false)
				n.EndBatch()
				c.Equal(1, n.BatchLevel())
				c.Equal([]bool{true}, target.getBatchStarts(), "no BatchMode delivery until the outermost batch ends")
				n.EndBatch()
				c.Equal(0, n.BatchLevel())
				c.Equal([]bool{true, false}, target.getBatchStarts())
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n := notifier.New(testRecoveryHandler)
			target := &mockBatchTarget{}
			n.Register(target, 0, "test")
			tc.run(check.New(t), n, target)
		})
	}
}

func TestRegisterFromNotifier(t *testing.T) {
	for _, tc := range []struct {
		run  func(c check.Checker, src, dst *notifier.Notifier)
		name string
	}{
		{
			name: "targets new to the destination",
			run: func(c check.Checker, src, dst *notifier.Notifier) {
				target1 := &mockTarget{}
				target2 := &mockBatchTarget{}

				src.Register(target1, 5, "foo", "bar")
				src.Register(target2, 10, "baz")

				dst.RegisterFromNotifier(src)

				dst.Notify("foo", "producer1")
				dst.Notify("baz", "producer2")

				c.Equal(1, len(target1.getNotifications()))
				c.Equal(1, len(target2.getNotifications()))
			},
		},
		{
			name: "target already in the destination keeps its names and gains new ones",
			run: func(c check.Checker, src, dst *notifier.Notifier) {
				target := &mockTarget{}

				src.Register(target, 5, "foo")
				dst.Register(target, 10, "bar")

				dst.RegisterFromNotifier(src)

				dst.Notify("foo", "producer1")
				dst.Notify("bar", "producer2")

				notifications := target.getNotifications()
				c.Equal(2, len(notifications))

				names := make(map[string]bool)
				for _, notif := range notifications {
					names[notif.name] = true
				}
				c.True(names["foo"])
				c.True(names["bar"])
			},
		},
		{
			name: "from itself is a no-op",
			run: func(c check.Checker, n, _ *notifier.Notifier) {
				target := &mockTarget{}

				n.Register(target, 0, "test")

				n.RegisterFromNotifier(n)

				n.Notify("test", "producer")
				c.Equal(1, len(target.getNotifications()))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(check.New(t), notifier.New(testRecoveryHandler), notifier.New(testRecoveryHandler))
		})
	}
}

func TestReset(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}
	batchTarget := &mockBatchTarget{}

	n.Register(target, 0, "test")
	n.Register(batchTarget, 0, "test")
	n.StartBatch()

	n.Reset()

	c.Equal(0, n.BatchLevel())

	n.Notify("test", "producer")
	c.Equal(0, len(target.getNotifications()))
	c.Equal(0, len(batchTarget.getNotifications()))
}

func TestConcurrentAccess(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)
	target := &mockTarget{}

	n.Register(target, 0, "test")

	var wg sync.WaitGroup
	numGoroutines := 10
	numNotifications := 100

	for range numGoroutines {
		wg.Go(func() {
			for j := range numNotifications {
				n.Notify("test", j)
			}
		})
	}

	wg.Go(func() {
		for range 10 {
			newTarget := &mockTarget{}
			n.Register(newTarget, 0, "test")
			time.Sleep(time.Millisecond)
			n.Unregister(newTarget)
		}
	})

	wg.Wait()

	notifications := target.getNotifications()
	c.Equal(numGoroutines*numNotifications, len(notifications))
}

func TestPanicRecovery(t *testing.T) {
	for _, tc := range []struct {
		target    notifier.Target
		trigger   func(n *notifier.Notifier)
		name      string
		recovered int
	}{
		{
			name:      "HandleNotification",
			target:    &panicTarget{},
			trigger:   func(n *notifier.Notifier) { n.Notify("test", "producer") },
			recovered: 1,
		},
		{
			name:   "BatchMode on start and end",
			target: &panicBatchTarget{},
			trigger: func(n *notifier.Notifier) {
				n.StartBatch()
				n.EndBatch()
			},
			recovered: 2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := check.New(t)

			recoveredErrors := make([]error, 0)
			recoveryHandler := func(err error) {
				recoveredErrors = append(recoveredErrors, err)
			}

			n := notifier.New(recoveryHandler)
			n.Register(tc.target, 0, "test")

			tc.trigger(n)

			c.Equal(tc.recovered, len(recoveredErrors))
		})
	}
}

func TestMultipleTargetsPriority(t *testing.T) {
	c := check.New(t)

	n := notifier.New(testRecoveryHandler)

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

	n.Register(target1, 1, "test")  // low priority
	n.Register(target2, 10, "test") // high priority
	n.Register(target3, 5, "test")  // medium priority

	n.Notify("test", "producer")

	c.Equal(3, len(executionOrder))
	c.Equal("high", executionOrder[0])
	c.Equal("medium", executionOrder[1])
	c.Equal("low", executionOrder[2])
}

type panicTarget struct{}

func (pt *panicTarget) HandleNotification(_ string, _, _ any) {
	panic("test panic")
}

type panicBatchTarget struct {
	panicTarget
}

func (pbt *panicBatchTarget) BatchMode(_ bool) {
	panic("batch panic")
}

type orderTarget struct {
	onNotify func()
	name     string
}

func (ot *orderTarget) HandleNotification(_ string, _, _ any) {
	ot.onNotify()
}
