// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes_test

import (
	"io/fs"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestLineWriter(t *testing.T) {
	lines := make([]string, 0)
	appender := func(line []byte) { lines = append(lines, string(line)) }
	w := xbytes.NewLineWriter(appender)
	n, err := w.Write([]byte{})
	c := check.New(t)
	c.Equal(0, n)
	c.NoError(err)
	c.Equal(0, len(lines))
	n, err = w.Write([]byte{'\n'})
	c.Equal(1, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(1, len(lines))
	c.Equal("", lines[0])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'\n', '\n'})
	c.Equal(2, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(2, len(lines))
	c.Equal("", lines[0])
	c.Equal("", lines[1])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	c.Equal(3, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(2, len(lines))
	c.Equal("", lines[0])
	c.Equal("a", lines[1])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'\n', 'a', '\n'})
	c.Equal(3, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(2, len(lines))
	c.Equal("", lines[0])
	c.Equal("a", lines[1])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'a', '\n', '\n'})
	c.Equal(3, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(2, len(lines))
	c.Equal("a", lines[0])
	c.Equal("", lines[1])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'a', '\n', '\n', 'b'})
	c.Equal(4, n)
	c.NoError(err)
	c.NoError(w.Close())
	c.Equal(3, len(lines))
	c.Equal("a", lines[0])
	c.Equal("", lines[1])
	c.Equal("b", lines[2])

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.Write([]byte{'a', '\n'})
	c.Equal(2, n)
	c.NoError(err)
	c.Equal(1, len(lines))
	c.Equal("a", lines[0])
	n, err = w.Write([]byte{'\n'})
	c.Equal(1, n)
	c.NoError(err)
	c.Equal(2, len(lines))
	c.Equal("a", lines[0])
	c.Equal("", lines[1])
	n, err = w.Write([]byte{'b'})
	c.Equal(1, n)
	c.NoError(err)
	c.Equal(2, len(lines))
	c.Equal("a", lines[0])
	c.Equal("", lines[1])
	c.NoError(w.Close())
	c.Equal(3, len(lines))
	c.Equal("a", lines[0])
	c.Equal("", lines[1])
	c.Equal("b", lines[2])

	n, err = w.Write([]byte{'c'})
	c.HasError(err)
	c.Equal(0, n)

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	n, err = w.WriteString("")
	c.Equal(0, n)
	c.NoError(err)
	c.Equal(0, len(lines))
	n, err = w.WriteString("\n")
	c.Equal(1, n)
	c.NoError(err)
	c.Equal(1, len(lines))
	c.Equal("", lines[0])
	lines = make([]string, 0)
	n, err = w.WriteString("\n\n")
	c.Equal(2, n)
	c.NoError(err)
	c.Equal(2, len(lines))
	c.Equal("", lines[0])
	c.Equal("", lines[1])

	lines = make([]string, 0)
	n, err = w.WriteString("abc\ndef\n")
	c.Equal(8, n)
	c.NoError(err)
	c.Equal(2, len(lines))
	c.Equal("abc", lines[0])
	c.Equal("def", lines[1])

	lines = make([]string, 0)
	n, err = w.WriteString("hello")
	c.Equal(5, n)
	c.NoError(err)
	c.Equal(0, len(lines))
	c.NoError(w.Close())
	c.Equal(1, len(lines))
	c.Equal("hello", lines[0])

	n, err = w.WriteString("test")
	c.Equal(0, n)
	c.Equal(fs.ErrClosed, err)

	lines = make([]string, 0)
	w = xbytes.NewLineWriter(appender)
	err = w.WriteByte('a')
	c.NoError(err)
	c.Equal(0, len(lines))

	err = w.WriteByte('\n')
	c.NoError(err)
	c.Equal(1, len(lines))
	c.Equal("a", lines[0])

	err = w.WriteByte('b')
	c.NoError(err)
	err = w.WriteByte('c')
	c.NoError(err)
	err = w.WriteByte('\n')
	c.NoError(err)
	c.Equal(2, len(lines))
	c.Equal("bc", lines[1])

	c.NoError(w.Close())
	err = w.WriteByte('x')
	c.Equal(fs.ErrClosed, err)
}
