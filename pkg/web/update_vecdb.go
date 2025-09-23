package web

import (
	"context"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
)

func (srv *Server) schedulePeriodicVecDBUpdates(ctx context.Context) error {
	updateIntervall := viper.GetDuration(cfg.VecDBUpdateIntervall)
	if updateIntervall < time.Hour {
		slog.Warn("Not starting periodic vector DB updates since update intervall is too short", "updateIntervall", updateIntervall)
		return nil
	}
	ticker := time.NewTicker(updateIntervall)
	go func() {
		start := time.Now()
		if err := srv.rag.Embbed(ctx); err != nil {
			srv.slog.Error("Cannot embedd confluence", "err", err)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
		slog.Warn("Finished vector DB update", "duration", time.Since(start))
	}()
	return nil
}
