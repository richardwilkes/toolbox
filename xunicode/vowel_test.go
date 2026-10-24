// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xunicode_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xunicode"
)

func TestIsVowel(t *testing.T) {
	c := check.New(t)
	for _, r := range []rune{
		'a', 'e', 'i', 'o', 'u',
		'A', 'E', 'I', 'O', 'U',
		'à', 'á', 'â', 'ä', 'æ', 'ã', 'å', 'ā',
		'À', 'Á', 'Â', 'Ä', 'Æ', 'Ã', 'Å', 'Ā',
		'è', 'é', 'ê', 'ë', 'ē', 'ė', 'ę',
		'È', 'É', 'Ê', 'Ë', 'Ē', 'Ė', 'Ę',
		'î', 'ï', 'í', 'ī', 'į', 'ì',
		'Î', 'Ï', 'Í', 'Ī', 'Į', 'Ì',
		'ô', 'ö', 'ò', 'ó', 'œ', 'ø', 'ō', 'õ',
		'Ô', 'Ö', 'Ò', 'Ó', 'Œ', 'Ø', 'Ō', 'Õ',
		'û', 'ü', 'ù', 'ú', 'ū',
		'Û', 'Ü', 'Ù', 'Ú', 'Ū',
	} {
		c.True(xunicode.IsVowel(r), "IsVowel(%q)", r)
	}
	for _, r := range []rune{
		'y', 'Y',
		'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z',
		'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		' ', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+', '[',
		']', '{', '}', '\\', '|', ';', ':', '\'', '"', ',', '.', '<', '>', '/', '?',
		'ñ', 'Ñ', 'ç', 'Ç', 'ß', 'ł',
		'ă', 'ą', 'ÿ',
		'α', 'ε', 'ι', 'ο', 'υ', // Greek
		'😀', '🚀', '🎉',
		'\n', '\t', '\r', '\x00',
	} {
		c.False(xunicode.IsVowel(r), "IsVowel(%q)", r)
	}
}

func TestIsVowely(t *testing.T) {
	c := check.New(t)

	c.True(xunicode.IsVowely('a'))
	c.True(xunicode.IsVowely('e'))
	c.True(xunicode.IsVowely('i'))
	c.True(xunicode.IsVowely('o'))
	c.True(xunicode.IsVowely('u'))
	c.True(xunicode.IsVowely('A'))
	c.True(xunicode.IsVowely('E'))
	c.True(xunicode.IsVowely('I'))
	c.True(xunicode.IsVowely('O'))
	c.True(xunicode.IsVowely('U'))

	c.True(xunicode.IsVowely('y'))
	c.True(xunicode.IsVowely('Y'))

	c.True(xunicode.IsVowely('ÿ'))
	c.True(xunicode.IsVowely('Ÿ'))

	c.True(xunicode.IsVowely('à'))
	c.True(xunicode.IsVowely('é'))
	c.True(xunicode.IsVowely('î'))
	c.True(xunicode.IsVowely('ô'))
	c.True(xunicode.IsVowely('ü'))
	c.True(xunicode.IsVowely('À'))
	c.True(xunicode.IsVowely('É'))
	c.True(xunicode.IsVowely('Î'))
	c.True(xunicode.IsVowely('Ô'))
	c.True(xunicode.IsVowely('Ü'))

	c.False(xunicode.IsVowely('b'))
	c.False(xunicode.IsVowely('c'))
	c.False(xunicode.IsVowely('d'))
	c.False(xunicode.IsVowely('f'))
	c.False(xunicode.IsVowely('g'))
	c.False(xunicode.IsVowely('h'))
	c.False(xunicode.IsVowely('j'))
	c.False(xunicode.IsVowely('k'))
	c.False(xunicode.IsVowely('l'))
	c.False(xunicode.IsVowely('m'))
	c.False(xunicode.IsVowely('n'))
	c.False(xunicode.IsVowely('p'))
	c.False(xunicode.IsVowely('q'))
	c.False(xunicode.IsVowely('r'))
	c.False(xunicode.IsVowely('s'))
	c.False(xunicode.IsVowely('t'))
	c.False(xunicode.IsVowely('v'))
	c.False(xunicode.IsVowely('w'))
	c.False(xunicode.IsVowely('x'))
	c.False(xunicode.IsVowely('z'))

	c.False(xunicode.IsVowely('B'))
	c.False(xunicode.IsVowely('C'))
	c.False(xunicode.IsVowely('D'))
	c.False(xunicode.IsVowely('F'))
	c.False(xunicode.IsVowely('G'))
	c.False(xunicode.IsVowely('H'))
	c.False(xunicode.IsVowely('J'))
	c.False(xunicode.IsVowely('K'))
	c.False(xunicode.IsVowely('L'))
	c.False(xunicode.IsVowely('M'))
	c.False(xunicode.IsVowely('N'))
	c.False(xunicode.IsVowely('P'))
	c.False(xunicode.IsVowely('Q'))
	c.False(xunicode.IsVowely('R'))
	c.False(xunicode.IsVowely('S'))
	c.False(xunicode.IsVowely('T'))
	c.False(xunicode.IsVowely('V'))
	c.False(xunicode.IsVowely('W'))
	c.False(xunicode.IsVowely('X'))
	c.False(xunicode.IsVowely('Z'))

	c.False(xunicode.IsVowely('0'))
	c.False(xunicode.IsVowely('1'))
	c.False(xunicode.IsVowely('2'))
	c.False(xunicode.IsVowely('3'))
	c.False(xunicode.IsVowely('4'))
	c.False(xunicode.IsVowely('5'))
	c.False(xunicode.IsVowely('6'))
	c.False(xunicode.IsVowely('7'))
	c.False(xunicode.IsVowely('8'))
	c.False(xunicode.IsVowely('9'))

	c.False(xunicode.IsVowely(' '))
	c.False(xunicode.IsVowely('!'))
	c.False(xunicode.IsVowely('@'))
	c.False(xunicode.IsVowely('#'))
	c.False(xunicode.IsVowely('$'))
	c.False(xunicode.IsVowely('%'))

	c.False(xunicode.IsVowely('😀'))
	c.False(xunicode.IsVowely('🚀'))
	c.False(xunicode.IsVowely('🎉'))
}
