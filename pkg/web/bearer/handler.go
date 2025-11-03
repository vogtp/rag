package bearer

import (
	"context"
	"net/http"
	"strings"
)

type ctxValue string

const bearerCtxValue = ctxValue("bearer")

func HandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, ok := readBearer(w, r)
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), bearerCtxValue, t)
		r = r.WithContext(ctx)
		http.Handler(handler).ServeHTTP(w, r)
	}
}

func readBearer(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r == nil {
		return "", false
	}
	if t, ok := r.Context().Value(bearerCtxValue).(string); ok {
		return t, ok
	}
	// Authorization: Bearer <token>
	ba := r.Header.Get("Authorization")
	if !strings.HasPrefix(ba, "Bearer") {
		http.Error(w, "No bearer authorization header found", http.StatusUnauthorized)
		return "", false
	}
	s := strings.Split(ba, " ")
	if len(s) < 2 {
		http.Error(w, "No bearer token header found", http.StatusUnauthorized)
		return "", false
	}
	return s[1], true
}

func Get(r *http.Request) (string, bool) {
	if r == nil {
		return "", false
	}
	if s, ok := r.Context().Value(bearerCtxValue).(string); ok {
		return s, ok
	}
	// Authorization: Bearer <token>
	ba := r.Header.Get("Authorization")
	if !strings.HasPrefix(ba, "Bearer") {
		return "", false
	}
	s := strings.Split(ba, " ")
	if len(s) < 2 {
		return "", false
	}
	return s[1], true
}
