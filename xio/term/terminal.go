// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package term provides terminal utilities.
package term

import (
	"fmt"
	"io"
	"strings"
)

const (
	defColumns = 80
	defRows    = 24
)

// WrapText prints the 'prefix' to 'out' and then wraps 'text' in the remaining space.
func WrapText(out io.Writer, prefix, text string) {
	fmt.Fprint(out, prefix)
	avail, _ := Size()
	avail -= 1 + len(prefix)
	if avail < 1 {
		avail = 1
	}
	remaining := avail
	indent := strings.Repeat(" ", len(prefix))
	for _, line := range strings.Split(text, "\n") {
		for _, ch := range line {
			if ch == ' ' {
				fmt.Fprint(out, " ")
				remaining--
			} else {
				break
			}
		}
		for i, token := range strings.Fields(line) {
			length := len(token) + 1
			if i != 0 {
				if length > remaining {
					fmt.Fprintln(out)
					fmt.Fprint(out, indent)
					remaining = avail
				} else {
					fmt.Fprint(out, " ")
				}
			}
			fmt.Fprint(out, token)
			remaining -= length
		}
		fmt.Fprintln(out)
	}
}
