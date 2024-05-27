/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

// Package taskqueue provides a simple asynchronous task queue.
package taskqueue

import (
	"runtime"

	"github.com/richardwilkes/toolbox/errs"
)

// Logger provides a way to log panics caused by workers in a queue.
type Logger func(v ...any)

// Task defines a unit of work.
type Task func()

// Option defines an option for the queue.
type Option func(*Queue)

// Queue holds the queue information.
type Queue struct {
	in              chan Task
	done            chan bool
	depth           int
	workers         int
	recoveryHandler errs.RecoveryHandler
}

// RecoveryHandler sets the recovery handler to use for tasks that panic. Defaults to none, which silently ignores the
// panic.
func RecoveryHandler(recoveryHandler errs.RecoveryHandler) Option {
	return func(q *Queue) { q.recoveryHandler = recoveryHandler }
}

// Depth sets the depth of the queue. Calls to Submit() will block when this number of tasks are already pending
// execution. Pass in a negative number to use an unbounded queue. Defaults to unbounded.
func Depth(depth int) Option {
	return func(q *Queue) { q.depth = depth }
}

// Workers sets the number of workers that will simultaneously process tasks. If this is set to 1, tasks submitted to
// the queue will be executed serially. Defaults to one plus the number of CPUs.
func Workers(workers int) Option {
	return func(q *Queue) { q.workers = workers }
}

// New creates a queue which executes the tasks submitted to it.
func New(options ...Option) *Queue {
	numCPU := runtime.NumCPU()
	q := &Queue{
		in:    make(chan Task, numCPU*2),
		done:  make(chan bool),
		depth: -1,
	}
	for _, option := range options {
		option(q)
	}
	if q.workers < 1 {
		q.workers = 1 + numCPU
	}
	go q.process()
	return q
}

// Submit a task to be run.
func (q *Queue) Submit(task Task) {
	q.in <- task
}

// Shutdown the queue. Does not return until all pending tasks have completed.
func (q *Queue) Shutdown() {
	close(q.in)
	<-q.done
}

func (q *Queue) process() {
	var received, processed uint64

	// Setup backlog
	var backlog []Task
	if q.depth > 0 {
		backlog = make([]Task, 0, q.depth)
	}

	// Setup workers
	ready := make(chan bool, q.workers)
	tasks := make(chan Task, q.workers)
	for i := 0; i < q.workers; i++ {
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
			if q.depth < 0 || len(backlog) < q.depth {
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

func (q *Queue) work(tasks <-chan Task, ready chan<- bool) {
	for task := range tasks {
		q.runTask(task)
		ready <- true
	}
}

func (q *Queue) runTask(task Task) {
	defer errs.Recovery(q.recoveryHandler)
	task()
}
