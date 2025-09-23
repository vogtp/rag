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

	//SessionMaxAge is the max age of a internal session, after it gets revalidated (possible transparent)
	SessionMaxAge = 5 * time.Minute
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
	User    *oidc.UserInfo
	Created time.Time
}

func (om *mux) setSession(w http.ResponseWriter, info *oidc.UserInfo) error {
	session := Session{
		User:    info,
		Created: time.Now(),
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&session); err != nil {
		return err
	}
	om.slog.Info("New session info", slog.Group("session", slog.Any("orig", session), slog.String("encoded", buf.String())))
	return om.cookieHandler.SetCookie(w, sessionCookieName, buf.String())
}

// GetSession returns session information with user info
func (om *mux) GetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	s, err := om.cookieHandler.CheckCookie(r, sessionCookieName)
	if err != nil {
		return nil, fmt.Errorf("get session cookie: %w", err)
	}
	session := &Session{
		User: &oidc.UserInfo{},
	}
	if err := json.Unmarshal([]byte(s), session); err != nil {
		return nil, fmt.Errorf("unmarshal session cookie: %w", err)
	}
	if time.Since(session.Created) > SessionMaxAge {
		a := time.Since(session.Created).Truncate(time.Second).String()
		om.cookieHandler.DeleteCookie(w, sessionCookieName)
		om.slog.Info("Session expired", "age", a, "maxAge", SessionMaxAge.String(), "session", session)
		return nil, fmt.Errorf("session expired. age: %s", a)
	}
	return session, err
}
