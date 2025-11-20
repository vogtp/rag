package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	sl "log/slog"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/langchaingo/llms"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/types"
)

func (srv *Server) completionHandler(w http.ResponseWriter, r *http.Request) {
	var req openai.CompletionRequest
	srv.slog.Info("completion request", "url", r.URL.String(), "remote", r.RemoteAddr)
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
	ctx := r.Context()
	slog := srv.slog.With("url", r.URL.String(), "remote", r.RemoteAddr)
	slog.Info("chatCompletion request")
	var req openai.ChatCompletionRequest
	start := time.Now()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog = slog.With("model", req.Model, "stream", req.Stream)
	ragModel, err := srv.rag(r).Model(req.Model)
	if err != nil {
		slog.Warn("Cannot get model", "err", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !ragModel.BearerAuth().Authorise(w, r) {
		slog.Warn("Not authorised", "host", r.Host, "remote", r.RemoteAddr, "URL", r.URL.String())
		return
	}
	defer func(t time.Time) {
		slog.Info("Chat completion request finished", sl.GroupAttrs("duration", sl.String("human", time.Since(t).String()), sl.Int64("millies", time.Since(t).Milliseconds()), sl.Duration("duration", time.Since(t))))
	}(start)
	if req.Stream {
		srv.handleCompletionStream(&req, ragModel, w, r)
		return
	}
	msgs := make([]llms.MessageContent, 0, len(req.Messages)*3)
	//choices := make([]openai.ChatCompletionChoice, len(req.Messages))
	for i, m := range req.Messages {
		srv.slog.Info("Chat message", "role", m.Role, "content", m.Content, "idx", i)
		role := rag.RoleOpenAI2langchain(m.Role)
		msgs = append(msgs, llms.TextParts(role, m.Content))
		// choices[i] = openai.ChatCompletionChoice{
		// 	Index: i,
		// 	Message: openai.ChatCompletionMessage{
		// 		Role:         m.Role,
		// 		Content:      m.Content,
		// 		Refusal:      m.Refusal,
		// 		MultiContent: m.MultiContent,
		// 		Name:         m.Name,
		// 		FunctionCall: m.FunctionCall,
		// 		ToolCalls:    m.ToolCalls,
		// 		ToolCallID:   m.ToolCallID,
		// 	},
		// }
	}
	content, err := ragModel.GenerateContent(ctx, slog, msgs, 0.001, func(ctx context.Context, chunk []byte) error { return nil })
	if err != nil {
		slog.Warn("Internal server error: Cannot generate content", "err", err, "ragModel", ragModel)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// choices = append(choices, openai.ChatCompletionChoice{
	// 	Index: len(choices),
	// 	Message: openai.ChatCompletionMessage{
	// 		Role:    openai.ChatMessageRoleAssistant,
	// 		Content: content,
	// 	},
	// 	FinishReason: openai.FinishReasonStop,
	// })
	id := prefixID("chatcmpl-")
	resp := openai.ChatCompletionResponse{
		ID:      id,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []openai.ChatCompletionChoice{{
			Index: 0,
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			},
			FinishReason: openai.FinishReasonStop,
		}},
	}

	//slog.Debug("Answer", "content", content, "response", resp)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Warn("Internal server error: encode chat response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (srv *Server) handleCompletionStream(req *openai.ChatCompletionRequest, ragModel types.Model, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog := srv.slog.With("request", r.URL.String())
	msgs := make([]llms.MessageContent, 0, len(req.Messages)*3)
	for i, m := range req.Messages {
		srv.slog.Info("Chat message", "role", m.Role, "content", m.Content, "idx", i)
		role := rag.RoleOpenAI2langchain(m.Role)
		msgs = append(msgs, llms.TextParts(role, m.Content))
	}

	resChan := make(chan []byte, 5)
	go func() {
		defer close(resChan)

		resp, err := ragModel.GenerateContent(ctx, srv.slog, msgs, 0.001, func(ctx context.Context, chunk []byte) error {
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
	stream(ctx, slog, w, func(w io.Writer) bool {
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
				slog.Error("Streaming function cannot decode json", "err", err, "chunk", string(chunk))
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
				slog.Error("Streaming function cannot write chunck to response", "err", err, "chunk", string(chunk))
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

func stream(ctx context.Context, slog *slog.Logger, w http.ResponseWriter, step func(w io.Writer) bool) bool {
	for {
		select {
		case <-ctx.Done():
			slog.Warn("stream func cancled by context", "err", ctx.Err().Error())
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

func generateChatStreamResponse(ragModel types.Model, chunk []byte) *openai.ChatCompletionStreamResponse {
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
		//slog.Info("ollama DONE")
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
