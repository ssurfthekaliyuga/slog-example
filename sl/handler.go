package sl

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := extractAttrs(ctx)

	for _, attr := range attrs {
		record.AddAttrs(attr)
	}

	return h.Handler.Handle(ctx, record)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler { //todo важно
	return &ContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler { //todo важно
	return &ContextHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
