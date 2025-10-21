package noopfilter

import (
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type noopfilter struct{}

// New creates a filter allowing all documents
// used to disable the default history filter
func New() vecdb.Filter {
	return &noopfilter{}
}

func (eh *noopfilter) ShouldEmbedd(_ *vecdb.EmbeddDocument) bool {
	return true
}

func (eh *noopfilter) ReqisterEmedded(d *vecdb.EmbeddDocument) {
	// satisfy the interface
}
