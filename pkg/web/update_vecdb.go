package web

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
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
		if err := srv.embeddConfluence(ctx); err != nil {
			srv.slog.Error("Cannot embedd confluence", "err", err)
		}
		slog.Warn("Finished vector DB update", "duration", time.Since(start))
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}()
	return nil
}

func (srv *Server) embeddConfluence(ctx context.Context) error {
	collectionName := viper.GetString(cfg.VecDBColName)
	last := time.Now()
	for _, v := range srv.lastEmbedd {
		if last.After(v) {
			last = v
		}
	}
	if time.Since(last) < viper.GetDuration(cfg.VecDBUpdateIntervall) {
		return fmt.Errorf("Not updating collection %s since it was updated %v ago", collectionName, time.Since(srv.lastEmbedd[collectionName]))
	}
	return confluence.Embbed(ctx, srv.slog, collectionName)
}
