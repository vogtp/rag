package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	chromaDefaultPort           = "8000"
	chromaDefaultContainerImage = "chromadb/chroma:0.5.23"
)

var chromaDefaultURL string

func init() {
	chromaDefaultURL = fmt.Sprintf("http://localhost:%v", chromaDefaultPort)
}

func ChromaUrl() string {
	port := viper.GetString(ChromaContainerPort)
	url := viper.GetString(chromaUrl)
	if port != chromaDefaultPort && url == chromaDefaultURL {
		return fmt.Sprintf("http://localhost:%v", port)
	}
	return url
}
