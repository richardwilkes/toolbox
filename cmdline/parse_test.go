// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/cmdline"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"hello world", []string{"hello", "world"}},
		{`hello "world hello"`, []string{"hello", "world hello"}},
		{`'hello again' "world hello"`, []string{"hello again", "world hello"}},
		{`\"hello\ world\"`, []string{`"hello world"`}},
		{"hello 世界", []string{"hello", "世界"}},
		{`hello\ world`, []string{"hello world"}},
	}
	for i, one := range tests {
		parts, err := cmdline.Parse(one.input)
		check.NoError(t, err, i)
		check.Equal(t, one.expected, parts, i)
	}
}
