package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
	"github.com/vogtp/rag/pkg/vecDB/filesystem"
)

var vecDbEmbbedCmd = &cobra.Command{
	Use:     "embedd",
	Short:   "Embbed to content to a collection",
	Aliases: []string{"e", "emb", "embbed"},
}

var vecDbEmbbedPathCmd = &cobra.Command{
	Use:   "path <collection> <path>",
	Short: "Embbed the content of a path to a collection",

	Aliases: []string{"path", "dir"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}
		collectionName := args[0]
		path := args[1]
		start := time.Now()
		defer func(t time.Time) {
			slog.Info(fmt.Sprintf("Updating collection %s took %s", collectionName, time.Since(t)))
		}(start)
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		_, err = client.Embedd(ctx, collectionName, filesystem.Generate(ctx, path))
		return err
	},
}

var vecDbEmbbedConfluenceCmd = &cobra.Command{
	Use:     "confluence <collection_name>",
	Short:   "Embbed confluence spaces into a collection",
	Aliases: []string{"conf", "c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		ctx := cmd.Context()
		if len(args) < 1 {
			return cmd.Usage()
		}
		collectionName := args[0]

		return confluence.Embbed(ctx, slog, collectionName)
	},
}
