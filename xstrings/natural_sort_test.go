// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings_test

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestNaturalLess(t *testing.T) {
	c := check.New(t)

	// Test numeric strings with different lengths
	c.Equal(true, xstrings.NaturalLess("0", "00", false))
	c.Equal(false, xstrings.NaturalLess("00", "0", false))

	// Test basic alphabetical comparison
	c.Equal(true, xstrings.NaturalLess("aa", "ab", false))
	c.Equal(true, xstrings.NaturalLess("ab", "abc", false))
	c.Equal(true, xstrings.NaturalLess("abc", "ad", false))

	// Test strings with numbers - natural ordering
	c.Equal(true, xstrings.NaturalLess("ab1", "ab2", false))
	c.Equal(false, xstrings.NaturalLess("ab1c", "ab1c", false))
	c.Equal(true, xstrings.NaturalLess("ab12", "abc", false))
	c.Equal(true, xstrings.NaturalLess("ab2a", "ab10", false))

	// Test numeric strings with leading zeros
	c.Equal(true, xstrings.NaturalLess("a0001", "a0000001", false))

	// Test mixed alphanumeric comparison
	c.Equal(true, xstrings.NaturalLess("a10", "abcdefgh2", false))

	// Test unicode characters with numbers
	c.Equal(true, xstrings.NaturalLess("🚀2🚀", "🚀10🚀", false))
	c.Equal(true, xstrings.NaturalLess("2🚀", "3🚀", false))

	// Test leading zero handling
	c.Equal(true, xstrings.NaturalLess("a1b", "a01b", false))
	c.Equal(false, xstrings.NaturalLess("a01b", "a1b", false))
	c.Equal(true, xstrings.NaturalLess("ab01b", "ab010b", false))
	c.Equal(false, xstrings.NaturalLess("ab010b", "ab01b", false))

	// Test complex leading zero scenarios
	c.Equal(true, xstrings.NaturalLess("a01b001", "a001b01", false))
	c.Equal(false, xstrings.NaturalLess("a001b01", "a01b001", false))

	// Test string length differences
	c.Equal(true, xstrings.NaturalLess("a1", "a1x", false))
	c.Equal(true, xstrings.NaturalLess("a1t", "a1tx", false))
	c.Equal(true, xstrings.NaturalLess("1ax", "1b", false))
	c.Equal(false, xstrings.NaturalLess("1b", "1ax", false))

	// Test numeric comparison with leading zeros
	c.Equal(true, xstrings.NaturalLess("082", "83", false))
	c.Equal(false, xstrings.NaturalLess("083a", "9a", false))
	c.Equal(true, xstrings.NaturalLess("9a", "083a", false))

	// Test case insensitive comparison
	c.Equal(true, xstrings.NaturalLess("a123", "A0123", true))
	c.Equal(true, xstrings.NaturalLess("A123", "a0123", true))
	c.Equal(false, xstrings.NaturalLess("ab010b", "ab01B", true))

	// Test version-like strings
	c.Equal(false, xstrings.NaturalLess("1.12.34", "1.2", false))
	c.Equal(true, xstrings.NaturalLess("1.2.34", "1.11.11", false))
}

