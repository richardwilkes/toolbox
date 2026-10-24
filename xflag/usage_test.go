// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xflag_test

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xflag"
)

// usageOutput defines a flag set with the given flags, applies SetUsage with hiddenFlags, and returns the text the
// resulting Usage function emits. Writing to a bytes.Buffer is detected as a dumb terminal, so no escape codes appear.
func usageOutput(t *testing.T, hiddenFlags ...string) string {
	t.Helper()
	fs := flag.NewFlagSet("xflag-test", flag.ContinueOnError)
	var buffer bytes.Buffer
	fs.SetOutput(&buffer)
	fs.String("visible", "", "A `path` that should be listed")
	fs.Bool("secret", false, "Should not be listed")
	fs.Int("internal", 0, "Should not be listed either")
	xflag.SetUsage(fs, "Test description.", "<arg>", hiddenFlags...)
	fs.Usage()
	return buffer.String()
}

func TestSetUsageHidesNamedFlags(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		hidden   []string
		contains []string
		omits    []string
	}{
		{contains: []string{"[options]", "Options:", "-visible <path>", "-secret", "-internal"}},
		{
			hidden:   []string{"secret", "internal", "does-not-exist"},
			contains: []string{"[options]", "Options:", "-visible <path>"},
			omits:    []string{"-secret", "-internal", "Should not be listed"},
		},
		{
			hidden:   []string{"visible", "secret", "internal"},
			contains: []string{"Test description.", "<arg>"},
			omits:    []string{"[options]", "Options:", "-visible", "-secret", "-internal"},
		},
	} {
		out := usageOutput(t, tc.hidden...)
		for _, s := range tc.contains {
			c.Contains(out, s, "hidden %q", tc.hidden)
		}
		for _, s := range tc.omits {
			c.NotContains(out, s, "hidden %q", tc.hidden)
		}
	}
}

func TestSetUsageHiddenFlagsStillParse(t *testing.T) {
	c := check.New(t)
	fs := flag.NewFlagSet("xflag-test", flag.ContinueOnError)
	fs.SetOutput(&strings.Builder{})
	secret := fs.Bool("secret", false, "Should not be listed")
	xflag.SetUsage(fs, "", "", "secret")
	c.NoError(fs.Parse([]string{"-secret"}))
	c.True(*secret)
}
