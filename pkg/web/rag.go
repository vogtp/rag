package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/rag"
)

// TODO fix summary handler
func (srv *Server) rag(r *http.Request) rag.Instance {
	rag := srv.ragMgr.AllFromRequest(r)
	return rag
}
