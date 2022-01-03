// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/softref"
	"github.com/stretchr/testify/assert"
)

type res struct {
	key      string
	released chan<- string
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
	p := softref.NewPool(&jot.Logger{})
	ch := make(chan string, 128)
	sr1, existed := p.NewSoftRef(newRes("1", ch))
	assert.False(t, existed)
	_, existed = p.NewSoftRef(newRes("2", ch))
	assert.False(t, existed)
	sr3, existed := p.NewSoftRef(newRes("3", ch))
	assert.False(t, existed)
	r4 := newRes("4", ch)
	sr4a, existed := p.NewSoftRef(r4)
	assert.False(t, existed)
	sr4b, existed := p.NewSoftRef(r4)
	assert.True(t, existed)
	lookFor(t, "2", ch)
	key := sr3.Resource.(*res).key
	lookFor(t, key, ch)
	key = sr1.Resource.(*res).key
	lookFor(t, key, ch)
	key = sr4a.Resource.(*res).key
	assert.Equal(t, key, sr4b.Resource.(*res).key)
	lookForExpectingTimeout(t, ch)
	assert.Equal(t, "4", sr4b.Key) // Keeps refs to r4 alive for the above call
	lookFor(t, key, ch)
}

func lookFor(t *testing.T, key string, ch <-chan string) {
	t.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
		assert.Failf(t, "timed out waiting for %s", key)
	case k := <-ch:
		assert.Equal(t, key, k)
	}
}

func lookForExpectingTimeout(t *testing.T, ch <-chan string) {
	t.Helper()
	runtime.GC()
	select {
	case <-time.After(time.Second):
	case k := <-ch:
		assert.Failf(t, "received key '%s' when none expected", k)
	}
}
