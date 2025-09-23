package oidc

import (
	"net/http"
)

// oidcMux can hold a oidc authenticated mux or a standard one
type Mux interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
	Handle(string, http.Handler)
}

// Handle makes shure it is authenticated
func (om *mux) Handle(pattern string, handler http.Handler) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		_, err := om.GetSession(w, r)
		if err != nil {
			om.slog.Warn("Not authorised", "err", err)
			http.Redirect(w, r, om.loginPath, http.StatusTemporaryRedirect)
			return
		}
		handler.ServeHTTP(w, r)
	}
	om.serveMux.HandleFunc(pattern, handleFunc)
}

// HandleFunc makes shure it is authenticated
func (om *mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	om.Handle(pattern, http.HandlerFunc(handler))
}
