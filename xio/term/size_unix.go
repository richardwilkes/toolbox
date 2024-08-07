// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package term

import (
	"syscall"
	"unsafe"
)

// Size returns the number of columns and rows comprising the terminal.
func Size() (columns, rows int) {
	var ws [4]uint16
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws[0]))); errno == 0 {
		return int(ws[1]), int(ws[0])
	}
	return defColumns, defRows
}
