package rotation_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/log/rotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	maxSize    = 100
	maxBackups = 2
)

func TestRotator(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rotator_test_")
	require.NoError(t, err)
	defer cleanup(t, tmpdir)

	logFiles := []string{filepath.Join(tmpdir, "test.log")}
	for i := 1; i <= maxBackups; i++ {
		logFiles = append(logFiles, fmt.Sprintf("%s-%d", logFiles[0], i))
	}

	r, err := rotation.New(rotation.Path(logFiles[0]), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups))
	require.NoError(t, err)
	_, err = os.Stat(logFiles[0])
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	for i := 0; i < maxSize*(2+maxBackups); i++ {
		_, err = fmt.Fprintln(r, i)
		require.NoError(t, err)
	}
	_, err = fmt.Fprintln(r, "goodbye")
	assert.NoError(t, err)
	require.NoError(t, r.Close())
	for _, f := range logFiles {
		fi, ferr := os.Stat(f)
		require.NoError(t, ferr)
		assert.True(t, fi.Size() <= maxSize)
	}

	r, err = rotation.New(rotation.Path(logFiles[0]), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups))
	require.NoError(t, err)
	_, err = fmt.Fprintln(r, "hello")
	assert.NoError(t, err)
	require.NoError(t, r.Close())
}

func cleanup(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.RemoveAll(path))
}
