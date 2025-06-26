// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f128_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128"
)

func TestFraction(t *testing.T) {
	check.Equal(t, f128.FromStringForced[fixed.D4]("0.3333"), f128.NewFraction[fixed.D4]("1/3").Value())
	check.Equal(t, f128.FromStringForced[fixed.D4]("0.3333"), f128.NewFraction[fixed.D4]("1 / 3").Value())
	check.Equal(t, f128.FromStringForced[fixed.D4]("0.3333"), f128.NewFraction[fixed.D4]("-1/-3").Value())
	check.Equal(t, f128.From[fixed.D4, int](0), f128.NewFraction[fixed.D4]("5/0").Value())
	check.Equal(t, f128.From[fixed.D4, int](5), f128.NewFraction[fixed.D4]("5/1").Value())
	check.Equal(t, f128.From[fixed.D4, int](-5), f128.NewFraction[fixed.D4]("-5/1").Value())
	check.Equal(t, f128.From[fixed.D4, int](-5), f128.NewFraction[fixed.D4]("5/-1").Value())
}
