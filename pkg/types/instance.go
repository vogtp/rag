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
	DisplayName() string

	Model(name string) (Model, error)

	Models(ctx context.Context) []Model
	ListCollections(ctx context.Context, slog *slog.Logger, ) ([]*chroma.Collection, error)

	LLM() string

	UpdateIntervall() time.Duration
	Embbed(ctx context.Context, slog *slog.Logger, filters ...vecdb.Filter) error
	SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error)

	DocRetriver

	bearer.Auth
}
