package rag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vogtp/rag/pkg/cfg"
)

type Handler interface {
	Public() Instance
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

func (h handler) Public() Instance {
	return newInstanceList(h.slog, "public", h.globalRags...)
}

// GetAllRags returns a slice of all instance without any authentical
// for internal use only
func (h handler) GetAllRags() []Instance {
	return h.globalRags
}
