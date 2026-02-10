// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xcrc64_test

import (
	"hash/crc64"
	"math"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xcrc64"
)

func TestBool(t *testing.T) {
	c := check.New(t)

	// Test with initial CRC of 0
	crcTrue := xcrc64.Bool(0, true)
	crcFalse := xcrc64.Bool(0, false)

	// Different boolean values should produce different CRCs
	c.NotEqual(crcTrue, crcFalse)

	// Test that the same boolean value produces the same CRC
	c.Equal(crcTrue, xcrc64.Bool(0, true))
	c.Equal(crcFalse, xcrc64.Bool(0, false))

	// Test with non-zero initial CRC
	initialCRC := uint64(0x123456789ABCDEF0)
	crcTrue2 := xcrc64.Bool(initialCRC, true)
	crcFalse2 := xcrc64.Bool(initialCRC, false)

	c.NotEqual(crcTrue2, crcFalse2)
	c.NotEqual(crcTrue, crcTrue2) // Different initial CRC should produce different result
	c.NotEqual(crcFalse, crcFalse2)
}

func TestBytes(t *testing.T) {
	c := check.New(t)

	// Test with empty byte slice
	emptyBytes := []byte{}
	crcEmpty := xcrc64.Bytes(0, emptyBytes)
	c.Equal(uint64(0), crcEmpty) // Empty bytes with CRC 0 should remain 0

	// Test with simple byte data
	data1 := []byte("My")
	data2 := []byte("world")
	data3 := []byte("My")

	crc1 := xcrc64.Bytes(0, data1)
	crc2 := xcrc64.Bytes(0, data2)
	crc3 := xcrc64.Bytes(0, data3)

	// Different data should produce different CRCs
	c.NotEqual(crc1, crc2)
	// Same data should produce same CRC
	c.Equal(crc1, crc3)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x1234567890ABCDEF)
	crcWithInitial := xcrc64.Bytes(initialCRC, data1)
	c.NotEqual(crc1, crcWithInitial)

	// Test with binary data
	binaryData := []byte{0x00, 0xFF, 0x80, 0x7F, 0x01, 0xFE}
	crcBinary := xcrc64.Bytes(0, binaryData)
	c.NotEqual(uint64(0), crcBinary)

	// Verify consistency with standard library
	table := crc64.MakeTable(crc64.ECMA)
	expectedCRC := crc64.Update(0, table, data1)
	c.Equal(expectedCRC, crc1)
}

func TestString(t *testing.T) {
	c := check.New(t)

	// Test with empty string
	crcEmpty := xcrc64.String(0, "")
	c.Equal(uint64(0), crcEmpty)

	// Test with simple strings
	str1 := "Yo"
	str2 := "world"
	str3 := "Yo"

	crc1 := xcrc64.String(0, str1)
	crc2 := xcrc64.String(0, str2)
	crc3 := xcrc64.String(0, str3)

	// Different strings should produce different CRCs
	c.NotEqual(crc1, crc2)
	// Same strings should produce same CRC
	c.Equal(crc1, crc3)

	// Test with Unicode strings
	unicode1 := "héllo"
	unicode2 := "wörld"
	unicode3 := "🚀🌟"

	crcUni1 := xcrc64.String(0, unicode1)
	crcUni2 := xcrc64.String(0, unicode2)
	crcUni3 := xcrc64.String(0, unicode3)

	c.NotEqual(crcUni1, crcUni2)
	c.NotEqual(crcUni1, crcUni3)
	c.NotEqual(crcUni2, crcUni3)

	// Test with non-zero initial CRC
	initialCRC := uint64(0xFEDCBA9876543210)
	crcWithInitial := xcrc64.String(initialCRC, str1)
	c.NotEqual(crc1, crcWithInitial)

	// Verify that String and Bytes produce the same result for equivalent data
	strData := "test string"
	byteData := []byte(strData)
	crcFromString := xcrc64.String(0, strData)
	crcFromBytes := xcrc64.Bytes(0, byteData)
	c.Equal(crcFromString, crcFromBytes)
}

