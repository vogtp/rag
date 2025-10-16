package model

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/langchaingo/chains"
	"github.com/vogtp/langchaingo/embeddings"
	"github.com/vogtp/langchaingo/llms"
	"github.com/vogtp/langchaingo/llms/ollama"
	llmopenai "github.com/vogtp/langchaingo/llms/openai"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	"github.com/vogtp/rag/pkg/web/bearer"
)

var _ types.Model = (*Ollama)(nil)

type Ollama struct {
	Name    string
	LLMName string

	OwnedBy string
}

func (m Ollama) GetName() string {
	return m.Name
}

func (m Ollama) String() string {
	return m.GetName()
}

func (m Ollama) GetLLMName() string {
	return m.LLMName
}

func (m Ollama) BearerAuth() bearer.Auth {
	return bearer.NoAuth()
}

func (m Ollama) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.Name,
		Object:  "OllamaModel",
		OwnedBy: m.OwnedBy,
		Parent:  m.LLMName,
	}
}

func (m Ollama) GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc types.StreamingFunc) (string, error) {
	llm, err := GetOllamaClient(ctx, m.LLMName)
	if err != nil {
		return "", fmt.Errorf("cannot get ollama client: %w", err)
	}

	chain := chains.LoadCondenseQuestionGenerator(llm)

	resp, err := chain.LLM.GenerateContent(ctx, messages, llms.WithTemperature(temperature), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		return streamingFunc(ctx, chunk)
	}))

	respString := ""
	if len(resp.Choices) > 0 {
		respString = resp.Choices[0].Content
	}

	return respString, err
}

// BackendModel allows generate and embedding
type BackendModel interface {
	llms.Model
	llms.ReasoningModel
	embeddings.EmbedderClient
}

// GetOllamaClient returns a ollama client
// it is used not only in the OllamaModel
func GetOllamaClient(ctx context.Context, llmName string) (BackendModel, error) {
	url := cfg.GetOllamaHost(ctx)
	slog.Info("connecting to ollama", "LLM", llmName, "url", url)
	return ollama.New(
		ollama.WithModel(llmName),
		ollama.WithServerURL(url),
	)
}

func getOpenAIClient(ctx context.Context, llmName string) (BackendModel, error) {
	url := cfg.GetOllamaHost(ctx)
	slog.Info("connecting to ollama", "LLM", llmName, "url", url)
	return llmopenai.New(
		llmopenai.WithModel(llmName),
		llmopenai.WithBaseURL("FIXME"),
	)
}
