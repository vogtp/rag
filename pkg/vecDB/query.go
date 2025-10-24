package vecdb

import (
	"context"
	"fmt"
	"log/slog"
)

// QueryResult is the result of a vectorDB search it contains one or more documents
type QueryResult struct {
	Question  string
	Documents []QueryDocument
}

// QueryDocument is a document found in the vectorDB
type QueryDocument struct {
	EmbedContent string // Content is the part of the document used for the embedding
	Document     string // Document is the original
	Modified     string
	URL          string
	Title        string
	IDField      string
	Distance     float32
}

func (qd QueryDocument) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", qd.Title),
		slog.String("URL", qd.URL),
	)
}

// Query searches the vectorDB
func (v *VecDB) Query(ctx context.Context, collection string, queryTexts []string, nResults int32) ([]QueryResult, error) {
	v.slog.Info("Query vecDB", "collection", collection, "queryTexts", queryTexts, "embeddingsModel", v.embeddingsModel, "nResults", nResults)
	col, err := v.GetCollection(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("get collection %s: %w", collection, err)
	}
	qr, err := col.Query(ctx, queryTexts, nResults, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("query colection: %w", err)
	}
	results := make([]QueryResult, len(queryTexts))
	for idx, question := range queryTexts {
		res := QueryResult{
			Question:  question,
			Documents: make([]QueryDocument, len(qr.Documents[idx])),
		}
		for i := range qr.Documents[idx] {
			doc := QueryDocument{
				EmbedContent: qr.Documents[idx][i],
				Distance:     qr.Distances[idx][i],
			}
			metaData := qr.Metadatas[idx][i]
			if m, ok := metaData[MetaUpdated].(string); ok {
				doc.Modified = m
			}
			if u, ok := metaData[MetaURL].(string); ok {
				doc.URL = u
			}
			if t, ok := metaData[MetaTitle].(string); ok {
				doc.Title = t
			}
			if t, ok := metaData[MetaIDKey].(string); ok {
				doc.IDField = metaData[t].(string)
			}
			if d, ok := metaData[MetaOrigDoc].(string); ok {
				doc.Document = d
			}

			res.Documents[i] = doc
		}
		results[idx] = res
	}

	return results, nil
}
