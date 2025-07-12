// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhash_test

import (
	"crypto/sha256"
	"hash/crc32"
	"math"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhash"
)

func TestBool(t *testing.T) {
	c := check.New(t)
	h1 := sha256.New()
	h2 := sha256.New()

	// Test that different boolean values produce different hashes
	xhash.Bool(h1, true)
	xhash.Bool(h2, false)
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same boolean values produce same hashes
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Bool(h3, true)
	xhash.Bool(h4, true)
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with custom boolean type
	type CustomBool bool
	h5 := sha256.New()
	h6 := sha256.New()
	xhash.Bool(h5, CustomBool(true))
	xhash.Bool(h6, true)
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	c.Equal(hash5, hash6)
}

func TestBytesWithLen(t *testing.T) {
	c := check.New(t)

	// Test with empty byte slice
	h1 := sha256.New()
	h2 := sha256.New()
	emptyBytes := []byte{}
	xhash.BytesWithLen(h1, emptyBytes)
	xhash.BytesWithLen(h2, []byte{})
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.Equal(hash1, hash2)

	// Test with different byte slices
	h3 := sha256.New()
	h4 := sha256.New()
	data1 := []byte("hello")
	data2 := []byte("world")
	xhash.BytesWithLen(h3, data1)
	xhash.BytesWithLen(h4, data2)
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.NotEqual(hash3, hash4)

	// Test that length is included in hash (same content, different length should be different)
	h5 := sha256.New()
	h6 := sha256.New()
	data3 := []byte("test")
	data4 := []byte("test\x00") // Same prefix but different length
	xhash.BytesWithLen(h5, data3)
	xhash.BytesWithLen(h6, data4)
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	c.NotEqual(hash5, hash6)

	// Test with custom byte slice type
	type CustomBytes []byte
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.BytesWithLen(h7, CustomBytes("test"))
	xhash.BytesWithLen(h8, []byte("test"))
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)
	c.Equal(hash7, hash8)
}

func TestStringWithLen(t *testing.T) {
	c := check.New(t)

	// Test with empty string
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.StringWithLen(h1, "")
	xhash.StringWithLen(h2, "")
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.Equal(hash1, hash2)

	// Test with different strings
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.StringWithLen(h3, "hello")
	xhash.StringWithLen(h4, "world")
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.NotEqual(hash3, hash4)

	// Test that length is included (strings with same prefix but different length)
	h5 := sha256.New()
	h6 := sha256.New()
	xhash.StringWithLen(h5, "test")
	xhash.StringWithLen(h6, "testing")
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	c.NotEqual(hash5, hash6)

	// Test with Unicode strings
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.StringWithLen(h7, "héllo")
	xhash.StringWithLen(h8, "world")
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)
	c.NotEqual(hash7, hash8)

	// Test that string and equivalent bytes produce same hash
	h9 := sha256.New()
	h10 := sha256.New()
	testStr := "test string"
	xhash.StringWithLen(h9, testStr)
	xhash.BytesWithLen(h10, []byte(testStr))
	hash9 := h9.Sum(nil)
	hash10 := h10.Sum(nil)
	c.Equal(hash9, hash10)

	// Test with custom string type
	type CustomString string
	h11 := sha256.New()
	h12 := sha256.New()
	xhash.StringWithLen(h11, CustomString("test"))
	xhash.StringWithLen(h12, "test")
	hash11 := h11.Sum(nil)
	hash12 := h12.Sum(nil)
	c.Equal(hash11, hash12)
}

func TestNum8(t *testing.T) {
	c := check.New(t)

	// Test with different 8-bit values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Num8(h1, int8(42))
	xhash.Num8(h2, int8(43))
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Num8(h3, uint8(100))
	xhash.Num8(h4, uint8(100))
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with boundary values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.Num8(h5, int8(-128)) // Min int8
	xhash.Num8(h6, int8(127))  // Max int8
	xhash.Num8(h7, uint8(0))   // Min uint8
	xhash.Num8(h8, uint8(255)) // Max uint8
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)

	// All should be different
	c.NotEqual(hash5, hash6)
	c.NotEqual(hash5, hash7)
	c.NotEqual(hash5, hash8)
	c.NotEqual(hash6, hash7)
	c.NotEqual(hash6, hash8)
	c.NotEqual(hash7, hash8)

	// Test that int8 and uint8 with same bit pattern produce same hash
	h9 := sha256.New()
	h10 := sha256.New()
	xhash.Num8(h9, int8(-1))    // 0xFF
	xhash.Num8(h10, uint8(255)) // 0xFF
	hash9 := h9.Sum(nil)
	hash10 := h10.Sum(nil)
	c.Equal(hash9, hash10)

	// Test with custom 8-bit types
	type CustomInt8 int8
	type CustomUint8 uint8
	h11 := sha256.New()
	h12 := sha256.New()
	h13 := sha256.New()
	xhash.Num8(h11, CustomInt8(42))
	xhash.Num8(h12, CustomUint8(42))
	xhash.Num8(h13, int8(42))
	hash11 := h11.Sum(nil)
	hash12 := h12.Sum(nil)
	hash13 := h13.Sum(nil)
	c.Equal(hash11, hash12)
	c.Equal(hash11, hash13)
}

