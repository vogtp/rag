package noopfilter

import (
	"github.com/vogtp/rag/pkg/types"
)

type noopfilter struct{}

// New creates a filter allowing all documents
// used to disable the default history filter
func New() types.Filter {
	return &noopfilter{}
}

func (eh *noopfilter) ShouldEmbedd(_ *types.EmbeddDocument) bool {
	return true
}

func (eh *noopfilter) ReqisterEmedded(d *types.EmbeddDocument) {
	// satisfy the interface
}
