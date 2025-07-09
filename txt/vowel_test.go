package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestIsVowel(t *testing.T) {
	c := check.New(t)

	// Test basic English vowels - lowercase
	c.True(IsVowel('a'))
	c.True(IsVowel('e'))
	c.True(IsVowel('i'))
	c.True(IsVowel('o'))
	c.True(IsVowel('u'))

	// Test basic English vowels - uppercase
	c.True(IsVowel('A'))
	c.True(IsVowel('E'))
	c.True(IsVowel('I'))
	c.True(IsVowel('O'))
	c.True(IsVowel('U'))

	// Test 'y' is not a vowel in IsVowel
	c.False(IsVowel('y'))
	c.False(IsVowel('Y'))

	// Test accented 'a' vowels
	c.True(IsVowel('à'))
	c.True(IsVowel('á'))
	c.True(IsVowel('â'))
	c.True(IsVowel('ä'))
	c.True(IsVowel('æ'))
	c.True(IsVowel('ã'))
	c.True(IsVowel('å'))
	c.True(IsVowel('ā'))

	// Test accented 'a' vowels - uppercase
	c.True(IsVowel('À'))
	c.True(IsVowel('Á'))
	c.True(IsVowel('Â'))
	c.True(IsVowel('Ä'))
	c.True(IsVowel('Æ'))
	c.True(IsVowel('Ã'))
	c.True(IsVowel('Å'))
	c.True(IsVowel('Ā'))

	// Test accented 'e' vowels
	c.True(IsVowel('è'))
	c.True(IsVowel('é'))
	c.True(IsVowel('ê'))
	c.True(IsVowel('ë'))
	c.True(IsVowel('ē'))
	c.True(IsVowel('ė'))
	c.True(IsVowel('ę'))

	// Test accented 'e' vowels - uppercase
	c.True(IsVowel('È'))
	c.True(IsVowel('É'))
	c.True(IsVowel('Ê'))
	c.True(IsVowel('Ë'))
	c.True(IsVowel('Ē'))
	c.True(IsVowel('Ė'))
	c.True(IsVowel('Ę'))

	// Test accented 'i' vowels
	c.True(IsVowel('î'))
	c.True(IsVowel('ï'))
	c.True(IsVowel('í'))
	c.True(IsVowel('ī'))
	c.True(IsVowel('į'))
	c.True(IsVowel('ì'))

	// Test accented 'i' vowels - uppercase
	c.True(IsVowel('Î'))
	c.True(IsVowel('Ï'))
	c.True(IsVowel('Í'))
	c.True(IsVowel('Ī'))
	c.True(IsVowel('Į'))
	c.True(IsVowel('Ì'))

	// Test accented 'o' vowels
	c.True(IsVowel('ô'))
	c.True(IsVowel('ö'))
	c.True(IsVowel('ò'))
	c.True(IsVowel('ó'))
	c.True(IsVowel('œ'))
	c.True(IsVowel('ø'))
	c.True(IsVowel('ō'))
	c.True(IsVowel('õ'))

	// Test accented 'o' vowels - uppercase
	c.True(IsVowel('Ô'))
	c.True(IsVowel('Ö'))
	c.True(IsVowel('Ò'))
	c.True(IsVowel('Ó'))
	c.True(IsVowel('Œ'))
	c.True(IsVowel('Ø'))
	c.True(IsVowel('Ō'))
	c.True(IsVowel('Õ'))

	// Test accented 'u' vowels
	c.True(IsVowel('û'))
	c.True(IsVowel('ü'))
	c.True(IsVowel('ù'))
	c.True(IsVowel('ú'))
	c.True(IsVowel('ū'))

	// Test accented 'u' vowels - uppercase
	c.True(IsVowel('Û'))
	c.True(IsVowel('Ü'))
	c.True(IsVowel('Ù'))
	c.True(IsVowel('Ú'))
	c.True(IsVowel('Ū'))

	// Test consonants - lowercase
	c.False(IsVowel('b'))
	c.False(IsVowel('c'))
	c.False(IsVowel('d'))
	c.False(IsVowel('f'))
	c.False(IsVowel('g'))
	c.False(IsVowel('h'))
	c.False(IsVowel('j'))
	c.False(IsVowel('k'))
	c.False(IsVowel('l'))
	c.False(IsVowel('m'))
	c.False(IsVowel('n'))
	c.False(IsVowel('p'))
	c.False(IsVowel('q'))
	c.False(IsVowel('r'))
	c.False(IsVowel('s'))
	c.False(IsVowel('t'))
	c.False(IsVowel('v'))
	c.False(IsVowel('w'))
	c.False(IsVowel('x'))
	c.False(IsVowel('z'))

	// Test consonants - uppercase
	c.False(IsVowel('B'))
	c.False(IsVowel('C'))
	c.False(IsVowel('D'))
	c.False(IsVowel('F'))
	c.False(IsVowel('G'))
	c.False(IsVowel('H'))
	c.False(IsVowel('J'))
	c.False(IsVowel('K'))
	c.False(IsVowel('L'))
	c.False(IsVowel('M'))
	c.False(IsVowel('N'))
	c.False(IsVowel('P'))
	c.False(IsVowel('Q'))
	c.False(IsVowel('R'))
	c.False(IsVowel('S'))
	c.False(IsVowel('T'))
	c.False(IsVowel('V'))
	c.False(IsVowel('W'))
	c.False(IsVowel('X'))
	c.False(IsVowel('Z'))

	// Test numbers
	c.False(IsVowel('0'))
	c.False(IsVowel('1'))
	c.False(IsVowel('2'))
	c.False(IsVowel('3'))
	c.False(IsVowel('4'))
	c.False(IsVowel('5'))
	c.False(IsVowel('6'))
	c.False(IsVowel('7'))
	c.False(IsVowel('8'))
	c.False(IsVowel('9'))

	// Test special characters
	c.False(IsVowel(' '))
	c.False(IsVowel('!'))
	c.False(IsVowel('@'))
	c.False(IsVowel('#'))
	c.False(IsVowel('$'))
	c.False(IsVowel('%'))
	c.False(IsVowel('^'))
	c.False(IsVowel('&'))
	c.False(IsVowel('*'))
	c.False(IsVowel('('))
	c.False(IsVowel(')'))
	c.False(IsVowel('-'))
	c.False(IsVowel('_'))
	c.False(IsVowel('='))
	c.False(IsVowel('+'))
	c.False(IsVowel('['))
	c.False(IsVowel(']'))
	c.False(IsVowel('{'))
	c.False(IsVowel('}'))
	c.False(IsVowel('\\'))
	c.False(IsVowel('|'))
	c.False(IsVowel(';'))
	c.False(IsVowel(':'))
	c.False(IsVowel('\''))
	c.False(IsVowel('"'))
	c.False(IsVowel(','))
	c.False(IsVowel('.'))
	c.False(IsVowel('<'))
	c.False(IsVowel('>'))
	c.False(IsVowel('/'))
	c.False(IsVowel('?'))

	// Test Unicode characters that are not vowels
	c.False(IsVowel('ñ'))
	c.False(IsVowel('ç'))
	c.False(IsVowel('ß'))
	c.False(IsVowel('ł'))

	// Test emoji
	c.False(IsVowel('😀'))
	c.False(IsVowel('🚀'))
	c.False(IsVowel('🎉'))

	// Test whitespace characters
	c.False(IsVowel('\n'))
	c.False(IsVowel('\t'))
	c.False(IsVowel('\r'))

	// Test null character
	c.False(IsVowel('\x00'))
}

