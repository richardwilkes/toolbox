// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio

import (
	"bytes"
)

// LineWriter buffers its input into lines before sending each line to an output function without the trailing line
// feed.
type LineWriter struct {
	buffer *bytes.Buffer
	out    func([]byte)
}

// NewLineWriter creates a new LineWriter.
func NewLineWriter(out func([]byte)) *LineWriter {
	return &LineWriter{buffer: &bytes.Buffer{}, out: out}
}

// Write implements the io.Writer interface.
func (w *LineWriter) Write(data []byte) (n int, err error) {
	n = len(data)
	for len(data) > 0 {
		i := bytes.IndexByte(data, '\n')
		if i == -1 {
			_, err = w.buffer.Write(data)
			return n, err
		}
		if i > 0 {
			if _, err = w.buffer.Write(data[:i]); err != nil {
				return n, err
			}
		}
		w.out(w.buffer.Bytes())
		w.buffer.Reset()
		data = data[i+1:]
	}
	return n, nil
}

// Close implements the io.Closer interface.
func (w *LineWriter) Close() error {
	if w.buffer.Len() > 0 {
		w.out(w.buffer.Bytes())
		w.buffer.Reset()
	}
	return nil
}
