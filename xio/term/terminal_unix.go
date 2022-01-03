// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package term

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// IsTerminal returns true if the writer's file descriptor is a terminal.
func IsTerminal(f io.Writer) bool {
	var termios syscall.Termios
	switch v := f.(type) {
	case *os.File:
		_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, v.Fd(), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
		return errno == 0
	default:
		return false
	}
}
