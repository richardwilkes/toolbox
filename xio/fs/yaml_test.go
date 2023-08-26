// Copyright ©2016-2023 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package fs_test

import (
	"os"
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xio/fs"
)

func TestLoadSaveYAML(t *testing.T) {
	type data struct {
		Name  string
		Count int
	}
	value := &data{
		Name:  "Rich",
		Count: 22,
	}
	f, err := os.CreateTemp("", "yaml_test")
	check.NoError(t, err)
	check.NoError(t, f.Close())
	check.NoError(t, fs.SaveYAMLWithMode(f.Name(), value, 0o600))
	var value2 data
	check.NoError(t, fs.LoadYAML(f.Name(), &value2))
	check.NoError(t, os.Remove(f.Name()))
	check.Equal(t, value, &value2)
}
