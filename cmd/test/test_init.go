package test

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

var testDataFile = "cmd/test/ignore_testdata.yml"

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(tstCmd)

	tstCmd.AddCommand(createTstStartCmd)
	tstCmd.AddCommand(deleteTstStartCmd)
	tstCmd.AddCommand(searchTstStartCmd)
	tstCmd.AddCommand(allInOneTstStartCmd)
}

var tstCmd = &cobra.Command{
	Use:     "test",
	Short:   "Manage RAG web server",
	Aliases: []string{"t", "tst"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var allInOneTstStartCmd = &cobra.Command{
	Use:          "all",
	Short:        "Create, search and delete the test data collections",
	Aliases:      []string{"a", "full"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func(start time.Time) { fmt.Printf("Duration %s", time.Since(start)) }(time.Now())
		tt, err := loadTestData()
		if err != nil {
			return err
		}
		if err := createVecDBCollecions(cmd.Context(), getSlog(), tt); err != nil {
			return err
		}
		retErr := searchTestData(cmd.Context(), getSlog(), tt)
		if err := deleteVecDBCollecions(cmd.Context(), getSlog(), tt); err != nil {
			if retErr == nil {
				retErr = err
			}
		}
		return retErr
	},
}

var createTstStartCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create the test data collections",
	Aliases:      []string{"c"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func(start time.Time) { fmt.Printf("Duration %s", time.Since(start)) }(time.Now())
		tt, err := loadTestData()
		if err != nil {
			return err
		}
		return createVecDBCollecions(cmd.Context(), getSlog(), tt)
	},
}

var deleteTstStartCmd = &cobra.Command{
	Use:          "delete",
	Short:        "Delete the test data collections",
	Aliases:      []string{"rm", "del"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func(start time.Time) { fmt.Printf("Duration %s", time.Since(start)) }(time.Now())
		tt, err := loadTestData()
		if err != nil {
			return err
		}
		return deleteVecDBCollecions(cmd.Context(), getSlog(), tt)
	},
}

func createVecDBCollecions(ctx context.Context, slog *slog.Logger, tt *testData) error {
	fmt.Println("Creating collections")
	for _, col := range tt.Collections() {
		if err := col.Embbed(ctx, slog); err != nil {
			return err
		}
		vecDB, err := vecdb.New(ctx, slog, col.ModelEmbedding())
		if err != nil {
			return err
		}
		c, err := vecDB.GetCollection(ctx, col.CollectionName())
		if err != nil {
			return err
		}
		cnt, err := c.Count(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Embedded %s (%s) -> %v documents\n", col.CollectionName(), col.Source.Parts, cnt)
	}
	return nil
}

func deleteVecDBCollecions(ctx context.Context, slog *slog.Logger, tt *testData) error {
	fmt.Println("Deleting collections")
	for _, c := range tt.Setup.Collections {
		vecDB, err := vecdb.New(ctx, slog, c.Model.Embedding)
		if err != nil {
			return err
		}
		if err := vecDB.DeleteCollection(ctx, c.Name); err != nil {
			slog.Warn("Cannot delete collection", "collection", c.Name, "err", err)
			continue
		}
		fmt.Printf("Deleted %s\n", c.Name)
	}
	return nil
}
