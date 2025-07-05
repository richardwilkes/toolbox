// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xio/fs"
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
	c := check.New(t)
	c.NoError(err)
	c.NoError(f.Close())
	c.NoError(fs.SaveYAMLWithMode(f.Name(), value, 0o600))
	var value2 data
	c.NoError(fs.LoadYAML(f.Name(), &value2))
	c.NoError(os.Remove(f.Name()))
	c.Equal(value, &value2)
}
