// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xio_test

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xio"
)

type mockCloser struct {
	shouldError bool
	closed      bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	if m.shouldError {
		return errors.New("mock close error")
	}
	return nil
}

type mockReadCloser struct {
	data   []byte
	pos    int
	closer mockCloser
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockReadCloser) Close() error {
	return m.closer.Close()
}

func TestCloseIgnoringErrors(t *testing.T) {
	c := check.New(t)

	closer := &mockCloser{shouldError: false}
	xio.CloseIgnoringErrors(closer)
	c.True(closer.closed)

	closerWithError := &mockCloser{shouldError: true}
	xio.CloseIgnoringErrors(closerWithError)
	c.True(closerWithError.closed)
}

func TestDiscardAndCloseIgnoringErrors(t *testing.T) {
	c := check.New(t)
	for _, tc := range []struct {
		size        int
		wantRead    int
		closeErrors bool
	}{
		{size: 9, wantRead: 9},
		{size: 9, wantRead: 9, closeErrors: true},
		{size: 4 * 1024 * 1024, wantRead: 256 * 1024}, // Stops at the drain byte bound.
	} {
		rc := &mockReadCloser{data: make([]byte, tc.size), closer: mockCloser{shouldError: tc.closeErrors}}
		xio.DiscardAndCloseIgnoringErrors(rc)
		c.True(rc.closer.closed)
		c.Equal(tc.wantRead, rc.pos)
	}
}

// blockingReadCloser blocks in Read until it is closed, simulating a source that stalls without delivering data or EOF.
type blockingReadCloser struct {
	closeCh chan struct{}
	once    sync.Once
	closed  atomic.Bool
}

func (r *blockingReadCloser) Read(_ []byte) (int, error) {
	<-r.closeCh
	return 0, io.EOF
}

func (r *blockingReadCloser) Close() error {
	r.once.Do(func() { close(r.closeCh) })
	r.closed.Store(true)
	return nil
}

func TestDiscardAndCloseIgnoringErrors_TimeBound(t *testing.T) {
	c := check.New(t)
	rc := &blockingReadCloser{closeCh: make(chan struct{})}
	start := time.Now()
	xio.DiscardAndCloseIgnoringErrors(rc)
	elapsed := time.Since(start)
	// The time bound must fire and close the stalled reader rather than hanging indefinitely.
	c.True(rc.closed.Load())
	c.True(elapsed < 5*time.Second)
}

func TestCloseLoggingAnyError(t *testing.T) {
	c := check.New(t)

	closer := &mockCloser{shouldError: false}
	xio.CloseLoggingErrors(closer)
	c.True(closer.closed)

	oldLogger := slog.Default()
	defer func() { slog.SetDefault(oldLogger) }()
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

	closerWithError := &mockCloser{shouldError: true}
	xio.CloseLoggingErrors(closerWithError)
	c.True(closerWithError.closed)
	c.Contains(buf.String(), "mock close error")
}

func TestCloseLoggingAnyErrorTo(t *testing.T) {
	c := check.New(t)

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	closer := &mockCloser{shouldError: false}
	xio.CloseLoggingErrorsTo(logger, closer)
	c.True(closer.closed)
	c.Equal("", buf.String())

	buf.Reset()
	closerWithError := &mockCloser{shouldError: true}
	xio.CloseLoggingErrorsTo(logger, closerWithError)
	c.True(closerWithError.closed)
	logOutput := buf.String()
	c.True(strings.Contains(logOutput, "mock close error"))
}
