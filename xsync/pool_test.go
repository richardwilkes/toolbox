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
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestPoolWithNil(t *testing.T) {
	defer func() { check.NotNil(t, recover()) }()
	NewPool[any](nil)
}

func TestEmptyPoolCallsNew(t *testing.T) {
	var i int
	p := NewPool(func() int {
		i++
		return i
	})
	check.Equal(t, 1, p.Get())
	check.Equal(t, 2, p.Get())
	check.Equal(t, 2, i, "Should be the number of times Get was called")
}

func TestPoolPutGet(t *testing.T) {
	var i int
	p := NewPool(func() int {
		i++
		return i
	})
	p.Put(10)
	p.Put(20)
	g1 := p.Get()
	g2 := p.Get()
	check.True(t, (g1 == 10 && g2 == 20) || (g1 == 20 && g2 == 10), "Any order is fine")
	check.Equal(t, 1, p.Get(), "Getting from an empty pool should call the new function")
	p.Put(30)
	check.Equal(t, 30, p.Get(), "Should return the only value in the pool, which we just placed there")
}
