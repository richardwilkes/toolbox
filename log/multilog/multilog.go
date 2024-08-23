package multilog

import (
	"context"
	"log/slog"

	"github.com/richardwilkes/toolbox/errs"
)

var _ slog.Handler = &Handler{}

// Handler is a slog.Handler that sends log records to multiple other handlers.
type Handler struct {
	handlers []slog.Handler
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
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
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
