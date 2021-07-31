// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package rate provides rate limiting which supports a hierarchy of limiters,
// each capped by their parent.
package rate

import (
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

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
	lock    sync.RWMutex
	waiting []*request
}

type request struct {
	limiter *limiter
	amount  int
	done    chan error
}

// New creates a new top-level rate limiter. 'capacity' is the number of units
// (bytes, for example) allowed to be used in a particular time 'period'.
func New(capacity int, period time.Duration) Limiter {
	c := &controller{
		ticker: time.NewTicker(period),
		done:   make(chan bool),
	}
	l := &limiter{
		controller: c,
		capacity:   capacity,
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
					if req.amount > req.limiter.capacity {
						req.done <- errs.Newf("Amount (%d) is greater than capacity (%d)", req.amount, req.limiter.capacity)
						continue
					}
					if c.root.capacity-c.root.used > 0 {
						available := req.limiter.capacity - req.limiter.used
						p := req.limiter.parent
						for p != nil {
							pa := p.capacity - p.used
							if pa < available {
								available = pa
							}
							p = p.parent
						}
						if available >= req.amount {
							req.limiter.used += req.amount
							p = req.limiter.parent
							for p != nil {
								p.used += req.amount
								p = p.parent
							}
							req.done <- nil
							continue
						}
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
		return nil
	}
	child := &limiter{
		controller: l.controller,
		parent:     l,
		capacity:   capacity,
	}
	l.children = append(l.children, child)
	return child
}

func (l *limiter) Cap(applyParentCaps bool) int {
	l.controller.lock.RLock()
	defer l.controller.lock.RUnlock()
	capacity := l.capacity
	if applyParentCaps {
		p := l.parent
		for p != nil {
			if p.capacity < capacity {
				capacity = p.capacity
			}
			p = p.parent
		}
	}
	return capacity
}

func (l *limiter) SetCap(capacity int) {
	l.controller.lock.Lock()
	l.capacity = capacity
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
	if amount > l.capacity {
		capacity := l.capacity
		l.controller.lock.Unlock()
		done <- errs.Newf("Amount (%d) is greater than capacity (%d)", amount, capacity)
		return done
	}
	available := l.capacity - l.used
	p := l.parent
	for p != nil {
		pa := p.capacity - p.used
		if pa < available {
			available = pa
		}
		p = p.parent
	}
	if available >= amount {
		l.used += amount
		p = l.parent
		for p != nil {
			p.used += amount
			p = p.parent
		}
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
				if child == l { //nolint:gocritic
					j := len(l.parent.children) - 1
					l.parent.children[i] = l.parent.children[j]
					l.parent.children[j] = nil
					l.parent.children = l.parent.children[:j]
					break
				}
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
