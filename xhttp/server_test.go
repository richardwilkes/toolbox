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
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

// recordingWriter is an http.ResponseWriter that records every status code passed to WriteHeader, so tests can detect a
// missing or superfluous WriteHeader call.
type recordingWriter struct {
	header       http.Header
	writeHeaders []int
	body         bytes.Buffer
}

func (r *recordingWriter) Header() http.Header {
	if r.header == nil {
		r.header = make(http.Header)
	}
	return r.header
}

func (r *recordingWriter) Write(data []byte) (int, error) { return r.body.Write(data) }

func (r *recordingWriter) WriteHeader(status int) { r.writeHeaders = append(r.writeHeaders, status) }

func newTestServer(c check.Checker, log *bytes.Buffer, handler http.HandlerFunc) *xhttp.Server {
	s, err := xhttp.NewServer(&xhttp.ServerConfig{
		Logger:  slog.New(slog.NewTextHandler(log, nil)),
		Handler: handler,
	})
	c.NoError(err)
	return s
}

// TestServerPanicBeforeWriteRecordsStatus verifies that when a handler panics before writing anything, the recovery
// sends a 500 through the StatusWriter so both the client and the access log see status 500. Previously the 500 was
// written to the raw writer, leaving the StatusWriter (and therefore the log) reporting 200.
func TestServerPanicBeforeWriteRecordsStatus(t *testing.T) {
	c := check.New(t)
	var log bytes.Buffer
	s := newTestServer(c, &log, func(_ http.ResponseWriter, _ *http.Request) {
		panic("boom before write")
	})
	rec := &recordingWriter{}
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	c.Equal([]int{http.StatusInternalServerError}, rec.writeHeaders)
	c.Contains(log.String(), "status=500")
}

// TestServerPanicAfterWriteDoesNotResend verifies that when a handler panics after committing a response, the recovery
// does not call WriteHeader again (which the underlying writer would ignore and log as superfluous), and the access log
// reports the status the handler actually sent.
func TestServerPanicAfterWriteDoesNotResend(t *testing.T) {
	c := check.New(t)
	var log bytes.Buffer
	s := newTestServer(c, &log, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, "partial") //nolint:errcheck // Ignored for testing
		panic("boom after write")
	})
	rec := &recordingWriter{}
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	// Only the handler's own header should have been written; no superfluous 500 from recovery.
	c.Equal([]int{http.StatusAccepted}, rec.writeHeaders)
	c.Equal("partial", rec.body.String())
	c.Contains(log.String(), "status=202")
}

// TestServerRunReturnsServeError verifies that a serve-time failure (here, a TLS server pointed at non-existent
// certificate files) is returned from Run itself, not merely stashed away for Error to report. Previously Run always
// returned nil, so callers doing `if err := srv.Run(); err != nil` never observed serve failures.
func TestServerRunReturnsServeError(t *testing.T) {
	c := check.New(t)
	var log bytes.Buffer
	s, err := xhttp.NewServer(&xhttp.ServerConfig{
		Logger:   slog.New(slog.NewTextHandler(&log, nil)),
		CertFile: filepath.Join(t.TempDir(), "missing-cert.pem"),
		KeyFile:  filepath.Join(t.TempDir(), "missing-key.pem"),
	})
	c.NoError(err)
	c.Equal(xhttp.ProtocolHTTPS, s.Protocol())

	runErr := s.Run()
	c.HasError(runErr)
	// The same error must also be observable via Error.
	c.Equal(runErr, s.Error())
}

