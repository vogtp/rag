package usercfg

import (
	"context"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
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

func NewENT(ctx context.Context, slog *slog.Logger, driverName string, dataSourceName string) (*DB, error) {
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
