// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
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
	_ http.ResponseWriter = &StatusWriter{}
	_ http.Flusher        = &StatusWriter{}
	_ http.Hijacker       = &StatusWriter{}
	_ http.Pusher         = &StatusWriter{}
)

// StatusWriter wraps an http.ResponseWriter and provides methods to retrieve the status code and number of bytes
// written.
type StatusWriter struct {
	w      http.ResponseWriter
	status int
	bytes  int
	headOp bool
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

// Header implements http.ResponseWriter.
func (w *StatusWriter) Header() http.Header {
	return w.w.Header()
}

// Write implements http.ResponseWriter.
func (w *StatusWriter) Write(data []byte) (int, error) {
	if w.headOp {
		return len(data), nil
	}
	n, err := w.w.Write(data)
	w.bytes += n
	return n, err
}

// WriteHeader implements http.ResponseWriter.
func (w *StatusWriter) WriteHeader(status int) {
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
