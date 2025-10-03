package usercfg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/collection"
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

func (db *DB) SaveUser(ctx context.Context, usr *ent.User) (err error) {
	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		slog := slog.With("user", usr.Name)
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
		if err != nil {
			slog.Error("SaveUser error doing rollbarg", "err", err)
			if err := tx.Rollback(); err != nil {
				slog.Error("Rollback error", "err", err)
				return
			}
			return
		}
		if err := tx.Commit(); err != nil {
			slog.Error("Commit error", "err", err)
			return
		}
		slog.Info("Saved user settings", "user", usr.Name)
	}()

	usrDB := tx.User.Query().WithCollections(func(cq *ent.CollectionQuery) { cq.WithSources() }).Where(user.Name(usr.Name)).FirstX(ctx)

	colsUsrDB, err := usrDB.Collections(ctx)
	if err != nil {
		return err
	}
	colsMap := make(map[int]*ent.Collection, len(colsUsrDB))
	for _, c := range colsUsrDB {
		colsMap[c.ID] = c
	}
	userUp := usrDB.Update().SetUser(usr)
	cols := usr.Edges.Collections
	for _, c := range cols {
		keys := tx.Collection.Query().Where(collection.APIKey(c.APIKey)).AllX(ctx)
		for _, k := range keys {
			if k.ID != c.ID {
				return fmt.Errorf("an API must not be used twice: %s", c.APIKey)
			}
		}

		col, ok := colsMap[c.ID]
		if !ok {
			col = tx.Collection.Create().SaveX(ctx)
			userUp.AddCollections(col)
		}
		colUp := col.Update().SetCollection(c)
		srcsDB, err := col.Sources(ctx)
		if err != nil {
			return err
		}
		srcMap := make(map[int]*ent.SourceSystem, len(srcsDB))
		for _, s := range srcsDB {
			srcMap[s.ID] = s
		}
		for _, s := range c.Edges.Sources {
			src, ok := srcMap[s.ID]
			if !ok {
				src = tx.SourceSystem.Create().SaveX(ctx)
				colUp.AddSources(src)
			}
			if err := src.Update().SetSourceSystem(s).Exec(ctx); err != nil {
				return fmt.Errorf("update source %q of user %q col %q: %w", s.Name, usr.Name, c.Name, err)
			}
		}
		if err := colUp.Exec(ctx); err != nil {
			return fmt.Errorf("update col %q of user %q: %w", c.Name, usr.Name, err)
		}
	}
	if err := userUp.Exec(ctx); err != nil {
		return fmt.Errorf("update user %q: %w", usr.Name, err)
	}

	return nil
}
