package chroma_test

import (
	"strings"
	"testing"

	"github.com/vogtp/rag/pkg/vecDB/chroma"
)

func TestFixCollectionName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "test-test", want: "test-test"},
		{name: "test-test1a", want: "test-test1a"},
		{name: "test-test1", want: "test-test"},
		{name: "test-test (gone)", want: "test-test__gone"},
		{name: "test-test..test", want: "test-test__test"},
		{name: "test-test...test", want: "test-test___test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chroma.FixCollectionName(tt.name)
			// TODO: update the condition below to compare got with tt.want.
			if !strings.EqualFold(got, tt.want) {
				t.Errorf("FixCollectionName() = %q, want %q", got, tt.want)
			}
		})
	}
}
