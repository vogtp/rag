package usercfg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/model"
	"github.com/vogtp/rag/pkg/types"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
	"github.com/vogtp/rag/pkg/vecDB/history"
	"github.com/vogtp/rag/pkg/web/bearer"
	"gorm.io/gorm"
)

var (
	ErrorEmbedAlreadyRunning     = errors.New("Embedding already running.")
	DefaultCollectionDisplayName = "New Collection (please change)"
)

type Collection struct {
	// gorm.Model in order to avoid name clashes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	owner          string
	UserID         uint   `gorm:"index;column:user_id"`
	Displayname    string `gorm:"not null;column:display_name"`        // DisplayName is the name displayed to the user
	Collectionname string `gorm:"index;unique;column:collection_name"` // CollectionName is the internal unique name of the collection
	APIKey         string `gorm:"index;column:api_key"`

	Genmodel          string        `gorm:"column:gen_model"`
	Embedmodel        string        `gorm:"column:embed_model"`
	DBUpdateIntervall time.Duration `gorm:"column:update_intervall"`
	NextDBUpdate      time.Time     `gorm:"column:update_next"`
	StartDBUpdate     time.Time     `gorm:"column:update_start"`

	Source SourceSystem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	vecDB *vecdb.VecDB
}

var _ (types.Instance) = (*Collection)(nil)

func (c Collection) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("DisplayName", c.Displayname),
		slog.String("CollecitonName", c.Collectionname),
		slog.Int("UserID", int(c.UserID)),
	)
}

func (c *Collection) DisplayName() string {
	return c.Displayname
}

func (c *Collection) CollectionName() string {
	return c.Collectionname
}

func (c *Collection) Owner() string {
	if len(c.owner) > 0 {
		return c.owner
	}
	if dbInstance != nil {
		u, err := gorm.G[User](dbInstance.db).Where("id = ?", c.UserID).First(dbInstance.srvCtx)
		if err != nil {
			dbInstance.slog.Warn("Cannot get user of collection", "col", c, "err", err)
			return c.Collectionname
		}
		c.owner = u.Name
		return c.owner
	}
	return c.Collectionname
}

func (c *Collection) Model(name string) (types.Model, error) {
	if !strings.EqualFold(name, c.Displayname) && !strings.EqualFold(name, c.Collectionname) {
		return nil, fmt.Errorf("model %q not found, collection %q or %q", name, c.Collectionname, c.Displayname)
	}
	return model.NewInstanceModel(c), nil
}

func (c *Collection) Models(ctx context.Context) []types.Model {
	m, _ := c.Model(c.Collectionname)
	return []types.Model{m}
}

func (c *Collection) LLM() string {
	if len(c.Genmodel) < 1 {
		c.Genmodel = viper.GetString(cfg.ModelLLM)
	}
	return c.Genmodel
}

func (c *Collection) ModelEmbedding() string {
	if len(c.Embedmodel) < 1 {
		c.Embedmodel = viper.GetString(cfg.ModelEmbedding)
	}
	return c.Embedmodel
}

func (c *Collection) UpdateIntervall() time.Duration {
	if c.DBUpdateIntervall == 0 {
		c.DBUpdateIntervall = cfg.DefaultVecDBUpdateIntervall
	}
	return c.DBUpdateIntervall
}

func (c *Collection) getVecDb(ctx context.Context, slog *slog.Logger) (*vecdb.VecDB, error) {
	if c.vecDB != nil {
		return c.vecDB, nil
	}
	v, err := vecdb.New(ctx, slog.With("collection", c), c.ModelEmbedding())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to chroma: %w", err)
	}
	c.vecDB = v
	return c.vecDB, nil
}

func (c *Collection) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	slog = slog.With("collection", c)
	vecDB, err := c.getVecDb(ctx, slog)
	if err != nil {
		return nil, fmt.Errorf("cannot create vecDB for %q: %w", c.Collectionname, err)
	}
	res, err := vecDB.Query(ctx, collection, []string{query}, int32(maxResults))
	if err != nil {
		return nil, fmt.Errorf("query vector DB: %w", err)
	}
	return res[0].Documents, nil
}

func (c *Collection) Authorise(w http.ResponseWriter, r *http.Request) bool {
	return bearer.TokenAuth(c.APIKey).Authorise(w, r)
}

func (c *Collection) Embbed(ctx context.Context, slog *slog.Logger, filters ...vecdb.Filter) error {
	if !c.StartDBUpdate.IsZero() && time.Since(c.StartDBUpdate) < c.UpdateIntervall()/3 {
		slog.Info("Embed allready running", "startdate", c.StartDBUpdate)
		return ErrorEmbedAlreadyRunning
	}
	if dbInstance != nil {
		// this should only be nil in  tests
		c.StartDBUpdate = time.Now()
		col := Collection{StartDBUpdate: c.StartDBUpdate}
		if _, err := gorm.G[Collection](dbInstance.db).Where("id = ?", c.ID).Updates(ctx, col); err != nil {
			slog.Error("Cannot set last update on collecion", "err", err)
		}
	}
	slog = slog.With("collection", c)
	start := time.Now()
	docsChan, err := c.GetDocuments(ctx, slog)
	if len(filters) == 0 {
		filters = append(filters, history.New(slog, c.CollectionName(), c.UpdateIntervall()))
	}
	if err != nil {
		return err
	}
	vecDB, err := c.getVecDb(ctx, slog)
	if err != nil {
		return err
	}
	cnt, err := vecDB.Embedd(ctx, slog, c.CollectionName(), docsChan, filters...)
	if err != nil {
		return err
	}
	log := slog.Warn
	if cnt < 10 {
		log = slog.Error
	}
	log("Finished embedding", "doc.count", cnt, "duration", time.Since(start).String())
	if dbInstance != nil {
		// this should only be nil in  tests
		c.NextDBUpdate = time.Now().Add(c.UpdateIntervall())
		c.StartDBUpdate = time.Time{}
		col := Collection{StartDBUpdate: c.StartDBUpdate, NextDBUpdate: c.NextDBUpdate}
		if _, err := gorm.G[Collection](dbInstance.db).Where("id = ?", c.ID).Updates(ctx, col); err != nil {
			slog.Error("Cannot set next update on collecion", "err", err)
		}
	}
	return nil
}

func (c *Collection) GetDocuments(ctx context.Context, slog *slog.Logger) (chan *vecdb.EmbeddDocument, error) {
	slog = slog.With("collection", c)
	confl, err := confluence.New(ctx, slog, *c.Source.confluence())
	if err != nil {
		return nil, fmt.Errorf("create confluence: %w", err)
	}
	return confl.GetDocuments(ctx, slog)
}

func (c *Collection) ListCollections(ctx context.Context, slog *slog.Logger) ([]*chroma.Collection, error) {
	slog = slog.With("collection", c)
	vecDB, err := c.getVecDb(ctx, slog)
	if err != nil {
		return nil, fmt.Errorf("cannot create vecDB for %q: %w", c.Collectionname, err)
	}
	cols, err := vecDB.ListCollections(ctx, c.Collectionname)
	if err != nil {
		return nil, err
	}
	collections := make([]*chroma.Collection, 0, len(cols))
	for _, col := range cols {
		if !strings.EqualFold(col.Name, c.Collectionname) {
			slog.Debug("Not a valid collection", "vecDB.name", col.Name, "collectionName", c.Collectionname)
			continue
		}
		collections = append(collections, col)
	}
	return collections, nil
}
