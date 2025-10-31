package rag

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
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

func (i *instanceList) DisplayName() string {
	return i.name
}

func (i *instanceList) CollectionName() string {
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

func (i *instanceList) Owner() string {
	for _, r := range i.rags {
		o := r.Owner()
		if len(o) > 0 && o != i.CollectionName() {
			return o
		}
	}
	return i.CollectionName()
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

func (i *instanceList) Embbed(ctx context.Context, _ *slog.Logger, filters ...vecdb.Filter) error {
	for _, r := range i.rags {
		if err := r.Embbed(ctx, i.slog, filters...); err != nil {
			i.slog.Warn("Cannot embed rag list member", "err", err, "rag.name", r.DisplayName())
		}
	}
	return nil
}

func (i *instanceList) GetDocuments(ctx context.Context, slog *slog.Logger) (chan *vecdb.EmbeddDocument, error) {
	c := make(chan *vecdb.EmbeddDocument, 10)
	go func() {
		defer close(c)
		var wg sync.WaitGroup
		for _, r := range i.rags {
			wg.Go(func() {
				cc, err := r.GetDocuments(ctx, slog)
				if err != nil {
					slog.Warn("Cannot get documents of rag list member", "err", err, "rag.name", r.DisplayName())
					return
				}
				for doc := range cc {
					c <- doc
				}
			})
		}
		wg.Wait()
	}()
	return c, nil
}

func (i *instanceList) ListCollections(ctx context.Context, _ *slog.Logger) ([]*chroma.Collection, error) {
	cols := make([]*chroma.Collection, 0)
	for _, r := range i.rags {
		c, err := r.ListCollections(ctx, i.slog)
		if err != nil {
			i.slog.Warn("Cannot list collections of rag list member", "err", err, "rag.name", r.DisplayName())
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
			i.slog.Error("Cannot search vecDD of rag list member", "err", err, "rag.name", r.DisplayName(), "collection", collection)
		}
		docs = append(docs, c...)
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
