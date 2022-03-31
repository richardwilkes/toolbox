// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fixed

import (
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d1"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d10"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d11"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d12"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d13"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d14"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d15"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d16"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d2"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d3"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d4"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d5"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d6"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d7"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d8"
	"github.com/richardwilkes/toolbox/xmath/fixed/f128d9"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d1"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d2"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d3"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d5"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d6"
)

// F64 defines the set of 64-bit fixed-point types.
type F64 interface {
	f64d1.Int | f64d2.Int | f64d3.Int | f64d4.Int | f64d5.Int | f64d6.Int
}

// F128 defines the set of 128-bit fixed-point types.
type F128 interface {
	f128d1.Int | f128d2.Int | f128d3.Int | f128d4.Int | f128d5.Int | f128d6.Int | f128d7.Int | f128d8.Int | f128d9.Int |
		f128d10.Int | f128d11.Int | f128d12.Int | f128d13.Int | f128d14.Int | f128d15.Int | f128d16.Int
}
