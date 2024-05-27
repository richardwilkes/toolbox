// Copyright (c) 2016-2022 by Richard A. Wilkes. All rights reserved.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/log/rotation"
)

const (
	maxSize    = 100
	maxBackups = 2
)

func TestRotator(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "rotator_test_")
	check.NoError(t, err)
	defer cleanup(t, tmpdir)

	logFiles := []string{filepath.Join(tmpdir, "test.log")}
	for i := 1; i <= maxBackups; i++ {
		logFiles = append(logFiles, fmt.Sprintf("%s-%d", logFiles[0], i))
	}

	r, err := rotation.New(rotation.Path(logFiles[0]), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups))
	check.NoError(t, err)
	_, err = os.Stat(logFiles[0])
	check.Error(t, err)
	check.True(t, os.IsNotExist(err))
	for i := 0; i < maxSize*(2+maxBackups); i++ {
		_, err = fmt.Fprintln(r, i)
		check.NoError(t, err)
	}
	_, err = fmt.Fprintln(r, "goodbye")
	check.NoError(t, err)
	check.NoError(t, r.Close())
	for _, f := range logFiles {
		fi, fErr := os.Stat(f)
		check.NoError(t, fErr)
		check.True(t, fi.Size() <= maxSize)
	}

	r, err = rotation.New(rotation.Path(logFiles[0]), rotation.MaxSize(maxSize), rotation.MaxBackups(maxBackups))
	check.NoError(t, err)
	_, err = fmt.Fprintln(r, "hello")
	check.NoError(t, err)
	check.NoError(t, r.Close())
}

func cleanup(t *testing.T, path string) {
	t.Helper()
	check.NoError(t, os.RemoveAll(path))
}
