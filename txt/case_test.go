// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/txt"
)

func TestToCamelCase(t *testing.T) {
	c := check.New(t)
	c.Equal("SnakeCase", txt.ToCamelCase("snake_case"))
	c.Equal("SnakeCase", txt.ToCamelCase("snake__case"))
	c.Equal("CamelCase", txt.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	c := check.New(t)
	c.Equal("ID", txt.ToCamelCaseWithExceptions("id", txt.StdAllCaps))
	c.Equal("世界ID", txt.ToCamelCaseWithExceptions("世界_id", txt.StdAllCaps))
	c.Equal("OneID", txt.ToCamelCaseWithExceptions("one_id", txt.StdAllCaps))
	c.Equal("IDOne", txt.ToCamelCaseWithExceptions("id_one", txt.StdAllCaps))
	c.Equal("OneIDTwo", txt.ToCamelCaseWithExceptions("one_id_two", txt.StdAllCaps))
	c.Equal("OneIDTwoID", txt.ToCamelCaseWithExceptions("one_id_two_id", txt.StdAllCaps))
	c.Equal("OneIDID", txt.ToCamelCaseWithExceptions("one_id_id", txt.StdAllCaps))
	c.Equal("Orchid", txt.ToCamelCaseWithExceptions("orchid", txt.StdAllCaps))
	c.Equal("OneURLTwo", txt.ToCamelCaseWithExceptions("one_url_two", txt.StdAllCaps))
	c.Equal("URLID", txt.ToCamelCaseWithExceptions("url_id", txt.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	c := check.New(t)
	c.Equal("snake_case", txt.ToSnakeCase("snake_case"))
	c.Equal("camel_case", txt.ToSnakeCase("CamelCase"))
}
