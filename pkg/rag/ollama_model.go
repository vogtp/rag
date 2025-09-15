package rag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/vogtp/rag/pkg/cfg"
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

// getOllamaClient returns a ollama client
// it is used not only in the OllamaModel
func getOllamaClient(ctx context.Context, llmName string) (*ollama.LLM, error) {
	url := cfg.GetOllamaHost(ctx)
	slog.Info("connecting to ollama", "OllamaModel", llmName, "url", url)
	return ollama.New(
		ollama.WithModel(llmName),
		ollama.WithServerURL(url),
	)
}
