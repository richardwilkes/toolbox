package fs_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/richardwilkes/gokit/xio/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSaveJSON(t *testing.T) {
	type data struct {
		Name  string
		Count int
	}
	value := &data{
		Name:  "Rich",
		Count: 22,
	}
	f, err := ioutil.TempFile("", "jsontest")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	require.NoError(t, fs.SaveJSON(f.Name(), value, false))
	var value2 data
	require.NoError(t, fs.LoadJSON(f.Name(), &value2))
	require.NoError(t, os.Remove(f.Name()))
	assert.Equal(t, value, &value2)
}
