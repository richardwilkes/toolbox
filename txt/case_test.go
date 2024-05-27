/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package txt_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/txt"
)

func TestToCamelCase(t *testing.T) {
	check.Equal(t, "SnakeCase", txt.ToCamelCase("snake_case"))
	check.Equal(t, "SnakeCase", txt.ToCamelCase("snake__case"))
	check.Equal(t, "CamelCase", txt.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	check.Equal(t, "ID", txt.ToCamelCaseWithExceptions("id", txt.StdAllCaps))
	check.Equal(t, "世界ID", txt.ToCamelCaseWithExceptions("世界_id", txt.StdAllCaps))
	check.Equal(t, "OneID", txt.ToCamelCaseWithExceptions("one_id", txt.StdAllCaps))
	check.Equal(t, "IDOne", txt.ToCamelCaseWithExceptions("id_one", txt.StdAllCaps))
	check.Equal(t, "OneIDTwo", txt.ToCamelCaseWithExceptions("one_id_two", txt.StdAllCaps))
	check.Equal(t, "OneIDTwoID", txt.ToCamelCaseWithExceptions("one_id_two_id", txt.StdAllCaps))
	check.Equal(t, "OneIDID", txt.ToCamelCaseWithExceptions("one_id_id", txt.StdAllCaps))
	check.Equal(t, "Orchid", txt.ToCamelCaseWithExceptions("orchid", txt.StdAllCaps))
	check.Equal(t, "OneURLTwo", txt.ToCamelCaseWithExceptions("one_url_two", txt.StdAllCaps))
	check.Equal(t, "URLID", txt.ToCamelCaseWithExceptions("url_id", txt.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	check.Equal(t, "snake_case", txt.ToSnakeCase("snake_case"))
	check.Equal(t, "camel_case", txt.ToSnakeCase("CamelCase"))
}
