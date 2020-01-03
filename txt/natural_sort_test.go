// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/stretchr/testify/assert"
)

var benchSet []string

func TestMain(m *testing.M) {
	flag.Parse()
	if f := flag.Lookup("test.bench"); f != nil && f.Value.String() != "" {
		initBenchSet()
	}
	atexit.Exit(m.Run())
}

func TestNaturalLess(t *testing.T) {
	testset := []struct {
		s1              string
		s2              string
		caseInsensitive bool
		less            bool
	}{
		{"0", "00", false, true},
		{"00", "0", false, false},
		{"aa", "ab", false, true},
		{"ab", "abc", false, true},
		{"abc", "ad", false, true},
		{"ab1", "ab2", false, true},
		{"ab1c", "ab1c", false, false},
		{"ab12", "abc", false, true},
		{"ab2a", "ab10", false, true},
		{"a0001", "a0000001", false, true},
		{"a10", "abcdefgh2", false, true},
		{"аб2аб", "аб10аб", false, true},
		{"2аб", "3аб", false, true},
		{"a1b", "a01b", false, true},
		{"a01b", "a1b", false, false},
		{"ab01b", "ab010b", false, true},
		{"ab010b", "ab01b", false, false},
		{"a01b001", "a001b01", false, true},
		{"a001b01", "a01b001", false, false},
		{"a1", "a1x", false, true},
		{"1ax", "1b", false, true},
		{"1b", "1ax", false, false},
		{"082", "83", false, true},
		{"083a", "9a", false, false},
		{"9a", "083a", false, true},
		{"a123", "A0123", true, true},
		{"A123", "a0123", true, true},
		{"ab010b", "ab01B", true, false},
		{"1.12.34", "1.2", false, false},
		{"1.2.34", "1.11.11", false, true},
	}
	for _, v := range testset {
		assert.Equal(t, v.less, txt.NaturalLess(v.s1, v.s2, v.caseInsensitive), fmt.Sprintf("%q < %q", v.s1, v.s2))
	}
}

func BenchmarkStdStringLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range benchSet {
			_ = benchSet[j] < benchSet[(j+1)%len(benchSet)]
		}
	}
}

func BenchmarkNaturalLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range benchSet {
			_ = txt.NaturalLess(benchSet[j], benchSet[(j+1)%len(benchSet)], false)
		}
	}
}

func BenchmarkNaturalLessCaseInsensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range benchSet {
			_ = txt.NaturalLess(benchSet[j], benchSet[(j+1)%len(benchSet)], true)
		}
	}
}

func initBenchSet() {
	rnd := rand.New(rand.NewSource(22))
	benchSet = make([]string, 20000)
	for i := range benchSet {
		strlen := rnd.Intn(6) + 3
		numlen := rnd.Intn(3) + 1
		numpos := rnd.Intn(strlen + 1)
		var num string
		for j := 0; j < numlen; j++ {
			num += strconv.Itoa(rnd.Intn(10))
		}
		var str string
		for j := 0; j < strlen+1; j++ {
			if j == numpos {
				str += num
			} else {
				str += string('a' + rnd.Intn(16))
			}
		}
		benchSet[i] = str
	}
}
