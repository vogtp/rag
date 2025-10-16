package usercfg

import (
	"context"
	"log/slog"
)

func Migrate2Gorm(ctx context.Context, slog *slog.Logger) error {
	dbEnt, err := newENT(ctx, slog, Dialect, DBFileNameEnt)
	if err != nil {
		return err
	}
	usrs, err := dbEnt.GetUserQuery(ctx).All(ctx)
	if err != nil {
		return err
	}
	dbGorm, err := Create(ctx, slog, DBFileName)
	if err != nil {
		return err
	}

	for _, u := range usrs {
		if eu, err := dbGorm.User(ctx, u.Name); err == nil && eu != nil {
			slog.Info("User already migrated", "user", u.Name)
			continue
		}
		slog.Warn("Migration DB from ent to gorm", "user", u.Name)
		gu := User{
			Name:        u.Name,
			APIKey:      u.APIKey,
			Collections: make([]Collection, len(u.Edges.Collections)),
		}

		for i, c := range u.Edges.Collections {
			var src SourceSystem
			if len(c.Edges.Sources) > 0 {
				s := c.Edges.Sources[0]
				src = SourceSystem{
					Name:  s.Name,
					Type:  SourceConfluence,
					URL:   s.URL,
					Key:   s.Key,
					Parts: s.Parts,
				}
			}
			gc := Collection{
				Displayname:   c.Name,
				Collectioname: c.CollectionName,
				Source:        src,
			}

			gu.Collections[i] = gc
		}
		if err := dbGorm.Add(ctx, &gu); err != nil {
			return err
		}
	}

	return nil
}
