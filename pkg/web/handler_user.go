package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/web/oidc"
)

type Sessioner interface {
	GetSession(w http.ResponseWriter, r *http.Request) (*oidc.Session, error)
}

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
	userName, err := srv.getUserName(w, r)
	if len(userName) < 1 {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
		return

	}
	usr := &ent.User{}
	if err:=json.NewDecoder(r.Body).Decode(usr); err!=nil{
		http.Error(w,fmt.Sprintf("decode user settings: %v",err),http.StatusInternalServerError)
	}
	srv.slog.Debug("Saved user settings","data",usr)
	if !strings.EqualFold(userName, usr.Name) {
		http.Error(w, fmt.Sprintf("User setting %q does not match oidc user %q", usr.Name, userName), http.StatusNotAcceptable)
		return
	}

	if err := srv.usercfg.SaveUser(r.Context(), usr);err != nil {
		srv.slog.Warn("Cannot save user setting", "usersetting", usr, "err",err)
		http.Error(w, fmt.Sprintf("cannot save user setting: %v", err), http.StatusInternalServerError)
		return
	}
}

func (srv *Server) loadUser(w http.ResponseWriter, r *http.Request) {
	userName, err := srv.getUserName(w, r)
	if len(userName) < 1 {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
		return

	}

	user, err := srv.usercfg.ByName(r.Context(), userName)
	if err != nil {
		srv.slog.Info("Creating user config", "userName", userName)
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
		// return "", fmt.Errorf("no session found")
		srv.slog.Error("USING HARDCODED USER")
		return "vogtp", nil // FIXME Debug only
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
