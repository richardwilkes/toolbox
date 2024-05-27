/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/txt"
)

func TestCollapseSpaces(t *testing.T) {
	data := []string{
		"123", "123",
		" 123", "123",
		" 123 ", "123",
		"    abc  ", "abc",
		"  a b c   d", "a b c d",
		"", "",
		" ", "",
	}
	for i := 0; i < len(data); i += 2 {
		check.Equal(t, data[i+1], txt.CollapseSpaces(data[i]))
	}
}
