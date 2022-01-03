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

// CollapseSpaces removes leading and trailing spaces and reduces any runs of two or more spaces to a single space.
func CollapseSpaces(in string) string {
	var buffer strings.Builder
	lastWasSpace := false
	for i, r := range in {
		if r == ' ' {
			if !lastWasSpace {
				if i != 0 {
					buffer.WriteByte(' ')
				}
				lastWasSpace = true
			}
		} else {
			buffer.WriteRune(r)
			lastWasSpace = false
		}
	}
	str := buffer.String()
	if lastWasSpace && len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
}
