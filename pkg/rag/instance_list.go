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
	"github.com/vogtp/rag/pkg/types"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

var _ (types.Instance) = (*instanceList)(nil)

type instanceList struct {
	slog *slog.Logger
	name string
	rags []types.Instance
}

func newInstanceList(slog *slog.Logger, name string, rags ...types.Instance) *instanceList {
	if rags == nil {
		rags = make([]types.Instance, 0)
	}
	return &instanceList{
		slog: slog,
		name: name,
		rags: rags,
	}
}

func (i *instanceList) Add(rags ...types.Instance) {
	i.rags = append(i.rags, rags...)
}

func (i *instanceList) GetName() string {
	return i.name
}

func (i *instanceList) Model(name string) (m types.Model, err error) {
	for _, r := range i.rags {
		m, err = r.Model(name)
		if m != nil {
			return m, err
		}
	}
	return nil, err
}

func (i *instanceList) Models(ctx context.Context) []types.Model {
	m := make([]types.Model, 0)
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
			i.slog.Warn("Cannot embed rag list member", "err", err, "rag.name", r.GetName())
		}
	}
	return nil
}

func (i *instanceList) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	cols := make([]*chroma.Collection, 0)
	for _, r := range i.rags {
		c, err := r.ListCollections(ctx)
		if err != nil {
			i.slog.Warn("Cannot list collections of rag list member", "err", err, "rag.name", r.GetName())
		}
		cols = append(cols, c...)
	}
	// if len(cols) < 1 {
	// 	return cols, fmt.Errorf("no collections found")
	// }
	return cols, nil
}

func (i *instanceList) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	docs := make([]vecdb.QueryDocument, 0)
	for _, r := range i.rags {
		c, err := r.SearchVecDB(ctx, slog, collection, query, maxResults)
		if err != nil {
			i.slog.Error("Cannot search vecDD of rag list member", "err", err, "rag.name", r.GetName(), "collection", collection)
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
