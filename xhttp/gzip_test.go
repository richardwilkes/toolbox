// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp_test

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

func gzipAcceptingRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("Accept-Encoding", "gzip")
	return req
}

// plainWriter is an http.ResponseWriter that implements none of the optional streaming/upgrade interfaces.
type plainWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (p *plainWriter) Header() http.Header {
	if p.header == nil {
		p.header = make(http.Header)
	}
	return p.header
}

func (p *plainWriter) Write(data []byte) (int, error) { return p.body.Write(data) }

func (p *plainWriter) WriteHeader(status int) { p.status = status }

// hijackWriter adds http.Hijacker support to a recorder.
type hijackWriter struct {
	*httptest.ResponseRecorder
	hijacked bool
}

func (h *hijackWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.hijacked = true
	return nil, nil, nil
}

// pushWriter adds http.Pusher support to a recorder.
type pushWriter struct {
	*httptest.ResponseRecorder
	pushed string
}

func (p *pushWriter) Push(target string, _ *http.PushOptions) error {
	p.pushed = target
	return nil
}

// deadlineWriter is an http.ResponseWriter that records the deadlines set via the per-request deadline methods, which
// http.ResponseController reaches by descending through Unwrap.
type deadlineWriter struct {
	http.ResponseWriter
	readDeadline  time.Time
	writeDeadline time.Time
}

func (d *deadlineWriter) SetReadDeadline(deadline time.Time) error {
	d.readDeadline = deadline
	return nil
}

func (d *deadlineWriter) SetWriteDeadline(deadline time.Time) error {
	d.writeDeadline = deadline
	return nil
}

// TestGZipWrapUnwrapReachesDeadlines verifies that http.ResponseController can reach the underlying writer's deadline
// methods through the gzip wrapper, which it can only do if gzipResponseWriter implements Unwrap.
func TestGZipWrapUnwrapReachesDeadlines(t *testing.T) {
	c := check.New(t)
	underlying := &deadlineWriter{ResponseWriter: httptest.NewRecorder()}
	deadline := time.Now().Add(time.Minute)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		rc := http.NewResponseController(w)
		c.NoError(rc.SetReadDeadline(deadline))
		c.NoError(rc.SetWriteDeadline(deadline))
	}))
	handler.ServeHTTP(underlying, gzipAcceptingRequest())
	c.True(ran)
	c.Equal(deadline, underlying.readDeadline)
	c.Equal(deadline, underlying.writeDeadline)
}

// TestGZipWrapSniffsOnlyUnsetContentType verifies that a handler which writes a body but never sets a Content-Type
// still gets a correctly sniffed one, while a Content-Type the handler set is kept even when the sniffed type would
// differ. Once Content-Encoding: gzip is set, net/http stops sniffing (Go issue #31753), so the wrapper must sniff the
// uncompressed body before installing gzip.
func TestGZipWrapSniffsOnlyUnsetContentType(t *testing.T) {
	c := check.New(t)
	const body = "<!DOCTYPE html><html><head><title>hi</title></head><body>hello</body></html>"
	for _, tc := range []struct {
		contentType     string // Set by the handler when non-empty.
		wantContentType string
	}{
		{wantContentType: "text/html; charset=utf-8"},
		{contentType: "application/json", wantContentType: "application/json"},
	} {
		handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if tc.contentType != "" {
				w.Header().Set("Content-Type", tc.contentType)
			}
			_, err := io.WriteString(w, body)
			c.NoError(err)
		}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, gzipAcceptingRequest())
		c.Equal("gzip", rec.Header().Get("Content-Encoding"))
		c.Equal(tc.wantContentType, rec.Header().Get("Content-Type"), "handler Content-Type %q", tc.contentType)
		gr, err := gzip.NewReader(rec.Body)
		c.NoError(err)
		data, err := io.ReadAll(gr)
		c.NoError(err)
		c.NoError(gr.Close())
		c.Equal(body, string(data))
	}
}

// TestGZipWrapClearsContentLength verifies that when the response is compressed, a Content-Length the handler set for
// the uncompressed body is removed, since it no longer matches the compressed bytes.
func TestGZipWrapClearsContentLength(t *testing.T) {
	c := check.New(t)
	const body = "the quick brown fox jumps over the lazy dog, repeatedly and at length"
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, err := io.WriteString(w, body)
		c.NoError(err)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal("gzip", rec.Header().Get("Content-Encoding"))
	c.Equal("", rec.Header().Get("Content-Length"), "stale uncompressed Content-Length must be removed")
	gr, err := gzip.NewReader(rec.Body)
	c.NoError(err)
	data, err := io.ReadAll(gr)
	c.NoError(err)
	c.NoError(gr.Close())
	c.Equal(body, string(data))
}

