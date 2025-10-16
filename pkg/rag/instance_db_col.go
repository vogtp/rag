package rag

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
)

var _ (types.Instance) = (*instanceDBCol)(nil)
var _ (cfg.RagConfig) = (*instanceDBCol)(nil)

type instanceDBCol struct {
	*instanceCfg
	collection *usercfg.Collection

	muEmbed sync.Mutex
}

func newInstanceDBCol(ctx context.Context, slog *slog.Logger, col *usercfg.Collection) (*instanceDBCol, error) {
	src := col.Source
	spaces := strings.Split(src.Parts, ",")
	if len(spaces) < 1 {
		spaces = strings.Split(src.Parts, " ")
	}
	if len(spaces) < 1 {
		return nil, fmt.Errorf("no spaces found: %q", src.Parts)
	}
	ucfg := cfg.RagConfigInteral{
		NameInt: col.DisplayName,
		ModelInt: cfg.ModelCfg{
			Embedding: viper.GetString(cfg.ModelEmbedding),
			LLM:       viper.GetString(cfg.ModelLLM),
		},
		VecdbInt: cfg.VecDbCfg{
			CollectionName:  col.CollectionName,
			UpdateIntervall: "24h",
		},
		ConfluenceCfg: cfg.ConfluenceCfg{
			Key:     src.Key,
			BaseURL: src.URL,
			Spaces:  spaces,
		},
	}
	ic, err := newInstanceCfg(ctx, slog, ucfg)
	if err != nil {
		return nil, err
	}
	i := &instanceDBCol{
		instanceCfg: ic,
		collection:  col,
	}
	return i, nil
}

func (i *instanceDBCol) CollectionName() string {
	return i.collection.CollectionName
}

func (i *instanceDBCol) Embbed(ctx context.Context) error {
	if !i.muEmbed.TryLock() {
		return fmt.Errorf("Embedding (%q) already running!", i.config.Name())
	}
	defer i.muEmbed.Unlock()
	return confluence.Embed(ctx, i.slog, i)
}

func (i *instanceDBCol) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	if len(i.CollectionName()) < 1 {
		return nil, fmt.Errorf("no collection name given")

	}
	cols, err := i.vecDB.ListAllCollections(ctx)
	if err != nil {
		return nil, err
	}
	collections := make([]*chroma.Collection, 0, len(cols))
	for _, col := range cols {
		if !strings.EqualFold(col.Name, i.CollectionName()) {
			continue
		}
		collections = append(collections, col)
	}
	return collections, nil
}
