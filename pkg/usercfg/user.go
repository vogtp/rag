package usercfg

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model in order to avoid name clashes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string       `gorm:"index;unique;column:name"`
	APIKey      string       `gorm:"index;column:api_key"`
	Collections []Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

var _ (types.Instance) = (*User)(nil)
var _ (cfg.RagConfig) = (*User)(nil)

// Collection returns a collection by CollectionName or DisplayName
func (u *User) Collection(n string) *Collection {
	for _, c := range u.Collections {
		if c.Collectionname == n {
			return &c
		}
		if c.Displayname == n {
			return &c
		}
	}
	return nil
}

func (u *User) DisplayName() string {
	return u.Name
}

func (u *User) CollectionName() string {
	return u.Name
}

func (u *User) Model(name string) (m types.Model, err error) {
	var retErr error
	for _, c := range u.Collections {
		m, err = c.Model(name)
		if err == nil {
			return m, nil
		}
		retErr = fmt.Errorf("%v\n%w", retErr, err)
	}
	return nil, retErr
}

func (u *User) Models(ctx context.Context) []types.Model {
	m := make([]types.Model, 0)
	for _, c := range u.Collections {
		m = append(m, c.Models(ctx)...)
	}
	return m
}

func (u *User) LLM() string {
	for _, c := range u.Collections {
		llm := c.LLM()
		if len(llm) > 0 {
			return llm
		}
	}
	return viper.GetString(cfg.ModelLLM)
}

func (u *User) ModelEmbedding() string {
	for _, c := range u.Collections {
		llm := c.ModelEmbedding()
		if len(llm) > 0 {
			return llm
		}
	}
	return viper.GetString(cfg.ModelEmbedding)
}

func (u *User) UpdateIntervall() time.Duration {
	d := cfg.DefaultVecDBUpdateIntervall
	for _, c := range u.Collections {
		intervall := c.UpdateIntervall()
		if intervall < cfg.MinVecDBUpdateIntervall {
			continue
		}
		d = min(d, intervall)
	}
	return d
}

// FIXME remove
func (u *User) VecDBUpdateIntervall() time.Duration {
	return u.UpdateIntervall()
}

func (u *User) Confluence() *cfg.ConfluenceCfg {
	return &cfg.ConfluenceCfg{}
}

func (u *User) Embbed(ctx context.Context, slog *slog.Logger) error {
	slog = slog.With("username", u.Name)
	for _, c := range u.Collections {
		go func(c *Collection) {
			if err := c.Embbed(ctx, slog); err != nil {
				slog.Warn("Cannot embed user collection", "err", err, "collection", c.DisplayName())
			}
		}(&c)
	}
	return nil
}

func (u *User) ListCollections(ctx context.Context, slog *slog.Logger) ([]*chroma.Collection, error) {
	slog = slog.With("username", u.Name)
	cols := make([]*chroma.Collection, 0)
	for _, c := range u.Collections {
		c, err := c.ListCollections(ctx, slog)
		if err != nil {
			slog.Warn("Cannot list collections of user", "err", err)
		}
		cols = append(cols, c...)
	}
	// if len(cols) < 1 {
	// 	return cols, fmt.Errorf("no collections found")
	// }
	return cols, nil
}

func (u *User) SearchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	slog = slog.With("username", u.Name)
	docs := make([]vecdb.QueryDocument, 0)
	for _, c := range u.Collections {
		c, err := c.SearchVecDB(ctx, slog, collection, query, maxResults)
		if err != nil {
			slog.Error("Cannot search vecDB of user", "err", err, "collection", collection)
		}
		docs = append(docs, c...)
	}
	if len(docs) < 1 {
		return docs, fmt.Errorf("no documents found")
	}
	return docs, nil
}

func (u *User) Authorise(w http.ResponseWriter, r *http.Request) bool {
	for _, c := range u.Collections {
		if c.Authorise(w, r) {
			return true
		}
	}
	return false
}
