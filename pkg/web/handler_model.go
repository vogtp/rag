package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func (srv *Server) modelsHandler(w http.ResponseWriter, r *http.Request) {
	var ret any
	rag:=srv.rag(r)
	if !rag.BearerAuth().Authorise(w, r) {
		slog.Warn("Not authorised", "host", r.Host, "remote", r.RemoteAddr, "URL", r.URL.String())
		return
	}
	if name := r.PathValue("model"); len(name) > 0 {
		rm, err := rag.Model(name)
		if err != nil {
			http.Error(w, fmt.Sprintf("model %s not found", name), http.StatusNotAcceptable)
			return
		}
		ret = rm.ToOpenAI()
	} else {
		mdls := openai.ModelsList{}
		for _, m := range rag.Models(r.Context()) {
			mdls.Models = append(mdls.Models, m.ToOpenAI())
		}
		ret = mdls
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		slog.Warn("cannot encode models json", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("returned /models", "Models", ret)
}
