package cfg

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetBackends(t *testing.T) {

	tests := []struct {
		name string
		want struct {
			api    BackendApi
			models []BackendModel
		}
	}{
		{
			name: "ollama",
			want: struct {
				api    BackendApi
				models []BackendModel
			}{
				api: BackendApi{URL: "https://ollama.example.com:11434", APIType: BackendApiTypeOllama},
				models: []BackendModel{
					{Name: "bge-m3", Type: "embedding"},
					{Name: "qwen3", Type: "llm"},
					{Name: "mistral", Type: "llm"},
					{Name: "mxbai-embed-large", Type: "embedding"},
				},
			},
		},
		{
			name: "LiteLLM",
			want: struct {
				api    BackendApi
				models []BackendModel
			}{
				api: BackendApi{URL: "https://llm.example.net/litellm", Key: "myLiteLLMKey", APIType: BackendApiTypeOpenai},
				models: []BackendModel{
					{Name: "qwen3-235b-fp8", Type: "llm"},
					{Name: "GLM-4.5V-FP8", Type: "image2txt"},
				},
			},
		},
	}

	viper.Set(CfgFile, "./testdata/backend_test.yml")
	Parse()

	backends, err := GetBackends()
	if err != nil {
		t.Fatalf("Cannot read test data config: %v", err)
	}
	if len(backends.backends) < 1 {
		t.Fatalf("No backends found: %v", backends)
	}
	be, _ := GetBackends()
	if be != backends {
		t.Errorf("Backend caching not working")
	}

	if backends.Get("notToFind") != nil {
		t.Error("Found backend that is not in there")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bes, err := GetBackends()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			be := bes.Get(tt.name)
			if be == nil {
				t.Errorf("Backend %q not found", tt.name)
				return
			}
			if be.API.URL != tt.want.api.URL {
				t.Errorf("Wrong URL: got %q want %q", be.API.URL, tt.want.api.URL)
			}
			if be.API.Key != tt.want.api.Key {
				t.Errorf("Wrong Key: got %q want %q", be.API.Key, tt.want.api.Key)
			}
			if be.API.APIType != tt.want.api.APIType {
				t.Errorf("Wrong APIType: got %q want %q", be.API.APIType, tt.want.api.APIType)
			}
			for i, w := range tt.want.models {
				if be.Models[i].Name != w.Name {
					t.Errorf("Wrong model name: got %q want %q", be.Models[i].Name, w.Name)
				}
				if be.Models[i].Type != w.Type {
					t.Errorf("Wrong model Type: got %q want %q", be.Models[i].Type, w.Type)
				}
				if backends.Model(w.Name).Name != be.Models[i].Name {
					t.Errorf("Wrong model returned by Models")
				}
			}
		})
	}
}
