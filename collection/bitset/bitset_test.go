// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

// TestBitSetClearRangeBeyondStorage verifies that ClearRange fully clears its final allocated word when 'end' lies
// beyond the backing storage. Previously the loop bound was capped to the last allocated word but the per-word "last
// bit" was still derived from 'end', so the capped final word was only partially cleared, leaving set bits behind.
func TestBitSetClearRangeBeyondStorage(t *testing.T) {
	c := check.New(t)
	var bs BitSet
	bs.Set(230) // Allocates words 0..3; the highest allocated bit index is 255.
	bs.ClearRange(168, 328)
	c.False(bs.State(230), "bit 230 must be cleared")
	c.Equal(0, bs.Count())

	// A single set bit in the capped final word, cleared by a range that both starts and ends beyond storage in the
	// same final word's span.
	bs.Reset()
	bs.Set(200)
	bs.ClearRange(200, 100000)
	c.False(bs.State(200))
	c.Equal(0, bs.Count())
}

// TestBitSetClearRangeBruteForce differentially checks ClearRange against a straightforward per-bit reference across
// many start/end pairs, including ranges that extend well past the allocated storage.
func TestBitSetClearRangeBruteForce(t *testing.T) {
	c := check.New(t)
	// Keep the set bits within a few words (highest is 255, so words 0..3 are allocated) but drive 'end' far past that
	// allocation so the capped-final-word path is exercised with set bits beyond (end & bitIndexMask).
	const setLimit = 255
	const endLimit = 600
	for start := 0; start <= setLimit; start += 7 {
		for end := 0; end <= endLimit; end += 11 {
			var bs BitSet
			ref := make(map[int]bool)
			for bit := 0; bit <= setLimit; bit += 3 {
				bs.Set(bit)
				ref[bit] = true
			}
			bs.ClearRange(start, end)
			lo, hi := start, end
			if lo > hi {
				lo, hi = hi, lo
			}
			for bit := lo; bit <= hi; bit++ {
				delete(ref, bit)
			}
			for bit := 0; bit <= endLimit+dataBitsPerWord; bit++ {
				c.Equal(ref[bit], bs.State(bit), "bit %d after ClearRange(%d, %d)", bit, start, end)
			}
			c.Equal(len(ref), bs.Count(), "Count after ClearRange(%d, %d)", start, end)
		}
	}
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

// checkBitSetContents verifies that 'bs' has exactly the bits listed in 'indexes' set and no others. It checks the
// state of each expected bit, verifies Count() against both the expected number of bits and an independent tally made
// by walking the set bits with NextSet(), and compares 'bs' to an expected BitSet built purely from Set() calls. The
// Equal() comparison is made on trimmed copies, since Equal() also compares the underlying storage lengths, which
// legitimately vary with how the storage happened to be grown. 'indexes' must not contain duplicates.
func checkBitSetContents(c check.Checker, bs *BitSet, label string, indexes ...int) {
	c.Helper()
	var expected BitSet
	for _, index := range indexes {
		expected.Set(index)
	}
	for _, index := range indexes {
		c.True(bs.State(index), "%s: bit %d should be set", label, index)
	}
	tally := 0
	for i := bs.FirstSet(); i != -1; i = bs.NextSet(i + 1) {
		tally++
		c.True(expected.State(i), "%s: bit %d should not be set", label, i)
	}
	c.Equal(len(indexes), tally, "%s: number of set bits found by iterating with NextSet", label)
	c.Equal(len(indexes), bs.Count(), "%s: Count", label)
	if len(indexes) == 0 {
		c.Equal(-1, bs.FirstSet(), "%s: FirstSet of an empty BitSet", label)
		c.Equal(-1, bs.LastSet(), "%s: LastSet of an empty BitSet", label)
	}
	actual := bs.Clone()
	actual.Trim()
	expected.Trim()
	c.True(actual.Equal(&expected), "%s: result should equal the expected BitSet", label)
}

func TestBitSetAnd(t *testing.T) {
	c := check.New(t)

	// Overlapping sets within a single word
	var bs1, bs2 BitSet
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)
	bs1.Set(4)
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)
	bs1.And(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, single word", 3, 4)

	// Overlapping sets spanning multiple words
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{0, 65, 70, 130, 200} {
		bs1.Set(index)
	}
	for _, index := range []int{65, 130, 131, 255} {
		bs2.Set(index)
	}
	bs1.And(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, multiple words", 65, 130)

	// Disjoint sets
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{1, 3, 70, 129} {
		bs1.Set(index)
	}
	for _, index := range []int{2, 4, 71, 130} {
		bs2.Set(index)
	}
	bs1.And(&bs2)
	checkBitSetContents(c, &bs1, "disjoint")

	// The receiver has more words than 'other', so the bits beyond the last word of 'other' must be cleared
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{10, 100, 150, 250} {
		bs1.Set(index)
	}
	bs2.Set(10)
	bs2.Set(100)
	c.True(len(bs1.data) > len(bs2.data), "the receiver should have more words than 'other'")
	bs1.And(&bs2)
	c.False(bs1.State(150), "bit 150 lies beyond the last word of 'other' and must be cleared")
	c.False(bs1.State(250), "bit 250 lies beyond the last word of 'other' and must be cleared")
	checkBitSetContents(c, &bs1, "receiver larger than 'other'", 10, 100)

	// 'other' has more words than the receiver; the receiver must not grow
	bs1.Reset()
	bs2.Reset()
	bs1.Set(5)
	bs1.Set(70)
	for _, index := range []int{5, 70, 200} {
		bs2.Set(index)
	}
	words := len(bs1.data)
	c.True(len(bs2.data) > words, "'other' should have more words than the receiver")
	bs1.And(&bs2)
	c.Equal(words, len(bs1.data), "And must not grow the receiver")
	checkBitSetContents(c, &bs1, "'other' larger than receiver", 5, 70)

	// An empty receiver
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 200} {
		bs2.Set(index)
	}
	bs1.And(&bs2)
	c.Nil(bs1.data, "And must not allocate storage for an empty receiver")
	checkBitSetContents(c, &bs1, "empty receiver")

	// An empty 'other'
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 130, 200} {
		bs1.Set(index)
	}
	bs1.And(&bs2)
	checkBitSetContents(c, &bs1, "empty 'other'")

	// Both empty
	bs1.Reset()
	bs2.Reset()
	bs1.And(&bs2)
	checkBitSetContents(c, &bs1, "both empty")

	// Aliased call
	bs1.Reset()
	for _, index := range []int{3, 70, 130, 255} {
		bs1.Set(index)
	}
	bs1.And(&bs1)
	checkBitSetContents(c, &bs1, "aliased", 3, 70, 130, 255)
}

