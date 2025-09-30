package web

import (
	"encoding/json"
	"net/http"
)

func (srv *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	srv.slog.Info("Summary request", "url", r.URL.String(), "remote", r.RemoteAddr)
	ctx := r.Context()
	userName := r.PathValue("name")
	if len(userName) < 1 {
		//FIXME get username from oidc
		userName = "vogtp"
	}

	user, err := srv.usercfg.ByName(ctx, userName)
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