func TestByte(t *testing.T) {
	c := check.New(t)

	// Test with different byte values
	b1 := byte(0x00)
	b2 := byte(0xFF)
	b3 := byte(0x80)
	b4 := byte(0x7F)
	b5 := byte(0x00) // Same as b1

	crc1 := xcrc64.Num8(0, b1)
	crc2 := xcrc64.Num8(0, b2)
	crc3 := xcrc64.Num8(0, b3)
	crc4 := xcrc64.Num8(0, b4)
	crc5 := xcrc64.Num8(0, b5)

	// Different bytes should produce different CRCs
	c.NotEqual(crc1, crc2)
	c.NotEqual(crc1, crc3)
	c.NotEqual(crc1, crc4)
	c.NotEqual(crc2, crc3)
	c.NotEqual(crc2, crc4)
	c.NotEqual(crc3, crc4)

	// Same byte should produce same CRC
	c.Equal(crc1, crc5)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x0123456789ABCDEF)
	crcWithInitial := xcrc64.Num8(initialCRC, b1)
	c.NotEqual(crc1, crcWithInitial)

	// Verify consistency with Bytes function
	byteSlice := []byte{b2}
	crcFromBytes := xcrc64.Bytes(0, byteSlice)
	c.Equal(crc2, crcFromBytes)
}

func TestNumber(t *testing.T) {
	c := check.New(t)

	// Test with different number types and values
	testInt64 := int64(0x123456789ABCDEF0)
	testUint64 := uint64(0x123456789ABCDEF0)
	testInt := int(0x12345678)
	testUint := uint(0x12345678)

	crcInt64 := xcrc64.Num64(0, testInt64)
	crcUint64 := xcrc64.Num64(0, testUint64)
	crcInt := xcrc64.Num64(0, testInt)
	crcUint := xcrc64.Num64(0, testUint)

	// Same numeric value in different types should produce same CRC
	c.Equal(crcInt64, crcUint64)
	c.Equal(crcInt, crcUint)

	// Different values should produce different CRCs
	c.NotEqual(crcInt64, crcInt)

	// Test with zero
	crcZero1 := xcrc64.Num64(0, int64(0))
	crcZero2 := xcrc64.Num64(0, uint64(0))
	c.Equal(crcZero1, crcZero2)

	// Test with negative numbers (for signed types)
	crcPos := xcrc64.Num64(0, int64(42))
	crcNeg := xcrc64.Num64(0, int64(-42))
	c.NotEqual(crcPos, crcNeg)

	// Test with max values
	crcMaxInt64 := xcrc64.Num64(0, int64(9223372036854775807))    // Max int64
	crcMaxUint64 := xcrc64.Num64(0, uint64(18446744073709551615)) // Max uint64
	c.NotEqual(crcMaxInt64, crcMaxUint64)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x5555555555555555)
	crcWithInitial := xcrc64.Num64(initialCRC, testInt64)
	c.NotEqual(crcInt64, crcWithInitial)

	// Test little-endian byte order consistency
	testValue := uint64(0x0123456789ABCDEF)
	expectedBytes := []byte{
		0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01,
	}
	crcFromNumber := xcrc64.Num64(0, testValue)
	crcFromBytes := xcrc64.Bytes(0, expectedBytes)
	c.Equal(crcFromNumber, crcFromBytes)
}

func TestChaining(t *testing.T) {
	c := check.New(t)

	// Test that chaining operations produces consistent results
	data := "consistent data"
	number := uint64(42)
	flag := true
	singleByte := byte(0xFF)

	// Calculate CRC by chaining all operations
	crc := uint64(0)
	crc = xcrc64.String(crc, data)
	crc = xcrc64.Num64(crc, number)
	crc = xcrc64.Bool(crc, flag)
	crc = xcrc64.Num8(crc, singleByte)

	// Calculate the same CRC in a different order
	crc2 := uint64(0)
	crc2 = xcrc64.Num8(crc2, singleByte)
	crc2 = xcrc64.Bool(crc2, flag)
	crc2 = xcrc64.Num64(crc2, number)
	crc2 = xcrc64.String(crc2, data)

	// Different order should produce different CRC (order matters)
	c.NotEqual(crc, crc2)

	// But the same order should produce the same result
	crc3 := uint64(0)
	crc3 = xcrc64.String(crc3, data)
	crc3 = xcrc64.Num64(crc3, number)
	crc3 = xcrc64.Bool(crc3, flag)
	crc3 = xcrc64.Num8(crc3, singleByte)

	c.Equal(crc, crc3)
}

