package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/usercfg/db/ent"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/collection"
	"github.com/vogtp/rag/pkg/usercfg/db/ent/user"
	"github.com/vogtp/rag/pkg/web/bearer"
	"github.com/vogtp/rag/pkg/web/oidc"
)

type Handler interface {
	FromRequest(r *http.Request) Instance
}

var _ (Handler) = (*handler)(nil)

type handler struct {
	slog       *slog.Logger
	usercfg    *usercfg.DB
	globalRags []Instance
}

func New(ctx context.Context, slog *slog.Logger, usercfg *usercfg.DB) (Handler, error) {
	ragCfgs, err := cfg.GetRagConfig()
	if err != nil {
		return nil, fmt.Errorf("read RAG config: %w", err)
	}
	rags := make([]Instance, len(ragCfgs))
	for i, ragCfg := range ragCfgs {
		r, err := newInstanceCfg(ctx, slog, ragCfg)
		if err != nil {
			return nil, fmt.Errorf("start rag %q backend: %w", ragCfg.Name, err)
		}
		rags[i] = r
	}
	m := &handler{
		slog:       slog,
		usercfg:    usercfg,
		globalRags: rags,
	}
	return m, nil
}

func (h handler) publicInstances() *instanceList {
	return newInstanceList(h.slog, "public", h.globalRags...)
}

// GetAllRags returns a slice of all instance without any authentical
// for internal use only
func (h handler) GetAllRags(ctx context.Context) []Instance {
	rags := make([]Instance, 0)
	rags = append(rags, h.globalRags...)
	usrs, err := h.usercfg.User.Query().All(ctx)
	if err != nil {
		h.slog.Warn("Cannot query users rags", "err", err)
	}
	rags = append(rags, h.getUserInstances(ctx, usrs...)...)
	return rags
}

func (h handler) FromRequest(r *http.Request) Instance {
	rags := newInstanceList(h.slog, "global")

	rags.Add(h.publicInstances())

	// by authenticated User
	username := oidc.UserName(r)

	// by Bearer Token
	bt, _ := bearer.Get(r)

	if len(username) > 1 {
		if u, err := h.usercfg.ByName(r.Context(), username); err == nil {
			rags.Add(h.getUserInstances(r.Context(), u)...)
		} else {
			h.slog.Warn("Cannot query user by name", "err", err, "username", username)
		}
	}

	if len(bt) > 1 {
		if usrs, err := h.usercfg.GetUserQuery(r.Context()).Where(user.APIKey(bt)).All(r.Context()); err != nil {
			rags.Add(h.getUserInstances(r.Context(), usrs...)...)
		} else {
			h.slog.Warn("Cannot query users by api key", "err", err, "apikey", bt)
		}
		if cols, err := h.usercfg.Collection.Query().Where(collection.APIKey(bt)).All(r.Context()); err != nil {
			rags.Add(h.getCollectionInstances(r.Context(), cols...)...)
		} else {
			h.slog.Warn("Cannot query collections by api key", "err", err, "apikey", bt)
		}
	}

	return rags
}

func (h handler) getUserInstances(ctx context.Context, usrs ...*ent.User) []Instance {
	rags := make([]Instance, 0)
	for _, u := range usrs {
		if u == nil {
			continue
		}
		cols, err := u.Collections(ctx)
		if err != nil {
			h.slog.Warn("Cannot get collections from user", "err", err, "username", u.Name)
		}
		rags = append(rags, h.getCollectionInstances(ctx, cols...)...)
	}
	return rags
}

func (h handler) getCollectionInstances(ctx context.Context, usrs ...*ent.Collection) []Instance {
	rags := make([]Instance, 0)
	for _, c := range usrs {
		idc, err := newInstanceDBCol(ctx, h.slog, c)
		if err != nil {
			h.slog.Warn("Could not create rag instance from collection %q: %w", c.Name, err)
			continue
		}
		rags = append(rags, idc)
	}
	return rags
}
