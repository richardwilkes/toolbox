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

func TestToSnakeCase(t *testing.T) {
	require.Equal(t, "snake_case", txt.ToSnakeCase("snake_case"))
	require.Equal(t, "camel_case", txt.ToSnakeCase("CamelCase"))
}
