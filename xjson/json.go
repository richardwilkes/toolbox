// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xjson

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// Load JSON data from the specified path.
func Load(path string, data any) error {
	f, err := os.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return load(f, path, data)
}

// LoadFS loads JSON data from the specified filesystem path.
func LoadFS(fsys fs.FS, path string, data any) error {
	f, err := fsys.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return load(f, path, data)
}

func load(r io.ReadCloser, path string, data any) error {
	defer xio.CloseIgnoringErrors(r)
	if err := json.NewDecoder(bufio.NewReader(r)).Decode(data); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}

// Save JSON data to the specified path.
func Save(path string, data any, format bool) error {
	if err := xos.WriteSafeFile(path, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		if format {
			encoder.SetIndent("", "  ")
		}
		return errs.Wrap(encoder.Encode(data))
	}); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}
