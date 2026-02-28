// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestStripBOM(t *testing.T) {
	c := check.New(t)

	// Test with UTF-8 BOM present
	bomBytes := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(expected, xbytes.StripBOM(bomBytes))

	// Test with UTF-8 BOM and empty content
	bomOnly := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, xbytes.StripBOM(bomOnly))

	// Test with UTF-8 BOM and single character
	bomSingle := []byte{0xef, 0xbb, 0xbf, 'A'}
	c.Equal([]byte{'A'}, xbytes.StripBOM(bomSingle))

	// Test with UTF-8 BOM and Unicode characters
	bomUnicode := []byte{0xef, 0xbb, 0xbf, 0xc3, 0xa9} // BOM + é
	expectedUnicode := []byte{0xc3, 0xa9}
	c.Equal(expectedUnicode, xbytes.StripBOM(bomUnicode))

	// Test with UTF-8 BOM and emoji
	bomEmoji := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80} // BOM + 😀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80}
	c.Equal(expectedEmoji, xbytes.StripBOM(bomEmoji))

	// Test without BOM - regular text
	noBom := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(noBom, xbytes.StripBOM(noBom))

	// Test without BOM - empty slice
	empty := []byte{}
	c.Equal(empty, xbytes.StripBOM(empty))

	// Test without BOM - single byte
	single := []byte{'A'}
	c.Equal(single, xbytes.StripBOM(single))

	// Test without BOM - two bytes
	twoBytes := []byte{'A', 'B'}
	c.Equal(twoBytes, xbytes.StripBOM(twoBytes))

	// Test with partial BOM - only first byte matches
	partialBom1 := []byte{0xef, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom1, xbytes.StripBOM(partialBom1))

	// Test with partial BOM - first two bytes match
	partialBom2 := []byte{0xef, 0xbb, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom2, xbytes.StripBOM(partialBom2))

	// Test with BOM-like bytes in the middle (should not be stripped)
	bomInMiddle := []byte{'h', 'e', 0xef, 0xbb, 0xbf, 'l', 'o'}
	c.Equal(bomInMiddle, xbytes.StripBOM(bomInMiddle))

	// Test with BOM-like bytes at the end (should not be stripped)
	bomAtEnd := []byte{'h', 'e', 'l', 'l', 'o', 0xef, 0xbb, 0xbf}
	c.Equal(bomAtEnd, xbytes.StripBOM(bomAtEnd))

	// Test with multiple BOM sequences (only first should be stripped)
	multipleBom := []byte{0xef, 0xbb, 0xbf, 0xef, 0xbb, 0xbf, 'h', 'i'}
	expectedMultiple := []byte{0xef, 0xbb, 0xbf, 'h', 'i'}
	c.Equal(expectedMultiple, xbytes.StripBOM(multipleBom))

	// Test with similar but different byte sequences
	notBom1 := []byte{0xef, 0xbb, 0xbe, 'h', 'e', 'l', 'l', 'o'} // Last byte different
	c.Equal(notBom1, xbytes.StripBOM(notBom1))

	notBom2 := []byte{0xef, 0xbc, 0xbf, 'h', 'e', 'l', 'l', 'o'} // Middle byte different
	c.Equal(notBom2, xbytes.StripBOM(notBom2))

	notBom3 := []byte{0xee, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'} // First byte different
	c.Equal(notBom3, xbytes.StripBOM(notBom3))

	// Test with null bytes
	withNull := []byte{0xef, 0xbb, 0xbf, 0x00, 'h', 'e', 'l', 'l', 'o'}
	expectedNull := []byte{0x00, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(expectedNull, xbytes.StripBOM(withNull))

	// Test with binary data after BOM
	binaryData := []byte{0xef, 0xbb, 0xbf, 0x01, 0x02, 0x03, 0xff, 0xfe}
	expectedBinary := []byte{0x01, 0x02, 0x03, 0xff, 0xfe}
	c.Equal(expectedBinary, xbytes.StripBOM(binaryData))

	// Test with JSON after BOM
	jsonWithBom := []byte{0xef, 0xbb, 0xbf, '{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	expectedJSON := []byte{'{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	c.Equal(expectedJSON, xbytes.StripBOM(jsonWithBom))

	// Test with XML after BOM
	xmlWithBom := []byte{0xef, 0xbb, 0xbf, '<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	expectedXML := []byte{'<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	c.Equal(expectedXML, xbytes.StripBOM(xmlWithBom))
}

func TestStripBOMEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test that original slice is not modified when BOM is present
	original := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := xbytes.StripBOM(original)

	// Original should be unchanged
	c.Equal(originalCopy, original)
	// Result should be the content without BOM
	c.Equal([]byte{'h', 'e', 'l', 'l', 'o'}, result)
	// Result should be a different slice (not sharing memory with original)
	c.NotEqual(&original[0], &result[0])

	// Test that original slice is returned when no BOM is present
	noBomOriginal := []byte{'h', 'e', 'l', 'l', 'o'}
	noBomResult := xbytes.StripBOM(noBomOriginal)

	// Should return the same slice reference when no BOM
	c.Equal(&noBomOriginal[0], &noBomResult[0])
	c.Equal(noBomOriginal, noBomResult)

	// Test with exactly 3 bytes that form BOM
	exactBom := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, xbytes.StripBOM(exactBom))

	// Test with less than 3 bytes
	oneByte := []byte{0xef}
	c.Equal(oneByte, xbytes.StripBOM(oneByte))

	twoBytes := []byte{0xef, 0xbb}
	c.Equal(twoBytes, xbytes.StripBOM(twoBytes))

	// Test with nil slice
	var nilSlice []byte
	c.Equal(nilSlice, xbytes.StripBOM(nilSlice))

	// Test with very large content after BOM
	largeBom := make([]byte, 10000)
	largeBom[0] = 0xef
	largeBom[1] = 0xbb
	largeBom[2] = 0xbf
	for i := 3; i < len(largeBom); i++ {
		largeBom[i] = byte(i % 256)
	}

	largeResult := xbytes.StripBOM(largeBom)
	c.Equal(len(largeBom)-3, len(largeResult))
	c.Equal(largeBom[3:], largeResult)
}

func TestStripBOMUnicodeContent(t *testing.T) {
	c := check.New(t)

	// Test with various Unicode characters after BOM

	// Chinese characters
	chineseBom := []byte{0xef, 0xbb, 0xbf, 0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd} // BOM + 中国
	expectedChinese := []byte{0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd}
	c.Equal(expectedChinese, xbytes.StripBOM(chineseBom))

	// Arabic characters
	arabicBom := []byte{0xef, 0xbb, 0xbf, 0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9} // BOM + العربية
	expectedArabic := []byte{0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9}
	c.Equal(expectedArabic, xbytes.StripBOM(arabicBom))

	// Emoji sequence
	emojiBom := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80} // BOM + 😀🎉🚀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80}
	c.Equal(expectedEmoji, xbytes.StripBOM(emojiBom))

	// Mathematical symbols
	mathBom := []byte{0xef, 0xbb, 0xbf, 0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80} // BOM + ∑∞π
	expectedMath := []byte{0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80}
	c.Equal(expectedMath, xbytes.StripBOM(mathBom))

	// Mixed content: ASCII + Unicode
	mixedBom := []byte{0xef, 0xbb, 0xbf, 'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	expectedMixed := []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	c.Equal(expectedMixed, xbytes.StripBOM(mixedBom))
}
