package model

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/langchaingo/llms"
	"github.com/vogtp/langchaingo/memory"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/types"
	"github.com/vogtp/rag/pkg/web/bearer"
)

var _ types.Model = (*instanceModel)(nil)

type instanceModel struct {
	backend types.Instance
}

func NewInstanceModel(backend types.Instance) types.Model {
	return &instanceModel{
		backend: backend,
	}
}

func (m instanceModel) GetName() string {
	return m.backend.DisplayName()
}

func (m instanceModel) String() string {
	return m.GetName()
}

func (m instanceModel) GetLLMName() string {
	return m.backend.LLM()
}

func (m instanceModel) BearerAuth() bearer.Auth {
	return m.backend
}

func (m instanceModel) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.backend.CollectionName(),
		Object:  "VectorStoreModel",
		OwnedBy: m.backend.CollectionName(), // FIXME get name from DB
		Parent:  m.backend.LLM(),
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

func (m instanceModel) GenerateContent(ctx context.Context, slog *slog.Logger, messages []llms.MessageContent, temperature float64, streamingFunc types.StreamingFunc) (string, error) {
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
		var err error
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

	//res, err := m.vecDB.Query(ctx, m.Collection, []string{text}, cfg.CntSeachResults)
	res, err := m.backend.SearchVecDB(ctx, slog, m.backend.CollectionName(), text, cfg.CntSeachResults)
	if err != nil {
		return "", fmt.Errorf("query vector DB: %w", err)
	}

	var knowledge strings.Builder

	knowledge.WriteString("Anwser the next question based the the following knowledge:\n")
	for _, d := range res {
		slog.Info("Adding knowledge", "knowledge", d)
		//knowledge.WriteString(fmt.Sprintf("<knowledge href=%q >%s</knowledge>\n", d.URL, d.Document))
		knowledge.WriteString(fmt.Sprintf("<knowledge href=%q >%s</knowledge>\n", d.URL, d.EmbedContent))
	}
	knowledge.WriteString("Always reference the used knowledge by the name and link to the href tag\n")

	last := messages[len(messages)-1]
	messages[len(messages)-1] = llms.MessageContent{
		Role: llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{
			Text: knowledge.String(),
		}},
	}

	messages = append(messages, last)

	llm, err := GetOllamaClient(ctx, slog, m.backend.LLM())
	if err != nil {
		return "", fmt.Errorf("cannot get ollama: %w", err)
	}

	resp, err := llm.GenerateContent(ctx, messages, llms.WithStreamingFunc(streamingFunc))
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Content, nil
}

// func (m *VectorStoreModel) getEmbedder(ctx context.Context, slog *slog.Logger) (*embeddings.EmbedderImpl, error) {
// 	if m.embedder != nil {
// 		return m.embedder, nil
// 	}
// 	llm, err := model.GetOllamaClient(ctx, slog, m.config.ModelEmbedding())
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create llm client: %w", err)
// 	}

// 	e, err := embeddings.NewEmbedder(llm)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create embedder: %w", err)
// 	}
// 	m.embedder = e
// 	return m.embedder, nil
// }
