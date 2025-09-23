package cmd

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
)

func Test_WebServer(t *testing.T) {
	viper.AddConfigPath("..")
	cfg.Parse()
	viper.Set(cfg.LogJson, true)
	viper.Set(cfg.LogLevel, "warn")
	logger.New()

	if err := startWeb(t.Context()); err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}
