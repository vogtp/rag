package web

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/vogtp/rag/internal/testhelper"
	"github.com/vogtp/rag/pkg/rag"
)

func TestServer_rag(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		slog *slog.Logger
		// Named input parameters for target function.
		r    *http.Request
		want rag.Instance
	}{
		// TODO: Add test cases.
	}
	db, slog := testhelper.GetDB(t)
	srv, err := New(t.Context(), slog)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := srv.rag(tt.r)
			db.Users(t.Context())
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("rag() = %v, want %v", got, tt.want)
			}
		})
	}
}
