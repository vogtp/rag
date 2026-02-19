package model

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/amikos-tech/chroma-go/types"
	lcEmbedd "github.com/vogtp/langchaingo/embeddings"
)

func NewEmbedder(m GenericModel) (types.EmbeddingFunction, error) {
	ef, err := lcEmbedd.NewEmbedder(m)
	if err != nil {
		return nil, err
	}
	return &embedWrapper{
		e: ef,
	}, nil
}

type embedWrapper struct {
	e lcEmbedd.Embedder
}

// EmbedDocuments returns a vector for each text.
func (ew embedWrapper) EmbedDocuments(ctx context.Context, texts []string) ([]*types.Embedding, error) {

	ebds, err := ew.e.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}
	ret := make([]*types.Embedding, len(ebds))
	for i, e := range ebds {
		ret[i] = types.NewEmbeddingFromFloat32(e)
	}
	return ret, nil
}

// EmbedQuery embeds a single text.
func (ew embedWrapper) EmbedQuery(ctx context.Context, text string) (*types.Embedding, error) {
	e, err := ew.e.EmbedQuery(ctx, text)
	if err != nil {
		if ctx.Err() != nil {
			return nil, err
		}
		//FIXME sort of hackish error handling
		d := time.Second
		if strings.Contains(err.Error(), "status code: 429") {
			d = time.Second * 10
		}
		slog.Info("EmbedQuery got error, will retry shortly", "err", err, "delay", d.String())
		time.Sleep(d)
		e, err = ew.e.EmbedQuery(ctx, text)
	}
	if err != nil {
		return nil, err
	}
	return types.NewEmbeddingFromFloat32(e), nil
}

func (ew embedWrapper) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	var err error
	var e *types.Embedding
	for i, r := range records {
		e, err = ew.EmbedQuery(ctx, r.Document)
		if err != nil {
			slog.Warn("EmbedRecords got error embedding", "err", err)
			if !force {
				return err
			}
		}
		r.Embedding = *e
		records[i] = r
	}
	return err
}
