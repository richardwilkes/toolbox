// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs

import "os"

// FileIsReadable returns true if the path points to a regular file that we have permission to read.
func FileIsReadable(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		mode := fi.Mode()
		return !mode.IsDir() && mode.IsRegular() && mode.Perm()&0o400 != 0
	}
	return false
}
