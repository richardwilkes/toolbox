// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

// AddVersionFlags adds -v and -version flags to flag.CommandLine, along with a post-parse function that prints the
// short or full version and exits if either is set. This only works when xflag.Parse is used rather than flag.Parse.
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
