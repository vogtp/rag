package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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

var _ (Instance) = (*instanceCfg)(nil)

type instanceCfg struct {
	slog   *slog.Logger
	config cfg.RagConfig

	vecDB  *vecdb.VecDB
	models []Model

	bearerAuth bearer.Auth
}

func newInstanceCfg(ctx context.Context, slog *slog.Logger, config cfg.RagConfig) (*instanceCfg, error) {
	m := instanceCfg{
		slog:       slog.With("rag.name", config.Name, "collection.name", config.Vecdb.CollectionName),
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

func (i instanceCfg) Name() string {
	return i.config.Name
}

func (i instanceCfg) UpdateIntervall() time.Duration {
	return i.config.VecDBUpdateIntervall()
}

func (i instanceCfg) Authorise(w http.ResponseWriter, r *http.Request) bool {
	return i.bearerAuth.Authorise(w, r)
}

// LLM returns the name of the LLM that is used for generation
func (i instanceCfg) LLM() string {
	return i.config.Model.LLM
}

func (i *instanceCfg) updateModelsFromChroma(ctx context.Context) error {

	collections, err := i.vecDB.ListCollections(ctx, &i.config)
	if err != nil {
		return fmt.Errorf("cannot list chroma collections: %w", err)
	}

	model := i.config.Model.LLM
	for _, c := range collections {
		i.models = append(i.models, VectorStoreModel{Name: c.Name, vecDB: i.vecDB, Collection: c.Name, LLMName: model, config: i.config, bearerAuth: i.bearerAuth})
	}
	i.slog.Debug("Models raw ", "models", i.models)
	slices.SortFunc(i.models, func(a, b Model) int { return strings.Compare(a.GetName(), b.GetName()) })
	i.slog.Debug("Models sort", "models", i.models)
	i.models = slices.CompactFunc(i.models, func(a, b Model) bool { return strings.EqualFold(a.GetName(), b.GetName()) })
	i.slog.Debug("Models comp", "models", i.models)
	return nil
}

func (i *instanceCfg) Models(ctx context.Context) []Model {
	if err := i.updateModelsFromChroma(ctx); err != nil {
		i.slog.WarnContext(ctx, "Cannot update models from chroma", "err", err)
	}
	return i.models
}

func (i instanceCfg) Model(name string) (Model, error) {
	i.slog.Debug("Query model", "model", name)
	decoded, err := url.QueryUnescape(name)
	if err != nil {
		decoded = name
	}
	for _, model := range i.models {
		i.slog.Debug("looking for model", "model", decoded, "cur", model.GetName())
		if strings.EqualFold(model.GetName(), decoded) {
			i.slog.Debug("found model", "model", decoded)
			return model, nil
		}
	}
	return nil, fmt.Errorf("model %s not found", name)
}
func (i instanceCfg) Embbed(ctx context.Context) error {
	if len(i.config.Vecdb.CollectionName) < 1 {
		return nil
	}
	return confluence.Embed(ctx, i.slog, i.config)
}

func (i instanceCfg) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	res, err := i.vecDB.Query(ctx, collection, []string{query}, int32(maxResults))
	if err != nil {
		return nil, fmt.Errorf("query vector DB: %w", err)
	}
	return res[0].Documents, nil
}
func (i instanceCfg) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	return i.vecDB.ListCollections(ctx, &i.config)
}
