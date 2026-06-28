// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestLocalization(t *testing.T) {
	c := check.New(t)
	de := make(map[string]string)
	de["a"] = "1"
	langMap["de"] = de
	deDE := make(map[string]string)
	deDE["a"] = "2"
	langMap["de_dn"] = deDE
	Language = "de_dn.UTF-8"
	c.Equal("2", Text("a"))
	Language = "de_dn"
	c.Equal("2", Text("a"))
	Language = "de"
	c.Equal("1", Text("a"))
	Language = "xx"
	c.Equal("a", Text("a"))
	delete(langMap, "de_dn")
	Language = "de"
	c.Equal("1", Text("a"))
}

func TestEmptyTranslation(t *testing.T) {
	c := check.New(t)
	savedLanguage := Language
	savedLanguages := Languages
	defer func() {
		Language = savedLanguage
		Languages = savedLanguages
	}()
	langMap["fr"] = map[string]string{"Hello": ""}
	defer delete(langMap, "fr")
	// A key deliberately translated to "" must return "", not fall back to the original text.
	Language = "fr"
	Languages = nil
	c.Equal("", Text("Hello"))
	// A key with no translation in any language still falls back to the original text.
	c.Equal("Missing", Text("Missing"))
	// An empty translation in the primary language must not be skipped in favor of a fallback language.
	langMap["fr_fr"] = map[string]string{"Hello": "Bonjour"}
	defer delete(langMap, "fr_fr")
	Language = "fr"
	Languages = []string{"fr_fr"}
	c.Equal("", Text("Hello"))
}

func TestAltLocalization(t *testing.T) {
	c := check.New(t)
	c.Equal("Hello!", Text("Hello!"))
	SetLocalizer(func(_ string) string { return "Bonjour!" })
	c.Equal("Bonjour!", Text("Hello!"))
	SetLocalizer(nil)
	c.Equal("Hello!", Text("Hello!"))
}
