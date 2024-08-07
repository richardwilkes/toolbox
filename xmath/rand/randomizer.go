// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package rand

// Randomizer defines a source of random integer values.
type Randomizer interface {
	// Intn returns a non-negative random number from 0 to n-1. If n <= 0, the implementation should return 0.
	Intn(n int) int
}
