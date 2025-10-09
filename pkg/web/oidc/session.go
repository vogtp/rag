package oidc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var (
	// ParamOrigPath is the original path
	ParamOrigPath = "OrigPath"
	//DefaultPath is the default path to redirect to
	DefaultPath = "/"
)

func getOrigPath(r *http.Request) string {
	if err := r.ParseForm(); err != nil {
		return DefaultPath
	}
	op := r.Form.Get(ParamOrigPath)
	if len(op) < 1 {
		return DefaultPath
	}
	return op
}

const (
	sessionCookieName = "session"
)

// Session represents an authorised OIDC session
type Session struct {
	User     *oidc.UserInfo
	UserName string
	Created  time.Time
}

func (om *mux) setSession(w http.ResponseWriter, info *oidc.UserInfo) error {
	session := Session{
		User:     info,
		UserName: fmt.Sprintf("%v", info.Claims["subname"]),
		Created:  time.Now(),
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&session); err != nil {
		return err
	}
	om.slog.Info("New session info", slog.Group("session", slog.Any("orig", session), slog.String("encoded", buf.String())))
	return cookieHandler.SetCookie(w, sessionCookieName, buf.String())
}

// GetSession returns session information with user info
func GetSession(r *http.Request) (*Session, error) {
	s, err := cookieHandler.CheckCookie(r, sessionCookieName)
	if err != nil {
		return nil, fmt.Errorf("get session cookie: %w", err)
	}
	session := &Session{
		User: &oidc.UserInfo{},
	}
	if err := json.Unmarshal([]byte(s), session); err != nil {
		return nil, fmt.Errorf("unmarshal session cookie: %w", err)
	}

	return session, err
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookieHandler.DeleteCookie(w, sessionCookieName)
}

func UserName(r *http.Request) (string, error) {
	sess, err := GetSession(r)
	if err != nil {
		slog.Info("Cannot get session to get username", "err", err)
		return "", fmt.Errorf("get user from session: %w", err)
	}
	return sess.UserName, nil
}
