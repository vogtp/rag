package rag

import (
	"github.com/sashabaranov/go-openai"
	"github.com/vogtp/langchaingo/llms"
)

func RoleOpenAI2langchain(role string) llms.ChatMessageType {
	switch role {
	case openai.ChatMessageRoleUser:
		return llms.ChatMessageTypeHuman
	case openai.ChatMessageRoleAssistant:
		return llms.ChatMessageTypeAI
	case openai.ChatMessageRoleSystem:
		return llms.ChatMessageTypeSystem
	case openai.ChatMessageRoleFunction:
		return llms.ChatMessageTypeFunction
	case openai.ChatMessageRoleTool:
		return llms.ChatMessageTypeTool
	default:
		return llms.ChatMessageTypeHuman
	}
}
