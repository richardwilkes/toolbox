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
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/v2/xio"
)

var (
	_ http.ResponseWriter                       = &gzipResponseWriter{}
	_ http.Flusher                              = &gzipResponseWriter{}
	_ http.Hijacker                             = &gzipResponseWriter{}
	_ http.Pusher                               = &gzipResponseWriter{}
	_ interface{ Unwrap() http.ResponseWriter } = &gzipResponseWriter{}
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
		if acceptsGzip(req.Header.Get("Accept-Encoding")) {
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

// acceptsGzip reports whether the Accept-Encoding header value indicates the client accepts the gzip encoding. It
// parses the comma-separated codings and their optional q-values per RFC 7231, so an explicit refusal ("gzip;q=0") is
// honored and gzip is only matched as a whole coding (not as a substring of, e.g., "x-gzip"). When gzip is not named
// explicitly, a wildcard ("*") with a non-zero q-value also makes it acceptable.
func acceptsGzip(acceptEncoding string) bool {
	gzipQ := -1.0 // q-value of an explicit "gzip" coding, or -1 if absent
	starQ := -1.0 // q-value of a "*" wildcard coding, or -1 if absent
	for part := range strings.SplitSeq(acceptEncoding, ",") {
		switch coding, q := parseCoding(part); coding {
		case "gzip":
			gzipQ = q
		case "*":
			starQ = q
		}
	}
	if gzipQ >= 0 {
		return gzipQ > 0
	}
	return starQ > 0
}

// parseCoding parses a single Accept-Encoding element into its lowercased coding name and q-value (defaulting to 1.0
// per RFC 7231 when no valid q parameter is present).
func parseCoding(part string) (coding string, q float64) {
	segments := strings.Split(part, ";")
	coding = strings.ToLower(strings.TrimSpace(segments[0]))
	q = 1
	for _, seg := range segments[1:] {
		if value, ok := strings.CutPrefix(strings.ToLower(strings.TrimSpace(seg)), "q="); ok {
			if parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64); err == nil {
				q = parsed
			}
		}
	}
	return coding, q
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
// which must not carry a body (1xx informational, 204 No Content, 304 Not Modified), or which the handler has already
// encoded itself (a non-empty Content-Encoding), are neither advertised as gzip nor emitted as a gzip stream, both of
// which would produce an invalid or doubly-encoded response.
func (w *gzipResponseWriter) WriteHeader(status int) {
	// 1xx responses are interim; the final status arrives in a later call, so forward them without committing.
	if status >= http.StatusContinue && status < http.StatusOK {
		w.w.WriteHeader(status)
		return
	}
	if !w.wroteHeader {
		w.wroteHeader = true
		// Skip a body-less status, and skip a response the handler already encoded: overwriting its Content-Encoding
		// and wrapping its bytes in a second gzip stream would corrupt the content.
		if status != http.StatusNoContent && status != http.StatusNotModified &&
			w.w.Header().Get("Content-Encoding") == "" {
			w.w.Header().Set("Content-Encoding", "gzip")
			// The handler's Content-Length describes the uncompressed body and no longer matches the compressed bytes
			// we are about to emit, so drop it and let the response be sent with chunked transfer encoding.
			w.w.Header().Del("Content-Length")
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

// Unwrap returns the wrapped http.ResponseWriter so that http.ResponseController can reach optional interfaces this
// writer does not implement itself, such as deadline control (SetReadDeadline/SetWriteDeadline) and EnableFullDuplex.
// Flush is intentionally implemented here (rather than delegated via Unwrap) so the gzip stream is flushed correctly.
func (w *gzipResponseWriter) Unwrap() http.ResponseWriter {
	return w.w
}
