// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package taskqueue_test

import (
	"sync/atomic"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/taskqueue"
)

const (
	parallelWorkSubmissions = 10000
	workTotal               = 49995000
)

var (
	prev    int
	counter int
	total   int32
	count   int32
)

func TestSerialQueue(t *testing.T) {
	q := taskqueue.New(taskqueue.Depth(100), taskqueue.Workers(1))
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

func submitSerial(q *taskqueue.Queue, i int) {
	q.Submit(func() {
		if i-1 == prev {
			prev = i
			counter++
		}
	})
}

func TestParallelQueue(t *testing.T) {
	q := taskqueue.New(taskqueue.Workers(5))
	total = 0
	count = 0
	for i := range parallelWorkSubmissions {
		submitParallel(q, i)
	}
	q.Shutdown()
	c := check.New(t)
	c.Equal(parallelWorkSubmissions, int(count))
	c.Equal(workTotal, int(total))
}

func submitParallel(q *taskqueue.Queue, i int) {
	q.Submit(func() {
		atomic.AddInt32(&total, int32(i))
		atomic.AddInt32(&count, 1)
	})
}

func TestRecovery(t *testing.T) {
	c := check.New(t)
	c.Panics(boom)
	logged := false
	c.NotPanics(func() {
		q := taskqueue.New(taskqueue.RecoveryHandler(func(_ error) { logged = true }))
		q.Submit(boom)
		q.Shutdown()
	})
	c.True(logged)
}

func TestRecoveryWithBadLogger(t *testing.T) {
	c := check.New(t)
	c.Panics(boom)
	c.NotPanics(func() {
		q := taskqueue.New(taskqueue.RecoveryHandler(func(_ error) { boom() }))
		q.Submit(boom)
		q.Shutdown()
	})
}

func boom() {
	var bad *int
	*bad = 1 //nolint:govet // Yes, this is an intentional store to a nil pointer
}
