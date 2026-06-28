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
	"compress/gzip"
	"net"
	"net/http"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xio"
)

var (
	_ http.ResponseWriter = &gzipResponseWriter{}
	_ http.Flusher        = &gzipResponseWriter{}
	_ http.Hijacker       = &gzipResponseWriter{}
	_ http.Pusher         = &gzipResponseWriter{}
)

type gzipResponseWriter struct {
	w           http.ResponseWriter
	gw          *gzip.Writer
	wroteHeader bool
}

// GZipWrap wraps the given handler, providing automatic gzip compression when requests advertise that they accept the
// gzip encoding.
func GZipWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			gw := &gzipResponseWriter{w: w}
			defer func() {
				if gw.gw != nil {
					xio.CloseLoggingErrorsTo(LoggerForRequest(req), gw.gw)
				}
			}()
			w = gw
		}
		next.ServeHTTP(w, req)
	})
}

// Header implements http.ResponseWriter.
func (w *gzipResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Write implements http.ResponseWriter.
func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.gw != nil {
		return w.gw.Write(data)
	}
	return w.w.Write(data)
}

// WriteHeader implements http.ResponseWriter. The decision to compress is deferred until this point so that responses
// which must not carry a body (1xx informational, 204 No Content, 304 Not Modified) are neither advertised as gzip nor
// emitted as a gzip stream, both of which would produce an invalid response.
func (w *gzipResponseWriter) WriteHeader(status int) {
	// 1xx responses are interim; the final status arrives in a later call, so forward them without committing.
	if status >= http.StatusContinue && status < http.StatusOK {
		w.w.WriteHeader(status)
		return
	}
	if !w.wroteHeader {
		w.wroteHeader = true
		if status != http.StatusNoContent && status != http.StatusNotModified {
			w.w.Header().Set("Content-Encoding", "gzip")
			w.gw = gzip.NewWriter(w.w)
		}
	}
	w.w.WriteHeader(status)
}

// Flush implements http.Flusher. It flushes any buffered compressed data to the underlying writer before flushing the
// underlying writer itself, so streaming responses such as Server-Sent Events reach the client promptly.
func (w *gzipResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.gw != nil {
		if err := w.gw.Flush(); err != nil {
			return
		}
	}
	if f, ok := w.w.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker.
func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.w.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Push implements http.Pusher.
func (w *gzipResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
