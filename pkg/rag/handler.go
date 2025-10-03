package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vogtp/rag/pkg/cfg"
)

type Handler interface {
	FromRequest(r *http.Request) Instance
}

var _ (Handler) = (*handler)(nil)

type handler struct {
	slog       *slog.Logger
	globalRags []Instance
}

func New(ctx context.Context, slog *slog.Logger) (Handler, error) {
	ragCfgs, err := cfg.GetRagConfig()
	if err != nil {
		return nil, fmt.Errorf("read RAG config: %w", err)
	}
	rags := make([]Instance, len(ragCfgs))
	for i, ragCfg := range ragCfgs {
		r, err := newInstanceCfg(ctx, slog, ragCfg)
		if err != nil {
			return nil, fmt.Errorf("start rag %q backend: %w", ragCfg.Name, err)
		}
		rags[i] = r
	}
	m := &handler{
		slog:       slog,
		globalRags: rags,
	}
	return m, nil
}

func (h handler) FromRequest(r *http.Request) Instance {
	insts := h.globalInstances()

	// by authenticated User
	//	u, err := oidc.getUserName(r)

	// by Bearer Token

	// fix summary handler

	return insts
}

func (h handler) globalInstances() Instance {
	return newInstanceList(h.slog, "public", h.globalRags...)
}

// GetAllRags returns a slice of all instance without any authentical
// for internal use only
func (h handler) GetAllRags() []Instance {
	return h.globalRags
}
