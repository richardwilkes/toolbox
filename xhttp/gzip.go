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
	status      int
	wroteHeader bool
	committed   bool
}

// GZipWrap wraps the given handler, providing automatic gzip compression when requests advertise that they accept the
// gzip encoding.
func GZipWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if acceptsGzip(req.Header.Get("Accept-Encoding")) {
			gw := &gzipResponseWriter{w: w}
			defer gw.finish(req)
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
	if !w.committed {
		// Defer committing (and thus the compress-or-not decision) until there is actual body data. This keeps an
		// empty-body response from gaining a Content-Encoding: gzip header and an (empty) gzip stream in place of its
		// Content-Length, and ensures a leading zero-length write does not disable compression for the bytes that
		// follow it.
		if len(data) == 0 {
			return 0, nil
		}
		w.commit(true, data)
	}
	if w.gw != nil {
		return w.gw.Write(data)
	}
	return w.w.Write(data)
}

// WriteHeader implements http.ResponseWriter. It records the final status but defers forwarding it to the underlying
// writer until the response is committed, which happens on the first non-empty Write, a Flush, or when the handler
// returns. Deferring the commit lets the compress-or-not decision be made once it is known whether the response
// actually carries a body, so a body-less response is sent unmodified rather than being advertised as gzip and emitted
// as a gzip stream.
func (w *gzipResponseWriter) WriteHeader(status int) {
	// 1xx responses are interim; the final status arrives in a later call, so forward them without committing.
	if status >= http.StatusContinue && status < http.StatusOK {
		w.w.WriteHeader(status)
		return
	}
	if !w.wroteHeader {
		w.wroteHeader = true
		w.status = status
	}
}

// commit forwards the recorded status and headers to the underlying writer exactly once. When useGzip is true and the
// response is permitted to carry a body that the handler has not already encoded (not a 204 No Content or 304 Not
// Modified, and no existing Content-Encoding), it installs gzip compression: overwriting an existing Content-Encoding
// or wrapping already-encoded bytes in a second gzip stream would corrupt the content, so those cases are left
// untouched. When gzip is installed, the handler's Content-Length describes the uncompressed body and no longer matches
// the compressed bytes, so it is dropped and the response is sent with chunked transfer encoding. The body slice, when
// non-nil, is the first chunk of the uncompressed body and is used to sniff a Content-Type before gzip is installed.
func (w *gzipResponseWriter) commit(useGzip bool, body []byte) {
	if w.committed {
		return
	}
	w.committed = true
	if useGzip && w.status != http.StatusNoContent && w.status != http.StatusNotModified &&
		w.w.Header().Get("Content-Encoding") == "" {
		header := w.w.Header()
		// net/http skips its automatic Content-Type sniffing once Content-Encoding is set (see Go issue #31753), and
		// sniffing the compressed bytes would be wrong anyway, so sniff the uncompressed body ourselves here to
		// preserve the Content-Type a handler relies on being detected automatically. http.DetectContentType only
		// examines the first 512 bytes, so the first body chunk is enough.
		if len(body) > 0 && header.Get("Content-Type") == "" {
			header.Set("Content-Type", http.DetectContentType(body))
		}
		header.Set("Content-Encoding", "gzip")
		header.Del("Content-Length")
		w.gw = gzip.NewWriter(w.w)
	}
	w.w.WriteHeader(w.status)
}

// finish runs after the wrapped handler returns. It forwards the status of a response the handler committed to (via
// WriteHeader) but never gave a body, doing so without a gzip stream so the response keeps its Content-Length, and
// closes the gzip writer if one was installed.
func (w *gzipResponseWriter) finish(req *http.Request) {
	if w.wroteHeader && !w.committed {
		w.commit(false, nil)
	}
	if w.gw != nil {
		xio.CloseLoggingErrorsTo(LoggerForRequest(req), w.gw)
	}
}

// Flush implements http.Flusher. A flush signals a streaming body, so it commits the response (installing gzip
// compression) if that has not already happened, then flushes any buffered compressed data to the underlying writer
// before flushing the underlying writer itself, so streaming responses such as Server-Sent Events reach the client
// promptly.
func (w *gzipResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if !w.committed {
		w.commit(true, nil)
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
