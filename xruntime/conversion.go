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

// PtrFromUintptr converts a uintptr value to a pointer of the specified type. The type parameter U is constrained to be
// a uintptr or a type that is based on uintptr. This function is useful for converting uintptr values obtained from
// unsafe operations back into pointers of the desired type without the linter complaining.
func PtrFromUintptr[T any, U ~uintptr](v U) *T {
	return (*T)(unsafe.Pointer(v))
}