func TestConsistencyWithStandardLibrary(t *testing.T) {
	c := check.New(t)

	// Create the same table used by the package
	table := crc64.MakeTable(crc64.ECMA)

	// Test various data types against standard library
	testData := []byte("test data for consistency check")

	// Using xcrc64
	crcXcrc64 := xcrc64.Bytes(0, testData)

	// Using standard library
	crcStdLib := crc64.Update(0, table, testData)

	// Should be identical
	c.Equal(crcStdLib, crcXcrc64)

	// Test with initial CRC
	initialCRC := uint64(0x123456789ABCDEF0)
	crcXcrc64Initial := xcrc64.Bytes(initialCRC, testData)
	crcStdLibInitial := crc64.Update(initialCRC, table, testData)

	c.Equal(crcStdLibInitial, crcXcrc64Initial)
}

func TestEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test with maximum uint64 as initial CRC
	maxCRC := uint64(0xFFFFFFFFFFFFFFFF)
	crc1 := xcrc64.String(maxCRC, "test")
	crc2 := xcrc64.String(maxCRC, "test")
	c.Equal(crc1, crc2)

	// Test with very long string
	c.NotEqual(uint64(0), xcrc64.String(0, strings.Repeat("ab", 10000)))

	// Test with large byte slice
	largeBytes := make([]byte, 100000)
	for i := range largeBytes {
		largeBytes[i] = byte(i % 256)
	}
	crcLargeBytes := xcrc64.Bytes(0, largeBytes)
	c.NotEqual(uint64(0), crcLargeBytes)

	// Test with special Unicode characters
	specialChars := "™©®℠℗§¶†‡•‰‱⁂⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞"
	crcSpecial := xcrc64.String(0, specialChars)
	c.NotEqual(uint64(0), crcSpecial)
}

func TestBooleanConsistency(t *testing.T) {
	c := check.New(t)

	// Test that Bool function is consistent with Byte function
	// true should be equivalent to byte(1)
	crcBoolTrue := xcrc64.Bool(0, true)
	crcByte1 := xcrc64.Num8(0, byte(1))
	c.Equal(crcBoolTrue, crcByte1)

	// false should be equivalent to byte(0)
	crcBoolFalse := xcrc64.Bool(0, false)
	crcByte0 := xcrc64.Num8(0, byte(0))
	c.Equal(crcBoolFalse, crcByte0)
}

func TestBytesWithLen(t *testing.T) {
	c := check.New(t)

	// Test with empty byte slice
	emptyBytes := []byte{}
	crcEmpty := xcrc64.BytesWithLen(0, emptyBytes)
	// Should include length (0) in the CRC calculation
	expectedCRC := xcrc64.Num64(0, 0) // CRC of length 0
	expectedCRC = xcrc64.Bytes(expectedCRC, emptyBytes)
	c.Equal(expectedCRC, crcEmpty)

	// Test with simple byte data
	data1 := []byte("hello")
	data2 := []byte("world")
	data3 := []byte("hello")

	crc1 := xcrc64.BytesWithLen(0, data1)
	crc2 := xcrc64.BytesWithLen(0, data2)
	crc3 := xcrc64.BytesWithLen(0, data3)

	// Different data should produce different CRCs
	c.NotEqual(crc1, crc2)
	// Same data should produce same CRC
	c.Equal(crc1, crc3)

	// BytesWithLen should differ from Bytes for the same data
	crcBytesOnly := xcrc64.Bytes(0, data1)
	c.NotEqual(crc1, crcBytesOnly)

	// Test that length is actually included
	// Two byte slices with same content but called separately should match BytesWithLen
	manualCRC := xcrc64.Num64(0, len(data1))
	manualCRC = xcrc64.Bytes(manualCRC, data1)
	c.Equal(manualCRC, crc1)

	// Test with different lengths but same initial content
	shortData := []byte("test")
	longData := []byte("test with more content")
	crcShort := xcrc64.BytesWithLen(0, shortData)
	crcLong := xcrc64.BytesWithLen(0, longData)
	c.NotEqual(crcShort, crcLong)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x123456789ABCDEF0)
	crcWithInitial := xcrc64.BytesWithLen(initialCRC, data1)
	c.NotEqual(crc1, crcWithInitial)
}

