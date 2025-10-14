package rag

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
	"github.com/vogtp/rag/pkg/web/bearer"
)

var _ Model = (*OllamaModel)(nil)

type OllamaModel struct {
	Name    string
	LLMName string

	OwnedBy string
}

func (m OllamaModel) GetName() string {
	return m.Name
}

func (m OllamaModel) String() string {
	return m.GetName()
}

func (m OllamaModel) GetLLMName() string {
	return m.LLMName
}

func (m OllamaModel) BearerAuth() bearer.Auth {
	return bearer.NoAuth()
}

func (m OllamaModel) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.Name,
		Object:  "OllamaModel",
		OwnedBy: m.OwnedBy,
		Parent:  m.LLMName,
	}
}

func (m OllamaModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc StreamingFunc) (string, error) {
	llm, err := getOllamaClient(ctx, m.LLMName)
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

// getOllamaClient returns a ollama client
// it is used not only in the OllamaModel
func getOllamaClient(ctx context.Context, llmName string) (BackendModel, error) {
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
