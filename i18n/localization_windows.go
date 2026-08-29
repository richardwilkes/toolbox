// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
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

// Locale returns the user's default locale, falling back to the system default locale and then to "en_US.UTF-8".
func Locale() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultLocaleName")
	buffer := make([]uint16, 128)
	if ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 { //nolint:errcheck // ret is the error code
		proc = kernel32.NewProc("GetSystemDefaultLocaleName")
		if ret, _, _ = proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); ret == 0 { //nolint:errcheck // ret is the error code
			return "en_US.UTF-8"
		}
	}
	return syscall.UTF16ToString(buffer)
}
