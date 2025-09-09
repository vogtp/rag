package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/vogtp/rag/pkg/rag"
)

func (srv *Server) completionHandler(w http.ResponseWriter, r *http.Request) {
	var req openai.CompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog := slog.With("model", req.Model)
	slog.Info("Completition Request")
	// model, err := a.rag.Model(req.Model)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	// if req.Stream {
	// 	a.handleCompletionStream(&req, model, w, r)
	// 	return
	// }
	// rag.handleCompletion(&req, ragModel, w, r)
}

func (srv *Server) chatCompletionHandler(w http.ResponseWriter, r *http.Request) {
	var req openai.ChatCompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog := slog.With("model", req.Model)
	slog.Info("Completition Request")
	ragModel, err := srv.rag.Model(req.Model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if req.Stream {
		srv.handleCompletionStream(&req, ragModel, w, r)
		return
	}
	//a.handleChatCompletion(&req, ragModel, w, r)
}

func (srv *Server) handleCompletionStream(req *openai.ChatCompletionRequest, ragModel rag.Model, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	msgs := make([]llms.MessageContent, 0, len(req.Messages)*3)
	for i, m := range req.Messages {
		srv.slog.Info("Chat message", "role", m.Role, "content", m.Content, "idx", i)
		role := rag.RoleOpenAI2langchain(m.Role)
		// if len(ragModel.Collection) > 0 && role == llms.ChatMessageTypeHuman {
		// 	docs, err := getDocs(ctx, ragModel.Collection, m.Content)
		// 	if err == nil {
		// 		for _, doc := range docs {
		// 			msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeSystem, doc.PageContent))
		// 			a.slog.Info("Added doc", "doc_start", doc.PageContent[:60], "chat_sequence", i, "collection", ragModel.Collection)
		// 		}
		// 	} else {
		// 		slog.Warn("Cannot query docs", "err", err)
		// 	}
		// }
		msgs = append(msgs, llms.TextParts(role, m.Content))
	}

	resChan := make(chan []byte, 5)
	go func() {
		defer close(resChan)

		resp, err := ragModel.GenerateContent(ctx, msgs, 0.001, func(ctx context.Context, chunk []byte) error {
			if ctx.Err() != nil {
				slog.Error("GenerateContent call with a canceled context", "context error", ctx.Err())
				return ctx.Err()
			}
			resChan <- chunk
			slog.Debug("stream response", "chunk", string(chunk))
			return nil
		})
		if err != nil {
			slog.Error("llm backend error", "err", err)
			http.Error(w, fmt.Sprintf("llm backend error: %v", err), http.StatusInternalServerError)
			return
		}
		slog.Debug("Generate content finished", "resp", resp)
	}()

	srv.setStreamHeaders(w)
	stream(ctx, w, func(w io.Writer) bool {
		data := []byte("data: ")
		// chunk data
		if chunk, ok := <-resChan; ok {
			// chunk, err := json.Marshal(res)
			if chunk == nil {
				srv.slog.Warn("Stream error data is nil")
				if _, err := w.Write([]byte("data: [ERROR]\n\n")); err != nil {
					slog.Warn("Cannot write streaming bytes", "err", err)
					return false
				}
				return false
			}

			res := generateChatStreamResponse(ragModel, chunk)
			paypload, err := json.Marshal(res)
			if err != nil {
				if _, err := w.Write([]byte("data: [ERROR]\n\n")); err != nil {
					slog.Warn("Cannot write streaming ERROR", "err", err)
					return false
				}
				return false
			}
			// write
			srv.slog.Debug("http stram", "chunk", chunk, "out", res.Choices[0].Delta.Content)
			data = append(data, paypload...)
			data = append(data, []byte("\n\n")...)
			_, err = w.Write(data)
			if err != nil {
				if _, err := w.Write([]byte("data: [ERROR]\n\n")); err != nil {
					slog.Warn("Cannot write streaming ERROR", "err", err)
					return false
				}
				return false
			}
			return true
		}
		// done
		if _, err := w.Write([]byte("data: [DONE]\n\n")); err != nil {
			slog.Warn("Cannot write streaming DONE", "err", err)
			return false
		}
		srv.slog.Debug("Finished streaming")
		return false
	})
}

func stream(ctx context.Context, w http.ResponseWriter, step func(w io.Writer) bool) bool {
	for {
		select {
		case <-ctx.Done():
			return true
		default:
			keepOpen := step(w)
			w.(http.Flusher).Flush()
			if !keepOpen {
				return false
			}
		}
	}
}

func generateChatStreamResponse(ragModel rag.Model, chunk []byte) *openai.ChatCompletionStreamResponse {
	id := prefixID("chatcmpl-")
	res := openai.ChatCompletionStreamResponse{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   ragModel.GetLLMName(),
	}
	choice := openai.ChatCompletionStreamChoice{
		Delta: openai.ChatCompletionStreamChoiceDelta{
			Content: string(chunk),
			Role:    openai.ChatMessageRoleAssistant,
		},
	}
	if len(chunk) < 1 {
		choice.Delta = openai.ChatCompletionStreamChoiceDelta{}
		choice.FinishReason = openai.FinishReasonStop
		slog.Info("ollama DONE")
	}
	res.Choices = append(res.Choices, choice)
	return &res
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func prefixID(prefix string, length ...int) string {
	l := 29
	if len(length) > 0 {
		l = length[0]
	}
	id, err := gonanoid.Generate(alphabet, l)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s", prefix, id)
}
