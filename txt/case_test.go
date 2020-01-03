// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
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

	"github.com/richardwilkes/toolbox/txt"
	"github.com/stretchr/testify/require"
)

func TestToCamelCase(t *testing.T) {
	require.Equal(t, "SnakeCase", txt.ToCamelCase("snake_case"))
	require.Equal(t, "SnakeCase", txt.ToCamelCase("snake__case"))
	require.Equal(t, "CamelCase", txt.ToCamelCase("CamelCase"))
}

func TestToCamelCaseWithExceptions(t *testing.T) {
	require.Equal(t, "ID", txt.ToCamelCaseWithExceptions("id", txt.StdAllCaps))
	require.Equal(t, "世界ID", txt.ToCamelCaseWithExceptions("世界_id", txt.StdAllCaps))
	require.Equal(t, "OneID", txt.ToCamelCaseWithExceptions("one_id", txt.StdAllCaps))
	require.Equal(t, "IDOne", txt.ToCamelCaseWithExceptions("id_one", txt.StdAllCaps))
	require.Equal(t, "OneIDTwo", txt.ToCamelCaseWithExceptions("one_id_two", txt.StdAllCaps))
	require.Equal(t, "OneIDTwoID", txt.ToCamelCaseWithExceptions("one_id_two_id", txt.StdAllCaps))
	require.Equal(t, "OneIDID", txt.ToCamelCaseWithExceptions("one_id_id", txt.StdAllCaps))
	require.Equal(t, "Orchid", txt.ToCamelCaseWithExceptions("orchid", txt.StdAllCaps))
	require.Equal(t, "OneURLTwo", txt.ToCamelCaseWithExceptions("one_url_two", txt.StdAllCaps))
	require.Equal(t, "URLID", txt.ToCamelCaseWithExceptions("url_id", txt.StdAllCaps))
}

func TestToSnakeCase(t *testing.T) {
	require.Equal(t, "snake_case", txt.ToSnakeCase("snake_case"))
	require.Equal(t, "camel_case", txt.ToSnakeCase("CamelCase"))
}
