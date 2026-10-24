// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xos"
)

var (
	prev    int
	counter int
)

func TestSerialQueue(t *testing.T) {
	q := xos.NewTaskQueue(&xos.TaskQueueConfig{Depth: 100, Workers: 1})
	prev = -1
	counter = 0
	for i := range 200 {
		submitSerial(q, i)
	}
	q.Shutdown()
	c := check.New(t)
	c.Equal(199, prev)
	c.Equal(200, counter)
}

func submitSerial(q *xos.TaskQueue, i int) {
	q.Submit(func() {
		if i-1 == prev {
			prev = i
			counter++
		}
	})
}

func TestQueueRunsEveryTask(t *testing.T) {
	for _, tc := range []struct {
		name   string
		config xos.TaskQueueConfig
		tasks  int
		delay  time.Duration
	}{
		{name: "depth one", config: xos.TaskQueueConfig{Depth: 1, Workers: 1}, tasks: 2000, delay: time.Microsecond},
		{name: "parallel", config: xos.TaskQueueConfig{Workers: 5}, tasks: 10000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := check.New(t)
			c.NotPanics(func() {
				q := xos.NewTaskQueue(&tc.config)
				var count, total int32
				for i := range tc.tasks {
					q.Submit(func() {
						time.Sleep(tc.delay)
						atomic.AddInt32(&total, int32(i))
						atomic.AddInt32(&count, 1)
					})
				}
				q.Shutdown()
				c.Equal(int32(tc.tasks), atomic.LoadInt32(&count))
				c.Equal(int32(tc.tasks*(tc.tasks-1)/2), atomic.LoadInt32(&total))
			})
		})
	}
}

func TestRecovery(t *testing.T) {
	c := check.New(t)
	c.Panics(boom)
	logged := false
	c.NotPanics(func() {
		q := xos.NewTaskQueue(&xos.TaskQueueConfig{RecoveryHandler: func(_ error) { logged = true }})
		q.Submit(boom)
		q.Shutdown()
	})
	c.True(logged)
}

func TestRecoveryWithBadLogger(t *testing.T) {
	c := check.New(t)
	c.Panics(boom)
	c.NotPanics(func() {
		q := xos.NewTaskQueue(&xos.TaskQueueConfig{RecoveryHandler: func(_ error) { boom() }})
		q.Submit(boom)
		q.Shutdown()
	})
}

func TestSubmitAfterShutdown(t *testing.T) {
	c := check.New(t)
	c.NotPanics(func() {
		var ran int32
		q := xos.NewTaskQueue(&xos.TaskQueueConfig{Workers: 2})
		c.True(q.Submit(func() { atomic.AddInt32(&ran, 1) }))
		q.Shutdown()

		// Submitting after Shutdown must not panic ("send on closed channel"); it is rejected and dropped.
		var late int32
		c.False(q.Submit(func() { atomic.AddInt32(&late, 1) }))
		c.Equal(int32(1), atomic.LoadInt32(&ran))
		c.Equal(int32(0), atomic.LoadInt32(&late))
	})
}

func TestSubmitConcurrentWithShutdown(t *testing.T) {
	c := check.New(t)
	c.NotPanics(func() {
		q := xos.NewTaskQueue(&xos.TaskQueueConfig{Workers: 4})
		var wg sync.WaitGroup
		for range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range 100 {
					// Whether accepted or rejected, this must never panic on a closed channel.
					q.Submit(func() {})
				}
			}()
		}
		q.Shutdown()
		wg.Wait()
		// Any submission that races in after Shutdown is simply rejected.
		c.False(q.Submit(func() {}))
	})
}

func TestShutdownIdempotent(t *testing.T) {
	c := check.New(t)
	c.NotPanics(func() {
		q := xos.NewTaskQueue(&xos.TaskQueueConfig{Workers: 2})
		c.True(q.Submit(func() {}))
		var wg sync.WaitGroup
		for range 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				q.Shutdown() // Every concurrent caller must return without panicking on a double close.
			}()
		}
		wg.Wait()
		c.False(q.Submit(func() {}))
	})
}

func boom() {
	var bad *int
	*bad = 1 //nolint:govet // Yes, this is an intentional store to a nil pointer
}
