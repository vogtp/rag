package confluence

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func Embed(ctx context.Context, slog *slog.Logger, config cfg.RagConfig) error {
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
		slog := slog.With("space", space)
		slog.Info("Starting confluence embdding")
		start := time.Now()
		cnt, err := client.Embedd(ctx, slog, config, c)
		if err != nil {
			slog.Warn("Embedding returned an error", "err", err)
		}
		d := time.Since(start)
		slog.Info("Embebbing finished", "document.count", cnt, "duration", d.String(), "duration_ms", d.Milliseconds())

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}
