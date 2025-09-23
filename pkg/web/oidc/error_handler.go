package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (om *mux) oidcErrorHandler(w http.ResponseWriter, r *http.Request, errorType string, errorDesc string, state string) {

	info := struct {
		Hanlder   string
		ErrorType string
		ErrorDesc string
		State     string
	}{
		Hanlder:   "Error",
		ErrorType: errorType,
		ErrorDesc: errorDesc,
		State:     state,
	}
	om.slog.Warn("OIDC error", "errorDesc", errorDesc, errorType, "errorType", "state", state)
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(&info)

	if err != nil {
		http.Error(w, fmt.Sprintf("oidcErrorHandler cannot marshal json: %v", err), http.StatusInternalServerError)
		return
	}
}
