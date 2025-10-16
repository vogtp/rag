package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/types"
)

// TODO fix summary handler
func (srv *Server) rag(r *http.Request) types.Instance {
	rag := srv.ragMgr.AllFromRequest(r.Context(), r)
	return rag
}
