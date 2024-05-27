/*
 * Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package xio_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/check"
	"github.com/richardwilkes/toolbox/xio"
)

func TestLineWriter(t *testing.T) {
	lines := make([]string, 0)
	w := xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err := w.Write([]byte{})
	check.Equal(t, 0, n)
	check.NoError(t, err)
	check.Equal(t, 0, len(lines))
	n, err = w.Write([]byte{'\n'})
	check.Equal(t, 1, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 1, len(lines))
	check.Equal(t, "", lines[0])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', '\n'})
	check.Equal(t, 2, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 2, len(lines))
	check.Equal(t, "", lines[0])
	check.Equal(t, "", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	check.Equal(t, 3, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 2, len(lines))
	check.Equal(t, "", lines[0])
	check.Equal(t, "a", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	check.Equal(t, 3, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 2, len(lines))
	check.Equal(t, "", lines[0])
	check.Equal(t, "a", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n', '\n'})
	check.Equal(t, 3, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 2, len(lines))
	check.Equal(t, "a", lines[0])
	check.Equal(t, "", lines[1])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n', '\n', 'b'})
	check.Equal(t, 4, n)
	check.NoError(t, err)
	check.NoError(t, w.Close())
	check.Equal(t, 3, len(lines))
	check.Equal(t, "a", lines[0])
	check.Equal(t, "", lines[1])
	check.Equal(t, "b", lines[2])

	lines = make([]string, 0)
	w = xio.NewLineWriter(func(line []byte) {
		lines = append(lines, string(line))
	})
	n, err = w.Write([]byte{'a', '\n'})
	check.Equal(t, 2, n)
	check.NoError(t, err)
	check.Equal(t, 1, len(lines))
	check.Equal(t, "a", lines[0])
	n, err = w.Write([]byte{'\n'})
	check.Equal(t, 1, n)
	check.NoError(t, err)
	check.Equal(t, 2, len(lines))
	check.Equal(t, "a", lines[0])
	check.Equal(t, "", lines[1])
	n, err = w.Write([]byte{'b'})
	check.Equal(t, 1, n)
	check.NoError(t, err)
	check.Equal(t, 2, len(lines))
	check.Equal(t, "a", lines[0])
	check.Equal(t, "", lines[1])
	check.NoError(t, w.Close())
	check.Equal(t, 3, len(lines))
	check.Equal(t, "a", lines[0])
	check.Equal(t, "", lines[1])
	check.Equal(t, "b", lines[2])
}
