package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/web/oidc"
)

func (srv *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	srv.slog.Info("User config request", "url", r.URL.String(), "remote", r.RemoteAddr, "method", r.Method)
	switch r.Method {
	case http.MethodGet:
		srv.loadUser(w, r)
		return
	case http.MethodPut:
		srv.saveUser(w, r)
		return
	default:
		http.Error(w, fmt.Sprintf("Unsupported Method %s", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func (srv *Server) saveUser(w http.ResponseWriter, r *http.Request) {
	userName, err := oidc.UserName(r)
	if len(userName) < 1 {
		srv.slog.Warn("No user to save to found", "err", err)
		http.Error(w, "No authenticated user found", http.StatusUnauthorized)
		return

	}
	usr := &usercfg.User{}
	if err := json.NewDecoder(r.Body).Decode(usr); err != nil {
		srv.slog.Warn("Internal server error: decode ent user", "err", err)
		http.Error(w, fmt.Sprintf("decode user settings: %v", err), http.StatusInternalServerError)
	}
	srv.slog.Debug("Saved user settings", "data", usr)
	if !strings.EqualFold(userName, usr.Name) {
		http.Error(w, fmt.Sprintf("User setting %q does not match oidc user %q", usr.Name, userName), http.StatusNotAcceptable)
		return
	}

	if err := srv.usercfg.Add(r.Context(), usr); err != nil {
		srv.slog.Warn("Cannot save user setting", "usersetting", usr, "err", err)
		http.Error(w, fmt.Sprintf("cannot save user setting: %v", err), http.StatusInternalServerError)
		return
	}
	go func() {
		srv.slog.Warn("Updating user rag after config saved", "user", userName, "ctx.err", srv.srvCtx.Err())
		rag := srv.ragMgr.UserFromRequest(srv.srvCtx, r)
		if err := rag.Embbed(srv.srvCtx, srv.slog); err != nil {
			srv.slog.Warn("Failed embed user rag", "err", err, "user", userName)
		}
	}()
}

func (srv *Server) loadUser(w http.ResponseWriter, r *http.Request) {
	userName, err := oidc.UserName(r)
	if len(userName) < 1 {
		srv.slog.Warn("No user to load from found", "err", err)
		http.Error(w, "No authenticated user found", http.StatusUnauthorized)
		return

	}

	user, err := srv.usercfg.User(r.Context(), userName)
	if err != nil {
		srv.slog.Info("Creating user config", "userName", userName)
		user = &usercfg.User{
			Name: userName,
			Collections: []usercfg.Collection{
				{
					DisplayName: "Intranet",
					Source: usercfg.SourceSystem{
						Name: viper.GetString(cfg.ConfluenceBaseURL),
						URL:  viper.GetString(cfg.ConfluenceBaseURL),
						Type: usercfg.SourceConfluence,
					},
				},
			},
		}
		//user, err = srv.usercfgEnt.CreateUser(r.Context(), userName)
	}
	if err != nil {
		srv.slog.Warn("Internal server error: get user ", "err", err)
		http.Error(w, fmt.Sprintf("could not create user config: %v %T", err, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		srv.slog.Warn("Internal server error: encode user response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
