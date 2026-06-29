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
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

//nolint:goconst // The tests are more readable without constants for duplicated string
func TestColorSupport(t *testing.T) {
	for _, tc := range []struct {
		name        string
		colorTerm   string
		termProgram string
		termVersion string
		term        string
		want        Kind
	}{
		{name: "COLORTERM=truecolor", colorTerm: "truecolor", term: "xterm", want: Color24},
		{name: "COLORTERM=24bit", colorTerm: "24bit", term: "xterm", want: Color24},
		{name: "COLORTERM=24bit not downgraded by 256 TERM", colorTerm: "24bit", term: "xterm-256color", want: Color24},
		{name: "COLORTERM hint with 256 TERM stays 256", colorTerm: "yes", term: "xterm-256color", want: Color8},
		{name: "COLORTERM hint with 16 TERM", colorTerm: "yes", term: "xterm", want: Color4},
		{name: "COLORTERM hint with unknown TERM", colorTerm: "1", term: "frobozz", want: Color4},
		{name: "empty COLORTERM is not a hint", colorTerm: "", term: "frobozz", want: Dumb},
		{name: "empty COLORTERM with 256 TERM", colorTerm: "", term: "xterm-256color", want: Color8},
		{name: "empty COLORTERM with 16 TERM", colorTerm: "", term: "xterm", want: Color4},
		{name: "empty COLORTERM and empty TERM", colorTerm: "", term: "", want: Dumb},
		{name: "iTerm.app v3 is truecolor", termProgram: "iTerm.app", termVersion: "3.4.1", term: "frobozz", want: Color24},
		{name: "iTerm.app v2 is 256", termProgram: "iTerm.app", termVersion: "2.9", term: "frobozz", want: Color8},
		{name: "Apple_Terminal is 256", termProgram: "Apple_Terminal", term: "frobozz", want: Color8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := check.New(t)
			t.Setenv("COLORTERM", tc.colorTerm)
			t.Setenv("TERM_PROGRAM", tc.termProgram)
			t.Setenv("TERM_PROGRAM_VERSION", tc.termVersion)
			c.Equal(tc.want, colorSupport(tc.term))
		})
	}
}
