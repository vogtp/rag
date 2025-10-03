package bearer

import (
	"net/http"
)

type Auth interface {
	// Authorise reads the Bearer token from the Authorization header
	// sets the response message and code
	// returns true if authorised and false otherwise
	Authorise(w http.ResponseWriter, r *http.Request) bool
}

func NoAuth() Auth { return noAuth{} }

type noAuth struct {
}

func (noAuth) Authorise(_ http.ResponseWriter, _ *http.Request) bool {
	return true
}

func TokenAuth(token string) Auth {
	return tokenAuth{token: token}
}

type tokenAuth struct {
	token string
}

func (ta tokenAuth) Authorise(w http.ResponseWriter, r *http.Request) bool {
	if len(ta.token) < 1 {
		return true
	}
	t, ok := readBearer(w, r)
	// slog.Warn(" auth check", "ok", ok, "t", t, "token", ta.token)
	// if t != ta.token {
	// 	http.Error(w, "Could not authorise request", http.StatusUnauthorized)
	// }
	return ok && t == ta.token
}
