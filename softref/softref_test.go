// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package softref_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/softref"
)

type res struct {
	released chan<- string
	key      string
}

func newRes(key string, released chan<- string) *res {
	return &res{
		key:      key,
		released: released,
	}
}

func (r *res) Key() string {
	return r.key
}

func (r *res) Release() {
	r.released <- r.key
}

func TestSoftRef(t *testing.T) {
	c := check.New(t)
	p := softref.NewPool()
	ch := make(chan string, 128)
	sr1, existed := p.NewSoftRef(newRes("1", ch))
	c.False(existed)
	_, existed = p.NewSoftRef(newRes("2", ch))
	c.False(existed)
	sr3, existed := p.NewSoftRef(newRes("3", ch))
	c.False(existed)
	r4 := newRes("4", ch)
	sr4a, existed := p.NewSoftRef(r4)
	c.False(existed)
	sr4b, existed := p.NewSoftRef(r4)
	c.True(existed)
	lookFor(c, "2", ch)
	get, existed := sr3.Resource.(*res)
	c.True(existed)
	lookFor(c, get.key, ch)
	get, existed = sr1.Resource.(*res)
	c.True(existed)
	lookFor(c, get.key, ch)
	get, existed = sr4a.Resource.(*res)
	c.True(existed)
	get2, existed2 := sr4b.Resource.(*res)
	c.True(existed2)
	c.Equal(get.key, get2.key)
	lookForExpectingTimeout(c, ch)
	c.Equal("4", sr4b.Key) // Keeps refs to r4 alive for the above call
	lookFor(c, get.key, ch)
}

// inertRes is a resource whose Release does nothing.
type inertRes struct {
	key string
}

func (r *inertRes) Key() string { return r.key }
func (r *inertRes) Release()    {}

// reentrantRes re-enters its pool from within Release by acquiring another SoftRef. With the old implementation this
// ran while the finalizer goroutine held the pool lock, deadlocking on the non-reentrant mutex.
type reentrantRes struct {
	pool *softref.Pool
	done chan struct{}
	kept *softref.SoftRef // Keeps the re-entrant ref alive so it isn't finalized during the test.
	key  string
}

func (r *reentrantRes) Key() string { return r.key }

func (r *reentrantRes) Release() {
	r.kept, _ = r.pool.NewSoftRef(&inertRes{key: r.key + "-inner"})
	close(r.done)
}

func TestSoftRefReleaseReentrantNoDeadlock(t *testing.T) {
	c := check.New(t)
	p := softref.NewPool()
	done := make(chan struct{})

	// Create and immediately drop the only reference, so the finalizer will run Release(), which re-enters the pool.
	func() {
		sr, existed := p.NewSoftRef(&reentrantRes{pool: p, key: "reentrant", done: done})
		c.False(existed)
		runtime.KeepAlive(sr)
	}()

	deadline := time.After(5 * time.Second)
	for {
		runtime.GC()
		select {
		case <-done:
			return // Release() re-entered the pool and completed without deadlocking.
		case <-deadline:
			c.Errorf("Release() that re-entered the pool deadlocked; the finalizer held the pool lock across Release()")
			return
		case <-time.After(20 * time.Millisecond):
			// The finalizer may not have run yet; loop and GC again.
		}
	}
}

func lookFor(c check.Checker, key string, ch <-chan string) {
	c.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
		c.Errorf("timed out waiting for %s", key)
	case k := <-ch:
		c.Equal(key, k)
	}
}

func lookForExpectingTimeout(c check.Checker, ch <-chan string) {
	c.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
	case k := <-ch:
		c.Errorf("received key '%s' when none expected", k)
	}
}
