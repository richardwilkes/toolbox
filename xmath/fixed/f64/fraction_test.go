// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package f64_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xmath/fixed"
	"github.com/richardwilkes/toolbox/v2/xmath/fixed/f64"
)

func TestFraction(t *testing.T) {
	c := check.New(t)
	c.Equal(f64.FromStringForced[fixed.D4]("0.3333"), f64.NewFraction[fixed.D4]("1/3").Value())
	c.Equal(f64.FromStringForced[fixed.D4]("0.3333"), f64.NewFraction[fixed.D4]("1 / 3").Value())
	c.Equal(f64.FromStringForced[fixed.D4]("0.3333"), f64.NewFraction[fixed.D4]("-1/-3").Value())
	c.Equal(f64.From[fixed.D4, int](0), f64.NewFraction[fixed.D4]("5/0").Value())
	c.Equal(f64.From[fixed.D4, int](5), f64.NewFraction[fixed.D4]("5/1").Value())
	c.Equal(f64.From[fixed.D4, int](-5), f64.NewFraction[fixed.D4]("-5/1").Value())
	c.Equal(f64.From[fixed.D4, int](-5), f64.NewFraction[fixed.D4]("5/-1").Value())
}
