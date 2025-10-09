package vecdb

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	historyCacheDirName = ".embedd_history"
	historyCacheFilePrefix = ""
)

type emeddHistory struct {
	slog                 *slog.Logger
	collectionName       string
	history              map[string]time.Time
	minAge               time.Duration
	vecDBUpdateIntervall time.Duration
}

func (eh *emeddHistory) shouldEmbedd(d *EmbeddDocument) bool {
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
	log := eh.slog.Info
	if b {
		log = eh.slog.Debug
	}
	log("should embedd", "should", b, "age", last.String(), "last", h.String())
	return b
}

func (eh *emeddHistory) reqisterEmedded(d *EmbeddDocument) {
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
	if err:=os.MkdirAll(historyCacheDirName, os.ModePerm); err != nil {
		slog.Error("cannot create history cache dir", "err",err,"HistoryCacheDirName",historyCacheDirName)
	}
	interval := eh.vecDBUpdateIntervall
	if interval > 24*time.Hour || interval < time.Hour {
		interval = time.Hour
	}
	eh.minAge = interval
	eh.history = make(map[string]time.Time)
	return eh.load()
}

func (eh emeddHistory) key(d *EmbeddDocument) string {
	return fmt.Sprintf("%s.%s", d.IDMetaKey, d.IDMetaValue)
}

func (eh emeddHistory) filename() string {
	return fmt.Sprintf("%s/%s_%s.json",historyCacheDirName,historyCacheFilePrefix, eh.collectionName)
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
