package types

import (
	"context"
	"log/slog"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/web/bearer"
)

type Instance interface {
	Displayname() string

	Model(name string) (Model, error)

	//FIXME models and list collections should be the same
	Models(ctx context.Context) []Model
	ListCollections(ctx context.Context, slog *slog.Logger) ([]*chroma.Collection, error)

	LLM() string

	UpdateIntervall() time.Duration
	Embbed(ctx context.Context, slog *slog.Logger) error
	SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error)

	// CollectionName() string
	// Confluence() *cfg.ConfluenceCfg
	// ModelEmbedding() string
	// VecDBUpdateIntervall() time.Duration
	bearer.Auth
}
