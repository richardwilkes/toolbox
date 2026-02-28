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

	// Test basic English vowels - lowercase
	c.True(xunicode.IsVowel('a'))
	c.True(xunicode.IsVowel('e'))
	c.True(xunicode.IsVowel('i'))
	c.True(xunicode.IsVowel('o'))
	c.True(xunicode.IsVowel('u'))

	// Test basic English vowels - uppercase
	c.True(xunicode.IsVowel('A'))
	c.True(xunicode.IsVowel('E'))
	c.True(xunicode.IsVowel('I'))
	c.True(xunicode.IsVowel('O'))
	c.True(xunicode.IsVowel('U'))

	// Test 'y' is not a vowel in IsVowel
	c.False(xunicode.IsVowel('y'))
	c.False(xunicode.IsVowel('Y'))

	// Test accented 'a' vowels
	c.True(xunicode.IsVowel('à'))
	c.True(xunicode.IsVowel('á'))
	c.True(xunicode.IsVowel('â'))
	c.True(xunicode.IsVowel('ä'))
	c.True(xunicode.IsVowel('æ'))
	c.True(xunicode.IsVowel('ã'))
	c.True(xunicode.IsVowel('å'))
	c.True(xunicode.IsVowel('ā'))

	// Test accented 'a' vowels - uppercase
	c.True(xunicode.IsVowel('À'))
	c.True(xunicode.IsVowel('Á'))
	c.True(xunicode.IsVowel('Â'))
	c.True(xunicode.IsVowel('Ä'))
	c.True(xunicode.IsVowel('Æ'))
	c.True(xunicode.IsVowel('Ã'))
	c.True(xunicode.IsVowel('Å'))
	c.True(xunicode.IsVowel('Ā'))

	// Test accented 'e' vowels
	c.True(xunicode.IsVowel('è'))
	c.True(xunicode.IsVowel('é'))
	c.True(xunicode.IsVowel('ê'))
	c.True(xunicode.IsVowel('ë'))
	c.True(xunicode.IsVowel('ē'))
	c.True(xunicode.IsVowel('ė'))
	c.True(xunicode.IsVowel('ę'))

	// Test accented 'e' vowels - uppercase
	c.True(xunicode.IsVowel('È'))
	c.True(xunicode.IsVowel('É'))
	c.True(xunicode.IsVowel('Ê'))
	c.True(xunicode.IsVowel('Ë'))
	c.True(xunicode.IsVowel('Ē'))
	c.True(xunicode.IsVowel('Ė'))
	c.True(xunicode.IsVowel('Ę'))

	// Test accented 'i' vowels
	c.True(xunicode.IsVowel('î'))
	c.True(xunicode.IsVowel('ï'))
	c.True(xunicode.IsVowel('í'))
	c.True(xunicode.IsVowel('ī'))
	c.True(xunicode.IsVowel('į'))
	c.True(xunicode.IsVowel('ì'))

	// Test accented 'i' vowels - uppercase
	c.True(xunicode.IsVowel('Î'))
	c.True(xunicode.IsVowel('Ï'))
	c.True(xunicode.IsVowel('Í'))
	c.True(xunicode.IsVowel('Ī'))
	c.True(xunicode.IsVowel('Į'))
	c.True(xunicode.IsVowel('Ì'))

	// Test accented 'o' vowels
	c.True(xunicode.IsVowel('ô'))
	c.True(xunicode.IsVowel('ö'))
	c.True(xunicode.IsVowel('ò'))
	c.True(xunicode.IsVowel('ó'))
	c.True(xunicode.IsVowel('œ'))
	c.True(xunicode.IsVowel('ø'))
	c.True(xunicode.IsVowel('ō'))
	c.True(xunicode.IsVowel('õ'))

	// Test accented 'o' vowels - uppercase
	c.True(xunicode.IsVowel('Ô'))
	c.True(xunicode.IsVowel('Ö'))
	c.True(xunicode.IsVowel('Ò'))
	c.True(xunicode.IsVowel('Ó'))
	c.True(xunicode.IsVowel('Œ'))
	c.True(xunicode.IsVowel('Ø'))
	c.True(xunicode.IsVowel('Ō'))
	c.True(xunicode.IsVowel('Õ'))

	// Test accented 'u' vowels
	c.True(xunicode.IsVowel('û'))
	c.True(xunicode.IsVowel('ü'))
	c.True(xunicode.IsVowel('ù'))
	c.True(xunicode.IsVowel('ú'))
	c.True(xunicode.IsVowel('ū'))

	// Test accented 'u' vowels - uppercase
	c.True(xunicode.IsVowel('Û'))
	c.True(xunicode.IsVowel('Ü'))
	c.True(xunicode.IsVowel('Ù'))
	c.True(xunicode.IsVowel('Ú'))
	c.True(xunicode.IsVowel('Ū'))

	// Test consonants - lowercase
	c.False(xunicode.IsVowel('b'))
	c.False(xunicode.IsVowel('c'))
	c.False(xunicode.IsVowel('d'))
	c.False(xunicode.IsVowel('f'))
	c.False(xunicode.IsVowel('g'))
	c.False(xunicode.IsVowel('h'))
	c.False(xunicode.IsVowel('j'))
	c.False(xunicode.IsVowel('k'))
	c.False(xunicode.IsVowel('l'))
	c.False(xunicode.IsVowel('m'))
	c.False(xunicode.IsVowel('n'))
	c.False(xunicode.IsVowel('p'))
	c.False(xunicode.IsVowel('q'))
	c.False(xunicode.IsVowel('r'))
	c.False(xunicode.IsVowel('s'))
	c.False(xunicode.IsVowel('t'))
	c.False(xunicode.IsVowel('v'))
	c.False(xunicode.IsVowel('w'))
	c.False(xunicode.IsVowel('x'))
	c.False(xunicode.IsVowel('z'))

	// Test consonants - uppercase
	c.False(xunicode.IsVowel('B'))
	c.False(xunicode.IsVowel('C'))
	c.False(xunicode.IsVowel('D'))
	c.False(xunicode.IsVowel('F'))
	c.False(xunicode.IsVowel('G'))
	c.False(xunicode.IsVowel('H'))
	c.False(xunicode.IsVowel('J'))
	c.False(xunicode.IsVowel('K'))
	c.False(xunicode.IsVowel('L'))
	c.False(xunicode.IsVowel('M'))
	c.False(xunicode.IsVowel('N'))
	c.False(xunicode.IsVowel('P'))
	c.False(xunicode.IsVowel('Q'))
	c.False(xunicode.IsVowel('R'))
	c.False(xunicode.IsVowel('S'))
	c.False(xunicode.IsVowel('T'))
	c.False(xunicode.IsVowel('V'))
	c.False(xunicode.IsVowel('W'))
	c.False(xunicode.IsVowel('X'))
	c.False(xunicode.IsVowel('Z'))

	// Test numbers
	c.False(xunicode.IsVowel('0'))
	c.False(xunicode.IsVowel('1'))
	c.False(xunicode.IsVowel('2'))
	c.False(xunicode.IsVowel('3'))
	c.False(xunicode.IsVowel('4'))
	c.False(xunicode.IsVowel('5'))
	c.False(xunicode.IsVowel('6'))
	c.False(xunicode.IsVowel('7'))
	c.False(xunicode.IsVowel('8'))
	c.False(xunicode.IsVowel('9'))

	// Test special characters
	c.False(xunicode.IsVowel(' '))
	c.False(xunicode.IsVowel('!'))
	c.False(xunicode.IsVowel('@'))
	c.False(xunicode.IsVowel('#'))
	c.False(xunicode.IsVowel('$'))
	c.False(xunicode.IsVowel('%'))
	c.False(xunicode.IsVowel('^'))
	c.False(xunicode.IsVowel('&'))
	c.False(xunicode.IsVowel('*'))
	c.False(xunicode.IsVowel('('))
	c.False(xunicode.IsVowel(')'))
	c.False(xunicode.IsVowel('-'))
	c.False(xunicode.IsVowel('_'))
	c.False(xunicode.IsVowel('='))
	c.False(xunicode.IsVowel('+'))
	c.False(xunicode.IsVowel('['))
	c.False(xunicode.IsVowel(']'))
	c.False(xunicode.IsVowel('{'))
	c.False(xunicode.IsVowel('}'))
	c.False(xunicode.IsVowel('\\'))
	c.False(xunicode.IsVowel('|'))
	c.False(xunicode.IsVowel(';'))
	c.False(xunicode.IsVowel(':'))
	c.False(xunicode.IsVowel('\''))
	c.False(xunicode.IsVowel('"'))
	c.False(xunicode.IsVowel(','))
	c.False(xunicode.IsVowel('.'))
	c.False(xunicode.IsVowel('<'))
	c.False(xunicode.IsVowel('>'))
	c.False(xunicode.IsVowel('/'))
	c.False(xunicode.IsVowel('?'))

	// Test Unicode characters that are not vowels
	c.False(xunicode.IsVowel('ñ'))
	c.False(xunicode.IsVowel('ç'))
	c.False(xunicode.IsVowel('ß'))
	c.False(xunicode.IsVowel('ł'))

	// Test emoji
	c.False(xunicode.IsVowel('😀'))
	c.False(xunicode.IsVowel('🚀'))
	c.False(xunicode.IsVowel('🎉'))

	// Test whitespace characters
	c.False(xunicode.IsVowel('\n'))
	c.False(xunicode.IsVowel('\t'))
	c.False(xunicode.IsVowel('\r'))

	// Test null character
	c.False(xunicode.IsVowel('\x00'))
}

