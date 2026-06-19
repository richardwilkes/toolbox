// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xrand_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xrand"
)

func TestNewReturnsUsableRandomizer(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	c.NotNil(r)
}

func TestIntnNonPositive(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	c.Equal(0, r.Intn(0))
	c.Equal(0, r.Intn(-1))
	c.Equal(0, r.Intn(-1000))
}

func TestIntnOne(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	for range 1000 {
		c.Equal(0, r.Intn(1))
	}
}

func TestIntnWithinRange(t *testing.T) {
	c := check.New(t)
	r := xrand.New()
	// Exercise values that span the different byte-size selections in the implementation (1, 2, 3+ bytes).
	for _, n := range []int{2, 7, 255, 256, 1000, 65535, 65536, 1 << 20, 1 << 30} {
		seen := make(map[int]bool)
		for range 10000 {
			v := r.Intn(n)
			c.True(v >= 0 && v < n, "Intn(%d) returned out-of-range value %d", n, v)
			seen[v] = true
		}
		// With this many samples over a small range, we expect more than one distinct value, which guards against a
		// stuck generator returning a constant.
		if n <= 1000 {
			c.True(len(seen) > 1, "Intn(%d) only ever produced %d distinct value(s)", n, len(seen))
		}
	}
}
