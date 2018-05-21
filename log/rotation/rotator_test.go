package rotation_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/richardwilkes/toolbox/log/rotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRotator(t *testing.T) {
	cleanup(t)
	r, err := rotation.New(rotation.Path("test"), rotation.MaxSize(100), rotation.MaxBackups(2))
	require.NoError(t, err)
	_, err = os.Stat("test.log")
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
	for i := 0; i < 300; i++ {
		_, err = fmt.Fprintln(r, i)
		require.NoError(t, err)
	}
	require.NoError(t, r.Close())
	fi, err := os.Stat("test.log")
	require.NoError(t, err)
	assert.True(t, fi.Size() <= 100)
	fi, err = os.Stat("test.1.log")
	require.NoError(t, err)
	assert.True(t, fi.Size() <= 100)
	fi, err = os.Stat("test.2.log")
	require.NoError(t, err)
	assert.True(t, fi.Size() <= 100)
	cleanup(t)
}

func cleanup(t *testing.T) {
	t.Helper()
	for _, f := range []string{"test.log", "test.1.log", "test.2.log"} {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			t.Fatal("Unable to remove " + f)
		}
	}
}