func TestIsVowely(t *testing.T) {
	c := check.New(t)

	// Test that IsVowely includes all regular vowels
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

	// Test 'y' is considered a vowel in IsVowely
	c.True(xunicode.IsVowely('y'))
	c.True(xunicode.IsVowely('Y'))

	// Test accented 'y'
	c.True(xunicode.IsVowely('ÿ'))
	c.True(xunicode.IsVowely('Ÿ'))

	// Test all accented vowels work in IsVowely
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

	// Test consonants are still not vowels
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

	// Test consonants - uppercase
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

	// Test numbers
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

	// Test special characters
	c.False(xunicode.IsVowely(' '))
	c.False(xunicode.IsVowely('!'))
	c.False(xunicode.IsVowely('@'))
	c.False(xunicode.IsVowely('#'))
	c.False(xunicode.IsVowely('$'))
	c.False(xunicode.IsVowely('%'))

	// Test emoji
	c.False(xunicode.IsVowely('😀'))
	c.False(xunicode.IsVowely('🚀'))
	c.False(xunicode.IsVowely('🎉'))
}

func TestVowelCheckerInterface(t *testing.T) {
	c := check.New(t)

	// Test that IsVowel can be used as VowelChecker
	var vowelChecker xunicode.VowelChecker = xunicode.IsVowel
	c.True(vowelChecker('a'))
	c.True(vowelChecker('E'))
	c.False(vowelChecker('y'))
	c.False(vowelChecker('b'))

	// Test that IsVowely can be used as VowelChecker
	var vowelyChecker xunicode.VowelChecker = xunicode.IsVowely
	c.True(vowelyChecker('a'))
	c.True(vowelyChecker('E'))
	c.True(vowelyChecker('y'))
	c.True(vowelyChecker('Y'))
	c.False(vowelyChecker('b'))

	// Test custom VowelChecker that only accepts 'a' and 'e'
	customChecker := func(ch rune) bool {
		return ch == 'a' || ch == 'e' || ch == 'A' || ch == 'E'
	}
	c.True(customChecker('a'))
	c.True(customChecker('e'))
	c.True(customChecker('A'))
	c.True(customChecker('E'))
	c.False(customChecker('i'))
	c.False(customChecker('o'))
	c.False(customChecker('u'))
	c.False(customChecker('y'))
}

func TestVowelEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test case sensitivity is properly handled for all vowels
	vowelPairs := map[rune]rune{
		'a': 'A', 'e': 'E', 'i': 'I', 'o': 'O', 'u': 'U',
		'à': 'À', 'á': 'Á', 'â': 'Â', 'ä': 'Ä', 'æ': 'Æ', 'ã': 'Ã', 'å': 'Å', 'ā': 'Ā',
		'è': 'È', 'é': 'É', 'ê': 'Ê', 'ë': 'Ë', 'ē': 'Ē', 'ė': 'Ė', 'ę': 'Ę',
		'î': 'Î', 'ï': 'Ï', 'í': 'Í', 'ī': 'Ī', 'į': 'Į', 'ì': 'Ì',
		'ô': 'Ô', 'ö': 'Ö', 'ò': 'Ò', 'ó': 'Ó', 'œ': 'Œ', 'ø': 'Ø', 'ō': 'Ō', 'õ': 'Õ',
		'û': 'Û', 'ü': 'Ü', 'ù': 'Ù', 'ú': 'Ú', 'ū': 'Ū',
	}

	for lower, upper := range vowelPairs {
		c.True(xunicode.IsVowel(lower))
		c.True(xunicode.IsVowel(upper))
		c.True(xunicode.IsVowely(lower))
		c.True(xunicode.IsVowely(upper))
	}

	// Test 'y' case sensitivity for IsVowely
	c.True(xunicode.IsVowely('y'))
	c.True(xunicode.IsVowely('Y'))
	c.True(xunicode.IsVowely('ÿ'))
	c.True(xunicode.IsVowely('Ÿ'))

	// Test that non-vowel accented characters are still not vowels
	c.False(xunicode.IsVowel('ñ'))
	c.False(xunicode.IsVowel('Ñ'))
	c.False(xunicode.IsVowel('ç'))
	c.False(xunicode.IsVowel('Ç'))
	c.False(xunicode.IsVowely('ñ'))
	c.False(xunicode.IsVowely('Ñ'))
	c.False(xunicode.IsVowely('ç'))
	c.False(xunicode.IsVowely('Ç'))

	// Test some similar-looking but different Unicode characters
	c.True(xunicode.IsVowel('ā'))  // This is in the vowel list
	c.False(xunicode.IsVowel('ă')) // This is not in the vowel list
	c.False(xunicode.IsVowel('ą')) // This is not in the vowel list

	// Test mathematical and other Unicode symbols
	c.False(xunicode.IsVowel('α')) // Greek alpha
	c.False(xunicode.IsVowel('ε')) // Greek epsilon
	c.False(xunicode.IsVowel('ι')) // Greek iota
	c.False(xunicode.IsVowel('ο')) // Greek omicron
	c.False(xunicode.IsVowel('υ')) // Greek upsilon
	c.False(xunicode.IsVowely('α'))
	c.False(xunicode.IsVowely('ε'))
	c.False(xunicode.IsVowely('ι'))
	c.False(xunicode.IsVowely('ο'))
	c.False(xunicode.IsVowely('υ'))
}

