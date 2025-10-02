package rag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vogtp/rag/pkg/cfg"
)

type manager struct {
}

func New(ctx context.Context, slog *slog.Logger) ([]Manager, error) {
	ragCfgs, err := cfg.GetRagConfig()
	if err != nil {
		return nil, fmt.Errorf("read RAG config: %w", err)
	}
	rags := make([]Manager, len(ragCfgs))
	for i, ragCfg := range ragCfgs {
		r, err := newInstance(ctx, slog, ragCfg)
		if err != nil {
			return nil, fmt.Errorf("start rag %q backend: %w", ragCfg.Name, err)
		}
		rags[i] = r
	}
	return rags, nil
}
