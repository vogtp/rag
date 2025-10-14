package usercfg

import (
	"context"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	//_ "github.com/xiaoqidun/entps" // needed to acces sqlite
)

const (
	Dialect    = dialect.SQLite
	DBFileName = "rag.sqlite"
)

type DB struct {
	*ent.Client
	srvCtx context.Context
	slog   *slog.Logger
}

func newENT(ctx context.Context, slog *slog.Logger, driverName string, dataSourceName string) (*DB, error) {
	if driverName == dialect.SQLite {
		dataSourceName = fmt.Sprintf("file:%s?&cache=shared&_fk=1", dataSourceName)
	}
	entClient, err := ent.Open(dialect.SQLite, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open ent db client %s: %w", dataSourceName, err)
	}
	if err := entClient.Schema.Create(ctx, schema.WithGlobalUniqueID(true)); err != nil { //
		return nil, fmt.Errorf("creating schema resources: %v", err)
	}
	dbClient := &DB{
		srvCtx: ctx,
		slog:   slog,
		Client: entClient,
	}

	//dbClient.User.Query().Where(user.Name())
	return dbClient, nil
}

func (db *DB) CleanupUserCollectionsEnt(ctx context.Context, usr *ent.User) error {
	slog := db.slog.With("cleanup", usr.Name)
	slog.Info("Cleanup user collecions")
	colMap := make(map[string]bool)
	for _, col := range usr.Edges.Collections {
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
