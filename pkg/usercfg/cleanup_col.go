package usercfg

import (
	"context"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func (db *DB) CleanupUserCollections(ctx context.Context, usr *ent.User) error {
	slog := db.slog.With("cleanup", usr.Name)
	slog.Info("Cleanup user collecions")
	colMap := make(map[string]bool)
	for _, col := range usr.Edges.Collections {
		for _, src := range col.Edges.Sources {
			spaces := strings.Split(src.Parts, ",")
			if len(spaces) < 1 {
				spaces = strings.Split(src.Parts, " ")
			}
			spaces = append(spaces, "all")
			colMap[col.CollectionName] = true
			// for _, s := range spaces {
			// 	colName := strings.ToLower(fmt.Sprintf("%s-%s-%s", usr.Name, col.Name, s))
			// 	colMap[colName] = true
			// }
		}
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
		if colMap[strings.ToLower(col.Name)] {
			continue
		}
		slog.Info("Removing unused user collection", "collection", col.Name)
		if err := vecDb.DeleteCollection(ctx, col.Name); err != nil {
			slog.Warn("Cannot delete collection", "collection", col.Name)
		}
	}
	return nil
}
