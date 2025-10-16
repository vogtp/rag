package usercfg

import (
	"context"
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
	"github.com/vogtp/rag/pkg/web/bearer"
	"gorm.io/gorm"
)

type Collection struct {
	// gorm.Model in order to avoid name clashes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID        uint   `gorm:"index,column:user_id"`
	Displayname   string `gorm:"column:display_name"`                 // DisplayName is the name displayed to the user
	Collectioname string `gorm:"unique,index,column:collection_name"` // CollectionName is the internal unique name of the collection
	APIKey        string `gorm:"index,column:api_key"`

	Genmodel          string        `gorm:"column:gen_model"`
	Embedmodel        string        `gorm:"column:embed_model"`
	DBUpdateIntervall time.Duration `gorm:"column:update_intervall"`

	Source SourceSystem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	vecDB *vecdb.VecDB
}

var _ (types.Instance) = (*Collection)(nil)
var _ (cfg.RagConfig) = (*Collection)(nil)

func (c *Collection) DisplayName() string {
	return c.Displayname
}

func (c *Collection) CollectionName() string {
	return c.Collectioname
}

func (c *Collection) Model(name string) (types.Model, error) {
	if !strings.EqualFold(name, c.Displayname) && !strings.EqualFold(name, c.Collectioname) {
		return nil, fmt.Errorf("model %q not found, collection %q or %q", name, c.Collectioname, c.Displayname)
	}
	m := model.Ollama{
		Name:    c.LLM(),
		LLMName: c.LLM(),
	}
	return m, nil
}

func (c *Collection) Models(ctx context.Context) []types.Model {
	m, _ := c.Model(c.Collectioname)
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

// FIXME remove
func (c *Collection) VecDBUpdateIntervall() time.Duration {
	return c.UpdateIntervall()
}

func (c *Collection) Confluence() *cfg.ConfluenceCfg {
	src := c.Source
	return &cfg.ConfluenceCfg{
		BaseURL: src.URL,
		Key:     src.Key,
		Spaces:  src.splitParts(),
	}
}

func (c *Collection) getVecDb(ctx context.Context, slog *slog.Logger) (*vecdb.VecDB, error) {
	if c.vecDB != nil {
		return c.vecDB, nil
	}
	v, err := vecdb.New(ctx, slog.With("collection", c.Collectioname), c.ModelEmbedding())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to chroma: %w", err)
	}
	c.vecDB = v
	return c.vecDB, nil
}

func (c *Collection) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	vecDB, err := c.getVecDb(ctx, slog)
	if err != nil {
		return nil, fmt.Errorf("cannot create vecDB for %q: %w", c.Collectioname, err)
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

func (c *Collection) Embbed(ctx context.Context, slog *slog.Logger) error {
	return confluence.Embed(ctx, slog, c)
}

func (c *Collection) ListCollections(ctx context.Context, slog *slog.Logger) ([]*chroma.Collection, error) {
	vecDB, err := c.getVecDb(ctx, slog)
	if err != nil {
		return nil, fmt.Errorf("cannot create vecDB for %q: %w", c.Collectioname, err)
	}
	cols, err := vecDB.ListAllCollections(ctx)
	if err != nil {
		return nil, err
	}
	collections := make([]*chroma.Collection, 0, len(cols))
	for _, col := range cols {
		if !strings.EqualFold(col.Name, c.Collectioname) {
			continue
		}
		collections = append(collections, col)
	}
	return collections, nil
}
