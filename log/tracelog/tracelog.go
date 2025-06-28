// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tracelog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xreflect"
)

var _ slog.Handler = &Handler{}

// Config is used to configure a Handler.
type Config struct {
	// Level is the minimum log level that will be emitted. Defaults to slog.LevelInfo if not set.
	Level slog.Leveler
	// LevelNames is an optional map of slog.Level to string that will be used to format the log level. If not set, the
	// default slog level names will be used.
	LevelNames map[slog.Level]string
	// Sink is the io.Writer that will receive the formatted log output. Defaults to os.Stderr if not set.
	Sink io.Writer
	// BufferDepth greater than 0 will enable asynchronous delivery of log messages to the sink. When enabled, if there
	// is no room remaining in the buffer, the message will be discarded rather than waiting for room to become
	// available. Defaults to 0 for synchronous delivery.
	BufferDepth int
}

// Normalize ensures that the Config is valid.
func (c *Config) Normalize() {
	if xreflect.IsNil(c.Level) {
		c.Level = slog.LevelInfo
	}
	if c.Sink == nil {
		c.Sink = os.Stderr
	}
	if c.BufferDepth < 0 {
		c.BufferDepth = 0
	}
}

// Handler provides a formatted text output that may include a stack trace on separate lines. The stack trace is
// formatted such that most IDEs will auto-generate links for it within their consoles. Note that this slog.Handler is
// not optimized for performance, as I expect those that need to run this is environments where that matters will use
// one of the implementations provided by slog itself.
type Handler struct {
	level      slog.Leveler
	levelNames map[slog.Level]string
	delivery   chan []byte
	lock       *sync.Mutex
	sink       io.Writer
	list       []entry
}

type entry struct {
	group string
	attrs []slog.Attr
}

// New creates a new Handler. May pass nil for cfg to use the defaults.
func New(cfg *Config) *Handler {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.Normalize()
	h := Handler{
		level:      cfg.Level,
		levelNames: cfg.LevelNames,
		lock:       &sync.Mutex{},
		sink:       cfg.Sink,
	}
	if cfg.BufferDepth > 0 {
		h.delivery = make(chan []byte, cfg.BufferDepth)
		go h.backgroundDelivery()
	}
	return &h
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(entry{group: name})
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(entry{attrs: attrs})
}

func (h *Handler) withGroupOrAttrs(ga entry) *Handler {
	other := *h
	other.list = make([]entry, len(h.list)+1)
	copy(other.list, h.list)
	other.list[len(other.list)-1] = ga
	return &other
}

// Handle implements slog.Handler.
func (h *Handler) Handle(_ context.Context, r slog.Record) error { //nolint:gocritic // Must use defined API
	var buffer bytes.Buffer
	if name, ok := h.levelNames[r.Level]; ok {
		buffer.WriteString(name)
	} else {
		switch r.Level {
		case slog.LevelDebug:
			buffer.WriteString("DBG")
		case slog.LevelInfo:
			buffer.WriteString("INF")
		case slog.LevelWarn:
			buffer.WriteString("WRN")
		case slog.LevelError:
			buffer.WriteString("ERR")
		default:
			fmt.Fprintf(&buffer, "%3d", r.Level)
		}
	}
	buffer.WriteString(r.Time.Round(0).Format(" | 2006-01-02 | 15:04:05.000 | "))
	buffer.WriteString(r.Message)

	s := &state{buffer: &buffer, needBar: true}
	for _, ga := range h.list {
		s.append(ga)
	}
	r.Attrs(func(attr slog.Attr) bool {
		s.appendAttr(attr)
		return true
	})
	buffer.WriteByte('\n')
	if s.stackErr != nil {
		buffer.WriteString(s.stackErr.StackTrace(true))
		buffer.WriteByte('\n')
	}
	if h.delivery != nil {
		select {
		case h.delivery <- buffer.Bytes():
		default:
		}
		return nil
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	_, err := h.sink.Write(buffer.Bytes())
	return err
}

func (h *Handler) backgroundDelivery() {
	for data := range h.delivery {
		_, _ = h.sink.Write(data) //nolint:errcheck // We don't care about errors here
	}
}

type state struct {
	buffer   *bytes.Buffer
	stackErr errs.StackError
	group    string
	needBar  bool
}

func (s *state) append(ga entry) {
	if ga.group != "" {
		s.addGroup(ga.group)
		return
	}
	for _, attr := range ga.attrs {
		s.appendAttr(attr)
	}
}

func (s *state) appendAttr(attr slog.Attr) {
	if s.group == "" && attr.Key == errs.StackTraceKey {
		if embedded, ok := attr.Value.Any().(interface{ StackError() errs.StackError }); ok {
			s.stackErr = embedded.StackError()
			return
		}
	}
	attr.Value = attr.Value.Resolve()
	if !attr.Equal(slog.Attr{}) {
		switch attr.Value.Kind() {
		case slog.KindString:
			s.addBarIfNeeded()
			s.writeGroupAndKey(attr.Key)
			_, _ = fmt.Fprintf(s.buffer, "%q", attr.Value.String())
		case slog.KindTime:
			s.addBarIfNeeded()
			s.writeGroupAndKey(attr.Key)
			_, _ = s.buffer.WriteString(attr.Value.Time().Format(time.RFC3339Nano))
		case slog.KindGroup:
			attrs := attr.Value.Group()
			if len(attrs) == 0 {
				return
			}
			s.addBarIfNeeded()
			savedPrefix := s.group
			s.addGroup(attr.Key)
			for _, attr = range attrs {
				s.appendAttr(attr)
			}
			s.group = savedPrefix
		default:
			s.addBarIfNeeded()
			s.writeGroupAndKey(attr.Key)
			_, _ = s.buffer.WriteString(attr.Value.String())
		}
	}
}

func (s *state) writeGroupAndKey(key string) {
	_ = s.buffer.WriteByte(' ')
	_, _ = s.buffer.WriteString(s.group)
	_, _ = s.buffer.WriteString(key)
	_ = s.buffer.WriteByte('=')
}

func (s *state) addGroup(group string) {
	group += "."
	if s.group == "" {
		s.group = group
	} else {
		s.group += group
	}
}

func (s *state) addBarIfNeeded() {
	if s.needBar {
		s.buffer.WriteString(" |")
		s.needBar = false
	}
}
