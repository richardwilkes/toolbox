// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
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

// Split a path into its component parts. In the case of a full path, the first element will be filepath.Separator,
// possibly prefixed by a volume name. In the case of a relative path, the first element will be ".".
func Split(path string) []string {
	var parts []string
	path = filepath.Clean(path)
	parts = append(parts, filepath.Base(path))
	sep := string(filepath.Separator)
	volName := filepath.VolumeName(path)
	path = path[len(volName):]
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
	if volName != "" && result[0] == sep {
		result[0] = volName + sep
	}
	return result
}
