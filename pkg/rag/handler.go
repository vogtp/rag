package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vogtp/rag/pkg/types"
	"github.com/vogtp/rag/pkg/usercfg"
	"github.com/vogtp/rag/pkg/web/bearer"
	"github.com/vogtp/rag/pkg/web/oidc"
)

type Handler interface {
	FromRequest(ctx context.Context, r *http.Request) types.Instance
}

var _ (Handler) = (*handler)(nil)

type GetAllRager interface {
	GetAllRags(context.Context) []types.Instance
}

type handler struct {
	slog    *slog.Logger
	usercfg *usercfg.DataBase
}

func New(ctx context.Context, slog *slog.Logger, usercfg *usercfg.DataBase) (Handler, error) {
	m := &handler{
		slog:    slog,
		usercfg: usercfg,
	}
	return m, nil
}

func (h handler) FromRequest(ctx context.Context, r *http.Request) types.Instance {
	username, _ := oidc.UserName(r)

	// by Bearer Token
	bt, _ := bearer.Get(r)
	rags := newInstanceList(h.slog, fmt.Sprintf("%s|%s", username, bt))

	if len(username)+len(bt) < 1 {
		h.slog.Warn("Neither user nor token authentication found in request", "request", r.Header, "username", username, "bearer", bt)
	}

	if len(username) > 1 {
		if u, err := h.usercfg.User(ctx, username); err == nil {
			rags.Add(u)
		} else {
			h.slog.Warn("Cannot query user by name", "err", err, "username", username)
		}
	}

	if len(bt) > 1 {
		if usrs, err := h.usercfg.UserByAPIKey(ctx, bt); err == nil {
			for _, u := range usrs {
				rags.Add(&u)
			}
		} else {
			h.slog.Warn("Cannot query users by api key", "err", err, "apikey", bt)
		}
		if cols, err := h.usercfg.CollectionByAPIKey(ctx, bt); err == nil {
			for _, c := range cols {
				rags.Add(&c)
			}
		} else {
			h.slog.Warn("Cannot query collections by api key", "err", err, "apikey", bt)
		}
	}

	return rags
}
