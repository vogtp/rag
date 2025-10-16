package confluence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func Embed(ctx context.Context, slog *slog.Logger, config cfg.RagConfig) error {
	collectionName := config.Collectionname()
	client, err := vecdb.New(ctx, slog, config.ModelEmbedding(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return fmt.Errorf("Failed to create vector DB: %w", err)
	}

	for _, space := range config.Confluence().Spaces {
		space = strings.TrimSpace(space)
		if len(space) < 1 {
			continue
		}
		c, err := GetDocuments(ctx, slog, config, space)
		if err != nil {
			return err
		}
		// o1, o2 := fanOut(c)

		// go func() {
		// 	colName := fmt.Sprintf("%s-%s", collectionName, "all")
		// 	cfg := config
		// 	cfg.Vecdb.CollectionName = colName
		// 	slog := slog.With("collectionName", colName, "space", "all")
		// 	slog.Info("Starting confluence embdding")
		// 	if err := embedd(ctx, client, cfg, o1); err != nil {
		// 		slog.Warn("Embedding returned an error", "err", err)
		// 	}
		// }()
		// go func() {
		colName := fmt.Sprintf("%s-%s", collectionName, strings.ToLower(space))
		cfg := config
		// cfg.Vecdb.CollectionName = colName
		slog := slog.With("collectionName", colName, "space", space)
		slog.Info("Starting confluence embdding")
		if err := embedd(ctx, client, cfg, c); err != nil {
			slog.Warn("Embedding returned an error", "err", err)
		}
		// }()
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

func embedd(ctx context.Context, client *vecdb.VecDB, config cfg.RagConfig, c chan *vecdb.EmbeddDocument) error {
	slog := slog.With("collection", config.Collectionname())
	slog.Info("Embebbing start")
	start := time.Now()
	cnt, err := client.Embedd(ctx, config, c)
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("confluence embedding failed: %w", err)
	}
	d := time.Since(start)
	slog.Info("Embebbing finished", "document.count", cnt, "duration", d.String(), "duration_ms", d.Milliseconds())
	return nil
}
