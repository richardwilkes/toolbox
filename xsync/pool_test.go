// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xsync

import (
	"runtime/debug"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestPoolWithNil(t *testing.T) {
	c := check.New(t)
	defer func() { c.NotNil(recover()) }()
	NewPool[any](nil)
}

func TestEmptyPoolCallsNew(t *testing.T) {
	var i int
	p := NewPool(func() int {
		i++
		return i
	})
	c := check.New(t)
	c.Equal(1, p.Get())
	c.Equal(2, p.Get())
	c.Equal(2, i, "Should be the number of times Get was called")
}

func TestPoolPutGet(t *testing.T) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "-race" && setting.Value == "true" {
				// When -race is enabled, one quarter of the Put requests get randomly dropped, so can't reliably test
				return
			}
		}
	}
	var i int
	p := NewPool(func() int {
		i++
		return i
	})
	p.Put(10)
	p.Put(20)
	g1 := p.Get()
	g2 := p.Get()
	c := check.New(t)
	c.True((g1 == 10 && g2 == 20) || (g1 == 20 && g2 == 10), "Any order is fine")
	c.Equal(1, p.Get(), "Getting from an empty pool should call the new function")
	p.Put(30)
	c.Equal(30, p.Get(), "Should return the only value in the pool, which we just placed there")
}
