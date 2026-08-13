// Copyright (c) 2021-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package uti_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/uti"
)

func TestConformsTo(t *testing.T) {
	c := check.New(t)

	c.True(uti.PlainText.ConformsTo(uti.Content))
	c.False(uti.Content.ConformsTo(uti.PlainText))
	c.True(uti.UTF8PlainText.ConformsTo(uti.Content))
	c.True(uti.UTF8PlainText.ConformsTo(uti.PlainText))
	c.False(uti.PlainText.ConformsTo(uti.UTF8PlainText))
}

func TestByExtension(t *testing.T) {
	c := check.New(t)

	// The built-in data types must be reachable by their file extensions. This was broken when byExtension was
	// populated from MimeTypes rather than Extensions.
	c.Equal([]*uti.DataType{uti.JSON}, uti.ByExtension(".json"))
	c.Equal([]*uti.DataType{uti.PlainText}, uti.ByExtension(".txt"))
	c.Equal([]*uti.DataType{uti.PNG}, uti.ByExtension(".png"))

	// Lookup is case-insensitive.
	c.Equal([]*uti.DataType{uti.JSON}, uti.ByExtension(".JSON"))

	// Unknown extensions, and the MIME-type strings the old buggy code mistakenly used as keys, return nil.
	c.Nil(uti.ByExtension(".no-such-extension"))
	c.Nil(uti.ByExtension("application/json"))
}

func TestRegisterUnregisterExtension(t *testing.T) {
	c := check.New(t)

	c.Nil(uti.ByExtension(".uti-test-ext"))

	dt := uti.Register(&uti.DataType{
		UTI:        "test.uti.custom",
		Parents:    []*uti.DataType{uti.Data},
		MimeTypes:  []string{"application/x-uti-test"},
		Extensions: []string{".uti-test-ext"},
	})
	// Register must index the extension (exercises the byExtension build loop).
	c.Equal([]*uti.DataType{dt}, uti.ByExtension(".uti-test-ext"))
	c.Equal([]*uti.DataType{dt}, uti.ByMimeType("application/x-uti-test"))
	c.Equal(dt, uti.ByUTI("test.uti.custom"))

	uti.Unregister(dt)
	// Unregister must remove the extension entry (exercises the byExtension cleanup loop).
	c.Equal(0, len(uti.ByExtension(".uti-test-ext")))
	c.Equal(0, len(uti.ByMimeType("application/x-uti-test")))
	c.Nil(uti.ByUTI("test.uti.custom"))
}
