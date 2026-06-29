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
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xhttp"
)

// recordingWriter is an http.ResponseWriter that records every status code passed to WriteHeader, so tests can detect a
// missing or superfluous WriteHeader call.
type recordingWriter struct {
	header       http.Header
	body         bytes.Buffer
	writeHeaders []int
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
		_, _ = io.WriteString(w, "partial")
		panic("boom after write")
	})
	rec := &recordingWriter{}
	s.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", http.NoBody))

	// Only the handler's own header should have been written; no superfluous 500 from recovery.
	c.Equal([]int{http.StatusAccepted}, rec.writeHeaders)
	c.Equal("partial", rec.body.String())
	c.Contains(log.String(), "status=202")
}
