package test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
)

var (
	testFailedDocuemntNotFoundError = fmt.Errorf("Tests failed: Document not found")
	testFailedKeywordsNotFoundError = fmt.Errorf("Tests failed: Keyword not found")
)

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
	maxResult := cfg.CntSeachResults
	var retErr error
	var resultSummary strings.Builder
	for _, col := range tt.Collections() {
		if excludeCollection(&col) {
			continue
		}
		fmt.Printf("* %s\n", col.DisplayName())
		for _, t := range tt.Tests {
			if !slices.Contains(t.Collections, col.CollectionName()) {
				continue
			}
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
				tstErr = testFailedDocuemntNotFoundError
			} else {
				for _, r := range t.Results {
					if err := ensureTitle(r, docs); err != nil {
						fmt.Fprintf(&resOut, "%s%v\n", intent, err)
						tstErr = testFailedDocuemntNotFoundError
					}
					if err := ensureURL(r, docs); err != nil {
						fmt.Fprintf(&resOut, "%s%v\n", intent, err)
						tstErr = testFailedDocuemntNotFoundError
					}
					for _, k := range r.Keywords {
						if err := ensureKeyword(k, docs); err != nil {
							fmt.Fprintf(&resOut, "%s%v\n", intent, err)
							tstErr = testFailedKeywordsNotFoundError
						}
					}
					if tstErr == testFailedKeywordsNotFoundError {
						for _, d := range docs {
							if strings.HasSuffix(d.URL, r.URL) {
								filename := fmt.Sprintf("ignore_doc_content_%v.txt", d.Title)
								f, err := os.Create(filename)
								if err != nil {
									fmt.Printf("Cannot create file %q: %v", filename, err)
									continue
								}
								fmt.Fprintf(f, "Document of %q:\nURL: %q\nKeywords: %+v\n%s\n", r.Title, r.URL, r.Keywords, d.Document)
								f.Close()
								fmt.Fprintf(&resOut, "Saved %s with content of %s\n", filename, d.Title)
							}
						}
					}
				}
				if tstErr != nil {
					fmt.Fprintf(&resOut, "%sFound docs:\n", intent)
					intent := fmt.Sprintf("  %s", intent)
					for _, d := range docs {
						url := d.URL
						maxLen := 100
						if len(url) > maxLen {
							spacer := "[...]"
							maxLen = maxLen - len(spacer)
							url = fmt.Sprintf("%s%s%s", url[:maxLen/2], spacer, url[len(url)-maxLen/2:])
							// fmt.Fprintf(&resOut, "%s %v\n%s %v\n", d.URL, len(d.URL), url, len(url))
						}
						title := d.Title
						if len(title) > 40 {
							title = title[:40]
						}
						// title = d.Title
						// url = d.URL
						fmt.Fprintf(&resOut, "%s %-40s %s\n", intent, title, url)
					}
				}
			}
			tick := "✅"
			if tstErr != nil {
				tick = "⚠"
				retErr = tstErr
			}
			fmt.Printf("  %q in %s: %s\n%s\n", t.Question, col.CollectionName(), tick, resOut.String())
			fmt.Fprintf(&resultSummary, "%q in %s: %s\n", t.Question, col.CollectionName(), tick)
		}
	}
	fmt.Println(resultSummary.String())
	return retErr
}

func ensureTitle(r testDataResult, docs []types.QueryDocument) error {
	for i, d := range docs {
		if strings.HasPrefix(d.Title, r.Title) {
			fmt.Printf("          Found title in doc %v dist: %v\n", i, d.Distance)
			return nil
		}
	}
	return fmt.Errorf("Title not found: %q", r.Title)
}

func ensureURL(r testDataResult, docs []types.QueryDocument) error {
	for i, d := range docs {
		if strings.HasSuffix(d.URL, r.URL) {
			fmt.Printf("          Found URL in doc %v dist: %v\n", i, d.Distance)
			return nil
		}
	}
	return fmt.Errorf("URL not found: %q", r.URL)
}

func ensureKeyword(keyword string, docs []types.QueryDocument) error {
	for _, d := range docs {
		if strings.Contains(d.Document, keyword) {
			return nil
		}
	}
	return fmt.Errorf("Keyword not found: %q", keyword)
}
