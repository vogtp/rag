package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func generateState() string {
	return uuid.NewString()
}

func (om *mux) init(ctx context.Context, slog *slog.Logger) error {
	client := &http.Client{
		Timeout:   time.Minute,
		Transport: getSocksProxy(om.slog),
	}
	// enable outgoing request logging
	logging.EnableHTTPClient(client,
		logging.WithClientGroup("client"),
	)

	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
		rp.WithHTTPClient(client),
		rp.WithLogger(slog),
		rp.WithSigningAlgsFromDiscovery(),
		rp.WithErrorHandler(om.oidcErrorHandler),
		rp.WithUnauthorizedHandler(om.oidcUnauthHandler),
	}
	if om.cfg.ClientSecret == "" {
		options = append(options, rp.WithPKCE(cookieHandler))
	}
	if om.cfg.KeyPath != "" {
		options = append(options, rp.WithJWTProfile(rp.SignerFromKeyPath(om.cfg.KeyPath)))
	}

	// One can add a logger to the context,
	// pre-defining log attributes as required.
	ctx = logging.ToContext(ctx, slog)
	var err error
	om.providerOIDC, err = rp.NewRelyingPartyOIDC(ctx, om.cfg.Issuer, om.cfg.ClientID, om.cfg.ClientSecret, om.cfg.RedirectURI, om.scopes, options...)
	if err != nil {
		return fmt.Errorf("cannot register OIDC relying party: %w", err)
	}

	om.urlOptionsOIDC = []rp.URLParamOpt{
		rp.WithPromptURLParam("Welcome back!"),
	}

	if om.responseMode != "" {
		om.urlOptionsOIDC = append(om.urlOptionsOIDC, rp.WithResponseModeURLParam(oidc.ResponseMode(om.responseMode)))
	}
	om.reqisterLoginHandler()

	return nil
}