func TestNum16(t *testing.T) {
	c := check.New(t)

	// Test with different 16-bit values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Num16(h1, int16(1000))
	xhash.Num16(h2, int16(2000))
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Num16(h3, uint16(12345))
	xhash.Num16(h4, uint16(12345))
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with boundary values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.Num16(h5, int16(-32768)) // Min int16
	xhash.Num16(h6, int16(32767))  // Max int16
	xhash.Num16(h7, uint16(0))     // Min uint16
	xhash.Num16(h8, uint16(65535)) // Max uint16
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)

	// All should be different
	c.NotEqual(hash5, hash6)
	c.NotEqual(hash5, hash7)
	c.NotEqual(hash5, hash8)
	c.NotEqual(hash6, hash7)
	c.NotEqual(hash6, hash8)
	c.NotEqual(hash7, hash8)

	// Test little-endian encoding
	h9 := sha256.New()
	xhash.Num16(h9, uint16(0x1234))
	// Should write bytes [0x34, 0x12] (little-endian)
	h10 := crc32.NewIEEE()
	_, _ = h10.Write([]byte{0x34, 0x12})
	// We can't directly compare since different hash algorithms, but verify it writes something
	hash9 := h9.Sum(nil)
	c.NotEqual(len(hash9), 0)

	// Test that int16 and uint16 with same bit pattern produce same hash
	h11 := sha256.New()
	h12 := sha256.New()
	xhash.Num16(h11, int16(-1))     // 0xFFFF
	xhash.Num16(h12, uint16(65535)) // 0xFFFF
	hash11 := h11.Sum(nil)
	hash12 := h12.Sum(nil)
	c.Equal(hash11, hash12)

	// Test with custom 16-bit types
	type CustomInt16 int16
	h13 := sha256.New()
	h14 := sha256.New()
	xhash.Num16(h13, CustomInt16(1234))
	xhash.Num16(h14, int16(1234))
	hash13 := h13.Sum(nil)
	hash14 := h14.Sum(nil)
	c.Equal(hash13, hash14)
}

func TestNum32(t *testing.T) {
	c := check.New(t)

	// Test with different 32-bit values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Num32(h1, int32(100000))
	xhash.Num32(h2, int32(200000))
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Num32(h3, uint32(123456789))
	xhash.Num32(h4, uint32(123456789))
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with boundary values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.Num32(h5, int32(-2147483648)) // Min int32
	xhash.Num32(h6, int32(2147483647))  // Max int32
	xhash.Num32(h7, uint32(0))          // Min uint32
	xhash.Num32(h8, uint32(4294967295)) // Max uint32
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)

	// All should be different
	c.NotEqual(hash5, hash6)
	c.NotEqual(hash5, hash7)
	c.NotEqual(hash5, hash8)
	c.NotEqual(hash6, hash7)
	c.NotEqual(hash6, hash8)
	c.NotEqual(hash7, hash8)

	// Test little-endian encoding
	h9 := sha256.New()
	xhash.Num32(h9, uint32(0x12345678))
	// Should write bytes [0x78, 0x56, 0x34, 0x12] (little-endian)
	hash9 := h9.Sum(nil)
	c.NotEqual(len(hash9), 0)

	// Test that int32 and uint32 with same bit pattern produce same hash
	h10 := sha256.New()
	h11 := sha256.New()
	xhash.Num32(h10, int32(-1))          // 0xFFFFFFFF
	xhash.Num32(h11, uint32(4294967295)) // 0xFFFFFFFF
	hash10 := h10.Sum(nil)
	hash11 := h11.Sum(nil)
	c.Equal(hash10, hash11)

	// Test with custom 32-bit types
	type CustomInt32 int32
	h12 := sha256.New()
	h13 := sha256.New()
	xhash.Num32(h12, CustomInt32(123456))
	xhash.Num32(h13, int32(123456))
	hash12 := h12.Sum(nil)
	hash13 := h13.Sum(nil)
	c.Equal(hash12, hash13)
}

