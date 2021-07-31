// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
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
	"strings"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/yookoala/realpath"
)

// UniquePaths returns a list of unique paths from the given paths, pruning
// out paths that are a subset of another.
func UniquePaths(paths ...string) ([]string, error) {
	set := make(map[string]bool)
	for _, path := range paths {
		actual, err := realpath.Realpath(path)
		if err != nil {
			return nil, errs.NewWithCause(path, err)
		}
		if _, exists := set[actual]; !exists {
			add := true
			for one := range set {
				var p1, p2 string
				if p1, err = filepath.Rel(one, actual); err != nil {
					return nil, errs.NewWithCause(path, err)
				}
				if p2, err = filepath.Rel(actual, one); err != nil {
					return nil, errs.NewWithCause(path, err)
				}
				prefixed := strings.HasPrefix(p1, "..")
				if prefixed != strings.HasPrefix(p2, "..") {
					if prefixed {
						delete(set, one)
					} else {
						add = false
						break
					}
				}
			}
			if add {
				set[actual] = true
			}
		}
	}
	result := make([]string, 0, len(set))
	for p := range set {
		result = append(result, p)
	}
	return result, nil
}