func TestStringWithLen(t *testing.T) {
	c := check.New(t)

	// Test with empty string
	crcEmpty := xcrc64.StringWithLen(0, "")
	// Should include length (0) in the CRC calculation
	expectedCRC := xcrc64.Num64(0, 0) // CRC of length 0
	expectedCRC = xcrc64.String(expectedCRC, "")
	c.Equal(expectedCRC, crcEmpty)

	// Test with simple strings
	str1 := "what a"
	str2 := "world"
	str3 := "what a"

	crc1 := xcrc64.StringWithLen(0, str1)
	crc2 := xcrc64.StringWithLen(0, str2)
	crc3 := xcrc64.StringWithLen(0, str3)

	// Different strings should produce different CRCs
	c.NotEqual(crc1, crc2)
	// Same strings should produce same CRC
	c.Equal(crc1, crc3)

	// StringWithLen should differ from String for the same data
	crcStringOnly := xcrc64.String(0, str1)
	c.NotEqual(crc1, crcStringOnly)

	// Test that length is actually included
	// Manual calculation should match StringWithLen
	manualCRC := xcrc64.Num64(0, len(str1))
	manualCRC = xcrc64.String(manualCRC, str1)
	c.Equal(manualCRC, crc1)

	// Test with Unicode strings of different byte lengths
	ascii := "hello"   // 5 bytes
	unicode := "héllo" // 6 bytes (é is 2 bytes in UTF-8)
	emoji := "hello🚀"  // 9 bytes (🚀 is 4 bytes in UTF-8)

	createASCII := xcrc64.StringWithLen(0, ascii)
	crcUnicode := xcrc64.StringWithLen(0, unicode)
	crcEmoji := xcrc64.StringWithLen(0, emoji)

	c.NotEqual(createASCII, crcUnicode)
	c.NotEqual(createASCII, crcEmoji)
	c.NotEqual(crcUnicode, crcEmoji)

	// Verify StringWithLen matches BytesWithLen for equivalent data
	testStr := "test string"
	testBytes := []byte(testStr)
	crcFromString := xcrc64.StringWithLen(0, testStr)
	crcFromBytes := xcrc64.BytesWithLen(0, testBytes)
	c.Equal(crcFromString, crcFromBytes)

	// Test with non-zero initial CRC
	initialCRC := uint64(0xFEDCBA9876543210)
	crcWithInitial := xcrc64.StringWithLen(initialCRC, str1)
	c.NotEqual(crc1, crcWithInitial)
}

func TestNum16(t *testing.T) {
	c := check.New(t)

	// Test with different 16-bit values
	val1 := uint16(0x0000)
	val2 := uint16(0xFFFF)
	val3 := uint16(0x8000)
	val4 := uint16(0x7FFF)
	val5 := uint16(0x1234)
	val6 := int16(-1)     // 0xFFFF when cast to uint16
	val7 := int16(-32768) // 0x8000 when cast to uint16

	crc1 := xcrc64.Num16(0, val1)
	crc2 := xcrc64.Num16(0, val2)
	crc3 := xcrc64.Num16(0, val3)
	crc4 := xcrc64.Num16(0, val4)
	crc5 := xcrc64.Num16(0, val5)
	crc6 := xcrc64.Num16(0, val6)
	crc7 := xcrc64.Num16(0, val7)

	// Different values should produce different CRCs
	c.NotEqual(crc1, crc2)
	c.NotEqual(crc1, crc3)
	c.NotEqual(crc1, crc4)
	c.NotEqual(crc1, crc5)
	c.NotEqual(crc2, crc3)
	c.NotEqual(crc2, crc4)
	c.NotEqual(crc2, crc5)

	// Test signed vs unsigned with same bit pattern
	c.Equal(crc2, crc6) // uint16(0xFFFF) == int16(-1)
	c.Equal(crc3, crc7) // uint16(0x8000) == int16(-32768)

	// Test little-endian byte order consistency
	testValue := uint16(0x1234)
	expectedBytes := []byte{0x34, 0x12} // little-endian
	crcFromNumber := xcrc64.Num16(0, testValue)
	crcFromBytes := xcrc64.Bytes(0, expectedBytes)
	c.Equal(crcFromNumber, crcFromBytes)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x0123456789ABCDEF)
	crcWithInitial := xcrc64.Num16(initialCRC, val5)
	c.NotEqual(crc5, crcWithInitial)

	// Test edge values
	crcMin := xcrc64.Num16(0, uint16(0))
	crcMax := xcrc64.Num16(0, uint16(0xFFFF))
	c.NotEqual(crcMin, crcMax)
}

