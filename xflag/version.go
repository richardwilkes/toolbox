// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xflag

import (
	"flag"
	"fmt"

	"github.com/richardwilkes/toolbox/v2/i18n"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// VersionFlagPriority is the priority used when adding the automatic post-parse function for the version flags.
const VersionFlagPriority = -10

var (
	shortVersionFlag bool
	longVersionFlag  bool
)

// AddVersionFlags adds flags for showing the short or long version information to flag.CommandLine. Adds a post-parse
// function that will print the version information and exit if either of the flags is set. Note that this automatic
// handling only works if xflag.Parse is used and not flag.Parse directly.
func AddVersionFlags() {
	flag.BoolVar(&shortVersionFlag, "v", false, i18n.Text("Show the short version and exit"))
	flag.BoolVar(&longVersionFlag, "version", false, i18n.Text("Show the full version and exit"))
	AddPostParseFunc(VersionFlagPriority, func() {
		if longVersionFlag {
			fmt.Println(xos.LongAppVersion())
			xos.Exit(0)
		}
		if shortVersionFlag {
			fmt.Println(xos.ShortAppVersion())
			xos.Exit(0)
		}
	})
}
