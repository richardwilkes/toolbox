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
	"github.com/richardwilkes/toolbox/v2/errs"
)

// SplitCommandLine splits a command line string into its component parts.
func SplitCommandLine(command string) ([]string, error) {
	var args []string
	var lookingForQuote rune
	var escapeNext bool
	current := make([]rune, 0, len(command))
	for _, ch := range command {
		switch {
		case escapeNext:
			current = append(current, ch)
			escapeNext = false
		case ch == '\\':
			escapeNext = true
		case lookingForQuote == ch:
			args = append(args, string(current))
			current = current[:0]
			lookingForQuote = 0
		case lookingForQuote != 0:
			current = append(current, ch)
		case ch == '"' || ch == '\'':
			lookingForQuote = ch
		case ch == ' ' || ch == '\t':
			if len(current) != 0 {
				args = append(args, string(current))
				current = current[:0]
			}
		default:
			current = append(current, ch)
		}
	}
	if escapeNext {
		return nil, errs.Newf("escape at end of command line:\n%s", command)
	}
	if lookingForQuote != 0 {
		return nil, errs.Newf("unclosed quote in command line:\n%s", command)
	}
	if len(current) != 0 {
		args = append(args, string(current))
	}
	return args, nil
}

// SplitCommandLineWithoutEscapes splits a command line string into its component parts without considering escape
// sequences.
func SplitCommandLineWithoutEscapes(command string) ([]string, error) {
	var args []string
	var lookingForQuote rune
	current := make([]rune, 0, len(command))
	for _, ch := range command {
		switch {
		case lookingForQuote == ch:
			args = append(args, string(current))
			current = current[:0]
			lookingForQuote = 0
		case lookingForQuote != 0:
			current = append(current, ch)
		case ch == '"' || ch == '\'':
			lookingForQuote = ch
		case ch == ' ' || ch == '\t':
			if len(current) != 0 {
				args = append(args, string(current))
				current = current[:0]
			}
		default:
			current = append(current, ch)
		}
	}
	if lookingForQuote != 0 {
		return nil, errs.Newf("unclosed quote in command line:\n%s", command)
	}
	if len(current) != 0 {
		args = append(args, string(current))
	}
	return args, nil
}
