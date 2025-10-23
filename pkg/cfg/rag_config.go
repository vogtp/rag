package cfg

import (
	"time"
)

const (
	DefaultVecDBUpdateIntervall = 24 * time.Hour
	MinVecDBUpdateIntervall     = time.Hour

	CntSeachResults = 7
)

type RagConfig interface {
	CollectionName() string
	DisplayName() string
	Confluence() *ConfluenceCfg
	ModelEmbedding() string
	UpdateIntervall() time.Duration
}

type ConfluenceCfg struct {
	Key     string   `yaml:"key"`
	BaseURL string   `yaml:"baseURL"`
	Spaces  []string `yaml:"spaces"`
}
