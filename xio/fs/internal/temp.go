// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package internal

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xmath/rand"
)

// CreateTemp is essentially the same as os.CreateTemp, except it allows you to specify the file mode of the newly
// created file. This is here solely because having it in the fs package would cause circular references.
func CreateTemp(dir, pattern string, perm os.FileMode) (*os.File, error) {
	if dir == "" {
		dir = os.TempDir()
	}
	for i := range len(pattern) {
		if os.IsPathSeparator(pattern[i]) {
			return nil, &os.PathError{Op: "createtemp", Path: pattern, Err: errors.New("pattern contains path separator")}
		}
	}
	var prefix, suffix string
	if pos := strings.LastIndexByte(pattern, '*'); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	if dir != "" && os.IsPathSeparator(dir[len(dir)-1]) {
		prefix = dir + prefix
	} else {
		prefix = filepath.Join(dir, prefix)
	}
	try := 0
	for {
		f, err := os.OpenFile(prefix+strconv.Itoa(rand.NewCryptoRand().Intn(math.MaxInt))+suffix,
			os.O_RDWR|os.O_CREATE|os.O_EXCL, perm)
		if os.IsExist(err) {
			try++
			if try < 1000 {
				continue
			}
			return nil, &os.PathError{Op: "createtemp", Path: prefix + "*" + suffix, Err: os.ErrExist}
		}
		return f, err
	}
}
