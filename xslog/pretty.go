// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xslog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/richardwilkes/toolbox/v2/xio/term"
	"github.com/richardwilkes/toolbox/v2/xruntime"
	"github.com/richardwilkes/toolbox/v2/xsync"
)

var _ slog.Handler = &PrettyHandler{}

// StackKey is the key for the stack trace attribute.
const StackKey = "stack" // Keep in sync with errs.StackTraceKey

// PrettyHandler is an slog.Handler that outputs a "pretty" format: colorful and supporting formatted stack traces.
type PrettyHandler struct {
	handler    slog.Handler
	sharedLock *sync.Mutex
	buffer     *bytes.Buffer
	w          io.Writer
	stack      []string
	kind       term.Kind
}

// PrettyOptions is used to configure the PrettyHandler.
type PrettyOptions struct {
	slog.HandlerOptions
	ColorSupportOverride term.Kind
}

var poolBuffer = xsync.NewPool[[]byte](func() []byte { return make([]byte, 0, 1024) })

// NewPrettyHandler creates a new handler with "pretty" output.
func NewPrettyHandler(w io.Writer, opts *PrettyOptions) *PrettyHandler {
	h := &PrettyHandler{
		sharedLock: &sync.Mutex{},
		buffer:     &bytes.Buffer{},
		w:          w,
	}
	var jsonHandlerOpts slog.HandlerOptions
	if opts != nil {
		jsonHandlerOpts = opts.HandlerOptions
		h.kind = opts.ColorSupportOverride
	}
	if h.kind == term.InvalidKind {
		h.kind = term.DetectKind(w)
	}
	next := jsonHandlerOpts.ReplaceAttr
	jsonHandlerOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey ||
			a.Key == slog.SourceKey {
			return slog.Attr{}
		}
		if a.Key == StackKey {
			if s, ok := a.Value.Any().([]string); ok {
				h.stack = s
			}
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
	h.handler = slog.NewJSONHandler(h.buffer, &jsonHandlerOpts)
	return h
}

// Handle implements slog.Handler interface.
//
//nolint:gocritic // The API cannot be changed
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := poolBuffer.Get()
	defer poolBuffer.Put(buf)
	buf = h.writeLevel(buf, r.Level)
	if !r.Time.IsZero() {
		buf = h.writeDivider(buf)
		buf = append(buf, r.Time.Format("2006-01-02")...)
		buf = h.writeDivider(buf)
		buf = append(buf, r.Time.Format("15:04:05.000")...)
	}
	buf = h.writeCaller(buf, r.PC)
	if r.Message != "" {
		buf = h.writeDivider(buf)
		buf = append(buf, h.kind.Green()...)
		buf = append(buf, r.Message...)
		buf = append(buf, h.kind.Reset()...)
	}
	attrs, stack, err := h.collectAttrs(ctx, &r)
	if err != nil {
		return err
	}
	if attrs != "" && attrs != "{}" {
		buf = h.writeDivider(buf)
		buf = append(buf, attrs...)
	}
	if len(stack) > 0 {
		buf = append(buf, '\n', '\t')
		buf = append(buf, h.kind.Dim()...)
		buf = append(buf, h.kind.Yellow()...)
		buf = append(buf, strings.Join(stack, "\n\t")...)
		buf = append(buf, h.kind.Reset()...)
	}
	buf = append(buf, '\n')
	h.sharedLock.Lock()
	defer h.sharedLock.Unlock()
	_, err = h.w.Write(buf)
	return err
}

func (h *PrettyHandler) writeDivider(buf []byte) []byte {
	return append(buf, " "+h.kind.Grey()+"|"+h.kind.Reset()+" "...)
}

func (h *PrettyHandler) writeLevel(buf []byte, level slog.Level) []byte {
	var prefix string
	var base slog.Level
	switch {
	case level < slog.LevelInfo:
		prefix = h.kind.Cyan() + "DBG"
		base = slog.LevelDebug
	case level < slog.LevelWarn:
		prefix = h.kind.Green() + "INF"
		base = slog.LevelInfo
	case level < slog.LevelError:
		prefix = h.kind.Color256(214) + "WRN" // Orange
		base = slog.LevelWarn
	default:
		prefix = h.kind.Red() + "ERR"
		base = slog.LevelError
	}
	buf = append(buf, prefix...)
	if val := int(level - base); val != 0 {
		if val >= 0 {
			buf = append(buf, '+')
		}
		buf = append(buf, strconv.Itoa(val)...)
	}
	return append(buf, h.kind.Reset()...)
}

func (h *PrettyHandler) writeCaller(buf []byte, pc uintptr) []byte {
	if pc == 0 {
		return buf
	}
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	file := xruntime.StackTracePath(f.File)
	if file == "" {
		return buf
	}
	buf = h.writeDivider(buf)
	return append(buf, h.kind.Dim()+h.kind.Yellow()+file+":"+strconv.Itoa(f.Line)+h.kind.Reset()...)
}

func (h *PrettyHandler) collectAttrs(ctx context.Context, r *slog.Record) (textAttr string, stack []string, err error) {
	h.sharedLock.Lock()
	defer func() {
		h.stack = nil
		h.buffer.Reset()
		h.sharedLock.Unlock()
	}()
	if err = h.handler.Handle(ctx, *r); err != nil {
		return "", nil, err
	}
	return strings.TrimRight(h.buffer.String(), "\n"), h.stack, nil
}

// Enabled implements slog.Handler interface.
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs implements slog.Handler interface.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		handler:    h.handler.WithAttrs(attrs),
		sharedLock: h.sharedLock,
		buffer:     h.buffer,
		w:          h.w,
	}
}

// WithGroup implements slog.Handler interface.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		handler:    h.handler.WithGroup(name),
		sharedLock: h.sharedLock,
		buffer:     h.buffer,
		w:          h.w,
	}
}
