package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vogtp/rag/pkg/web/oidc"
)

type Sessioner interface {
	GetSession(w http.ResponseWriter, r *http.Request) (*oidc.Session, error)
}

func (srv *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	srv.slog.Info("Summary request", "url", r.URL.String(), "remote", r.RemoteAddr)

	userName, err := srv.getUserName(w, r)
	if len(userName) < 1 {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
		return

	}

	user, err := srv.usercfg.ByName(r.Context(), userName)
	if err != nil {
		srv.slog.Warn("Creating user config", "userName", userName)
		user, err = srv.usercfg.CreateUser(r.Context(), userName)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create user config: %v %T", err, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (srv *Server) getUserName(w http.ResponseWriter, r *http.Request) (string, error) {
	sessioner, ok := srv.oidcMux.(Sessioner)
	if !ok {
		return "", fmt.Errorf("no session found")
		// srv.slog.Error("USING HARDCODED USER")
		// return  "vogtp", nil // FIXME Debug only
	}
	sess, err := sessioner.GetSession(w, r)
	if err != nil {
		return "", err
	}
	userName := sess.UserName
	if len(userName) < 1 {
		return "", err

	}
	srv.slog.Info("found user session", "session", sess, "userName", userName)
	return userName, nil
}
