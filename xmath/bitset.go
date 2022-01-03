// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xmath

import (
	"fmt"
	"math"

	"github.com/richardwilkes/toolbox/atexit"
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

func countSetBits(x uint64) int {
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return int(x >> 56)
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
	j := bitIndexForMask(wordMask(start))
	for i := i1; i <= i2; i++ {
		if i != i1 && i != i2 {
			b.set += dataBitsPerWord - countSetBits(b.data[i])
			b.data[i] = math.MaxUint64
		} else {
			var last int
			if i == i2 {
				last = bitIndexForMask(wordMask(end)) + 1
			} else {
				last = dataBitsPerWord
			}
			for j < last {
				mask := wordMask(j)
				if b.data[i]&mask == 0 {
					b.data[i] |= mask
					b.set++
				}
				j++
			}
			j = 0
		}
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
	max := len(b.data) - 1
	i1 := start >> addressBitsPerWord
	if i1 > max {
		return
	}
	i2 := end >> addressBitsPerWord
	if i2 > max {
		i2 = max
	}
	j := bitIndexForMask(wordMask(start))
	for i := i1; i <= i2; i++ {
		if i != i1 && i != i2 {
			b.set -= countSetBits(b.data[i])
			b.data[i] = 0
		} else {
			var last int
			if i == i2 {
				last = bitIndexForMask(wordMask(end)) + 1
			} else {
				last = dataBitsPerWord
			}
			for j < last {
				mask := wordMask(j)
				if b.data[i]&mask == mask {
					b.data[i] &= ^mask
					b.set--
				}
				j++
			}
			j = 0
		}
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
	j := bitIndexForMask(wordMask(start))
	for i := i1; i <= i2; i++ {
		if i != i1 && i != i2 {
			b.set += dataBitsPerWord - 2*countSetBits(b.data[i])
			b.data[i] ^= math.MaxUint64
		} else {
			var last int
			if i == i2 {
				last = bitIndexForMask(wordMask(end)) + 1
			} else {
				last = dataBitsPerWord
			}
			for j < last {
				mask := wordMask(j)
				b.data[i] ^= mask
				if b.data[i]&mask == mask {
					b.set++
				} else {
					b.set--
				}
				j++
			}
			j = 0
		}
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
	var firstBit int
	if max := len(b.data) - 1; i > max {
		i = max
		firstBit = 63
	} else {
		firstBit = bitIndexForMask(wordMask(start))
	}
	for i >= 0 {
		word := b.data[i]
		if word != 0 {
			for j := firstBit; j >= 0; j-- {
				mask := wordMask(j)
				if word&mask == mask {
					return i<<addressBitsPerWord + j
				}
			}
		}
		firstBit = 63
		i--
	}
	return -1
}

// NextSet returns the next set bit starting from 'start'. If no bits are set at or beyond 'start', then -1 is returned.
func (b *BitSet) NextSet(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	firstBit := bitIndexForMask(wordMask(start))
	max := len(b.data)
	for i < max {
		word := b.data[i]
		if word != 0 {
			for j := firstBit; j < dataBitsPerWord; j++ {
				mask := wordMask(j)
				if word&mask == mask {
					return i<<addressBitsPerWord + j
				}
			}
		}
		firstBit = 0
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
	firstBit := bitIndexForMask(wordMask(start))
	for i >= 0 {
		word := b.data[i]
		if word != math.MaxUint64 {
			for j := firstBit; j >= 0; j-- {
				mask := wordMask(j)
				if word&mask == 0 {
					return i<<addressBitsPerWord + j
				}
			}
		}
		firstBit = 63
		i--
	}
	return -1
}

// NextClear returns the next clear bit starting from 'start'.
func (b *BitSet) NextClear(start int) int {
	validateBitSetIndex(start)
	i := start >> addressBitsPerWord
	firstBit := bitIndexForMask(wordMask(start))
	max := len(b.data)
	for i < max {
		word := b.data[i]
		if word != math.MaxUint64 {
			for j := firstBit; j < dataBitsPerWord; j++ {
				mask := wordMask(j)
				if word&mask == 0 {
					return i<<addressBitsPerWord + j
				}
			}
		}
		firstBit = 0
		i++
	}
	return MaxInt(max*dataBitsPerWord, start)
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
		i--
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
	for i := len(b.data) - 1; i >= 0; i-- {
		word := data[i]
		if word != 0 {
			for j := 0; j < dataBitsPerWord; j++ {
				mask := wordMask(j)
				if word&mask == mask {
					b.set++
				}
			}
		}
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

func bitIndexForMask(mask uint64) int {
	for i := 0; i < dataBitsPerWord; i++ {
		if mask == wordMask(i) {
			return i
		}
	}
	fmt.Printf("Unable to determine bit index for mask %064b\n", mask)
	atexit.Exit(1)
	return 0
}

func validateBitSetIndex(index int) {
	if index < 0 {
		fmt.Printf("Index must be positive (was %d)\n", index)
		atexit.Exit(1)
	}
}
