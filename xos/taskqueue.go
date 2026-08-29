// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import (
	"runtime"
	"sync"
)

// TaskQueueConfig provides configuration for a TaskQueue.
type TaskQueueConfig struct {
	// RecoveryHandler receives, as an error, any panic raised by a task. If nil, the panic is logged.
	RecoveryHandler func(error)
	// Depth limits the number of tasks queued awaiting a worker; beyond it, Submit blocks once a small internal buffer
	// fills. Zero or less means unbounded.
	Depth int
	// Workers is the number of workers processing tasks concurrently; 1 executes tasks serially. Less than 1 means to
	// use the number of logical CPUs + 1.
	Workers int
}

// TaskQueue executes submitted tasks asynchronously on a pool of workers.
type TaskQueue struct {
	in              chan func()
	done            chan bool
	recoveryHandler func(error)
	lock            sync.RWMutex
	depth           int
	workers         int
	closed          bool
}

// NewTaskQueue creates a queue that asynchronously executes the tasks submitted to it. A nil config uses the defaults.
func NewTaskQueue(config *TaskQueueConfig) *TaskQueue {
	if config == nil {
		config = &TaskQueueConfig{}
	}
	numCPU := runtime.NumCPU()
	q := &TaskQueue{
		in:              make(chan func(), numCPU*2),
		done:            make(chan bool),
		recoveryHandler: config.RecoveryHandler,
		depth:           config.Depth,
		workers:         config.Workers,
	}
	if q.workers < 1 {
		q.workers = 1 + numCPU
	}
	go q.process()
	return q
}

// Submit queues a task to be run. Returns true if it was accepted, or false if the queue has already been shut down, in
// which case the task is dropped. Safe to call concurrently, including concurrently with or after Shutdown.
func (q *TaskQueue) Submit(task func()) bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if q.closed {
		return false
	}
	q.in <- task
	return true
}

// Shutdown closes the queue and returns once all pending tasks have completed, after which Submit rejects tasks. Safe
// to call more than once and concurrently; every call blocks until completion.
func (q *TaskQueue) Shutdown() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		close(q.in)
	}
	q.lock.Unlock()
	<-q.done
}

func (q *TaskQueue) process() {
	var received, processed uint64

	var backlog []func()
	if q.depth > 1 {
		backlog = make([]func(), 0, q.depth-1)
	}

	ready := make(chan bool, q.workers)
	tasks := make(chan func(), q.workers)
	for range q.workers {
		go q.work(tasks, ready)
	}

outer:
	for {
	inner:
		select {
		case task := <-q.in:
			if task == nil {
				break outer
			}
			received++
			if len(backlog) == 0 {
				select {
				case tasks <- task:
					break inner
				default:
				}
			}
			if q.depth <= 0 || 1+len(backlog) < q.depth {
				backlog = append(backlog, task)
			} else {
				<-ready
				processed++
				if len(backlog) == 0 {
					tasks <- task
				} else {
					tasks <- backlog[0]
					copy(backlog, backlog[1:])
					backlog[len(backlog)-1] = task
				}
			}
		case <-ready:
			processed++
			if len(backlog) > 0 {
				tasks <- backlog[0]
				copy(backlog, backlog[1:])
				backlog[len(backlog)-1] = nil
				backlog = backlog[:len(backlog)-1]
			}
		}
	}

	// Finish any remaining tasks
	for _, task := range backlog {
	drain:
		for {
			select {
			case tasks <- task:
				break drain
			case <-ready:
				processed++
			}
		}
	}
	for received != processed {
		<-ready
		processed++
	}
	close(tasks)
	close(q.done)
}

func (q *TaskQueue) work(tasks <-chan func(), ready chan<- bool) {
	for task := range tasks {
		SafeCall(task, q.recoveryHandler)
		ready <- true
	}
}
