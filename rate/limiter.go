// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package rate provides rate limiting which supports a hierarchy of limiters, each capped by their parent.
package rate

import (
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// Limiter provides a rate limiter.
type Limiter interface {
	// New returns a child limiter whose effective capacity is also capped by this limiter's. A capacity less than 1 is
	// treated as 1. If this limiter is closed, it returns itself so that chained calls report the closed state rather
	// than panicking on a nil interface.
	New(capacity int) Limiter

	// Cap returns the capacity per time period. If applyParentCaps is true, the result is also capped by every
	// ancestor's capacity.
	Cap(applyParentCaps bool) int

	// SetCap sets the capacity. A capacity less than 1 is treated as 1.
	SetCap(capacity int)

	// LastUsed returns the capacity used in the last time period.
	LastUsed() int

	// Use returns a channel that receives nil once the request is granted, or an error if it cannot be fulfilled.
	Use(amount int) <-chan error

	// Closed returns true if the limiter is closed.
	Closed() bool

	// Close this limiter and any children it may have.
	Close()
}

type limiter struct {
	controller *controller
	parent     *limiter
	children   []*limiter
	last       int
	capacity   int
	used       int
	closed     bool
}

type controller struct {
	root    *limiter
	ticker  *time.Ticker
	done    chan bool
	waiting []*request
	lock    sync.RWMutex
}

type request struct {
	limiter *limiter
	done    chan error
	amount  int
}

// New creates a top-level rate limiter that allows 'capacity' units (bytes, for example) per 'period'. A capacity less
// than 1 is treated as 1.
func New(capacity int, period time.Duration) Limiter {
	c := &controller{
		ticker: time.NewTicker(period),
		done:   make(chan bool, 1),
	}
	l := &limiter{
		controller: c,
		capacity:   max(capacity, 1),
	}
	c.root = l
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.lock.Lock()
				c.root.reset()
				remaining := make([]*request, 0, len(c.waiting))
				for _, req := range c.waiting {
					if req.limiter.closed {
						req.done <- errs.New("Limiter is closed")
						continue
					}
					if capped := req.limiter.cappedCapacity(); req.amount > capped {
						req.done <- errs.Newf("Amount (%d) is greater than capacity (%d)", req.amount, capped)
						continue
					}
					if req.limiter.tryConsume(req.amount) {
						req.done <- nil
						continue
					}
					remaining = append(remaining, req)
				}
				c.waiting = remaining
				c.lock.Unlock()
			case <-c.done:
				c.ticker.Stop()
				c.lock.Lock()
				for _, req := range c.waiting {
					req.done <- errs.New("Limiter is closed")
				}
				c.waiting = make([]*request, 0)
				c.lock.Unlock()
				return
			}
		}
	}()
	return l
}

func (l *limiter) New(capacity int) Limiter {
	l.controller.lock.Lock()
	defer l.controller.lock.Unlock()
	if l.closed {
		// Return the closed limiter itself rather than nil so that chained calls such as parent.New(n).Use(...) report
		// the closed state instead of panicking on a nil interface.
		return l
	}
	child := &limiter{
		controller: l.controller,
		parent:     l,
		capacity:   max(capacity, 1),
	}
	l.children = append(l.children, child)
	return child
}

func (l *limiter) Cap(applyParentCaps bool) int {
	l.controller.lock.RLock()
	defer l.controller.lock.RUnlock()
	if applyParentCaps {
		return l.cappedCapacity()
	}
	return l.capacity
}

// cappedCapacity returns this limiter's capacity capped by every ancestor's capacity. The controller lock must be held.
func (l *limiter) cappedCapacity() int {
	capacity := l.capacity
	for p := l.parent; p != nil; p = p.parent {
		if p.capacity < capacity {
			capacity = p.capacity
		}
	}
	return capacity
}

func (l *limiter) SetCap(capacity int) {
	l.controller.lock.Lock()
	l.capacity = max(capacity, 1)
	l.controller.lock.Unlock()
}

func (l *limiter) LastUsed() int {
	l.controller.lock.RLock()
	defer l.controller.lock.RUnlock()
	return l.last
}

func (l *limiter) Use(amount int) <-chan error {
	done := make(chan error, 1)
	if amount < 0 {
		done <- errs.Newf("Amount (%d) must be positive", amount)
		return done
	}
	if amount == 0 {
		done <- nil
		return done
	}
	l.controller.lock.Lock()
	if l.closed {
		l.controller.lock.Unlock()
		done <- errs.New("Limiter is closed")
		return done
	}
	if capacity := l.cappedCapacity(); amount > capacity {
		l.controller.lock.Unlock()
		done <- errs.Newf("Amount (%d) is greater than capacity (%d)", amount, capacity)
		return done
	}
	// Take the fast path only when nothing is queued, preserving FIFO order. Otherwise a steady stream of smaller
	// requests could keep consuming capacity ahead of an earlier, larger queued request and starve it.
	if len(l.controller.waiting) == 0 && l.tryConsume(amount) {
		l.controller.lock.Unlock()
		done <- nil
		return done
	}
	l.controller.waiting = append(l.controller.waiting, &request{
		limiter: l,
		amount:  amount,
		done:    done,
	})
	l.controller.lock.Unlock()
	return done
}

// tryConsume adds amount to the used count of this limiter and all of its ancestors if every one of them has room,
// reporting whether it did so. The controller lock must be held.
func (l *limiter) tryConsume(amount int) bool {
	available := l.capacity - l.used
	for p := l.parent; p != nil; p = p.parent {
		if pa := p.capacity - p.used; pa < available {
			available = pa
		}
	}
	if available < amount {
		return false
	}
	l.used += amount
	for p := l.parent; p != nil; p = p.parent {
		p.used += amount
	}
	return true
}

func (l *limiter) reset() {
	l.last = l.used
	l.used = 0
	for _, child := range l.children {
		child.reset()
	}
}

func (l *limiter) Closed() bool {
	l.controller.lock.RLock()
	defer l.controller.lock.RUnlock()
	return l.closed
}

func (l *limiter) Close() {
	l.controller.lock.Lock()
	if !l.closed {
		l.close()
		if l.parent == nil {
			l.controller.done <- true
		} else {
			for i, child := range l.parent.children {
				if child != l {
					continue
				}
				j := len(l.parent.children) - 1
				l.parent.children[i] = l.parent.children[j]
				l.parent.children[j] = nil
				l.parent.children = l.parent.children[:j]
				break
			}
			l.closed = true
		}
	}
	l.controller.lock.Unlock()
}

func (l *limiter) close() {
	l.closed = true
	for _, child := range l.children {
		child.close()
	}
}
