package vecdb

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	ollamaEmbedd "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
)

// VecDB is a wrapper of a vectoDB
type VecDB struct {
	slog            *slog.Logger
	chromaAddr      string
	chroma          *chroma.Client
	embedFunc       *ollamaEmbedd.OllamaEmbeddingFunction
	ollamaAddr      string
	embeddingsModel string
}

// New creates a vectorDB
func New(ctx context.Context, slog *slog.Logger, opts ...Option) (*VecDB, error) {
	v := &VecDB{
		slog:       slog,
		chromaAddr: cfg.ChromaUrl(),
		//embeddingsModel: "nomic-embed-text",
		embeddingsModel: viper.GetString(cfg.ModelEmbedding),
	}
	for _, o := range opts {
		o(v)
	}
	if len(v.ollamaAddr) < 1 {
		v.ollamaAddr = cfg.GetOllamaHost(ctx)
	}
	if len(v.ollamaAddr) < 1 {
		return nil, fmt.Errorf("no running ollama found: %q", v.ollamaAddr)
	}
	if len(v.chromaAddr) < 1 {
		return nil, fmt.Errorf("no chroma address given")
	}
	v.slog = slog.With("chroma_addr", v.chromaAddr, "ollama_addr", v.ollamaAddr)

	client, err := chroma.NewClient(v.chromaAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}
	v.slog.Debug("Connected to chroma")
	v.chroma = client
	return v, nil
}

// CreateCollection create a collection
func (v *VecDB) CreateCollection(ctx context.Context, name string, metadata map[string]interface{}) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return nil, err
	}
	return v.chroma.CreateCollection(ctx, name, nil, true, embedFunc, types.L2)
}

// GetCollection returns a collection
func (v *VecDB) GetCollection(ctx context.Context, name string) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return nil, err
	}
	return v.chroma.GetCollection(ctx, name, embedFunc)
}

// GetEmbeddingFunc load the embedding function from the llm
func (v *VecDB) GetEmbeddingFunc() (*ollamaEmbedd.OllamaEmbeddingFunction, error) {
	if v.embedFunc != nil {
		return v.embedFunc, nil
	}
	v.slog.Debug("Loading embedding function", "embeddingsModel", v.embeddingsModel)
	embedFunc, err := ollamaEmbedd.NewOllamaEmbeddingFunction(
		ollamaEmbedd.WithBaseURL(v.ollamaAddr),
		ollamaEmbedd.WithModel(v.embeddingsModel),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating ollama embedding function: %w", err)
	}
	v.embedFunc = embedFunc
	return embedFunc, nil
}

// DeleteCollection delete a collection
func (v *VecDB) DeleteCollection(ctx context.Context, collectionName string) error {
	_, err := v.chroma.DeleteCollection(ctx, collectionName)
	if err != nil {
		err = fmt.Errorf("cannot delete collection %s: %w", collectionName, err)
	}
	return err
}

// ListCollections lists all colletions
func (v *VecDB) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	prefix := viper.GetString(cfg.VecDBColName)
	cols, err := v.chroma.ListCollections(ctx)
	if err != nil {
		return nil, err
	}
	collections := make([]*chroma.Collection, 0, len(cols))
	for _, col := range cols {
		if !strings.HasPrefix(col.Name, prefix) {
			continue
		}
		collections = append(collections, col)
	}
	return collections, nil
}
