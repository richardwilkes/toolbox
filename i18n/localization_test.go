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

	"github.com/richardwilkes/toolbox/check"
)

func TestLocalization(t *testing.T) {
	de := make(map[string]string)
	de["a"] = "1"
	langMap["de"] = de
	deDE := make(map[string]string)
	deDE["a"] = "2"
	langMap["de_dn"] = deDE
	Language = "de_dn.UTF-8"
	check.Equal(t, "2", Text("a"))
	Language = "de_dn"
	check.Equal(t, "2", Text("a"))
	Language = "de"
	check.Equal(t, "1", Text("a"))
	Language = "xx"
	check.Equal(t, "a", Text("a"))
	delete(langMap, "de_dn")
	Language = "de"
	check.Equal(t, "1", Text("a"))
}

func TestAltLocalization(t *testing.T) {
	check.Equal(t, "Hello!", Text("Hello!"))
	SetLocalizer(func(_ string) string { return "Bonjour!" })
	check.Equal(t, "Bonjour!", Text("Hello!"))
	SetLocalizer(nil)
	check.Equal(t, "Hello!", Text("Hello!"))
}
