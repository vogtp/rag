package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/web/bearer"
	"github.com/vogtp/rag/pkg/web/oidc"
)

type Handler interface {
	AllFromRequest(ctx context.Context, r *http.Request) types.Instance
	UserFromRequest(ctx context.Context, r *http.Request) types.Instance
}

var _ (Handler) = (*handler)(nil)

type GetAllRager interface {
	GetAllRags(context.Context) []types.Instance
}

var _ (GetAllRager) = (*handler)(nil)

type handler struct {
	slog       *slog.Logger
	usercfg    *usercfg.DataBase
	globalRags []types.Instance
}

func New(ctx context.Context, slog *slog.Logger, usercfg *usercfg.DataBase) (Handler, error) {
	ragCfgs, err := cfg.GetRagConfig()
	if err != nil {
		return nil, fmt.Errorf("read RAG config: %w", err)
	}
	rags := make([]types.Instance, len(ragCfgs))
	for i, ragCfg := range ragCfgs {
		r, err := newInstanceCfg(ctx, slog, ragCfg)
		if err != nil {
			return nil, fmt.Errorf("start rag %q backend: %w", ragCfg.GetName(), err)
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
func (h handler) GetAllRags(ctx context.Context) []types.Instance {
	rags := make([]types.Instance, 0)
	rags = append(rags, h.globalRags...)
	usrs, err := h.usercfg.Users(ctx)
	if err != nil {
		h.slog.Warn("Cannot query users rags", "err", err)
	}
	rags = append(rags, h.getUserInstances(ctx, usrs...)...)
	return rags
}

func (h handler) AllFromRequest(ctx context.Context, r *http.Request) types.Instance {
	rags := newInstanceList(h.slog, "global")

	rags.Add(h.publicInstances())
	rags.Add(h.UserFromRequest(ctx, r))
	return rags
}

func (h handler) UserFromRequest(ctx context.Context, r *http.Request) types.Instance {
	username, _ := oidc.UserName(r)

	// by Bearer Token
	bt, _ := bearer.Get(r)
	rags := newInstanceList(h.slog, fmt.Sprintf("%s|%s", username, bt))

	if len(username)+len(bt) < 1 {
		h.slog.Warn("Neither user nor token authentication found in request", "request", r.Header)
	}

	if len(username) > 1 {
		if u, err := h.usercfg.User(ctx, username); err == nil {
			rags.Add(h.getUserInstances(ctx, *u)...)
		} else {
			h.slog.Warn("Cannot query user by name", "err", err, "username", username)
		}
	}

	if len(bt) > 1 {
		if usrs, err := h.usercfg.UserByAPIKey(ctx, bt); err == nil {
			rags.Add(h.getUserInstances(ctx, usrs...)...)
		} else {
			h.slog.Warn("Cannot query users by api key", "err", err, "apikey", bt)
		}
		if cols, err := h.usercfg.CollectionByAPIKey(ctx, bt); err == nil {
			rags.Add(h.getCollectionInstances(ctx, cols...)...)
		} else {
			h.slog.Warn("Cannot query collections by api key", "err", err, "apikey", bt)
		}
	}

	return rags
}

func (h handler) getUserInstances(ctx context.Context, usrs ...usercfg.User) []types.Instance {
	rags := make([]types.Instance, 0)
	for _, u := range usrs {
		rags = append(rags, h.getCollectionInstances(ctx, u.Collections...)...)
	}
	return rags
}

func (h handler) getCollectionInstances(ctx context.Context, usrs ...usercfg.Collection) []types.Instance {
	rags := make([]types.Instance, 0)
	for _, c := range usrs {
		idc, err := newInstanceDBCol(ctx, h.slog, &c)
		if err != nil {
			h.slog.Warn("Could not create rag instance from collection %q: %w", c.DisplayName, err)
			continue
		}
		rags = append(rags, idc)
	}
	return rags
}
