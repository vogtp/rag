package cfg

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

const (
	ragConfigKey = "rag_model"

	DefaultVecDBUpdateIntervall = 24 * time.Hour
	MinVecDBUpdateIntervall     = time.Hour
)

// type RagConfig struct {
// 	Name     string `yaml:"Name"`
// 	APIToken string `yaml:"api_token"`
// 	Model    struct {
// 		Embedding string `yaml:"embedding"`
// 		Default   string `yaml:"default"`
// 	} `yaml:"model"`
// 	Confluence struct {
// 		Key     string   `yaml:"key"`
// 		BaseURL string   `yaml:"baseURL"`
// 		Spaces  []string `yaml:"spaces"`
// 	} `yaml:"confluence"`
// 	Vecdb struct {
// 		UpdateIntervall string `yaml:"update_intervall"`
// 		CollectionName  string `yaml:"collection_name"`
// 	} `yaml:"vecdb"`
// }

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

func GetRagConfig() ([]RagConfigInteral, error) {
	ragCfg := make([]RagConfigInteral, 0)
	if err := viper.UnmarshalKey(ragConfigKey, &ragCfg); err != nil {
		return nil, fmt.Errorf("read rag config: %v", err)
	}
	//FIXME this is a hack since UnmarshalKey does not parse the vecDB stuff
	raw := viper.Get(ragConfigKey)
	for i, l := range raw.([]any) {
		r := l.(map[string]any)
		v := r["vecdb"].(map[string]any)
		ragCfg[i].VecdbInt.CollectionName = v["collection_name"].(string)
		ragCfg[i].VecdbInt.UpdateIntervall = v["update_intervall"].(string)
	}
	return ragCfg, nil
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
