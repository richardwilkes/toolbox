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
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

// LoadJSON data from the specified path.
func LoadJSON(path string, data interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(f)
	return errs.Wrap(json.NewDecoder(bufio.NewReader(f)).Decode(data))
}

// SaveJSON data to the specified path.
func SaveJSON(path string, data interface{}, format bool) error {
	return SaveJSONWithMode(path, data, format, 0644) //nolint:gocritic // File modes are octal
}

// SaveJSONWithMode data to the specified path.
func SaveJSONWithMode(path string, data interface{}, format bool, mode os.FileMode) error {
	return safe.WriteFileWithMode(path, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		if format {
			encoder.SetIndent("", "  ")
		}
		return errs.Wrap(encoder.Encode(data))
	}, mode)
}
