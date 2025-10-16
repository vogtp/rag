package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
)

func (srv *Server) vecDBlist(w http.ResponseWriter, r *http.Request) {
	slog := slog.With("url", r.URL)
	slog.Info("Collection list requested")
	var data = struct {
		*commonData
		Collections []*chroma.Collection
		Path        string
	}{
		commonData: srv.common("List collections", r),
		Path:       r.URL.Path,
	}
	ctx := r.Context()

	cols, err := srv.rag(r).ListCollections(ctx, slog)
	if err != nil {
		slog.Error("Error listing vectorDB collections", "err", err)
		srv.Error(w, r, fmt.Sprintf("Error listing vectorDB collections: %v", err), http.StatusInternalServerError)
		return
	}

	sort.Slice(cols, func(i, j int) bool {
		return strings.Compare(cols[i].Name, cols[j].Name) < 0
	})

	data.Collections = cols
	srv.render(w, r, "vecdb_list.gohtml", data)
}