// TestGZipWrapSkipsAlreadyEncoded verifies that a response the handler has already encoded (a non-empty
// Content-Encoding) is left untouched: its header is not overwritten to gzip and its bytes are not wrapped in a second
// gzip stream.
func TestGZipWrapSkipsAlreadyEncoded(t *testing.T) {
	c := check.New(t)
	const body = "pretend-brotli-bytes"
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "br")
		_, err := io.WriteString(w, body)
		c.NoError(err)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal("br", rec.Header().Get("Content-Encoding"), "handler's Content-Encoding must be preserved")
	c.Equal(body, rec.Body.String(), "already-encoded body must be passed through verbatim, not re-gzipped")
}

// TestGZipWrapAcceptEncodingNegotiation verifies that the Accept-Encoding header is parsed per RFC 7231: q-values are
// honored (an explicit gzip;q=0 refusal is not compressed), gzip is matched only as a whole coding (not as a substring
// of x-gzip), and a wildcard makes gzip acceptable only when gzip is not explicitly refused.
func TestGZipWrapAcceptEncodingNegotiation(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		acceptEncoding string
		wantGzip       bool
	}{
		{acceptEncoding: "gzip", wantGzip: true},
		{acceptEncoding: "GZIP", wantGzip: true}, // Codings are case-insensitive.
		{acceptEncoding: "gzip;q=0.5", wantGzip: true},
		{acceptEncoding: "deflate, gzip", wantGzip: true},
		{acceptEncoding: "*", wantGzip: true},
		{acceptEncoding: "br, *", wantGzip: true},
		{acceptEncoding: "gzip;q=0", wantGzip: false},   // Explicit refusal.
		{acceptEncoding: "gzip ; q=0", wantGzip: false}, // Refusal with whitespace.
		{acceptEncoding: "deflate, gzip;q=0", wantGzip: false},
		{acceptEncoding: "gzip;q=0, *", wantGzip: false}, // Explicit refusal overrides the wildcard.
		{acceptEncoding: "x-gzip", wantGzip: false},      // Must not match as a substring.
		{acceptEncoding: "identity", wantGzip: false},
		{acceptEncoding: "*;q=0", wantGzip: false},
		{acceptEncoding: "", wantGzip: false}, // Header absent.
	} {
		handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := io.WriteString(w, "body")
			c.NoError(err)
		}))
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		if tc.acceptEncoding != "" {
			req.Header.Set("Accept-Encoding", tc.acceptEncoding)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		gotGzip := rec.Header().Get("Content-Encoding") == "gzip"
		c.Equal(tc.wantGzip, gotGzip, "Accept-Encoding: %q", tc.acceptEncoding)
		if !tc.wantGzip {
			c.Equal("body", rec.Body.String(), "Accept-Encoding: %q must pass the body through uncompressed",
				tc.acceptEncoding)
		}
	}
}

// TestGZipWrapLeavesEmptyBodyUncompressed verifies that a response without a body is not advertised as gzip and does
// not emit an (~20-byte) empty gzip stream, whether or not the handler sets a status, and that a Content-Length the
// handler set is kept.
func TestGZipWrapLeavesEmptyBodyUncompressed(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		contentLength string // Set by the handler when non-empty.
		status        int    // Zero means the handler never calls WriteHeader.
	}{
		{status: http.StatusOK},
		{status: http.StatusOK, contentLength: "0"},
		{status: http.StatusNoContent},
		{status: http.StatusNoContent, contentLength: "0"},
		{status: http.StatusNotModified},
		{},
	} {
		handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if tc.contentLength != "" {
				w.Header().Set("Content-Length", tc.contentLength)
			}
			if tc.status != 0 {
				w.WriteHeader(tc.status)
			}
		}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, gzipAcceptingRequest())
		if tc.status != 0 {
			c.Equal(tc.status, rec.Code)
		}
		c.Equal("", rec.Header().Get("Content-Encoding"), "status %d, Content-Length %q", tc.status, tc.contentLength)
		if tc.contentLength != "" {
			c.Equal(tc.contentLength, rec.Header().Get("Content-Length"), "status %d", tc.status)
		}
		c.Equal(0, rec.Body.Len(), "status %d, Content-Length %q", tc.status, tc.contentLength)
	}
}