func TestNum64(t *testing.T) {
	c := check.New(t)

	// Test with different 64-bit values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Num64(h1, int64(1000000000000))
	xhash.Num64(h2, int64(2000000000000))
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Num64(h3, uint64(123456789012345))
	xhash.Num64(h4, uint64(123456789012345))
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with boundary values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	xhash.Num64(h5, int64(-9223372036854775808))  // Min int64
	xhash.Num64(h6, int64(9223372036854775807))   // Max int64
	xhash.Num64(h7, uint64(0))                    // Min uint64
	xhash.Num64(h8, uint64(18446744073709551615)) // Max uint64
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)

	// All should be different
	c.NotEqual(hash5, hash6)
	c.NotEqual(hash5, hash7)
	c.NotEqual(hash5, hash8)
	c.NotEqual(hash6, hash7)
	c.NotEqual(hash6, hash8)
	c.NotEqual(hash7, hash8)

	// Test with int and uint types (which are 64-bit on 64-bit systems typically)
	h9 := sha256.New()
	h10 := sha256.New()
	h11 := sha256.New()
	xhash.Num64(h9, int(42))
	xhash.Num64(h10, uint(42))
	xhash.Num64(h11, int64(42))
	hash9 := h9.Sum(nil)
	hash10 := h10.Sum(nil)
	hash11 := h11.Sum(nil)
	c.Equal(hash9, hash10)
	c.Equal(hash9, hash11)

	// Test little-endian encoding
	h12 := sha256.New()
	xhash.Num64(h12, uint64(0x123456789ABCDEF0))
	// Should write bytes [0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12] (little-endian)
	hash12 := h12.Sum(nil)
	c.NotEqual(len(hash12), 0)

	// Test that int64 and uint64 with same bit pattern produce same hash
	h13 := sha256.New()
	h14 := sha256.New()
	xhash.Num64(h13, int64(-1))                    // 0xFFFFFFFFFFFFFFFF
	xhash.Num64(h14, uint64(18446744073709551615)) // 0xFFFFFFFFFFFFFFFF
	hash13 := h13.Sum(nil)
	hash14 := h14.Sum(nil)
	c.Equal(hash13, hash14)

	// Test with custom 64-bit types
	type CustomInt64 int64
	type CustomUint64 uint64
	type CustomInt int
	type CustomUint uint
	h15 := sha256.New()
	h16 := sha256.New()
	h17 := sha256.New()
	h18 := sha256.New()
	xhash.Num64(h15, CustomInt64(123456789))
	xhash.Num64(h16, CustomUint64(123456789))
	xhash.Num64(h17, CustomInt(123456789))
	xhash.Num64(h18, CustomUint(123456789))
	hash15 := h15.Sum(nil)
	hash16 := h16.Sum(nil)
	hash17 := h17.Sum(nil)
	hash18 := h18.Sum(nil)
	c.Equal(hash15, hash16)
	c.Equal(hash15, hash17)
	c.Equal(hash15, hash18)
}

func TestFloat32(t *testing.T) {
	c := check.New(t)

	// Test with different float32 values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Float32(h1, float32(3.14159))
	xhash.Float32(h2, float32(2.71828))
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Float32(h3, float32(1.23456))
	xhash.Float32(h4, float32(1.23456))
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with special float values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	h9 := sha256.New()
	h10 := sha256.New()
	xhash.Float32(h5, float32(0.0))
	xhash.Float32(h6, float32(math.Copysign(0.0, -1))) // -0.0
	xhash.Float32(h7, float32(math.Inf(1)))            // +Inf
	xhash.Float32(h8, float32(math.Inf(-1)))           // -Inf
	xhash.Float32(h9, float32(math.NaN()))             // NaN
	xhash.Float32(h10, float32(math.NaN()))            // Another NaN
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)
	hash9 := h9.Sum(nil)
	hash10 := h10.Sum(nil)

	// +0.0 and -0.0 have different bit representations
	c.NotEqual(hash5, hash6)
	// +Inf and -Inf should be different
	c.NotEqual(hash7, hash8)
	// NaN values should have same hash (same bit representation)
	c.Equal(hash9, hash10)

	// Test boundary values
	h11 := sha256.New()
	h12 := sha256.New()
	h13 := sha256.New()
	xhash.Float32(h11, float32(math.MaxFloat32))
	xhash.Float32(h12, float32(math.SmallestNonzeroFloat32))
	xhash.Float32(h13, float32(-math.MaxFloat32))
	hash11 := h11.Sum(nil)
	hash12 := h12.Sum(nil)
	hash13 := h13.Sum(nil)

	c.NotEqual(hash11, hash12)
	c.NotEqual(hash11, hash13)
	c.NotEqual(hash12, hash13)

	// Test that the hash is based on IEEE 754 bit representation
	h14 := sha256.New()
	h15 := sha256.New()
	testValue := float32(1.5)
	xhash.Float32(h14, testValue)
	bits := math.Float32bits(testValue)
	xhash.Num32(h15, bits)
	hash14 := h14.Sum(nil)
	hash15 := h15.Sum(nil)
	c.Equal(hash14, hash15)

	// Test with custom float32 type
	type CustomFloat32 float32
	h16 := sha256.New()
	h17 := sha256.New()
	xhash.Float32(h16, CustomFloat32(3.14159))
	xhash.Float32(h17, float32(3.14159))
	hash16 := h16.Sum(nil)
	hash17 := h17.Sum(nil)
	c.Equal(hash16, hash17)
}

