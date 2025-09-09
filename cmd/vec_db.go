package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/chroma"
)

func addVecDB() {
	rootCmd.AddCommand(vecDbCmd)
	vecDbCmd.AddCommand(vecDbStartChromaCmd)
	vecDbCmd.AddCommand(vecDbStopChromaCmd)
	vecDbCmd.AddCommand(vecDbRmCmd)
	vecDbCmd.AddCommand(vecDbLsCmd)
	vecDbCmd.AddCommand(vecDbSearchCmd)
	vecDbCmd.AddCommand(vecDbColLsCmd)

	vecDbCmd.AddCommand(vecDbEmbbedCmd)
	vecDbEmbbedCmd.AddCommand(vecDbEmbbedPathCmd)
	vecDbEmbbedCmd.AddCommand(vecDbEmbbedConfluenceCmd)
	vecDbEmbbedCmd.AddCommand(vecDbEmbbedScrapCmd)
}

var vecDbCmd = &cobra.Command{
	Use:   "vecdb",
	Short: "manage the vector DB",

	Aliases: []string{"v", "vec", "vecDB", "vecDb"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var vecDbStopChromaCmd = &cobra.Command{
	Use:   "stop chroma",
	Short: "stop a chroma container",

	Aliases: []string{},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := chroma.NewContainer(slog.Default())
		if err != nil {
			return err
		}
		return c.EnsureStopped(cmd.Context())
	},
}
var vecDbStartChromaCmd = &cobra.Command{
	Use:   "start chroma",
	Short: "start a chroma container",

	Aliases: []string{"run", "r"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return chroma.EnsureStarted(slog.Default(), cmd.Context(), viper.GetInt(cfg.ChromaPort))
	},
}

var vecDbSearchCmd = &cobra.Command{
	Use:   "search <collection> <query>",
	Short: "Search in a collection",

	Aliases: []string{"s"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}
		collectionName := args[0]
		search := args[1]
		start := time.Now()
		defer func(t time.Time) {
			fmt.Printf("Searching collection %s took %s\n", collectionName, time.Since(t))
		}(start)
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		res, err := client.Query(ctx, collectionName, []string{search}, 5)
		if err != nil {
			return fmt.Errorf("Failed to query vector DB: %w", err)
		}
		for i, r := range res[0].Documents {
			fmt.Printf("\n\nDocument %v: %+v\n", i, r)
		}
		fmt.Printf("Found %v documents for %q\n", len(res), search)
		return nil
	},
}

var vecDbRmCmd = &cobra.Command{
	Use:   "rm <collection> [all]",
	Short: "Delete collections.  Sepatate by space or use all to delete all",

	Aliases: []string{"del", "remove", "delete"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		if args[0] == "all" {
			cols, err := client.ListCollections(ctx)
			if err != nil {
				return err
			}
			for _, c := range cols {
				slog.Info("Deleting collection", "name", c.Name)
				if err := client.DeleteCollection(ctx, c.Name); err != nil {
					return err
				}
			}
			return nil
		}
		for _, a := range args {
			if err := client.DeleteCollection(ctx, a); err != nil {
				return err
			}
		}
		return nil
	},
}

var vecDbLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all collection",

	Aliases: []string{"list", "show"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default())
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		cols, err := client.ListCollections(ctx)
		if err != nil {
			return fmt.Errorf("cannot list collections: %w", err)
		}
		for _, c := range cols {
			fmt.Printf("%s\n", c.Name)
		}
		return nil
	},
}

var vecDbColLsCmd = &cobra.Command{
	Use:   "col <collection_name>",
	Short: "List collection documents",

	Aliases: []string{"c", "collection"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		colName := args[0]
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		col, err := client.GetCollection(ctx, colName)
		if err != nil {
			return err
		}

		res, err := col.GetWithOptions(ctx)
		if err != nil {
			return fmt.Errorf("cannot get collection documents: %w", err)
		}
		for i, d := range res.Metadatas {
			// fmt.Printf("%s %s\n", d[vecdb.MetaPath], d[vecdb.MetaUpdated])
			fmt.Printf("  ID: %v Len: %v Meta: %v\n", res.Ids[i], len(res.Documents[i]), d)
		}
		fmt.Printf("Found %v docs in collection %s\n", len(res.Documents), colName)
		return nil
	},
}
