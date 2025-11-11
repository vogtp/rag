package usercfg

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/vecDB/chroma"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	//Dialect    = dialect.SQLite
	DBFileName = "rag_users.sqlite"
)

type DataBase struct {
	db     *gorm.DB
	srvCtx context.Context
	slog   *slog.Logger
}

var dbInstance *DataBase

func Create(ctx context.Context, sl *slog.Logger, name string) (*DataBase, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	sl = sl.With(slog.String("database", name))
	logOpts := []slogGorm.Option{slogGorm.WithHandler(sl.Handler())}
	if sl.Enabled(ctx, slog.LevelDebug) {
		logOpts = append(logOpts, slogGorm.WithTraceAll()) // trace all messages
		logOpts = append(logOpts, slogGorm.SetLogLevel(slogGorm.DefaultLogType, logger.Level()))
	}
	gormSlog := slogGorm.New(logOpts...)
	backend := sqlite.Open(fmt.Sprintf("file:%s?&cache=shared&_fk=1", name))
	db, err := gorm.Open(backend, &gorm.Config{Logger: gormSlog})
	if err != nil {
		return nil, fmt.Errorf("create db: %w", err)
	}
	if err := db.AutoMigrate(&User{}, &Collection{}, &SourceSystem{}); err != nil {
		return nil, fmt.Errorf("automigration DB: %w", err)
	}
	dbInstance = &DataBase{db: db, slog: sl, srvCtx: ctx}
	return dbInstance, nil
}

func (d *DataBase) Add(ctx context.Context, u *User) error {
	for i, c := range u.Collections {
		u.Collections[i].Collectionname = fmt.Sprintf("%s-%s", u.Name, chroma.FixCollectionName(c.Displayname))
	}
	// get primary keys from DB
	if eu, err := d.User(ctx, u.Name); err == nil && eu != nil {
		u.ID = eu.ID
		for i, c := range u.Collections {
			c.UserID = u.ID
			if ec := eu.Collection(c.Collectionname); ec != nil {
				c.ID = ec.ID
				c.Source.ID = ec.Source.ID
				if c.Source.Parts != ec.Source.Parts {
					c.NextDBUpdate = time.Now()
				}
			}
			c.Source.CollectionID = c.ID
			u.Collections[i] = c
		}
	}
	if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(u).Error; err != nil {
		return fmt.Errorf("adding user %q: %w", u.Name, err)
	}
	for _, c := range u.Collections {
		if strings.EqualFold(c.Displayname, DefaultCollectionDisplayName) {
			continue
		}
		if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c).Error; err != nil {
			return fmt.Errorf("adding user %q collection %q: %w", u.Name, c.Displayname, err)
		}
		if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c.Source).Error; err != nil {
			return fmt.Errorf("adding user %q collection %q source %q: %w", u.Name, c.Displayname, c.Source.Name, err)
		}
	}
	if err := d.CleanupUserCollections(ctx, u); err != nil {
		d.slog.Warn("Cannot cleanup user collections", "username", u.Name, "err", err)
	}
	return nil
}

var nilQuery = func(db gorm.PreloadBuilder) error { return nil }

func (d *DataBase) usr() gorm.ChainInterface[User] {
	return gorm.G[User](d.db).Preload("Collections.Source", nilQuery)
}

func (d *DataBase) col() gorm.ChainInterface[Collection] {
	return gorm.G[Collection](d.db).Preload("Source", nilQuery)
}

func (d *DataBase) Users(ctx context.Context) ([]User, error) {
	return d.usr().Find(ctx)
}

func (d *DataBase) User(ctx context.Context, name string) (*User, error) {
	u, err := d.usr().Where("name = ?", name).First(ctx)
	for i,c:=range u.Collections{
		if c.DBUpdateIntervall < time.Hour {
			c.DBUpdateIntervall = viper.GetDuration(cfg.VecDBUpdateIntervall)
		}
		if len(c.Genmodel) < 1 {
			c.Genmodel = viper.GetString(cfg.ModelLLM)
		}
		if len(c.Embedmodel) < 1 {
			c.Embedmodel = viper.GetString(cfg.ModelEmbedding)
		}
		u.Collections[i] = c
	}
	return &u, err
}

func (d *DataBase) UserByAPIKey(ctx context.Context, key string) ([]User, error) {
	return d.usr().Where("api_key = ?", key).Find(ctx)
}

func (d *DataBase) CollectionByAPIKey(ctx context.Context, key string) ([]Collection, error) {
	return d.col().Where("api_key = ?", key).Find(ctx)
}

func (d *DataBase) CollectionsToUpdate(ctx context.Context, t time.Time) ([]Collection, error) {
	return d.col().Where("update_next < ?", t.Format(time.DateTime)).Order("update_next asc").Find(ctx)
}
