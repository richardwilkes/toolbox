// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package toolbox

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/richardwilkes/toolbox/errs"
)

// AppDir returns the logical directory the application resides within. For macOS, this means the directory where the
// .app bundle resides, not the binary that's tucked down inside it.
func AppDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", errs.Wrap(err)
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return "", errs.Wrap(err)
	}
	if path, err = filepath.Abs(path); err != nil {
		return "", errs.Wrap(err)
	}
	path = filepath.Dir(path)
	if runtime.GOOS == MacOS {
		// Account for macOS bundles
		if i := strings.LastIndex(path, ".app/"); i != -1 {
			path = filepath.Dir(path[:i])
		}
	}
	return path, nil
}
