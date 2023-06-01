// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalization(t *testing.T) {
	de := make(map[string]string)
	de["a"] = "1"
	langMap["de"] = de
	deDE := make(map[string]string)
	deDE["a"] = "2"
	langMap["de_dn"] = deDE
	Language = "de_dn.UTF-8"
	assert.Equal(t, "2", Text("a"))
	Language = "de_dn"
	assert.Equal(t, "2", Text("a"))
	Language = "de"
	assert.Equal(t, "1", Text("a"))
	Language = "xx"
	assert.Equal(t, "a", Text("a"))
	delete(langMap, "de_dn")
	Language = "de"
	assert.Equal(t, "1", Text("a"))
}
