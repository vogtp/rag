package cfg

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultVecDBUpdateIntervall = 24 * time.Hour
	MinVecDBUpdateIntervall     = time.Hour
)

type RagConfig interface {
	Collectionname() string
	Displayname() string
	Confluence() *ConfluenceCfg
	ModelEmbedding() string
	VecDBUpdateIntervall() time.Duration
}

var _ RagConfig = (*RagConfigInteral)(nil)

func (r RagConfigInteral) Displayname() string {
	return r.NameInt
}
func (r RagConfigInteral) Collectionname() string {
	return r.VecdbInt.CollectionName
}
func (r RagConfigInteral) Confluence() *ConfluenceCfg {
	return &r.ConfluenceCfg
}
func (r RagConfigInteral) ModelEmbedding() string {
	return r.ModelInt.Embedding
}

type RagConfigInteral struct {
	NameInt       string        `yaml:"Name"`
	APITokenInt   string        `yaml:"api_token"`
	ModelInt      ModelCfg      `yaml:"model"`
	VecdbInt      VecDbCfg      `yaml:"vecdb"`
	ConfluenceCfg ConfluenceCfg `yaml:"confluence"`
}

type ModelCfg struct {
	Embedding string `yaml:"embedding"`
	LLM       string `yaml:"llm"`
}

type VecDbCfg struct {
	UpdateIntervall string `yaml:"update_intervall"`
	CollectionName  string `yaml:"collection_name"`
}

type ConfluenceCfg struct {
	Key     string   `yaml:"key"`
	BaseURL string   `yaml:"baseURL"`
	Spaces  []string `yaml:"spaces"`
}

func (r RagConfigInteral) VecDBUpdateIntervall() time.Duration {
	d, err := time.ParseDuration(r.VecdbInt.UpdateIntervall)
	if err != nil {
		slog.Warn("Cannot parse update intervall of RAG", "name", r.Displayname(), "update_intervall", r.VecdbInt.UpdateIntervall, "err", err)
		return DefaultVecDBUpdateIntervall
	}
	return d
}

var defaultRagCfg *RagConfigInteral

func RagCfgFIXME() RagConfig {
	return DefaultRagCfg()
}

func DefaultRagCfg() RagConfigInteral {
	if defaultRagCfg == nil {
		defaultRagCfg = &RagConfigInteral{
			NameInt: "Default",
			ModelInt: ModelCfg{
				Embedding: viper.GetString(ModelEmbedding),
				LLM:       viper.GetString(ModelLLM),
			},
			VecdbInt: VecDbCfg{
				CollectionName:  viper.GetString(VecDBColName),
				UpdateIntervall: viper.GetString(VecDBUpdateIntervall),
			},
			ConfluenceCfg: ConfluenceCfg{
				Key:     viper.GetString(ConfluenceKey),
				BaseURL: viper.GetString(ConfluenceBaseURL),
				Spaces:  viper.GetStringSlice(ConfluenceSpaces),
			},
		}
	}
	return *defaultRagCfg
}
