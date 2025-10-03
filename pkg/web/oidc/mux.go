package oidc

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// NewMux creates a new OIDC authenticated mux
func NewMux(ctx context.Context, slog *slog.Logger, serveMux *http.ServeMux, addr string, cfg Config) (Mux, error) {
	om := &mux{
		slog:          slog.With("oidc", "oidc"),
		serveMux:      serveMux,
		sessionMaxAge: 24 * time.Hour,
		addr:          addr,
		cfg:           cfg,
	}
	redURI, err := url.Parse(cfg.RedirectURI)
	if err != nil {
		return nil, err
	}
	om.callbackPath = redURI.Path
	om.loginPath = cfg.LoginPath
	if len(cfg.LoginPath) < 1 {
		om.loginPath = "/login"
	}
	om.scopes = cfg.Scopes
	if len(om.scopes) < 1 {
		om.scopes = []string{oidc.ScopeOpenID, oidc.ScopeProfile}
	}
	om.responseMode = cfg.ResponseMode
	if len(om.responseMode) < 1 {
		om.responseMode = "code token"
	}
	if cfg.sessionMaxAge > 0 {
		om.sessionMaxAge = cfg.sessionMaxAge
	}
	if err := om.init(ctx, slog); err != nil {
		return nil, err
	}
	cookieHandler = httphelper.NewCookieHandler(securecookie.GenerateRandomKey(24), securecookie.GenerateRandomKey(24), httphelper.WithMaxAge(int(om.sessionMaxAge.Seconds())))
	return om, nil
}

var cookieHandler *httphelper.CookieHandler

// mux OIDC authenticated mux
type mux struct {
	slog *slog.Logger

	serveMux *http.ServeMux
	addr     string
	cfg      Config

	loginPath    string
	callbackPath string
	scopes       []string
	responseMode string

	sessionMaxAge  time.Duration
	providerOIDC   rp.RelyingParty
	stateOIDC      func() string
	urlOptionsOIDC []rp.URLParamOpt
}