func TestNaturalCmp(t *testing.T) {
	c := check.New(t)

	// Test equal strings
	c.Equal(0, xstrings.NaturalCmp("abc", "abc", false))
	c.Equal(0, xstrings.NaturalCmp("123", "123", false))
	c.Equal(0, xstrings.NaturalCmp("", "", false))

	// Test less than cases (should return -1)
	c.Equal(-1, xstrings.NaturalCmp("a", "b", false))
	c.Equal(-1, xstrings.NaturalCmp("1", "2", false))
	c.Equal(-1, xstrings.NaturalCmp("a1", "a2", false))
	c.Equal(-1, xstrings.NaturalCmp("a2", "a10", false))
	c.Equal(-1, xstrings.NaturalCmp("0", "00", false))

	// Test greater than cases (should return 1)
	c.Equal(1, xstrings.NaturalCmp("b", "a", false))
	c.Equal(1, xstrings.NaturalCmp("2", "1", false))
	c.Equal(1, xstrings.NaturalCmp("a2", "a1", false))
	c.Equal(1, xstrings.NaturalCmp("a10", "a2", false))
	c.Equal(1, xstrings.NaturalCmp("00", "0", false))

	// Test case sensitivity
	c.Equal(1, xstrings.NaturalCmp("a", "A", false))
	c.Equal(-1, xstrings.NaturalCmp("A", "a", false))

	// Test case insensitive comparison (this actually falls back to case sensitive if they would otherwise be equal)
	c.Equal(1, xstrings.NaturalCmp("a", "A", true))
	c.Equal(-1, xstrings.NaturalCmp("ABC", "abc", true))
	c.Equal(-1, xstrings.NaturalCmp("a", "B", true))
	c.Equal(1, xstrings.NaturalCmp("B", "a", true))

	// Test mixed case with numbers
	c.Equal(-1, xstrings.NaturalCmp("a1", "A2", true))
	c.Equal(1, xstrings.NaturalCmp("A2", "a1", true))

	// Test leading zeros
	c.Equal(1, xstrings.NaturalCmp("a001", "a01", false))
	c.Equal(-1, xstrings.NaturalCmp("a01", "a001", false))

	// Test version-like strings
	c.Equal(1, xstrings.NaturalCmp("1.12.34", "1.2", false))
	c.Equal(-1, xstrings.NaturalCmp("1.2.34", "1.11.11", false))

	// Test empty strings
	c.Equal(-1, xstrings.NaturalCmp("", "a", false))
	c.Equal(1, xstrings.NaturalCmp("a", "", false))

	// Test strings with different lengths
	c.Equal(-1, xstrings.NaturalCmp("a", "ab", false))
	c.Equal(1, xstrings.NaturalCmp("ab", "a", false))
}

func TestSortStringsNaturalAscending(t *testing.T) {
	c := check.New(t)

	// Test empty slice
	empty := []string{}
	xstrings.SortStringsNaturalAscending(empty)
	c.Equal(0, len(empty))

	// Test single element
	single := []string{"hello"}
	xstrings.SortStringsNaturalAscending(single)
	c.Equal([]string{"hello"}, single)

	// Test basic alphabetical sorting
	basic := []string{"c", "a", "b"}
	xstrings.SortStringsNaturalAscending(basic)
	c.Equal([]string{"a", "b", "c"}, basic)

	// Test natural number sorting
	numbers := []string{"item10", "item2", "item1", "item20"}
	xstrings.SortStringsNaturalAscending(numbers)
	c.Equal([]string{"item1", "item2", "item10", "item20"}, numbers)

	// Test version-like strings
	versions := []string{"v1.12.0", "v1.2.0", "v1.11.0", "v1.1.0"}
	xstrings.SortStringsNaturalAscending(versions)
	c.Equal([]string{"v1.1.0", "v1.2.0", "v1.11.0", "v1.12.0"}, versions)

	// Test with leading zeros
	zeros := []string{"file001", "file01", "file1", "file10"}
	xstrings.SortStringsNaturalAscending(zeros)
	c.Equal([]string{"file1", "file01", "file001", "file10"}, zeros)

	// Test mixed case (case insensitive)
	mixed := []string{"File2", "file1", "FILE10", "file3"}
	xstrings.SortStringsNaturalAscending(mixed)
	c.Equal([]string{"file1", "File2", "file3", "FILE10"}, mixed)

	// Test comp mixed content
	comp := []string{"a10b", "a2b", "a1b", "ab", "a", "b10", "b2"}
	xstrings.SortStringsNaturalAscending(comp)
	c.Equal([]string{"a", "a1b", "a2b", "a10b", "ab", "b2", "b10"}, comp)

	// Test identical strings
	identical := []string{"same", "same", "same"}
	xstrings.SortStringsNaturalAscending(identical)
	c.Equal([]string{"same", "same", "same"}, identical)

	// Test unicode characters
	unicode := []string{"📝10", "📝2", "📝1"}
	xstrings.SortStringsNaturalAscending(unicode)
	c.Equal([]string{"📝1", "📝2", "📝10"}, unicode)
}

