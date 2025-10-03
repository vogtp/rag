package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
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

var _ (Instance) = (*instanceList)(nil)

type instanceList struct {
	slog *slog.Logger
	name string
	rags []Instance
}

func newInstanceList(slog *slog.Logger, name string, rags ...Instance) Instance {
	return &instanceList{
		slog: slog,
		name: name,
		rags: rags,
	}
}

func (i *instanceList) Name() string {
	return i.name
}

func (i *instanceList) Model(name string) (m Model, err error) {
	for _, r := range i.rags {
		m, err = r.Model(name)
		if m != nil {
			return m, err
		}
	}
	return nil, err
}

func (i *instanceList) Models(ctx context.Context) []Model {
	m := make([]Model, 0)
	for _, r := range i.rags {
		m = append(m, r.Models(ctx)...)
	}
	return m
}

func (i *instanceList) LLM() string {
	for _, r := range i.rags {
		llm := r.LLM()
		if len(llm) > 0 {
			return llm
		}
	}
	return viper.GetString(cfg.ModelLLM)
}

func (i *instanceList) UpdateIntervall() time.Duration {
	d := cfg.DefaultVecDBUpdateIntervall
	for _, r := range i.rags {
		intervall := r.UpdateIntervall()
		if intervall < cfg.MinVecDBUpdateIntervall {
			continue
		}
		d = min(d, intervall)
	}
	return d
}

func (i *instanceList) Embbed(ctx context.Context) error {
	for _, r := range i.rags {
		if err := r.Embbed(ctx); err != nil {
			i.slog.Error("Cannot embed rag list member", "err", err, "rag.name", r.Name())
		}
	}
	return nil
}

func (i *instanceList) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	cols := make([]*chroma.Collection, 0)
	for _, r := range i.rags {
		c, err := r.ListCollections(ctx)
		if err != nil {
			i.slog.Error("Cannot list collections of rag list member", "err", err, "rag.name", r.Name())
		}
		cols = append(cols, c...)
	}
	if len(cols) < 1 {
		return cols, fmt.Errorf("no collections found")
	}
	return cols, nil
}

func (i *instanceList) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	docs := make([]vecdb.QueryDocument, 0)
	for _, r := range i.rags {
		c, err := r.SearchVecDB(ctx, slog, collection, query, maxResults)
		if err != nil {
			i.slog.Error("Cannot search vecDV of rag list member", "err", err, "rag.name", r.Name())
		}
		docs = append(docs, c...)
	}
	if len(docs) < 1 {
		return docs, fmt.Errorf("no documents found")
	}
	return docs, nil
}

func (i *instanceList) Authorise(w http.ResponseWriter, r *http.Request) bool {
	for _, rag := range i.rags {
		if rag.Authorise(w, r) {
			return true
		}
	}
	return false
}
