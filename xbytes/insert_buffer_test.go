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
	"bytes"
	"io"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xbytes"
)

func TestInsertBuffer_Basic(t *testing.T) {
	c := check.New(t)
	var bp *xbytes.InsertBuffer
	c.Equal("<nil>", bp.String())
	var b xbytes.InsertBuffer
	c.Equal(0, b.Len())
	c.True(b.Cap() >= 0)
	c.Equal(b.Cap(), b.Available())
	c.NotEqual("<nil>", b.String())
	_, err := b.Write([]byte("abc"))
	c.NoError(err)
	c.Equal("abc", string(b.Bytes()))
	b.Reset()
	c.Equal(0, b.Len())
}

func TestInsertBuffer_Insert(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("ace")
	c.NoError(err)
	err = b.Insert(1, []byte("b"))
	c.NoError(err)
	c.Equal("abce", string(b.Bytes()))
	err = b.Insert(-1, []byte("x"))
	c.HasError(err)
	err = b.Insert(100, []byte("x"))
	c.HasError(err)
}

func TestInsertBuffer_InsertByte(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("ac")
	c.NoError(err)
	err = b.InsertByte(1, 'b')
	c.NoError(err)
	c.Equal("abc", string(b.Bytes()))
	err = b.InsertByte(-1, 'x')
	c.HasError(err)
	err = b.InsertByte(100, 'x')
	c.HasError(err)
}

func TestInsertBuffer_InsertRune(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("a")
	c.NoError(err)
	err = b.InsertRune(1, 'ß')
	c.NoError(err)
	c.Equal("aß", string(b.Bytes()))
	err = b.InsertRune(1, 'c')
	c.NoError(err)
	c.Equal("acß", string(b.Bytes()))
}

func TestInsertBuffer_InsertString(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("ac")
	c.NoError(err)
	err = b.InsertString(1, "b")
	c.NoError(err)
	c.Equal("abc", string(b.Bytes()))
}

func TestInsertBuffer_WriteMethods(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	n, err := b.Write([]byte("foo"))
	c.NoError(err)
	c.Equal(3, n)
	err = b.WriteByte('b')
	c.NoError(err)
	n, err = b.WriteRune('ß')
	c.NoError(err)
	c.Equal(2, n)
	n, err = b.WriteRune('c')
	c.NoError(err)
	c.Equal(1, n)
	n, err = b.WriteString("ar")
	c.NoError(err)
	c.Equal(2, n)
	c.Equal("foobßcar", string(b.Bytes()))
}

func TestInsertBuffer_Truncate(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("abcdef")
	c.NoError(err)
	b.Truncate(3)
	c.Equal("abc", string(b.Bytes()))
}

func TestInsertBuffer_WriteTo(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("hello")
	c.NoError(err)
	var buf bytes.Buffer
	n, err := b.WriteTo(&buf)
	c.NoError(err)
	c.Equal(int64(5), n)
	c.Equal("hello", buf.String())
}

func TestInsertBuffer_WriteToShortWrite(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("hello")
	c.NoError(err)
	w := &shortWriter{n: 2}
	n, err := b.WriteTo(w)
	c.Equal(io.ErrShortWrite, err)
	c.Equal(int64(2), n)
}

func TestInsertBuffer_WriteToError(t *testing.T) {
	c := check.New(t)
	var b xbytes.InsertBuffer
	_, err := b.WriteString("hello")
	c.NoError(err)
	w := &shortWriter{err: io.ErrClosedPipe, n: 2}
	n, err := b.WriteTo(w)
	c.Equal(io.ErrClosedPipe, err)
	c.Equal(int64(0), n)
}

type shortWriter struct {
	err error
	n   int
}

func (w *shortWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if len(p) > w.n {
		p = p[:w.n]
	}
	return len(p), nil
}
