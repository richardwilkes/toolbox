// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import "strings"

// NormalizeLineEndings converts CRLF and CR into LF.
func NormalizeLineEndings(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "\r\n", "\n"), "\r", "\n")
}
