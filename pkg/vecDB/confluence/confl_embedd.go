package confluence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

var wg sync.WaitGroup

func Embbed(ctx context.Context, slog *slog.Logger, collectionName string) error {
	client, err := vecdb.New(ctx, slog, vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return fmt.Errorf("Failed to create vector DB: %w", err)
	}

	for _, space := range viper.GetStringSlice(cfg.ConfluenceSpaces) {
		c, err := GetDocuments(ctx, slog, space)
		if err != nil {
			return err
		}
		o1, o2 := fanOut(c)

		go func() {
			if err := embedd(ctx, client, fmt.Sprintf("%s-%s", collectionName, "all"), o1); err != nil {
				slog.Warn("Embedding returned an error", "collectionName", collectionName, "space", "all", "err", err)
			}
		}()
		go func() {
			if err := embedd(ctx, client, fmt.Sprintf("%s-%s", collectionName, strings.ToLower(space)), o2); err != nil {
				slog.Warn("Embedding returned an error", "collectionName", collectionName, "space", space, "err", err)
			}
		}()
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	wg.Wait()
	return nil
}

func embedd(ctx context.Context, client *vecdb.VecDB, collectionName string, c chan vecdb.EmbeddDocument) error {
	wg.Add(1)
	defer wg.Done()
	cnt, err := client.Embedd(ctx, collectionName, c)
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("confluence embedding failed: %w", err)
	}
	slog.Info("Embebbing finished", "collection", collectionName, "document.count", cnt)
	return nil
}

func fanOut(in chan vecdb.EmbeddDocument) (chan vecdb.EmbeddDocument, chan vecdb.EmbeddDocument) {
	o1 := make(chan vecdb.EmbeddDocument)
	o2 := make(chan vecdb.EmbeddDocument)
	go func() {
		defer close(o1)
		defer close(o2)
		for d := range in {
			o1 <- d
			o2 <- d
		}
	}()
	return o1, o2
}
