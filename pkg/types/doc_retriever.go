package types

import (
	"context"
	"log/slog"
)

type DocRetriver interface {
	GetDocuments(context.Context, *slog.Logger) (chan *EmbeddDocument, error)
}
