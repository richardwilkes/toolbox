// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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
	var buffer [8]byte
	size := 8
	n64 := int64(n)
	for i := 1; i < 8; i++ {
		if n64 < int64(1)<<uint(i*8) {
			size = i
			break
		}
	}
	if _, err := rand.Read(buffer[:size]); err != nil {
		return mrnd.IntN(n) //nolint:gosec // Yes, it is ok to use a weak prng here
	}
	var v int
	for i := size - 1; i >= 0; i-- {
		v |= int(buffer[i]) << uint(i*8)
	}
	if v < 0 {
		v = -v
	}
	return v % n
}
