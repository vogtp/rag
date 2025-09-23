package oidc

import (
	"fmt"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func (om *mux) reqisterLoginHandler() {
	om.serveMux.Handle(om.loginPath, rp.AuthURLHandler(
		generateState,
		om.providerOIDC,
		om.urlOptionsOIDC...,
	))
	callbackHandler := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty, info *oidc.UserInfo) {
		// fmt.Println("access token", tokens.AccessToken)
		// fmt.Println("refresh token", tokens.RefreshToken)
		// fmt.Println("id token", tokens.IDToken)

		// data, err := json.Marshal(info)
		// if err != nil {
		// 	http.Error(w, fmt.Sprintf("should redirect to login: %v", err), http.StatusInternalServerError)
		// 	return
		// }
		// w.Header().Set("content-type", "application/json")
		// w.Write(data)
		if err := om.setSession(w, info); err != nil {
			om.slog.Warn("Cannot set session cookie", "err", err)
			http.Error(w, fmt.Sprintf("cannot save session: %v", err), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, getOrigPath(r), http.StatusTemporaryRedirect)
	}

	om.serveMux.Handle(om.callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(callbackHandler), om.providerOIDC))

}
