package scraper

import (
	"context"
	"fmt"

	"github.com/vogtp/langchaingo/documentloaders"
	"github.com/vogtp/langchaingo/schema"
	"github.com/vogtp/langchaingo/textsplitter"
)

var _ documentloaders.Loader = Scraper{}

// Load loads from a source and returns documents.
func (s Scraper) Load(ctx context.Context) ([]schema.Document, error) {
	docsChannel, err := s.Call(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot scrap: %s", err)
	}
	docs := make([]schema.Document, 0)
	for d := range docsChannel {
		docs = append(docs,
			schema.Document{
				PageContent: d.Document,
				Metadata:    d.MetaData,
			},
		)
	}
	return docs, nil
}

// LoadAndSplit loads from a source and splits the documents using a text splitter.
func (s Scraper) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := s.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot load scrapped documents: %w", err)
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
