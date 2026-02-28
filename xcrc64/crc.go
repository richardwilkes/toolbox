// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xcrc64

import (
	"hash/crc64"
	"math"
)

var crcTable = crc64.MakeTable(crc64.ECMA)

// Bool returns the CRC-64 value for the given data starting with the given crc value.
func Bool(crc uint64, data bool) uint64 {
	var buffer [1]byte
	if data {
		buffer[0] = 1
	}
	return crc64.Update(crc, crcTable, buffer[:])
}

// Bytes returns the CRC-64 value for the given data starting with the given crc value.
func Bytes(crc uint64, data []byte) uint64 {
	return crc64.Update(crc, crcTable, data)
}

// BytesWithLen returns the CRC-64 value for the given data's length + data starting with the given crc value.
func BytesWithLen(crc uint64, data []byte) uint64 {
	return crc64.Update(Num64(crc, len(data)), crcTable, data)
}

// String returns the CRC-64 value for the given data starting with the given crc value.
func String(crc uint64, data string) uint64 {
	return crc64.Update(crc, crcTable, []byte(data))
}

// StringWithLen returns the CRC-64 value for the given data's length + data starting with the given crc value.
func StringWithLen(crc uint64, data string) uint64 {
	return crc64.Update(Num64(crc, len(data)), crcTable, []byte(data))
}

// Num8 returns the CRC-64 value for the given data starting with the given crc value.
func Num8[T ~int8 | ~uint8](crc uint64, data T) uint64 {
	var buffer [1]byte
	buffer[0] = byte(data)
	return crc64.Update(crc, crcTable, buffer[:])
}

// Num16 returns the CRC-64 value for the given data starting with the given crc value.
func Num16[T ~int16 | ~uint16](crc uint64, data T) uint64 {
	var buffer [2]byte
	d := uint16(data)
	buffer[0] = byte(d)
	buffer[1] = byte(d >> 8)
	return crc64.Update(crc, crcTable, buffer[:])
}

// Num32 returns the CRC-64 value for the given data starting with the given crc value.
func Num32[T ~int32 | ~uint32](crc uint64, data T) uint64 {
	var buffer [4]byte
	d := uint32(data)
	buffer[0] = byte(d)
	buffer[1] = byte(d >> 8)
	buffer[2] = byte(d >> 16)
	buffer[3] = byte(d >> 24)
	return crc64.Update(crc, crcTable, buffer[:])
}

// Num64 returns the CRC-64 value for the given data starting with the given crc value.
func Num64[T ~int64 | ~uint64 | ~int | ~uint](crc uint64, data T) uint64 {
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
	return crc64.Update(crc, crcTable, buffer[:])
}

// Float32 writes the given 64-bit float to the hash.
func Float32[T ~float32](crc uint64, data T) uint64 {
	return Num32(crc, math.Float32bits(float32(data)))
}

// Float64 writes the given 64-bit float to the hash.
func Float64[T ~float64](crc uint64, data T) uint64 {
	return Num64(crc, math.Float64bits(float64(data)))
}
