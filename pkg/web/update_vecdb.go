package web

import (
	"context"
	"time"
)

func (srv *Server) schedulePeriodicVecDBUpdates(ctx context.Context) {
	for _, rag := range srv.ragManagers {
		updateIntervall := rag.UpdateIntervall()
		slog := srv.slog.With("rag", rag.Name(), "updateIntervall", updateIntervall.String())
		if updateIntervall < time.Hour {
			slog.Warn("Not starting periodic vector DB updates since update intervall is too short")
			return
		}
		ticker := time.NewTicker(updateIntervall)
		go func() {
			start := time.Now()
			if err := rag.Embbed(ctx); err != nil {
				slog.Error("Cannot embedd", "err", err)
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
