package test

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/vogtp/rag/pkg/logger"
)

func getSlog() *slog.Logger {
	var logWriter io.Writer = os.Stdout
	logOpts := slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey && len(groups) == 0 {
				return logger.ProcessSourceField(a, true)
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(logWriter, &logOpts)
	testHanlder := testSlogHandler{
		Handler: handler,
	}
	return slog.New(&testHanlder)
}

type testSlogHandler struct {
	slog.Handler
}

func (h *testSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *testSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &testSlogHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *testSlogHandler) WithGroup(name string) slog.Handler {
	return &testSlogHandler{Handler: h.Handler.WithGroup(name)}
}
func (h *testSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// return nil
	return h.Handler.Handle(ctx, r)
}
