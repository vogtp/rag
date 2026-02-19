package cfg

//go:generate stringer -type=BackendApiType --trimprefix BackendApiType
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/graze/go-throttled"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

const backendsKey = "backends"

type Backend struct {
	Name   string         `yaml:"name"`
	API    *BackendApi    `yaml:"api"`
	Models []BackendModel `yaml:"models"`
}

type BackendApiType int

func (bat BackendApiType) FromString(s string) BackendApiType {
	switch strings.ToLower(s) {
	case strings.ToLower(BackendApiTypeOllama.String()):
		return BackendApiTypeOllama
	case strings.ToLower(BackendApiTypeOpenai.String()):
		return BackendApiTypeOpenai
	default:
		return BackendApiTypeUnknown
	}
}

const (
	BackendApiTypeOllama BackendApiType = iota
	BackendApiTypeOpenai
	BackendApiTypeUnknown
)

type BackendApi struct {
	URL              string  `yaml:"URL"`
	Key              string  `yaml:"key"`
	Type             string  `yaml:"type"`
	Requests_per_sec float64 `yaml:"requests_per_sec"`
	APIType          BackendApiType
}

type BackendModel struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	API  *BackendApi
}

func GetBackends() (*Backends, error) {
	if backends != nil {
		return backends, nil
	}
	beLst := make([]Backend, 0)
	if err := viper.UnmarshalKey(backendsKey, &beLst); err != nil {
		return nil, fmt.Errorf("read %q config: %v", backendsKey, err)
	}
	var err error
	for i, be := range beLst {
		api := be.API
		api.APIType = api.APIType.FromString(be.API.Type)
		api.URL = strings.TrimRight(api.URL, "/")
		beLst[i].API = api
		if be.API.APIType == BackendApiTypeUnknown {
			err = fmt.Errorf("unknown backend api type %q: %w", be.API.Type, err)
		}
		for j := range be.Models {
			be.Models[j].API = api
		}
	}
	backends = &Backends{backends: beLst}
	return backends, err
}

var backends *Backends

type Backends struct {
	backends []Backend
}

func (bes Backends) Get(name string) *Backend {
	for _, b := range bes.backends {
		if strings.EqualFold(b.Name, name) {
			return &b
		}
	}
	return nil
}

func (bes Backends) Model(name string) *BackendModel {
	for _, be := range bes.backends {
		for _, m := range be.Models {
			if strings.EqualFold(m.Name, name) {
				return &m
			}
		}
	}
	return nil
}

func (ba BackendApi) RateLimitedHTTPClient() *http.Client {
	if ba.Requests_per_sec <= 0 {
		return http.DefaultClient
	}
	return throttled.Client(rate.NewLimiter(rate.Limit(ba.Requests_per_sec), 1))
}
