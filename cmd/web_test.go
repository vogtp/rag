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
	logger.New()

	if err := startWeb(t.Context()); err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}
