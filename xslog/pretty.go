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
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xruntime"
	"github.com/richardwilkes/toolbox/v2/xsync"
	"github.com/richardwilkes/toolbox/v2/xterm"
)

var _ slog.Handler = &PrettyHandler{}

// PrettyOptions is used to configure the PrettyHandler.
type PrettyOptions struct {
	slog.HandlerOptions
	ColorSupportOverride xterm.Kind
}

// PrettyHandler is an slog.Handler that outputs a "pretty" format: colorful and supporting formatted stack traces.
type PrettyHandler struct {
	handler          slog.Handler
	sharedBufferLock *sync.Mutex
	buffer           *bytes.Buffer // Protected by sharedBufferLock
	sharedWriterLock *sync.Mutex
	w                io.Writer // Protected by sharedWriterLock
	stack            []string  // Protected by sharedBufferLock
	kind             xterm.Kind
	addSource        bool
}

var poolBuffer = xsync.NewPool(func() []byte { return make([]byte, 0, 1024) })

// NewPrettyHandler creates a new handler with "pretty" output.
func NewPrettyHandler(w io.Writer, opts *PrettyOptions) *PrettyHandler {
	h := &PrettyHandler{
		sharedBufferLock: &sync.Mutex{},
		buffer:           &bytes.Buffer{},
		sharedWriterLock: &sync.Mutex{},
		w:                w,
	}
	var jsonHandlerOpts slog.HandlerOptions
	if opts != nil {
		jsonHandlerOpts = opts.HandlerOptions
		h.kind = opts.ColorSupportOverride
	}
	h.addSource = jsonHandlerOpts.AddSource
	if h.kind == xterm.InvalidKind {
		h.kind = xterm.DetectKind(w)
	}
	next := jsonHandlerOpts.ReplaceAttr
	jsonHandlerOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey ||
			a.Key == slog.SourceKey {
			return slog.Attr{}
		}
		if a.Key == errs.StackTraceKey {
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
	defer func() {
		if cap(buf) <= 8192 { // Only return buffers <= 8 KiB to keep peak allocation low
			buf = buf[:0]
			poolBuffer.Put(buf)
		}
	}()
	buf = h.writeLevel(buf, r.Level)
	buf = h.writeDateTime(buf, r.Time)
	buf = h.writeMessage(buf, r.Message)
	buf = h.writeCaller(buf, r.PC)
	attrs, stack, err := h.collectAttrs(ctx, &r)
	if err != nil {
		return err
	}
	buf = h.writeAttributes(buf, attrs)
	buf = h.writeStack(buf, stack)
	buf = append(buf, '\n')
	h.sharedWriterLock.Lock()
	defer h.sharedWriterLock.Unlock()
	_, err = h.w.Write(buf)
	return err
}

func (h *PrettyHandler) writeDivider(buf []byte) []byte {
	return append(buf, " "+h.kind.Grey()+"|"+h.kind.Reset()+" "...)
}

func (h *PrettyHandler) writeLevel(buf []byte, level slog.Level) []byte {
	var color string
	switch {
	case level < slog.LevelInfo:
		color = h.kind.Cyan()
	case level < slog.LevelWarn:
		color = h.kind.Green()
	case level < slog.LevelError:
		color = h.kind.Color256(214) // Orange
	default:
		color = h.kind.Red()
	}
	buf = append(buf, color...)
	buf = append(buf, level.String()...)
	return append(buf, h.kind.Reset()...)
}

func (h *PrettyHandler) writeDateTime(buf []byte, t time.Time) []byte {
	if !t.IsZero() {
		buf = h.writeDate(buf, t)
		buf = h.writeTime(buf, t)
	}
	return buf
}

func (h *PrettyHandler) writeDate(buf []byte, t time.Time) []byte {
	buf = h.writeDivider(buf)
	buf = append(buf, t.Format("2006-01-02")...)
	return buf
}

func (h *PrettyHandler) writeTime(buf []byte, t time.Time) []byte {
	buf = h.writeDivider(buf)
	buf = append(buf, t.Format("15:04:05.000")...)
	return buf
}

func (h *PrettyHandler) writeMessage(buf []byte, msg string) []byte {
	if msg == "" {
		return buf
	}
	buf = h.writeDivider(buf)
	buf = append(buf, h.kind.Green()...)
	if strings.Contains(msg, "\n") {
		buf = append(buf, strings.ReplaceAll(msg, "\n", "\n    ")...)
	} else {
		buf = append(buf, msg...)
	}
	buf = append(buf, h.kind.Reset()...)
	return buf
}

func (h *PrettyHandler) writeCaller(buf []byte, pc uintptr) []byte {
	if pc == 0 || !h.addSource {
		return buf
	}
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	file := xruntime.StackTracePath(f.Function, f.File)
	if file == "" {
		return buf
	}
	buf = h.writeDivider(buf)
	return append(buf, h.kind.Dim()+h.kind.Yellow()+file+":"+strconv.Itoa(f.Line)+h.kind.Reset()...)
}

func (h *PrettyHandler) collectAttrs(ctx context.Context, r *slog.Record) (textAttr string, stack []string, err error) {
	h.sharedBufferLock.Lock()
	defer func() {
		h.stack = nil
		h.buffer.Reset()
		h.sharedBufferLock.Unlock()
	}()
	if err = h.handler.Handle(ctx, *r); err != nil {
		return "", nil, err
	}
	return strings.TrimRight(h.buffer.String(), "\n"), h.stack, nil
}

func (h *PrettyHandler) writeAttributes(buf []byte, attrs string) []byte {
	if attrs == "" || attrs == "{}" {
		return buf
	}
	buf = h.writeDivider(buf)
	buf = append(buf, attrs...)
	return buf
}

func (h *PrettyHandler) writeStack(buf []byte, stack []string) []byte {
	if len(stack) == 0 {
		return buf
	}
	buf = append(buf, "\n    "...)
	buf = append(buf, h.kind.Dim()...)
	buf = append(buf, h.kind.Yellow()...)
	buf = append(buf, strings.Join(stack, "\n    ")...)
	buf = append(buf, h.kind.Reset()...)
	return buf
}

// Enabled implements slog.Handler interface.
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs implements slog.Handler interface.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &PrettyHandler{
		handler:          h.handler.WithAttrs(attrs),
		sharedBufferLock: h.sharedBufferLock,
		buffer:           h.buffer,
		sharedWriterLock: h.sharedWriterLock,
		w:                h.w,
	}
}

// WithGroup implements slog.Handler interface.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &PrettyHandler{
		handler:          h.handler.WithGroup(name),
		sharedBufferLock: h.sharedBufferLock,
		buffer:           h.buffer,
		sharedWriterLock: h.sharedWriterLock,
		w:                h.w,
	}
}
