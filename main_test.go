package main

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/model"
)

func Test_main(t *testing.T) {
	cfg.Parse()
	m, err := model.GetBackendModel(t.Context(), slog.Default(), "hosted_vllm/qwen3-embedding-0.6b")
	if err != nil {
		t.Fatalf("Loading model: %v", err)
	}
	e, err := m.CreateEmbedding(t.Context(), []string{"test"})
	if err != nil {
		t.Errorf("emedding: %v", err)
	}
	fmt.Printf("Embedding: %+v\n", e)
}
