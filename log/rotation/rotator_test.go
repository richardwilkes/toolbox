// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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
		fi, fErr := os.Stat(f)
		require.NoError(t, fErr)
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
