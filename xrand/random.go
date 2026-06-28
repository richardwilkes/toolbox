// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package rand provides a Randomizer based upon the crypto/rand package.
package xrand

import (
	"crypto/rand"
	"math/bits"
	mrnd "math/rand/v2"
)

var cryptoRandInstance = &cryptoRand{}

// Randomizer defines a source of random integer values.
type Randomizer interface {
	// Intn returns a non-negative random number from 0 to n-1. If n <= 0, the implementation should return 0.
	Intn(n int) int
}

// New returns a Randomizer based on the crypto/rand package. This method returns a shared singleton instance
// and does not allocate.
func New() Randomizer {
	return cryptoRandInstance
}

type cryptoRand struct{}

func (r *cryptoRand) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	un := uint64(n)
	// Read only as many bits as are needed to represent [0, n) and reject any sample that lands at or above n. Working
	// in unsigned space avoids the signed-overflow that made -v wrap to a negative result, and rejection sampling keeps
	// the distribution uniform rather than introducing the modulo bias of a bare v % n.
	bitLen := bits.Len64(un - 1)
	byteLen := (bitLen + 7) / 8
	mask := (uint64(1) << uint(bitLen)) - 1
	var buffer [8]byte
	for {
		if _, err := rand.Read(buffer[:byteLen]); err != nil {
			return mrnd.IntN(n) //nolint:gosec // Yes, it is ok to use a weak prng here
		}
		var v uint64
		for i := range byteLen {
			v |= uint64(buffer[i]) << uint(i*8)
		}
		if v &= mask; v < un {
			return int(v) //nolint:gosec // v < un <= math.MaxInt, so the result is non-negative and fits in an int
		}
	}
}
