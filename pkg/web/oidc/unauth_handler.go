package oidc

import (
	"net/http"
)

func (om *mux) oidcUnauthHandler(w http.ResponseWriter, r *http.Request, desc string, state string) {
	om.slog.Warn("OIDC unauthorised", "desc", desc, "state", state, "url", r.URL.String())
	http.Redirect(w, r, om.loginPath, http.StatusTemporaryRedirect)
}
