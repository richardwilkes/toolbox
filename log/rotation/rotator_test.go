package rotation_test

import (
	"fmt"
	"os"
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
	logFiles := []string{"test.log"}
	for i := 1; i <= maxBackups; i++ {
		logFiles = append(logFiles, fmt.Sprintf("%s-%d", logFiles[0], i))
	}

	cleanup(t, logFiles)
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
		fi, err := os.Stat(f)
		require.NoError(t, err)
		assert.True(t, fi.Size() <= maxSize)
	}

	r, err = rotation.New(rotation.Path(logFiles[0]), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups))
	require.NoError(t, err)
	_, err = fmt.Fprintln(r, "hello")
	assert.NoError(t, err)
	require.NoError(t, r.Close())

	cleanup(t, logFiles)
}

func cleanup(t *testing.T, logFiles []string) {
	t.Helper()
	for _, f := range logFiles {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			t.Fatal("Unable to remove " + f)
		}
	}
}
