package usercfg

import (
	"context"

	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/sourcesystem"
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
func (db *DB) CreateUser(ctx context.Context, name string) (*ent.User, error) {
	uc := db.User.Create().SetName(name)
	conflDefault, err := db.SourceSystem.Create().SetName("Intranet").SetURL("https://intranet.unibas.ch/").SetType(sourcesystem.TypeConfluence).Save(ctx)
	if err != nil {
		return nil, err
	}
	coll, err := db.Collection.Create().SetName("Collection").AddSources(conflDefault).Save(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := uc.AddCollections(coll).Save(ctx); err != nil {
		return nil, err
	}
	return db.ByName(ctx, name)
}
