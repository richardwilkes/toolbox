// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp

import (
	"bufio"
	"net"
	"net/http"
)

var (
	_ http.ResponseWriter                       = &StatusWriter{}
	_ http.Flusher                              = &StatusWriter{}
	_ http.Hijacker                             = &StatusWriter{}
	_ http.Pusher                               = &StatusWriter{}
	_ interface{ Unwrap() http.ResponseWriter } = &StatusWriter{}
)

// StatusWriter wraps an http.ResponseWriter and provides methods to retrieve the status code and number of bytes
// written.
type StatusWriter struct {
	w           http.ResponseWriter
	status      int
	bytes       int
	headOp      bool
	wroteHeader bool
}

// NewStatusWriter creates a new StatusWriter.
func NewStatusWriter(w http.ResponseWriter, req *http.Request) *StatusWriter {
	return &StatusWriter{
		w:      w,
		status: http.StatusOK,
		headOp: req.Method == http.MethodHead,
	}
}

// Status returns the status that was set, or http.StatusOK if no call to WriteHeader() was made.
func (w *StatusWriter) Status() int {
	return w.status
}

// BytesWritten returns the number of bytes written.
func (w *StatusWriter) BytesWritten() int {
	return w.bytes
}

// HeaderWritten returns true if a response header has been committed, either explicitly via WriteHeader() or implicitly
// by the first call to Write(). Once this is true, a further WriteHeader() call would be ignored by the underlying
// writer (and logged as superfluous), so callers such as panic recovery can use it to decide whether it is still
// possible to set a status.
func (w *StatusWriter) HeaderWritten() bool {
	return w.wroteHeader
}

// Header implements http.ResponseWriter.
func (w *StatusWriter) Header() http.Header {
	return w.w.Header()
}

// Write implements http.ResponseWriter.
func (w *StatusWriter) Write(data []byte) (int, error) {
	if w.headOp {
		return len(data), nil
	}
	w.wroteHeader = true // The first write commits the implicit 200 header on the underlying writer.
	n, err := w.w.Write(data)
	w.bytes += n
	return n, err
}

// WriteHeader implements http.ResponseWriter.
func (w *StatusWriter) WriteHeader(status int) {
	w.wroteHeader = true
	w.status = status
	w.w.WriteHeader(status)
}

// Flush implements http.Flusher.
func (w *StatusWriter) Flush() {
	if f, ok := w.w.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker.
func (w *StatusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.w.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Push implements http.Pusher.
func (w *StatusWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Unwrap returns the wrapped http.ResponseWriter so that http.ResponseController can reach optional interfaces it does
// not implement itself, such as deadline control (SetReadDeadline/SetWriteDeadline), EnableFullDuplex, and FlushError.
func (w *StatusWriter) Unwrap() http.ResponseWriter {
	return w.w
}
