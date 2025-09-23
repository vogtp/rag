package filesystem_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/filesystem"
)

func TestGenerate(t *testing.T) {
	ctx := context.Background()
	client, err := vecdb.New(ctx, slog.Default(), cfg.DefaultRagCfg(), vecdb.WithChromaAddress("http://localhost:8000"), vecdb.WithOllamaAddress("https://llama-1.its.unibas.ch"))
	if err != nil {
		t.Fatalf("Failed to create vector DB: %v", err)
	}
	dir := "../../../ignore_hr_pdf"
	ls, err := os.ReadDir(dir)
	if err != nil {
		t.Errorf("Cannot read dir %q: %v", dir, err)
	}
	cnt, err := client.Embedd(ctx, "test-Test", filesystem.Generate(ctx, dir))
	if err != nil {
		t.Fatalf("Embedding: %v", err)
	}
	if len(ls) != cnt {
		t.Fatalf("Did not embedd all (%v) docs (%v)", len(ls), cnt)
	}
}