func TestBitSetOr(t *testing.T) {
	c := check.New(t)

	// Overlapping sets within a single word
	var bs1, bs2 BitSet
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)
	bs1.Or(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, single word", 1, 2, 3, 4, 5)

	// Overlapping sets spanning multiple words
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{0, 65, 70, 130, 200} {
		bs1.Set(index)
	}
	for _, index := range []int{65, 130, 131, 255} {
		bs2.Set(index)
	}
	bs1.Or(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, multiple words", 0, 65, 70, 130, 131, 200, 255)

	// Disjoint sets
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{0, 64, 129} {
		bs1.Set(index)
	}
	for _, index := range []int{1, 65, 130} {
		bs2.Set(index)
	}
	bs1.Or(&bs2)
	checkBitSetContents(c, &bs1, "disjoint", 0, 1, 64, 65, 129, 130)

	// The receiver has more words than 'other'; the bits beyond the last word of 'other' must be left alone
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 200} {
		bs1.Set(index)
	}
	bs2.Set(6)
	words := len(bs1.data)
	c.True(words > len(bs2.data), "the receiver should have more words than 'other'")
	bs1.Or(&bs2)
	c.Equal(words, len(bs1.data), "Or must not grow the receiver when 'other' is smaller")
	checkBitSetContents(c, &bs1, "receiver larger than 'other'", 5, 6, 70, 200)

	// 'other' has more words than the receiver, so the receiver must grow
	bs1.Reset()
	bs2.Reset()
	bs1.Set(5)
	bs2.Set(5)
	bs2.Set(130)
	c.True(len(bs2.data) > len(bs1.data), "'other' should have more words than the receiver")
	bs1.Or(&bs2)
	c.True(len(bs1.data) >= len(bs2.data), "Or must grow the receiver when 'other' has more words")
	c.True(bs1.State(130), "bit 130 came from the extra words of 'other'")
	checkBitSetContents(c, &bs1, "'other' larger than receiver", 5, 130)

	// An empty receiver
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 200} {
		bs2.Set(index)
	}
	bs1.Or(&bs2)
	checkBitSetContents(c, &bs1, "empty receiver", 5, 70, 200)

	// An empty 'other'
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 130, 200} {
		bs1.Set(index)
	}
	words = len(bs1.data)
	bs1.Or(&bs2)
	c.Equal(words, len(bs1.data), "Or with an empty 'other' must not grow the receiver")
	checkBitSetContents(c, &bs1, "empty 'other'", 5, 70, 130, 200)

	// Both empty
	bs1.Reset()
	bs2.Reset()
	bs1.Or(&bs2)
	checkBitSetContents(c, &bs1, "both empty")

	// Aliased call
	bs1.Reset()
	for _, index := range []int{3, 70, 130, 255} {
		bs1.Set(index)
	}
	bs1.Or(&bs1)
	checkBitSetContents(c, &bs1, "aliased", 3, 70, 130, 255)
}

