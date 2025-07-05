// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

import "runtime"

// TaskQueueConfig provides configuration for a TaskQueue.
type TaskQueueConfig struct {
	// RecoveryHandler is the recovery handler to use for tasks that panic. If the handler is nil, the panic will be
	// logged as an error.
	RecoveryHandler func(error)
	// Depth controls the maximum depth of the queue. Calls to Submit() will block when this number of tasks are already
	// pending execution. Zero or less means to use an unbounded queue.
	Depth int
	// Workers controls the number of workers that will simultaneously process tasks. If set to 1, tasks submitted to
	// the queue will be executed serially. If set to less than 1, the number of logical CPUs + 1 will be used instead.
	Workers int
}

// TaskQueue holds the queue information.
type TaskQueue struct {
	in              chan func()
	done            chan bool
	recoveryHandler func(error)
	depth           int
	workers         int
}

// NewTaskQueue creates an asynchronous queue which executes the tasks submitted to it.
func NewTaskQueue(config *TaskQueueConfig) *TaskQueue {
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

// Submit a task to be run.
func (q *TaskQueue) Submit(task func()) {
	q.in <- task
}

// Shutdown the queue. Does not return until all pending tasks have completed.
func (q *TaskQueue) Shutdown() {
	close(q.in)
	<-q.done
}

func (q *TaskQueue) process() {
	var received, processed uint64

	// Setup backlog
	var backlog []func()
	if q.depth > 1 {
		backlog = make([]func(), 0, q.depth-1)
	}

	// Setup workers
	ready := make(chan bool, q.workers)
	tasks := make(chan func(), q.workers)
	for range q.workers {
		go q.work(tasks, ready)
	}

	// Main processing loop
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
				tasks <- backlog[0]
				copy(backlog, backlog[1:])
				backlog[len(backlog)-1] = task
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
	q.done <- true
}

func (q *TaskQueue) work(tasks <-chan func(), ready chan<- bool) {
	for task := range tasks {
		SafeCall(task, q.recoveryHandler)
		ready <- true
	}
}
