package vecdb

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/model"
)

// VecDB is a wrapper of a vectoDB
type VecDB struct {
	slog            *slog.Logger
	chromaAddr      string
	chroma          *chroma.Client
	embedFunc       types.EmbeddingFunction
	embeddingsModel string
}

// New creates a vectorDB
func New(ctx context.Context, slog *slog.Logger, embeddingsModel string, opts ...Option) (*VecDB, error) {
	v := &VecDB{
		slog:            slog,
		chromaAddr:      cfg.ChromaUrl(),
		embeddingsModel: embeddingsModel,
	}
	for _, o := range opts {
		o(v)
	}
	if len(v.chromaAddr) < 1 {
		return nil, fmt.Errorf("no chroma address given")
	}
	v.slog = slog.With("chroma_addr", v.chromaAddr)

	client, err := chroma.NewClient(chroma.WithBasePath(v.chromaAddr))
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}
	v.slog.Debug("Connected to chroma")
	v.chroma = client
	return v, nil
}

// CreateCollection create a collection
func (v *VecDB) CreateCollection(ctx context.Context, name string, metadata map[string]interface{}) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc(ctx)
	if err != nil {
		return nil, err
	}

	// types.IP     -> doc 0 dist: 0.3435346
	// types.COSINE -> doc 0 dist: 0.34353453
	// types.L2     -> doc 0 dist: 0.68706906
	return v.chroma.CreateCollection(ctx, name, nil, true, embedFunc, types.COSINE)
}

// GetCollection returns a collection
func (v *VecDB) GetCollection(ctx context.Context, name string) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc(ctx)
	if err != nil {
		return nil, err
	}
	return v.chroma.GetCollection(ctx, name, embedFunc)
}

// GetEmbeddingFunc load the embedding function from the llm
func (v *VecDB) GetEmbeddingFunc(ctx context.Context) (types.EmbeddingFunction, error) {
	genModel, err := model.GetBackendModel(ctx, v.slog, v.embeddingsModel)
	if err != nil {
		return nil, fmt.Errorf("get backend models: %w", err)
	}
	embedder, err := model.NewEmbedder(genModel)
	if err != nil {
		return nil, fmt.Errorf("creating embedder from model: %w", err)
	}
	return embedder, nil
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
func (v *VecDB) ListAllCollections(ctx context.Context) ([]*chroma.Collection, error) {
	return v.chroma.ListCollections(ctx)
}

// ListCollections lists all colletions
func (v *VecDB) ListCollections(ctx context.Context, collectionName string) ([]*chroma.Collection, error) {
	if len(collectionName) < 1 {
		return nil, fmt.Errorf("no collection name given")

	}
	cols, err := v.ListAllCollections(ctx)
	if err != nil {
		return nil, err
	}
	prefix := collectionName
	collections := make([]*chroma.Collection, 0, len(cols))
	for _, col := range cols {
		if !strings.HasPrefix(col.Name, prefix) {
			continue
		}
		collections = append(collections, col)
	}
	return collections, nil
}
