// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import "strings"

// IsTruthy returns true for "truthy" values, i.e. ones that should be interpreted as true.
func IsTruthy(in string) bool {
	switch strings.ToLower(in) {
	case "1", "t", "y", "true", "yes", "on":
		return true
	default:
		return false
	}
}
