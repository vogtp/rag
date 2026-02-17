package history

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/vogtp/rag/pkg/types"
)

const (
	HistoryCacheDirName = ".embedd_history"
)

type emeddHistory struct {
	slog            *slog.Logger
	collectionName  string
	history         map[string]time.Time
	minAge          time.Duration
	updateIntervall time.Duration
}

func New(slog *slog.Logger, collectionName string, updateIntervall time.Duration) types.Filter {
	return &emeddHistory{
		slog:            slog,
		collectionName:  collectionName,
		updateIntervall: updateIntervall,
	}
}

func (eh *emeddHistory) ShouldEmbedd(d *types.EmbeddDocument) bool {
	if err := eh.init(); err != nil {
		slog.Warn("Error init the embedd history", "err", err)
		return true
	}
	h, ok := eh.history[eh.key(d)]
	if !ok {
		return true
	}
	last := time.Since(h)
	b := last > eh.minAge
	eh.slog.Debug("should embedd", "should", b, "age", last.String(), "last", h.String())
	return b
}

func (eh *emeddHistory) ReqisterEmedded(d *types.EmbeddDocument) {
	if err := eh.init(); err != nil {
		slog.Warn("Error init the embedd history", "err", err)
		return
	}
	eh.history[eh.key(d)] = time.Now()
	if err := eh.save(); err != nil {
		slog.Warn("Cannot save embedd history", "err", err)
	}
}

func (eh *emeddHistory) init() error {
	if eh.history != nil {
		return nil
	}
	if eh.slog == nil {
		eh.slog = slog.Default()
	}
	if err := os.MkdirAll(HistoryCacheDirName, os.ModePerm); err != nil {
		slog.Error("cannot create history cache dir", "err", err, "HistoryCacheDirName", HistoryCacheDirName)
	}
	interval := eh.updateIntervall
	if interval > 24*time.Hour || interval < time.Hour {
		interval = time.Hour
	}
	eh.minAge = interval
	eh.history = make(map[string]time.Time)
	return eh.load()
}

func (eh emeddHistory) key(d *types.EmbeddDocument) string {
	return fmt.Sprintf("%s.%s", d.IDMetaKey, d.IDMetaValue)
}

func (eh emeddHistory) filename() string {
	return fmt.Sprintf("%s/%s.json", HistoryCacheDirName, eh.collectionName)
}

func (eh *emeddHistory) load() error {
	jsonData, err := os.ReadFile(eh.filename())
	if err != nil {
		eh.slog.Info("Cannot read embed history", "err", err)
		return nil
	}
	return json.Unmarshal(jsonData, &eh.history)
}

func (eh *emeddHistory) save() error {
	data, err := json.Marshal(eh.history)
	if err != nil {
		return fmt.Errorf("saving collection history file: %w", err)
	}
	return os.WriteFile(eh.filename(), data, 0660)
}
