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

func TestIntnNonPositive(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	c.Equal(0, r.Intn(0))
	c.Equal(0, r.Intn(-1))
	c.Equal(0, r.Intn(-1000))
}

func TestIntnWithinRange(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	// Exercise values that span the different byte-size selections in the implementation (0, 1, 2, 3+ bytes).
	for _, n := range []int{1, 2, 7, 255, 256, 1000, 65535, 65536, 1 << 20, 1 << 30} {
		seen := make(map[int]bool)
		for range 10000 {
			v := r.Intn(n)
			c.True(v >= 0 && v < n, "Intn(%d) returned out-of-range value %d", n, v)
			seen[v] = true
		}
		// Over a small range, this many samples should produce more than one distinct value, guarding against a stuck
		// generator returning a constant.
		if n > 1 && n <= 1000 {
			c.True(len(seen) > 1, "Intn(%d) only ever produced %d distinct value(s)", n, len(seen))
		}
	}
}

// TestIntnMapsBytesIntoRange feeds fixed bytes in place of crypto/rand.Reader and checks the value Intn derives.
func TestIntnMapsBytesIntoRange(t *testing.T) {
	c := check.New(t)
	orig := rand.Reader
	defer func() { rand.Reader = orig }()
	shift := 55 // A variable shift keeps 3<<55 from being a constant that overflows int on 32-bit builds.
	for _, tc := range []struct {
		data   []byte
		n      int
		want   int
		need64 bool
	}{
		// For n == 255 a single byte is read. A bare v % 255 would fold 255 onto 0, over-representing the low values,
		// so it must be rejected in favor of the next byte.
		{data: []byte{255, 5}, n: 255, want: 5},
		// The little-endian bytes of int64's most-negative value, which a signed reconstruction would negate back onto
		// itself and return as negative. 3<<55 is >= 2^56, so all 8 bytes are read, and is not a divisor of 2^63.
		{data: []byte{0, 0, 0, 0, 0, 0, 0, 0x80}, n: 3 << shift, want: 0, need64: true},
	} {
		if tc.need64 && bits.UintSize < 64 {
			continue
		}
		rand.Reader = &cycleReader{data: tc.data}
		c.Equal(tc.want, xrand.New().Intn(tc.n), "Intn(%d) reading %v", tc.n, tc.data)
	}
}