func TestIsVowely(t *testing.T) {
	c := check.New(t)

	// Test that IsVowely includes all regular vowels
	c.True(IsVowely('a'))
	c.True(IsVowely('e'))
	c.True(IsVowely('i'))
	c.True(IsVowely('o'))
	c.True(IsVowely('u'))
	c.True(IsVowely('A'))
	c.True(IsVowely('E'))
	c.True(IsVowely('I'))
	c.True(IsVowely('O'))
	c.True(IsVowely('U'))

	// Test 'y' is considered a vowel in IsVowely
	c.True(IsVowely('y'))
	c.True(IsVowely('Y'))

	// Test accented 'y'
	c.True(IsVowely('ÿ'))
	c.True(IsVowely('Ÿ'))

	// Test all accented vowels work in IsVowely
	c.True(IsVowely('à'))
	c.True(IsVowely('é'))
	c.True(IsVowely('î'))
	c.True(IsVowely('ô'))
	c.True(IsVowely('ü'))
	c.True(IsVowely('À'))
	c.True(IsVowely('É'))
	c.True(IsVowely('Î'))
	c.True(IsVowely('Ô'))
	c.True(IsVowely('Ü'))

	// Test consonants are still not vowels
	c.False(IsVowely('b'))
	c.False(IsVowely('c'))
	c.False(IsVowely('d'))
	c.False(IsVowely('f'))
	c.False(IsVowely('g'))
	c.False(IsVowely('h'))
	c.False(IsVowely('j'))
	c.False(IsVowely('k'))
	c.False(IsVowely('l'))
	c.False(IsVowely('m'))
	c.False(IsVowely('n'))
	c.False(IsVowely('p'))
	c.False(IsVowely('q'))
	c.False(IsVowely('r'))
	c.False(IsVowely('s'))
	c.False(IsVowely('t'))
	c.False(IsVowely('v'))
	c.False(IsVowely('w'))
	c.False(IsVowely('x'))
	c.False(IsVowely('z'))

	// Test consonants - uppercase
	c.False(IsVowely('B'))
	c.False(IsVowely('C'))
	c.False(IsVowely('D'))
	c.False(IsVowely('F'))
	c.False(IsVowely('G'))
	c.False(IsVowely('H'))
	c.False(IsVowely('J'))
	c.False(IsVowely('K'))
	c.False(IsVowely('L'))
	c.False(IsVowely('M'))
	c.False(IsVowely('N'))
	c.False(IsVowely('P'))
	c.False(IsVowely('Q'))
	c.False(IsVowely('R'))
	c.False(IsVowely('S'))
	c.False(IsVowely('T'))
	c.False(IsVowely('V'))
	c.False(IsVowely('W'))
	c.False(IsVowely('X'))
	c.False(IsVowely('Z'))

	// Test numbers
	c.False(IsVowely('0'))
	c.False(IsVowely('1'))
	c.False(IsVowely('2'))
	c.False(IsVowely('3'))
	c.False(IsVowely('4'))
	c.False(IsVowely('5'))
	c.False(IsVowely('6'))
	c.False(IsVowely('7'))
	c.False(IsVowely('8'))
	c.False(IsVowely('9'))

	// Test special characters
	c.False(IsVowely(' '))
	c.False(IsVowely('!'))
	c.False(IsVowely('@'))
	c.False(IsVowely('#'))
	c.False(IsVowely('$'))
	c.False(IsVowely('%'))

	// Test emoji
	c.False(IsVowely('😀'))
	c.False(IsVowely('🚀'))
	c.False(IsVowely('🎉'))
}

