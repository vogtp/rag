package usercfg

import (
	"context"

	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/user"
)

func (db *DB) All(ctx context.Context) ([]*ent.User, error) {
	return db.GetUserQuery(ctx).All(ctx)
}

func (db *DB) ByName(ctx context.Context, username string) (*ent.User, error) {
	return db.GetUserQuery(ctx).Where(user.Name(username)).First(ctx)
}

func (db *DB) GetUserQuery(ctx context.Context) *ent.UserQuery {
	uq := db.User.Query()
	// get collections
	uq = uq.WithCollections(func(cq *ent.CollectionQuery) {
		cq.WithSources()
	})
	return uq
}
