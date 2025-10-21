package cfg

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	ollamaAPI "github.com/ollama/ollama/api"
	"github.com/spf13/viper"
)

var _ollamaHost string

// GetOllamaHost returns a active ollama host
func GetOllamaHost(ctx context.Context, slog *slog.Logger) string {
	if checkOllama(ctx, _ollamaHost) {
		slog.Debug("Using prevously found ollala host", "url", _ollamaHost)
		return _ollamaHost
	}
	for _, o := range viper.GetStringSlice(ollamaHosts) {
		if checkOllama(ctx, o) {
			_ollamaHost = o
			slog.Info("Found running ollama", "url", o)
			return _ollamaHost
		}
	}
	slog.Warn("No running ollama found")
	return ""
}

func checkOllama(ctx context.Context, urlStr string) bool {
	slog.Debug("Checking ollama host", "url", urlStr)
	if len(urlStr) < 7 { // 7 -> http://
		return false
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		slog.Warn("Cannot parse ollama url", "url", urlStr, "err", err)
		return false
	}
	c := ollamaAPI.NewClient(u, http.DefaultClient)
	if err := c.Heartbeat(ctx); err != nil {
		slog.Warn("Cannot connect to ollama", "url", u, "err", err)
		return false
	}
	return true
}
