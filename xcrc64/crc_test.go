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
	data1 := []byte("hello")
	data2 := []byte("world")
	data3 := []byte("hello")

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
	str1 := "hello"
	str2 := "world"
	str3 := "hello"

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
	longString := ""
	for range 10000 {
		longString += "a"
	}
	crcLong := xcrc64.String(0, longString)
	c.NotEqual(uint64(0), crcLong)

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
