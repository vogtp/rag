package web

import (
	"net/http"

	"github.com/vogtp/rag/pkg/rag"
)

func (srv *Server) rag(r *http.Request) rag.Instance {

	// by authenticated User

	// by Bearer Token

	// fix summary handler

	return srv.ragHandlers.Public()
}
