// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xruntime

import "unsafe"

// PtrFromUintptr converts a uintptr value into a *T. Unlike a direct (*T)(unsafe.Pointer(v)) conversion — which
// checkptr instrumentation (enabled by -race) treats as pointer arithmetic with no known source allocation, aborting
// whenever the result lands inside a Go heap object — reading the pointer out of the address of a local involves no
// instrumented conversion, so it stays valid under checkptr. That matters when the value round-trips one of our own
// pinned Go objects through the OS, as COM does. The caller remains responsible for ensuring the value designates
// memory that is still live and, when it refers to Go memory, pinned or otherwise immovable.
func PtrFromUintptr[T any, U ~uintptr](v U) *T {
	return *(**T)(unsafe.Pointer(&v))
}
