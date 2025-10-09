package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type queryDoc struct {
	*vecdb.QueryDocument
	UUID   uuid.UUID
	access time.Time
}

type docChace struct {
	mu    sync.RWMutex
	cache map[uuid.UUID]queryDoc
}

func newDocCache() docChace {
	return docChace{
		cache: make(map[uuid.UUID]queryDoc),
	}
}

func (dc *docChace) add(d *queryDoc) {
	d.access = time.Now()
	dc.mu.RLock()
	dc.cache[d.UUID] = *d
	dc.mu.RUnlock()
}

func (dc *docChace) get(id uuid.UUID) (*queryDoc, error) {
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

func (dc *docChace) cleanup() {
	for uuid, d := range dc.cache {
		if time.Since(d.access) < 10*time.Minute {
			continue
		}
		dc.mu.Lock()
		delete(dc.cache, uuid)
		dc.mu.Unlock()
	}
}
