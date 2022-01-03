// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//go:build !windows

package term

import (
	"github.com/pkg/term"
)

// Read a byte from the terminal.
func Read() (ch byte, ok bool) {
	t, err := term.Open("/dev/tty")
	if err != nil {
		return 0, false
	}
	err = term.RawMode(t)
	if err != nil {
		return 0, false
	}
	bytes := make([]byte, 1)
	numRead, err := t.Read(bytes)
	if altErr := t.Restore(); altErr != nil && err == nil {
		err = altErr
	}
	if altErr := t.Close(); altErr != nil && err == nil {
		err = altErr
	}
	if err != nil || numRead == 0 {
		return 0, false
	}
	return bytes[0], true
}
