package safe_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/richardwilkes/gokit/xio/fs/safe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbortNonExisting(t *testing.T) {
	filename := "abort.txt"
	f, err := safe.Create(filename)
	require.NoError(t, err)
	n, err := f.WriteString("abort")
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	err = f.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filename)
	assert.Error(t, err)
}

func TestCommitNonExisting(t *testing.T) {
	filename := "commit.txt"
	f, err := safe.Create(filename)
	require.NoError(t, err)
	n, err := f.WriteString("commit")
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	err = f.Commit()
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	_, err = os.Stat(filename)
	assert.NoError(t, err)
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestAbortExisting(t *testing.T) {
	filename := "safe.txt"
	originalData := ([]byte)("safe")
	err := ioutil.WriteFile(filename, originalData, 0600)
	require.NoError(t, err)
	f, err := safe.Create(filename)
	require.NoError(t, err)
	n, err := f.WriteString("bad")
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	err = f.Close()
	assert.NoError(t, err)
	data, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, originalData, data)
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestCommitExisting(t *testing.T) {
	filename := "safe.txt"
	originalData := ([]byte)("safe")
	replacement := ([]byte)("replaced")
	err := ioutil.WriteFile(filename, originalData, 0600)
	require.NoError(t, err)
	f, err := safe.Create(filename)
	require.NoError(t, err)
	n, err := f.Write(replacement)
	assert.NoError(t, err)
	assert.Equal(t, len(replacement), n)
	err = f.Commit()
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	data, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, replacement, data)
	err = os.Remove(filename)
	assert.NoError(t, err)
}
