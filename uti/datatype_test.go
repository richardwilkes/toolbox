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
