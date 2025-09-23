package vecdb

import (
	"context"
	"log/slog"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/amikos-tech/chroma-go/types"
)

type embeddingFunctionWrapper struct {
	ef embeddings.EmbeddingFunction
}

// EmbedDocuments returns a vector for each text.
func (efw embeddingFunctionWrapper) EmbedDocuments(ctx context.Context, texts []string) ([]*types.Embedding, error) {
	ebds, err := efw.ef.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}
	ret := make([]*types.Embedding, len(ebds))
	for i, e := range ebds {
		ret[i] = types.NewEmbeddingFromFloat32(e.ContentAsFloat32())
	}
	return ret, nil
}

// EmbedQuery embeds a single text.
func (efw embeddingFunctionWrapper) EmbedQuery(ctx context.Context, text string) (*types.Embedding, error) {
	e, err := efw.ef.EmbedQuery(ctx, text)
	if err != nil {
		return nil, err
	}
	return types.NewEmbeddingFromFloat32(e.ContentAsFloat32()), nil
}

func (efw embeddingFunctionWrapper) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	var err error
	var e *types.Embedding
	for i, r := range records {
		e, err = efw.EmbedQuery(ctx, r.Document)
		if err != nil {
			slog.Warn("EmedRecords got error embedding", "err", err)
			if !force {
				return err
			}
		}
		r.Embedding = *e
		records[i] = r
	}
	return err
}
