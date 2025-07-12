// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xjson_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xjson"
)

func TestLoadSaveJSON(t *testing.T) {
	type data struct {
		Name  string
		Count int
	}
	value := &data{
		Name:  "Ziggy",
		Count: 12345,
	}
	c := check.New(t)
	p := filepath.Join(t.TempDir(), "test1.json")
	c.NoError(xjson.Save(p, value, false))
	var value2 data
	c.NoError(xjson.Load(p, &value2))
	c.Equal(value, &value2)
	p = filepath.Join(t.TempDir(), "test2.json")
	c.NoError(xjson.Save(p, value, true))
	var value3 data
	c.NoError(xjson.Load(p, &value3))
	c.Equal(value, &value3)
	var value4 data
	c.NoError(xjson.LoadFS(os.DirFS(filepath.Dir(p)), "test2.json", &value4))
	c.Equal(value, &value4)
}
