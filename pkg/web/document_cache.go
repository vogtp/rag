package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/rag/pkg/types"
)

type queryDocument struct {
	*types.QueryDocument
	UUID   uuid.UUID
	access time.Time
}

type docCache struct {
	mu    sync.RWMutex
	cache map[uuid.UUID]queryDocument
}

func newDocCache() docCache {
	return docCache{
		cache: make(map[uuid.UUID]queryDocument),
	}
}

func (dc *docCache) add(d *queryDocument) {
	d.access = time.Now()
	dc.mu.RLock()
	dc.cache[d.UUID] = *d
	dc.mu.RUnlock()
}

func (dc *docCache) get(id uuid.UUID) (*queryDocument, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	d, ok := dc.cache[id]
	if !ok {
		return nil, fmt.Errorf("cannot find document for %v", id)
	}
	d.access = time.Now()
	dc.cache[id] = d
	go dc.cleanup()
	return &d, nil
}

func (dc *docCache) cleanup() {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	for uuid, d := range dc.cache {
		if time.Since(d.access) < 10*time.Minute {
			continue
		}
		delete(dc.cache, uuid)
	}
}
