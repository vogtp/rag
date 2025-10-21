package test

import (
	"os"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
	"github.com/vogtp/rag/pkg/usercfg"
)

type testData struct {
	Setup struct {
		Collections []testDataCol `yaml:"Collections"`
	} `yaml:"setup"`
	Tests []testDataTest `yaml:"tests"`
}

type testDataTest struct {
	Question    string           `yaml:"Question"`
	Collections []string         `yaml:"Collections"`
	Results     []testDataResult `yaml:"Results"`
}

type testDataResult struct {
	URL      string   `yaml:"URL"`
	Title    string   `yaml:"Title"`
	Keywords []string `yaml:"Keywords"`
}

type testDataCol struct {
	Name  string `yaml:"Name"`
	Model struct {
		Embedding string `yaml:"embedding"`
		Llm       string `yaml:"llm"`
	} `yaml:"model"`
	Confluence struct {
		Key     string   `yaml:"key"`
		BaseURL string   `yaml:"baseURL"`
		Spaces  []string `yaml:"spaces"`
	} `yaml:"confluence"`
}

func (td testData) Collections() []usercfg.Collection {
	cols := make([]usercfg.Collection, len(td.Setup.Collections))
	for i, c := range td.Setup.Collections {
		cols[i] = *tdCol2DB(c)
	}
	return cols
}

func loadTestData() (*testData, error) {
	b, err := os.ReadFile(testDataFile)
	if err != nil {
		return nil, err
	}
	td := &testData{}
	err = yaml.Unmarshal(b, td)
	return td, err
}

func tdCol2DB(c testDataCol) *usercfg.Collection {
	return &usercfg.Collection{
		Displayname:    c.Name,
		Collectionname: c.Name,
		Genmodel:       c.Model.Llm,
		Embedmodel:     c.Model.Embedding,
		Source: usercfg.SourceSystem{
			URL:   c.Confluence.BaseURL,
			Key:   c.Confluence.Key,
			Parts: strings.Join(c.Confluence.Spaces, ","),
			QueryRetryMax: 1,
		},
	}
}
