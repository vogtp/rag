package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
)

type commonData struct {
	Title         string
	Baseurl       string
	Version       string
	StatusMessage string
	Query         string
}

func (srv *Server) common(t string, r *http.Request) *commonData {
	if err := r.ParseForm(); err != nil {
		srv.slog.Warn("Cannot parse form", "err", err)
	}
	q := ""
	if len(r.URL.RawQuery) > 0 {
		q = fmt.Sprintf("?%s", r.URL.RawQuery)
	}
	cd := &commonData{
		Title:   t,
		Baseurl: srv.baseURL,
		Version: cfg.Version,
		Query:   q,
		//Theme:   "light",
	}
	// if theme, err := r.Cookie("theme"); err == nil && theme.Value == "dark" {
	// 	cd.Theme = theme.Value
	// }
	srv.slog.Debug("Prepaired common data", "title", t)
	return cd
}

func (srv *Server) render(w http.ResponseWriter, r *http.Request, templateName string, data any) {
	ah := r.Header.Get("Accept")
	srv.slog.Debug("Render page", "accept_header", ah)

	if strings.Contains(ah, "html") {
		if err := templates.ExecuteTemplate(w, templateName, data); err != nil {
			srv.slog.Error("cannot render template", "template", templateName, "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if strings.Contains(ah, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			srv.slog.Error("cannot encode data to json", "template", templateName, "err", err, logger.Stacktrace())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err := fmt.Errorf("unsupported content-type: %v", ah)
	srv.slog.Warn("Cannot render", "template", templateName, "err", err)
	http.Error(w, err.Error(), http.StatusBadRequest)

}
