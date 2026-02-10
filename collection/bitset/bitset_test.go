// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package bitset

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
	c.Equal(7, bs.FirstSet())
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
	bs.SetRange(300, 10)
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
	bs.ClearRange(295, 15)
	c.Equal(10, bs.Count())
	for i := 15; i < 296; i++ {
		c.False(bs.State(i))
	}
	bs.FlipRange(10, 300)
	c.Equal(281, bs.Count())
	c.Equal(295, bs.LastSet())
	bs.FlipRange(300, 10)
	c.Equal(10, bs.Count())
	c.Equal(300, bs.LastSet())
	c.Equal(-1, bs.NextSet(301))
}

func TestBitSetEqual(t *testing.T) {
	c := check.New(t)

	// Test nil comparison
	var bs1 BitSet
	c.False(bs1.Equal(nil))

	// Test empty BitSets
	var bs2 BitSet
	c.True(bs1.Equal(&bs2))

	// Test same BitSets with same bits set
	bs1.Set(5)
	bs1.Set(10)
	bs1.Set(100)

	bs2.Set(5)
	bs2.Set(10)
	bs2.Set(100)
	c.True(bs1.Equal(&bs2))

	// Test different set counts
	bs2.Set(200)
	c.False(bs1.Equal(&bs2))

	// Test same count but different bits
	bs1.Clear(5)
	bs1.Set(200)
	bs2.Clear(10)
	c.False(bs1.Equal(&bs2))

	// Test different underlying data lengths but same logical content
	var bs3, bs4 BitSet
	bs3.Set(5)
	bs4.Set(5)
	bs4.Set(100)   // This will expand the data array
	bs4.Clear(100) // Clear it but data array remains larger
	c.False(bs3.Equal(&bs4))
	bs3.EnsureCapacity(2)
	c.True(bs3.Equal(&bs4))

	// Test self equality
	c.True(bs1.Equal(&bs1))

	// Test cloned BitSets
	bs5 := bs1.Clone()
	c.True(bs1.Equal(bs5))

	// Test copied BitSets
	var bs6 BitSet
	bs6.Copy(&bs1)
	c.True(bs1.Equal(&bs6))
}

func TestBitSetLoad(t *testing.T) {
	c := check.New(t)

	// Test loading empty data
	var bs BitSet
	bs.Set(5) // Set some bits first
	bs.Load([]uint64{})
	c.Equal(0, bs.Count())
	c.Nil(bs.data)

	// Test loading single word with some bits set
	bs.Load([]uint64{0b1010001})
	c.Equal(3, bs.Count())
	c.True(bs.State(0))
	c.False(bs.State(1))
	c.False(bs.State(2))
	c.False(bs.State(3))
	c.True(bs.State(4))
	c.False(bs.State(5))
	c.True(bs.State(6))

	// Test loading multiple words
	data := []uint64{
		0b1100000000000000000000000000000000000000000000000000000000000001, // word 0: bits 0 and 62, 63
		0b0000000000000000000000000000000000000000000000000000000000001010, // word 1: bits 65 and 67
		0, // word 2: no bits set
		0b1000000000000000000000000000000000000000000000000000000000000000, // word 3: bit 63 (bit 255 overall)
	}
	bs.Load(data)
	c.Equal(6, bs.Count())
	c.True(bs.State(0))
	c.True(bs.State(62))
	c.True(bs.State(63))
	c.True(bs.State(65))
	c.True(bs.State(67))
	c.True(bs.State(255))
	c.False(bs.State(1))
	c.False(bs.State(64))
	c.False(bs.State(66))
	c.False(bs.State(128))
	c.False(bs.State(254))
	c.Equal(data, bs.Data())

	// Test loading data with trailing zeros (should be trimmed)
	dataWithZeros := []uint64{0b101, 0, 0, 0}
	bs.Load(dataWithZeros)
	c.Equal(2, bs.Count())
	c.True(bs.State(0))
	c.True(bs.State(2))
	c.Equal(1, len(bs.data)) // Should be trimmed to 1 word

	// Test loading all zeros
	bs.Load([]uint64{0, 0, 0})
	c.Equal(0, bs.Count())
	c.Nil(bs.data) // Should be trimmed to nil

	// Test loading max uint64
	bs.Load([]uint64{^uint64(0)})
	c.Equal(64, bs.Count())
	for i := range 64 {
		c.True(bs.State(i))
	}

	// Test that Load replaces existing data completely
	bs.Set(100)
	bs.Set(200)
	originalCount := bs.Count()
	c.True(originalCount > 64) // Should be more than 64 from previous test
	bs.Load([]uint64{0b11})
	c.Equal(2, bs.Count())
	c.True(bs.State(0))
	c.True(bs.State(1))
	c.False(bs.State(2))
	c.False(bs.State(100))
	c.False(bs.State(200))

	// Test loading nil data (should behave like empty slice)
	bs.Load(nil)
	c.Equal(0, bs.Count())
	c.Nil(bs.data)
}
