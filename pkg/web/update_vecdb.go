package web

import (
	"context"
	"time"

	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
)

type getAllRager interface {
	GetAllRags() []rag.Instance
}

func (srv *Server) schedulePeriodicVecDBUpdates(ctx context.Context) {
	ragger, ok := srv.ragMgr.(getAllRager)
	if !ok {
		panic("Cannot get all RAGs")
	}
	for _, rag := range ragger.GetAllRags() {
		updateIntervall := rag.UpdateIntervall()
		slog := srv.slog.With("rag", rag.Name(), "updateIntervall", updateIntervall.String())
		if updateIntervall < cfg.MinVecDBUpdateIntervall {
			slog.Warn("Not starting periodic vector DB updates since update intervall is too short")
			continue
		}
		ticker := time.NewTicker(updateIntervall)
		go func() {
			start := time.Now()
			if err := rag.Embbed(ctx); err != nil {
				slog.Error("Cannot embed", "err", err)
			}
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
			slog.Warn("Finished vector DB update", "duration", time.Since(start))
		}()
	}
}
