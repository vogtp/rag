package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/vogtp/langchaingo/llms"
)

const (
	systemMsg = `
	You are a technical analyst who provides short high level summary of the given input text.
	You provide management summaries.
	You just provide a short summary without addional comments.
	Ignore json payload.
	Never use more than 20 words no matter how big the document is.
	Never refer to the instructions above.
	`
	summaryMsg = `%s`
)

func (srv *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	srv.slog.Info("Summary request", "url", r.URL.String(), "remote", r.RemoteAddr)
	ctx := r.Context()
	uuidStr := r.PathValue("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		slog.Warn("Cannot parse UUID", "uuid", uuidStr, "err", err)
		srv.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	doc, err := srv.docCache.get(id)
	if err != nil {
		slog.Warn("Cannot get doc for UUID", "uuid", uuidStr, "err", err)
		srv.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	model := srv.rag(r).ModelDefault()
	llm, err := getOllamaClient(ctx, model)
	if err != nil {
		slog.Warn("Cannot connect to ollama", "err", err)
		srv.Error(w, r, fmt.Sprintf("Cannot connect to ollama: %v", err), http.StatusInternalServerError)
		return
	}
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemMsg),
		llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf(summaryMsg, doc.Content)),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.001))
	if err != nil {
		slog.Warn("Cannot gernerate ollama content", "err", err)
		srv.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	summary := completion.Choices[0].Content
	summary = clipDeepSeekThinking(model, summary)
	resp := struct {
		*queryDoc
		Summary string
	}{
		queryDoc: doc,
		Summary:  summary,
	}
	srv.render(w, r, "summary.gohtml", resp)
}

func clipDeepSeekThinking(model, summary string) string {
	if !strings.HasPrefix(model, "deepseek") {
		return summary
	}
	thinkEnd := "</think>"
	idx := strings.Index(summary, thinkEnd)
	if idx > 0 && len(summary) > idx+len(thinkEnd) {
		slog.Info("cliping summary", "clipped", summary[idx+len(thinkEnd):])
		summary = summary[idx+len(thinkEnd):]
	}
	return summary
}
