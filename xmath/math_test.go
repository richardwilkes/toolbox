// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xmath_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xmath"
)

func TestGCD(t *testing.T) {
	c := check.New(t)
	c.Equal(6, xmath.GCD(54, 24))
	c.Equal(1, xmath.GCD(17, 13))
	c.Equal(5, xmath.GCD(-5, 0))
	c.Equal(5, xmath.GCD(0, -5))
	c.Equal(4, xmath.GCD(-8, -12))
}
