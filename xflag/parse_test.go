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
	"flag"
	"os"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xflag"
)

// TestParsePostParseOrdering verifies that Parse runs registered post-parse functions ordered by ascending priority,
// and in registration order within a single priority.
func TestParsePostParseOrdering(t *testing.T) {
	c := check.New(t)

	// Parse operates on the global flag.CommandLine and os.Args, so save and restore them.
	savedArgs := os.Args
	savedCommandLine := flag.CommandLine
	defer func() {
		os.Args = savedArgs
		flag.CommandLine = savedCommandLine
	}()
	flag.CommandLine = flag.NewFlagSet("xflag-test", flag.ContinueOnError)
	flagValue := flag.CommandLine.String("name", "default", "")
	os.Args = []string{"xflag-test", "-name", "supplied"}

	var order []string
	xflag.AddPostParseFunc(10, func() { order = append(order, "p10-a") })
	xflag.AddPostParseFunc(-5, func() { order = append(order, "p-5") })
	xflag.AddPostParseFunc(10, func() { order = append(order, "p10-b") })
	xflag.AddPostParseFunc(0, func() { order = append(order, "p0") })

	xflag.Parse()

	// Flag parsing happened before the post-parse functions ran.
	c.Equal("supplied", *flagValue)
	c.Equal([]string{"p-5", "p0", "p10-a", "p10-b"}, order)
}
