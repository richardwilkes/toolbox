// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xos

// Must is a helper function that takes a value of any type and an error. If the error is nil, it returns the value; if
// the error is non-nil, it panics.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
