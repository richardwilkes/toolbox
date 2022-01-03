// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import "strings"

// IsTruthy returns true for "truthy" values, i.e. ones that should be interpreted as true.
func IsTruthy(in string) bool {
	in = strings.ToLower(in)
	return in == "1" || in == "true" || in == "yes" || in == "on"
}
