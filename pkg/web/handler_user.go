package web

import (
	"encoding/json"
	"net/http"

	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/user"
)

func (srv *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	srv.slog.Info("Summary request", "url", r.URL.String(), "remote", r.RemoteAddr)
	ctx := r.Context()
	userName := r.PathValue("name")

	user, err := srv.usercfg.User.Query().WithConfluence(func(cq *ent.ConfluenceQuery) { cq.WithSpaces() }).Where(user.Name(userName)).First(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
