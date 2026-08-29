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

	bomBytes := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(expected, xbytes.StripBOM(bomBytes))

	bomOnly := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, xbytes.StripBOM(bomOnly))

	bomSingle := []byte{0xef, 0xbb, 0xbf, 'A'}
	c.Equal([]byte{'A'}, xbytes.StripBOM(bomSingle))

	bomUnicode := []byte{0xef, 0xbb, 0xbf, 0xc3, 0xa9} // BOM + é
	expectedUnicode := []byte{0xc3, 0xa9}
	c.Equal(expectedUnicode, xbytes.StripBOM(bomUnicode))

	bomEmoji := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80} // BOM + 😀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80}
	c.Equal(expectedEmoji, xbytes.StripBOM(bomEmoji))

	noBom := []byte{'h', 'e', 'l', 'l', 'o'}
	c.Equal(noBom, xbytes.StripBOM(noBom))

	empty := []byte{}
	c.Equal(empty, xbytes.StripBOM(empty))

	single := []byte{'A'}
	c.Equal(single, xbytes.StripBOM(single))

	twoBytes := []byte{'A', 'B'}
	c.Equal(twoBytes, xbytes.StripBOM(twoBytes))

	partialBom1 := []byte{0xef, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom1, xbytes.StripBOM(partialBom1))

	partialBom2 := []byte{0xef, 0xbb, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(partialBom2, xbytes.StripBOM(partialBom2))

	bomInMiddle := []byte{'h', 'e', 0xef, 0xbb, 0xbf, 'l', 'o'}
	c.Equal(bomInMiddle, xbytes.StripBOM(bomInMiddle))

	bomAtEnd := []byte{'h', 'e', 'l', 'l', 'o', 0xef, 0xbb, 0xbf}
	c.Equal(bomAtEnd, xbytes.StripBOM(bomAtEnd))

	multipleBom := []byte{0xef, 0xbb, 0xbf, 0xef, 0xbb, 0xbf, 'h', 'i'}
	expectedMultiple := []byte{0xef, 0xbb, 0xbf, 'h', 'i'}
	c.Equal(expectedMultiple, xbytes.StripBOM(multipleBom))

	notBom1 := []byte{0xef, 0xbb, 0xbe, 'h', 'e', 'l', 'l', 'o'} // Last byte different
	c.Equal(notBom1, xbytes.StripBOM(notBom1))

	notBom2 := []byte{0xef, 0xbc, 0xbf, 'h', 'e', 'l', 'l', 'o'} // Middle byte different
	c.Equal(notBom2, xbytes.StripBOM(notBom2))

	notBom3 := []byte{0xee, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'} // First byte different
	c.Equal(notBom3, xbytes.StripBOM(notBom3))

	withNull := []byte{0xef, 0xbb, 0xbf, 0x00, 'h', 'e', 'l', 'l', 'o'}
	expectedNull := []byte{0x00, 'h', 'e', 'l', 'l', 'o'}
	c.Equal(expectedNull, xbytes.StripBOM(withNull))

	binaryData := []byte{0xef, 0xbb, 0xbf, 0x01, 0x02, 0x03, 0xff, 0xfe}
	expectedBinary := []byte{0x01, 0x02, 0x03, 0xff, 0xfe}
	c.Equal(expectedBinary, xbytes.StripBOM(binaryData))

	jsonWithBom := []byte{0xef, 0xbb, 0xbf, '{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	expectedJSON := []byte{'{', '"', 'k', 'e', 'y', '"', ':', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	c.Equal(expectedJSON, xbytes.StripBOM(jsonWithBom))

	xmlWithBom := []byte{0xef, 0xbb, 0xbf, '<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	expectedXML := []byte{'<', '?', 'x', 'm', 'l', ' ', 'v', 'e', 'r', 's', 'i', 'o', 'n', '=', '"', '1', '.', '0', '"', '?', '>'}
	c.Equal(expectedXML, xbytes.StripBOM(xmlWithBom))
}

func TestStripBOMEdgeCases(t *testing.T) {
	c := check.New(t)

	original := []byte{0xef, 0xbb, 0xbf, 'h', 'e', 'l', 'l', 'o'}
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	result := xbytes.StripBOM(original)

	c.Equal(originalCopy, original)
	c.Equal([]byte{'h', 'e', 'l', 'l', 'o'}, result)
	// The result is a sub-slice of the original starting past the BOM.
	c.NotEqual(&original[0], &result[0])

	noBomOriginal := []byte{'h', 'e', 'l', 'l', 'o'}
	noBomResult := xbytes.StripBOM(noBomOriginal)

	c.Equal(&noBomOriginal[0], &noBomResult[0])
	c.Equal(noBomOriginal, noBomResult)

	exactBom := []byte{0xef, 0xbb, 0xbf}
	c.Equal([]byte{}, xbytes.StripBOM(exactBom))

	oneByte := []byte{0xef}
	c.Equal(oneByte, xbytes.StripBOM(oneByte))

	twoBytes := []byte{0xef, 0xbb}
	c.Equal(twoBytes, xbytes.StripBOM(twoBytes))

	var nilSlice []byte
	c.Equal(nilSlice, xbytes.StripBOM(nilSlice))

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

	chineseBom := []byte{0xef, 0xbb, 0xbf, 0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd} // BOM + 中国
	expectedChinese := []byte{0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd}
	c.Equal(expectedChinese, xbytes.StripBOM(chineseBom))

	arabicBom := []byte{0xef, 0xbb, 0xbf, 0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9} // BOM + العربية
	expectedArabic := []byte{0xd8, 0xa7, 0xd9, 0x84, 0xd8, 0xb9, 0xd8, 0xb1, 0xd8, 0xa8, 0xd9, 0x8a, 0xd8, 0xa9}
	c.Equal(expectedArabic, xbytes.StripBOM(arabicBom))

	emojiBom := []byte{0xef, 0xbb, 0xbf, 0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80} // BOM + 😀🎉🚀
	expectedEmoji := []byte{0xf0, 0x9f, 0x98, 0x80, 0xf0, 0x9f, 0x8e, 0x89, 0xf0, 0x9f, 0x9a, 0x80}
	c.Equal(expectedEmoji, xbytes.StripBOM(emojiBom))

	mathBom := []byte{0xef, 0xbb, 0xbf, 0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80} // BOM + ∑∞π
	expectedMath := []byte{0xe2, 0x88, 0x91, 0xe2, 0x88, 0x9e, 0xcf, 0x80}
	c.Equal(expectedMath, xbytes.StripBOM(mathBom))

	mixedBom := []byte{0xef, 0xbb, 0xbf, 'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	expectedMixed := []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xf0, 0x9f, 0x8c, 0x8d, '!'}
	c.Equal(expectedMixed, xbytes.StripBOM(mixedBom))
}
