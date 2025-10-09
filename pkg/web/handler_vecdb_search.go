package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (srv *Server) vecDBsearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	collection := r.PathValue("collection")
	if err := r.ParseForm(); err != nil {
		srv.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	query := r.FormValue("query")
	maxResStr := r.FormValue("maxResults")
	slog := srv.slog.With("collection", collection, "query", query, "maxResults", maxResStr)
	maxResults, err := strconv.Atoi(maxResStr)
	if err != nil {
		slog.Info("Cannot convert max Results to int", "err", err)
		maxResults = 10
	}

	slog = srv.slog.With("collection", collection, "query", query, "maxResults", maxResults)
	slog.Info("Collection search requested")
	var data = struct {
		*commonData
		Collection string
		Query      string
		Documents  []queryDoc
	}{
		commonData: srv.common(fmt.Sprintf("Search: %s", collection), r),
		Collection: collection,
		Query:      query,
	}
	if !srv.lastEmbedd[collection].IsZero() {
		data.StatusMessage = fmt.Sprintf("Last %s update: %v", collection, srv.lastEmbedd[collection])
	}
	if len(query) > 0 {
		ctx := r.Context()
		docs, err := srv.rag(r).SearchVecDB(ctx, slog, collection, query, maxResults)
		if err != nil {
			slog.Error("Cannot query vecDB", "err", err)
			srv.Error(w, r, err.Error(), http.StatusInternalServerError)
			return
		}
		sizeOrig := len(docs)
		// srtFunc := func(a, b vecdb.QueryDocument) int {
		// 	return strings.Compare(a.IDField, b.IDField)
		// }
		// slices.SortFunc(docs, srtFunc)
		// cmpFunc := func(a, b vecdb.QueryDocument) bool {
		// 	return strings.EqualFold(a.IDField, b.IDField)
		// }
		// docs = slices.CompactFunc(docs, cmpFunc)
		keyMap := make(map[string]int)
		data.Documents = make([]queryDoc, 0, len(docs))
		for _, d := range docs {
			c, found := keyMap[d.IDField]
			keyMap[d.IDField] = c + 1
			if found {
				slog.Debug("removing document", "title", d.Title, "URL", d.URL, "docCnt", c)
				continue
			}

			qd := queryDoc{
				QueryDocument: &d,
				UUID:          uuid.New(),
			}
			data.Documents = append(data.Documents, qd)
			srv.docCache.add(&qd)
		}

		if sizeOrig > len(docs) {
			slog.Warn("removed doubs", "orig", sizeOrig, "now", len(docs))
		}
	}
	data.StatusMessage = fmt.Sprintf("Duration %v - %v", time.Since(start).Truncate(time.Second), data.StatusMessage)
	srv.render(w, r, "vecdb_search.gohtml", data)
}
