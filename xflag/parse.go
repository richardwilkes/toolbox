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
	"maps"
	"slices"
)

var postParseFuncs = make(map[int][]func())

// Parse parses the command-line flags from os.Args[1:], then calls any post-parse functions added with
// AddPostParseFunc. Must be called after all flags are defined and before flags are accessed by the program.
func Parse() {
	flag.Parse()
	if len(postParseFuncs) > 0 {
		for _, priority := range slices.Sorted(maps.Keys(postParseFuncs)) {
			for _, f := range postParseFuncs[priority] {
				f()
			}
		}
	}
}

// AddPostParseFunc adds a function to be called after Parse has parsed the flags. Functions run in ascending priority
// order, and in the order they were added within a priority.
func AddPostParseFunc(priority int, f func()) {
	postParseFuncs[priority] = append(postParseFuncs[priority], f)
}
