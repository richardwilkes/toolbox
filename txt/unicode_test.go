package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestStripBOM(t *testing.T) {
	c := check.New(t)

	// Test with UTF-8 BOM present
	bomBytes := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(expected, StripBOM(bomBytes))

	// Test with UTF-8 BOM and empty content
	bomOnly := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, StripBOM(bomOnly))

	// Test with UTF-8 BOM and single character
	bomSingle := []byte{0xef, 0xbb, 0xbf, 'A'}
	c.Equal([]byte{'A'}, StripBOM(bomSingle))

	// Test with UTF-8 BOM and Unicode characters
	bomUnicode := []byte{0xef, 0xbb, 0xbf, 0xc3, 0xa9} // BOM + é
	expectedUnicode := []byte{0xc3, 0xa9}
	c.Equal(expectedUnicode, StripBOM(bomUnicode))

	// Test with UTF-8 BOM and emoji
	bomEmoji := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80} // BOM + 😀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80}
	c.Equal(expectedEmoji, StripBOM(bomEmoji))

	// Test without BOM - regular text
	noBom := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(noBom, StripBOM(noBom))

	// Test without BOM - empty slice
	empty := []byte{}
	c.Equal(empty, StripBOM(empty))

	// Test without BOM - single byte
	single := []byte{'A'}
	c.Equal(single, StripBOM(single))

	// Test without BOM - two bytes
	twoBytes := []byte{'A', 'B'}
	c.Equal(twoBytes, StripBOM(twoBytes))

	// Test with partial BOM - only first byte matches
	partialBom1 := []byte{0xef, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom1, StripBOM(partialBom1))

	// Test with partial BOM - first two bytes match
	partialBom2 := []byte{0xef, 0xbb, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom2, StripBOM(partialBom2))

	// Test with BOM-like bytes in the middle (should not be stripped)
	bomInMiddle := []byte{'h', 'e', 0xef, 0xbb, 0xbf, 'l', 'o'}
	c.Equal(bomInMiddle, StripBOM(bomInMiddle))

	// Test with BOM-like bytes at the end (should not be stripped)
	bomAtEnd := []byte{'h', 'e', 'l', 'l', 'o', 0xef, 0xbb, 0xbf}
	c.Equal(bomAtEnd, StripBOM(bomAtEnd))

	// Test with multiple BOM sequences (only first should be stripped)
	multipleBom := []byte{0xef, 0xbb, 0xbf, 0xef, 0xbb, 0xbf, 'h', 'i'}
	expectedMultiple := []byte{0xef, 0xbb, 0xbf, 'h', 'i'}
	c.Equal(expectedMultiple, StripBOM(multipleBom))

	// Test with similar but different byte sequences
	notBom1 := []byte{0xef, 0xbb, 0xbe, 'h', 'e', 'l', 'l', 'o'} // Last byte different
	c.Equal(notBom1, StripBOM(notBom1))

	notBom2 := []byte{0xef, 0xbc, 0xbf, 'h', 'e', 'l', 'l', 'o'} // Middle byte different
	c.Equal(notBom2, StripBOM(notBom2))

	notBom3 := []byte{0xee, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'} // First byte different
	c.Equal(notBom3, StripBOM(notBom3))

	// Test with null bytes
	withNull := []byte{0xef, 0xbb, 0xbf, 0x00, 'h', 'e', 'l', 'l', 'o'}
	expectedNull := []byte{0x00, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(expectedNull, StripBOM(withNull))

	// Test with binary data after BOM
	binaryData := []byte{0xef, 0xbb, 0xbf, 0x01, 0x02, 0x03, 0xff, 0xfe}
	expectedBinary := []byte{0x01, 0x02, 0x03, 0xff, 0xfe}
	c.Equal(expectedBinary, StripBOM(binaryData))

	// Test with JSON after BOM
	jsonWithBom := []byte{0xef, 0xbb, 0xbf, '{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	expectedJSON := []byte{'{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	c.Equal(expectedJSON, StripBOM(jsonWithBom))

	// Test with XML after BOM
	xmlWithBom := []byte{0xef, 0xbb, 0xbf, '<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	expectedXML := []byte{'<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	c.Equal(expectedXML, StripBOM(xmlWithBom))
}

func TestStripBOMEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test that original slice is not modified when BOM is present
	original := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := StripBOM(original)

	// Original should be unchanged
	c.Equal(originalCopy, original)
	// Result should be the content without BOM
	c.Equal([]byte{'h', 'e', 'l', 'l', 'o'}, result)
	// Result should be a different slice (not sharing memory with original)
	c.NotEqual(&original[0], &result[0])

	// Test that original slice is returned when no BOM is present
	noBomOriginal := []byte{'h', 'e', 'l', 'l', 'o'}
	noBomResult := StripBOM(noBomOriginal)

	// Should return the same slice reference when no BOM
	c.Equal(&noBomOriginal[0], &noBomResult[0])
	c.Equal(noBomOriginal, noBomResult)

	// Test with exactly 3 bytes that form BOM
	exactBom := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, StripBOM(exactBom))

	// Test with less than 3 bytes
	oneByte := []byte{0xef}
	c.Equal(oneByte, StripBOM(oneByte))

	twoBytes := []byte{0xef, 0xbb}
	c.Equal(twoBytes, StripBOM(twoBytes))

	// Test with nil slice
	var nilSlice []byte
	c.Equal(nilSlice, StripBOM(nilSlice))

	// Test with very large content after BOM
	largeBom := make([]byte, 10000)
	largeBom[0] = 0xef
	largeBom[1] = 0xbb
	largeBom[2] = 0xbf
	for i := 3; i < len(largeBom); i++ {
		largeBom[i] = byte(i % 256)
	}

	largeResult := StripBOM(largeBom)
	c.Equal(len(largeBom)-3, len(largeResult))
	c.Equal(largeBom[3:], largeResult)
}

func TestStripBOMUnicodeContent(t *testing.T) {
	c := check.New(t)

	// Test with various Unicode characters after BOM

	// Chinese characters
	chineseBom := []byte{0xef, 0xbb, 0xbf, 0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd} // BOM + 中国
	expectedChinese := []byte{0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd}
	c.Equal(expectedChinese, StripBOM(chineseBom))

	// Arabic characters
	arabicBom := []byte{0xef, 0xbb, 0xbf, 0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9} // BOM + العربية
	expectedArabic := []byte{0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9}
	c.Equal(expectedArabic, StripBOM(arabicBom))

	// Emoji sequence
	emojiBom := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80} // BOM + 😀🎉🚀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80}
	c.Equal(expectedEmoji, StripBOM(emojiBom))

	// Mathematical symbols
	mathBom := []byte{0xef, 0xbb, 0xbf, 0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80} // BOM + ∑∞π
	expectedMath := []byte{0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80}
	c.Equal(expectedMath, StripBOM(mathBom))

	// Mixed content: ASCII + Unicode
	mixedBom := []byte{0xef, 0xbb, 0xbf, 'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	expectedMixed := []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	c.Equal(expectedMixed, StripBOM(mixedBom))
}