func TestVowelComprehensiveCoverage(t *testing.T) {
	c := check.New(t)

	// Test every vowel explicitly listed in the switch cases

	// All 'a' variants
	aVowels := []rune{'a', 'à', 'á', 'â', 'ä', 'æ', 'ã', 'å', 'ā'}
	for _, vowel := range aVowels {
		c.True(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}

	// All 'e' variants
	eVowels := []rune{'e', 'è', 'é', 'ê', 'ë', 'ē', 'ė', 'ę'}
	for _, vowel := range eVowels {
		c.True(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}

	// All 'i' variants
	iVowels := []rune{'i', 'î', 'ï', 'í', 'ī', 'į', 'ì'}
	for _, vowel := range iVowels {
		c.True(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}

	// All 'o' variants
	oVowels := []rune{'o', 'ô', 'ö', 'ò', 'ó', 'œ', 'ø', 'ō', 'õ'}
	for _, vowel := range oVowels {
		c.True(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}

	// All 'u' variants
	uVowels := []rune{'u', 'û', 'ü', 'ù', 'ú', 'ū'}
	for _, vowel := range uVowels {
		c.True(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}

	// 'y' variants (only in xunicode.IsVowely)
	yVowels := []rune{'y', 'ÿ'}
	for _, vowel := range yVowels {
		c.False(xunicode.IsVowel(vowel))
		c.True(xunicode.IsVowely(vowel))
	}
}
