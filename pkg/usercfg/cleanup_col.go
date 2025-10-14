package usercfg

import (
	"context"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func (d *DataBase) CleanupUserCollections(ctx context.Context, usr *User) error {
	slog := d.slog.With("cleanup", usr.Name)
	slog.Info("Cleanup user collecions")
	colMap := make(map[string]bool)
	for _, col := range usr.Collections {
		colMap[col.CollectionName] = true
	}
	vecDb, err := vecdb.New(ctx, slog, viper.GetString(cfg.ModelEmbedding))
	if err != nil {
		return err
	}
	cols, err := vecDb.ListCollections(ctx, usr.Name)
	if err != nil {
		return err
	}
	if len(cols) < 1 {
		slog.Info("No collections to cleanup")
		return nil
	}
	for _, col := range cols {
		if colMap[col.Name] {
			continue
		}
		slog.Info("Removing unused user collection", "collection", col.Name)
		if err := vecDb.DeleteCollection(ctx, col.Name); err != nil {
			slog.Warn("Cannot delete collection", "collection", col.Name)
		}
	}
	return nil
}
