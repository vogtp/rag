package web

import (
	"net/http"

	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

// FIXME use or remove
// Error is a wrapper for srv.Error
func (srv *Server) Error(w http.ResponseWriter, r *http.Request, errStr string, code int) {
	data := struct {
		*commonData
		Error     string
		Code      int
		Documents []vecdb.QueryDocument
	}{
		commonData: srv.common(errStr, r),
		Error:      errStr,
		Code:       code,
	}
	h := w.Header()
	h.Del("Content-Length")
	h.Set("X-Content-Type-Options", "nosniff")
	// http.Error(w, errStr, code)
	w.WriteHeader(code)
	srv.render(w, r, "error.gohtml", data)
}