func TestNum32(t *testing.T) {
	c := check.New(t)

	// Test with different 32-bit values
	val1 := uint32(0x00000000)
	val2 := uint32(0xFFFFFFFF)
	val3 := uint32(0x80000000)
	val4 := uint32(0x7FFFFFFF)
	val5 := uint32(0x12345678)
	val6 := int32(-1)          // 0xFFFFFFFF when cast to uint32
	val7 := int32(-2147483648) // 0x80000000 when cast to uint32

	crc1 := xcrc64.Num32(0, val1)
	crc2 := xcrc64.Num32(0, val2)
	crc3 := xcrc64.Num32(0, val3)
	crc4 := xcrc64.Num32(0, val4)
	crc5 := xcrc64.Num32(0, val5)
	crc6 := xcrc64.Num32(0, val6)
	crc7 := xcrc64.Num32(0, val7)

	// Different values should produce different CRCs
	c.NotEqual(crc1, crc2)
	c.NotEqual(crc1, crc3)
	c.NotEqual(crc1, crc4)
	c.NotEqual(crc1, crc5)
	c.NotEqual(crc2, crc3)
	c.NotEqual(crc2, crc4)
	c.NotEqual(crc2, crc5)

	// Test signed vs unsigned with same bit pattern
	c.Equal(crc2, crc6) // uint32(0xFFFFFFFF) == int32(-1)
	c.Equal(crc3, crc7) // uint32(0x80000000) == int32(-2147483648)

	// Test little-endian byte order consistency
	testValue := uint32(0x12345678)
	expectedBytes := []byte{0x78, 0x56, 0x34, 0x12} // little-endian
	crcFromNumber := xcrc64.Num32(0, testValue)
	crcFromBytes := xcrc64.Bytes(0, expectedBytes)
	c.Equal(crcFromNumber, crcFromBytes)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x5555555555555555)
	crcWithInitial := xcrc64.Num32(initialCRC, val5)
	c.NotEqual(crc5, crcWithInitial)

	// Test edge values
	crcMin := xcrc64.Num32(0, uint32(0))
	crcMax := xcrc64.Num32(0, uint32(0xFFFFFFFF))
	c.NotEqual(crcMin, crcMax)
}

func TestFloat32(t *testing.T) {
	c := check.New(t)

	// Test with different float32 values
	val1 := float32(0.0)
	val2 := float32(1.0)
	val3 := float32(-1.0)
	val4 := float32(3.14159)
	val5 := float32(-3.14159)
	val6 := float32(1.23456789e10)
	val7 := float32(1.23456789e-10)

	crc1 := xcrc64.Float32(0, val1)
	crc2 := xcrc64.Float32(0, val2)
	crc3 := xcrc64.Float32(0, val3)
	crc4 := xcrc64.Float32(0, val4)
	crc5 := xcrc64.Float32(0, val5)
	crc6 := xcrc64.Float32(0, val6)
	crc7 := xcrc64.Float32(0, val7)

	// Different values should produce different CRCs
	c.NotEqual(crc1, crc2)
	c.NotEqual(crc1, crc3)
	c.NotEqual(crc2, crc3)
	c.NotEqual(crc4, crc5)
	c.NotEqual(crc4, crc6)
	c.NotEqual(crc4, crc7)

	// Test special float values
	posInf := float32(math.Inf(1))
	negInf := float32(math.Inf(-1))
	nan := float32(math.NaN())

	crcPosInf := xcrc64.Float32(0, posInf)
	crcNegInf := xcrc64.Float32(0, negInf)
	crcNaN := xcrc64.Float32(0, nan)

	c.NotEqual(crcPosInf, crcNegInf)
	c.NotEqual(crcPosInf, crcNaN)
	c.NotEqual(crcNegInf, crcNaN)

	// Test positive and negative zero
	posZero := float32(0.0)
	negZero := float32(math.Copysign(0.0, -1.0))
	crcPosZero := xcrc64.Float32(0, posZero)
	crcNegZero := xcrc64.Float32(0, negZero)
	// In IEEE 754, +0.0 and -0.0 have different bit representations
	c.NotEqual(crcPosZero, crcNegZero)

	// Test consistency with Num32 using Float32bits
	testFloat := float32(42.5)
	crcFloat := xcrc64.Float32(0, testFloat)
	crcBits := xcrc64.Num32(0, math.Float32bits(testFloat))
	c.Equal(crcFloat, crcBits)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x1111111111111111)
	crcWithInitial := xcrc64.Float32(initialCRC, val4)
	c.NotEqual(crc4, crcWithInitial)

	// Test that same float value produces same CRC
	val8 := float32(3.14159)
	crc8 := xcrc64.Float32(0, val8)
	c.Equal(crc4, crc8)
}

