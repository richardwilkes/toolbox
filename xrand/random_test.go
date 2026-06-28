// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xrand_test

import (
	"crypto/rand"
	"math/bits"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xrand"
)

// cycleReader yields a fixed byte sequence, repeating it indefinitely, so tests can feed deterministic "random" bytes
// to the Randomizer by overriding crypto/rand.Reader.
type cycleReader struct {
	data []byte
	pos  int
}

func (r *cycleReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = r.data[r.pos%len(r.data)]
		r.pos++
	}
	return len(p), nil
}

func TestNewReturnsUsableRandomizer(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	c.NotNil(r)
}

func TestIntnNonPositive(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	c.Equal(0, r.Intn(0))
	c.Equal(0, r.Intn(-1))
	c.Equal(0, r.Intn(-1000))
}

func TestIntnOne(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	for range 1000 {
		c.Equal(0, r.Intn(1))
	}
}

func TestIntnWithinRange(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	// Exercise values that span the different byte-size selections in the implementation (1, 2, 3+ bytes).
	for _, n := range []int{2, 7, 255, 256, 1000, 65535, 65536, 1 << 20, 1 << 30} {
		seen := make(map[int]bool)
		for range 10000 {
			v := r.Intn(n)
			c.True(v >= 0 && v < n, "Intn(%d) returned out-of-range value %d", n, v)
			seen[v] = true
		}
		// With this many samples over a small range, we expect more than one distinct value, which guards against a
		// stuck generator returning a constant.
		if n <= 1000 {
			c.True(len(seen) > 1, "Intn(%d) only ever produced %d distinct value(s)", n, len(seen))
		}
	}
}

// TestIntnNeverNegativeWithAdversarialBytes feeds the byte pattern for int64's most-negative value (0x80 in the high
// byte). The old implementation reconstructed this as a signed int, negated it (which overflows back to the same
// negative value), and returned a negative result for n >= 2^56; the unsigned implementation must stay in range.
func TestIntnNeverNegativeWithAdversarialBytes(t *testing.T) {
	c := check.New(t)
	if bits.UintSize < 64 {
		t.Skip("the negative-overflow case only arises for n >= 2^56, which requires a 64-bit int")
	}
	orig := rand.Reader
	defer func() { rand.Reader = orig }()
	// Little-endian bytes for 0x8000000000000000: only the most-significant byte's high bit is set.
	rand.Reader = &cycleReader{data: []byte{0, 0, 0, 0, 0, 0, 0, 0x80}}
	shift := 55     // A variable shift keeps this from being a constant that overflows int on 32-bit builds.
	n := 3 << shift // 3 * 2^55: >= 2^56 (so a full 8 bytes are consumed) and not a divisor of 2^63.
	got := xrand.New().Intn(n)
	c.True(got >= 0 && got < n, "Intn(%d) returned out-of-range value %d", n, got)
}

// TestIntnRejectionSamplingRemovesBias verifies that out-of-range samples are rejected rather than folded back with a
// bare modulo, which would over-represent the low values. For n == 255 a single byte is read; the value 255 is the one
// byte that v % 255 would map onto 0, so it must be discarded in favor of the next byte.
func TestIntnRejectionSamplingRemovesBias(t *testing.T) {
	c := check.New(t)
	orig := rand.Reader
	defer func() { rand.Reader = orig }()
	rand.Reader = &cycleReader{data: []byte{255, 5}}
	c.Equal(5, xrand.New().Intn(255))
}
