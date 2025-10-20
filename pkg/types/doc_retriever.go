package types

import (
	"context"
	"log/slog"

	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type DocRetriver interface {
	GetDocuments(context.Context, *slog.Logger) (chan *vecdb.EmbeddDocument, error)
}
