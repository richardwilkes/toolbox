// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xfilepath

import (
	"path/filepath"
	"strings"
)

// Split a path into its component parts. In the case of a full path, the first element will be filepath.Separator,
// possibly prefixed by a volume name. In the case of a relative path, the first element will be ".".
func Split(path string) []string {
	var parts []string
	path = filepath.Clean(path)
	parts = append(parts, filepath.Base(path))
	sep := string(filepath.Separator)
	volName := filepath.VolumeName(path)
	if volName != "" && strings.Trim(volName, sep) == "" {
		// On Windows, filepath.Clean() reduces a path made up of nothing but separators (e.g. "//") to `\\`, which
		// filepath.VolumeName() then reports as a volume name even though it names neither a host nor a share. It
		// isn't a real volume, so discard it and treat what remains as a plain root.
		volName = ""
		path = sep
	} else {
		path = path[len(volName):]
		if path == "" {
			// The path was nothing but a volume name (e.g. `\\host\share`), which denotes that volume's root.
			path = sep
		}
	}
	for path != "." && path != sep {
		path = filepath.Dir(path)
		parts = append(parts, filepath.Base(path))
	}
	result := make([]string, len(parts))
	for i := range parts {
		result[len(parts)-(i+1)] = parts[i]
	}
	if volName != "" && result[0] == sep {
		result[0] = volName + sep
	}
	return result
}
