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

	// Test successful close
	closer := &mockCloser{shouldError: false}
	xio.CloseIgnoringErrors(closer)
	c.True(closer.closed)

	// Test close with error (should be ignored)
	closerWithError := &mockCloser{shouldError: true}
	xio.CloseIgnoringErrors(closerWithError)
	c.True(closerWithError.closed)
}

func TestDiscardAndCloseIgnoringErrors(t *testing.T) {
	c := check.New(t)

	// Test with data to discard
	data := []byte("test data")
	rc := &mockReadCloser{
		data:   data,
		closer: mockCloser{shouldError: false},
	}
	xio.DiscardAndCloseIgnoringErrors(rc)
	c.True(rc.closer.closed)
	c.Equal(len(data), rc.pos)

	// Test with close error (should be ignored)
	rc2 := &mockReadCloser{
		data:   []byte("more data"),
		closer: mockCloser{shouldError: true},
	}
	xio.DiscardAndCloseIgnoringErrors(rc2)
	c.True(rc2.closer.closed)
}

// countingReadCloser serves up to remaining bytes and records how many were actually read.
type countingReadCloser struct {
	remaining int
	read      int
	closed    bool
}

func (r *countingReadCloser) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, io.EOF
	}
	n := min(len(p), r.remaining)
	r.remaining -= n
	r.read += n
	return n, nil
}

func (r *countingReadCloser) Close() error {
	r.closed = true
	return nil
}

func TestDiscardAndCloseIgnoringErrors_ByteBound(t *testing.T) {
	c := check.New(t)
	// A source with far more than the drain byte bound available; only the bound's worth should be read before closing.
	rc := &countingReadCloser{remaining: 4 * 1024 * 1024}
	xio.DiscardAndCloseIgnoringErrors(rc)
	c.True(rc.closed)
	c.Equal(256*1024, rc.read)
	c.True(rc.remaining > 0) // The source was left partially unread rather than drained fully.
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

	// Test successful close (no logging)
	closer := &mockCloser{shouldError: false}
	xio.CloseLoggingErrors(closer)
	c.True(closer.closed)

	// Create a buffer to capture log output
	oldLogger := slog.Default()
	defer func() { slog.SetDefault(oldLogger) }()
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// Test close with error
	closerWithError := &mockCloser{shouldError: true}
	xio.CloseLoggingErrors(closerWithError)
	c.True(closerWithError.closed)
	c.Contains(buf.String(), "mock close error")
}

func TestCloseLoggingAnyErrorTo(t *testing.T) {
	c := check.New(t)

	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Test successful close (no logging)
	closer := &mockCloser{shouldError: false}
	xio.CloseLoggingErrorsTo(logger, closer)
	c.True(closer.closed)
	c.Equal("", buf.String())

	// Test close with error (should log)
	buf.Reset()
	closerWithError := &mockCloser{shouldError: true}
	xio.CloseLoggingErrorsTo(logger, closerWithError)
	c.True(closerWithError.closed)
	logOutput := buf.String()
	c.True(strings.Contains(logOutput, "mock close error"))
}
