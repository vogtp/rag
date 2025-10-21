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

func searchTestData(ctx context.Context, log *slog.Logger, tt *testData) error {
	if h, ok := log.Handler().(testSlogHandler); ok {
		h.doDot = false
		log = slog.New(h)
	}
	maxResult := 7
	var retErr error
	for _, col := range tt.Collections() {
		fmt.Printf("* %s\n", col.DisplayName())
		for _, t := range tt.Tests {
			var resOut strings.Builder
			var tstErr error
			docs, err := col.SearchVecDB(ctx, log, col.CollectionName(), t.Question, maxResult)
			if err != nil {
				// if collection does not exist recreate
				if strings.Contains(err.Error(), fmt.Sprintf("InvalidCollection: Collection %s does not exist.", col.CollectionName())) {
					if err := createVecDBCollecions(ctx, log, tt); err != nil {
						return err
					}
					return searchTestData(ctx, log, tt)
				}
				return err
			}
			intent := "      "
			if len(docs) < 1 {
				fmt.Fprintf(&resOut, "%s%q no docs found\n", intent, col.Displayname)
				tstErr = testFailedError
			} else {
				for _, r := range t.Results {
					if err := ensureTitle(r, docs); err != nil {
						fmt.Fprintf(&resOut, "%s%v\n", intent, err)
						tstErr = testFailedError
					}
					if err := ensureURL(r, docs); err != nil {
						fmt.Fprintf(&resOut, "%s%v\n", intent, err)
						tstErr = testFailedError
					}
					for _, k := range r.Keywords {
						if err := ensureKeyword(k, docs); err != nil {
							fmt.Fprintf(&resOut, "%s%v\n", intent, err)
							tstErr = testFailedError
						}
					}
				}
			}
			tick := "✅"
			if tstErr != nil {
				tick = "⚠"
				retErr = tstErr
			}
			fmt.Printf("  %q: %s\n%s\n", t.Question, tick, resOut.String())
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
	return fmt.Errorf("Title not found: %q", r.Title)
}

func ensureURL(r testDataResult, docs []vecdb.QueryDocument) error {
	for _, d := range docs {
		if strings.HasSuffix(d.URL, r.URL) {
			return nil
		}
	}
	return fmt.Errorf("URL not found: %q", r.URL)
}

func ensureKeyword(keyword string, docs []vecdb.QueryDocument) error {
	for _, d := range docs {
		if strings.Contains(d.Content, keyword) {
			return nil
		}
	}
	return fmt.Errorf("Keyword not found: %q", keyword)
}
