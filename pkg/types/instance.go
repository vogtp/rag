package types

import (
	"context"
	"log/slog"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/vogtp/rag/pkg/web/bearer"
)

type Instance interface {
	DisplayName() string
	CollectionName() string
	Owner() string

	Model(name string) (Model, error)

	Models(ctx context.Context) []Model
	ListCollections(ctx context.Context, slog *slog.Logger) ([]*chroma.Collection, error)

	LLM() string

	UpdateIntervall() time.Duration
	Embbed(ctx context.Context, slog *slog.Logger, filters ...Filter) error
	SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]QueryDocument, error)

	DocRetriver

	bearer.Auth
}

type Filter interface {
	// ShouldEmbedd returns true if the document should be embedded
	ShouldEmbedd(*EmbeddDocument) bool
	// ReqisterEmedded tels the Filter that a document has been embedded
	ReqisterEmedded(*EmbeddDocument)
}

// QueryDocument is a document found in the vectorDB
type QueryDocument struct {
	EmbedContent string // Content is the part of the document used for the embedding
	Document     string // Document is the original
	Modified     string
	URL          string
	Title        string
	IDField      string
	Distance     float32
}

func (qd QueryDocument) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", qd.Title),
		slog.String("URL", qd.URL),
	)
}
