package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/rag"
)

func (srv *Server) rag(r *http.Request) rag.Manager {
	return srv.ragManagers[0]
}
