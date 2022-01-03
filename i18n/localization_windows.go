// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package i18n

import (
	"syscall"
	"unsafe"
)

// Locale returns the locale set for the user. If that has not been set, then it falls back to the locale set for the
// system. If that is also unset, then it return "en_US.UTF-8".
func Locale() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")
	buffer := make([]uint16, 128)
	if ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 {
		proc = kernel32.NewProc("GetSystemDefaultLocaleName")
		if ret, _, _ = proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 {
			return "en_US.UTF-8"
		}
	}
	return syscall.UTF16ToString(buffer)
}