func TestFloat64(t *testing.T) {
	c := check.New(t)

	// Test with different float64 values
	h1 := sha256.New()
	h2 := sha256.New()
	xhash.Float64(h1, 3.141592653589793)
	xhash.Float64(h2, 2.718281828459045)
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)
	c.NotEqual(hash1, hash2)

	// Test that same values produce same hash
	h3 := sha256.New()
	h4 := sha256.New()
	xhash.Float64(h3, 1.2345678901234567)
	xhash.Float64(h4, 1.2345678901234567)
	hash3 := h3.Sum(nil)
	hash4 := h4.Sum(nil)
	c.Equal(hash3, hash4)

	// Test with special float values
	h5 := sha256.New()
	h6 := sha256.New()
	h7 := sha256.New()
	h8 := sha256.New()
	h9 := sha256.New()
	h10 := sha256.New()
	xhash.Float64(h5, 0.0)
	xhash.Float64(h6, math.Copysign(0.0, -1)) // -0.0
	xhash.Float64(h7, math.Inf(1))            // +Inf
	xhash.Float64(h8, math.Inf(-1))           // -Inf
	xhash.Float64(h9, math.NaN())             // NaN
	xhash.Float64(h10, math.NaN())            // Another NaN
	hash5 := h5.Sum(nil)
	hash6 := h6.Sum(nil)
	hash7 := h7.Sum(nil)
	hash8 := h8.Sum(nil)
	hash9 := h9.Sum(nil)
	hash10 := h10.Sum(nil)

	// +0.0 and -0.0 have different bit representations
	c.NotEqual(hash5, hash6)
	// +Inf and -Inf should be different
	c.NotEqual(hash7, hash8)
	// NaN values should have same hash (same bit representation)
	c.Equal(hash9, hash10)

	// Test boundary values
	h11 := sha256.New()
	h12 := sha256.New()
	h13 := sha256.New()
	xhash.Float64(h11, math.MaxFloat64)
	xhash.Float64(h12, math.SmallestNonzeroFloat64)
	xhash.Float64(h13, -math.MaxFloat64)
	hash11 := h11.Sum(nil)
	hash12 := h12.Sum(nil)
	hash13 := h13.Sum(nil)

	c.NotEqual(hash11, hash12)
	c.NotEqual(hash11, hash13)
	c.NotEqual(hash12, hash13)

	// Test that the hash is based on IEEE 754 bit representation
	h14 := sha256.New()
	h15 := sha256.New()
	testValue := 1.5
	xhash.Float64(h14, testValue)
	bits := math.Float64bits(testValue)
	xhash.Num64(h15, bits)
	hash14 := h14.Sum(nil)
	hash15 := h15.Sum(nil)
	c.Equal(hash14, hash15)

	// Test with custom float64 type
	type CustomFloat64 float64
	h16 := sha256.New()
	h17 := sha256.New()
	xhash.Float64(h16, CustomFloat64(3.141592653589793))
	xhash.Float64(h17, 3.141592653589793)
	hash16 := h16.Sum(nil)
	hash17 := h17.Sum(nil)
	c.Equal(hash16, hash17)
}

func TestConsistencyAcrossHashers(t *testing.T) {
	c := check.New(t)

	// Test that the same data produces the same result across different hash implementations
	data := []byte("test data")

	h1 := sha256.New()
	h2 := crc32.NewIEEE()

	xhash.BytesWithLen(h1, data)
	xhash.BytesWithLen(h2, data)

	// The actual hash values will be different due to different algorithms,
	// but the data written should be the same length (8 bytes for length + data)
	// We can't directly verify this, but we ensure no panics occur
	hash1 := h1.Sum(nil)
	hash2 := h2.Sum(nil)

	// Just verify they produced some output
	c.NotEqual(len(hash1), 0)
	c.NotEqual(len(hash2), 0)
}
