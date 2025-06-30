// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs

import (
	"os"

	"github.com/richardwilkes/toolbox/v2/xio/fs/internal"
)

// CreateTemp is essentially the same as os.CreateTemp, except it allows you to specify the file mode of the newly
// created file.
func CreateTemp(dir, pattern string, perm os.FileMode) (*os.File, error) {
	return internal.CreateTemp(dir, pattern, perm)
}
