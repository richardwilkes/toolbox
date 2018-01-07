package fs_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	f, err := ioutil.TempFile("", "yamltest")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	require.NoError(t, fs.SaveYAML(f.Name(), value))
	var value2 data
	require.NoError(t, fs.LoadYAML(f.Name(), &value2))
	require.NoError(t, os.Remove(f.Name()))
	assert.Equal(t, value, &value2)
}
