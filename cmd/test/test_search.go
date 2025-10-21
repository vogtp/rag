package test

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/cobra"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

var testFailedError = fmt.Errorf("Tests failed")

var searchTstStartCmd = &cobra.Command{
	Use:          "search",
	Short:        "Only do the search part of the tests",
	Aliases:      []string{"s"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer func(start time.Time) { fmt.Printf("Duration %s\n", time.Since(start)) }(time.Now())
		ctx := cmd.Context()
		slog := getSlog()
		tt, err := loadTestData()
		if err != nil {
			return err
		}
		return searchTestData(ctx, slog, tt)

	},
}

func searchTestData(ctx context.Context, slog *slog.Logger, tt *testData) error {
	maxResult := 7
	var retErr error
	for _, col := range tt.Collections() {
		fmt.Printf("* %s\n", col.DisplayName())
		for _, t := range tt.Tests {
			fmt.Printf("  %q:\n", t.Question)
			docs, err := col.SearchVecDB(ctx, slog, col.CollectionName(), t.Question, maxResult)
			if err != nil {
				// if collection does not exist recreate
				if strings.Contains(err.Error(), fmt.Sprintf("InvalidCollection: Collection %s does not exist.", col.CollectionName())) {
					if err := createVecDBCollecions(ctx, slog, tt); err != nil {
						return err
					}
					return searchTestData(ctx, slog, tt)
				}
				return err
			}
			if len(docs) < 1 {
				fmt.Printf("   %q no docs found\n", col.Displayname)
				retErr = testFailedError
				continue
			}
			for _, r := range t.Results {
				if err := ensureTitle(r, docs); err != nil {
					fmt.Printf("    %v\n", err)
					retErr = testFailedError
				}
				if err := ensureURL(r, docs); err != nil {
					fmt.Printf("    %v\n", err)
					retErr = testFailedError
				}
				if err := ensureKeywords(r, docs); err != nil {
					fmt.Printf("    %v\n", err)
					retErr = testFailedError
				}
			}
		}
	}
	return retErr
}

func ensureTitle(r testDataResult, docs []vecdb.QueryDocument) error {
	for _, d := range docs {
		if d.Title == r.Title {
			return nil
		}
	}
	return fmt.Errorf("Title %q not found", r.Title)
}

func ensureURL(r testDataResult, docs []vecdb.QueryDocument) error {
	for _, d := range docs {
		if d.URL == r.URL {
			return nil
		}
	}
	return fmt.Errorf("URL %q not found", r.URL)
}

func ensureKeywords(r testDataResult, docs []vecdb.QueryDocument) error {
	for _, k := range r.Keywords {
		if err := ensureKeyword(k, docs); err != nil {
			return err
		}
	}
	return nil
}

func ensureKeyword(keyword string, docs []vecdb.QueryDocument) error {
	for _, d := range docs {
		if strings.Contains(d.Content, keyword) {
			return nil
		}
	}
	return fmt.Errorf("Keyword %q not found", keyword)
}