// TestServerRunReturnsNilOnCleanShutdown verifies that a normal shutdown via Stop causes Run to return nil rather than
// leaking http.ErrServerClosed to the caller.
func TestServerRunReturnsNilOnCleanShutdown(t *testing.T) {
	c := check.New(t)
	var log bytes.Buffer
	s := newTestServer(c, &log, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	runErr := make(chan error, 1)
	go func() { runErr <- s.Run() }()
	s.WaitForStart()
	s.Stop()
	c.NoError(<-runErr)
	c.NoError(s.Error())
}

// newTestServerWithWriteTimeout builds a server with the given WriteTimeout so tests can exercise the request-context
// deadline that mirrors the connection's write deadline.
func newTestServerWithWriteTimeout(c check.Checker, writeTimeout time.Duration, handler http.HandlerFunc) *xhttp.Server {
	var log bytes.Buffer
	s, err := xhttp.NewServer(&xhttp.ServerConfig{
		Logger:       slog.New(slog.NewTextHandler(&log, nil)),
		Handler:      handler,
		WriteTimeout: writeTimeout,
	})
	c.NoError(err)
	return s
}

// TestServerWriteTimeoutAnchorsContextDeadlineAtRequestStart verifies that, when WriteTimeout is configured, the
// request context handed to the handler carries a deadline of WriteTimeout anchored to the moment the request is
// received (which is when net/http arms the socket's write deadline), not some later or earlier instant.
func TestServerWriteTimeoutAnchorsContextDeadlineAtRequestStart(t *testing.T) {
	c := check.New(t)
	const writeTimeout = 200 * time.Millisecond
	var gotDeadline time.Time
	var gotOK bool
	s := newTestServerWithWriteTimeout(c, writeTimeout, func(_ http.ResponseWriter, req *http.Request) {
		gotDeadline, gotOK = req.Context().Deadline()
	})

	before := time.Now()
	s.ServeHTTP(&recordingWriter{}, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	after := time.Now()

	c.True(gotOK, "handler context should have a deadline when WriteTimeout is set")
	// The anchor is the request-start instant captured inside ServeHTTP, which necessarily falls within [before,
	// after], so the deadline must fall within [before+WriteTimeout, after+WriteTimeout]. This pins the anchor to
	// request entry rather than to some later point after this method's setup work.
	c.False(gotDeadline.Before(before.Add(writeTimeout)), "deadline anchored before the request was received")
	c.False(gotDeadline.After(after.Add(writeTimeout)), "deadline anchored after the request completed")
}

// TestServerNoWriteTimeoutLeavesContextWithoutDeadline verifies that when WriteTimeout is not configured, the handler's
// request context is left without an added deadline.
func TestServerNoWriteTimeoutLeavesContextWithoutDeadline(t *testing.T) {
	c := check.New(t)
	var gotOK bool
	s := newTestServerWithWriteTimeout(c, 0, func(_ http.ResponseWriter, req *http.Request) {
		_, gotOK = req.Context().Deadline()
	})
	s.ServeHTTP(&recordingWriter{}, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	c.False(gotOK, "handler context should not have a deadline when WriteTimeout is unset")
}

// TestServerWriteTimeoutResetPerRequest verifies that the context deadline is derived fresh for each request rather
// than drifting across requests. This mirrors a reused keep-alive connection, where net/http re-arms the socket write
// deadline for every request: a request that arrives after an idle gap must still get the full WriteTimeout budget, not
// a shrunken remainder.
func TestServerWriteTimeoutResetPerRequest(t *testing.T) {
	c := check.New(t)
	const writeTimeout = 200 * time.Millisecond
	var remaining []time.Duration
	s := newTestServerWithWriteTimeout(c, writeTimeout, func(_ http.ResponseWriter, req *http.Request) {
		deadline, ok := req.Context().Deadline()
		c.True(ok)
		remaining = append(remaining, time.Until(deadline))
	})

	s.ServeHTTP(&recordingWriter{}, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	// Simulate an idle gap between two requests on the same keep-alive connection.
	time.Sleep(writeTimeout * 3 / 4)
	s.ServeHTTP(&recordingWriter{}, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	c.Equal(2, len(remaining))
	// If the deadline drifted (anchored once and shared), the second request's remaining budget would be roughly a
	// quarter of WriteTimeout after the sleep. A fresh per-request deadline leaves nearly the full budget.
	c.True(remaining[1] > writeTimeout*3/4,
		"second request should get a fresh WriteTimeout budget, got", remaining[1])
}

// TestServerWriteTimeoutCancelsSlowHandler verifies the deadline actually fires: a handler that outlives WriteTimeout
// observes ctx.Done() with a DeadlineExceeded error rather than blocking indefinitely.
func TestServerWriteTimeoutCancelsSlowHandler(t *testing.T) {
	c := check.New(t)
	const writeTimeout = 100 * time.Millisecond
	var ctxErr error
	s := newTestServerWithWriteTimeout(c, writeTimeout, func(_ http.ResponseWriter, req *http.Request) {
		select {
		case <-req.Context().Done():
			ctxErr = req.Context().Err()
		case <-time.After(5 * time.Second):
			ctxErr = errors.New("handler was not canceled by the write timeout")
		}
	})
	s.ServeHTTP(&recordingWriter{}, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	c.True(errors.Is(ctxErr, context.DeadlineExceeded), "expected DeadlineExceeded, got", ctxErr)
}
