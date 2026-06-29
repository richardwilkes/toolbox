// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//go:build !windows

package xterm

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	term256Matcher = regexp.MustCompile("(?i)-256(color)?$")
	term16Matcher  = regexp.MustCompile("(?i)^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux")
)

func enableColor() bool {
	return true
}

func colorSupport(envTerm string) Kind {
	// "truecolor" and "24bit" are the two de-facto values terminals advertise in COLORTERM for 24-bit support.
	envColorTerm := os.Getenv("COLORTERM")
	if envColorTerm == "truecolor" || envColorTerm == "24bit" { //nolint:goconst // Only duplicated in test code
		return Color24
	}
	if exe, ok := os.LookupEnv("TERM_PROGRAM"); ok {
		switch exe {
		case "iTerm.app": //nolint:goconst // Only duplicated in test code
			if version, err := strconv.Atoi(strings.Split(os.Getenv("TERM_PROGRAM_VERSION"), ".")[0]); err == nil && version >= 3 {
				return Color24
			}
			return Color8
		case "Apple_Terminal":
			return Color8
		}
	}
	if term256Matcher.MatchString(envTerm) {
		return Color8
	}
	// A non-empty COLORTERM (any other value) is a hint of at least basic color support. An empty value carries no
	// information, so it is not treated as a hint.
	if envColorTerm != "" || term16Matcher.MatchString(envTerm) {
		return Color4
	}
	return Dumb
}
