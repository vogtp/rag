package usercfg

//go:generate stringer -type=SourceSystemType --trimprefix Source

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	slogGorm "github.com/orandin/slog-gorm"
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

type SourceSystemType int

const (
	SourceConfluence SourceSystemType = iota
	SourceHTTP
)

type SourceSystem struct {
	// gorm.Model in order to avoid name clashes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CollectionID uint             `gorm:"column:collection_id"`
	Name         string           `gorm:"column:name"`
	Type         SourceSystemType `gorm:"column:src_sys_type"`
	URL          string           `gorm:"column:url"`
	Key          string           `gorm:"column:key"`
	Parts        string           `gorm:"column:parts"`
}

func (s SourceSystem) splitParts() []string {
	parts := strings.Split(s.Parts, ",")
	if len(parts) < 1 {
		parts = strings.Split(s.Parts, " ")
	}
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

type DataBase struct {
	db   *gorm.DB
	slog *slog.Logger
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
	return &DataBase{db: db, slog: sl}, nil
}

func (d *DataBase) Add(ctx context.Context, u *User) error {
	// get primary keys from DB
	if eu, err := d.User(ctx, u.Name); err == nil && eu != nil {
		u.ID = eu.ID
		for i, c := range u.Collections {
			if len(c.CollectionName) < len(c.DisplayName) {
				c.CollectionName = fmt.Sprintf("%s-%s", u.Name, chroma.FixCollectionName(c.DisplayName))
			}
			if ec := eu.Collection(c.CollectionName); ec != nil {
				c.ID = ec.ID
				c.Source.ID = ec.Source.ID
				u.Collections[i] = c
			}
		}
	}
	if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(u).Error; err != nil {
		return fmt.Errorf("adding user %q: %w", u.Name, err)
	}
	for _, c := range u.Collections {
		if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c).Error; err != nil {
			return fmt.Errorf("adding user %q collection %q: %w", u.Name, c.DisplayName, err)
		}
		if err := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c.Source).Error; err != nil {
			return fmt.Errorf("adding user %q collection %q source %q: %w", u.Name, c.DisplayName, c.Source.Name, err)
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
	return &u, err
}

func (d *DataBase) UserByAPIKey(ctx context.Context, key string) ([]User, error) {
	usrs, err := d.usr().Where("api_key = ?", key).Find(ctx)
	return usrs, err
}

func (d *DataBase) CollectionByAPIKey(ctx context.Context, key string) ([]Collection, error) {
	cols, err := d.col().Where("api_key = ?", key).Find(ctx)
	return cols, err
}
