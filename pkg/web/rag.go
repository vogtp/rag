package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/types"
)

// TODO fix summary handler
func (srv *Server) rag(r *http.Request) types.Instance {
	rag := srv.ragMgr.FromRequest(r.Context(), r)
	return rag
}