func TestFloat64(t *testing.T) {
	c := check.New(t)

	// Test with different float64 values
	val1 := float64(0.0)
	val2 := float64(1.0)
	val3 := float64(-1.0)
	val4 := 3.14159265358979323846
	val5 := -3.14159265358979323846
	val6 := 1.23456789e100
	val7 := 1.23456789e-100

	crc1 := xcrc64.Float64(0, val1)
	crc2 := xcrc64.Float64(0, val2)
	crc3 := xcrc64.Float64(0, val3)
	crc4 := xcrc64.Float64(0, val4)
	crc5 := xcrc64.Float64(0, val5)
	crc6 := xcrc64.Float64(0, val6)
	crc7 := xcrc64.Float64(0, val7)

	// Different values should produce different CRCs
	c.NotEqual(crc1, crc2)
	c.NotEqual(crc1, crc3)
	c.NotEqual(crc2, crc3)
	c.NotEqual(crc4, crc5)
	c.NotEqual(crc4, crc6)
	c.NotEqual(crc4, crc7)

	// Test special float values
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	nan := math.NaN()

	crcPosInf := xcrc64.Float64(0, posInf)
	crcNegInf := xcrc64.Float64(0, negInf)
	crcNaN := xcrc64.Float64(0, nan)

	c.NotEqual(crcPosInf, crcNegInf)
	c.NotEqual(crcPosInf, crcNaN)
	c.NotEqual(crcNegInf, crcNaN)

	// Test positive and negative zero
	posZero := 0.0
	negZero := math.Copysign(0.0, -1.0)
	crcPosZero := xcrc64.Float64(0, posZero)
	crcNegZero := xcrc64.Float64(0, negZero)
	// In IEEE 754, +0.0 and -0.0 have different bit representations
	c.NotEqual(crcPosZero, crcNegZero)

	// Test consistency with Num64 using Float64bits
	testFloat := 42.123456789
	crcFloat := xcrc64.Float64(0, testFloat)
	crcBits := xcrc64.Num64(0, math.Float64bits(testFloat))
	c.Equal(crcFloat, crcBits)

	// Test with non-zero initial CRC
	initialCRC := uint64(0x9999999999999999)
	crcWithInitial := xcrc64.Float64(initialCRC, val4)
	c.NotEqual(crc4, crcWithInitial)

	// Test that same float value produces same CRC
	val8 := 3.14159265358979323846
	crc8 := xcrc64.Float64(0, val8)
	c.Equal(crc4, crc8)

	// Test precision differences
	lowPrecision := float64(float32(3.14159265358979323846)) // Converted through float32
	highPrecision := 3.14159265358979323846
	crcLowPrec := xcrc64.Float64(0, lowPrecision)
	crcHighPrec := xcrc64.Float64(0, highPrecision)
	// Due to precision loss, these should be different
	c.NotEqual(crcLowPrec, crcHighPrec)
}

func TestNumericTypesConsistency(t *testing.T) {
	c := check.New(t)

	// Test that different numeric functions produce consistent results for same bit patterns
	value8 := uint8(0xFF)

	num8 := xcrc64.Num8(0, value8)
	num16 := xcrc64.Num16(0, uint16(value8))
	num32 := xcrc64.Num32(0, uint32(value8))
	num64 := xcrc64.Num64(0, uint64(value8))

	// These should all be different because they produce different byte sequences
	c.NotEqual(num8, num16)
	c.NotEqual(num8, num32)
	c.NotEqual(num8, num64)
	c.NotEqual(num16, num32)
	c.NotEqual(num16, num64)
	c.NotEqual(num32, num64)

	// But Num16 with 0x00FF should be the same as byte sequence [0xFF, 0x00]
	expectedBytes16 := []byte{0xFF, 0x00}
	crcBytes16 := xcrc64.Bytes(0, expectedBytes16)
	c.Equal(num16, crcBytes16)
}
