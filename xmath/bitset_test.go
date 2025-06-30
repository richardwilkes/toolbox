// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xmath

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestBitSet(t *testing.T) {
	var bs BitSet
	c := check.New(t)
	c.Equal(0, bs.Count())
	bs.Set(0)
	c.Equal(1, bs.Count())
	bs.Set(7)
	c.Equal(2, bs.Count())
	bs.Set(dataBitsPerWord - 1)
	c.Equal(3, bs.Count())
	bs.Set(dataBitsPerWord)
	c.Equal(4, bs.Count())
	bs.Set(dataBitsPerWord + 1)
	c.Equal(5, bs.Count())
	bs.Set(0)
	c.Equal(5, bs.Count())
	bs.Clear(0)
	c.Equal(4, bs.Count())
	bs.Clear(1)
	c.Equal(4, bs.Count())
	bs.Clear(1000)
	c.Equal(4, bs.Count())
	c.False(bs.State(0))
	c.False(bs.State(1))
	c.True(bs.State(7))
	c.False(bs.State(77))
	c.True(bs.State(dataBitsPerWord))
	bs.Flip(22)
	c.True(bs.State(22))
	bs.Flip(22)
	c.False(bs.State(22))
	c.Equal(7, bs.NextSet(0))
	c.Equal(7, bs.NextSet(7))
	c.Equal(dataBitsPerWord-1, bs.NextSet(8))
	c.Equal(dataBitsPerWord, bs.NextSet(dataBitsPerWord))
	bs.Set(1234)
	c.Equal(1234, bs.NextSet(dataBitsPerWord+2))
	c.Equal(0, bs.NextClear(0))
	c.Equal(dataBitsPerWord+2, bs.NextClear(dataBitsPerWord-1))
	c.Equal(1235, bs.NextClear(1234))
	bs.Set(dataBitsPerWord*100 - 1)
	c.Equal(dataBitsPerWord*100, bs.NextClear(dataBitsPerWord*100-1))
	c.Equal(dataBitsPerWord*100-1, bs.PreviousSet(dataBitsPerWord*100))
	c.Equal(1234, bs.PreviousSet(dataBitsPerWord*100-2))
	c.Equal(-1, bs.PreviousSet(0))
	c.Equal(dataBitsPerWord*1000, bs.PreviousClear(dataBitsPerWord*1000))
	c.Equal(dataBitsPerWord*100-2, bs.PreviousClear(dataBitsPerWord*100-1))
	c.Equal(0, bs.PreviousClear(0))
	bs.Set(0)
	c.Equal(-1, bs.PreviousClear(0))

	bs.Reset()
	bs.Set(65)
	bs.SetRange(10, 300)
	c.Equal(291, bs.Count())
	for i := 10; i < 301; i++ {
		c.True(bs.State(i))
	}
	c.Equal(301, bs.NextClear(10))
	c.Equal(9, bs.PreviousClear(300))
	c.Equal(10, bs.NextSet(0))
	c.Equal(300, bs.PreviousSet(1000))
	bs.ClearRange(15, 295)
	c.Equal(10, bs.Count())
	for i := 15; i < 296; i++ {
		c.False(bs.State(i))
	}
	bs.FlipRange(10, 300)
	c.Equal(281, bs.Count())
}
