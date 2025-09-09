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
		if err := srv.embeddConfluence(ctx); err != nil {
			srv.slog.Error("Cannot embedd confluence", "err", err)
		}
		slog.Info("Finished vector DB update")
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}()
	return nil
}

func (srv *Server) embeddConfluence(ctx context.Context) error {
	collectionName := "intranet"
	if time.Since(srv.lastEmbedd[collectionName]) < viper.GetDuration(cfg.VecDBUpdateIntervall) {
		return fmt.Errorf("Not updating collection %s since it was updated %v ago", collectionName, time.Since(srv.lastEmbedd[collectionName]))
	}
	return confluence.Embbed(ctx, srv.slog, collectionName)
}
