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

func TestGZipWrapRoundTrip(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := io.WriteString(w, "hello, gzip")
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
	c.Equal("hello, gzip", string(data))
}

func TestGZipWrapNoGzipWhenNotAccepted(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := io.WriteString(w, "hello, plain")
		c.NoError(err)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	c.Equal("", rec.Header().Get("Content-Encoding"))
	c.Equal("hello, plain", rec.Body.String())
}

// TestGZipWrapClearsContentLength verifies that when the response is compressed, a Content-Length the handler set for
// the uncompressed body is removed, since it no longer matches the (smaller) compressed bytes. Leaving it would make
// the client wait for bytes that never arrive or treat the response as truncated.
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

// TestGZipWrapKeepsContentLengthWhenNotCompressing verifies that a Content-Length set on a response that is not
// compressed (here a 204 that must not carry a body) is left untouched.
func TestGZipWrapKeepsContentLengthWhenNotCompressing(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal("", rec.Header().Get("Content-Encoding"))
	c.Equal("0", rec.Header().Get("Content-Length"))
}

// TestGZipWrapNoContent verifies that a 204 No Content response is neither advertised as gzip nor emitted as a gzip
// stream, since responses with that status must not carry a body.
func TestGZipWrapNoContent(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal(http.StatusNoContent, rec.Code)
	c.Equal("", rec.Header().Get("Content-Encoding"))
	c.Equal(0, rec.Body.Len())
}

// TestGZipWrapNotModified verifies that a 304 Not Modified response is neither advertised as gzip nor emitted as a gzip
// stream, since responses with that status must not carry a body.
func TestGZipWrapNotModified(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal(http.StatusNotModified, rec.Code)
	c.Equal("", rec.Header().Get("Content-Encoding"))
	c.Equal(0, rec.Body.Len())
}

// TestGZipWrapExplicitStatusRoundTrip verifies that a handler which explicitly writes a body-bearing status still gets
// its body compressed and the encoding advertised.
func TestGZipWrapExplicitStatusRoundTrip(t *testing.T) {
	c := check.New(t)
	handler := xhttp.GZipWrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, "hello, gzip")
		c.NoError(err)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, gzipAcceptingRequest())
	c.Equal(http.StatusOK, rec.Code)
	c.Equal("gzip", rec.Header().Get("Content-Encoding"))
	gr, err := gzip.NewReader(rec.Body)
	c.NoError(err)
	data, err := io.ReadAll(gr)
	c.NoError(err)
	c.NoError(gr.Close())
	c.Equal("hello, gzip", string(data))
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
