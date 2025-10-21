package test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/vogtp/rag/pkg/logger"
)

func getSlog() *slog.Logger {
	var logWriter io.Writer = os.Stderr
	logOpts := slog.HandlerOptions{
		Level:     slog.LevelInfo,
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
		w:       logWriter,
		doDot:   true,
	}
	logFile, err := os.Create("ignore_test.log")
	if err == nil {
		testHanlder.logFileHandler = slog.NewJSONHandler(logFile, &logOpts)
		testHanlder.logFile = logFile
	} else {
		fmt.Printf("Cannot create test log: %v", err)
	}
	return slog.New(&testHanlder)
}

type testSlogHandler struct {
	slog.Handler
	doDot          bool
	logFileHandler slog.Handler
	logFile        *os.File
	w              io.Writer
}

func (h testSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h testSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &testSlogHandler{Handler: h.Handler.WithAttrs(attrs), w: h.w}
}

func (h testSlogHandler) WithGroup(name string) slog.Handler {
	return &testSlogHandler{Handler: h.Handler.WithGroup(name), w: h.w}
}

func (h testSlogHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	if h.doDot {
		switch r.Level {
		case slog.LevelInfo:
			_, err = h.w.Write([]byte("."))
		case slog.LevelWarn:
			_, err = h.w.Write([]byte("!"))
		}
	}
	//write log file
	if h.logFileHandler != nil {
		if err := h.logFileHandler.Handle(ctx, r); err != nil {
			fmt.Fprintf(h.logFile, "logFileHandler err: %v", err)
		}
		slog.New(h.logFileHandler).Error("Start")
	}
	if r.Level == slog.LevelError {
		return h.Handler.Handle(ctx, r)
	}
	return
}
