// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings

import (
	"strings"
)

// CapitalizeWords capitalizes the first letter of each word in a string, and lowercases the rest. This also trims
// leading and trailing whitespace, and replaces multiple spaces with a single space.
func CapitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = FirstToUpper(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}
