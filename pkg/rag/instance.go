package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strings"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
	"github.com/vogtp/rag/pkg/web/bearer"
)

type Manager interface {
	Name() string

	Model(name string) (Model, error)
	Models(ctx context.Context) []Model

	LLM() string

	UpdateIntervall() time.Duration
	Embbed(ctx context.Context) error
	ListCollections(ctx context.Context) ([]*chroma.Collection, error)
	SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error)

	BearerAuth() bearer.Auth
}

var _ (Manager) = (*instance)(nil)

type instance struct {
	slog   *slog.Logger
	config cfg.RagConfig

	vecDB  *vecdb.VecDB
	models []Model

	bearerAuth bearer.Auth
}

func newInstance(ctx context.Context, slog *slog.Logger, config cfg.RagConfig) (Manager, error) {
	m := instance{
		slog:       slog,
		config:     config,
		bearerAuth: bearer.TokenAuth(config.APIToken),
		models: []Model{
			OllamaModel{
				Name:    config.Model.LLM,
				LLMName: config.Model.LLM,
			},
		},
	}
	v, err := vecdb.New(ctx, slog, config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to chroma: %w", err)
	}
	m.vecDB = v
	if err := m.updateModelsFromChroma(ctx); err != nil {
		return nil, fmt.Errorf("cannot get collections from chroma: %w", err)
	}

	return &m, nil
}

func (m instance) Name() string {
	return m.config.Name
}

func (m instance) UpdateIntervall() time.Duration {
	return m.config.VecDBUpdateIntervall()
}

func (m instance) BearerAuth() bearer.Auth {
	return m.bearerAuth
}

// LLM returns the name of the LLM that is used for generation
func (m instance) LLM() string {
	return m.config.Model.LLM
}

func (m *instance) updateModelsFromChroma(ctx context.Context) error {

	collections, err := m.vecDB.ListCollections(ctx, &m.config)
	if err != nil {
		return fmt.Errorf("cannot list chroma collections: %w", err)
	}

	model := m.config.Model.LLM
	for _, c := range collections {
		m.models = append(m.models, VectorStoreModel{Name: c.Name, vecDB: m.vecDB, Collection: c.Name, LLMName: model, config: m.config, bearerAuth: m.bearerAuth})
	}
	m.slog.Info("Models raw ", "models", m.models)
	slices.SortFunc(m.models, func(a, b Model) int { return strings.Compare(a.GetName(), b.GetName()) })
	m.slog.Info("Models sort", "models", m.models)
	m.models = slices.CompactFunc(m.models, func(a, b Model) bool { return strings.EqualFold(a.GetName(), b.GetName()) })
	m.slog.Info("Models comp", "models", m.models)
	return nil
}

func (m *instance) Models(ctx context.Context) []Model {
	if err := m.updateModelsFromChroma(ctx); err != nil {
		m.slog.WarnContext(ctx, "Cannot update models from chroma", "err", err)
	}
	return m.models
}

func (m instance) Model(name string) (Model, error) {
	m.slog.Debug("Query model", "model", name)
	decoded, err := url.QueryUnescape(name)
	if err != nil {
		decoded = name
	}
	for _, model := range m.models {
		m.slog.Debug("looking for model", "model", decoded, "cur", model.GetName())
		if strings.EqualFold(model.GetName(), decoded) {
			m.slog.Debug("found model", "model", decoded)
			return model, nil
		}
	}
	return nil, fmt.Errorf("model %s not found", name)
}
func (m instance) Embbed(ctx context.Context) error {
	return confluence.Embed(ctx, m.slog, m.config)
}

func (m instance) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	res, err := m.vecDB.Query(ctx, collection, []string{query}, int32(maxResults))
	if err != nil {
		return nil, fmt.Errorf("query vector DB: %w", err)
	}
	return res[0].Documents, nil
}
func (m instance) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	return m.vecDB.ListCollections(ctx, &m.config)
}
