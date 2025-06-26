// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package multilog

import (
	"context"
	"log/slog"

	"github.com/richardwilkes/toolbox/errs"
)

var _ slog.Handler = &Handler{}

// Handler is a slog.Handler that fans out log records to multiple other handlers.
type Handler struct {
	handlers []slog.Handler
}

// New creates a new Handler that fans out log records to the provided handlers.
func New(handlers ...slog.Handler) *Handler {
	return &Handler{handlers: handlers}
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range h.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	handlers := make([]slog.Handler, len(h.handlers))
	for i, one := range h.handlers {
		handlers[i] = one.WithGroup(name)
	}
	return &Handler{handlers: handlers}
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, one := range h.handlers {
		handlers[i] = one.WithAttrs(attrs)
	}
	return &Handler{handlers: handlers}
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic // Must use defined API
	var result error
	for _, one := range h.handlers {
		if one.Enabled(ctx, r.Level) {
			result = errs.Append(result, runHandler(ctx, &r, one))
		}
	}
	return result
}

func runHandler(ctx context.Context, r *slog.Record, h slog.Handler) (err error) {
	defer errs.Recovery(func(rerr error) { err = rerr })
	err = h.Handle(ctx, r.Clone())
	return err
}
