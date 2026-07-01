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
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
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

// TestLoadRejectsTrailingContent verifies that a key or value line with content after the closing quote is rejected
// with a warning rather than silently truncated, while well-formed lines still round-trip.
func TestLoadRejectsTrailingContent(t *testing.T) {
	c := check.New(t)

	// Capture the warnings load() emits.
	var logBuf bytes.Buffer
	restore := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelWarn})))
	defer slog.SetDefault(restore)

	dir := t.TempDir()
	savedDir := Dir
	Dir = dir
	defer func() { Dir = savedDir }()

	const name = "trailing.i18n"
	content := "k:\"foo\" bar\n" + // trailing "bar" after the closing quote makes this an invalid key
		"v:\"dropped\"\n" + // value for the rejected key; there is no live key, so it is discarded too
		"k:\"good\"\n" + // a well-formed pair that must still round-trip
		"v:\"translated\"\n"
	c.NoError(os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))

	load(name)
	defer delete(langMap, "trailing")

	translations := langMap["trailing"]
	c.NotNil(translations)
	// The malformed line must not register a truncated "foo" key.
	_, hasFoo := translations["foo"]
	c.False(hasFoo)
	// The well-formed pair still round-trips.
	c.Equal("translated", translations["good"])
	// A warning was logged for the invalid key.
	c.Contains(logBuf.String(), "ignoring invalid key")

	// Sanity check: a value line with trailing content is likewise rejected, leaving the key without a value.
	logBuf.Reset()
	const name2 = "badvalue.i18n"
	c.NoError(os.WriteFile(filepath.Join(dir, name2), []byte("k:\"key\"\nv:\"val\" extra\n"), 0o600))
	load(name2)
	defer delete(langMap, "badvalue")
	badValue := langMap["badvalue"]
	c.NotNil(badValue)
	_, hasKey := badValue["key"]
	c.False(hasKey)
	c.Contains(logBuf.String(), "ignoring invalid value")
}
