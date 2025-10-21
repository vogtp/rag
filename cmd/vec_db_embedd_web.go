package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/scraper"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/history"
)

var vecDbEmbbedScrapCmd = &cobra.Command{
	Use:     "scrap collection_name https://example.net",
	Short:   "scrap webpage",
	Aliases: []string{"s", "scraper"},
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionName := "anleitungen"
		url := "https://its.unibas.ch/de/anleitungen/"
		if len(args) > 0 {
			collectionName = args[0]
		}
		if len(args) > 1 {
			url = args[1]
		}
		fmt.Printf("Scrapping %s into %s\n", url, collectionName)
		return scapper2vecDB(cmd.Context(), url, collectionName)
	},
}

func scapper2vecDB(ctx context.Context, url string, collectionName string) error {
	scrap, err := scraper.New(url,
		scraper.WithBlacklist([]string{
			"ueber-uns", "about-us",
			"aktuelles",
			"servicekatalog",
			"news",
			"shared-elements",
			"event",
		}),
	)
	if err != nil {
		return fmt.Errorf("cannot create scrapper: %w", err)
	}

	docsChannel, err := scrap.Call(ctx)
	if err != nil {
		return err
	}

	slog := logger.New()
	client, err := vecdb.New(ctx, slog, viper.GetString(cfg.ModelEmbedding), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx, slog)))
	if err != nil {
		return fmt.Errorf("Failed to create vector DB: %w", err)
	}
	_, err = client.Embedd(ctx, slog, collectionName, docsChannel, history.New(slog, collectionName, viper.GetDuration(cfg.VecDBUpdateIntervall)))
	return err
}
