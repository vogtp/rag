package web

import (
	"context"
	"time"
	sl "log/slog"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
)

// schedulePeriodicVecDBUpdates starts periodic vec db updates
// must be started as goroutine
func (srv *Server) schedulePeriodicVecDBUpdates(ctx context.Context) {
	delay := viper.GetDuration(cfg.VecDBUpdateDelay)
	if delay > 23*time.Hour {
		srv.slog.Warn("Not doing periodic updates since delay is langer than 23h", "delay", delay.String())
		return
	}
	time.Sleep(delay)
	ragger, ok := srv.ragMgr.(rag.GetAllRager)
	if !ok {
		panic("Cannot get all RAGs")
	}
	for _, rag := range ragger.GetAllRags(ctx) {
		updateIntervall := rag.UpdateIntervall()
		slog := srv.slog.With(sl.Group("rag", "displayName", rag.DisplayName(), "updateIntervall", updateIntervall.String() ))
		if updateIntervall < cfg.MinVecDBUpdateIntervall {
			slog.Warn("Not starting periodic vector DB updates since update intervall is too short")
			continue
		}
		ticker := time.NewTicker(updateIntervall)
		go func() {
			start := time.Now()
			if err := rag.Embbed(ctx, slog); err != nil {
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
