package model

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vogtp/langchaingo/embeddings"
	"github.com/vogtp/langchaingo/llms"
	"github.com/vogtp/langchaingo/llms/ollama"
	llmopenai "github.com/vogtp/langchaingo/llms/openai"
	"github.com/vogtp/rag/pkg/cfg"
)

type GenericModel interface {
	llms.Model
	llms.ReasoningModel
	embeddings.EmbedderClient
}

func GetBackendModel(ctx context.Context, slog *slog.Logger, llmName string) (GenericModel, error) {
	backends, err := cfg.GetBackends()
	if err != nil {
		return nil, fmt.Errorf("get backends: %w", err)
	}
	bm := backends.Model(llmName)
	if bm == nil {
		return nil, fmt.Errorf("no backend model %q found", llmName)
	}
	switch bm.API.APIType {
	case cfg.BackendApiTypeOllama:
		return getOllamaClient(ctx, slog, bm)
	case cfg.BackendApiTypeOpenai:
		return getOpenAIClient(ctx, slog, bm)
	default:
		slog.Error("Cannot get the backend model, api type not found", "APIType", bm.API.Type)
		return nil, fmt.Errorf("API type %q not found", bm.API.Type)
	}
}

func getOllamaClient(_ context.Context, slog *slog.Logger, bm *cfg.BackendModel) (GenericModel, error) {
	slog.Info("connecting to ollama", "LLM", bm.Name, "url", bm.API.URL)
	opts := []ollama.Option{
		ollama.WithModel(bm.Name),
		ollama.WithServerURL(bm.API.URL),
	}
	// if len(bm.API.Key) > 0 {
	// 	opts = append(opts, ollama.With)
	// }
	return ollama.New(opts...)
}

func getOpenAIClient(_ context.Context, slog *slog.Logger, bm *cfg.BackendModel) (GenericModel, error) {
	slog.Info("connecting to an openAI compatible API", "LLM", bm.Name, "url", bm.API.URL)
	opts := []llmopenai.Option{
		llmopenai.WithModel(bm.Name),
		llmopenai.WithBaseURL(bm.API.URL),
	}
	if len(bm.API.Key) > 0 {
		opts = append(opts, llmopenai.WithToken(bm.API.Key))
	}
	return llmopenai.New(opts...)
}
