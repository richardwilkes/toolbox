// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
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
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

var _ slog.Handler = &Handler{}

// Handler provides a formatted text output that may include a stack trace on separate lines. The stack trace is
// formatted such that most IDEs will auto-generate links for it within their consoles. Note that this slog.Handler is
// not optimized for performance, as I expect those that need to run this is environments where that matters will use
// one of the implementations provided by slog itself.
type Handler struct {
	level slog.Leveler
	lock  *sync.Mutex
	out   io.Writer
	list  []entry
}

type entry struct {
	group string
	attrs []slog.Attr
}

type embeddedStackError interface {
	StackError() errs.StackError
}

// New creates a new Handler. Only log levels >= the provided level will be emitted.
func New(w io.Writer, level slog.Leveler) *Handler {
	if level == nil {
		level = slog.LevelInfo
	}
	return &Handler{
		level: level,
		lock:  &sync.Mutex{},
		out:   w,
	}
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
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var buffer bytes.Buffer
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

	h.lock.Lock()
	defer h.lock.Unlock()
	_, err := h.out.Write(buffer.Bytes())
	return err
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
		if embedded, ok := attr.Value.Any().(embeddedStackError); ok {
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
