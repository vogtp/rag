package cmd

import (
	"fmt"
	"log/slog"
	"strings"
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
		dcfg := cfg.DefaultRagCfg()
		dcfg.Vecdb.CollectionName = collectionName
		client, err := vecdb.New(ctx, slog.Default(), dcfg, vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
		if err != nil {
			return fmt.Errorf("create vector DB: %w", err)
		}

		_, err = client.Embedd(ctx, dcfg, filesystem.Generate(ctx, path))
		return err
	},
}

var vecDbEmbbedConfluenceCmd = &cobra.Command{
	Use:     "confluence <rag_name>",
	Short:   "Embbed confluence spaces into a collection",
	Aliases: []string{"conf", "c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		ctx := cmd.Context()
		ragCfgs, err := cfg.GetRagConfig()
		if err != nil {
			return err
		}
		notUsed := ""
		for _, ragCfg := range ragCfgs {
			if len(args) > 0 && !strings.EqualFold(args[0], ragCfg.Name) {
				notUsed = fmt.Sprintf("%s %s", notUsed, ragCfg.Name)
				continue
			}
			if err := confluence.Embed(ctx, slog, ragCfg); err != nil {
				fmt.Printf("Embed confluence: %v", err)
			}
		}
		if len(notUsed) > 0 {
			return fmt.Errorf("RAG %q not found, possible RAGs:%s", args[0], notUsed)
		}
		return nil
	},
}
