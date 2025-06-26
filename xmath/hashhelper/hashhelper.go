// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package hashhelper

import (
	"hash"
	"math"
)

// String writes the given string to the hash.
func String[T ~string](h hash.Hash, data T) {
	Num64(h, len(data))
	_, _ = h.Write([]byte(data))
}

// Num64 writes the given 64-bit number to the hash.
func Num64[T ~int64 | ~uint64 | ~int | ~uint](h hash.Hash, data T) {
	var buffer [8]byte
	d := uint64(data)
	buffer[0] = byte(d)
	buffer[1] = byte(d >> 8)
	buffer[2] = byte(d >> 16)
	buffer[3] = byte(d >> 24)
	buffer[4] = byte(d >> 32)
	buffer[5] = byte(d >> 40)
	buffer[6] = byte(d >> 48)
	buffer[7] = byte(d >> 56)
	_, _ = h.Write(buffer[:])
}

// Num32 writes the given 32-bit number to the hash.
func Num32[T ~int32 | ~uint32](h hash.Hash, data T) {
	var buffer [4]byte
	d := uint32(data)
	buffer[0] = byte(d)
	buffer[1] = byte(d >> 8)
	buffer[2] = byte(d >> 16)
	buffer[3] = byte(d >> 24)
	_, _ = h.Write(buffer[:])
}

// Num16 writes the given 16-bit number to the hash.
func Num16[T ~int16 | ~uint16](h hash.Hash, data T) {
	var buffer [2]byte
	d := uint16(data)
	buffer[0] = byte(d)
	buffer[1] = byte(d >> 8)
	_, _ = h.Write(buffer[:])
}

// Num8 writes the given 8-bit number to the hash.
func Num8[T ~int8 | ~uint8](h hash.Hash, data T) {
	_, _ = h.Write([]byte{byte(data)})
}

// Bool writes the given boolean to the hash.
func Bool[T ~bool](h hash.Hash, data T) {
	var b byte
	if data {
		b = 1
	}
	_, _ = h.Write([]byte{b})
}

// Float64 writes the given 64-bit float to the hash.
func Float64[T ~float64](h hash.Hash, data T) {
	Num64(h, math.Float64bits(float64(data)))
}

// Float32 writes the given 64-bit float to the hash.
func Float32[T ~float32](h hash.Hash, data T) {
	Num32(h, math.Float32bits(float32(data)))
}
