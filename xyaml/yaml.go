// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xyaml

import (
	"bufio"
	"io"
	"io/fs"
	"os"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio"
	"github.com/richardwilkes/toolbox/v2/xos"

	"gopkg.in/yaml.v3"
)

// Load YAML data from the specified path.
func Load(path string, data any) error {
	f, err := os.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return load(f, path, data)
}

// LoadFS YAML data from the specified filesystem path.
func LoadFS(fsys fs.FS, path string, data any) error {
	f, err := fsys.Open(path)
	if err != nil {
		return errs.NewWithCause(path, err)
	}
	return load(f, path, data)
}

func load(r io.ReadCloser, path string, data any) error {
	defer xio.CloseIgnoringErrors(r)
	if err := yaml.NewDecoder(bufio.NewReader(r)).Decode(data); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}

// Save YAML data to the specified path. This will use xos.WriteSafeFile so that a failure does not overwrite any
// original file that may have been present.
func Save(path string, data any) error {
	if err := xos.WriteSafeFile(path, func(w io.Writer) error {
		encoder := yaml.NewEncoder(w)
		encoder.SetIndent(2)
		if err := encoder.Encode(data); err != nil {
			return errs.Wrap(err)
		}
		return errs.Wrap(encoder.Close())
	}); err != nil {
		return errs.NewWithCause(path, err)
	}
	return nil
}
