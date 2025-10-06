package rag

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
)

var _ (Instance) = (*instanceDBCol)(nil)

type instanceDBCol struct {
	*instanceCfg
}

func newInstanceDBCol(ctx context.Context, slog *slog.Logger, col *ent.Collection) (Instance, error) {
	srcs, err := col.Sources(ctx)
	if err != nil {
		return nil, fmt.Errorf("get source system from collection %q: %w", col.Name, err)
	}
	if len(srcs) < 1 || srcs[0] == nil {
		return nil, fmt.Errorf("no source system from collection %q found", col.Name)
	}
	src := srcs[0]
	spaces := strings.Split(src.Parts, ",")
	if len(spaces) < 1 {
		spaces = strings.Split(src.Parts, " ")
	}
	if len(spaces) < 1 {
		return nil, fmt.Errorf("no spaces found: %q", src.Parts)
	}
	ucfg := cfg.RagConfig{
		Name: col.Name,
		Model: cfg.ModelCfg{
			Embedding: viper.GetString(cfg.ModelEmbedding),
			LLM:       viper.GetString(cfg.ModelLLM),
		},
		Vecdb: cfg.VecDbCfg{
			CollectionName:  col.CollectionName,
			UpdateIntervall: "24h",
		},
		Confluence: cfg.ConfluenceCfg{
			Key:     src.Key,
			BaseURL: src.URL,
			Spaces:  spaces,
		},
	}
	return newInstanceCfg(ctx, slog, ucfg)
}