func TestSortStringsNaturalDescending(t *testing.T) {
	c := check.New(t)

	// Test empty slice
	empty := []string{}
	xstrings.SortStringsNaturalDescending(empty)
	c.Equal(0, len(empty))

	// Test single element
	single := []string{"hello"}
	xstrings.SortStringsNaturalDescending(single)
	c.Equal([]string{"hello"}, single)

	// Test basic alphabetical sorting (reverse)
	basic := []string{"a", "b", "c"}
	xstrings.SortStringsNaturalDescending(basic)
	c.Equal([]string{"c", "b", "a"}, basic)

	// Test natural number sorting (reverse)
	numbers := []string{"item1", "item2", "item10", "item20"}
	xstrings.SortStringsNaturalDescending(numbers)
	c.Equal([]string{"item20", "item10", "item2", "item1"}, numbers)

	// Test version-like strings (reverse)
	versions := []string{"v1.1.0", "v1.2.0", "v1.11.0", "v1.12.0"}
	xstrings.SortStringsNaturalDescending(versions)
	c.Equal([]string{"v1.12.0", "v1.11.0", "v1.2.0", "v1.1.0"}, versions)

	// Test with leading zeros (reverse)
	zeros := []string{"file1", "file01", "file001", "file10"}
	xstrings.SortStringsNaturalDescending(zeros)
	c.Equal([]string{"file10", "file001", "file01", "file1"}, zeros)

	// Test mixed case (case insensitive, reverse)
	mixed := []string{"file1", "File2", "file3", "FILE10"}
	xstrings.SortStringsNaturalDescending(mixed)
	c.Equal([]string{"FILE10", "file3", "File2", "file1"}, mixed)

	// Test comp mixed content (reverse)
	comp := []string{"a", "a1b", "a2b", "a10b", "ab", "b2", "b10"}
	xstrings.SortStringsNaturalDescending(comp)
	c.Equal([]string{"b10", "b2", "ab", "a10b", "a2b", "a1b", "a"}, comp)

	// Test identical strings
	identical := []string{"same", "same", "same"}
	xstrings.SortStringsNaturalDescending(identical)
	c.Equal([]string{"same", "same", "same"}, identical)

	// Test unicode characters (reverse)
	unicode := []string{"📝1", "📝2", "📝10"}
	xstrings.SortStringsNaturalDescending(unicode)
	c.Equal([]string{"📝10", "📝2", "📝1"}, unicode)
}

func BenchmarkStdStringLess(b *testing.B) {
	benchSet := createBenchSet()
	for b.Loop() {
		for j := range benchSet {
			_ = benchSet[j] < benchSet[(j+1)%len(benchSet)]
		}
	}
}

func BenchmarkNaturalLess(b *testing.B) {
	benchSet := createBenchSet()
	for b.Loop() {
		for j := range benchSet {
			_ = xstrings.NaturalLess(benchSet[j], benchSet[(j+1)%len(benchSet)], false)
		}
	}
}

func BenchmarkNaturalLessCaseInsensitive(b *testing.B) {
	benchSet := createBenchSet()
	for b.Loop() {
		for j := range benchSet {
			_ = xstrings.NaturalLess(benchSet[j], benchSet[(j+1)%len(benchSet)], true)
		}
	}
}

func createBenchSet() []string {
	rnd := rand.New(rand.NewPCG(22, 1967)) //nolint:gosec // Use of weak prng is fine here
	benchSet := make([]string, 20000)
	for i := range benchSet {
		strlen := rnd.IntN(6) + 3
		numlen := rnd.IntN(3) + 1
		numpos := rnd.IntN(strlen + 1)
		var num string
		for range numlen {
			num += strconv.Itoa(rnd.IntN(10))
		}
		var str string
		for j := range strlen + 1 {
			if j == numpos {
				str += num
			} else {
				str += string(rune('a' + rnd.IntN(16)))
			}
		}
		benchSet[i] = str
	}
	return benchSet
}
