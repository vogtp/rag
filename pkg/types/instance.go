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
	Name() string

	Model(name string) (Model, error)
	Models(ctx context.Context) []Model

	LLM() string

	UpdateIntervall() time.Duration
	Embbed(ctx context.Context) error
	ListCollections(ctx context.Context) ([]*chroma.Collection, error)
	SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error)

	bearer.Auth
}
