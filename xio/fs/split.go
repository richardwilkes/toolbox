// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs

import (
	"path/filepath"
)

// Split a path into its component parts.
func Split(path string) []string {
	var parts []string
	path = filepath.Clean(path)
	parts = append(parts, filepath.Base(path))
	sep := string(filepath.Separator)
	for {
		path = filepath.Dir(path)
		parts = append(parts, filepath.Base(path))
		if path == "." || path == sep {
			break
		}
	}
	result := make([]string, len(parts))
	for i := 0; i < len(parts); i++ {
		result[len(parts)-(i+1)] = parts[i]
	}
	return result
}
