package usercfg

import (
	"context"
	"log/slog"
)

func Migrate2Gorm(ctx context.Context, slog *slog.Logger) error {
	dbEnt, err := NewENT(ctx, slog, Dialect, DBFileName)
	if err != nil {
		return err
	}
	usrs, err := dbEnt.GetUserQuery(ctx).All(ctx)
	if err != nil {
		return err
	}
	gDB, err := Create(ctx, slog)
	if err != nil {
		return err
	}

	for _, u := range usrs {
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
				DisplayName:    c.Name,
				CollectionName: c.CollectionName,
				Source:         src,
			}

			gu.Collections[i] = gc
		}
		if err := gDB.Add(ctx, &gu); err != nil {
			return err
		}
	}

	return nil
}
