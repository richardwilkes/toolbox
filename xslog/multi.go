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
	"context"
	"log/slog"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xos"
)

var _ slog.Handler = &MultiHandler{}

// MultiHandler is a slog.Handler that fans out log records to multiple other handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a new handler that fans out log records to the provided handlers.
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	list := make([]slog.Handler, 0, len(handlers))
	for _, h := range handlers {
		if h != nil {
			if multi, ok := h.(*MultiHandler); ok {
				list = append(list, multi.handlers...)
			} else {
				list = append(list, h)
			}
		}
	}
	return &MultiHandler{handlers: list}
}

// Enabled implements slog.Handler.
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range h.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithGroup implements slog.Handler.
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	handlers := make([]slog.Handler, len(h.handlers))
	for i, one := range h.handlers {
		handlers[i] = one.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

// WithAttrs implements slog.Handler.
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	handlers := make([]slog.Handler, len(h.handlers))
	for i, one := range h.handlers {
		handlers[i] = one.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

// Handle implements slog.Handler interface.
//
//nolint:gocritic // The API cannot be changed
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var result error
	for _, one := range h.handlers {
		if one.Enabled(ctx, r.Level) {
			xos.SafeCall(func() {
				if err := one.Handle(ctx, r.Clone()); err != nil {
					result = errs.Append(result, err)
				}
			}, func(err error) { result = errs.Append(result, err) })
		}
	}
	return result
}