// TestGZipWrapEmptyBodyServerSetsContentLength drives the wrapper through a real server to verify that an empty-body
// response is sent with Content-Length: 0 and no Content-Encoding, rather than as a chunked gzip stream.
func TestGZipWrapEmptyBodyServerSetsContentLength(t *testing.T) {
	c := check.New(t)
	server := httptest.NewServer(xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	c.NoError(err)
	req.Header.Set("Accept-Encoding", "gzip") // Set manually so the transport will not transparently decompress.
	resp, err := http.DefaultClient.Do(req)
	c.NoError(err)
	defer func() { c.NoError(resp.Body.Close()) }()
	body, err := io.ReadAll(resp.Body)
	c.NoError(err)
	c.Equal(http.StatusOK, resp.StatusCode)
	c.Equal("", resp.Header.Get("Content-Encoding"), "empty body must not be advertised as gzip")
	c.Equal(int64(0), resp.ContentLength, "empty body must be sent with Content-Length: 0")
	c.Equal(0, len(body), "empty body must not be sent as a gzip stream")
}

// TestGZipWrapLeadingEmptyWriteStillCompresses verifies that a leading zero-length write does not disable compression:
// a subsequent non-empty write is still gzip-compressed.
func TestGZipWrapLeadingEmptyWriteStillCompresses(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n, err := w.Write(nil)
		c.NoError(err)
		c.Equal(0, n)
		_, err = io.WriteString(w, "hello after empty")
		c.NoError(err)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal("gzip", rec.Header().Get("Content-Encoding"))
	gr, err := gzip.NewReader(rec.Body)
	c.NoError(err)
	data, err := io.ReadAll(gr)
	c.NoError(err)
	c.NoError(gr.Close())
	c.Equal("hello after empty", string(data))
}

// TestGZipWrapInterimResponse verifies that a 1xx interim response is forwarded without committing the gzip decision,
// so a body-bearing final status that follows is still compressed. A plainWriter is used because httptest's recorder
// does not model a real server's interim-response handling (it commits the first status code it sees).
func TestGZipWrapInterimResponse(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusEarlyHints)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, "after interim")
		c.NoError(err)
	}))
	pw := &plainWriter{}
	handler.ServeHTTP(pw, gzipAcceptingRequest())
	c.Equal(http.StatusOK, pw.status)
	c.Equal("gzip", pw.Header().Get("Content-Encoding"))
	gr, err := gzip.NewReader(&pw.body)
	c.NoError(err)
	data, err := io.ReadAll(gr)
	c.NoError(err)
	c.NoError(gr.Close())
	c.Equal("after interim", string(data))
}

// TestGZipWrapFlush verifies that the wrapped writer exposes http.Flusher and that a flush both pushes the buffered
// compressed data out and propagates to the underlying writer, which is what lets Server-Sent Events reach the client.
func TestGZipWrapFlush(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		f, ok := w.(http.Flusher)
		c.True(ok, "wrapped ResponseWriter must implement http.Flusher")
		_, err := io.WriteString(w, "stream me")
		c.NoError(err)
		if ok {
			f.Flush()
		}
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.True(ran)
	c.True(rec.Flushed, "flush must propagate to the underlying http.Flusher")
	gr, err := gzip.NewReader(rec.Body)
	c.NoError(err)
	data, err := io.ReadAll(gr)
	c.NoError(err)
	c.NoError(gr.Close())
	c.Equal("stream me", string(data))
}

// TestGZipWrapFlushUnsupportedUnderlying ensures Flush is a safe no-op when the underlying writer is not a Flusher.
func TestGZipWrapFlushUnsupportedUnderlying(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		f, ok := w.(http.Flusher)
		c.True(ok)
		if ok {
			f.Flush() // Must not panic even though the underlying writer cannot flush.
		}
	}))
	handler.ServeHTTP(&plainWriter{}, gzipAcceptingRequest())
	c.True(ran)
}

func TestGZipWrapHijack(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		h, ok := w.(http.Hijacker)
		c.True(ok, "wrapped ResponseWriter must implement http.Hijacker")
		if ok {
			_, _, err := h.Hijack()
			c.NoError(err)
		}
	}))
	rec := &hijackWriter{ResponseRecorder: httptest.NewRecorder()}
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.True(ran)
	c.True(rec.hijacked, "hijack must delegate to the underlying http.Hijacker")
}

func TestGZipWrapHijackUnsupportedUnderlying(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		h, ok := w.(http.Hijacker)
		c.True(ok)
		if ok {
			_, _, err := h.Hijack()
			c.HasError(err)
		}
	}))
	handler.ServeHTTP(&plainWriter{}, gzipAcceptingRequest())
	c.True(ran)
}

func TestGZipWrapPush(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		p, ok := w.(http.Pusher)
		c.True(ok, "wrapped ResponseWriter must implement http.Pusher")
		if ok {
			c.NoError(p.Push("/style.css", nil))
		}
	}))
	rec := &pushWriter{ResponseRecorder: httptest.NewRecorder()}
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.True(ran)
	c.Equal("/style.css", rec.pushed)
}

func TestGZipWrapPushUnsupportedUnderlying(t *testing.T) {
	c := check.New(t)
	ran := false
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ran = true
		p, ok := w.(http.Pusher)
		c.True(ok)
		if ok {
			c.HasError(p.Push("/style.css", nil))
		}
	}))
	handler.ServeHTTP(&plainWriter{}, gzipAcceptingRequest())
	c.True(ran)
}
