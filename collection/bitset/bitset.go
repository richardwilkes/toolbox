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
	"math/bits"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xos"
)

const (
	addressBitsPerWord = 6
	dataBitsPerWord    = 1 << addressBitsPerWord
	bitIndexMask       = dataBitsPerWord - 1
)

// BitSet contains a set of bits.
type BitSet struct {
	data []uint64
	set  int
}

// Clone this BitSet.
func (b *BitSet) Clone() *BitSet {
	bs := &BitSet{data: make([]uint64, len(b.data)), set: b.set}
	copy(bs.data, b.data)
	return bs
}

// Copy the content of 'other' into this BitSet, making them equal.
func (b *BitSet) Copy(other *BitSet) {
	b.set = other.set
	b.data = make([]uint64, len(other.data))
	copy(b.data, other.data)
}

// Equal returns true if this BitSet is equal to 'other'.
func (b *BitSet) Equal(other *BitSet) bool {
	if other == nil {
		return false
	}
	if b.set != other.set {
		return false
	}
	if len(b.data) != len(other.data) {
		return false
	}
	for i := range b.data {
		if b.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

// Count returns the number of set bits.
func (b *BitSet) Count() int {
	return b.set
}

// State returns the state of the bit at 'index'.
func (b *BitSet) State(index int) bool {
	validateBitSetIndex(index)
	i := index >> addressBitsPerWord
	if i >= len(b.data) {
		return false
	}
	mask := wordMask(index)
	return b.data[i]&mask == mask
}

// Set the bit at 'index'.
func (b *BitSet) Set(index int) {
	validateBitSetIndex(index)
	i := index >> addressBitsPerWord
	b.EnsureCapacity(i + 1)
	mask := wordMask(index)
	if b.data[i]&mask == 0 {
		b.data[i] |= mask
		b.set++
	}
}

// SetRange sets the bits from 'start' to 'end', inclusive.
func (b *BitSet) SetRange(start, end int) {
	validateBitSetIndex(start)
	validateBitSetIndex(end)
	if start > end {
		start, end = end, start
	}
	i1 := start >> addressBitsPerWord
	i2 := end >> addressBitsPerWord
	b.EnsureCapacity(i2 + 1)
	for i := i1; i <= i2; i++ {
		lo := 0
		if i == i1 {
			lo = start & bitIndexMask
		}
		hi := bitIndexMask
		if i == i2 {
			hi = end & bitIndexMask
		}
		mask := (^uint64(0) << uint(lo)) & (^uint64(0) >> uint(bitIndexMask-hi))
		b.set += bits.OnesCount64(mask &^ b.data[i])
		b.data[i] |= mask
	}
}

// Clear the bit at 'index'.
func (b *BitSet) Clear(index int) {
	validateBitSetIndex(index)
	i := index >> addressBitsPerWord
	if i < len(b.data) {
		mask := wordMask(index)
		if b.data[i]&mask == mask {
			b.data[i] &= ^mask
			b.set--
		}
	}
}

// ClearRange clears the bits from 'start' to 'end', inclusive.
func (b *BitSet) ClearRange(start, end int) {
	validateBitSetIndex(start)
	validateBitSetIndex(end)
	if start > end {
		start, end = end, start
	}
	maximum := len(b.data) - 1
	i1 := start >> addressBitsPerWord
	if i1 > maximum {
		return
	}
	// endWord is the word actually containing 'end'. The loop bound i2 is capped to the allocated storage, so when
	// 'end' lies beyond it, i2 < endWord; in that case the capped final word is wholly inside the range and must be
	// cleared in full, i.e. its hi bound is bitIndexMask rather than (end & bitIndexMask). Comparing against endWord
	// (not i2) keeps that distinction.
	endWord := end >> addressBitsPerWord
	i2 := min(endWord, maximum)
	for i := i1; i <= i2; i++ {
		lo := 0
		if i == i1 {
			lo = start & bitIndexMask
		}
		hi := bitIndexMask
		if i == endWord {
			hi = end & bitIndexMask
		}
		mask := (^uint64(0) << uint(lo)) & (^uint64(0) >> uint(bitIndexMask-hi))
		b.set -= bits.OnesCount64(mask & b.data[i])
		b.data[i] &^= mask
	}
}

// Flip the bit at 'index'.
func (b *BitSet) Flip(index int) {
	validateBitSetIndex(index)
	i := index >> addressBitsPerWord
	b.EnsureCapacity(i + 1)
	mask := wordMask(index)
	b.data[i] ^= mask
	if b.data[i]&mask == mask {
		b.set++
	} else {
		b.set--
	}
}

// FlipRange flips the bits from 'start' to 'end', inclusive.
func (b *BitSet) FlipRange(start, end int) {
	validateBitSetIndex(start)
	validateBitSetIndex(end)
	if start > end {
		start, end = end, start
	}
	i1 := start >> addressBitsPerWord
	i2 := end >> addressBitsPerWord
	b.EnsureCapacity(i2 + 1)
	for i := i1; i <= i2; i++ {
		lo := 0
		if i == i1 {
			lo = start & bitIndexMask
		}
		hi := bitIndexMask
		if i == i2 {
			hi = end & bitIndexMask
		}
		mask := (^uint64(0) << uint(lo)) & (^uint64(0) >> uint(bitIndexMask-hi))
		b.set += bits.OnesCount64(mask&^b.data[i]) - bits.OnesCount64(mask&b.data[i])
		b.data[i] ^= mask
	}
}

// And changes this BitSet in place to the bitwise AND of itself and 'other'.
func (b *BitSet) And(other *BitSet) {
	limit := min(len(b.data), len(other.data))
	for i := range limit {
		b.set -= bits.OnesCount64(b.data[i] &^ other.data[i])
		b.data[i] &= other.data[i]
	}
	for i := limit; i < len(b.data); i++ {
		b.set -= bits.OnesCount64(b.data[i])
		b.data[i] = 0
	}
}

// Or changes this BitSet in place to the bitwise OR of itself and 'other'.
func (b *BitSet) Or(other *BitSet) {
	b.EnsureCapacity(len(other.data))
	for i, word := range other.data {
		b.set += bits.OnesCount64(word &^ b.data[i])
		b.data[i] |= word
	}
}

// Xor changes this BitSet in place to the bitwise XOR (exclusive OR) of itself and 'other'.
func (b *BitSet) Xor(other *BitSet) {
	b.EnsureCapacity(len(other.data))
	for i, word := range other.data {
		b.set += bits.OnesCount64(word&^b.data[i]) - bits.OnesCount64(word&b.data[i])
		b.data[i] ^= word
	}
}

// FirstSet returns the first set bit. If no bits are set, then -1 is returned.
func (b *BitSet) FirstSet() int {
	return b.NextSet(0)
}

// LastSet returns the last set bit. If no bits are set, then -1 is returned.
func (b *BitSet) LastSet() int {
	return b.PreviousSet(len(b.data) << addressBitsPerWord)
}

// PreviousSet returns the previous set bit starting from 'start'. If no bits are set at or before 'start', then -1 is
// returned.
func (b *BitSet) PreviousSet(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	firstBit := start & bitIndexMask
	if maximum := len(b.data) - 1; i > maximum {
		i = maximum
		firstBit = bitIndexMask
	}
	mask := ^uint64(0) >> uint(bitIndexMask-firstBit)
	for i >= 0 {
		if word := b.data[i] & mask; word != 0 {
			return i<<addressBitsPerWord + bitIndexMask - bits.LeadingZeros64(word)
		}
		mask = ^uint64(0)
		i--
	}
	return -1
}

// NextSet returns the next set bit starting from 'start'. If no bits are set at or beyond 'start', then -1 is returned.
func (b *BitSet) NextSet(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	mask := ^uint64(0) << uint(start&bitIndexMask)
	maximum := len(b.data)
	for i < maximum {
		if word := b.data[i] & mask; word != 0 {
			return i<<addressBitsPerWord + bits.TrailingZeros64(word)
		}
		mask = ^uint64(0)
		i++
	}
	return -1
}

// PreviousClear returns the previous clear bit starting from 'start'. If no bits are clear at or before 'start', then
// -1 is returned.
func (b *BitSet) PreviousClear(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	if i > len(b.data)-1 {
		return start
	}
	mask := ^uint64(0) >> uint(bitIndexMask-(start&bitIndexMask))
	for i >= 0 {
		if word := ^b.data[i] & mask; word != 0 {
			return i<<addressBitsPerWord + bitIndexMask - bits.LeadingZeros64(word)
		}
		mask = ^uint64(0)
		i--
	}
	return -1
}

// NextClear returns the next clear bit starting from 'start'.
func (b *BitSet) NextClear(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	mask := ^uint64(0) << uint(start&bitIndexMask)
	maximum := len(b.data)
	for i < maximum {
		if word := ^b.data[i] & mask; word != 0 {
			return i<<addressBitsPerWord + bits.TrailingZeros64(word)
		}
		mask = ^uint64(0)
		i++
	}
	return max(maximum*dataBitsPerWord, start)
}

// Trim the BitSet down to the minimum required to store the set bits.
func (b *BitSet) Trim() {
	size := len(b.data)
	for i := size - 1; i >= 0; i-- {
		if b.data[i] != 0 {
			i++
			if i != size {
				data := make([]uint64, i)
				copy(data, b.data)
				b.data = data
			}
			return
		}
	}
	b.data = nil
}

// EnsureCapacity ensures that the BitSet has enough underlying storage to accommodate setting a bit as high as index
// position 'words' x 64 - 1 without needing to allocate more storage.
func (b *BitSet) EnsureCapacity(words int) {
	size := len(b.data)
	if words > size {
		size *= 2
		if size < words {
			size = words
		}
		data := make([]uint64, size)
		copy(data, b.data)
		b.data = data
	}
}

// Data returns a copy of the underlying storage.
func (b *BitSet) Data() []uint64 {
	b.Trim()
	data := make([]uint64, len(b.data))
	copy(data, b.data)
	return data
}

// Load replaces the current data with the bits set in 'data'.
func (b *BitSet) Load(data []uint64) {
	b.data = make([]uint64, len(data))
	copy(b.data, data)
	b.Trim()
	b.set = 0
	for _, word := range b.data {
		b.set += bits.OnesCount64(word)
	}
}

// Reset the BitSet back to an empty state.
func (b *BitSet) Reset() {
	b.data = nil
	b.set = 0
}

func wordMask(index int) uint64 {
	return uint64(1) << uint(index&bitIndexMask)
}

func validateBitSetIndex(index int) {
	if index < 0 {
		xos.ExitWithErr(errs.Newf("Index must be positive (was %d)\n", index))
	}
}
