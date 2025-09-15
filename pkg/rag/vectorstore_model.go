package rag

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
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

var firstInstruction = llms.MessageContent{
	Role: llms.ChatMessageTypeSystem,
	Parts: []llms.ContentPart{llms.TextContent{
		Text: `You are the friendly assitant of the university of Basel in Switzerland.
Answer short an precise based on the provided knowledge.
Always give references to the used knowledge and if you cannot say "I do not know"`,
	}},
}

func (m VectorStoreModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc StreamingFunc) (string, error) {
	store, err := m.getChroma(ctx)
	if err != nil {
		return "", err
	}
	_ = store
	mem := memory.NewConversationBuffer()

	if !reflect.DeepEqual(messages[0], firstInstruction) {
		messages = append([]llms.MessageContent{firstInstruction}, messages...)
	}

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
	client, err := vecdb.New(ctx, logger.New(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return "", fmt.Errorf("create vector DB: %w", err)
	}
	res, err := client.Query(ctx, m.Collection, []string{text}, 5)
	if err != nil {
		return "", fmt.Errorf("query vector DB: %w", err)
	}
	for i, r := range res[0].Documents {
		fmt.Printf("\n\nDocument %v: %+v\n", i, r)
	}

	var knowledge strings.Builder

	knowledge.WriteString("Anwser the next question based the the following knowledge:\n")
	for _, d := range res[0].Documents {
		knowledge.WriteString(fmt.Sprintf("<knowledge reference=%q >%s</knowledge>\n", d.URL, d.Content))
	}

	last := messages[len(messages)-1]
	messages[len(messages)-1] = llms.MessageContent{
		Role: llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{
			Text: knowledge.String(),
		}},
	}

	messages = append(messages, last)
	// rec := vectorstores.ToRetriever(
	// 	store,
	// 	7,
	// 	// vectorstores.WithNameSpace(index),
	// 	//vectorstores.WithScoreThreshold(0.2),
	// )
	// recs, err := rec.GetRelevantDocuments(ctx, text)
	// if err != nil {
	// 	slog.Warn("cannot get relevant docs", "err", err)
	// }
	// for _, r := range recs {
	// 	slog.Warn(r.PageContent)
	// }
	llm, err := getOllamaClient(ctx, m.LLMName)
	if err != nil {
		return "", fmt.Errorf("cannot get ollama: %w", err)
	}

	resp, err := llm.GenerateContent(ctx, messages, llms.WithStreamingFunc(streamingFunc))
	return resp.Choices[0].Content, err
	// // chains.NewLLMChain(llm, prompts.NewChatPromptTemplate())
	// c := chains.NewConversationalRetrievalQAFromLLM(llm, rec, mem)
	// // input["question"] = text
	// // r, err := chains.Call(ctx, c, input, chains.WithStreamingFunc(streamingFunc))
	// // if err != nil {
	// // 	return "", fmt.Errorf("chains.chall error: %w", err)
	// // }
	// // for k, v := range r {
	// // 	slog.Info("Call response", "k", k, "v", v)
	// // }
	// // return "", nil
	// return chains.Run(ctx, c, text, chains.WithStreamingFunc(streamingFunc))
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
		chroma.WithChromaURL(cfg.ChromaUrl()),
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
