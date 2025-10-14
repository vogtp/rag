package logger

import (
	"io"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
)

func New() *slog.Logger {
	return Create(Level())
}

func Level() slog.Level {
	return LevelFromString(viper.GetString(cfg.LogLevel))
}

func Create(lvl slog.Level) *slog.Logger {
	logOpts := slog.HandlerOptions{
		Level: lvl,
	}
	logOpts.AddSource = viper.GetBool(cfg.LogSource)
	logJson := viper.GetBool(cfg.LogJson)
	if !(logJson || logOpts.AddSource) {
		slog.SetLogLoggerLevel(lvl)
		return slog.Default()
	}
	if logOpts.AddSource {
		logOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey && len(groups) == 0 {
				return ProcessSourceField(a, logJson)
			}
			return a
		}
	}
	var logWriter io.Writer = os.Stdout
	var handler slog.Handler
	handler = slog.NewTextHandler(logWriter, &logOpts)
	if logJson {
		handler = slog.NewJSONHandler(logWriter, &logOpts)
	}
	sl := slog.New(handler)
	slog.SetDefault(sl)
	return sl
}

func LevelFromString(levelStr string) slog.Level {
	// We don't care about case. Accept both "INFO" and "info".
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))
	switch levelStr {
	case "trace":
		return slog.LevelDebug
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "off":
		return slog.Level(88)
	default:
		return slog.LevelInfo
	}
}
