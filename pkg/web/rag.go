package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/rag"
)

// TODO fix summary handler
func (srv *Server) rag(r *http.Request) rag.Instance {
	rag, err := srv.ragMgr.AllFromRequest(r)
	if err != nil {
		srv.slog.Warn("Cannot load rag from request", "err", err)
	}
	return rag
}
