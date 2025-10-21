package vecdb

import (
	"context"
	"fmt"
	"log/slog"
	sl "log/slog"
	"strings"
	"time"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/vogtp/rag/pkg/logger"
)

const (
	// MetaIDKey is the name of the key which identifies the unique value
	MetaIDKey = "IDkey"

	MetaOrigDoc = "document_original"
	MetaPath    = "path"
	MetaIsRag   = "RAG"
	MetaCreated = "created"
	MetaUpdated = "updated"
	MetaURL     = "URL"
	MetaTitle   = "title"
)

const timeFormat = "2006-01-02 15:04:05.999999999 -0700"

// const timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
// 2025-09-11 16:30:36.875593874 +0200 CEST m=+207.045403821
// 2006-01-02 15:04:05.999999999 -0700
func parseTime(t string) (time.Time, error) {
	if len(t) > len(timeFormat) {
		t = t[:len(timeFormat)]
	}
	return time.Parse(timeFormat, t)
}

func (v *VecDB) Embedd(ctx context.Context, slog *slog.Logger, collectionName string, in <-chan *EmbeddDocument, filters ...Filter) (int, error) {
	if len(collectionName) < 1 {
		slog.Info("No collections name given")
		return 0, fmt.Errorf("no collection name to embed to")
	}
	slog = slog.With("collection", collectionName)
	slogBase := slog
	slog.Warn("Starting embedding", logger.Stacktrace())
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return 0, err
	}
	coll, err := v.GetCollection(ctx, collectionName)
	if err != nil {
		coll, err = v.CreateCollection(ctx, collectionName, map[string]interface{}{MetaIsRag: true, MetaCreated: time.Now().Unix})
		if err != nil {
			slog.Error("cannot create collection", "collectionName", collectionName, "err", err)
			return 0, fmt.Errorf("failed to create collection: %v", err)
		}
	}
	coll.Metadata[MetaUpdated] = time.Now()
	docUpdated := 0

	for d := range in {
		if strings.HasPrefix(d.Title, "Arbeitszeit, Ferien & unbezahlter Urlaub") {
			//FIXME only for debug
			slog.Error("DEBUG found Arbeitszeit, Ferien & unbezahlter Urlaub")
		}
		slog = slogBase.With(sl.Group("RecordID", sl.String(d.IDMetaKey, d.IDMetaValue)))
		for _, f := range filters {
			if !f.ShouldEmbedd(d) {
				continue
			}
		}
		res, err := coll.Get(ctx, map[string]interface{}{d.IDMetaKey: d.IDMetaValue}, nil, nil, nil)
		if err != nil {
			slog.Warn("cannot check for existing docs", "err", err, "title", d.Title)
			continue
		}
		existCnt := len(res.Documents)
		slog = slog.With("existing_records", existCnt)
		skipFile := existCnt > 0
		for _, m := range res.Metadatas {
			if u, ok := m[MetaUpdated].(string); ok {
				t, err := parseTime(u)
				if err != nil {
					slog.Info("Cannot parse update time", "time", u)
					skipFile = false
				}
				if d.Modified.After(t) {
					skipFile = false
					slog.Debug("Document was modified", "modification_time", d.Modified.String())
					break
				}
			} else {
				skipFile = false
				slog.Warn("cannot read meta data as string", "meta", MetaUpdated, "value", m[MetaUpdated])
			}
		}
		if skipFile {
			slog.Info("document allready exists and not updated")
			continue
		}
		for _, s := range d.Split(slog) {
			rs, err := types.NewRecordSet(
				types.WithEmbeddingFunction(embedFunc),
				types.WithIDGenerator(types.NewULIDGenerator()),
			)
			if err != nil {
				slog.Warn("cannot create record set", "err", err)
				continue
			}

			metadata := []types.Option{
				types.WithDocument(s),
				types.WithID(d.IDMetaValue),
				types.WithMetadata(MetaOrigDoc, d.Document),
				types.WithMetadata(d.IDMetaKey, d.IDMetaValue),
				types.WithMetadata(MetaIDKey, d.IDMetaKey),
				types.WithMetadata(MetaUpdated, d.Modified.String()),
			}
			if len(d.URL) > 0 {
				metadata = append(metadata, types.WithMetadata(MetaURL, d.URL))
			}
			if len(d.Title) > 0 {
				metadata = append(metadata, types.WithMetadata(MetaTitle, d.Title))
			}
			for k, v := range d.MetaData {
				metadata = append(metadata, types.WithMetadata(k, v))
			}
			rs.WithRecord(metadata...)

			_, err = rs.BuildAndValidate(ctx)
			if err != nil {
				slog.Debug("cannot validate document", "err", err, "rs", rs)
				if err.Error() == "document cannot be empty" {
					slog.Info("document not validated", "err", err)
					continue
				}
				slog.Warn("document not validated", "err", err)
				continue
				//return fmt.Errorf("error validating record set: %s \n", err)
			}
			// Add the records to the collection
			ids := rs.GetIDs()
			if len(ids) == len(res.Ids) {
				ids = res.Ids
				slog.Debug("Using IDs from existing document")
			}
			slog.DebugContext(ctx, "Embedding document")
			_, err = coll.Upsert(ctx, rs.GetEmbeddings(), rs.GetMetadatas(), rs.GetDocuments(), ids)
			if err != nil {
				slog.Warn("cannot add document", "err", err)
				continue
			}
		}
		for _, f := range filters {
			f.ReqisterEmedded(d)
		}
		docUpdated++
	}

	// Count the number of documents in the collection
	countDocs, qrerr := coll.Count(ctx)
	if qrerr != nil {
		return docUpdated, fmt.Errorf("error counting documents: %w", qrerr)
	}

	slog.Info("Finished embedding", "docsCount", countDocs, "docsUpdates", docUpdated)

	return docUpdated, nil
}
