// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package xmath provides math-related utilities.
package xmath

import (
	"math"
)

const (
	// DegreesToRadians converts a value in degrees to radians when multiplied with the value.
	DegreesToRadians = math.Pi / 180
	// RadiansToDegrees converts a value in radians to degrees when multiplied with the value.
	RadiansToDegrees = 180 / math.Pi
)
