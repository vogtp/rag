package usercfg

import (
	"context"
	"fmt"
	"strings"

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
			for _, s := range spaces {
				colName := strings.ToLower(fmt.Sprintf("%s-%s-%s", usr.Name, col.Name, s))
				colMap[colName] = true
			}
		}
	}
	vc := cfg.DefaultRagCfg()
	vc.Vecdb.CollectionName = usr.Name
	vecDb, err := vecdb.New(ctx, slog, vc)
	if err != nil {
		return err
	}
	cols, err := vecDb.ListCollections(ctx, vc)
	if err != nil {
		return err
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
