package web

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/vogtp/rag/pkg/cfg"
)

func getOllamaClient(ctx context.Context, model string) (*ollama.LLM, error) {
	url := cfg.GetOllamaHost(ctx)
	slog.Debug("connecting to ollama", "model", model, "url", url)
	return ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(url),
	)
}

func getEmbedding(ctx context.Context, model string) (llms.Model, *embeddings.EmbedderImpl, error) {
	llm, err := getOllamaClient(ctx, model)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create ollama client: %w", err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create ollama embedder: %w", err)
	}
	return llms.Model(llm), e, nil
}

func getDocs(ctx context.Context, index string, query string) ([]schema.Document, error) {

	_, e, err := getEmbedding(ctx, "mxbai-embed-large")
	if err != nil {
		return nil, err
	}
	store, err := chroma.New(
		chroma.WithChromaURL(cfg.ChromaUrl()),
		chroma.WithNameSpace(index),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction(types.COSINE),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create chroma client: %w", err)
	}
	return store.SimilaritySearch(ctx, query, 3, vectorstores.WithScoreThreshold(0.3))
}
