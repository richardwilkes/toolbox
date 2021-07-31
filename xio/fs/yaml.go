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
	"io"
	"os"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/safe"

	"gopkg.in/yaml.v2"
)

// LoadYAML data from the specified path.
func LoadYAML(path string, data interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(f)
	return errs.Wrap(yaml.NewDecoder(bufio.NewReader(f)).Decode(data))
}

// SaveYAML data to the specified path.
func SaveYAML(path string, data interface{}) error {
	return SaveYAMLWithMode(path, data, 0644) //nolint:gocritic // File modes are octal
}

// SaveYAMLWithMode data to the specified path.
func SaveYAMLWithMode(path string, data interface{}, mode os.FileMode) error {
	return safe.WriteFileWithMode(path, func(w io.Writer) error {
		encoder := yaml.NewEncoder(w)
		if err := encoder.Encode(data); err != nil {
			return errs.Wrap(err)
		}
		return errs.Wrap(encoder.Close())
	}, mode)
}
