// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	table := []struct {
		Prefix string
		Text   string
		Max    int
		Out    string
	}{
		{Prefix: "// ", Text: "short", Max: 78, Out: "// short"},
		{Prefix: "// ", Text: "some text that is longer", Max: 12, Out: "// some text\n// that is\n// longer"},
		{Prefix: "", Text: "some text that is longer", Max: 12, Out: "some text\nthat is\nlonger"},
		{Prefix: "", Text: "some text that is longer", Max: 4, Out: "some\ntext\nthat\nis\nlonger"},
		{Prefix: "", Text: "some text that is longer, yep", Max: 4, Out: "some\ntext\nthat\nis\nlonger,\nyep"},
	}
	for i, one := range table {
		assert.Equal(t, one.Out, txt.Wrap(one.Prefix, one.Text, one.Max), "#%d", i)
	}
}
