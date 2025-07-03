// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xbytes

import (
	"io"
	"unicode/utf8"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// InsertBuffer is a variable-sized buffer of bytes with Write and Insert methods. The zero value for InsertBuffer is an
// empty buffer ready to use.
type InsertBuffer struct {
	data []byte
}

// Bytes returns the underlying buffer of bytes.
func (b *InsertBuffer) Bytes() []byte {
	return b.data
}

// String returns the underlying buffer of bytes as a string. If the InsertBuffer is a nil pointer, it returns "<nil>".
func (b *InsertBuffer) String() string {
	if b == nil {
		return "<nil>"
	}
	return string(b.data)
}

// Len returns the number of bytes contained by the buffer.
func (b *InsertBuffer) Len() int {
	return len(b.data)
}

// Cap returns the capacity of the buffer.
func (b *InsertBuffer) Cap() int {
	return cap(b.data)
}

// Available returns how many bytes are unused in the buffer.
func (b *InsertBuffer) Available() int {
	return cap(b.data) - len(b.data)
}

// Truncate discards all but the first n bytes from the buffer.
func (b *InsertBuffer) Truncate(n int) {
	b.data = b.data[:n]
}

// Reset resets the buffer to be empty.
func (b *InsertBuffer) Reset() {
	b.data = b.data[:0]
}

// Insert data at the given offset.
func (b *InsertBuffer) Insert(index int, data []byte) error {
	if index < 0 || index > len(b.data) {
		return errs.New("invalid index")
	}
	if len(data) != 0 {
		b.data = append(b.data, data...)
		copy(b.data[index+len(data):], b.data[index:])
		copy(b.data[index:], data)
	}
	return nil
}

// InsertByte inserts a byte at the given offset.
func (b *InsertBuffer) InsertByte(index int, ch byte) error {
	if index < 0 || index > len(b.data) {
		return errs.New("invalid index")
	}
	b.data = append(b.data, 0)
	copy(b.data[index+1:], b.data[index:])
	b.data[index] = ch
	return nil
}

// InsertRune inserts the UTF-8 encoding of the rune at the given offset.
func (b *InsertBuffer) InsertRune(index int, r rune) error {
	if uint32(r) < utf8.RuneSelf {
		return b.InsertByte(index, byte(r))
	}
	var buffer [4]byte
	n := utf8.EncodeRune(buffer[:], r)
	return b.Insert(index, buffer[:n])
}

// InsertString inserts the string at the given offset.
func (b *InsertBuffer) InsertString(index int, s string) error {
	return b.Insert(index, []byte(s))
}

// Write appends the contents of data to the buffer.
func (b *InsertBuffer) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	return len(data), nil
}

// WriteByte appends the byte to the buffer.
func (b *InsertBuffer) WriteByte(ch byte) error {
	b.data = append(b.data, ch)
	return nil
}

// WriteRune appends the UTF-8 encoding of the rune to the buffer.
func (b *InsertBuffer) WriteRune(r rune) (int, error) {
	if uint32(r) < utf8.RuneSelf {
		b.data = append(b.data, byte(r))
		return 1, nil
	}
	i := len(b.data)
	b.data = append(b.data, 0, 0, 0, 0)
	n := utf8.EncodeRune(b.data[i:i+4], r)
	b.data = b.data[:i+n]
	return n, nil
}

// WriteString appends the string to the buffer.
func (b *InsertBuffer) WriteString(s string) (int, error) {
	b.data = append(b.data, []byte(s)...)
	return len(s), nil
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
func (b *InsertBuffer) WriteTo(w io.Writer) (int64, error) {
	var n int64
	if nBytes := b.Len(); nBytes > 0 {
		m, err := w.Write(b.data)
		n = int64(m)
		if err != nil {
			return n, err
		}
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	return n, nil
}
