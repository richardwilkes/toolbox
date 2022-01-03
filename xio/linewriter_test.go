// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/richardwilkes/toolbox/xio"
)

func TestLineWriter(t *testing.T) {
	lines := make([]string, 0)
	w := xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err := w.Write([]byte{})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(lines))
	n, err = w.Write([]byte{'\n'})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 1, len(lines))
	assert.Equal(t, "", lines[0])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', '\n'})
	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "", lines[0])
	assert.Equal(t, "", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "", lines[0])
	assert.Equal(t, "a", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "", lines[0])
	assert.Equal(t, "a", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n', '\n'})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "a", lines[0])
	assert.Equal(t, "", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n', '\n', 'b'})
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "a", lines[0])
	assert.Equal(t, "", lines[1])
	assert.Equal(t, "b", lines[2])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n'})
	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(lines))
	assert.Equal(t, "a", lines[0])
	n, err = w.Write([]byte{'\n'})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "a", lines[0])
	assert.Equal(t, "", lines[1])
	n, err = w.Write([]byte{'b'})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "a", lines[0])
	assert.Equal(t, "", lines[1])
	assert.NoError(t, w.Close())
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "a", lines[0])
	assert.Equal(t, "", lines[1])
	assert.Equal(t, "b", lines[2])
}