func TestVowelCheckerInterface(t *testing.T) {
	c := check.New(t)

	// Test that IsVowel can be used as VowelChecker
	var vowelChecker VowelChecker = IsVowel
	c.True(vowelChecker('a'))
	c.True(vowelChecker('E'))
	c.False(vowelChecker('y'))
	c.False(vowelChecker('b'))

	// Test that IsVowely can be used as VowelChecker
	var vowelyChecker VowelChecker = IsVowely
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
		c.True(IsVowel(lower))
		c.True(IsVowel(upper))
		c.True(IsVowely(lower))
		c.True(IsVowely(upper))
	}

	// Test 'y' case sensitivity for IsVowely
	c.True(IsVowely('y'))
	c.True(IsVowely('Y'))
	c.True(IsVowely('ÿ'))
	c.True(IsVowely('Ÿ'))

	// Test that non-vowel accented characters are still not vowels
	c.False(IsVowel('ñ'))
	c.False(IsVowel('Ñ'))
	c.False(IsVowel('ç'))
	c.False(IsVowel('Ç'))
	c.False(IsVowely('ñ'))
	c.False(IsVowely('Ñ'))
	c.False(IsVowely('ç'))
	c.False(IsVowely('Ç'))

	// Test some similar-looking but different Unicode characters
	c.True(IsVowel('ā'))  // This is in the vowel list
	c.False(IsVowel('ă')) // This is not in the vowel list
	c.False(IsVowel('ą')) // This is not in the vowel list

	// Test mathematical and other Unicode symbols
	c.False(IsVowel('α')) // Greek alpha
	c.False(IsVowel('ε')) // Greek epsilon
	c.False(IsVowel('ι')) // Greek iota
	c.False(IsVowel('ο')) // Greek omicron
	c.False(IsVowel('υ')) // Greek upsilon
	c.False(IsVowely('α'))
	c.False(IsVowely('ε'))
	c.False(IsVowely('ι'))
	c.False(IsVowely('ο'))
	c.False(IsVowely('υ'))
}

func TestVowelComprehensiveCoverage(t *testing.T) {
	c := check.New(t)

	// Test every vowel explicitly listed in the switch cases

	// All 'a' variants
	aVowels := []rune{'a', 'à', 'á', 'â', 'ä', 'æ', 'ã', 'å', 'ā'}
	for _, vowel := range aVowels {
		c.True(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}

	// All 'e' variants
	eVowels := []rune{'e', 'è', 'é', 'ê', 'ë', 'ē', 'ė', 'ę'}
	for _, vowel := range eVowels {
		c.True(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}

	// All 'i' variants
	iVowels := []rune{'i', 'î', 'ï', 'í', 'ī', 'į', 'ì'}
	for _, vowel := range iVowels {
		c.True(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}

	// All 'o' variants
	oVowels := []rune{'o', 'ô', 'ö', 'ò', 'ó', 'œ', 'ø', 'ō', 'õ'}
	for _, vowel := range oVowels {
		c.True(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}

	// All 'u' variants
	uVowels := []rune{'u', 'û', 'ü', 'ù', 'ú', 'ū'}
	for _, vowel := range uVowels {
		c.True(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}

	// 'y' variants (only in IsVowely)
	yVowels := []rune{'y', 'ÿ'}
	for _, vowel := range yVowels {
		c.False(IsVowel(vowel))
		c.True(IsVowely(vowel))
	}
}
