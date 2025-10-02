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

type RagConfig struct {
	Name       string        `yaml:"Name"`
	APIToken   string        `yaml:"api_token"`
	Model      modelCfg      `yaml:"model"`
	Vecdb      vecDbCfg      `yaml:"vecdb"`
	Confluence confluenceCfg `yaml:"confluence"`
}

type modelCfg struct {
	Embedding string `yaml:"embedding"`
	Default   string `yaml:"default"`
}

type vecDbCfg struct {
	UpdateIntervall string `yaml:"update_intervall"`
	CollectionName  string `yaml:"collection_name"`
}

type confluenceCfg struct {
	Key     string   `yaml:"key"`
	BaseURL string   `yaml:"baseURL"`
	Spaces  []string `yaml:"spaces"`
}

func (r RagConfig) VecDBUpdateIntervall() time.Duration {
	d, err := time.ParseDuration(r.Vecdb.UpdateIntervall)
	if err != nil {
		slog.Warn("Cannot parse update intervall of RAG", "name", r.Name, "update_intervall", r.Vecdb.UpdateIntervall, "err", err)
		return DefaultVecDBUpdateIntervall
	}
	return d
}

func GetRagConfig() ([]RagConfig, error) {
	ragCfg := make([]RagConfig, 0)
	if err := viper.UnmarshalKey(ragConfigKey, &ragCfg); err != nil {
		return nil, fmt.Errorf("read rag config: %v", err)
	}
	//FIXME this is a hack since UnmarshalKey does not parse the vecDB stuff
	raw := viper.Get(ragConfigKey)
	for i, l := range raw.([]any) {
		r := l.(map[string]any)
		v := r["vecdb"].(map[string]any)
		ragCfg[i].Vecdb.CollectionName = v["collection_name"].(string)
		ragCfg[i].Vecdb.UpdateIntervall = v["update_intervall"].(string)
	}
	return ragCfg, nil
}

var defaultRagCfg = RagConfig{
	Name: "Default",
	Model: modelCfg{
		Embedding: viper.GetString(ModelEmbedding),
		Default:   viper.GetString(ModelDefault),
	},
	Vecdb: vecDbCfg{
		CollectionName:  viper.GetString(VecDBColName),
		UpdateIntervall: viper.GetString(VecDBUpdateIntervall),
	},
	Confluence: confluenceCfg{
		Key:     viper.GetString(ConfluenceKey),
		BaseURL: viper.GetString(ConfluenceBaseURL),
		Spaces:  viper.GetStringSlice(ConfluenceSpaces),
	},
}

func RagCfgFIXME() RagConfig {
	return DefaultRagCfg()
}

func DefaultRagCfg() RagConfig {
	return defaultRagCfg
}
