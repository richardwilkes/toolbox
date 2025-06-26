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

	"github.com/richardwilkes/toolbox/check"
)

func TestBitSet(t *testing.T) {
	var bs BitSet
	check.Equal(t, 0, bs.Count())
	bs.Set(0)
	check.Equal(t, 1, bs.Count())
	bs.Set(7)
	check.Equal(t, 2, bs.Count())
	bs.Set(dataBitsPerWord - 1)
	check.Equal(t, 3, bs.Count())
	bs.Set(dataBitsPerWord)
	check.Equal(t, 4, bs.Count())
	bs.Set(dataBitsPerWord + 1)
	check.Equal(t, 5, bs.Count())
	bs.Set(0)
	check.Equal(t, 5, bs.Count())
	bs.Clear(0)
	check.Equal(t, 4, bs.Count())
	bs.Clear(1)
	check.Equal(t, 4, bs.Count())
	bs.Clear(1000)
	check.Equal(t, 4, bs.Count())
	check.False(t, bs.State(0))
	check.False(t, bs.State(1))
	check.True(t, bs.State(7))
	check.False(t, bs.State(77))
	check.True(t, bs.State(dataBitsPerWord))
	bs.Flip(22)
	check.True(t, bs.State(22))
	bs.Flip(22)
	check.False(t, bs.State(22))
	check.Equal(t, 7, bs.NextSet(0))
	check.Equal(t, 7, bs.NextSet(7))
	check.Equal(t, dataBitsPerWord-1, bs.NextSet(8))
	check.Equal(t, dataBitsPerWord, bs.NextSet(dataBitsPerWord))
	bs.Set(1234)
	check.Equal(t, 1234, bs.NextSet(dataBitsPerWord+2))
	check.Equal(t, 0, bs.NextClear(0))
	check.Equal(t, dataBitsPerWord+2, bs.NextClear(dataBitsPerWord-1))
	check.Equal(t, 1235, bs.NextClear(1234))
	bs.Set(dataBitsPerWord*100 - 1)
	check.Equal(t, dataBitsPerWord*100, bs.NextClear(dataBitsPerWord*100-1))
	check.Equal(t, dataBitsPerWord*100-1, bs.PreviousSet(dataBitsPerWord*100))
	check.Equal(t, 1234, bs.PreviousSet(dataBitsPerWord*100-2))
	check.Equal(t, -1, bs.PreviousSet(0))
	check.Equal(t, dataBitsPerWord*1000, bs.PreviousClear(dataBitsPerWord*1000))
	check.Equal(t, dataBitsPerWord*100-2, bs.PreviousClear(dataBitsPerWord*100-1))
	check.Equal(t, 0, bs.PreviousClear(0))
	bs.Set(0)
	check.Equal(t, -1, bs.PreviousClear(0))

	bs.Reset()
	bs.Set(65)
	bs.SetRange(10, 300)
	check.Equal(t, 291, bs.Count())
	for i := 10; i < 301; i++ {
		check.True(t, bs.State(i))
	}
	check.Equal(t, 301, bs.NextClear(10))
	check.Equal(t, 9, bs.PreviousClear(300))
	check.Equal(t, 10, bs.NextSet(0))
	check.Equal(t, 300, bs.PreviousSet(1000))
	bs.ClearRange(15, 295)
	check.Equal(t, 10, bs.Count())
	for i := 15; i < 296; i++ {
		check.False(t, bs.State(i))
	}
	bs.FlipRange(10, 300)
	check.Equal(t, 281, bs.Count())
}
