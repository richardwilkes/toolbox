// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rate_test

import (
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/rate"
)

func TestCap(t *testing.T) {
	c := check.New(t)
	rl := rate.New(50*1024, time.Second)
	c.Equal(50*1024, rl.Cap(true))
	sub := rl.New(100 * 1024)
	c.Equal(100*1024, sub.Cap(false))
	c.Equal(50*1024, sub.Cap(true))
	sub.SetCap(1024)
	c.Equal(1024, sub.Cap(true))
	rl.Close()
	c.True(sub.Closed())
	c.True(rl.Closed())
}

func TestUse(t *testing.T) {
	c := check.New(t)
	rl := rate.New(100, 100*time.Millisecond)
	endAfter := time.Now().Add(250 * time.Millisecond)
	for endAfter.After(time.Now()) {
		err := <-rl.Use(1)
		c.NoError(err)
	}
	c.Equal(100, rl.LastUsed())
	rl.Close()
	c.True(rl.Closed())
}
