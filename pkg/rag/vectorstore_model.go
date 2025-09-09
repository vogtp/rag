package rag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/vogtp/rag/pkg/cfg"
)

var _ Model = (*VectorStoreModel)(nil)

type VectorStoreModel struct {
	Name    string
	LLMName string

	OwnedBy string

	Collection  string
	vectorStore vectorstores.VectorStore
	embedder    *embeddings.EmbedderImpl
}

func (m VectorStoreModel) GetName() string {
	return m.Name
}

func (m VectorStoreModel) String() string {
	return m.GetName()
}

func (m VectorStoreModel) GetLLMName() string {
	return m.LLMName
}

func (m VectorStoreModel) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.Name,
		Object:  "VectorStoreModel",
		OwnedBy: m.OwnedBy,
		Parent:  m.LLMName,
	}
}

func (m VectorStoreModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc StreamingFunc) (string, error) {
	store, err := m.getChroma(ctx)
	if err != nil {
		return "", err
	}
	mem := memory.NewConversationBuffer()

	text := ""
	for _, m := range messages {
		slog.Info("Message", "m", m, "type", fmt.Sprintf("%T", m))
		for _, p := range m.Parts {
			slog.Info("Message Part", "part", p)
			if tp, ok := p.(llms.TextContent); ok {
				text = tp.Text
			}
		}
		switch m.Role {
		case llms.ChatMessageTypeAI:
			err = mem.ChatHistory.AddAIMessage(ctx, text)
		case llms.ChatMessageTypeHuman:
			err = mem.ChatHistory.AddUserMessage(ctx, text)
		case llms.ChatMessageTypeSystem:
			err = mem.ChatHistory.AddMessage(ctx, llms.SystemChatMessage{Content: text})
		default:
			err = mem.ChatHistory.AddMessage(ctx, llms.GenericChatMessage{Content: text})
		}
		if err != nil {
			slog.Warn("error adding chat memory", "err", err)
		}
	}
	if len(text) < 1 {
		slog.Warn("No question found", "messages", messages)
		return "", fmt.Errorf("no question found")
	}
	slog.Info("sending final question to vecDB", "question", text)
	if h, err := mem.ChatHistory.Messages(ctx); err == nil {
		slog.Info("Added history", "size", len(h))
	} else {
		slog.Warn("No history", "err", err)
	}

	rec := vectorstores.ToRetriever(
		store,
		7,
		// vectorstores.WithNameSpace(index),
		vectorstores.WithScoreThreshold(0.2),
	)
	llm, err := getOllamaClient(ctx, m.LLMName)
	if err != nil {
		return "", fmt.Errorf("cannot get ollama: %w", err)
	}
	c := chains.NewConversationalRetrievalQAFromLLM(llm, rec, mem)
	// input["question"] = text
	// r, err := chains.Call(ctx, c, input, chains.WithStreamingFunc(streamingFunc))
	// if err != nil {
	// 	return "", fmt.Errorf("chains.chall error: %w", err)
	// }
	// for k, v := range r {
	// 	slog.Info("Call response", "k", k, "v", v)
	// }
	// return "", nil
	return chains.Run(ctx, c, text, chains.WithStreamingFunc(streamingFunc))
}

func (m *VectorStoreModel) getChroma(ctx context.Context) (vectorstores.VectorStore, error) {
	if m.vectorStore != nil {
		return m.vectorStore, nil
	}
	e, err := m.getEmbedder(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create embedder: %w", err)
	}
	store, err := chroma.New(
		chroma.WithChromaURL(viper.GetString(cfg.ChromaUrl)),
		chroma.WithNameSpace(m.Collection),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction(types.COSINE),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create chroma client: %w", err)
	}
	return &store, nil
}

func (m *VectorStoreModel) getEmbedder(ctx context.Context) (*embeddings.EmbedderImpl, error) {
	if m.embedder != nil {
		return m.embedder, nil
	}
	model := viper.GetString(cfg.ModelEmbedding)
	llm, err := getOllamaClient(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("cannot create llm client: %w", err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("cannot create embedder: %w", err)
	}
	m.embedder = e
	return m.embedder, nil
}