func TestBitSetXor(t *testing.T) {
	c := check.New(t)

	// Overlapping sets within a single word
	var bs1, bs2 BitSet
	bs1.Set(1)
	bs1.Set(2)
	bs1.Set(3)
	bs2.Set(3)
	bs2.Set(4)
	bs2.Set(5)
	bs1.Xor(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, single word", 1, 2, 4, 5)

	// Overlapping sets spanning multiple words
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{0, 65, 70, 130, 200} {
		bs1.Set(index)
	}
	for _, index := range []int{65, 130, 131, 255} {
		bs2.Set(index)
	}
	bs1.Xor(&bs2)
	checkBitSetContents(c, &bs1, "overlapping, multiple words", 0, 70, 131, 200, 255)

	// Disjoint sets
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{0, 64, 129} {
		bs1.Set(index)
	}
	for _, index := range []int{1, 65, 130} {
		bs2.Set(index)
	}
	bs1.Xor(&bs2)
	checkBitSetContents(c, &bs1, "disjoint", 0, 1, 64, 65, 129, 130)

	// The receiver has more words than 'other'; the bits beyond the last word of 'other' must be left alone
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 200} {
		bs1.Set(index)
	}
	bs2.Set(5)
	bs2.Set(71)
	words := len(bs1.data)
	c.True(words > len(bs2.data), "the receiver should have more words than 'other'")
	bs1.Xor(&bs2)
	c.Equal(words, len(bs1.data), "Xor must not grow the receiver when 'other' is smaller")
	checkBitSetContents(c, &bs1, "receiver larger than 'other'", 70, 71, 200)

	// 'other' has more words than the receiver, so the receiver must grow
	bs1.Reset()
	bs2.Reset()
	bs1.Set(5)
	bs2.Set(5)
	bs2.Set(130)
	c.True(len(bs2.data) > len(bs1.data), "'other' should have more words than the receiver")
	bs1.Xor(&bs2)
	c.True(len(bs1.data) >= len(bs2.data), "Xor must grow the receiver when 'other' has more words")
	c.True(bs1.State(130), "bit 130 came from the extra words of 'other'")
	checkBitSetContents(c, &bs1, "'other' larger than receiver", 130)

	// An empty receiver
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 200} {
		bs2.Set(index)
	}
	bs1.Xor(&bs2)
	checkBitSetContents(c, &bs1, "empty receiver", 5, 70, 200)

	// An empty 'other'
	bs1.Reset()
	bs2.Reset()
	for _, index := range []int{5, 70, 130, 200} {
		bs1.Set(index)
	}
	words = len(bs1.data)
	bs1.Xor(&bs2)
	c.Equal(words, len(bs1.data), "Xor with an empty 'other' must not grow the receiver")
	checkBitSetContents(c, &bs1, "empty 'other'", 5, 70, 130, 200)

	// Both empty
	bs1.Reset()
	bs2.Reset()
	bs1.Xor(&bs2)
	checkBitSetContents(c, &bs1, "both empty")

	// Aliased call: a BitSet XOR'd with itself must be empty
	bs1.Reset()
	for _, index := range []int{3, 70, 130, 255} {
		bs1.Set(index)
	}
	bs1.Xor(&bs1)
	c.Equal(0, bs1.Count(), "aliased: Count")
	c.Equal(-1, bs1.FirstSet(), "aliased: FirstSet")
	checkBitSetContents(c, &bs1, "aliased")
}

// TestBitSetLogicalOpsBruteForce differentially checks And, Or and Xor against straightforward per-bit references
// across a variety of receiver and argument shapes, including empty sets and both orderings of "one side has more
// words than the other".
func TestBitSetLogicalOpsBruteForce(t *testing.T) {
	c := check.New(t)
	// The specs are (step, highest) pairs; a step of 0 produces an empty set. The highest values straddle the word
	// boundaries at 64, 128, 192 and 256 so that the various word-length combinations are exercised.
	specs := [][2]int{{0, 0}, {1, 5}, {3, 63}, {1, 64}, {7, 130}, {5, 255}, {64, 255}}
	build := func(spec [2]int) (*BitSet, map[int]bool) {
		bs := &BitSet{}
		ref := make(map[int]bool)
		if spec[0] > 0 {
			for bit := 0; bit <= spec[1]; bit += spec[0] {
				bs.Set(bit)
				ref[bit] = true
			}
		}
		return bs, ref
	}
	const limit = 5 * dataBitsPerWord // Deliberately beyond the largest allocation any spec produces
	for _, lhsSpec := range specs {
		for _, rhsSpec := range specs {
			for _, op := range []string{"And", "Or", "Xor"} {
				bs, ref := build(lhsSpec)
				other, otherRef := build(rhsSpec)
				expected := make(map[int]bool)
				for bit := range limit {
					switch op {
					case "And":
						if ref[bit] && otherRef[bit] {
							expected[bit] = true
						}
					case "Or":
						if ref[bit] || otherRef[bit] {
							expected[bit] = true
						}
					default:
						if ref[bit] != otherRef[bit] {
							expected[bit] = true
						}
					}
				}
				switch op {
				case "And":
					bs.And(other)
				case "Or":
					bs.Or(other)
				default:
					bs.Xor(other)
				}
				for bit := range limit {
					c.Equal(expected[bit], bs.State(bit), "bit %d after %v.%s(%v)", bit, lhsSpec, op, rhsSpec)
				}
				c.Equal(len(expected), bs.Count(), "Count after %v.%s(%v)", lhsSpec, op, rhsSpec)
			}
		}
	}
}
