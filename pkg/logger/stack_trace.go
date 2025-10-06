package logger

import (
	"log/slog"
	"regexp"
	"runtime"
	"strings"
)

var (
	reTrace = regexp.MustCompile(`.*/slog/logger\.go.*\n`)
	logkey  = "stacktrace"
)

func Stacktrace() slog.Attr {
	stackInfo := make([]byte, 1024*1024)
	stackSize := runtime.Stack(stackInfo, false)
	if stackSize < 1 {
		return slog.String(logkey, "no trace found")
	}
	traceLines := reTrace.Split(string(stackInfo[:stackSize]), -1)
	if len(traceLines) == 0 {
		return slog.String(logkey, "no trace found")
	}
	t:=strings.ReplaceAll(traceLines[len(traceLines)-1], "\t","  ")
	lines:=strings.Split( t, "\n")
	return slog.Any(logkey, lines[3:])

}
