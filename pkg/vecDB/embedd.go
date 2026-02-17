package vecdb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/vogtp/rag/pkg/logger"
	ragtypes "github.com/vogtp/rag/pkg/types"
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

const (
	timeFormat  = "2006-01-02 15:04:05.999999999 -0700"
	timeFormat2 = "2006-01-02 15:04:05"
)

// const timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
// 2025-09-11 16:30:36.875593874 +0200 CEST m=+207.045403821
// 2006-01-02 15:04:05.999999999 -0700
// 2025-10-23 10:29:07.4234687 +0200 CEST m=+133.822645301
func parseTime(t string) (time.Time, error) {
	if len(t) > len(timeFormat) {
		t = t[:len(timeFormat)]
	}
	t1, err := time.Parse(timeFormat, t)
	if err == nil {
		return t1, nil
	}
	t = t[:len(timeFormat)-1]
	t1, err = time.Parse(timeFormat, t)
	if err == nil {
		return t1, nil
	}
	if len(t) > len(timeFormat2) {
		t = t[:len(timeFormat2)]
	}
	return time.Parse(timeFormat2, t)
}

func (v *VecDB) Embedd(ctx context.Context, slog *slog.Logger, collectionName string, in <-chan *ragtypes.EmbeddDocument, filters ...ragtypes.Filter) (int, error) {
	if len(collectionName) < 1 {
		slog.Info("No collections name given")
		return 0, fmt.Errorf("no collection name to embed to")
	}
	slogBase := slog
	embedder, err := v.GetEmbeddingFunc(ctx)
	if err != nil {
		return 0, err
	}
	slogBase.Warn("Starting embedding", logger.Stacktrace())

	coll, err := v.GetCollection(ctx, collectionName)
	if err != nil {
		coll, err = v.CreateCollection(ctx, collectionName, map[string]interface{}{MetaIsRag: true, MetaCreated: time.Now().Unix})
		if err != nil {
			slogBase.Error("cannot create collection", "collectionName", collectionName, "err", err)
			return 0, fmt.Errorf("failed to create collection: %v", err)
		}
	}
	coll.Metadata[MetaUpdated] = time.Now()
	docUpdated := 0

	for d := range in {
		slog = slogBase.With("doc", d)
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
					slog.Info("Cannot parse update time", "update_time", u)
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
		rs, err := types.NewRecordSet(
			types.WithEmbeddingFunction(embedder),
			//types.WithIDGenerator(types.NewULIDGenerator()),
		)
		if err != nil {
			slog.Warn("cannot create record set", "err", err)
			continue
		}
		recOpts := []types.Option{
			types.WithDocument(d.Document),
			types.WithID(d.IDMetaValue),
			types.WithMetadata(MetaOrigDoc, d.Document),
			types.WithMetadata(d.IDMetaKey, d.IDMetaValue),
			types.WithMetadata(MetaIDKey, d.IDMetaKey),
			types.WithMetadata(MetaUpdated, d.Modified.String()),
		}
		if len(d.URL) > 0 {
			recOpts = append(recOpts, types.WithMetadata(MetaURL, d.URL))
		}
		if len(d.Title) > 0 {
			recOpts = append(recOpts, types.WithMetadata(MetaTitle, d.Title))
		}
		for k, v := range d.MetaData {
			recOpts = append(recOpts, types.WithMetadata(k, v))
		}
		rs.WithRecord(recOpts...)

		_, err = rs.BuildAndValidate(ctx)
		if err != nil {
			slog.Debug("cannot validate document", "err", err, "rs", rs)
			if err.Error() == "document cannot be empty" {
				slog.Info("document not validated", "err", err)
				continue
			}
			slog.Warn("document not validated", "err", err)
			//continue
			//return fmt.Errorf("error validating record set: %s \n", err)
		}
		for i, s := range d.Split(slog) {

			recOpts := []types.Option{
				types.WithDocument(s),
				types.WithID(fmt.Sprintf("%s_idx:%d", d.IDMetaValue, i)),
				types.WithMetadata(MetaOrigDoc, d.Document),
				types.WithMetadata(d.IDMetaKey, d.IDMetaValue),
				types.WithMetadata(MetaIDKey, d.IDMetaKey),
				types.WithMetadata(MetaUpdated, d.Modified.String()),
			}
			if len(d.URL) > 0 {
				recOpts = append(recOpts, types.WithMetadata(MetaURL, d.URL))
			}
			if len(d.Title) > 0 {
				recOpts = append(recOpts, types.WithMetadata(MetaTitle, d.Title))
			}
			for k, v := range d.MetaData {
				recOpts = append(recOpts, types.WithMetadata(k, v))
			}
			rs.WithRecord(recOpts...)

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
		}
		// Add the records to the collection
		ids := rs.GetIDs()
		if len(ids) == len(res.Ids) {
			ids = res.Ids
			slog.Debug("Using IDs from existing document")
		}
		slog.DebugContext(ctx, "Embedding document")
		embs := rs.GetEmbeddings()
		mds := rs.GetMetadatas()
		docs := rs.GetDocuments()
		_, err = coll.Upsert(ctx, embs, mds, docs, ids)
		if err != nil {
			slog.Warn("cannot add document", "err", err)
			continue
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

	slogBase.Info("Finished embedding", "docsCount", countDocs, "docsUpdates", docUpdated)

	return docUpdated, nil
}
