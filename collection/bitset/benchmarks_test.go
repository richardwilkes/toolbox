// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package bitset

import "testing"

const (
	benchWords = 1024
	benchBits  = benchWords * dataBitsPerWord
	benchStep  = 129
)

var (
	benchIntSink    int
	benchBitSetSink *BitSet
)

func newSparseBenchSet() *BitSet {
	var bs BitSet
	for i := 0; i < benchBits; i += benchStep {
		bs.Set(i)
	}
	return &bs
}

func newDenseBenchSet() *BitSet {
	var bs BitSet
	bs.SetRange(0, benchBits-1)
	for i := 0; i < benchBits; i += benchStep {
		bs.Clear(i)
	}
	return &bs
}

func BenchmarkLoad(b *testing.B) {
	data := make([]uint64, benchWords)
	for i := range data {
		data[i] = 0xaaaa5555aaaa5555
	}
	var bs BitSet
	for b.Loop() {
		bs.Load(data)
	}
	benchBitSetSink = &bs
}

func BenchmarkNextSet(b *testing.B) {
	bs := newSparseBenchSet()
	var count int
	for b.Loop() {
		count = 0
		for i := bs.NextSet(0); i != -1; i = bs.NextSet(i + 1) {
			count++
		}
	}
	benchIntSink = count
}

func BenchmarkPreviousSet(b *testing.B) {
	bs := newSparseBenchSet()
	last := bs.LastSet()
	var count int
	for b.Loop() {
		count = 0
		for i := last; i > 0; i = bs.PreviousSet(i - 1) {
			count++
		}
	}
	benchIntSink = count
}

func BenchmarkNextClear(b *testing.B) {
	bs := newDenseBenchSet()
	var count int
	for b.Loop() {
		count = 0
		for i := bs.NextClear(0); i < benchBits; i = bs.NextClear(i + 1) {
			count++
		}
	}
	benchIntSink = count
}

func BenchmarkPreviousClear(b *testing.B) {
	bs := newDenseBenchSet()
	last := bs.PreviousClear(benchBits - 1)
	var count int
	for b.Loop() {
		count = 0
		for i := last; i > 0; i = bs.PreviousClear(i - 1) {
			count++
		}
	}
	benchIntSink = count
}

func BenchmarkSetRange(b *testing.B) {
	var bs BitSet
	bs.EnsureCapacity(benchWords)
	for b.Loop() {
		bs.SetRange(3, benchBits-3)
	}
	benchBitSetSink = &bs
}

func BenchmarkClearRange(b *testing.B) {
	var bs BitSet
	bs.EnsureCapacity(benchWords)
	for b.Loop() {
		bs.ClearRange(3, benchBits-3)
	}
	benchBitSetSink = &bs
}

func BenchmarkFlipRange(b *testing.B) {
	var bs BitSet
	bs.EnsureCapacity(benchWords)
	for b.Loop() {
		bs.FlipRange(3, benchBits-3)
		bs.FlipRange(3, benchBits-3)
	}
	benchBitSetSink = &bs
}

func BenchmarkAnd(b *testing.B) {
	bs := newDenseBenchSet()
	other := newSparseBenchSet()
	// And is idempotent, so the receiver reaches its final state on the first call and every later call measures the
	// steady-state word loop over the full benchWords of storage.
	for b.Loop() {
		bs.And(other)
	}
	benchBitSetSink = bs
}

func BenchmarkOr(b *testing.B) {
	bs := newDenseBenchSet()
	other := newSparseBenchSet()
	// Or is idempotent, so the receiver reaches its final state on the first call and every later call measures the
	// steady-state word loop over the full benchWords of storage.
	for b.Loop() {
		bs.Or(other)
	}
	benchBitSetSink = bs
}

func BenchmarkXor(b *testing.B) {
	bs := newDenseBenchSet()
	other := newSparseBenchSet()
	// Xor is not idempotent, so it is called twice per iteration to restore the receiver's original state, which keeps
	// the starting point of each iteration identical.
	for b.Loop() {
		bs.Xor(other)
		bs.Xor(other)
	}
	benchBitSetSink = bs
}
