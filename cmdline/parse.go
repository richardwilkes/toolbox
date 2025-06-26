// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import (
	"fmt"

	"github.com/richardwilkes/toolbox/errs"
)

// Parse a command line string into its component parts.
func Parse(command string) ([]string, error) {
	var args []string
	var current []rune
	var lookingForQuote rune
	var escapeNext bool
	for _, ch := range command {
		switch {
		case escapeNext:
			current = append(current, ch)
			escapeNext = false
		case lookingForQuote == ch:
			args = append(args, string(current))
			current = nil
			lookingForQuote = 0
		case lookingForQuote != 0:
			current = append(current, ch)
		case ch == '\\':
			escapeNext = true
		case ch == '"' || ch == '\'':
			lookingForQuote = ch
		case ch == ' ' || ch == '\t':
			if len(current) != 0 {
				args = append(args, string(current))
				current = nil
			}
		default:
			current = append(current, ch)
		}
	}
	if lookingForQuote != 0 {
		return nil, errs.New(fmt.Sprintf("unclosed quote in command line:\n%s", command))
	}
	if len(current) != 0 {
		args = append(args, string(current))
	}
	return args, nil
}
