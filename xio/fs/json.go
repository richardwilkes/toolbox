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
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
)

// LoadJSON data from the specified path.
func LoadJSON(path string, data any) error {
	f, err := os.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return loadJSON(f, path, data)
}

// LoadJSONFromFS data from the specified filesystem path.
func LoadJSONFromFS(fsys fs.FS, path string, data any) error {
	f, err := fsys.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return loadJSON(f, path, data)
}

func loadJSON(r io.ReadCloser, path string, data any) error {
	defer xio.CloseIgnoringErrors(r)
	if err := json.NewDecoder(bufio.NewReader(r)).Decode(data); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}

// SaveJSON data to the specified path.
func SaveJSON(path string, data any, format bool) error {
	return SaveJSONWithMode(path, data, format, 0o644)
}

// SaveJSONWithMode data to the specified path.
func SaveJSONWithMode(path string, data any, format bool, mode os.FileMode) error {
	if err := safe.WriteFileWithMode(path, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		if format {
			encoder.SetIndent("", "  ")
		}
		return errs.Wrap(encoder.Encode(data))
	}, mode); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}
